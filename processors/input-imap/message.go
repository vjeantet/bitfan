package imap_input

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/mail"
	"time"

	"github.com/vjeantet/go.enmime"
)

type Msg struct {
	Header      map[string]interface{}
	UID         int
	Date        string
	Text        string // Plain utf8 text
	Html        string // Plain utf8 text
	Attachments []Part
	Inlines     []Part
	Others      []Part
}
type Part struct {
	Filename    string
	ContentType string
	Content     []byte
}

func getMsg(emailbody string) (m *Msg) {
	msg, _ := mail.ReadMessage(bytes.NewBufferString(emailbody))
	m = &Msg{}
	if enmime.IsMultipartMessage(msg) {
		mime, err := enmime.ParseMIMEBody(msg)
		if err != nil {
			//fmt.Println("Trying to read", id)
			//panic(err)
			m.Text = err.Error()
		} else {
			log.Printf("mime.Text = %#v", mime.Text)
			m.Text = mime.Text
			m.Html = mime.Html
			msg.Header["Subject"] = []string{mime.GetHeader("Subject")}
			m.Attachments = []Part{}
			for _, v := range mime.Attachments {
				m.Attachments = append(m.Attachments, Part{
					Filename:    v.FileName(),
					ContentType: v.ContentType(),
					Content:     v.Content(),
				})
			}
			m.Inlines = []Part{}
			for _, v := range mime.Inlines {
				m.Inlines = append(m.Inlines, Part{
					Filename:    v.FileName(),
					ContentType: v.ContentType(),
					Content:     v.Content(),
				})
			}
			m.Others = []Part{}
			for _, v := range mime.OtherParts {
				m.Others = append(m.Others, Part{
					Filename:    v.FileName(),
					ContentType: v.ContentType(),
					Content:     v.Content(),
				})
			}
		}
	} else {
		body, _ := ioutil.ReadAll(msg.Body)
		m.Text = string(body)
	}
	//m.UID = id

	m.Header = make(map[string]interface{})
	// We need to map values on the interface explicitly
	for k, v := range msg.Header {
		if k == "Date" {
			m.Date = v[0]
		} else {
			m.Header[k] = v
		}
	}

	t, err := time.Parse(time.RFC1123Z, m.Date)
	if err == nil {
		m.Date = t.Format(time.RFC3339)
	}

	for _, a := range []string{"To", "From", "CC"} {
		addrs, err := msg.Header.AddressList(a)
		if err != nil {
			continue
		}

		m.Header[a] = addrs
	}

	return m
}
