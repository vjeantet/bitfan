//go:generate bitfanDoc
// Periodically scan an POP3 mailbox for new emails.
package pop3processor

import (
	"bytes"
	"sync"
	"time"

	"github.com/awillis/bitfan/processors"
	"github.com/jhillyerd/enmime"
	"github.com/vjeantet/jodaTime"
)

func New() processors.Processor {
	return &processor{opt: &options{}}
}

type options struct {
	processors.CommonOptions `mapstructure:",squash"`

	// When new mail should be retreived from POP3 server ?
	// Nothing by default, as this processor can be used in filter
	// @Type Interval
	Interval string `mapstructure:"interval" `

	// POP3 host name
	Host string `mapstructure:"host" validate:"required"`

	// POP3 server's port.
	//
	// When empty and secure is true (pop3s) the default port number is 995
	// @Default 110
	Port int `mapstructure:"port"`

	// Use TLS POP3S connexion with server.
	// The default pop3s port is 995 in this case
	// @Default false
	Secure bool `mapstructure:"secure"`

	// POP3 mailbox Username
	Username string `mapstructure:"username" validate:"required"`

	// POP3 mailbox Password
	// you may use an env variable to pass value, like password => "${BITFAN_POP3_PASSWORD}"
	Password string `mapstructure:"password" validate:"required"`

	// How long to wait for the server to respond ?
	// (in second)
	// @Default 30
	DialTimeout int `mapstructure:"dial_timeout"`

	// Should delete message after retreiving it ?
	//
	// When false, this processor will use sinceDB to not retreive an already seen message
	// @Default true
	Delete bool `mapstructure:"delete"`

	// Add Attachements, Inlines, in the produced event ?
	// When false Parts are added like
	// ```
	//  "parts": {
	//   {
	//     "Size":        336303,
	//     "Content":     $$ContentAsBytes$$,
	//     "Type":        "inline",
	//     "ContentType": "image/png",
	//     "Disposition": "inline",
	//     "FileName":    "Capture d’écran 2018-01-12 à 12.11.52.png",
	//   },
	//   {
	//     "Content":     $$ContentAsBytes$$,
	//     "Type":        "attachement",
	//     "ContentType": "application/pdf",
	//     "Disposition": "attachment",
	//     "FileName":    "58831639.pdf",
	//     "Size":        14962,
	//   },
	// },
	// ```
	// @Default false
	StripAttachments bool `mapstructure:"strip_attachments"`

	// When using a secure pop connexion (POP3S) should server'cert be verified ?
	// @Default true
	VerifyCert bool `mapstructure:"verify_cert"`

	// Path of the sincedb database file
	//
	// The sincedb database keeps track of the last seen message
	//
	// Set it to `"/dev/null"` to not persist sincedb features
	//
	// Tracks are done by host and username combination, you can customize this if needed giving a specific path
	// @Default : Host@Username
	// @ExampleLS : sincedb_path => "/dev/null"
	SincedbPath string `mapstructure:"sincedb_path"`

	// Add a field to event with the raw message data ?
	// @Default false
	AddRawMessage bool `mapstructure:"add_raw_message"`

	// Add a field to event with all headers as hash ?
	// @Default false
	AddAllHeaders bool `mapstructure:"add_all_headers"`
}

// Read mails from POP3 server
type processor struct {
	processors.Base

	opt *options
	q   chan bool
	wg  *sync.WaitGroup

	sinceDB *processors.SinceDB
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {
	p.opt = &options{
		Secure:           false,
		Interval:         "@every 30m",
		Delete:           true,
		StripAttachments: false,
		VerifyCert:       true,
		DialTimeout:      30,
		SincedbPath:      "HOST@USERNAME",
		AddRawMessage:    false,
		AddAllHeaders:    false,
	}
	err := p.ConfigureAndValidate(ctx, conf, p.opt)
	if err != nil {
		return err
	}

	if p.opt.SincedbPath == "HOST@USERNAME" {
		p.opt.SincedbPath = p.opt.Host + "@" + p.opt.Username
	}

	// Default port number depends on SSL usage
	if p.opt.Port == 0 {
		p.opt.Port = 110
		if p.opt.Secure == true {
			p.opt.Port = 995
		}
	}

	// test connexion
	client, err := newPop3Client(p.opt.Host, p.opt.Port, p.opt.Username, p.opt.Password, p.opt.Secure, p.opt.VerifyCert, p.opt.DialTimeout)
	if err != nil {
		return err
	}

	defer func() {
		client.Quit()
		client.Close()
	}()

	var count int
	var size uint64

	if count, size, err = client.Stat(); err != nil {
		return err
	}

	p.Logger.Debugf("POP3 configuration OK  - Message Count: %d, Size: %d", count, size)

	return err
}

func (p *processor) MaxConcurent() int {
	return 1
}

func (p *processor) Start(e processors.IPacket) error {
	p.wg = new(sync.WaitGroup)
	p.sinceDB = processors.NewSinceDB(&processors.SinceDBOptions{
		Identifier:    p.opt.SincedbPath,
		WriteInterval: 30,
		Storage:       p.Store,
	})
	return nil
}

func (p *processor) Tick(e processors.IPacket) error {
	return p.Receive(e)
}

func (p *processor) Receive(e processors.IPacket) error {
	client, err := newPop3Client(p.opt.Host, p.opt.Port, p.opt.Username, p.opt.Password, p.opt.Secure, p.opt.VerifyCert, p.opt.DialTimeout)
	if err != nil {
		return err
	}
	p.wg.Add(1)
	defer func() {
		p.wg.Done()
		client.Quit()
		client.Close()
	}()

	// Connect
	lastSeenUID, err := p.sinceDB.Ressource("lastSeenUID")

	err = receiveMail(client, string(lastSeenUID), func(number int, uid, data string, err error) (bool, error) {
		p.Logger.Debugf("receiving mail number=%d, uid=%s", number, uid)

		env, err := enmime.ReadEnvelope(bytes.NewBufferString(data))
		if err != nil {
			p.Logger.Errorf("reading email envelope : %s", err.Error())
			return false, err
		}

		packetFields := map[string]interface{}{
			"uid":     uid,
			"subject": env.GetHeader("Subject"),
			"text":    env.Text,
			"html":    env.HTML,
		}

		if p.opt.AddAllHeaders == true {
			packetFields["headers"] = env.Root.Header
		}

		if p.opt.AddRawMessage == true {
			packetFields["raw"] = data
		}

		//Date
		if sentAt, err := jodaTime.Parse("E, d MMM y HH:mm:ss Z", env.GetHeader("Date")); err == nil {
			packetFields["sentAt"] = sentAt
		} else {
			p.Logger.Warnf("can not parse email date %s : %s", env.GetHeader("Date"), err.Error())
		}

		// Address headers
		enmime.AddressHeaders["return-path"] = true
		for hname, _ := range enmime.AddressHeaders {
			arr := map[string]interface{}{}
			alist, _ := env.AddressList(hname)
			for _, addr := range alist {
				arr[addr.Address] = addr.Name
			}
			if len(arr) > 0 {
				packetFields[hname] = arr
			}
		}

		// Inlines / Attachements / OtherPars
		if p.opt.StripAttachments == false {
			parts := []interface{}{}
			for _, v := range env.Inlines {
				parts = append(parts, map[string]interface{}{
					"ContentType": v.ContentType,
					"Disposition": v.Disposition,
					"FileName":    v.FileName,
					"Size":        len(v.Content),
					"Content":     v.Content,
					"Type":        "inline",
				})
			}
			for _, v := range env.Attachments {
				parts = append(parts, map[string]interface{}{
					"ContentType": v.ContentType,
					"Disposition": v.Disposition,
					"FileName":    v.FileName,
					"Size":        len(v.Content),
					"Content":     v.Content,
					"Type":        "attachement",
				})
			}
			for _, v := range env.OtherParts {
				parts = append(parts, map[string]interface{}{
					"ContentType": v.ContentType,
					"Disposition": v.Disposition,
					"FileName":    v.FileName,
					"Size":        len(v.Content),
					"Content":     v.Content,
					"Type":        "other",
				})
			}
			packetFields["parts"] = parts
		}

		ne := p.NewPacket(packetFields)
		p.opt.ProcessCommonOptions(e.Fields())
		p.Send(ne)

		time.Sleep(time.Second * 4)

		p.sinceDB.SetRessource("lastSeenUID", []byte(uid))
		return p.opt.Delete, nil
	})

	return err

}

func (p *processor) Stop(e processors.IPacket) error {
	// Wait for a maybe existing retrieving process to finish
	p.wg.Wait()
	return nil
}
