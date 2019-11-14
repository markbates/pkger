package here

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_nonGoDirRx(t *testing.T) {
	r := require.New(t)
	r.False(nonGoDirRx.MatchString(""))
	r.False(nonGoDirRx.MatchString("hello"))

	table := []string{
		"go: cannot find main module; see 'go help modules'",
		"go help modules",
		"go: ",
		"build .:",
		"no Go files",
		"can't load package: package",
	}

	for _, tt := range table {
		t.Run(tt, func(st *testing.T) {
			r := require.New(st)

			b := nonGoDirRx.MatchString(tt)
			r.True(b)

		})
	}

}
