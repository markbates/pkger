package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Parser(t *testing.T) {
	r := require.New(t)

	res, err := Parse("github.com/markbates/pkger:/internal/examples")
	r.NoError(err)

	// for _, pt := range res.Paths {
	// 	fmt.Println(pt)
	// }
	r.True(len(res.Paths) > 3)
}
