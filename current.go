// TODO: need to populate in memory cache when packed.
// you can't use go list, etc... in prod
package pkger

import (
	"github.com/gobuffalo/here"
)

func Pkg(p string) (here.Info, error) {
	return here.Cache(p, here.Package)
}

func Current() (here.Info, error) {
	return here.Cache("", func(string) (here.Info, error) {
		return here.Current()
	})
}
