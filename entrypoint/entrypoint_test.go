package entrypoint

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/awillis/bitfan/core"
	"github.com/stretchr/testify/assert"
)

var tsURL = ""

func TestMain(m *testing.M) {
	ts := httptest.NewServer(http.FileServer(http.Dir("testdata")))
	defer ts.Close()
	tsURL = ts.URL
	m.Run()
}

func TestNewInline(t *testing.T) {
	e, err := New("input{} filter{} output{}", "", CONTENT_INLINE)
	assert.NoError(t, err)
	assert.IsType(t, &Entrypoint{}, e)
}

func TestNewFSRelativePathContentRef(t *testing.T) {
	e, err := New("testdata/001.conf", "", CONTENT_REF)
	assert.NoError(t, err)
	assert.IsType(t, &Entrypoint{}, e)
	assert.Equal(t, e.Kind, CONTENT_REF_FS)
}
func TestNewFSRelativePathContentRefFS(t *testing.T) {
	e, err := New("testdata/001.conf", "", CONTENT_REF_FS)
	assert.NoError(t, err)
	assert.IsType(t, &Entrypoint{}, e)
	assert.Equal(t, e.Kind, CONTENT_REF_FS)
}

func TestNewFSWithWorkingDir(t *testing.T) {
	e, err := New("001.conf", "parser/logstash/testdata", CONTENT_REF)
	assert.NoError(t, err)
	assert.IsType(t, &Entrypoint{}, e)
	assert.Equal(t, "parser/logstash/testdata", e.Workingpath)
	assert.Equal(t, "001.conf", e.Path)
	assert.Equal(t, e.Kind, CONTENT_REF_FS)
}

func TestNewFSUnknowWorkingDirPathError(t *testing.T) {
	e, err := New("001.conf", "parser/logstash/unknow", CONTENT_REF)
	assert.Error(t, err)
	assert.Nil(t, e)
}

func TestNewFSUnknowPathError(t *testing.T) {
	e, err := New("testdata/001-does-not-exists.conf", "", CONTENT_REF)
	assert.Error(t, err)
	assert.Nil(t, e)
}

func TestNewURLUnknowPathError(t *testing.T) {
	e, err := New("http://127.0.0.1/nowhere.conf", "", CONTENT_REF)
	assert.NoError(t, err)
	assert.Equal(t, e.Kind, CONTENT_REF_URL)
}
func TestNewURLWithWorkingURL(t *testing.T) {
	e, err := New("001.conf", "https://bitfan.io/conf/", CONTENT_REF)
	assert.NoError(t, err)
	assert.IsType(t, &Entrypoint{}, e)
	assert.Equal(t, "https://bitfan.io/conf/", e.Workingpath)
	assert.Equal(t, "https://bitfan.io/conf/001.conf", e.Path)
	assert.Equal(t, e.Kind, CONTENT_REF_URL)
}

func TestPipelineFS(t *testing.T) {
	e, err := New("001.conf", "testdata", CONTENT_REF)
	assert.NoError(t, err)
	assert.IsType(t, &Entrypoint{}, e)
	assert.Equal(t, e.Kind, CONTENT_REF_FS)

	pipeline, err := e.Pipeline()

	assert.NoError(t, err)
	assert.IsType(t, &core.Pipeline{}, pipeline)
	assert.Regexp(t, ".*/testdata/001.conf", pipeline.ConfigLocation)
	assert.Equal(t, "001", pipeline.Label)
	assert.Equal(t, "", pipeline.Description)
	assert.NotEqual(t, "", pipeline.Uuid)
	assert.Equal(t, 5, len(pipeline.Agents()))
}

func TestPipelineFixedUUID(t *testing.T) {
	e, err := New("001.conf", "testdata", CONTENT_REF)
	e.PipelineUuid = "1234"
	e.PipelineName = "Label1234"

	pipeline, err := e.Pipeline()
	assert.NoError(t, err)
	assert.IsType(t, &Entrypoint{}, e)
	assert.Equal(t, "Label1234", pipeline.Label)
	assert.Equal(t, "", pipeline.Description)
	assert.Equal(t, "1234", pipeline.Uuid)
}

func TestPipelineINLINE(t *testing.T) {
	content := `input{stdin { }} 
	filter{
		date {
   		 match => [ "timestamp" , "dd/MMM/yyyy:HH:mm:ss Z" ]
  		}
	} 
	output{stdout { }}`
	e, err := New(content, "", CONTENT_INLINE)

	pipeline, err := e.Pipeline()
	assert.NoError(t, err)
	assert.IsType(t, &Entrypoint{}, e)
	assert.Equal(t, "inline", pipeline.Label)
	assert.Equal(t, "inline", pipeline.ConfigLocation)
	assert.Equal(t, 3, len(pipeline.Agents()))
}

func TestPipelineURL(t *testing.T) {
	e, err := New(tsURL+"/001.conf", "", CONTENT_REF)

	pipeline, err := e.Pipeline()
	assert.NoError(t, err)
	assert.IsType(t, &Entrypoint{}, e)
	assert.Regexp(t, "127.0.0.1:[0-9]+.001.conf", pipeline.Label)
	assert.Equal(t, "", pipeline.Description)
	assert.Equal(t, 5, len(pipeline.Agents()))
}

func TestEntryPointExpand(t *testing.T) {
	e, err := New("testdata/use/subs", "", CONTENT_REF)
	assert.NoError(t, err)

	es := &EntrypointList{}
	err = es.AddEntrypoint(e)
	assert.NoError(t, err)

	assert.Equal(t, 4, len(es.Items))
}

func TestEntryPointExpandWithOne(t *testing.T) {
	e, err := New("testdata", "", CONTENT_REF)
	assert.NoError(t, err)

	es := &EntrypointList{}
	err = es.AddEntrypoint(e)
	assert.NoError(t, err)

	assert.Equal(t, 1, len(es.Items))
}

func TestEntryPointExpandNoConf(t *testing.T) {
	e, err := New("testdata/nothing", "", CONTENT_REF)
	es := &EntrypointList{}
	err = es.AddEntrypoint(e)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(es.Items))
}

func TestAddEntryPointURL(t *testing.T) {
	e, err := New(tsURL+"/subs", "", CONTENT_REF)
	assert.NoError(t, err)

	es := &EntrypointList{}
	err = es.AddEntrypoint(e)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(es.Items))
}

func TestContentURL(t *testing.T) {
	e, err := New(tsURL+"/nothing/test.txt", "", CONTENT_REF)
	assert.NoError(t, err)

	options := map[string]interface{}{
		"hello": "world",
	}

	contentBytes, ewl, err := e.content(options)
	assert.NoError(t, err)
	assert.Equal(t, 4, len(contentBytes))
	assert.Equal(t, tsURL+"/nothing/", ewl)
}

func TestContentWithOptions(t *testing.T) {
	e, err := New("testdata/vars/002.conf", "testdata/vars", CONTENT_REF)
	assert.NoError(t, err)

	options := map[string]interface{}{
		"var": map[string]interface{}{
			"myhost": "lacol123",
			"var2":   "lorem",
			"number": 1231,
		},
	}

	contentBytes, ewl, err := e.content(options)
	assert.NoError(t, err)
	assert.Regexp(t, `.*\["lacol123:9200"\].*`, string(contentBytes))
	assert.Regexp(t, ".*/testdata/vars", ewl)
}

func TestContentWithOptionsDefault(t *testing.T) {
	e, err := New("testdata/vars/002.conf", "testdata/vars", CONTENT_REF)
	assert.NoError(t, err)

	options := map[string]interface{}{
		"var": map[string]interface{}{
			"var2":   "lorem",
			"number": 1231,
		},
	}

	contentBytes, _, err := e.content(options)
	assert.NoError(t, err)
	assert.Regexp(t, `.*\["local:9200"\].*`, string(contentBytes))
}

func TestContentWithOptionsEnv(t *testing.T) {
	e, err := New("testdata/vars/002.conf", "testdata/vars", CONTENT_REF)
	assert.NoError(t, err)

	options := map[string]interface{}{}

	os.Setenv("myhost", "coucou")
	contentBytes, _, err := e.content(options)
	assert.NoError(t, err)
	assert.Regexp(t, `.*\["coucou:9200"\].*`, string(contentBytes))
}

func TestPipelineWithUse(t *testing.T) {
	e, err := New("testdata/use/main.conf", "", CONTENT_REF)
	assert.NoError(t, err)
	pipeline, err := e.Pipeline()

	assert.NoError(t, err)
	assert.Equal(t, 15, len(pipeline.Agents()))
}

func TestPipelineWithURLUse(t *testing.T) {
	e, err := New(tsURL+"/use/main.conf", "", CONTENT_REF)
	assert.NoError(t, err)
	pipeline, err := e.Pipeline()

	assert.NoError(t, err)
	assert.Equal(t, 15, len(pipeline.Agents()))
}
