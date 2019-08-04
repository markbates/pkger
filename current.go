package pkger

import (
	"github.com/markbates/pkger/here"
)

func Info(p string) (here.Info, error) {
	return rootIndex.Info(p)
}

func Stat() (here.Info, error) {
	return rootIndex.Stat()
}
