package core

import (
	"net/http"

	"github.com/gosimple/slug"
)

var webHookMap = map[string]*webHook{}

type webHook struct {
	mux       *http.ServeMux
	namespace string
	uri       string
	Hooks     map[string]*hook
}

type hook struct {
	Url     string
	handler *func(http.ResponseWriter, *http.Request)
}

func newWebHook(nameSpace string) *webHook {
	return &webHook{namespace: nameSpace, mux: httpHookServerMux, Hooks: map[string]*hook{}}
}

// Add register a new route matcher linked to hh
func (w *webHook) Add(hookName string, hh func(http.ResponseWriter, *http.Request)) {
	h := &hook{}
	h.Url = slug.Make(w.namespace) + "/" + slug.Make(hookName)
	h.handler = &hh
	w.Hooks[hookName] = h
	w.mux.HandleFunc("/"+h.Url, *h.handler)
	Log().Infof("Hook %s => %s", hookName, "/"+h.Url)
}

func getAgentHooks(agentName string) (hooks map[string]*hook) {
	if _, ok := webHookMap[agentName]; ok {
		hooks = webHookMap[agentName].Hooks
	}
	return hooks
}

var httpHookServerMux *http.ServeMux

func listenAndServeWebHook(addr string) {
	httpHookServerMux = http.NewServeMux()

	Log().Infof("Agents webHook listening on %s", addr)
	go http.ListenAndServe(addr, httpHookServerMux)
}
