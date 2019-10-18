package parser

import (
	"sort"
	"testing"

	"github.com/markbates/pkger/pkging/pkgtest"
	"github.com/stretchr/testify/require"
)

func Test_Parser_App(t *testing.T) {
	r := require.New(t)

	app, err := pkgtest.App()
	r.NoError(err)

	res, err := Parse(app.Info)

	r.NoError(err)

	act := make([]string, len(res))
	for i := 0; i < len(res); i++ {
		act[i] = res[i].String()
	}

	sort.Strings(act)
	r.Equal(app.Paths.Parser, act)
}
