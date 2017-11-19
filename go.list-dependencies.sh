#!/usr/bin/env bash
go list -f '

{{$dir := ""}}
{{range $imp := .Deps}}
{{printf "%s %s\n" $imp $dir}}
{{end}}' ./... | sort | uniq | grep "\." 