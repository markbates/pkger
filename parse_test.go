package pkger

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const curPkg = "github.com/markbates/pkger"

func Test_Parse_Happy(t *testing.T) {
	table := []struct {
		in  string
		out Path
	}{
		{in: curPkg, out: Path{Pkg: curPkg, Name: "/"}},
		{in: curPkg + ":/foo.go", out: Path{Pkg: curPkg, Name: "/foo.go"}},
		{in: "/paths/parse_test.go", out: Path{Pkg: curPkg, Name: "/paths/parse_test.go"}},
		{in: "", out: Path{Pkg: curPkg, Name: "/"}},
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
