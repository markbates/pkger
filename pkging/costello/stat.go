package costello

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/markbates/pkger/pkging"
	"github.com/stretchr/testify/require"
)

func StatTest(t *testing.T, ref *Ref, pkg pkging.Pkger) {
	r := require.New(t)

	osi, err := os.Stat(filepath.Join(ref.Dir, "go.mod"))
	r.NoError(err)

	r.NoError(LoadRef(ref, pkg))
	psi, err := pkg.Stat("/go.mod")
	r.NoError(err)

	cmpFileInfo(t, osi, psi)
}
