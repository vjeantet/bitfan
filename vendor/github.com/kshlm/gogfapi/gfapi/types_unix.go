// +build ignore

/*
To prevent structure alignment issues, cgo -godefs is designed to explicitly
convert C structs into their equivalent go definitions.
Refer: src/syscalls in go source.

This file is given as input to cgo -godefs. Example:

# export GOOS=linux
# export GOARCH=amd64
# go tool cgo -godefs types_unix.go | gofmt > ztypes_${GOOS}_${GOARCH}.go
*/

package gfapi

/*
#include <sys/statvfs.h>
*/
import "C"

type Statvfs_t C.struct_statvfs
