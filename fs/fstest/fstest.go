package fstest

import (
	"github.com/markbates/pkger/fs"
)

func Files(fx fs.FileSystem) (TestFiles, error) {
	tf := TestFiles{}
	for _, f := range fileList {
		pt, err := fx.Parse(f)
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

var fileList = []string{
	"/main.go",
	"/go.mod",
	"/go.sum",
	"/public/index.html",
	"/public/images/mark.png",
	"/templates/a.txt",
	"/templates/b/b.txt",
}
