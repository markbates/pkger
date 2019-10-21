package mem_test

import (
	"testing"

	"github.com/markbates/pkger/pkging"
	"github.com/markbates/pkger/pkging/mem"
	"github.com/markbates/pkger/pkging/pkgtest"
)

func Test_Pkger(t *testing.T) {
	suite, err := pkgtest.NewSuite("memos", func() (pkging.Pkger, error) {
		app, err := pkgtest.App()
		if err != nil {
			return nil, err
		}

		pkg, err := mem.New(app.Info)
		if err != nil {
			return nil, err
		}

		return pkg, nil
	})
	if err != nil {
		t.Fatal(err)
	}

	suite.Test(t)
}
