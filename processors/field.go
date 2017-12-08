package processors

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/clbanning/mxj"
	"github.com/vjeantet/jodaTime"
)

const (
	mark      = "%{"
	mark_time = "%{+"
)

var maskPattern *regexp.Regexp
var maskTimePattern *regexp.Regexp

func init() {
	maskTimePattern, _ = regexp.Compile(`%\{\+([^\}]+)\}`)
	maskPattern, _ = regexp.Compile(`%\{\[?([^\}]*[^\]])\]?\}`)
}

// Dynamic includes field value in place of %{key.path}
// When no field is not found replace with ""
func Dynamic(str *string, fields *mxj.Map) {
	if true == strings.Contains(*str, mark_time) {
		// Search for all %{+word}
		for _, values := range maskTimePattern.FindAllStringSubmatch(*str, -1) {
			// Search matching value, when not found use ""
			t, err := fields.ValueForPath("@timestamp")
			if err != nil {
				continue
			}
			switch tt := t.(type) {
			case time.Time:
				*str = strings.Replace(*str, values[0], jodaTime.Format(values[1], tt), -1)
			}

		}
	}
	if true == strings.Contains(*str, mark) {
		// Search for all %{word}
		for _, values := range maskPattern.FindAllStringSubmatch(*str, -1) {
			values[1] = strings.Replace(values[1], `][`, `.`, -1)
			// Search matching value, when not found use ""
			i, _ := fields.ValueForPath(values[1])
			if i == nil {
				i = ""
			}
			*str = strings.Replace(*str, values[0], fmt.Sprintf("%v", i), -1)
		}
	}
}

func SetType(typevalue string, data *mxj.Map) {
	Dynamic(&typevalue, data)
	if !data.Exists("type") {
		data.SetValueForPath(typevalue, "type")
	}
}

func AddFields(fields map[string]interface{}, data *mxj.Map) {
	for k, v := range fields {
		Dynamic(&k, data)
		if !data.Exists(k) {
			switch v.(type) {
			case string:
				d := v.(string)
				Dynamic(&d, data)
				v = d
			}

			data.SetValueForPath(v, k)
		}
	}
}

func AddTags(tags []string, data *mxj.Map) {
	var currentTags []string

	currentTagsInterface, _ := data.ValueForPath("tags")

	switch v := currentTagsInterface.(type) {
	case string:
		currentTags = []string{v}
	case []string:
		currentTags = currentTagsInterface.([]string)
	default:
		currentTags = []string{}
	}

	tagsEval := []string{}
	for _, t := range tags {
		Dynamic(&t, data)
		if !isInSlice(t, currentTags) {
			tagsEval = append(tagsEval, t)
		}
	}

	newTags := append(currentTags, tagsEval...)
	data.SetValueForPath(newTags, "tags")
}
func isInSlice(needle string, candidates []string) bool {
	for _, symbolType := range candidates {
		if needle == symbolType {
			return true
		}
	}
	return false
}
func RemoveTags(tags []string, data *mxj.Map) {
	currenttags, err := data.ValueForPath("tags")
	if err != nil {
		return
	}

	ct := currenttags.([]string)
	for i, t := range ct {
		for _, ttodelete := range tags {
			Dynamic(&ttodelete, data)
			if ttodelete == t {
				//delete
				ct = append(ct[:i], ct[i+1:]...)
			}
		}

	}

	data.SetValueForPath(ct, "tags")
}

func RemoveFields(fields []string, data *mxj.Map) {
	for _, k := range fields {
		Dynamic(&k, data)
		data.Remove(k)
	}
}
