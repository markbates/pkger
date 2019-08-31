package memfs

import (
	"log"

	"github.com/markbates/pkger/fs/fstest"
	"github.com/markbates/pkger/here"
)

func NewFS() *FS {
	fs, err := New(here.Info{})
	if err != nil {
		log.Fatal(err)
	}
	return fs
}

var Folder = fstest.TestFiles{
	"/main.go":                {Data: []byte("!/main.go")},
	"/go.mod":                 {Data: []byte("!/go.mod")},
	"/go.sum":                 {Data: []byte("!/go.sum")},
	"/public/index.html":      {Data: []byte("!/public/index.html")},
	"/public/images/mark.png": {Data: []byte("!/public/images/mark.png")},
	"/templates/a.txt":        {Data: []byte("!/templates/a.txt")},
	"/templates/b/b.txt":      {Data: []byte("!/templates/b/b.txt")},
}
