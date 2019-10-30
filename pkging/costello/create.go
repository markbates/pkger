package costello

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/markbates/pkger/pkging"
	"github.com/stretchr/testify/require"
)

func CreateTest(t *testing.T, pkg pkging.Pkger) {
	r := require.New(t)

	ref, err := NewRef()
	r.NoError(err)

	const name = "create.test"

	fp := filepath.Join(ref.Dir, name)
	os.RemoveAll(fp)
	defer os.RemoveAll(fp)

	_, err = os.Stat(fp)
	r.Error(err)

	_, err = pkg.Stat(name)
	r.Error(err)

	data := []byte(strings.ToUpper(name))

	osf, err := os.Create(fp)
	r.NoError(err)

	_, err = osf.Write(data)
	r.NoError(err)
	r.NoError(osf.Close())

	psf, err := pkg.Create(fmt.Sprintf("/%s", name))
	r.NoError(err)

	_, err = psf.Write(data)
	r.NoError(err)
	r.NoError(psf.Close())
	openTest(name, t, pkg)
}
