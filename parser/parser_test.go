package parser

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/markbates/pkger/here"
	"github.com/stretchr/testify/require"
)

func Test_Parser(t *testing.T) {
	r := require.New(t)

	pwd, err := os.Getwd()
	r.NoError(err)

	ch := filepath.Join(pwd, "..",
		"examples",
		"complex")
	info := here.Info{
		Dir:        ch,
		ImportPath: "github.com/markbates/pkger/examples/complex",
	}

	res, err := Parse(info)

	r.NoError(err)

	exp := []string{
		"github.com/markbates/pkger/examples/complex:/",
		"github.com/markbates/pkger/examples/complex:/go.mod",
		// "github.com/markbates/pkger/examples/app:/public",
		// "github.com/markbates/pkger/examples/app:/public/images/mark-small.png",
		// "github.com/markbates/pkger/examples/app:/public/images/mark.png",
		// "github.com/markbates/pkger/examples/app:/public/images/mark_250px.png",
		// "github.com/markbates/pkger/examples/app:/public/images/mark_400px.png",
		// "github.com/markbates/pkger/examples/app:/public/index.html",
	}
	sort.Strings(exp)

	act := make([]string, len(res))
	for i := 0; i < len(res); i++ {
		act[i] = res[i].String()
	}

	sort.Strings(act)
	// fmt.Printf("%#v\n", act)
	r.Equal(exp, act)
}
