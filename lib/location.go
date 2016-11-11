package lib

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/veino/veino/config"
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
				Workingpath: cwl,
				Kind:        CONTENT_FS,
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
	var agents []config.Agent
	var content []byte
	var err error
	var cwd string

	content, cwd, err = l.Content()
	if err != nil {
		return agents, err
	}

	agents, err = ParseConfig(content, cwd)
	return agents, err
}

func (l *Location) Content() ([]byte, string, error) {
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
