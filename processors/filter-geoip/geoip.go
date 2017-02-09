//go:generate bitfanDoc
package geoip

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/hraban/lrucache"
	"github.com/oschwald/geoip2-golang"
	"github.com/vjeantet/bitfan/processors"
)

func New() processors.Processor {
	return &processor{opt: &options{}}
}

// Adds geographical information about an IP address
type processor struct {
	processors.Base

	opt      *options
	cache    *lrucache.Cache
	etag     string
	database *geoip2.Reader
}

type options struct {
	// If this filter is successful, add any arbitrary fields to this event.
	// Field names can be dynamic and include parts of the event using the %{field}.
	AddField map[string]interface{} `mapstructure:"add_field"`

	// If this filter is successful, add arbitrary tags to the event.
	// Tags can be dynamic and include parts of the event using the %{field} syntax.
	AddTag []string `mapstructure:"add_tag"`

	// Path or URL to the MaxMind GeoIP2 database.
	// Default value is "http://geolite.maxmind.com/download/geoip/database/GeoLite2-City.mmdb.gz"
	// Note that URL can point to gzipped database (*.mmdb.gz) but local path must point to an unzipped file.
	Database string `mapstructure:"database"`

	// Type of GeoIP database. Default value is "city"
	// Must be one of "city", "asn", "isp" or "organization".
	DatabaseType string `mapstructure:"database_type"`

	// GeoIP database refresh interval, in minutes. Default value is 60
	// If `database` field is a path, file will be reloaded from disk.
	// If it is an URL, database will be fetched (if ETAG differs) and reloaded.
	RefreshInterval time.Duration `mapstructure:"refresh_interval"`

	// An array of geoip fields to be included in the event.
	// Possible fields depend on the database type. By default, all geoip fields are included in the event.
	Fields []string `mapstructure:"fields"`

	// Cache size. Default value is 1000
	LruCacheSize int64 `mapstructure:"lru_cache_size"`

	// If this filter is successful, remove arbitrary fields from this event.
	RemoveField []string `mapstructure:"remove_field"`

	// If this filter is successful, remove arbitrary tags from the event.
	// Tags can be dynamic and include parts of the event using the %{field} syntax
	RemoveTag []string `mapstructure:"remove_tag"`

	// The field containing the IP address or hostname to map via geoip.
	Source string `mapstructure:"source" validate:"required"`

	// Append values to the tags field when there has been no successful match
	// Default value is ["_geoipparsefailure"]
	TagOnFailure []string `mapstructure:"tag_on_failure"`

	// Define the target field for placing the parsed data. If this setting is omitted,
	// the geoip data will be stored at the root (top level) of the event
	Target string `mapstructure:"target"`

	// Language to use for city/region/continent names
	Language string `mapstructure:"language"`
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) (err error) {
	defaults := options{
		Fields: []string{
			"city_name",
			"country_code",
			"country_name",
			"continent_code",
			"continent_name",
			"latitude",
			"longitude",
			"timezone",
			"postal_code",
			"region_code",
			"region_name",
			"is_anonymous_proxy",
			"is_satellite_provider",
			"asn",
			"autonomous_system_organization",
			"organization",
			"isp",
		},
		Language:        "en",
		LruCacheSize:    1000,
		Target:          "",
		Database:        "http://geolite.maxmind.com/download/geoip/database/GeoLite2-City.mmdb.gz",
		DatabaseType:    "city",
		RefreshInterval: 60,
		TagOnFailure:    []string{"_geoipparsefailure"},
	}
	p.opt = &defaults
	p.cache = lrucache.New(p.opt.LruCacheSize)

	err = p.ConfigureAndValidate(ctx, conf, p.opt)
	if err != nil {
		return err
	}

	err = p.refresh()
	ticker := time.NewTicker(p.opt.RefreshInterval * time.Minute)
	go func() {
		for range ticker.C {
			err := p.refresh()
			if err != nil {
				p.Logger.Error(err)
			}
		}
	}()

	return err
}

func (p *processor) Receive(e processors.IPacket) error {
	ip, err := e.Fields().ValueForPathString(p.opt.Source)
	if err != nil {
		processors.AddTags(p.opt.TagOnFailure, e.Fields())
		p.Send(e, 0)
		return nil
	}

	p.cache.OnMiss(p.getInfo())

	cache, err := p.cache.Get(ip)
	if err != nil {
		processors.AddTags(p.opt.TagOnFailure, e.Fields())
		p.Send(e, 0)
		return nil
	}

	data := make(map[string]interface{})

	switch cache.(type) {
	case *geoip2.City:
		city := cache.(*geoip2.City)
		for _, field := range p.opt.Fields {
			switch field {
			case "city_name":
				data["city_name"] = city.City.Names[p.opt.Language]
			case "country_code":
				data["country_code"] = city.Country.IsoCode
			case "country_name":
				data["country_name"] = city.Country.Names[p.opt.Language]
			case "continent_code":
				data["continent_code"] = city.Continent.Code
			case "continent_name":
				data["continent_name"] = city.Continent.Names[p.opt.Language]
			case "latitude":
				data["latitude"] = city.Location.Latitude
			case "longitude":
				data["longitude"] = city.Location.Longitude
			case "metro_code":
				data["metro_code"] = city.Location.MetroCode
			case "timezone":
				data["timezone"] = city.Location.TimeZone
			case "postal_code":
				data["postal_code"] = city.Postal.Code
			case "region_code":
				subdivisions := city.Subdivisions
				if len(subdivisions) > 0 {
					data["region_code"] = subdivisions[0].IsoCode
				}
			case "region_name":
				subdivisions := city.Subdivisions
				if len(subdivisions) > 0 {
					data["region_name"] = subdivisions[0].Names[p.opt.Language]
				}
			case "is_anonymous_proxy":
				data["is_anonymous_proxy"] = city.Traits.IsAnonymousProxy
			case "is_satellite_provider":
				data["is_satellite_provider"] = city.Traits.IsSatelliteProvider
			}
		}
	case *geoip2.ISP:
		isp := cache.(*geoip2.ISP)
		for _, field := range p.opt.Fields {
			switch field {
			case "asn":
				data["asn"] = isp.AutonomousSystemNumber
			case "autonomous_system_organization":
				data["autonomous_system_organization"] = isp.AutonomousSystemOrganization
			case "organization":
				data["organization"] = isp.Organization
			case "isp":
				data["isp"] = isp.ISP
			}
		}
	}

	if p.opt.Target != "" {
		e.Fields().SetValueForPath(data, p.opt.Target)
	} else {
		for k, v := range data {
			e.Fields().SetValueForPath(v, k)
		}
	}

	processors.ProcessCommonFields2(e.Fields(),
		p.opt.AddField,
		p.opt.AddTag,
		p.opt.RemoveField,
		p.opt.RemoveTag,
	)

	p.Send(e, 0)
	return nil
}

func (p *processor) getInfo() func(ip string) (lrucache.Cacheable, error) {
	return func(ip string) (record lrucache.Cacheable, err error) {
		netIP := net.ParseIP(ip)
		if netIP == nil {
			return nil, errors.New("no valid IP address found")
		}

		switch strings.ToLower(p.opt.DatabaseType) {
		case "city":
			record, err = p.database.City(netIP)
		case "isp":
			record, err = p.database.ISP(netIP)
		case "country":
			record, err = p.database.Country(netIP)
		case "domain":
			record, err = p.database.Domain(netIP)
		case "anonymousip":
			record, err = p.database.AnonymousIP(netIP)
		default:
			return nil, fmt.Errorf("Unknown database type: %s", p.opt.DatabaseType)
		}
		return record, err
	}
}

func (p *processor) refresh() error {
	// Parse database field to check either http url or local path
	url, err := url.Parse(p.opt.Database)

	if url.Scheme == "http" || url.Scheme == "https" {
		p.Logger.Infof("Geoip filter: downloading %s database from %s\n", p.opt.DatabaseType, p.opt.Database)

		// Creating http request with Etag check
		client := &http.Client{}
		request, err := http.NewRequest("GET", p.opt.Database, nil)
		if err != nil {
			return err
		}
		request.Header.Set("If-None-Match", p.etag)
		response, err := client.Do(request)
		if err != nil {
			return err
		}
		defer response.Body.Close()

		if response.StatusCode == http.StatusNotModified {
			// Database is up to date, no need to re-download
			p.Logger.Infof("Geoip filter: database is up to date, no refresh needed.")
			return nil
		}

		if response.StatusCode != http.StatusOK {
			return fmt.Errorf("bad http response: %v", response.Status)
		}

		// Save Etag for next refresh
		p.etag = response.Header["Etag"][0]

		var stream io.Reader
		// Check if URL is pointing to a gzipped file (ending in .gz)
		if path.Ext(url.Path) == ".gz" {
			stream, err = gzip.NewReader(response.Body)
			if err != nil {
				return err
			}
		} else {
			// Remote file is not gzipped, parse the stream directly
			stream = response.Body
		}

		db, err := ioutil.ReadAll(stream)
		if err != nil {
			return err
		}

		p.database, err = geoip2.FromBytes(db)
		p.Logger.Infof("Geoip filter: %s database successfuly downloaded.\n", p.opt.DatabaseType)
		return err
	}

	// database field is not an http URL, trying to open file from local path
	p.database, err = geoip2.Open(p.opt.Database)
	return err
}

func (p *processor) Tick(e processors.IPacket) error { return nil }
func (p *processor) Stop(e processors.IPacket) error { return nil }
