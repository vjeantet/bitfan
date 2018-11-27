//go:generate bitfanDoc
package beatsinput

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/elastic/go-lumber/server"
	"bitfan/processors"
)

func New() processors.Processor {
	return &processor{opt: &options{}}
}

type processor struct {
	processors.Base

	server server.Server
	opt    *options
}

type options struct {
	processors.CommonOptions `mapstructure:",squash"`

	// The number of seconds before we raise a timeout,
	// this option is useful to control how much time to wait if something is blocking
	// the pipeline
	CongestionThreshold int

	// The IP address to listen on
	Host string

	// The port to listen on (default 5044)
	Port int

	// Events are by default send in plain text,
	// you can enable encryption by using ssl to true and
	// configuring the ssl_certificate and ssl_key options
	SSL bool

	// SSL certificate to use (path)
	SSLCertificate string

	// Validate client certificates against theses authorities
	//  You can defined multiples files or path, all the certificates will be read
	//  and added to the trust store.
	//  You need to configure the ssl_verify_mode to peer or force_peer to enable
	//  the verification.
	// This feature only support certificate directly signed by your root ca.
	// Intermediate CA are currently not supported.
	SSlCertificateAuthorities []string

	// SSL key to use (path)
	SSlKey string

	// SSL key passphrase to use. (not yet implemented)
	SSlKeyPassphrase string

	// By default the server dont do any client verification,
	// peer will make the server ask the client to provide a certificate,
	//   if the client provide the certificate it will be validated.
	// force_peer will make the server ask the client for their certificate,
	//   if the clients doesn’t provide it the connection will be closed.
	// This option need to be used with ssl_certificate_authorities and a defined list of CA.
	// Value can be any of: none, peer, force_peer
	SSlVerifyMode string
}

func (p *processor) Configure(ctx processors.ProcessorContext, conf map[string]interface{}) error {
	p.opt.CongestionThreshold = 30
	p.opt.Host = "0.0.0.0"
	p.opt.Port = 5044
	p.opt.SSL = false
	p.opt.SSlVerifyMode = "none"

	return p.ConfigureAndValidate(ctx, conf, p.opt)
}

func (p *processor) Start(e processors.IPacket) error {

	options := make([]server.Option, 0)

	if p.opt.SSL == true {
		config := new(tls.Config)

		// Server Certificates
		cert, err := tls.LoadX509KeyPair(p.opt.SSLCertificate, p.opt.SSlKey)

		if err != nil {
			return fmt.Errorf("error loading keys: %v", err)
		}
		config.Certificates = []tls.Certificate{cert}

		// Certificate authority
		if len(p.opt.SSlCertificateAuthorities) > 0 {
			config.RootCAs = x509.NewCertPool()
			for _, pemCertPath := range p.opt.SSlCertificateAuthorities {
				pemCert, err := ioutil.ReadFile(pemCertPath)
				if err != nil {
					return fmt.Errorf("error loading certificate authorities: %v", err)
				}
				config.RootCAs.AppendCertsFromPEM(pemCert)
			}
		}

		// SSL Verification mode
		if p.opt.SSlVerifyMode == "peer" {
			config.ClientAuth = tls.VerifyClientCertIfGiven
		}
		if p.opt.SSlVerifyMode == "force_peer" {
			config.ClientAuth = tls.RequireAndVerifyClientCert
		}

		options = append(options, server.TLS(config))
	}

	options = append(options, server.Timeout(time.Second * time.Duration(p.opt.CongestionThreshold)))

	svr, err := server.ListenAndServe(fmt.Sprintf("%s:%d", p.opt.Host, p.opt.Port), options...)
	if err != nil {
		return err
	}

	p.server = svr

	go func(p *processor) {
		rchan := p.server.ReceiveChan()

		for {
			var closed bool
			select {
			case batch, ok := <-rchan:
				if ok {
					for _, evt := range batch.Events {

						fields := evt.(map[string]interface{})
						if val, ok := fields["@timestamp"]; !ok {
							fields["@timestamp"] = time.Now()
						} else {
							fields["@timestamp"], _ = time.Parse(time.RFC3339, val.(string))
						}

						ev := p.NewPacket(fields)
						p.opt.ProcessCommonOptions(ev.Fields())
						p.Send(ev)
					}
					batch.ACK()
				} else {
					closed = true
				}
			}
			if closed {
				break
			}
		}
	}(p)

	return nil
}

func (p *processor) Stop(e processors.IPacket) error {
	err := p.server.Close()
	if err != nil {
		return err
	}
	return nil
}
