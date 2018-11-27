package w3ccodec

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"bitfan/commons"
)

func TestBuffer(t *testing.T) {
	assert.Len(t, NewDecoder(strings.NewReader("")).Buffer(), 0)
}

func TestSetOptionsError(t *testing.T) {
	d := NewDecoder(strings.NewReader("data"))
	conf := map[string]interface{}{
		"separator": 4,
	}
	err := d.SetOptions(conf, nil, "")
	assert.Error(t, err)
}

func TestDefaultSettings(t *testing.T) {
	data :=
		`
#Software: Microsoft Internet Information Services 8.5
#Version: 1.0
#Date: 2017-10-04 19:00:00
#Fields: date time cs-uri-stem cs-username c-ip cs(User-Agent) sc-status
2017-10-04 19:00:00 /Microsoft-Server-ActiveSync/default.eas domain/user 8.8.8.8 Android/5.0.2-EAS-2.0 200
2017-10-04 22:12:13 /Microsoft-Server-ActiveSync/default.eas domain\iuser 44.44.88.88 Apple-iPhone6C2/1403.92 200
2017-10-05 00:45:01 /EWS/Exchange.asmx - 22.11.22.11 MS-WebServices/1.0 401
`

	expectData := []map[string]interface{}{
		map[string]interface{}{
			"date":           "2017-10-04",
			"time":           "19:00:00",
			"cs-uri-stem":    "/Microsoft-Server-ActiveSync/default.eas",
			"cs-username":    "domain/user",
			"c-ip":           "8.8.8.8",
			"cs(User-Agent)": "Android/5.0.2-EAS-2.0",
			"sc-status":      "200",
		},
		map[string]interface{}{
			"date":           "2017-10-04",
			"time":           "22:12:13",
			"cs-uri-stem":    "/Microsoft-Server-ActiveSync/default.eas",
			"cs-username":    "domain\\iuser",
			"c-ip":           "44.44.88.88",
			"cs(User-Agent)": "Apple-iPhone6C2/1403.92",
			"sc-status":      "200",
		},
		map[string]interface{}{
			"date":           "2017-10-05",
			"time":           "00:45:01",
			"cs-uri-stem":    "/EWS/Exchange.asmx",
			"cs-username":    "-",
			"c-ip":           "22.11.22.11",
			"cs(User-Agent)": "MS-WebServices/1.0",
			"sc-status":      "401",
		},
	}

	d := NewDecoder(strings.NewReader(data))
	var l commons.Logger
	err := d.SetOptions(map[string]interface{}{}, l, "")
	assert.NoError(t, err)

	var m interface{}

	for i := range expectData {
		err := d.Decode(&m)
		assert.NoError(t, err)
		assert.Equal(t, expectData[i], m)
	}

	err = d.Decode(&m)
	assert.EqualError(t, err, "EOF")
}

func TestTabSeparatedColumns(t *testing.T) {
	data := `
#Software: IIS Advanced Logging Module
#Version: 1.0
#Start-Date: 2017-09-23 00:00:00.064
#Fields:     date time s-computername s-ip cs-method cs-uri-stem cs-uri-query s-port cs-username c-ip cs(User-Agent) cs(Referer) cs(Host) sc-status sc-bytes cs-bytes TimeTakenMS
2017-09-23	00:03:11.235	"SP-WFE02"	192.168.65.82	GET	/	-	80	-	192.168.64.43	"Mozilla/4.0 (ISA Server Connectivity Check)"	-	"w3c.example.com"	401	378	186	0
2017-09-23	00:34:44.666	"SP-WFE02"	10.10.33.82	POST	/_vti_bin/Lists.asmx	-	80	-	192.168.66.178	"Mozilla/4.0 (compatible; MSIE 6.0; MS Web Services Client Protocol 4.0.30319.42000)"	-	"w3c.example.com"	201	378	686	0
`
	conf := map[string]interface{}{
		"separator": "\t",
	}

	expectData := []map[string]interface{}{
		map[string]interface{}{
			"date":           "2017-09-23",
			"time":           "00:03:11.235",
			"s-computername": "SP-WFE02",
			"s-ip":           "192.168.65.82",
			"cs-method":      "GET",
			"cs-uri-stem":    "/",
			"cs-uri-query":   "-",
			"s-port":         "80",
			"cs-username":    "-",
			"c-ip":           "192.168.64.43",
			"cs(User-Agent)": "Mozilla/4.0 (ISA Server Connectivity Check)",
			"cs(Referer)":    "-",
			"cs(Host)":       "w3c.example.com",
			"sc-status":      "401",
			"sc-bytes":       "378",
			"cs-bytes":       "186",
			"TimeTakenMS":    "0",
		},
		map[string]interface{}{
			"date":           "2017-09-23",
			"time":           "00:34:44.666",
			"s-computername": "SP-WFE02",
			"s-ip":           "10.10.33.82",
			"cs-method":      "POST",
			"cs-uri-stem":    "/_vti_bin/Lists.asmx",
			"cs-uri-query":   "-",
			"s-port":         "80",
			"cs-username":    "-",
			"c-ip":           "192.168.66.178",
			"cs(User-Agent)": "Mozilla/4.0 (compatible; MSIE 6.0; MS Web Services Client Protocol 4.0.30319.42000)",
			"cs(Referer)":    "-",
			"cs(Host)":       "w3c.example.com",
			"sc-status":      "201",
			"sc-bytes":       "378",
			"cs-bytes":       "686",
			"TimeTakenMS":    "0",
		},
	}

	d := NewDecoder(strings.NewReader(data))
	var l commons.Logger
	err := d.SetOptions(conf, l, "")
	assert.NoError(t, err)

	var m interface{}

	for i := range expectData {
		err = d.Decode(&m)
		assert.NoError(t, err)
		assert.Equal(t, expectData[i], m)
	}

	err = d.Decode(&m)
	assert.EqualError(t, err, "EOF")
}

func TestWithCustomColumns(t *testing.T) {
	data := `
#Version: 1.0
#Date: 12-Jan-1996 00:00:00
#Fields: time cs-method cs-uri
00:34:23 GET /foo/bar.html
12:21:16 POST /foo/ 201 "0"
12:57:34 HEAD
`
	conf := map[string]interface{}{
		"autogenerate_column_names": false,
		"columns":                   []string{"t", "method", "location"},
	}

	expectData := []map[string]interface{}{
		map[string]interface{}{
			"t":        "00:34:23",
			"method":   "GET",
			"location": "/foo/bar.html",
		},
		map[string]interface{}{
			"t":        "12:21:16",
			"method":   "POST",
			"location": "/foo/",
			"column4":  "201",
			"column5":  "0",
		},
		map[string]interface{}{
			"t":      "12:57:34",
			"method": "HEAD",
		},
	}

	d := NewDecoder(strings.NewReader(data))

	var l commons.Logger
	err := d.SetOptions(conf, l, "")
	assert.NoError(t, err)

	var m interface{}

	for i := range expectData {
		err = d.Decode(&m)
		assert.NoError(t, err)
		assert.Equal(t, expectData[i], m)
	}

	err = d.Decode(&m)
	assert.EqualError(t, err, "EOF")
}

func TestMore(t *testing.T) {
	data :=
		`
#Software: Microsoft Internet Information Services 8.5
#Version: 1.0
#Date: 2017-10-04 19:00:00
#Fields: date time cs-uri-stem cs-username c-ip cs(User-Agent) sc-status
2017-10-04 19:00:00 /Microsoft-Server-ActiveSync/default.eas domain/user 8.8.8.8 Android/5.0.2-EAS-2.0 200
2017-10-04 22:12:13 /Microsoft-Server-ActiveSync/default.eas domain\iuser 44.44.88.88 Apple-iPhone6C2/1403.92 200
2017-10-05 00:45:01 /EWS/Exchange.asmx - 22.11.22.11 MS-WebServices/1.0 401
`

	expectData := []map[string]interface{}{
		map[string]interface{}{
			"date":           "2017-10-04",
			"time":           "19:00:00",
			"cs-uri-stem":    "/Microsoft-Server-ActiveSync/default.eas",
			"cs-username":    "domain/user",
			"c-ip":           "8.8.8.8",
			"cs(User-Agent)": "Android/5.0.2-EAS-2.0",
			"sc-status":      "200",
		},
		map[string]interface{}{
			"date":           "2017-10-04",
			"time":           "22:12:13",
			"cs-uri-stem":    "/Microsoft-Server-ActiveSync/default.eas",
			"cs-username":    "domain\\iuser",
			"c-ip":           "44.44.88.88",
			"cs(User-Agent)": "Apple-iPhone6C2/1403.92",
			"sc-status":      "200",
		},
		map[string]interface{}{
			"date":           "2017-10-05",
			"time":           "00:45:01",
			"cs-uri-stem":    "/EWS/Exchange.asmx",
			"cs-username":    "-",
			"c-ip":           "22.11.22.11",
			"cs(User-Agent)": "MS-WebServices/1.0",
			"sc-status":      "401",
		},
	}

	d := NewDecoder(strings.NewReader(data))
	var l commons.Logger
	err := d.SetOptions(map[string]interface{}{}, l, "")
	assert.NoError(t, err)

	var m interface{}

	var i = 0
	for d.More() {
		err := d.Decode(&m)
		if i+1 <= len(expectData) {
			assert.NoError(t, err)
			assert.Equal(t, expectData[i], m)
			i = i + 1
		} else {
			assert.Error(t, err)
		}

	}
	assert.Equal(t, 3, i)
}
