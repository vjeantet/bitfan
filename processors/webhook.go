package processors

import "net/http"

type WebHook interface {
	Add(string, func(http.ResponseWriter, *http.Request))
}
