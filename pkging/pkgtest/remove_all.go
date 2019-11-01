package pkgtest

import (
	"testing"

	"github.com/markbates/pkger/pkging"
	"github.com/stretchr/testify/require"
)

func RemoveAllTest(t *testing.T, ref *Ref, pkg pkging.Pkger) {
	r := require.New(t)

	name := "/public/assets"
	r.NoError(pkg.MkdirAll(name, 0755))

	_, err := pkg.Stat(name)
	r.NoError(err)

	r.NoError(pkg.RemoveAll(name))

	_, err = pkg.Stat(name)
	r.Error(err)
}
