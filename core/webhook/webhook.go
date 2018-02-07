package webhook

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/sync/syncmap"

	fqdn "github.com/ShowMax/go-fqdn"
	"github.com/gosimple/slug"
	"github.com/justinas/alice"
	"github.com/vjeantet/bitfan/commons"
)

type webHook struct {
	pipelineLabel string
	namespace     string
	Hooks         []string
}

type Hook struct {
	h            http.HandlerFunc
	Description  string
	Namespace    string
	PipelineUUID string
	Url          string
}

var webHookMap = syncmap.Map{}
var baseURL = ""
var whPrefixURL = "/h/"
var Log commons.Logger

func New(pipelineLabel, nameSpace string) *webHook {
	return &webHook{pipelineLabel: pipelineLabel, namespace: nameSpace, Hooks: []string{}}
}

func WebHooks(uuid string) []Hook {
	urls := []Hook{}
	webHookMap.Range(func(key, value interface{}) bool {
		if uuid == value.(*Hook).PipelineUUID {
			urls = append(urls, Hook{
				Description: value.(*Hook).Description,
				Namespace:   value.(*Hook).Namespace,
				Url:         value.(*Hook).Url,
			})
		}
		return true
	})
	return urls
}

func (w *webHook) buildURL(hookName string) string {
	return strings.ToLower(whPrefixURL + slug.Make(w.pipelineLabel) + "/" + slug.Make(hookName))
}

// Add a new route to a given http.HandlerFunc
func (w *webHook) Add(hookName string, hf http.HandlerFunc) {

	hUrl := w.buildURL(hookName)
	w.Hooks = append(w.Hooks, hookName)

	webHookMap.Store(hUrl, &Hook{
		h:            hf,
		Description:  hookName,
		Namespace:    w.namespace,
		PipelineUUID: w.pipelineLabel,
		Url:          hUrl,
	})
	Log.Infof("Hook [%s - %s] %s", w.pipelineLabel, w.namespace, baseURL+hUrl)
}

// Delete a route
func (w *webHook) Delete(hookName string) {
	hUrl := w.buildURL(hookName)
	webHookMap.Delete(hUrl)
	Log.Debugf("WebHook unregisted [%s]", hUrl)
}

// Delete all routes belonging to webHook
func (w *webHook) Unregister() {
	for _, hookName := range w.Hooks {
		w.Delete(hookName)
	}
}

func Handler(host string) http.Handler {
	addrSpit := strings.Split(host, ":")
	if addrSpit[0] == "0.0.0.0" {
		addrSpit[0] = fqdn.Get()
	}
	baseURL = fmt.Sprintf("http://%s:%s", addrSpit[0], addrSpit[1])

	commonHandlers := alice.New(loggingHandler, recoverHandler)
	return commonHandlers.ThenFunc(routerHandler)
}

func routerHandler(w http.ResponseWriter, r *http.Request) {
	hUrl := strings.ToLower(r.URL.Path)
	if tHook, ok := webHookMap.Load(hUrl); ok {
		Log.Debugf("Webhook found for %s", hUrl)
		tHook.(*Hook).h(w, r)
	} else {
		Log.Warnf("Webhook not found for %s", hUrl)
		w.WriteHeader(404)
		fmt.Fprint(w, "Not Found !")
	}
}

func loggingHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
		Log.Debugf("Webhook [%s] %s", r.Method, r.URL.Path)
	}
	return http.HandlerFunc(fn)
}

func recoverHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				Log.Errorf("Webhook panic [%s] %s : %+v", r.Method, r.URL.Path, err)
				http.Error(w, http.StatusText(500), 500)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
