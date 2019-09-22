package parser

import (
	"fmt"
	"path/filepath"
	"sort"
	"testing"

	"github.com/markbates/pkger/here"
	"github.com/stretchr/testify/require"
)

func Test_Parser(t *testing.T) {
	r := require.New(t)

	ch := filepath.Join("..",
		"examples",
		"app")
	info := here.Info{
		Dir:        ch,
		ImportPath: "github.com/markbates/pkger/examples/app",
	}
	res, err := Parse(info)

	r.NoError(err)

	exp := []string{
		"github.com/markbates/pkger/examples/app:/",
		"github.com/markbates/pkger/examples/app:/public",
		"github.com/markbates/pkger/examples/app:/public/images/mark-small.png",
		"github.com/markbates/pkger/examples/app:/public/images/mark.png",
		"github.com/markbates/pkger/examples/app:/public/images/mark_250px.png",
		"github.com/markbates/pkger/examples/app:/public/images/mark_400px.png",
		"github.com/markbates/pkger/examples/app:/public/index.html",
	}
	sort.Strings(exp)

	act := make([]string, len(res.Paths))
	for i := 0; i < len(res.Paths); i++ {
		act[i] = res.Paths[i].String()
	}

	sort.Strings(act)
	fmt.Printf("%#v\n", act)
	r.Equal(exp, act)
}
