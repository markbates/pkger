package fstest

import (
	"fmt"
	"path"
	"strings"

	"github.com/markbates/pkger/fs"
)

func Files(fx fs.FileSystem) (TestFiles, error) {
	info, err := fx.Current()
	if err != nil {
		return nil, err
	}
	tf := TestFiles{}
	for _, f := range fileList {
		name := Path(fx, f)
		tf[name] = TestFile{
			Name: name,
			Path: fs.Path{
				Pkg:  info.ImportPath,
				Name: name,
			},
		}
	}

	return tf, nil
}

func Path(fx fs.FileSystem, ps ...string) string {
	name := path.Join("/.fstest", fmt.Sprintf("%T", fx))
	name = path.Join(name, strings.Join(ps, "/"))
	return name
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
