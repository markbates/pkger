package pkgutil

import (
	"io/ioutil"

	"github.com/markbates/pkger/pkging"
)

// ReadFile reads the file named by filename and returns the contents. A successful call returns err == nil, not err == EOF. Because ReadFile reads the whole file, it does not treat an EOF from Read as an error to be reported.
func ReadFile(pkg pkging.Pkger, s string) ([]byte, error) {
	f, err := pkg.Open(s)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}
