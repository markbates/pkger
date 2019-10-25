package pkging_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/markbates/pkger/pkging"
	"github.com/markbates/pkger/pkging/pkgtest"
	"github.com/stretchr/testify/require"
)

func Test_NewFileInfo(t *testing.T) {
	r := require.New(t)

	app, err := pkgtest.App()
	r.NoError(err)

	exp, err := os.Stat(filepath.Join(app.Info.Dir, "go.mod"))
	r.NoError(err)

	act := pkging.NewFileInfo(exp)

	r.Equal(exp.Name(), act.Name())
	r.Equal(exp.Size(), act.Size())
	r.Equal(exp.Mode(), act.Mode())
	r.Equal(exp.IsDir(), act.IsDir())
	r.Equal(exp.ModTime().Format(time.RFC3339), act.ModTime().Format(time.RFC3339))

}
