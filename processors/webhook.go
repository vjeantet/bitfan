package processors

import "net/http"

type WebHook interface {
	Add(string, http.HandlerFunc)
	AddShort(string, http.HandlerFunc)
	Unregister()
}
