package mod

import (
	"io/ioutil"

	"github.com/markbates/pkger"
)

func Mod() (string, error) {
	f, err := pkger.Open("/go.mod")
	if err != nil {
		return "", nil
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return "", nil
	}
	return string(b), nil
}
