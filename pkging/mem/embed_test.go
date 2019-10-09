package mem

import (
	"fmt"
	"os"
	"testing"

	"github.com/markbates/pkger/here"
	"github.com/stretchr/testify/require"
)

func Test_Pkger_MarshalEmbed(t *testing.T) {
	r := require.New(t)

	info, err := here.Current()
	r.NoError(err)

	pkg, err := New(info)
	r.NoError(err)

	const N = 10
	for i := 0; i < N; i++ {
		name := fmt.Sprintf("/%d.txt", i)
		f, err := pkg.Create(name)
		r.NoError(err)
		f.Write([]byte(name))
		r.NoError(f.Close())
	}

	em, err := pkg.MarshalEmbed()
	r.NoError(err)

	p2, err := UnmarshalEmbed(em)
	r.NoError(err)

	pinfo, err := p2.Current()
	r.NoError(err)
	r.Equal(info.ImportPath, pinfo.ImportPath)
	r.Equal(info.Name, pinfo.Name)
	r.Equal(info.Imports, pinfo.Imports)

	m1 := info.Module
	m2 := pinfo.Module
	r.Equal(m1.Main, m2.Main)
	r.Equal(m1.GoVersion, m2.GoVersion)

	var act []os.FileInfo
	err = p2.Walk("/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		act = append(act, info)

		return nil
	})
	r.NoError(err)
	r.Len(act, N+1) // +1 for the / dir
}
