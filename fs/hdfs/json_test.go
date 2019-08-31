package hdfs

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_File_JSON(t *testing.T) {
	r := require.New(t)

	fs := NewFS()

	f, err := fs.Create("/radio.radio")
	r.NoError(err)
	_, err = io.Copy(f, strings.NewReader(radio))
	r.NoError(err)
	r.NoError(f.Close())

	f, err = fs.Open("/radio.radio")
	r.NoError(err)
	bi, err := f.Stat()
	r.NoError(err)

	mj, err := json.Marshal(f)
	r.NoError(err)

	f2 := &File{}

	r.NoError(json.Unmarshal(mj, f2))

	ai, err := f2.Stat()
	r.NoError(err)

	r.Equal(bi.Size(), ai.Size())

	fd, err := ioutil.ReadAll(f2)
	r.NoError(err)
	r.Equal(radio, string(fd))
}
