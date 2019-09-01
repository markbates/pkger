package fstest

import (
	"path"

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

func Path(fx fs.FileSystem, p string) (fs.Path, error) {
	pt, err := fx.Parse(p)
	if err != nil {
		return pt, err
	}
	pt.Name = path.Join("/.fstest", pt.Name)
	return pt, nil
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
