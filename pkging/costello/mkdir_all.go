package costello

import (
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/markbates/pkger/pkging"
	"github.com/stretchr/testify/require"
)

func MkdirAllTest(t *testing.T, ref *Ref, pkg pkging.Pkger) {
	r := require.New(t)

	parts := []string{"all", "this", "useless", "beauty"}

	fp := ref.Dir
	for _, part := range parts {
		fp = filepath.Join(fp, part)
	}

	os.RemoveAll(fp)
	defer os.RemoveAll(fp)

	_, err := os.Stat(fp)
	r.Error(err)

	name := path.Join(parts...)
	if !strings.HasPrefix(name, "/") {
		name = "/" + name
	}

	_, err = pkg.Stat(name)
	r.Error(err)

	r.NoError(os.MkdirAll(fp, 0755))
	r.NoError(pkg.MkdirAll(name, 0755))

	openTest(name, t, ref, pkg)
}
