package parser

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Parser(t *testing.T) {
	r := require.New(t)

	res, err := Parse("/internal/app")

	r.NoError(err)

	exp := []string{"github.com/gobuffalo/buffalo:/logo.svg", "github.com/markbates/pkger:/internal/app/public", "github.com/markbates/pkger:/internal/app/templates/a.txt", "github.com/markbates/pkger:/internal/app/public/images/mark-small.png", "github.com/markbates/pkger:/internal/app/public/images/mark.png", "github.com/markbates/pkger:/internal/app/public/images/mark_250px.png", "github.com/markbates/pkger:/internal/app/public/images/mark_400px.png", "github.com/markbates/pkger:/internal/app/public/index.html"}
	sort.Strings(exp)

	act := make([]string, len(res.Paths))
	for i := 0; i < len(res.Paths); i++ {
		act[i] = res.Paths[i].String()
	}

	sort.Strings(act)
	r.Equal(exp, act)
}
