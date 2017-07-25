package core

import (
	"fmt"
	"net/http"

	"golang.org/x/sync/syncmap"

	"github.com/gosimple/slug"
	"github.com/justinas/alice"
)

type webHook struct {
	pipelineLabel string
	namespace     string
	Hooks         []string
}

var webHookMap = syncmap.Map{}
var httpHookServerMux *http.ServeMux

func newWebHook(pipelineLabel, nameSpace string) *webHook {
	return &webHook{pipelineLabel: pipelineLabel, namespace: nameSpace, Hooks: []string{}}
}

// Add a new route to a given http.HandlerFunc
func (w *webHook) Add(hookName string, hf http.HandlerFunc) {
	hUrl := "/" + slug.Make(w.pipelineLabel) + "/" + slug.Make(w.namespace) + "/" + slug.Make(hookName)
	w.Hooks = append(w.Hooks, hookName)
	webHookMap.Store(hUrl, hf)
	Log().Infof("Hook %s => %s", hookName, hUrl)
}

// Delete a route
func (w *webHook) Delete(hookName string) {
	hUrl := "/" + slug.Make(w.pipelineLabel) + "/" + slug.Make(w.namespace) + "/" + slug.Make(hookName)
	webHookMap.Delete(hUrl)
	Log().Debugf("WebHook unregisted [%s]", hUrl)
}

// Delete all routes belonging to webHook
func (w *webHook) Unregister() {
	for _, hookName := range w.Hooks {
		w.Delete(hookName)
	}
}

func routerHandler(w http.ResponseWriter, r *http.Request) {
	if hfi, ok := webHookMap.Load(r.URL.Path); ok {
		Log().Debugf("Webhook found for %s", r.URL.Path)
		hfi.(http.HandlerFunc)(w, r)
	} else {
		Log().Warnf("Webhook not found for %s", r.URL.Path)
		w.WriteHeader(404)
		fmt.Fprint(w, "Not Found !")
	}
}

func listenAndServeWebHook(addr string) {
	httpHookServerMux = http.NewServeMux()
	commonHandlers := alice.New(loggingHandler, recoverHandler)
	httpHookServerMux.Handle("/", commonHandlers.ThenFunc(routerHandler))
	go http.ListenAndServe(addr, httpHookServerMux)

	Log().Infof("Agents webHook listening on %s", addr)
}

func loggingHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		Log().Debugf("Webhook [%s] %s", r.Method, r.URL.Path)
	}
	return http.HandlerFunc(fn)
}

func recoverHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				Log().Errorf("Webhook panic [%s] %s : %+v", r.Method, r.URL.Path, err)
				http.Error(w, http.StatusText(500), 500)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
