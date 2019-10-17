package stdos

import (
	"testing"

	"github.com/markbates/pkger/here"
	"github.com/stretchr/testify/require"
)

func Test_File_Stat_No_Info(t *testing.T) {
	r := require.New(t)

	her, err := here.Current()
	r.NoError(err)
	pkg, err := New(her)
	r.NoError(err)

	f, err := pkg.Open(":/pkging/stdos/file_test.go")
	r.NoError(err)
	defer f.Close()

	sf, ok := f.(*File)
	r.True(ok)

	oi := sf.info
	sf.info = nil

	info, err := sf.Stat()
	r.NoError(err)
	r.Equal(oi.Name(), info.Name())
	// r.Equal("", f.Name())
}
