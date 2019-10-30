package costello

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
	openTest("go.mod", t, ref, pkg)
}

func openTest(name string, t *testing.T, ref *Ref, pkg pkging.Pkger) {
	r := require.New(t)

	osf, err := os.Open(filepath.Join(ref.Dir, name))
	r.NoError(err)

	osi, err := osf.Stat()
	r.NoError(err)

	r.NoError(LoadRef(ref, pkg))

	pf, err := pkg.Open(fmt.Sprintf("/%s", name))
	r.NoError(err)

	psi, err := pf.Stat()
	r.NoError(err)

	cmpFileInfo(t, osi, psi)

	if osi.IsDir() {
		return
	}

	osb, err := ioutil.ReadAll(osf)
	r.NoError(err)
	r.NoError(osf.Close())

	psb, err := ioutil.ReadAll(pf)
	r.NoError(err)
	r.NoError(pf.Close())

	r.Equal(osb, psb)
}
