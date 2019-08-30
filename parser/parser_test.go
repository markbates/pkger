package parser

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Parser(t *testing.T) {
	r := require.New(t)
	pwd, err := os.Getwd()
	r.NoError(err)
	r.NoError(os.Chdir(filepath.Join("..", "examples", "app")))
	defer os.Chdir(pwd)

	res, err := Parse("")

	r.NoError(err)

	exp := []string{
		"github.com/markbates/pkger/examples/app:/",
		"github.com/markbates/pkger/examples/app:/public",
		"github.com/markbates/pkger/examples/app:/templates",
		"github.com/markbates/pkger/examples/app:/templates/a.txt",
		"github.com/markbates/pkger/examples/app:/templates/b",
		"github.com/markbates/pkger/examples/app:/public/images/mark-small.png",
		"github.com/markbates/pkger/examples/app:/public/images",
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
	r.Equal(exp, act)
}
