package pathparser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Parser(t *testing.T) {
	table := []struct {
		in  string
		exp Path
		err bool
	}{
		{
			in: "/a/b/c.txt",
			exp: Path{
				Name: "/a/b/c.txt",
			},
		},
		{
			in: "github.com/markbates/pkger:/a/b/c.txt",
			exp: Path{
				Pkg: &Package{
					Name: "github.com/markbates/pkger",
				},
				Name: "/a/b/c.txt"},
		},
		{
			in: "github.com/markbates/pkger@v1.0.0:/a/b/c.txt",
			exp: Path{
				Pkg: &Package{
					Name:    "github.com/markbates/pkger",
					Version: "v1.0.0",
				},
				Name: "/a/b/c.txt"},
		},
		{
			in: "github.com/markbates/pkger@v1.0.0",
			exp: Path{
				Pkg: &Package{
					Name:    "github.com/markbates/pkger",
					Version: "v1.0.0",
				},
				Name: "/",
			},
		},
		{
			in: "github.com/markbates/pkger",
			exp: Path{
				Pkg: &Package{
					Name: "github.com/markbates/pkger",
				},
				Name: "/",
			},
		},
		{
			in: "app",
			exp: Path{
				Pkg: &Package{
					Name: "app",
				},
				Name: "/",
			},
		},
	}

	for _, tt := range table {

		t.Run(tt.in, func(st *testing.T) {
			r := require.New(st)

			res, err := Parse(tt.in, []byte(tt.in))

			if tt.err {
				r.Error(err)
				return
			}

			r.NoError(err)

			pt, ok := res.(*Path)
			r.True(ok)

			r.Equal(&tt.exp, pt)
		})
	}

}
