package mem

import (
	"os"
	"testing"

	"github.com/markbates/pkger/here"
	"github.com/markbates/pkger/pkging"
	"github.com/markbates/pkger/pkging/pkgtest"
	"github.com/markbates/pkger/pkging/stdos"
)

func Test_Pkger(t *testing.T) {
	suite, err := pkgtest.NewSuite("memos", func() (pkging.Pkger, error) {
		info, err := here.Current()
		if err != nil {
			return nil, err
		}

		pkg, err := New(info)
		if err != nil {
			return nil, err
		}

		disk, err := stdos.New(info)
		if err != nil {
			return nil, err
		}

		err = disk.Walk("/examples/app/public", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			pt, err := disk.Parse(path)
			if err != nil {
				return err
			}

			f, err := disk.Open(pt.String())
			if err != nil {
				return err
			}
			defer f.Close()
			if err := pkg.Add(f); err != nil {
				return err
			}

			return nil
		})

		return pkg, nil
	})
	if err != nil {
		t.Fatal(err)
	}

	suite.Test(t)
}
