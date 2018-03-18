package processors

import "net/http"

type WebHook interface {
	Add(string, http.HandlerFunc)
	AddNamed(string, http.HandlerFunc)
	Unregister()
}
