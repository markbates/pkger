package costello

import (
	"testing"

	"github.com/markbates/pkger/pkging"
	"github.com/stretchr/testify/require"
)

func RemoveTest(t *testing.T, ref *Ref, pkg pkging.Pkger) {
	r := require.New(t)

	r.NoError(LoadRef(ref, pkg))

	name := "/go.mod"

	_, err := pkg.Stat(name)
	r.NoError(err)

	r.NoError(pkg.Remove(name))

	_, err = pkg.Stat(name)
	r.Error(err)
}
