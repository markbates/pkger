package costello

import (
	"testing"

	"github.com/markbates/pkger/pkging"
	"github.com/stretchr/testify/require"
)

func CurrentTest(t *testing.T, ref *Ref, pkg pkging.Pkger) {
	r := require.New(t)

	cur, err := pkg.Current()
	r.NoError(err)
	r.Equal(ref.Info, cur)
}
