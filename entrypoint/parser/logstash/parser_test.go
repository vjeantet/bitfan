package logstash

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewParser(t *testing.T) {
	r := strings.NewReader("")
	p := NewParser(r)
	assert.NotNil(t, p.l)
}

func parseTestCase(confCaseName string) (*Configuration, error) {
	wd, _ := os.Getwd()
	confCaseNameFileName := filepath.Join(wd, "testdata", confCaseName+".conf")
	f, err := os.Open(confCaseNameFileName)
	if err != nil {
		log.Fatalf("missing fixture for testcase %s/%s", wd, confCaseNameFileName)
	}

	return NewParser(f).Parse()
}

func TestParseEmpty(t *testing.T) {
	r := strings.NewReader("")
	p := NewParser(r)

	conf, err := p.Parse()

	assert.NoError(t, err)
	assert.Equal(t, 0, len(conf.Sections))
	assert.Nil(t, conf.Sections["input"])
	assert.Nil(t, conf.Sections["filter"])
	assert.Nil(t, conf.Sections["output"])
}

func TestParseError(t *testing.T) {
	_, err := parseTestCase("001")

	assert.Error(t, err)
	assert.IsType(t, &ParseError{}, err)
	assert.Equal(t, 6, err.(*ParseError).Line)
	assert.Equal(t, 7, err.(*ParseError).Column)
	assert.Contains(t, err.Error(), "'outsput'")
}

func TestParseInputOnly(t *testing.T) {
	conf, err := parseTestCase("002")

	assert.NoError(t, err)
	assert.Equal(t, 1, len(conf.Sections))
	assert.NotNil(t, conf.Sections["input"])
	assert.Nil(t, conf.Sections["filter"])
	assert.Nil(t, conf.Sections["output"])
}

func TestParseFilterOnly(t *testing.T) {
	conf, err := parseTestCase("003")

	assert.NoError(t, err)
	assert.Equal(t, 1, len(conf.Sections))
	assert.Nil(t, conf.Sections["input"])
	assert.NotNil(t, conf.Sections["filter"])
	assert.Nil(t, conf.Sections["output"])
}

func TestParseOutputOnly(t *testing.T) {
	conf, err := parseTestCase("004")

	assert.NoError(t, err)
	assert.Equal(t, 1, len(conf.Sections))
	assert.Nil(t, conf.Sections["input"])
	assert.Nil(t, conf.Sections["filter"])
	assert.NotNil(t, conf.Sections["output"])
}

func TestParse1(t *testing.T) {
	_, err := parseTestCase("005")
	assert.NoError(t, err)
}

func TestParseErrorOnFirstToken(t *testing.T) {
	_, err := parseTestCase("006")

	assert.Error(t, err)
	assert.IsType(t, &ParseError{}, err)
	assert.Equal(t, 1, err.(*ParseError).Line)
	assert.Equal(t, 1, err.(*ParseError).Column)
	assert.Contains(t, err.Error(), "'{'")
}

func TestParseErrorSectionErrorToken(t *testing.T) {
	_, err := parseTestCase("007")

	assert.Error(t, err)
	assert.IsType(t, &ParseError{}, err)
	assert.Equal(t, 7, err.(*ParseError).Line)
	assert.Equal(t, 8, err.(*ParseError).Column)
	assert.Contains(t, err.Error(), "'}'")
}

func TestParseErrorSectionErrorToken2(t *testing.T) {
	_, err := parseTestCase("008")

	assert.Error(t, err)
	assert.IsType(t, &ParseError{}, err)
	assert.Equal(t, 9, err.(*ParseError).Line)
	assert.Equal(t, 0, err.(*ParseError).Column)
}

func TestParseErrorPluginError(t *testing.T) {
	_, err := parseTestCase("009")

	assert.Error(t, err)
	assert.IsType(t, &ParseError{}, err)
	assert.Equal(t, 8, err.(*ParseError).Line)
	assert.Equal(t, 8, err.(*ParseError).Column)
}

func TestParseErrorPluginIdentifierError(t *testing.T) {
	_, err := parseTestCase("010")

	assert.Error(t, err)
	assert.IsType(t, &ParseError{}, err)
	assert.Equal(t, 8, err.(*ParseError).Line)
	assert.Equal(t, 16, err.(*ParseError).Column)
}

func TestParseErrorWhenError(t *testing.T) {
	_, err := parseTestCase("011")

	assert.Error(t, err)
	assert.IsType(t, &ParseError{}, err)
	assert.Equal(t, 38, err.(*ParseError).Line)
	assert.Equal(t, 34, err.(*ParseError).Column)
}

func TestParseErrorIsNumericError(t *testing.T) {
	_, err := parseTestCase("012")

	assert.Error(t, err)

	assert.IsType(t, &ParseError{}, err)
	assert.Equal(t, 7, err.(*ParseError).Line)
	assert.Equal(t, 0, err.(*ParseError).Column)
}

func TestParseErrorAssignement(t *testing.T) {
	_, err := parseTestCase("013")

	assert.Error(t, err)

	assert.IsType(t, &ParseError{}, err)
	assert.Equal(t, 6, err.(*ParseError).Line)
	assert.Equal(t, 14, err.(*ParseError).Column)
}

func TestParseErrorElseIf(t *testing.T) {
	_, err := parseTestCase("014")

	assert.Error(t, err)

	assert.IsType(t, &ParseError{}, err)
	assert.Equal(t, 7, err.(*ParseError).Line)
	assert.Equal(t, 13, err.(*ParseError).Column)
}

func TestParseErrorFalse(t *testing.T) {
	_, err := parseTestCase("015")

	assert.Error(t, err)

	assert.IsType(t, &ParseError{}, err)
	assert.Equal(t, 12, err.(*ParseError).Line)
	assert.Equal(t, 0, err.(*ParseError).Column)
}

func TestParseErrorUnkowTokken(t *testing.T) {
	_, err := parseTestCase("016")

	assert.Error(t, err)

	assert.IsType(t, &ParseError{}, err)
	assert.Equal(t, 6, err.(*ParseError).Line)
	assert.Equal(t, 13, err.(*ParseError).Column)
}

func TestParseIssue75(t *testing.T) {
	conf, err := parseTestCase("issue-75")
	assert.Equal(t, "'message' =~ '^{'", conf.Sections["filter"].Plugins[0].When[0].Expression)
	assert.Equal(t, "'message' =~ '^{'", conf.Sections["filter"].Plugins[1].When[0].Expression)
	assert.Equal(t, "'{' in [message]", conf.Sections["filter"].Plugins[2].When[0].Expression)
	assert.NoError(t, err)
}
