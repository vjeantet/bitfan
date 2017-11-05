package server

import (
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin/render"
	"github.com/k0kubun/pp"
)

type Render struct {
	Templates          map[string]*template.Template
	AssetsDir          string
	TemplatesDir       string
	AssetsTemplatesDir string
	Layout             string
	Ext                string
	TemplateFuncMap    map[string]interface{}
	Debug              bool
	UseFS              bool
}

func NewRender() Render {
	r := Render{

		Templates: map[string]*template.Template{},

		AssetsDir: "webui/",
		// TemplatesDir holds the location of the templates
		TemplatesDir: "assets/views/",
		// Layout is the file name of the layout file
		Layout: "layout.html",
		// Ext is the file extension of the rendered templates
		Ext: ".html",
		// Template's function map
		TemplateFuncMap: nil,
		// Debug enables debug mode
		Debug: false,

		UseFS: false,
	}

	return r
}

func (r Render) Glob(pattern string) ([]string, error) {
	var matches []string
	if r.UseFS == true {
		return filepath.Glob(pattern)
	} else {
		for _, key := range AssetNames() {
			if ok, err := filepath.Match(pattern, key); ok && err == nil {
				matches = append(matches, key)
			}
		}
		return matches, nil
	}
}

func (r Render) Init() Render {
	tmpDir := filepath.Join(r.AssetsDir, r.TemplatesDir)
	if _, err := os.Stat(tmpDir); !os.IsNotExist(err) {
		// assets exists on disk, UseFS
		if r.UseFS == false {
			r.TemplatesDir = filepath.Join(r.AssetsDir, r.TemplatesDir) + string(os.PathSeparator)
			r.UseFS = true
		}
	} else {
		// assets from bindData
		r.UseFS = false
	}

	layout := r.TemplatesDir + r.Layout

	viewDirs, _ := r.Glob(r.TemplatesDir + "**" + string(os.PathSeparator) + "*" + r.Ext)
	partials, _ := r.Glob(r.TemplatesDir + "_*" + r.Ext)

	for _, view := range viewDirs {
		renderName := r.getRenderName(view)
		if r.Debug {
			log.Printf("[GIN-debug] %-6s %-25s --> %s\n", "LOAD", view, renderName)
		}
		views := append(partials, layout, view)
		r.AddFromFiles(renderName, views...)
	}

	return r
}

func (r Render) getRenderName(tpl string) string {
	tmp := strings.TrimPrefix(tpl, r.TemplatesDir)
	tmp = strings.TrimSuffix(tmp, r.Ext)
	tmp = strings.Replace(tmp, "\"", "/", -1)
	return tmp
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
	tmpl := template.New(filepath.Base(r.Layout)).Funcs(r.TemplateFuncMap)
	for _, f := range files {
		templateString, err := ioutil.ReadFile(f)
		if err != nil {
			if r.Debug {
				pp.Printf("[GIN-debug] reading from BindData --> %s (%s)\n", f, err.Error())
			}
			templateString, err = Asset(f)
		} else {
			if r.Debug {
				pp.Printf("[GIN-debug] reading from FS --> %s\n", f)
			}
		}

		if err != nil {
			log.Fatal(err)
		}

		tmpl, err = tmpl.Parse(string(templateString))
		if err != nil {
			log.Fatal(err)
		}
	}
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
