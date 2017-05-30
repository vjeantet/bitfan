//go:generate bitfanDoc
package inputsql

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/ShowMax/go-fqdn"
	_ "github.com/go-sql-driver/mysql"
	"github.com/vjeantet/bitfan/core/location"
	"github.com/vjeantet/bitfan/processors"
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
	Codec string

	// GOLANG driver class to load, for example, "mysql".
	// @ExampleLS driver => "mysql"
	Driver string `mapstructure:"driver" validate:"required"`

	// Send an event row by row or one event with all results
	// possible values "row", "result"
	// @Default "row"
	EventBy string `mapstructure:"event_by"`

	// SQL Statement
	// When there is more than 1 statement, only data from the last one will generate events.
	// @ExampleLS statement => "SELECT * FROM mytable"
	Statement string `mapstructure:"statement" validate:"required"`

	// Set an interval when this processor is used as a input
	// @ExampleLS interval => "10"
	Interval string `mapstructure:"interval"  validate:"required"`

	// @ExampleLS connection_string => "username:password@tcp(192.168.1.2:3306)/mydatabase?charset=utf8"
	ConnectionString string `mapstructure:"connection_string" validate:"required"`

	// You can set variable to be used in Statements by using ${var}.
	// each reference will be replaced by the value of the variable found in Statement's content
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
	db   *sql.DB
	opt  *options
	q    chan bool
	host string
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {
	defaults := options{
		EventBy: "row",
		Target:  "data",
	}

	p.opt = &defaults
	p.host = fqdn.Get()

	err := p.ConfigureAndValidate(ctx, conf, p.opt)
	if err != nil {
		return err
	}

	if p.opt.Interval == "" {
		p.Logger.Warningln("No interval set")
	}

	loc, err := location.NewLocation(p.opt.Statement, p.ConfigWorkingLocation)
	if err != nil {
		return err
	}

	content, _, err := loc.ContentWithOptions(p.opt.Var)
	if err != nil {
		return err
	}
	p.opt.Statement = string(content)

	p.db, err = sql.Open(p.opt.Driver, p.opt.ConnectionString)
	if err != nil {
		return err
	}

	err = p.db.Ping()
	if err != nil {
		return err
	}

	return nil
}

func (p *processor) Tick(e processors.IPacket) error {
	return p.Receive(e)
}

func (p *processor) Receive(e processors.IPacket) error {
	p.opt.Statement = strings.Trim(p.opt.Statement, ";")
	reqs := strings.Split(p.opt.Statement, ";")
	for _, r := range reqs[:len(reqs)-1] {
		p.Logger.Debugf("db.Exec - %s", r)
		p.db.Exec(r)
	}

	p.Logger.Debugf("db.Query - %s", reqs[len(reqs)-1])
	rows, err := p.db.Query(reqs[len(reqs)-1])
	if err != nil {
		return err
	}

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	var records []map[string]interface{}
	for rows.Next() {
		record := make(map[string]interface{})
		err = rows.Scan(scanArgs...)
		if err != nil {
			return err
		}

		for i, col := range values {
			if col != nil {
				// fmt.Printf("\n%s: type= %s\n", columns[i], reflect.TypeOf(col))
				switch t := col.(type) {
				default:
					fmt.Printf("Unexpected type %T\n", t)
				case bool:
					record[columns[i]] = col.(bool)
				case int:
					record[columns[i]] = col.(int)
				case int64:
					record[columns[i]] = col.(int64)
				case float64:
					record[columns[i]] = col.(float64)
				case string:
					record[columns[i]] = col.(string)
				case []byte: // -- all cases go HERE!
					record[columns[i]] = string(col.([]byte))
					//case time.Time:
					// record[columns[i]] = col.(string)
				}
			}
		}

		if p.opt.EventBy == "row" {
			var e2 processors.IPacket
			e2 = e.Clone()
			e2.Fields().SetValueForPath(p.host, "host")
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

	rows.Close()

	if p.opt.EventBy != "row" {
		e.Fields().SetValueForPath(p.host, "host")
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
	p.db.Close()
	return nil
}
