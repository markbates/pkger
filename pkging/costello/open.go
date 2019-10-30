package costello

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

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

	osb, err := ioutil.ReadAll(osf)
	r.NoError(err)
	r.NoError(osf.Close())

	r.NoError(LoadRef(ref, pkg))

	pf, err := pkg.Open(fmt.Sprintf("/%s", name))
	r.NoError(err)

	psi, err := pf.Stat()
	r.NoError(err)

	psb, err := ioutil.ReadAll(pf)
	r.NoError(err)
	r.NoError(pf.Close())

	r.Equal(osi.Name(), psi.Name())
	r.Equal(osi.Mode(), psi.Mode())
	r.Equal(osi.Size(), psi.Size())
	r.Equal(osi.ModTime().Format(time.RFC3339), psi.ModTime().Format(time.RFC3339))
	r.Equal(osb, psb)
}