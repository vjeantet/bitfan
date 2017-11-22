#!/usr/bin/env bash
go list -f '

{{$dir := ""}}
{{range $imp := .Deps}}
{{printf "%s %s\n" $imp $dir}}
{{end}}' $1 | sort | uniq | grep "\." 