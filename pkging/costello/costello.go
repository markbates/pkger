package costello

import (
	"testing"

	"github.com/markbates/pkger/pkging"
	"github.com/stretchr/testify/require"
)

type AllFn func(ref *Ref) (pkging.Pkger, error)

func All(t *testing.T, ref *Ref, fn AllFn) {
	r := require.New(t)

	type tf func(*testing.T, *Ref, pkging.Pkger)

	tests := map[string]tf{
		"OpenTest":   OpenTest,
		"StatTest":   StatTest,
		"CreateTest": CreateTest,
	}

	for n, tt := range tests {
		t.Run(n, func(st *testing.T) {
			pkg, err := fn(ref)
			r.NoError(err)

			tt(st, ref, pkg)
		})
	}

}
