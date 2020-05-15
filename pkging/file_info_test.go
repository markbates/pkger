package pkging_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/markbates/pkger/here"
	"github.com/markbates/pkger/pkging"
	"github.com/stretchr/testify/require"
)

func Test_NewFileInfo(t *testing.T) {
	r := require.New(t)

	her, err := here.Current()
	r.NoError(err)

	exp, err := os.Stat(filepath.Join(her.Dir, "go.mod"))
	r.NoError(err)

	act := pkging.NewFileInfo(exp)

	r.Equal(exp.Name(), act.Name())
	r.Equal(exp.Size(), act.Size())
	r.Equal(exp.Mode(), act.Mode())
	r.Equal(exp.IsDir(), act.IsDir())
	r.Equal(exp.ModTime().Format(time.RFC3339), act.ModTime().Format(time.RFC3339))

}
