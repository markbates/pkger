package pkgtest

import (
	"path/filepath"
	"sort"
	"testing"

	"github.com/markbates/pkger/pkging/pkgutil"
	"github.com/stretchr/testify/require"
)

func (s Suite) Test_File_Info(t *testing.T) {
	r := require.New(t)

	pkg, err := s.Make()
	r.NoError(err)

	cur, err := pkg.Current()
	r.NoError(err)

	ip := cur.ImportPath
	table := []struct {
		in string
	}{
		{in: mould},
		{in: ":" + mould},
		{in: ip + ":" + mould},
	}

	for _, tt := range table {
		s.Run(t, tt.in, func(st *testing.T) {
			r := require.New(st)

			r.NoError(pkg.RemoveAll(mould))
			r.NoError(pkg.MkdirAll(filepath.Dir(tt.in), 0755))
			err := pkgutil.WriteFile(pkg, tt.in, []byte(mould), 0644)
			r.NoError(err)

			f, err := pkg.Open(tt.in)
			r.NoError(err)
			r.Equal(mould, f.Name())
			r.Equal(cur.ImportPath, f.Info().ImportPath)
			r.NoError(f.Close())
		})
	}
}

// func (s Suite) Test_File_Read(t *testing.T) {
// 	panic("not implemented")
// }
//
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
		})
	}
}

//
// func (s Suite) Test_File_Seek(t *testing.T) {
// 	panic("not implemented")
// }
//
// func (s Suite) Test_File_Stat(t *testing.T) {
// 	panic("not implemented")
// }
//
// func (s Suite) Test_File_Write(t *testing.T) {
// 	panic("not implemented")
// }
