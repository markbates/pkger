package here_test

import (
	"testing"

	"github.com/markbates/pkger/pkging/pkgtest"
	"github.com/stretchr/testify/require"
)

func Test_Dir(t *testing.T) {
	r := require.New(t)

	ref, err := pkgtest.NewRef()
	r.NoError(err)

	table := []struct {
		in  string
		err bool
	}{
		{in: ref.Dir, err: false},
	}
	for _, tt := range table {
		t.Run(tt.in, func(st *testing.T) {
			r := require.New(st)

			if tt.err {
				r.Error(err)
			}
			r.NoError(err)

		})
	}
}
