package mem

import (
	"os"
	"testing"

	"github.com/markbates/pkger/pkging/costello"
	"github.com/stretchr/testify/require"
)

func Test_Pkger_Info(t *testing.T) {
	r := require.New(t)

	ref, err := costello.NewRef()
	r.NoError(err)
	defer os.RemoveAll(ref.Dir)

	pkg, err := New(ref.Info)
	r.NoError(err)

	costello.InfoTest(t, ref, pkg)
}
