package parser

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Parser(t *testing.T) {
	// unsure why, but it seems the gproxy
	// has trouble handing this request.
	// 	finding github.com/markbates/pkger v0.0.0
	// 	github.com/markbates/pkger@v0.0.0: unexpected status (https://proxy.golang.org/github.com/markbates/pkger/@v/v0.0.0.info): 410 Gone
	// turning it off works just fine, so this
	// might be a bug.
	// 	module github.com/markbates/pkger/internal/examples/app
	//
	// 	go 1.12
	//
	// 	require github.com/markbates/pkger v0.0.0-20190803203656-a4a55a52dc5d
	// const env = "GOPROXY"
	// prxy := os.Getenv(env)
	// os.Setenv(env, "")
	// defer func() {
	// 	os.Setenv(env, prxy)
	// }()

	r := require.New(t)
	pwd, err := os.Getwd()
	r.NoError(err)
	r.NoError(os.Chdir(filepath.Join("..", "internal", "examples", "app")))
	defer os.Chdir(pwd)

	res, err := Parse("")

	r.NoError(err)

	exp := []string{
		"github.com/markbates/pkger/internal/examples/app:/",
		"github.com/markbates/pkger/internal/examples/app:/public",
		"github.com/markbates/pkger/internal/examples/app:/templates",
		"github.com/markbates/pkger/internal/examples/app:/templates/a.txt",
		"github.com/markbates/pkger/internal/examples/app:/templates/b",
		"github.com/markbates/pkger/internal/examples/app:/public/images/mark-small.png",
		"github.com/markbates/pkger/internal/examples/app:/public/images",
		"github.com/markbates/pkger/internal/examples/app:/public/images/mark.png",
		"github.com/markbates/pkger/internal/examples/app:/public/images/mark_250px.png",
		"github.com/markbates/pkger/internal/examples/app:/public/images/mark_400px.png",
		"github.com/markbates/pkger/internal/examples/app:/public/index.html",
	}
	sort.Strings(exp)

	act := make([]string, len(res.Paths))
	for i := 0; i < len(res.Paths); i++ {
		act[i] = res.Paths[i].String()
	}

	sort.Strings(act)
	r.Equal(exp, act)
}
