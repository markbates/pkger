package pkgtest

import (
	"path/filepath"
	"testing"

	"github.com/markbates/pkger/pkging/pkgutil"
	"github.com/stretchr/testify/require"
)

func (s Suite) Test_File_Info(t *testing.T) {
	r := require.New(t)
	cur, err := s.Current()
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
		t.Run(tt.in, func(st *testing.T) {
			r := require.New(st)

			r.NoError(s.RemoveAll(mould))
			r.NoError(s.MkdirAll(filepath.Dir(tt.in), 0755))
			err := pkgutil.WriteFile(s, tt.in, []byte(mould), 0644)
			r.NoError(err)

			f, err := s.Open(tt.in)
			r.NoError(err)
			r.Equal(cur.ImportPath, f.Info().ImportPath)
		})
	}
}

func (s Suite) Test_File_Name(t *testing.T) {
	panic("not implemented")
}

func (s Suite) Test_File_Open(t *testing.T) {
	panic("not implemented")
}

func (s Suite) Test_File_Path(t *testing.T) {
	panic("not implemented")
}

func (s Suite) Test_File_Read(t *testing.T) {
	panic("not implemented")
}

func (s Suite) Test_File_Readdir(t *testing.T) {
	panic("not implemented")
}

func (s Suite) Test_File_Seek(t *testing.T) {
	panic("not implemented")
}

func (s Suite) Test_File_Stat(t *testing.T) {
	panic("not implemented")
}

func (s Suite) Test_File_Write(t *testing.T) {
	panic("not implemented")
}
