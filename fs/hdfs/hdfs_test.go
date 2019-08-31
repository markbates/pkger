package hdfs

import (
	"log"

	"github.com/markbates/pkger/fs/fstest"
)

func NewFS() *FS {
	fs, err := New()
	if err != nil {
		log.Fatal(err)
	}
	return fs
}

var Folder = fstest.TestFiles{
	"/testdata/hdfs_test/main.go":                {Data: []byte("!/testdata/hdfs_test/main.go")},
	"/testdata/hdfs_test/public/index.html":      {Data: []byte("!/testdata/hdfs_test/public/index.html")},
	"/testdata/hdfs_test/public/images/mark.png": {Data: []byte("!/testdata/hdfs_test/public/images/mark.png")},
	"/testdata/hdfs_test/templates/a.txt":        {Data: []byte("!/testdata/hdfs_test/templates/a.txt")},
	"/testdata/hdfs_test/templates/b/b.txt":      {Data: []byte("!/testdata/hdfs_test/templates/b/b.txt")},
}
