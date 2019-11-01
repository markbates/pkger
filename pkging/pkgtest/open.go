package pkgtest

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/markbates/pkger/pkging"
	"github.com/stretchr/testify/require"
)

func OpenTest(t *testing.T, ref *Ref, pkg pkging.Pkger) {
	r := require.New(t)

	name := "/go.mod"

	osf, err := os.Open(filepath.Join(ref.Dir, name))
	r.NoError(err)

	osi, err := osf.Stat()
	r.NoError(err)

	_, err = LoadFile(name, ref, pkg)
	r.NoError(err)

	pf, err := pkg.Open(fmt.Sprintf("/%s", name))
	r.NoError(err)

	psi, err := pf.Stat()
	r.NoError(err)

	CmpFileInfo(t, osi, psi)

	osb, err := ioutil.ReadAll(osf)
	r.NoError(err)
	r.NoError(osf.Close())

	psb, err := ioutil.ReadAll(pf)
	r.NoError(err)
	r.NoError(pf.Close())

	r.Equal(osb, psb)
}
