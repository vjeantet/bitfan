package field

import (
	"regexp"
	"strings"

	"github.com/clbanning/mxj"
)

const (
	mark          = "%{"
	mark_variable = "["
)

var maskPattern *regexp.Regexp
var maskPattern_variable *regexp.Regexp

func init() {
	maskPattern, _ = regexp.Compile(`%{([\w\.]+)}`)
	maskPattern_variable, _ = regexp.Compile(`\[([\w\.]+)\]`)
}

// Dynamic includes field value in place of %{key.path}
// When no field is not found replace with ""
func Dynamic(str *string, fields *mxj.Map) {
	// If %{ exists in value
	if true == strings.Contains(*str, mark) {
		// Search for all %{word}
		for _, values := range maskPattern.FindAllStringSubmatch(*str, -1) {
			// Search matching value, when not found use ""
			replaceBy := fields.ValueOrEmptyForPathString(values[1])
			*str = strings.Replace(*str, values[0], replaceBy, -1)
		}
	}
}
