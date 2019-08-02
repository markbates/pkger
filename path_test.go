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
	}

	for _, tt := range table {
		t.Run(tt.in.String(), func(st *testing.T) {
			r := require.New(st)
			r.Equal(tt.out, tt.in.String())
		})
	}
}
