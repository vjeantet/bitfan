//go:generate bitfanDoc
package beatsinput

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"sync"
	"time"

	"github.com/vjeantet/bitfan/processors"
)

func (p *processor) serve() error {
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", p.opt.Host, p.opt.Port))
	if err != nil {
		return fmt.Errorf("Listener failed: %v", err)
	}

	clientTerm := make(chan bool)
	var wg sync.WaitGroup
	for {
		select {
		case <-p.q:
			ln.Close()
			close(clientTerm)
			wg.Wait()
			close(p.q)
			return nil
		default:
		}

		if l, ok := ln.(*net.TCPListener); ok {
			l.SetDeadline(time.Now().Add(1 * time.Second))
		}

		conn, err := ln.Accept()
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			// log.Printf("Error accepting connection: %v", err)
			continue
		}

		if p.opt.Ssl == true {
			config := tls.Config{}

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

			conn = tls.Server(conn, &config)
		}

		wg.Add(1)
		go p.clientServe(conn, &wg, clientTerm)
	}

	return nil
}

// lumberConn handles an incoming connection from a lumberjack client
func (p *processor) clientServe(c net.Conn, wg *sync.WaitGroup, clientTerm chan bool) {
	defer wg.Done()
	defer c.Close()

	// log.Printf("[%s] accepting lumberjack connection", c.RemoteAddr().String())

	dataChan := make(chan map[string]interface{}, 3)
	go NewParser(c, dataChan).Parse()

	for {
		select {
		case fields := <-dataChan:
			if fields == nil {
				// log.Printf("[%s] closing lumberjack connection", c.RemoteAddr().String())
				return
			}
			msg := ""
			if txt, ok := fields["message"]; ok {
				msg = txt.(string)
			}
			e := p.NewPacket(msg, fields)
			processors.ProcessCommonFields(e.Fields(), p.opt.Add_field, p.opt.Tags, p.opt.Type)
			p.Send(e)
		case <-clientTerm:
			c.SetReadDeadline(time.Now())
		}
	}

}
