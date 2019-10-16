package api

import (
	"io/ioutil"

	"github.com/markbates/pkger"
)

func Version() (string, error) {
	f, err := pkger.Open("github.com/markbates/pkger/examples/complex/api:/version.txt")
	if err != nil {
		return "", err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	return string(b), err
}
