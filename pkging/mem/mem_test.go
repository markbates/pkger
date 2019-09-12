package mem

import (
	"testing"

	"github.com/markbates/pkger/here"
	"github.com/markbates/pkger/pkging"
	"github.com/markbates/pkger/pkging/pkgtest"
)

func Test_Pkger(t *testing.T) {
	suite, err := pkgtest.NewSuite("memos", func() (pkging.Pkger, error) {
		info, err := here.Current()
		if err != nil {
			return nil, err
		}

		wh, err := New(info)
		if err != nil {
			return nil, err
		}

		WithInfo(wh, info)
		return wh, nil
	})
	if err != nil {
		t.Fatal(err)
	}

	suite.Test(t)
}
