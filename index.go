package pkger

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/markbates/pkger/here"
	"github.com/markbates/pkger/internal/debug"
)

var filesCache = &filesMap{}
var infosCache = &infosMap{}
var pathsCache = &pathsMap{}
var curOnce = &sync.Once{}
var currentInfo here.Info

func dubeg(key, format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	debug.Debug("[%s|%s] %s", key, s)
}

func Parse(p string) (Path, error) {
	dubeg("Parse", p)
	pt, ok := pathsCache.Load(p)
	if ok {
		return pt, nil
	}
	if len(p) == 0 {
		return build(p, "", "")
	}

	res := pathrx.FindAllStringSubmatch(p, -1)
	if len(res) == 0 {
		return pt, fmt.Errorf("could not parse %q", p)
	}

	matches := res[0]

	if len(matches) != 4 {
		return pt, fmt.Errorf("could not parse %q", p)
	}

	return build(p, matches[1], matches[3])
}

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

func Stat() (here.Info, error) {
	var err error
	curOnce.Do(func() {
		if currentInfo.IsZero() {
			currentInfo, err = here.Current()
		}
	})

	return currentInfo, err
}

func Create(name string) (*File, error) {
	pt, err := Parse(name)
	if err != nil {
		return nil, err
	}

	her, err := Info(pt.Pkg)
	if err != nil {
		return nil, err
	}
	f := &File{
		path: pt,
		her:  her,
		info: &FileInfo{
			name:    pt.Name,
			mode:    0666,
			modTime: time.Now(),
			virtual: true,
		},
	}

	filesCache.Store(pt, f)

	dir := filepath.Dir(pt.Name)
	if err := MkdirAll(dir, 0644); err != nil {
		return nil, err
	}
	return f, nil
}

func openDisk(pt Path) (*File, error) {
	dubeg("openDisk", pt.String())
	info, err := Info(pt.Pkg)
	if err != nil {
		return nil, err
	}
	fp := info.Dir
	if len(pt.Name) > 0 {
		fp = filepath.Join(fp, pt.Name)
	}

	fi, err := os.Stat(fp)
	if err != nil {
		return nil, err
	}
	f := &File{
		info: WithName(strings.TrimPrefix(pt.Name, "/"), NewFileInfo(fi)),
		her:  info,
		path: pt,
	}
	return f, nil
}
