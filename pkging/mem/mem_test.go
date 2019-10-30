package mem_test

import (
	"testing"

	"github.com/markbates/pkger/pkging"
	"github.com/markbates/pkger/pkging/costello"
	"github.com/markbates/pkger/pkging/mem"
)

func Test_Pkger(t *testing.T) {
	costello.All(t, func(ref *costello.Ref) (pkging.Pkger, error) {
		return mem.New(ref.Info)
	})
}
