package pkgtest

import (
	"testing"

	"github.com/markbates/pkger/pkging"
	"github.com/stretchr/testify/require"
)

func RemoveTest(t *testing.T, ref *Ref, pkg pkging.Pkger) {
	r := require.New(t)

	name := "/go.mod"
	_, err := LoadFile(name, ref, pkg)
	r.NoError(err)

	_, err = pkg.Stat(name)
	r.NoError(err)

	r.NoError(pkg.Remove(name))

	_, err = pkg.Stat(name)
	r.Error(err)
}
