package stdos

import (
	"testing"

	"github.com/markbates/pkger/pkging"
	"github.com/markbates/pkger/pkging/costello"
)

func Test_Pkger(t *testing.T) {
	costello.All(t, func(ref *costello.Ref) (pkging.Pkger, error) {
		return New(ref.Info)
	})
}
