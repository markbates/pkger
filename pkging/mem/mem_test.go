package mem_test

import (
	"testing"

	"github.com/markbates/pkger/pkging"
	"github.com/markbates/pkger/pkging/mem"
	"github.com/markbates/pkger/pkging/pkgtest"
)

func Test_Pkger(t *testing.T) {
	pkgtest.All(t, func(ref *pkgtest.Ref) (pkging.Pkger, error) {
		return mem.New(ref.Info)
	})
}
