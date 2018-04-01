package xprocessor

import (
	"fmt"
	"strconv"
)

type Option struct {
	Name         string      `json:"name"`
	Alias        string      `json:"alias"`
	Doc          string      `json:"doc"`
	Required     bool        `json:"required"`
	Type         string      `json:"type"`
	DefaultValue interface{} `json:"default_value"`
	ExampleLS    string      `json:"example"`
	Value        interface{} `json:"-"`
}

func (o *Option) Default() []string {
	r := []string{}
	if o.DefaultValue == nil {
		return r
	}

	switch dv := o.DefaultValue.(type) {
	case string:
		r = append(r, dv)
	case bool:
		if dv {
			r = append(r, "true")
		} else {
			r = append(r, "false")
		}
	case int:
		r = append(r, strconv.Itoa(dv))
	case []string:
		r = append(r, dv...)
	case []int:
		for _, v := range dv {
			r = append(r, strconv.Itoa(v))
		}
	case map[string]string:
		for k, v := range dv {
			r = append(r, fmt.Sprintf("%s:%s", k, v))
		}
	}

	return r
}

type Options map[string]*Option

func (o Options) Value(name string) interface{} {
	if _, ok := o[name]; !ok {
		return nil
	}
	switch v := o[name].Value.(type) {
	case *bool:
		return *v
	case *int:
		return *v
	case *string:
		return *v
	case *[]string:
		return *v
	case *[]int:
		return *v
	}
	return nil
}

func (o Options) IntSlice(name string) []int {
	if _, ok := o[name]; !ok {
		return []int{}
	}
	switch v := o[name].Value.(type) {
	case *[]int:
		return *v
	}
	return []int{}
}

func (o Options) StringSlice(name string) []string {
	if _, ok := o[name]; !ok {
		return []string{}
	}
	switch v := o[name].Value.(type) {
	case *[]string:
		return *v
	}
	return []string{}
}

func (o Options) MapString(name string) map[string]string {
	if _, ok := o[name]; !ok {
		return map[string]string{}
	}
	switch v := o[name].Value.(type) {
	case *map[string]string:
		return *v
	}
	return map[string]string{}
}

func (o Options) String(name string) string {
	if _, ok := o[name]; !ok {
		return ""
	}
	switch v := o[name].Value.(type) {
	case *string:
		return *v
	}
	return ""
}

func (o Options) Int(name string) int {
	if _, ok := o[name]; !ok {
		return 0
	}
	switch v := o[name].Value.(type) {
	case *int:
		return *v
	}
	return 0
}

func (o Options) Bool(name string) bool {
	if _, ok := o[name]; !ok {
		return false
	}
	switch v := o[name].Value.(type) {
	case *bool:
		return *v
	}
	return false
}
