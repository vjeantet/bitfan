package geoip

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/vjeantet/bitfan/processors/testutils"
)

func TestCommonProcessorDetails(t *testing.T) {
	Convey("This processor... ", t, func() {
		p, _ := testutils.NewProcessor(New)

		Convey("is a valid bitfan processor", func() {
			_, ok := p.Processor.(*processor)
			So(ok, ShouldBeTrue)
		})

		Convey("does not have limit on concurent event processing", func() {
			So(p.MaxConcurent(), ShouldEqual, 0)
		})

		Convey("is self documented", func() {
			if p.Doc().Doc == "" {
				Println("Missing documentation for this processor")
			} else {
				So(true, ShouldBeTrue)
			}
		})
	})
}

func TestInvalidConfiguration(t *testing.T) {
	conf := map[string]interface{}{}

	Convey("When source is missing", t, func() {
		conf["database"] = "foo/bar"
		_, err := testutils.NewProcessor(New, conf)
		Convey("Then an error occurs", func() {
			So(err, ShouldBeError)
		})
	})
}

func TestNormalCases(t *testing.T) {

	Convey("Given an existing event with a valid ip", t, func() {
		event := testutils.NewPacket("", map[string]interface{}{})
		conf := map[string]interface{}{
			"database": setupTmpDatabase("GeoIP2-City-Test.mmdb"),
			"source":   "ip",
		}
		event.Fields().SetValueForPath(`89.160.20.113`, "ip")

		Convey("When using a city database", func() {
			Convey("Then 1 event is produced with city and country name : ", func() {
				p, _ := testutils.StartNewProcessor(New, conf)
				defer p.Stop(nil)
				p.Receive(event)

				So(p.SentPacketsCount(0), ShouldEqual, 1)
				pe := p.SentPackets(0)[0]
				So(pe.Fields().ValueOrEmptyForPathString("city_name"), ShouldEqual, "Linköping")
				So(pe.Fields().ValueOrEmptyForPathString("country_name"), ShouldEqual, "Sweden")
			})

			Convey("When target field is set ", func() {
				conf["target"] = "geo"

				Convey("Then geo data are valued in this target field ", func() {
					p, _ := testutils.StartNewProcessor(New, conf)
					defer p.Stop(nil)
					p.Receive(event)

					So(p.SentPacketsCount(0), ShouldEqual, 1)
					pe := p.SentPackets(0)[0]
					So(pe.Fields().Exists("city_name"), ShouldBeFalse)
					So(pe.Fields().ValueOrEmptyForPathString("geo.city_name"), ShouldEqual, "Linköping")
					So(pe.Fields().ValueOrEmptyForPathString("geo.country_name"), ShouldEqual, "Sweden")
				})
			})

		})

		Convey("When using a ISP database", func() {
			conf := map[string]interface{}{
				"database":      setupTmpDatabase("GeoIP2-ISP-Test.mmdb"),
				"source":        "ip",
				"database_type": "isp",
			}

			event.Fields().SetValueForPath(`84.128.23.4`, "ip")

			Convey("Then 1 event is produced with ISP name and organization : ", func() {
				p, _ := testutils.StartNewProcessor(New, conf)
				defer p.Stop(nil)
				p.Receive(event)

				So(p.SentPacketsCount(0), ShouldEqual, 1)
				pe := p.SentPackets(0)[0]

				So(pe.Fields().ValueOrEmptyForPathString("isp"), ShouldEqual, "Deutsche Telekom AG")
				So(pe.Fields().ValueOrEmptyForPathString("organization"), ShouldEqual, "Deutsche Telekom AG")
			})
		})

		Convey("When using a Country database", func() {
			conf := map[string]interface{}{
				"database":      setupTmpDatabase("GeoIP2-Country-Test.mmdb"),
				"source":        "ip",
				"database_type": "country",
			}

			event.Fields().SetValueForPath(`111.235.160.9`, "ip")

			Convey("Then 1 event is produced with country name", func() {
				p, _ := testutils.StartNewProcessor(New, conf)
				defer p.Stop(nil)
				p.Receive(event)

				So(p.SentPacketsCount(0), ShouldEqual, 1)

				pe := p.SentPackets(0)[0]
				So(pe.Fields().ValueOrEmptyForPathString("country_code"), ShouldEqual, "CN")
				So(pe.Fields().ValueOrEmptyForPathString("country_name"), ShouldEqual, "People's Republic of China")
				So(pe.Fields().ValueOrEmptyForPathString("continent_code"), ShouldEqual, "AS")
				So(pe.Fields().ValueOrEmptyForPathString("continent_name"), ShouldEqual, "Asia")
			})
		})

	})
}

func TestFallbackCases(t *testing.T) {
	// When the ip field value is empty
	// When the ip field does not exist
	// when ip field is valid
	//   When no geoIP

}

func setupTmpDatabase(geoFileName string) string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}

	srcFilePath := path.Dir(filename) + string(os.PathSeparator) + "testdata" + string(os.PathSeparator) + geoFileName
	destFilePath := os.TempDir() + string(os.PathSeparator) + geoFileName

	srcFile, err := os.Open(srcFilePath)
	if err != nil {
		fmt.Println(err)
	}
	defer srcFile.Close()

	destFile, err := os.Create(destFilePath) // creates if file doesn't exist
	if err != nil {
		fmt.Println(err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile) // check first var for number of bytes copied
	if err != nil {
		fmt.Println(err)
	}
	err = destFile.Sync()
	if err != nil {
		fmt.Println(err)
	}

	return destFilePath
}
