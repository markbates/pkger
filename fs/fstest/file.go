package fstest

import (
	"github.com/markbates/pkger/fs"
)

type TestFile struct {
	Name string
	Path fs.Path
}

type TestFiles map[string]TestFile
