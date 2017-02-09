package handler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type PostBackHandler struct {
	Url           string
	PostEncoded   bool
	PostParamName string
}

func (hnd *PostBackHandler) Deliver(message string) error {
	var err error

	req, err := newPostRequest(hnd.Url, hnd.getPostBody(message))
	if err != nil {
		return fmt.Errorf("Could not deliver: %s", err)
	}

	req.Header.Add("Content-Type", hnd.getContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Request into postback hook failed: %s", err)
	}

	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return fmt.Errorf("An error occurred while reading hook response: %s", err)
	}

	if !responseOk(resp.StatusCode) {
		return fmt.Errorf("Hook returned with error: %s\n%q", resp.Status, data)
	}

	return nil
}

func (hnd *PostBackHandler) Describe() string {
	desc := "plain"
	if hnd.PostEncoded {
		desc = fmt.Sprintf("urlencoded[%s]", hnd.PostParamName)
	}

	return fmt.Sprintf("PostbackHandler (url=%s, %s)", redactedURL(hnd.Url), desc)
}

func (hnd *PostBackHandler) getPostBody(raw string) string {
	if !hnd.PostEncoded {
		return raw
	}

	data := url.Values{}
	data.Set(hnd.PostParamName, raw)
	return data.Encode()
}

func (hnd *PostBackHandler) getContentType() string {
	if !hnd.PostEncoded {
		return "text/plain"
	}

	return "application/x-www-form-urlencoded"
}

func NewPostBackHandler(postUrl string, postEncoded bool, postParamName string) *PostBackHandler {
	return &PostBackHandler{
		Url:           postUrl,
		PostEncoded:   postEncoded,
		PostParamName: postParamName}
}

func newPostRequest(endpoint string, payload string) (*http.Request, error) {
	uri, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("Malformed postback hook url: %s", err)
	}

	buff := strings.NewReader(payload)
	req, err := http.NewRequest("POST", endpoint, buff)
	if err != nil {
		return nil, fmt.Errorf("Unable to build request object: %s", err)
	}

	req.Header.Add("Host", uri.Host)
	req.Header.Add("User-Agent", "postman-postback")
	req.Header.Add("Accept", "*/*")

	return req, nil
}

func responseOk(status int) bool {
	return !(status != 200 && status != 201 && status != 204)
}

func redactedURL(u string) string {
	uri, err := url.Parse(u)
	if err != nil {
		return u
	}

	return uri.Scheme + "://" + uri.Host + uri.Path
}
