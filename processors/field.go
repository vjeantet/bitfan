package processors

import (
	"fmt"
	"regexp"
	"strconv"
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
	maskPattern, _ = regexp.Compile(`%\{\[?([^\}]+[^\]])\]?\}`)
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
			*str = strings.Replace(*str, values[0], jodaTime.Format(values[1], t.(time.Time)), -1)
		}
	}
	if true == strings.Contains(*str, mark) {
		// Search for all %{word}
		for _, values := range maskPattern.FindAllStringSubmatch(*str, -1) {
			values[1] = strings.Replace(values[1], `][`, `.`, -1)
			// Search matching value, when not found use ""
			*str = strings.Replace(*str, values[0], fields.ValueOrEmptyForPathString(values[1]), -1)
		}
	}
}

func ProcessCommonFields(data *mxj.Map, add_fields map[string]interface{}, tags []string, typevalue string) {
	AddFields(add_fields, data)
	if typevalue != "" {
		SetType(typevalue, data)
	}
	if len(tags) > 0 {
		AddTags(tags, data)
	}
}

func ProcessCommonFields2(data *mxj.Map,
	add_fields map[string]interface{},
	tags []string,
	remove_field []string,
	remove_tag []string) {

	if len(add_fields) > 0 {
		AddFields(add_fields, data)
	}
	if len(tags) > 0 {
		AddTags(tags, data)
	}
	if len(remove_field) > 0 {
		RemoveFields(remove_field, data)
	}
	if len(remove_tag) > 0 {
		RemoveTags(remove_tag, data)
	}
}

func SetType(typevalue string, data *mxj.Map) {
	Dynamic(&typevalue, data)
	data.SetValueForPath(typevalue, "type")
}

func AddFields(fields map[string]interface{}, data *mxj.Map) {
	for k, v := range fields {
		Dynamic(&k, data)
		switch v.(type) {
		case string:
			d := v.(string)
			Dynamic(&d, data)
			v = d
		}

		data.SetValueForPath(v, k)
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
		tagsEval = append(tagsEval, t)
	}

	newTags := append(currentTags, tagsEval...)
	data.SetValueForPath(newTags, "tags")
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

func RemoveAllButFields(fields []string, data *mxj.Map) {
	if len(fields) > 0 {
		cp := mxj.New()
		for _, k := range fields {
			Dynamic(&k, data)
			if value, err := data.ValueForPath(k); err == nil {
				cp.SetValueForPath(value, k)
			}
		}
		*data = cp
	}
}

func UpdateFields(fields map[string]interface{}, data *mxj.Map) {
	for k, v := range fields {
		if data.Exists(k) {
			data.SetValueForPath(v, k)
		}
	}
}

func RenameFields(fields map[string]string, data *mxj.Map) {
	for k, v := range fields {
		if data.Exists(k) {
			data.RenameKey(k, v)
		}
	}
}

func UpperCaseFields(fields []string, data *mxj.Map) {
	for _, k := range fields {
		if value, err := data.ValueForPathString(k); err == nil {
			data.SetValueForPath(strings.ToUpper(value), k)
		}
	}
}

func LowerCaseFields(fields []string, data *mxj.Map) {
	for _, k := range fields {
		if value, err := data.ValueForPathString(k); err == nil {
			data.SetValueForPath(strings.ToLower(value), k)
		}
	}
}

func Join(fields map[string]string, data *mxj.Map) {
	for path, glue := range fields {
		if !data.Exists(path) {
			continue
		}
		value, _ := data.ValueForPath(path)

		switch value.(type) {
		case []string:
			a := []string{}
			for _, s := range value.([]interface{}) {
				a = append(a, s.(string))
			}
			newValue := strings.Join(a, glue)
			data.SetValueForPath(newValue, path)
		}
	}
}

func Split(fields map[string]string, data *mxj.Map) {
	for path, separator := range fields {
		if !data.Exists(path) {
			continue
		}
		value := data.ValueOrEmptyForPathString(path)
		newValue := strings.Split(value, separator)
		data.SetValueForPath(newValue, path)
	}
}

func Strip(fields []string, data *mxj.Map) {
	for _, path := range fields {
		if value, err := data.ValueForPathString(path); err == nil {
			newValue := strings.TrimSpace(value)
			data.SetValueForPath(newValue, path)
		}

	}
}

func Gsub(fields []string, data *mxj.Map) {
	for i := 0; i < len(fields); i++ {
		fieldname := fields[i]
		i++
		pattern := fields[i]
		i++
		replacement := fields[i]

		if value, err := data.ValueForPathString(fieldname); err == nil {
			r, _ := regexp.Compile(pattern)
			newValue := r.ReplaceAllString(value, replacement)
			data.SetValueForPath(newValue, fieldname)
		}

	}
}

func Convert(fields map[string]string, data *mxj.Map) {

	for path, kind := range fields {
		if !data.Exists(path) {
			continue
		}

		value, err := data.ValueForPath(path)
		if err != nil {
			continue
		}

		switch value.(type) {
		case string:
			switch kind {
			case "integer":
				newValue, err := strconv.Atoi(value.(string))
				if err != nil {
					continue
				}
				data.SetValueForPath(newValue, path)
			case "float":
				newValue, err := strconv.ParseFloat(value.(string), 64)
				if err != nil {
					continue
				}
				data.SetValueForPath(newValue, path)
			case "boolean":
				newValue := false
				value = strings.ToLower(value.(string))
				for _, b := range []string{"true", "t", "yes", "y", "1"} {
					if b == value {
						newValue = true
					}
				}
				data.SetValueForPath(newValue, path)
			}
		case int:
			switch kind {
			case "string":
				newValue := fmt.Sprintf("%d", value.(int))
				data.SetValueForPath(newValue, path)
			case "float":
				newValue := float64(value.(int))
				data.SetValueForPath(newValue, path)
			case "boolean":
				newValue := false
				if value.(int) > 0 {
					newValue = true
				}
				data.SetValueForPath(newValue, path)
			}
		case float64:
			switch kind {
			case "string":
				newValue := fmt.Sprintf("%f", value.(float64))
				data.SetValueForPath(newValue, path)
			case "integer":
				newValue := int(value.(float64))
				data.SetValueForPath(newValue, path)
			case "boolean":
				newValue := false
				if value.(float64) > 0 {
					newValue = true
				}
				data.SetValueForPath(newValue, path)
			}
		case bool:
			switch kind {
			case "string":
				newValue := "false"
				if value.(bool) == true {
					newValue = "true"
				}
				data.SetValueForPath(newValue, path)
			case "integer":
				var newValue int
				newValue = 0
				if value.(bool) == true {
					newValue = 1
				}
				data.SetValueForPath(newValue, path)
			case "float":
				var newValue float64
				newValue = 0
				if value.(bool) == true {
					newValue = 1
				}
				data.SetValueForPath(newValue, path)
			}
		}

	}

}

func Merge(fields map[string]string, data *mxj.Map) {
	for path_dst, path_src := range fields {
		if !data.Exists(path_dst) || !data.Exists(path_src) {
			continue
		}
		value_src, _ := data.ValueForPath(path_src)
		value_dst, _ := data.ValueForPath(path_dst)

		a := []string{}
		b := []string{}

		// newValue := append(value_dst, value_src...)
		switch value_src.(type) {
		case []string:

			for _, s := range value_src.([]string) {
				a = append(a, s)
			}
		default:
			continue
		}

		// newValue := append(value_dst, value_src...)
		switch value_dst.(type) {
		case []string:

			for _, s := range value_dst.([]string) {
				b = append(b, s)
			}
		default:
			continue
		}

		a = append(b, a...)

		//Remove duplicates
		result := []string{}
		seen := map[string]string{}
		for _, val := range a {
			if _, ok := seen[val]; !ok {
				result = append(result, val)
				seen[val] = val
			}
		}

		data.SetValueForPath(result, path_dst)
	}
}
