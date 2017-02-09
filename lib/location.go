package lib

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
)

const (
	CONTENT_FS = iota + 1
	CONTENT_URL
	CONTENT_INLINE
)

type Location struct {
	Path        string
	Kind        int
	Workingpath string
}

type Locations struct {
	Items []*Location
}

func (l *Locations) Add(ref string, cwl string) error {
	loc := &Location{}
	if v, _ := url.Parse(ref); v.Scheme == "http" || v.Scheme == "https" {
		loc.Kind = CONTENT_URL
		loc.Path = ref
	} else if _, err := os.Stat(ref); err == nil {
		loc.Kind = CONTENT_FS
		loc.Path, err = filepath.Abs(ref)
		if err != nil {
			return err
		}
	} else if _, err := os.Stat(filepath.Join(cwl, ref)); err == nil {
		loc.Kind = CONTENT_FS
		loc.Path = filepath.Join(cwl, ref)
	} else if v, _ := url.Parse(cwl); v.Scheme == "http" || v.Scheme == "https" {
		loc.Kind = CONTENT_URL
		loc.Path = cwl + ref
	} else {
		loc.Kind = CONTENT_INLINE
		loc.Path = ref

		// return fmt.Errorf("unknow location %s -- current working location is %s", ref, cwl)
	}

	loc.Workingpath = cwl

	// if it's a file try to expand
	if loc.Kind == CONTENT_FS {
		subpaths, err := expandFilePath(loc.Path)
		if err != nil {
			return err
		}

		for _, subpath := range subpaths {
			subloc := &Location{
				Path:        subpath,
				Workingpath: loc.Workingpath,
				Kind:        loc.Kind,
			}
			l.Items = append(l.Items, subloc)
		}
	} else {
		l.Items = append(l.Items, loc)
	}

	return nil
}

func (l *Location) ConfigPipeline() config.Pipeline {
	var pipeline *config.Pipeline

	switch l.Kind {
	case CONTENT_INLINE:
		pipeline = config.NewPipeline("inline", "nodescription", "stdin")
	case CONTENT_URL:
		uriSegments := strings.Split(l.Path, "/")
		pipelineName := strings.Join(uriSegments[2:], ".")
		pipeline = config.NewPipeline(pipelineName, "nodescription", l.Path)
	case CONTENT_FS:
		filename := filepath.Base(l.Path)
		extension := filepath.Ext(filename)
		pipelineName := filename[0 : len(filename)-len(extension)]
		pipeline = config.NewPipeline(pipelineName, "nodescription", l.Path)
	}

	return *pipeline
}

func (l *Location) ConfigAgents() ([]config.Agent, error) {
	return l.configAgentsWithOptions(map[string]interface{}{}, "input", "filter", "output")
}

func (l *Location) configAgentsWithOptions(options map[string]interface{}, pickSections ...string) ([]config.Agent, error) {
	var agents []config.Agent
	var content []byte
	var err error
	var cwd string
	content, cwd, err = l.content(options)
	if err != nil {
		return agents, err
	}

	agents, err = buildAgents(content, cwd, pickSections...)
	return agents, err
}

func (l *Location) content(options map[string]interface{}) ([]byte, string, error) {
	var content []byte
	var cwl string
	var err error

	switch l.Kind {
	case CONTENT_INLINE:
		content = []byte(l.Path)
		cwl = l.Workingpath

	case CONTENT_URL:
		response, err := http.Get(l.Path)
		if err != nil {
			return content, cwl, err
		} else {
			content, err = ioutil.ReadAll(response.Body)
			response.Body.Close()
			if err != nil {
				return content, cwl, err
			}
		}

		uriSegments := strings.Split(l.Path, "/")
		cwl = strings.Join(uriSegments[:len(uriSegments)-1], "/") + "/"

	case CONTENT_FS:

		// si location est relatif
		if false == filepath.IsAbs(l.Path) {
			l.Path = filepath.Join(l.Workingpath, l.Path)
		}

		content, err = ioutil.ReadFile(l.Path)
		if err != nil {
			return content, cwl, fmt.Errorf(`Error while reading "%s" [%s]`, l.Path, err)
		}
		cwl = filepath.Dir(l.Path)
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
			return locs, fmt.Errorf("error %s", err.Error())

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
