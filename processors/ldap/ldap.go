//go:generate bitfanDoc
// Performs a search for a specified filter on the directory and fire events with results
package ldapprocessor

import (
	"fmt"
	"strings"

	"github.com/awillis/bitfan/processors"
	ldap "gopkg.in/ldap.v2"
)

func New() processors.Processor {
	return &processor{opt: &options{}}
}

type options struct {
	processors.CommonOptions `mapstructure:",squash"`

	// ldap hostname
	// @ExampleLS host => "ldap.forumsys.com"
	Host string `mapstructure:"host" validate:"required"`

	// ldap port
	// @Default 389
	// @ExampleLS port => 389
	Port int `mapstructure:"port" validate:"required"`

	// Bind dn
	// @ExampleLS bind_dn => "cn=read-only-admin,dc=example,dc=com"
	BindDn string `mapstructure:"bind_dn"`

	// Bind password
	// @ExampleLS bind_password => "password"
	BindPassword string `mapstructure:"bind_password"`

	// Base DN
	// If bind_dn is not specified or is empty, an anonymous bind is attempted.
	// This is defined in https://tools.ietf.org/html/rfc2251#section-4.2.2
	// @ExampleLS base_dn => "dc=example,dc=com"
	BaseDn string `mapstructure:"base_dn" validate:"required"`

	// A search base (the distinguished name of the search base object) defines the
	// location in the directory from which the LDAP search begins.
	SearchBase string `mapstructure:"search_base"`

	// The search filter can be simple or advanced, using boolean operators in the format
	// described in the LDAP documentation (see [RFC4515](http://www.faqs.org/rfcs/rfc4515) for full information on filters).
	// @Default "(objectClass=*)"
	SearchFilter string `mapstructure:"search_filter" validate:"required"`

	// An array of the required attributes, e.g. ["mail", "sn", "cn"].
	//
	// Note that the "dn" is always returned irrespective of which attributes types are requested.
	//
	// Using this parameter is much more efficient than the default action (which is to return all attributes and their associated values).
	//
	// The use of this parameter should therefore be considered good practice.
	// @ExampleLS search_attributes => ["mail", "sn", "cn"]
	SearchAttributes []string `mapstructure:"search_attributes"`

	// The SCOPE setting is the starting point of an LDAP search and the depth from the
	// base DN to which the search should occur.
	//
	// There are three options (values) that can be assigned to the SCOPE parameter:
	//
	// * **base** : indicate searching only the entry at the base DN, resulting in only that entry being returned
	// * **one** : indicate searching all entries one level under the base DN - but not including the base DN and not including any entries under that one level under the base DN.
	// * **subtree** : indicate searching of all entries at all levels under and including the specified base DN
	//
	// ![scope](../ldapscope.gif)
	// @Default "subtree"
	SearchScope string `mapstructure:"search_scope"`

	// Maximum entries to return (leave empty to let the server decide)
	// @Default 0
	SizeLimit int `mapstructure:"size_limit"`

	// Desired page size in order to execute LDAP queries to fulfill the
	// search request.
	//
	// Set 0 to not use Paging
	// @Default 1000
	PagingSize int `mapstructure:"paging_size"`

	// TODO : Optional controls affect how the search is processed

	// Send an event row by row or one event with all results
	// possible values "entry", "result"
	// @Default "entry"
	EventBy string `mapstructure:"event_by"`

	// Set an interval when this processor is used as a input
	// @ExampleLS interval => "10"
	Interval string `mapstructure:"interval"`

	// You can set variable to be used in Search Query by using ${var}.
	// each reference will be replaced by the value of the variable found in search query content
	// The replacement is case-sensitive.
	// @ExampleLS var => {"hostname"=>"myhost","varname"=>"varvalue"}
	Var map[string]string `mapstructure:"var"`

	// Define the target field for placing the retrieved data. If this setting is omitted,
	// the data will be stored in the "data" field
	// Set the value to "." to store value to the root (top level) of the event
	// @ExampleLS target => "data"
	// @Default "data"
	Target string `mapstructure:"target"`
}

type processor struct {
	processors.Base
	l                *ldap.Conn
	opt              *options
	q                chan bool
	searchScopeConst int
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {
	defaults := options{
		EventBy:      "row",
		Target:       "data",
		SearchScope:  "subtree",
		Port:         389,
		SearchFilter: "(objectClass=*)",
		SizeLimit:    0,
		PagingSize:   1000,
	}

	p.opt = &defaults

	err := p.ConfigureAndValidate(ctx, conf, p.opt)
	if err != nil {
		return err
	}

	if p.opt.Interval == "" {
		p.Logger.Warningln("No interval set")
	}

	if err := p.initConn(); err != nil {
		return err
	}
	p.l.Close()

	if p.opt.SearchBase == "" {
		p.opt.SearchBase = p.opt.BaseDn
	}

	switch p.opt.SearchScope {
	case "subtree":
		p.searchScopeConst = ldap.ScopeWholeSubtree
	case "base":
		p.searchScopeConst = ldap.ScopeBaseObject
	case "one":
		p.searchScopeConst = ldap.ScopeSingleLevel
	}

	return nil
}

func (p *processor) initConn() error {
	var err error
	p.l, err = ldap.Dial("tcp", fmt.Sprintf("%s:%d", p.opt.Host, p.opt.Port))
	if err != nil {
		return err
	}

	if p.opt.BindDn != "" {
		err = p.l.Bind(p.opt.BindDn, p.opt.BindPassword)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *processor) MaxConcurent() int {
	return 1
}

func (p *processor) Tick(e processors.IPacket) error {
	return p.Receive(e)
}

func (p *processor) Receive(e processors.IPacket) error {
	p.opt.SearchFilter = strings.Trim(p.opt.SearchFilter, " ")

	p.Logger.Debugf("SearchBase - %s", p.opt.SearchBase)
	p.Logger.Debugf("searchScope - %s", p.opt.SearchScope)
	p.Logger.Debugf("SearchFilter - %s", p.opt.SearchFilter)
	p.Logger.Debugf("SearchAttributes - %s", p.opt.SearchAttributes)

	// Search for the given username
	searchRequest := ldap.NewSearchRequest(
		p.opt.SearchBase,
		p.searchScopeConst, ldap.NeverDerefAliases, p.opt.SizeLimit, 0, false,
		p.opt.SearchFilter,
		p.opt.SearchAttributes,
		nil,
	)

	var sr *ldap.SearchResult
	var err error
	p.initConn()
	defer p.l.Close()
	if p.opt.PagingSize > 0 {
		sr, err = p.l.SearchWithPaging(searchRequest, uint32(p.opt.PagingSize))
	} else {
		sr, err = p.l.Search(searchRequest)
	}

	if err != nil {
		p.Logger.Errorf("while searching.. %v", err)
	}

	var records []map[string]interface{}
	for _, entry := range sr.Entries {
		record := make(map[string]interface{})
		record["dn"] = entry.DN
		for _, attr := range entry.Attributes {
			if len(attr.Values) == 1 {
				record[attr.Name] = attr.Values[0]
			} else {
				record[attr.Name] = attr.Values
			}

		}

		if p.opt.EventBy == "row" {
			e2 := p.NewPacket(nil)
			e2.Fields().SetValueForPath(p.opt.Host, "host")
			if len(p.opt.Var) > 0 {
				e2.Fields().SetValueForPath(p.opt.Var, "var")
			}

			if p.opt.Target == "." {
				for k, v := range record {
					e2.Fields().SetValueForPath(v, k)
				}
			} else {
				e2.Fields().SetValueForPath(record, p.opt.Target)
			}

			p.opt.ProcessCommonOptions(e2.Fields())
			p.Send(e2)
		} else {
			records = append(records, record)
		}
	}

	if p.opt.EventBy != "row" {
		e.Fields().SetValueForPath(p.opt.Host, "host")
		if len(p.opt.Var) > 0 {
			e.Fields().SetValueForPath(p.opt.Var, "var")
		}
		e.Fields().SetValueForPath(records, p.opt.Target)

		p.opt.ProcessCommonOptions(e.Fields())
		p.Send(e)
	}

	return nil
}

func (p *processor) Stop(e processors.IPacket) error {
	p.l.Close()
	return nil
}
