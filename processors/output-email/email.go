//go:generate bitfanDoc
// Send email when an output is received. Alternatively, you may include or exclude the email output execution using conditionals.
package email

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	uuid "github.com/nu7hatch/gouuid"
	"github.com/vjeantet/bitfan/commons"
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
	HTMLBody string `mapstructure:"htmlbody"`

	// Body for the email - plain text only.
	// @ExampleLS body => "message : {{.message}}. from {{.host}}."
	// @Type Location
	Body string `mapstructure:"body"`

	// Attachments - specify the name(s) and location(s) of the files
	Attachments []string `mapstructure:"attachments"`

	// Use event field's values as attachment content
	// each pair is  : event field's path => attachment's name
	// @ExampleLS  attachments_with_event=>{"mydata"=>"myimage.jpg"}
	AttachEventData map[string]string `mapstructure:"attachments_with_event"`

	// Images - specify the name(s) and location(s) of the images
	Images []string `mapstructure:"images"`

	// Search for img:data in HTML body, and replace them to a reference to inline attachment
	// @Default false
	EmbedB64Images bool `mapstructure:"embed_b64_images"`
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {
	defaults := options{
		Host:           "localhost",
		From:           "bitfan@nowhere.com",
		Port:           25,
		EmbedB64Images: false,
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
		loc, err := commons.NewLocation(p.opt.Subject, p.ConfigWorkingLocation)
		if err != nil {
			p.Logger.Errorf("email subject template error : %v", err)
			return err
		}
		tmpl, _, err := loc.TemplateWithOptions(nil)
		if err != nil {
			p.Logger.Errorf("email subject template error : %v", err)
			return err
		}
		buff := bytes.NewBufferString("")
		tmpl.Execute(buff, e.Fields())
		m.SetHeader("Subject", buff.String())
	}

	if p.opt.HTMLBody != "" {
		loc, err := commons.NewLocation(p.opt.HTMLBody, p.ConfigWorkingLocation)
		if err != nil {
			p.Logger.Errorf("email subject template error : %v", err)
			return err
		}
		tmpl, _, err := loc.TemplateWithOptions(nil)
		if err != nil {
			p.Logger.Errorf("email template error : %v", err)
			return err
		}

		buff := bytes.NewBufferString("")
		tmpl.Execute(buff, e.Fields())

		if p.opt.EmbedB64Images == false {
			m.SetBody("text/html", buff.String())
		} else {
			content := buff.String()
			r, _ := regexp.Compile(`<img([^>]*)src=['"]data:image/(\w+);base64,([^"']*)['"]([^>]*)>`)

			for i, match := range r.FindAllStringSubmatch(content, -1) {
				imgTag := match[0]
				mBeforeAttr := match[1]
				mImageType := match[2]
				m64Data := match[3]
				mAfterAttr := match[4]

				imgUid := fmt.Sprintf("embed-%d."+mImageType, i)
				content = strings.Replace(content, imgTag, `<img`+mBeforeAttr+`src='cid:`+imgUid+`'`+mAfterAttr+`>`, 1)
				imgPath := filepath.Join(os.TempDir(), imgUid)

				sDec, err := base64.StdEncoding.DecodeString(m64Data)
				if err != nil {
					p.Logger.Errorf("error while decoding base64 %s", err.Error())
					continue
				}

				if err := ioutil.WriteFile(imgPath, sDec, 0644); err != nil {
					p.Logger.Errorf("error while writing image to %s", imgPath)
					continue
				}

				p.opt.Images = append(p.opt.Images, imgPath)
			}
			m.SetBody("text/html", content)
		}

	}

	if p.opt.Body != "" {
		loc, err := commons.NewLocation(p.opt.Body, p.ConfigWorkingLocation)
		if err != nil {
			p.Logger.Errorf("email subject template error : %v", err)
			return err
		}
		tmpl, _, err := loc.TemplateWithOptions(nil)
		if err != nil {
			p.Logger.Errorf("email template error : %v", err)
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
				p.Logger.Errorf("Image file error %v", err)
			}
			var err error
			f, err = filepath.Abs(ref)
			if err != nil {
				p.Logger.Errorf("Image file path error %v", err)
				continue
			}
		} else if _, err := os.Stat(filepath.Join(p.ConfigWorkingLocation, ref)); err == nil {
			f = filepath.Join(p.ConfigWorkingLocation, ref)
		} else {
			p.Logger.Errorf("Image file path unknow %v", ref)
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
				p.Logger.Errorf("Attachment file error %v", err)
			}
			var err error
			f, err = filepath.Abs(ref)
			if err != nil {
				p.Logger.Errorf("Attachment file path error %v", err)
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

	for path, attachmentName := range p.opt.AttachEventData {
		value, err := e.Fields().ValueForPath(path)
		if err != nil {
			p.Logger.Warningf("Attach event data failed for path %s : %s", path, err)
			continue
		}

		name, _ := uuid.NewV4()
		attachFileName := filepath.Join(os.TempDir(), "bitfan-tmp-email-"+name.String())
		switch vt := value.(type) {
		case string:
			sDec, err := base64.StdEncoding.DecodeString(vt)
			if err != nil {
				p.Logger.Debugf("base64 decode %s, %s", path, err)
				if err = ioutil.WriteFile(attachFileName, []byte(vt), 0644); err != nil {
					p.Logger.Warningf("write file error %s : %s", attachFileName, err)
					continue
				}
			}
			if err = ioutil.WriteFile(attachFileName, sDec, 0644); err != nil {
				p.Logger.Warningf("write file error %s : %s", attachFileName, err)
				continue
			}
		case []byte:
			if err = ioutil.WriteFile(attachFileName, vt, 0644); err != nil {
				p.Logger.Warningf("write file error %s : %s", attachFileName, err)
				continue
			}
		default:
			p.Logger.Warningf("Unknow content type for %V", value)
			continue
		}

		m.Attach(attachFileName, gomail.Rename(attachmentName))
	}

	d := gomail.NewDialer(p.opt.Host, p.opt.Port, p.opt.Username, p.opt.Password)
	if err := d.DialAndSend(m); err != nil {
		p.Logger.Errorf("email send error : %v", err)
		return err
	}

	return nil
}

func (p *processor) Stop(e processors.IPacket) error {
	//TODO disconnect from SMTP
	return nil
}
