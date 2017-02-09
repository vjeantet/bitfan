package beatsinput

import "github.com/vjeantet/bitfan/processors"

func New() processors.Processor {
	return &processor{opt: &options{}}
}

type processor struct {
	processors.Base

	opt *options
	q   chan bool
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
	p.opt.Congestion_threshold = 5
	p.opt.Host = "0.0.0.0"
	p.opt.Port = 5044
	p.opt.Ssl = false
	p.opt.Ssl_verify_mode = "none"

	return p.ConfigureAndValidate(ctx, conf, p.opt)
}

func (p *processor) Start(e processors.IPacket) error {
	p.q = make(chan bool)
	go p.serve()
	return nil
}

func (p *processor) Stop(e processors.IPacket) error {
	p.q <- true
	<-p.q
	return nil
}
