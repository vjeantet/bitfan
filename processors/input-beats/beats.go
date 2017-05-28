//go:generate bitfanDoc
package beatsinput

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	"time"

	"github.com/elastic/go-lumber/log"
	"github.com/elastic/go-lumber/server/v2"
	"github.com/vjeantet/bitfan/processors"
)

func New() processors.Processor {
	return &processor{opt: &options{}}
}

type processor struct {
	processors.Base

	server *v2.Server
	opt    *options
	q      chan bool
}

type options struct {
	Add_field map[string]interface{}
	Codec     string

	// The number of seconds before we raise a timeout,
	// this option is useful to control how much time to wait if something is blocking
	// the pipeline
	Congestion_threshold int

	// The IP address to listen on
	Host string

	// The port to listen on (default 5044)
	Port int

	// Events are by default send in plain text,
	// you can enable encryption by using ssl to true and
	// configuring the ssl_certificate and ssl_key options
	Ssl bool

	// SSL certificate to use (path)
	Ssl_certificate string

	// Validate client certificates against theses authorities
	//  You can defined multiples files or path, all the certificates will be read
	//  and added to the trust store.
	//  You need to configure the ssl_verify_mode to peer or force_peer to enable
	//  the verification.
	// This feature only support certificate directly signed by your root ca.
	// Intermediate CA are currently not supported.
	Ssl_certificate_authorities []string

	// SSL key to use (path)
	Ssl_key string

	// SSL key passphrase to use. (not yet implemented)
	Ssl_key_passphrase string

	// By default the server dont do any client verification,
	// peer will make the server ask the client to provide a certificate,
	//   if the client provide the certificate it will be validated.
	// force_peer will make the server ask the client for their certificate,
	//   if the clients doesn’t provide it the connection will be closed.
	// This option need to be used with ssl_certificate_authorities and a defined list of CA.
	// Value can be any of: none, peer, force_peer
	Ssl_verify_mode string

	// Add any number of arbitrary tags to your event
	Tags []string

	// Add a type field to all events handled by this input
	Type string
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {
	p.opt.Congestion_threshold = 30
	p.opt.Host = "0.0.0.0"
	p.opt.Port = 5044
	p.opt.Ssl = false
	p.opt.Ssl_verify_mode = "none"

	return p.ConfigureAndValidate(ctx, conf, p.opt)
}

func (p *processor) Start(e processors.IPacket) error {
	p.q = make(chan bool)

	options := []v2.Option{}

	if p.opt.Ssl == true {
		config := &tls.Config{}

		// Server Certificates
		cert, err := tls.LoadX509KeyPair(p.opt.Ssl_certificate, p.opt.Ssl_key)

		if err != nil {
			return fmt.Errorf("Error loading keys: %v", err)
		}
		config.Certificates = []tls.Certificate{cert}

		// Certificate authority
		if len(p.opt.Ssl_certificate_authorities) > 0 {
			config.RootCAs = x509.NewCertPool()
			for _, pemCertPath := range p.opt.Ssl_certificate_authorities {
				pemCert, err := ioutil.ReadFile(pemCertPath)
				if err != nil {
					return fmt.Errorf("Error loading certificate authorities: %v", err)
				}
				config.RootCAs.AppendCertsFromPEM(pemCert)
			}
		}

		// SSL Verification mode
		if p.opt.Ssl_verify_mode == "peer" {
			config.ClientAuth = tls.VerifyClientCertIfGiven
		}
		if p.opt.Ssl_verify_mode == "force_peer" {
			config.ClientAuth = tls.RequireAndVerifyClientCert
		}

		options = append(options, v2.TLS(config))
	}

	options = append(options, v2.Timeout(time.Second*time.Duration(p.opt.Congestion_threshold)))

	log.Logger = p.Logger
	server, err := v2.ListenAndServe(fmt.Sprintf("%s:%d", p.opt.Host, p.opt.Port), options...)
	if err != nil {
		return err
	}

	p.server = server

	go func() {
		for batch := range p.server.ReceiveChan() {
			batch.ACK()
			events := batch.Events
			for _, e := range events {
				fields := e.(map[string]interface{})
				if val, ok := fields["@timestamp"]; !ok {
					fields["@timestamp"] = time.Now()
				} else {
					fields["@timestamp"], _ = time.Parse("2006-01-02T15:04:05Z07:00", val.(string))
				}

				ev := p.NewPacket("", fields)
				processors.ProcessCommonFields(ev.Fields(), p.opt.Add_field, p.opt.Tags, p.opt.Type)
				p.Send(ev, 0)
			}
		}
		p.Logger.Debug("received events acked and drained")
		close(p.q)
	}()

	return nil
}

func (p *processor) Stop(e processors.IPacket) error {
	err := p.server.Close()
	if err != nil {
		return err
	}
	<-p.q
	return nil
}
