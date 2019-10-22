package parser_test

import (
	"fmt"
	"sort"
	"testing"

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
