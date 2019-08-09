package pkger

import "github.com/markbates/pkger/here"

func Stat() (here.Info, error) {
	var err error
	curOnce.Do(func() {
		if currentInfo.IsZero() {
			currentInfo, err = here.Current()
		}
	})

	return currentInfo, err
}
