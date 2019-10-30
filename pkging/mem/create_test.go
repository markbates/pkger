package mem

import (
	"testing"

	"github.com/markbates/pkger/pkging/costello"
	"github.com/stretchr/testify/require"
)

func Test_Pkger_Open(t *testing.T) {
	r := require.New(t)

	app, err := costello.NewRef()
	r.NoError(err)

	pkg, err := New(app.Info)
	r.NoError(err)

	err = costello.LoadRef(app, pkg)
	r.NoError(err)

	costello.OpenTest(t, pkg)
}
