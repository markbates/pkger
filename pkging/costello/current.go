package costello

import (
	"testing"

	"github.com/markbates/pkger/pkging"
)

func CurrentTest(t *testing.T, ref *Ref, pkg pkging.Pkger) {
	// panic("ref")
	cur, err := pkg.Current()
	if err != nil {
		t.Fatal(err)
	}
	cmpHereInfo(t, ref.Info, cur)
}
