package pkger

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Path_String(t *testing.T) {
	table := []struct {
		in  Path
		out string
	}{
		{in: Path{}, out: ":/"},
		{in: Path{Pkg: curPkg}, out: curPkg + ":/"},
		{in: Path{Pkg: curPkg, Name: "/foo.go"}, out: curPkg + ":/foo.go"},
		{in: Path{Name: "/foo.go"}, out: ":/foo.go"},
		{in: Path{Pkg: "github.com/markbates/pkger/internal/examples/app"}, out: "github.com/markbates/pkger/internal/examples/app:/"},
	}

	for _, tt := range table {
		t.Run(tt.in.String(), func(st *testing.T) {
			r := require.New(st)
			r.Equal(tt.out, tt.in.String())
		})
	}
}
func Test_Parse(t *testing.T) {
	table := []struct {
		in  string
		out Path
	}{
		{in: curPkg, out: Path{Pkg: curPkg, Name: "/"}},
		{in: curPkg + ":/foo.go", out: Path{Pkg: curPkg, Name: "/foo.go"}},
		{in: "/paths/parse_test.go", out: Path{Pkg: curPkg, Name: "/paths/parse_test.go"}},
		{in: `\windows\path.go`, out: Path{Pkg: curPkg, Name: "/windows/path.go"}},
		{in: "", out: Path{Pkg: curPkg, Name: "/"}},
		{in: "github.com/markbates/pkger/internal/examples/app", out: Path{Pkg: "github.com/markbates/pkger/internal/examples/app", Name: "/"}},
	}

	for _, tt := range table {
		t.Run(tt.in, func(st *testing.T) {
			r := require.New(st)

			pt, err := Parse(tt.in)
			r.NoError(err)

			r.Equal(tt.out.Pkg, pt.Pkg)
			r.Equal(tt.out.Name, pt.Name)
		})
	}
}
