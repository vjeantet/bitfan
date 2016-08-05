// Copyright (c) 2013, Kaushal M <kshlmster at gmail dot com>
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this
//    list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
// ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package gfapi

// This file includes some helper functions used internally by the package

import (
	"os"
	"path"
	"syscall"
	"time"
)

// posixMode() returns the posix specific  mode bits from Go's portable mode bits
//
// Copied from the syscallMode() function in file_posix.go in the Go source
func posixMode(i os.FileMode) (o uint32) {
	o |= uint32(i.Perm())
	if i&os.ModeSetuid != 0 {
		o |= syscall.S_ISUID
	}
	if i&os.ModeSetgid != 0 {
		o |= syscall.S_ISGID
	}
	if i&os.ModeSticky != 0 {
		o |= syscall.S_ISVTX
	}
	return
}

// fileInfo is an implementation of the os.FileInfo interface
//
// Based on the implementation of fileStat structure in the pkg/os/types_notwin.go file of the Go source
type fileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	sys     interface{}
}

func (fs *fileInfo) Size() int64 {
	return fs.size
}

func (fs *fileInfo) Name() string {
	return fs.name
}

func (fs *fileInfo) Mode() os.FileMode {
	return fs.mode
}

func (fs *fileInfo) ModTime() time.Time {
	return fs.modTime
}

func (fs *fileInfo) IsDir() bool {
	return fs.mode.IsDir()
}

func (fs *fileInfo) Sys() interface{} {
	return fs.sys
}

// fileInfoFromStat() returns an os.FileInfo struct from the given syscall.Stat_t struc
//
// Based on the fileInfoFromStat function in the pkg/os/stat_linux.go file in the Go source
func fileInfoFromStat(st *syscall.Stat_t, name string) os.FileInfo {
	fs := &fileInfo{
		name:    path.Base(name),
		size:    int64(st.Size),
		modTime: timespecToTime(getLastModification(st)),
		sys:     st,
	}
	fs.mode = os.FileMode(st.Mode & 0777)
	switch st.Mode & syscall.S_IFMT {
	case syscall.S_IFBLK:
		fs.mode |= os.ModeDevice
	case syscall.S_IFCHR:
		fs.mode |= os.ModeDevice | os.ModeCharDevice
	case syscall.S_IFDIR:
		fs.mode |= os.ModeDir
	case syscall.S_IFIFO:
		fs.mode |= os.ModeNamedPipe
	case syscall.S_IFLNK:
		fs.mode |= os.ModeSymlink
	case syscall.S_IFREG:
		// nothing to do
	case syscall.S_IFSOCK:
		fs.mode |= os.ModeSocket
	}
	if st.Mode&syscall.S_ISGID != 0 {
		fs.mode |= os.ModeSetgid
	}
	if st.Mode&syscall.S_ISUID != 0 {
		fs.mode |= os.ModeSetuid
	}
	if st.Mode&syscall.S_ISVTX != 0 {
		fs.mode |= os.ModeSticky
	}
	return fs
}

// timespecToTime() converts a given syscall.Timespec to time.Time
//
// Copied from pkg/os/stat_linux.go in the Go source
func timespecToTime(ts syscall.Timespec) time.Time {
	return time.Unix(int64(ts.Sec), int64(ts.Nsec))
}
