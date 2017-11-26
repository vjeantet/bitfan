package parser

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var tsURL = ""

func TestMain(m *testing.M) {
	ts := httptest.NewServer(http.FileServer(http.Dir("testdata")))
	defer ts.Close()
	tsURL = ts.URL
	m.Run()
}

func TestBuildAgentsSimpleURL(t *testing.T) {
	response, _ := http.Get(tsURL + "/002.conf")
	defer response.Body.Close()
	responseData, _ := ioutil.ReadAll(response.Body)

	ewl := tsURL + "/"

	agents, err := BuildAgents(responseData, ewl, nil)

	assert.NoError(t, err)
	assert.Equal(t, 31, len(agents))
}
func TestBuildAgentsComplex(t *testing.T) {
	f, err := os.Open("testdata/002.conf")
	defer f.Close()
	responseData, _ := ioutil.ReadAll(f)

	ewl := "./testdata"

	agents, err := BuildAgents(responseData, ewl, nil)

	assert.NoError(t, err)
	assert.Equal(t, 31, len(agents))
}
func TestBuildAgentsComplexURL(t *testing.T) {
	response, _ := http.Get(tsURL + "/002.conf")
	defer response.Body.Close()
	responseData, _ := ioutil.ReadAll(response.Body)

	ewl := tsURL + "/"

	agents, err := BuildAgents(responseData, ewl, nil)

	assert.NoError(t, err)
	assert.Equal(t, 31, len(agents))
}

func TestBuildAgentsUse(t *testing.T) {
	f, err := os.Open("testdata/use/main.conf")
	ewl, _ := filepath.Abs(filepath.Dir(f.Name()))
	defer f.Close()
	responseData, _ := ioutil.ReadAll(f)

	agents, err := BuildAgents(responseData, ewl, entrypointContentFS)
	// pp.Println(agents, err)

	assert.NoError(t, err)
	assert.Equal(t, 35, len(agents))
}

func entrypointContentFS(path string, cwl string, options map[string]interface{}) ([]byte, string, error) {
	var content []byte
	var ewl string
	var err error
	var pathFix string
	var absFix string

	absFix, _ = filepath.Abs(cwl)
	pathFix = path
	if filepath.IsAbs(path) == false {
		pathFix = filepath.Join(absFix, path)
	}

	f, err := os.Open(pathFix)
	defer f.Close()
	content, _ = ioutil.ReadAll(f)

	ewl = filepath.Dir(pathFix)

	return content, ewl, err
}
