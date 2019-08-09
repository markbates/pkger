package pkger

import "github.com/markbates/pkger/here"

func Info(p string) (here.Info, error) {
	info, ok := infosCache.Load(p)
	if ok {
		return info, nil
	}

	info, err := here.Package(p)
	if err != nil {
		return info, err
	}
	infosCache.Store(p, info)
	return info, nil
}
