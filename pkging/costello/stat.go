package costello

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/markbates/pkger/pkging"
	"github.com/stretchr/testify/require"
)

func StatTest(t *testing.T, pkg pkging.Pkger) {
	r := require.New(t)

	ref, err := NewRef()
	r.NoError(err)

	osi, err := os.Stat(filepath.Join(ref.Dir, "go.mod"))
	r.NoError(err)

	r.NoError(LoadRef(ref, pkg))
	psi, err := pkg.Stat("/go.mod")
	r.NoError(err)

	r.Equal(osi.Name(), psi.Name())
	r.Equal(osi.Mode(), psi.Mode())
	r.Equal(osi.Size(), psi.Size())
	r.Equal(osi.ModTime().Format(time.RFC3339), psi.ModTime().Format(time.RFC3339))
}
