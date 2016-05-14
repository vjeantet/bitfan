package runtime

import (
	"net/http"

	"github.com/gosimple/slug"
)

type HTTPHookHandler func(http.ResponseWriter, *http.Request)

type WebHook struct {
	mux       *http.ServeMux
	namespace string
	uri       string
	Hooks     map[string]*hook
}

type hook struct {
	Url     string
	handler *HTTPHookHandler
}

func NewWebHook(nameSpace string) *WebHook {
	return &WebHook{namespace: nameSpace, mux: httpHookServerMux, Hooks: map[string]*hook{}}
}

// Add register a new route matcher linked to hh
// TODO: the route matcher
func (w *WebHook) Add(hookName string, hh HTTPHookHandler) {
	h := &hook{}
	h.Url = slug.Make(w.namespace) + "/" + slug.Make(hookName)
	h.handler = &hh
	w.Hooks[hookName] = h
	w.mux.HandleFunc("/"+h.Url, *h.handler)
	Logger().Infof("Hook %s => %s", hookName, "/"+h.Url)
}

func GetAgentHooks(agentName string) (hooks map[string]*hook) {
	if _, ok := webHookMap[agentName]; ok {
		hooks = webHookMap[agentName].Hooks
	}
	return hooks
}

var httpHookServerMux *http.ServeMux

func listenAndServeHTTP(addr string) {
	httpHookServerMux = http.NewServeMux()

	Logger().Infof("Agents webHook listening on %s", addr)
	go http.ListenAndServe(addr, httpHookServerMux)
}
