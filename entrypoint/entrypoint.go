// Entrypoint manage pipeline's definitions to get Pipeline ready to be used by the core
package entrypoint

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/vjeantet/bitfan/core/config"
	"github.com/vjeantet/bitfan/entrypoint/parser"
)

const (
	CONTENT_REF     = iota // Content is a reference to something
	CONTENT_REF_FS         // Content is a reference to something on the filesystem
	CONTENT_REF_URL        // Content is a reference to something on the web (http, https)
	CONTENT_INLINE         // Content is a value
)

// Entrypoint is a the pipeline's definition ressource
type Entrypoint struct {
	Path         string
	Kind         int // Kind of content
	Workingpath  string
	Content      string
	PipelineName string
	PipelineUuid string
}

// List of Entrypoints
type EntrypointList struct {
	Items []*Entrypoint
}

// Create a new entrypoint (pipeline definition)
//
// - contentValue may be a filesystem path, a URL or a string,
//
// - cwl Working Location should be provided to the parser, it could be an filesystem dir, a baseUrl base path, this part is
// used when the entrypoint contains references to other configurations. @see use, route processors.
//
// - contentKind refer to the kind of contentValue @see CONTENT_* constants
func New(contentValue string, cwl string, contentKind int) (*Entrypoint, error) {
	loc := &Entrypoint{}

	if contentKind == CONTENT_INLINE {
		loc.Kind = CONTENT_INLINE
		loc.Content = contentValue
		return loc, nil
	}

	if v, _ := url.Parse(contentValue); v.Scheme == "http" || v.Scheme == "https" {
		loc.Kind = CONTENT_REF_URL
		loc.Path = contentValue
	} else if _, err := os.Stat(contentValue); err == nil {
		var err error
		loc.Kind = CONTENT_REF_FS
		loc.Path, err = filepath.Abs(contentValue)
		if err != nil {
			return loc, err
		}
	} else if _, err := os.Stat(filepath.Join(cwl, contentValue)); err == nil {
		loc.Kind = CONTENT_REF_FS
		loc.Path = filepath.Join(cwl, contentValue)
	} else if v, _ := url.Parse(cwl); v.Scheme == "http" || v.Scheme == "https" {
		loc.Kind = CONTENT_REF_URL
		loc.Path = cwl + contentValue
	} else {
		return nil, fmt.Errorf("can not find any configuration contentValue=%s, cwl=%s", contentValue, cwl)
	}

	loc.Workingpath = cwl
	return loc, nil
}

// AddEntrypoint add the provided entrypoint to the list
func (e *EntrypointList) AddEntrypoint(loc *Entrypoint) error {
	// if it's a file try to expand
	if loc.Kind == CONTENT_REF_FS {
		subpaths, err := expandFilePath(loc.Path)
		if err != nil {
			return err
		}
		if len(subpaths) == 1 {
			e.Items = append(e.Items, loc)
		} else {
			for _, subpath := range subpaths {
				subloc := &Entrypoint{
					Path:        subpath,
					Workingpath: loc.Workingpath,
					Kind:        loc.Kind,
				}
				e.Items = append(e.Items, subloc)
			}
		}

		return nil
	}

	e.Items = append(e.Items, loc)

	return nil
}

// ConfigPipeline returns a core Pipeline from entrypoint's definition
func (e *Entrypoint) ConfigPipeline() config.Pipeline {
	var pipeline *config.Pipeline

	switch e.Kind {
	case CONTENT_INLINE:
		pipeline = config.NewPipeline("inline", "nodescription", "inline")
	case CONTENT_REF_URL:
		uriSegments := strings.Split(e.Path, "/")
		pipelineName := strings.Join(uriSegments[2:], ".")
		pipeline = config.NewPipeline(pipelineName, "nodescription", e.Path)
	case CONTENT_REF_FS:
		filename := filepath.Base(e.Path)
		extension := filepath.Ext(filename)
		pipelineName := filename[0 : len(filename)-len(extension)]
		pipeline = config.NewPipeline(pipelineName, "nodescription", e.Path)
	}

	if e.PipelineName != "" {
		pipeline.Name = e.PipelineName
	}
	if e.PipelineUuid != "" {
		pipeline.Uuid = e.PipelineUuid
	}

	return *pipeline
}

// ConfigPipeline returns core agents from entrypoint's definition
func (e *Entrypoint) ConfigAgents() ([]config.Agent, error) {
	var agents []config.Agent
	var content []byte
	var err error
	var cwd string

	content, cwd, err = e.content(map[string]interface{}{})
	if err != nil {
		return agents, err
	}

	agents, err = parser.BuildAgents(content, cwd, entrypointContent)
	return agents, err
}

func entrypointContent(path string, cwl string, options map[string]interface{}) ([]byte, string, error) {
	e, err := New(path, cwl, CONTENT_REF)
	if err != nil {
		return nil, "", err
	}
	return e.content(options)
}

func (e *Entrypoint) content(options map[string]interface{}) ([]byte, string, error) {
	var content []byte
	var cwl string
	var err error

	switch e.Kind {
	case CONTENT_INLINE:
		content = []byte(e.Content)
		cwl = e.Workingpath

	case CONTENT_REF_URL:
		response, err := http.Get(e.Path)
		if err != nil {
			return content, cwl, err
		} else {
			content, err = ioutil.ReadAll(response.Body)
			response.Body.Close()
			if err != nil {
				return content, cwl, err
			}
		}

		uriSegments := strings.Split(e.Path, "/")
		cwl = strings.Join(uriSegments[:len(uriSegments)-1], "/") + "/"

	case CONTENT_REF_FS:

		// relative .Path ?
		if false == filepath.IsAbs(e.Path) {
			e.Path = filepath.Join(e.Workingpath, e.Path)
		}

		content, err = ioutil.ReadFile(e.Path)
		if err != nil {
			return content, cwl, fmt.Errorf(`Error while reading "%s" [%v]`, e.Path, err)
		}
		cwl = filepath.Dir(e.Path)
	}

	// find ${FOO:default value} and replace with
	// var["FOO"] if found
	// environnement variaable FOO if env variable exists
	// default value, empty when not provided
	contentString := string(content)
	r, _ := regexp.Compile(`\${([a-zA-Z_\-0-9]+):?([^"'}]*)}`)
	envVars := r.FindAllStringSubmatch(contentString, -1)
	for _, envVar := range envVars {
		varText := envVar[0]
		varName := envVar[1]
		varDefaultValue := envVar[2]

		if values, ok := options["var"]; ok {
			if value, ok := values.(map[string]interface{})[varName]; ok {
				contentString = strings.Replace(contentString, varText, value.(string), -1)
				continue
			}
		}
		// Lookup for env
		if value, found := os.LookupEnv(varName); found {
			contentString = strings.Replace(contentString, varText, value, -1)
			continue
		}
		// Set default value
		contentString = strings.Replace(contentString, varText, varDefaultValue, -1)
		continue
	}
	content = []byte(contentString)

	return content, cwl, err

}

func expandFilePath(path string) ([]string, error) {
	locs := []string{}
	if fi, err := os.Stat(path); err == nil {

		if false == fi.IsDir() {
			locs = append(locs, path)
			return locs, nil
		}
		files, err := filepath.Glob(filepath.Join(path, "*.*"))
		if err != nil {
			return locs, err

		}
		//use each file
		for _, file := range files {
			switch strings.ToLower(filepath.Ext(file)) {
			case ".conf":
				locs = append(locs, file)
				continue
			default:

			}
		}
	} else {
		return locs, fmt.Errorf("%s not found", path)
	}
	return locs, nil
}
