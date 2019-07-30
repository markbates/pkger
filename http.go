package pkger

import (
	"bytes"
	"io"
	"net/http"
	"os"
)

type crs interface {
	io.Closer
	io.Reader
	io.Seeker
}

type byteCRS struct {
	*bytes.Reader
}

func (byteCRS) Close() error {
	return nil
}

var _ crs = &byteCRS{}

type httpFile struct {
	File *File
	crs
}

func (h httpFile) Readdir(n int) ([]os.FileInfo, error) {
	if h.File == nil {
		return nil, os.ErrNotExist
	}
	return h.File.Readdir(n)
}

func (h httpFile) Stat() (os.FileInfo, error) {
	if h.File == nil {
		return nil, os.ErrNotExist
	}
	return h.File.Stat()
}

var _ http.File = &httpFile{}
