package pkgtest

import (
	"testing"

	"github.com/markbates/pkger/pkging"
	"github.com/stretchr/testify/require"
)

func MkdirAllTest(t *testing.T, ref *Ref, pkg pkging.Pkger) {
	r := require.New(t)

	name := "/all/this/useless/beauty"

	_, err := pkg.Stat(name)
	r.Error(err)

	r.NoError(pkg.MkdirAll(name, 0755))

	f, err := pkg.Open(name)
	r.NoError(err)

	info, err := f.Stat()
	r.NoError(err)

	r.Equal("app:"+name, f.Name())
	r.Equal("beauty", info.Name())
}
