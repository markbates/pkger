package costello

import (
	"testing"

	"github.com/markbates/pkger/pkging"
	"github.com/stretchr/testify/require"
)

func RemoveAllTest(t *testing.T, ref *Ref, pkg pkging.Pkger) {
	r := require.New(t)

	r.NoError(LoadRef(ref, pkg))

	name := "/public/assets"

	_, err := pkg.Stat(name)
	r.NoError(err)

	r.NoError(pkg.RemoveAll(name))

	_, err = pkg.Stat(name)
	r.Error(err)
}
