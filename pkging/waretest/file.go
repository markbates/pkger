package waretest

import (
	"github.com/markbates/pkger/pkging"
)

type TestFile struct {
	Name string
	Path pkging.Path
}

type TestFiles map[pkging.Path]TestFile
