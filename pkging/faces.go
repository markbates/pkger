package pkging

import (
	"io"

	"github.com/markbates/pkger/here"
)

type Adder interface {
	Add(files ...File) error
}

type Stuffer interface {
	Stuff(w io.Writer, paths []here.Path) error
}
