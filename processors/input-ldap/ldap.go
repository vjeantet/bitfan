//go:generate bitfanDoc
// Performs a search for a specified filter on the directory and fire events with results
package ldapinput

import (
	"fmt"
	"log"
	"strings"

	"github.com/vjeantet/bitfan/processors"
	ldap "gopkg.in/ldap.v2"
)

func New() processors.Processor {
	return &processor{opt: &options{}}
}

type options struct {
	// If this filter is successful, add any arbitrary fields to this event.
	Add_field map[string]interface{}

	// If this filter is successful, add arbitrary tags to the event. Tags can be dynamic
	// and include parts of the event using the %{field} syntax.
	Tags []string

	// Add a type field to all events handled by this input
	Type string

	// The codec used for input data. Input codecs are a convenient method for decoding
	// your data before it enters the input, without needing a separate filter in your bitfan pipeline
	// Codec string

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

	// TODO : Optional controls affect how the search is processed

	// Send an event row by row or one event with all results
	// possible values "entry", "result"
	// @Default "entry"
	EventBy string `mapstructure:"event_by"`

	// Set an interval when this processor is used as a input
	// @ExampleLS interval => "10"
	Interval string `mapstructure:"interval"  validate:"required"`

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
	}

	p.opt = &defaults

	err := p.ConfigureAndValidate(ctx, conf, p.opt)
	if err != nil {
		return err
	}

	if p.opt.Interval == "" {
		p.Logger.Warningln("No interval set")
	}

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
		p.searchScopeConst, ldap.NeverDerefAliases, 0, 0, false,
		p.opt.SearchFilter,
		p.opt.SearchAttributes,
		nil,
	)

	sr, err := p.l.Search(searchRequest)
	if err != nil {
		log.Fatal(err)
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
			e2 := e.Clone()
			e2 = p.NewPacket("", map[string]interface{}{})
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

			processors.ProcessCommonFields(e2.Fields(), p.opt.Add_field, p.opt.Tags, p.opt.Type)
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

		processors.ProcessCommonFields(e.Fields(), p.opt.Add_field, p.opt.Tags, p.opt.Type)
		p.Send(e)
	}

	return nil
}

func (p *processor) Stop(e processors.IPacket) error {
	p.l.Close()
	return nil
}
