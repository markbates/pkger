package pkger

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gobuffalo/here"
	"github.com/markbates/pkger/paths"
	"github.com/markbates/pkger/pkgs"
)

type index struct {
	Pkg   string
	Files map[paths.Path]*File
}

func (i *index) Create(pt paths.Path) (*File, error) {
	her, err := pkgs.Pkg(pt.Pkg)
	if err != nil {
		return nil, err
	}
	f := &File{
		path:  pt,
		index: newIndex(),
		her:   her,
		info: &FileInfo{
			name:    strings.TrimPrefix(pt.Name, "/"),
			mode:    0666,
			modTime: time.Now(),
		},
	}

	i.Files[pt] = f
	return f, nil
}

func (i index) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{
		"pkg": i.Pkg,
	}

	fm := map[string]File{}

	for k, v := range i.Files {
		fm[k.String()] = *v
	}

	m["files"] = fm

	return json.Marshal(m)
}

func (i index) Walk(pt paths.Path, wf WalkFunc) error {
	if len(pt.Pkg) == 0 {
		pt.Pkg = i.Pkg
	}
	if len(i.Files) > 0 {
		for k, v := range i.Files {
			if k.Pkg != pt.Pkg {
				continue
			}
			if err := wf(k, v.info); err != nil {
				return err
			}
		}
	}

	var info here.Info
	var err error
	if pt.Pkg == "." {
		info, err = pkgs.Current()
		if err != nil {
			return err
		}
		pt.Pkg = info.ImportPath
	}

	if info.IsZero() {
		info, err = pkgs.Pkg(pt.Pkg)
		if err != nil {
			return fmt.Errorf("%s: %s", pt, err)
		}
	}
	fp := filepath.Join(info.Dir, pt.Name)
	err = filepath.Walk(fp, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		path = strings.TrimPrefix(path, info.Dir)
		pt, err := paths.Parse(fmt.Sprintf("%s:%s", pt.Pkg, path))
		if err != nil {
			return err
		}
		return wf(pt, NewFileInfo(fi))
	})

	return err
}

func (i index) Open(pt paths.Path) (*File, error) {
	if len(pt.Pkg) == 0 {
		pt.Pkg = i.Pkg
	}
	f, ok := i.Files[pt]
	if !ok {
		return i.openDisk(pt)
	}
	return &File{
		info:  f.info,
		path:  f.path,
		data:  f.data,
		her:   f.her,
		index: newIndex(),
	}, nil
}

func (i index) openDisk(pt paths.Path) (*File, error) {
	if len(pt.Pkg) == 0 {
		pt.Pkg = i.Pkg
	}
	info, err := pkgs.Pkg(pt.Pkg)
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
		index: &index{
			Files: map[paths.Path]*File{},
		},
	}
	return f, nil
}

func (i index) Parse(p string) (paths.Path, error) {
	pt, err := paths.Parse(p)
	if err != nil {
		return pt, err
	}

	if len(pt.Pkg) == 0 {
		pt.Pkg = i.Pkg
	}

	return pt, nil
}

func newIndex() *index {
	return &index{
		Files: map[paths.Path]*File{},
	}
}

var rootIndex = newIndex()
