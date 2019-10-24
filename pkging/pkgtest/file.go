package pkgtest

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

func (s Suite) Test_File_Readdir(t *testing.T) {
	r := require.New(t)

	pkg, err := s.Make()
	r.NoError(err)

	cur, err := pkg.Current()
	r.NoError(err)

	ip := cur.ImportPath
	table := []struct {
		in string
	}{
		{in: ":/public"},
		{in: ip + ":/public"},
	}

	r.NoError(s.LoadFolder(pkg))

	for _, tt := range table {
		s.Run(t, tt.in, func(st *testing.T) {
			r := require.New(st)

			dir, err := pkg.Open(tt.in)
			r.NoError(err)
			defer dir.Close()

			infos, err := dir.Readdir(-1)
			r.NoError(err)
			r.Len(infos, 2)

			sort.Slice(infos, func(i, j int) bool {
				return infos[i].Name() < infos[j].Name()
			})

			r.Equal("images", infos[0].Name())
			r.Equal("index.html", infos[1].Name())

			dir, err = pkg.Open(tt.in + "/images")
			r.NoError(err)

			infos, err = dir.Readdir(-1)
			r.NoError(err)
			r.Len(infos, 2)

			sort.Slice(infos, func(i, j int) bool {
				return infos[i].Name() < infos[j].Name()
			})

			r.Equal("img1.png", infos[0].Name())
			r.Equal("img2.png", infos[1].Name())

		})
	}
}
