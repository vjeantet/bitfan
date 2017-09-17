//go:generate bitfanDoc
// Send email when an output is received. Alternatively, you may include or exclude the email output execution using conditionals.
package email

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"

	"github.com/vjeantet/bitfan/core/location"
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

	// The fully-qualified email address to send the email to.
	//
	// This field also accepts a comma-separated string of addresses, for example: `"me@host.com, you@host.com"`
	//
	// You can also use dynamic fields from the event with the %{fieldname} syntax
	// @ExampleLS to => "me@host.com, you@host.com"
	To string `mapstructure:"to" validate:"required"`

	// The fully-qualified email address(es) to include as cc: address(es).
	//
	// This field also accepts a comma-separated string of addresses, for example: `"me@host.com, you@host.com"`
	// @ExampleLS cc => "me@host.com, you@host.com"
	Cc string `mapstructure:"cc"`

	// The fully-qualified email address(es) to include as bcc: address(es).
	//
	// This field also accepts a comma-separated string of addresses, for example: `"me@host.com, you@host.com"`
	// @ExampleLS bcc => "me@host.com, you@host.com"
	Bcc string `mapstructure:"bcc"`

	// Subject: for the email
	//
	// You can use template
	// @ExampleLS subject => "message from {{.host}}"
	Subject string `mapstructure:"subject"`

	// Path to Subject template file for the email
	Subjectfile string `mapstructure:"subjectfile"`

	// HTML Body for the email, which may contain HTML markup
	// @ExampleLS htmlBody => "<h1>Hello</h1> message received : {{.message}}"
	// @Type Location
	HTMLBody string `mapstructure:"htmlBody"`

	// Body for the email - plain text only.
	// @ExampleLS body => "message : {{.message}}. from {{.host}}."
	// @Type Location
	Body string `mapstructure:"body"`

	// Attachments - specify the name(s) and location(s) of the files
	Attachments []string `mapstructure:"attachments"`

	// Images - specify the name(s) and location(s) of the images
	Images []string `mapstructure:"images"`
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
	toStr := p.opt.To
	processors.Dynamic(&toStr, e.Fields())
	to := strings.Split(toStr, ",")

	ccStr := p.opt.Cc
	processors.Dynamic(&ccStr, e.Fields())
	cc := strings.Split(ccStr, ",")

	bccStr := p.opt.Bcc
	processors.Dynamic(&bccStr, e.Fields())
	bcc := strings.Split(bccStr, ",")

	m := gomail.NewMessage()

	fromStr := p.opt.From
	processors.Dynamic(&fromStr, e.Fields())
	m.SetHeader("From", fromStr)

	if p.opt.Replyto != "" {
		replyToStr := p.opt.Replyto
		processors.Dynamic(&replyToStr, e.Fields())
		m.SetHeader("Reply-To", replyToStr)
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
		loc, err := location.NewLocation(p.opt.Subject, p.ConfigWorkingLocation)
		if err != nil {
			p.Logger.Errorf("email subject template error : %s", err)
			return err
		}
		tmpl, _, err := loc.TemplateWithOptions(nil)
		if err != nil {
			p.Logger.Errorf("email subject template error : %s", err)
			return err
		}
		buff := bytes.NewBufferString("")
		tmpl.Execute(buff, e.Fields())
		m.SetHeader("Subject", buff.String())
	}

	if p.opt.HTMLBody != "" {
		loc, err := location.NewLocation(p.opt.HTMLBody, p.ConfigWorkingLocation)
		if err != nil {
			p.Logger.Errorf("email subject template error : %s", err)
			return err
		}
		tmpl, _, err := loc.TemplateWithOptions(nil)
		if err != nil {
			p.Logger.Errorf("email template error : %s", err)
			return err
		}

		buff := bytes.NewBufferString("")
		tmpl.Execute(buff, e.Fields())
		m.SetBody("text/html", buff.String())
	}

	if p.opt.Body != "" {
		loc, err := location.NewLocation(p.opt.Body, p.ConfigWorkingLocation)
		if err != nil {
			p.Logger.Errorf("email subject template error : %s", err)
			return err
		}
		tmpl, _, err := loc.TemplateWithOptions(nil)
		if err != nil {
			p.Logger.Errorf("email template error : %s", err)
			return err
		}
		buff := bytes.NewBufferString("")
		tmpl.Execute(buff, e.Fields())
		if p.opt.HTMLBody != "" {
			m.AddAlternative("text/plain", buff.String())
		} else {
			m.SetBody("text/plain", buff.String())
		}
	}

	for _, ref := range p.opt.Images {
		// todo build image from event data images{}
		f := ""
		if _, err := os.Stat(ref); err == nil {
			if err != nil {
				p.Logger.Errorf("Image file error %s", err)
			}
			var err error
			f, err = filepath.Abs(ref)
			if err != nil {
				p.Logger.Errorf("Image file path error %s", err)
				continue
			}
		} else if _, err := os.Stat(filepath.Join(p.ConfigWorkingLocation, ref)); err == nil {
			f = filepath.Join(p.ConfigWorkingLocation, ref)
		} else {
			p.Logger.Errorf("Image file path unknow %s", ref)
			continue
		}

		m.Embed(f)
		// m.Embed(f, gomail.Rename("o.png"))
	}

	for _, ref := range p.opt.Attachments {

		// todo build attachments from event data attachments{}
		f := ""
		if _, err := os.Stat(ref); err == nil {
			if err != nil {
				p.Logger.Errorf("Attachment file error %s", err)
			}
			var err error
			f, err = filepath.Abs(ref)
			if err != nil {
				p.Logger.Errorf("Attachment file path error %s", err)
				continue
			}
		} else if _, err := os.Stat(filepath.Join(p.ConfigWorkingLocation, ref)); err == nil {
			f = filepath.Join(p.ConfigWorkingLocation, ref)
		} else {
			p.Logger.Errorf("Attachment file path unknow %s", ref)
			continue
		}

		m.Attach(f)
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
