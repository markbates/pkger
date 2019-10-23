package parser_test

import (
	"fmt"
	"path/filepath"
	"sort"
	"testing"

	"github.com/markbates/pkger/here"
	"github.com/markbates/pkger/parser"
	"github.com/markbates/pkger/pkging/pkgtest"
	"github.com/stretchr/testify/require"
)

func Test_Parser_App(t *testing.T) {
	r := require.New(t)

	app, err := pkgtest.App()
	r.NoError(err)

	res, err := parser.Parse(app.Info)

	r.NoError(err)

	files, err := res.Files()
	r.NoError(err)

	act := make([]string, len(files))
	for i := 0; i < len(files); i++ {
		act[i] = files[i].Path.String()
	}

	sort.Strings(act)

	for _, a := range act {
		fmt.Println(a)
	}
	r.Equal(app.Paths.Parser, act)
}

func Test_Parse_Dynamic_Files(t *testing.T) {
	r := require.New(t)

	app, err := dynamic()
	r.NoError(err)

	res, err := parser.Parse(app.Info)

	r.NoError(err)

	files, err := res.Files()
	r.NoError(err)

	r.Len(files, 1)

	f := files[0]
	r.Equal("/go.mod", f.Path.Name)
}

// dynamic returns here.info that represents the
// ./internal/testdata/app. This should be used
// by tests.
func dynamic() (pkgtest.AppDetails, error) {
	var app pkgtest.AppDetails

	her, err := here.Package("github.com/markbates/pkger")
	if err != nil {
		return app, err
	}

	info := here.Info{
		ImportPath: "dynamic",
	}

	ch := filepath.Join(
		her.Dir,
		"pkging",
		"pkgtest",
		"internal",
		"testdata",
		"dynamic")

	info.Dir = ch

	info, err = here.Cache(info.ImportPath, func(s string) (here.Info, error) {
		return info, nil
	})
	if err != nil {
		return app, err
	}
	app.Info = info
	return app, nil
}
