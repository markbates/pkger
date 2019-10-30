package costello

import (
	"testing"

	"github.com/markbates/pkger/pkging"
	"github.com/stretchr/testify/require"
)

type AllFn func(ref *Ref) (pkging.Pkger, error)

func All(t *testing.T, fn AllFn) {
	r := require.New(t)
	type tf func(*testing.T, pkging.Pkger)

	tests := map[string]tf{
		"OpenTest": OpenTest,
		"StatTest": StatTest,
	}

	ref, err := NewRef()
	r.NoError(err)

	for n, tt := range tests {
		t.Run(n, func(st *testing.T) {
			pkg, err := fn(ref)
			r.NoError(err)

			tt(st, pkg)
		})
	}

}
