//go:generate bitfanDoc
// Send email when an output is received. Alternatively, you may include or exclude the email output execution using conditionals.
package email

import (
	"bytes"
	"html/template"
	"path/filepath"
	"strings"

	"github.com/vjeantet/bitfan/processors"
	gomail "gopkg.in/gomail.v2"
)

func New() processors.Processor {
	return &processor{opt: &options{}}
}

// Sends email to a specified address when output is received
type processor struct {
	processors.Base

	opt *options
}

type options struct {
	// The address used to connect to the mail server
	// @Default "localhost"
	Host string `mapstructure:"address"`

	// Port used to communicate with the mail server
	// @Default 25
	Port int `mapstructure:"port"`

	// Username to authenticate with the server
	Username string `mapstructure:"username"`

	// Password to authenticate with the server
	Password string `mapstructure:"password"`

	// The fully-qualified email address for the From: field in the email
	// @Default "bitfan@nowhere.com"
	From string `mapstructure:"from"`

	// The fully qualified email address for the Reply-To: field
	// @ExampleLS replyto => "test@nowhere.com"
	Replyto string `mapstructure:"replyto"`

	//The fully-qualified email address to send the email to.
	// This field also accepts a comma-separated string of addresses, for example: `"me@host.com, you@host.com"`
	// You can also use dynamic fields from the event with the %{fieldname} syntax
	// @ExampleLS to => "me@host.com, you@host.com"
	To string `mapstructure:"to", validate:"required"`

	// The fully-qualified email address(es) to include as cc: address(es).
	// This field also accepts a comma-separated string of addresses, for example: `"me@host.com, you@host.com"`
	// @ExampleLS cc => "me@host.com, you@host.com"
	Cc string `mapstructure:"cc"`

	// The fully-qualified email address(es) to include as bcc: address(es).
	// This field also accepts a comma-separated string of addresses, for example: `"me@host.com, you@host.com"`
	// @ExampleLS bcc => "me@host.com, you@host.com"
	Bcc string `mapstructure:"bcc"`

	// Subject: for the email
	// You can use template
	// @ExampleLS subject => "message from {{.host}}"
	Subject string `mapstructure:"subject"`

	// Path to Subject template file for the email
	Subjectfile string `mapstructure:"subjectfile"`

	// HTML Body for the email, which may contain HTML markup
	// @ExampleLS htmlbody => "<h1>Hello</h1> message received : {{.message}}"
	Htmlbody string `mapstructure:"htmlbody"`

	// Local path to HTML Body template file for the email, which may contain HTML markup
	// can be relative to the configuration file
	Htmlbodyfile string `mapstructure:"htmlbodyfile"`

	// Body for the email - plain text only.
	// @ExampleLS body => "message : {{.message}}. from {{.host}}."
	Body string `mapstructure:"body"`

	// Path to Body template file for the email.
	Bodyfile string `mapstructure:"bodyfile"`

	// Attachments - specify the name(s) and location(s) of the files
	Attachments []string `mapstructure:"attachments"`
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {
	defaults := options{
		Host: "localhost",
		From: "bitfan@nowhere.com",
		Port: 25,
	}
	p.opt = &defaults
	err := p.ConfigureAndValidate(ctx, conf, p.opt)
	if err != nil {
		return err
	}

	if p.opt.Htmlbodyfile != "" {
		if !filepath.IsAbs(p.opt.Htmlbodyfile) {
			// Set HtmlBodyFile full path (assuming a relative local path was given)
			p.opt.Htmlbodyfile = filepath.Join(p.ConfigWorkingLocation, p.opt.Htmlbodyfile)
		}
	}

	return nil
}

func (p *processor) Start(e processors.IPacket) error {
	// TODO prepare TEMPLATE HTML et SIMPLE
	// connect
	return nil
}

func (p *processor) Receive(e processors.IPacket) error {
	// TODO use prepared template
	// connnect only if needed

	to := strings.Split(p.opt.To, ",")
	cc := strings.Split(p.opt.Cc, ",")
	bcc := strings.Split(p.opt.Bcc, ",")

	m := gomail.NewMessage()

	m.SetHeader("From", p.opt.From)

	if p.opt.Replyto != "" {
		m.SetHeader("Reply-To", p.opt.Replyto)
	}

	m.SetHeader("To", to...)

	if p.opt.Cc != "" {
		m.SetHeader("Cc", cc...)
	}
	if p.opt.Bcc != "" {
		m.SetHeader("Bcc", bcc...)
	}
	// m.SetAddressHeader("Cc", p.opt.Cc, "")

	// // Subject with Dynamics subject => "message from %{host}!"
	// subject := p.opt.Subject
	// processors.Dynamic(&subject, e.Fields())
	// m.SetHeader("Subject", subject)
	// pp.Println("subject-->", subject)

	if p.opt.Subject != "" {
		tmpl, err := template.New("email").Parse(p.opt.Subject)
		if err != nil {
			p.Logger.Errorf("email subject template error : %s", err)
			return err
		}
		buff := bytes.NewBufferString("")
		tmpl.Execute(buff, e.Fields())
		m.SetHeader("Subject", buff.String())
	}

	if p.opt.Htmlbody != "" {
		tmpl, err := template.New("email").Parse(p.opt.Htmlbody)
		if err != nil {
			p.Logger.Errorf("email template error : %s", err)
			return err
		}
		buff := bytes.NewBufferString("")
		tmpl.Execute(buff, e.Fields())
		m.SetBody("text/html", buff.String())
	}

	if p.opt.Htmlbodyfile != "" {
		tmpl, err := template.ParseGlob(p.opt.Htmlbodyfile)
		if err != nil {
			p.Logger.Errorf("email template error : %s", err)
			return err
		}

		buff := bytes.NewBufferString("")

		tmpl.Execute(buff, e.Fields())
		m.SetBody("text/html", buff.String())
	}

	if p.opt.Body != "" {
		tmpl, err := template.New("email").Parse(p.opt.Body)
		if err != nil {
			p.Logger.Errorf("email template error : %s", err)
			return err
		}
		buff := bytes.NewBufferString("")
		tmpl.Execute(buff, e.Fields())
		if p.opt.Htmlbody != "" || p.opt.Htmlbodyfile != "" {
			m.AddAlternative("text/plain", buff.String())
		} else {
			m.SetBody("text/plain", buff.String())
		}
	}

	for _, f := range p.opt.Attachments {
		// todo check if attachments exists
		// todo build attachments from event data attachments{}
		m.Attach(f)
		// todo handle error
	}

	d := gomail.NewDialer(p.opt.Host, p.opt.Port, p.opt.Username, p.opt.Password)
	if err := d.DialAndSend(m); err != nil {
		p.Logger.Errorf("email send error : %s", err)
		return err
	}

	return nil
}

func (p *processor) Stop(e processors.IPacket) error {
	//TODO disconnect from SMTP
	return nil
}
