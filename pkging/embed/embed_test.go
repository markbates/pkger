package embed

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Encoding(t *testing.T) {
	r := require.New(t)

	in := []byte("hi\n")

	enc, err := Encode(in)
	r.NoError(err)

	r.NotEqual(in, enc)

	dec, err := Decode(enc)
	r.NoError(err)
	r.Equal(in, dec)
}
