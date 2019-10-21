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

	app, err := App()
	r.NoError(err)

	ip := app.Info.ImportPath
	mould := "/public/index.html"

	table := []struct {
		in string
	}{
		{in: mould},
		{in: ip + ":" + mould},
	}

	for _, tt := range table {
		s.Run(t, tt.in, func(st *testing.T) {
			r := require.New(st)

			pkg, err := s.Make()
			r.NoError(err)

			pt, err := pkg.Parse(tt.in)
			r.NoError(err)

			r.NoError(pkg.RemoveAll(pt.String()))
			r.NoError(pkg.MkdirAll(filepath.Dir(pt.Name), 0755))

			f, err := pkg.Create(tt.in)
			r.NoError(err)

			body := "!" + pt.String()
			_, err = io.Copy(f, strings.NewReader(body))
			r.NoError(err)
			r.NoError(f.Close())

			b, err := pkgutil.ReadFile(pkg, tt.in)
			r.NoError(err)
			r.Equal(body, string(b))
		})
	}
}
