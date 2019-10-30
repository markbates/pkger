package costello

import (
	"fmt"
	"testing"

	"github.com/markbates/pkger/pkging"
	"github.com/stretchr/testify/require"
)

type AllFn func(ref *Ref) (pkging.Pkger, error)

func All(t *testing.T, ref *Ref, fn AllFn) {
	r := require.New(t)

	type tf func(*testing.T, *Ref, pkging.Pkger)

	tests := map[string]tf{
		"OpenTest":    OpenTest,
		"StatTest":    StatTest,
		"CreateTest":  CreateTest,
		"CurrentTest": CurrentTest,
		"InfoTest":    InfoTest,
		"MkdirAll":    MkdirAllTest,
	}

	pkg, err := fn(ref)
	r.NoError(err)

	for n, tt := range tests {
		t.Run(fmt.Sprintf("%T/%s", pkg, n), func(st *testing.T) {
			st.Parallel()
			pkg, err := fn(ref)
			r.NoError(err)

			tt(st, ref, pkg)
		})
	}

}
