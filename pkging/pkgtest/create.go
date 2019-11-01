package pkgtest

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/markbates/pkger/pkging"
	"github.com/stretchr/testify/require"
)

func CreateTest(t *testing.T, ref *Ref, pkg pkging.Pkger) {
	r := require.New(t)

	const name = "/create.test"

	_, err := pkg.Stat(name)
	r.Error(err)

	data := []byte(strings.ToUpper(name))

	f, err := pkg.Create(name)
	r.NoError(err)

	_, err = f.Write(data)
	r.NoError(err)
	r.NoError(f.Close())

	f, err = pkg.Open(name)
	r.NoError(err)

	info, err := f.Stat()
	r.NoError(err)

	b, err := ioutil.ReadAll(f)
	r.NoError(err)
	r.NoError(f.Close())

	r.Equal(data, b)
	r.Equal("create.test", info.Name())

}
