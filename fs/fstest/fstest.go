package fstest

import (
	"path"
	"strings"

	"github.com/markbates/pkger/fs"
)

func Files(fx fs.FileSystem) (TestFiles, error) {
	tf := TestFiles{}
	for _, f := range fileList {
		pt, err := Path(fx, f)
		if err != nil {
			return tf, err
		}
		tf[pt] = TestFile{
			Name: pt.Name,
			Path: pt,
		}
	}

	return tf, nil
}

func Path(fx fs.FileSystem, ps ...string) (fs.Path, error) {
	name := path.Join(ps...)
	name = path.Join(".fstest", name)
	if !strings.HasPrefix(name, "/") {
		name = "/" + name
	}

	return fx.Parse(name)
}

var fileList = []string{
	"/main.go",
	"/go.mod",
	"/go.sum",
	"/public/index.html",
	"/public/images/mark.png",
	"/templates/a.txt",
	"/templates/b/b.txt",
}
