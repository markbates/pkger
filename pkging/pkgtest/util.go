package pkgtest

import (
	"io"
	"path/filepath"
	"strings"
	"testing"

	"github.com/markbates/pkger/pkging/pkgutil"
	"github.com/stretchr/testify/require"
)

func (s Suite) Test_Util_ReadFile(t *testing.T) {
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
		{in: hart},
	}

	for _, tt := range table {
		t.Run(tt.in, func(st *testing.T) {

			r := require.New(st)

			pt, err := s.Parse(tt.in)
			r.NoError(err)

			r.NoError(s.RemoveAll(pt.String()))
			r.NoError(s.MkdirAll(filepath.Dir(pt.Name), 0755))

			f, err := s.Create(tt.in)
			r.NoError(err)

			body := "!" + pt.String()
			_, err = io.Copy(f, strings.NewReader(body))
			r.NoError(err)
			r.NoError(f.Close())

			b, err := pkgutil.ReadFile(s, tt.in)
			r.NoError(err)
			r.Equal(body, string(b))
		})
	}
}
