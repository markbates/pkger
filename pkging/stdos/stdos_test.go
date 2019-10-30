package stdos_test

import (
	"io/ioutil"
	"testing"

	"github.com/markbates/pkger/pkging"
	"github.com/markbates/pkger/pkging/pkgtest"
	"github.com/markbates/pkger/pkging/stdos"
)

func Test_Pkger(t *testing.T) {
	t.SkipNow()
	suite, err := pkgtest.NewSuite("stdos", func() (pkging.Pkger, error) {
		app, err := pkgtest.App()
		if err != nil {
			return nil, err
		}

		dir, err := ioutil.TempDir("", "stdos")
		if err != nil {
			return nil, err
		}

		app.Dir = dir

		mypkging, err := stdos.New(app.Info)
		if err != nil {
			return nil, err
		}

		return mypkging, nil
	})
	if err != nil {
		t.Fatal(err)
	}

	suite.Test(t)
}
