package processors

import "net/http"

type WebHook interface {
	Add(string, http.HandlerFunc)
	Unregister()
}
