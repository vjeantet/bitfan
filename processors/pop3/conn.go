package pop3processor

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"time"

	pop3 "github.com/taknb2nch/go-pop3"
)

var (
	EOF = errors.New("skip the all mail remaining")
)

func receiveMail(c *pop3.Client, lastSeenUID string, receiveFn receiveMailFunc) error {
	var err error
	defer func() {
		if err != nil && err != EOF {
			c.Rset()
		}

		c.Quit()
		c.Close()
	}()

	var mis []pop3.MessageInfo

	if mis, err = c.UidlAll(); err != nil {
		return err
	}

	// reverse order of mis id
	// create a new mis2
	mis_filtered := []pop3.MessageInfo{}
	for i := len(mis) - 1; i >= 0; i-- {
		mi := mis[i]
		if mi.Uid == lastSeenUID {
			break
		}
		mis_filtered = append(mis_filtered, mi)
	}

	// iterate over mis id and
	// push id to mis2
	// when id == lastseen ; break
	// iterate over mis

	// for _, mi := range mis_filtered { // natural order (FIRST OUT)
	for i := len(mis_filtered) - 1; i >= 0; i-- { // reverse order (LAST FIRST)
		mi := mis_filtered[i]

		var data string

		data, err = c.Retr(mi.Number)
		del, err := receiveFn(mi.Number, mi.Uid, data, err)

		if err != nil && err != EOF {
			return err
		}

		if del {
			if err = c.Dele(mi.Number); err != nil {
				return err
			}
		}

		if err == EOF {
			break
		}
	}

	return nil
}

type receiveMailFunc func(number int, uid, data string, err error) (bool, error)

func newPop3Client(host string, port int, username string, password string, secure bool, verifyCert bool, timeout int) (*pop3.Client, error) {
	var err error
	var client *pop3.Client
	defer func() {
		if err != nil && err != EOF {
			client.Rset()
			client.Quit()
			client.Close()
		}
	}()

	defaultTimeout := time.Duration(timeout) * time.Second

	var conn net.Conn
	conn, err = createConn(host, port, secure, !verifyCert, defaultTimeout)
	if err != nil {
		return nil, fmt.Errorf("create connection to %v:%v failed: %v", host, port, err)
	}

	client, err = pop3.NewClient(conn)

	if err = client.User(username); err != nil {
		return nil, fmt.Errorf("Username Error: %v\n", err)
	}

	if err = client.Pass(password); err != nil {
		return nil, fmt.Errorf("Password Error: %v\n", err)
	}

	return client, nil
}

func createConn(
	host string,
	port int,
	tlsActive bool,
	tlsSkipVerify bool,
	timeout time.Duration,
) (net.Conn, error) {
	address := fmt.Sprintf("%v:%v", host, port)
	dailer := &net.Dialer{Timeout: timeout}
	if tlsActive {
		tlsconfig := &tls.Config{
			InsecureSkipVerify: tlsSkipVerify,
			ServerName:         host,
		}
		return tls.DialWithDialer(dailer, "tcp", address, tlsconfig)
	}
	return dailer.Dial("tcp", address)
}
