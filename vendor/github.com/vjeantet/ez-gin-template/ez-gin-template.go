package eztemplate

import (
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin/render"
	"github.com/k0kubun/pp"
)

type Render struct {
	Templates       map[string]*template.Template
	TemplatesDir    string
	Layout          string
	Ext             string
	TemplateFuncMap map[string]interface{}
	Debug           bool
}

func New() Render {
	r := Render{

		Templates: map[string]*template.Template{},
		// TemplatesDir holds the location of the templates
		TemplatesDir: "app/views/",
		// Layout is the file name of the layout file
		Layout: "layouts/base",
		// Ext is the file extension of the rendered templates
		Ext: ".html",
		// Template's function map
		TemplateFuncMap: nil,
		// Debug enables debug mode
		Debug: false,
	}

	return r
}

func (r Render) Init() Render {
	layout := r.TemplatesDir + r.Layout + r.Ext

	viewDirs, _ := filepath.Glob(r.TemplatesDir + "**" + string(os.PathSeparator) + "*" + r.Ext)

	for _, view := range viewDirs {
		renderName := r.getRenderName(view)
		if r.Debug {
			log.Printf("[GIN-debug] %-6s %-25s --> %s\n", "LOAD", view, renderName)
		}
		r.AddFromFiles(renderName, layout, view)
	}

	return r
}

func (r Render) getRenderName(tpl string) string {
	dir, file := filepath.Split(tpl)
	dir = strings.Replace(dir, string(os.PathSeparator), "/", -1)
	tempdir := strings.Replace(r.TemplatesDir, string(os.PathSeparator), "/", -1)
	dir = strings.Replace(dir, tempdir, "", 1)
	file = strings.TrimSuffix(file, r.Ext)
	return dir + file
}

func (r Render) Add(name string, tmpl *template.Template) {
	if tmpl == nil {
		panic("template can not be nil")
	}
	if len(name) == 0 {
		panic("template name cannot be empty")
	}
	r.Templates[name] = tmpl
}

func (r Render) AddFromFiles(name string, files ...string) *template.Template {
	tmpl := template.Must(template.New(filepath.Base(r.Layout + r.Ext)).Funcs(r.TemplateFuncMap).ParseFiles(files...))
	r.Add(name, tmpl)
	return tmpl
}

func (r Render) Instance(name string, data interface{}) render.Render {
	if r.Debug == true {
		r.Init()
		pp.Println("data -->", data)
	}

	return render.HTML{
		Template: r.Templates[name],
		Data:     data,
	}
}
