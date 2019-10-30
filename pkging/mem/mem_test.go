package mem_test

import (
	"testing"

	"github.com/markbates/pkger/pkging"
	"github.com/markbates/pkger/pkging/costello"
	"github.com/markbates/pkger/pkging/mem"
)

func Test_Pkger(t *testing.T) {
	ref, err := costello.NewRef()
	if err != nil {
		t.Fatal(err)
	}
	costello.All(t, ref, func(ref *costello.Ref) (pkging.Pkger, error) {
		return mem.New(ref.Info)
	})
}
