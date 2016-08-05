// Created by cgo -godefs - DO NOT EDIT
// cgo -godefs types_unix.go

package gfapi

type Statvfs_t struct {
	Bsize      uint64
	Frsize     uint64
	Blocks     uint64
	Bfree      uint64
	Bavail     uint64
	Files      uint64
	Ffree      uint64
	Favail     uint64
	Fsid       uint64
	Flag       uint64
	Namemax    uint64
	X__f_spare [6]int32
}
