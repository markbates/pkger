package pkger

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gobuffalo/here"
	"github.com/markbates/hepa"
	"github.com/markbates/hepa/filters"
)

type index struct {
	Files   *filesMap
	Infos   *infosMap
	current here.Info
	once    sync.Once
}

func (i *index) Info(p string) (here.Info, error) {
	info, ok := i.Infos.Load(p)
	if ok {
		return info, nil
	}

	info, err := here.Cache(p, here.Package)
	if err != nil {
		return info, err
	}
	i.Infos.Store(p, info)
	return info, nil
}

func (i *index) Current() (here.Info, error) {
	var err error
	i.once.Do(func() {
		i.current, err = here.Cache("", func(string) (here.Info, error) {
			return here.Current()
		})
	})

	return i.current, err
}

func (i *index) Create(pt Path) (*File, error) {
	her, err := Info(pt.Pkg)
	if err != nil {
		return nil, err
	}
	f := &File{
		path: pt,
		her:  her,
		info: &FileInfo{
			name:    strings.TrimPrefix(pt.Name, "/"),
			mode:    0666,
			modTime: time.Now(),
		},
	}

	i.Files.Store(pt, f)
	return f, nil
}

func (i *index) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{}

	m["files"] = i.Files
	m["infos"] = i.Infos
	m["current"] = i.current

	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}

	hep := hepa.New()
	hep = hepa.With(hep, filters.Golang())
	hep = hepa.With(hep, filters.Secrets())
	return hep.Filter(b)
}

func (i *index) UnmarshalJSON(b []byte) error {
	m := map[string]json.RawMessage{}

	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}

	infos, ok := m["infos"]
	if !ok {
		return fmt.Errorf("missing infos")
	}
	i.Infos = &infosMap{}
	if err := json.Unmarshal(infos, i.Infos); err != nil {
		return err
	}

	files, ok := m["files"]
	if !ok {
		return fmt.Errorf("missing files")
	}

	i.Files = &filesMap{}
	if err := json.Unmarshal(files, i.Files); err != nil {
		return err
	}

	current, ok := m["current"]
	if !ok {
		return fmt.Errorf("missing current")
	}
	if err := json.Unmarshal(current, &i.current); err != nil {
		return err
	}
	return nil
}

func (i index) Walk(pt Path, wf WalkFunc) error {
	var err error
	i.Files.Range(func(k Path, v *File) bool {
		if k.Pkg != pt.Pkg {
			return true
		}
		if err = wf(k, v.info); err != nil {
			return false
		}
		return true
	})

	if err != nil {
		return err
	}

	var info here.Info
	if pt.Pkg == "." {
		info, err = Current()
		if err != nil {
			return err
		}
		pt.Pkg = info.ImportPath
	}

	if info.IsZero() {
		info, err = Info(pt.Pkg)
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
		pt, err := Parse(fmt.Sprintf("%s:%s", pt.Pkg, path))
		if err != nil {
			return err
		}
		return wf(pt, NewFileInfo(fi))
	})

	return err
}

func (i *index) Open(pt Path) (*File, error) {
	f, ok := i.Files.Load(pt)
	if !ok {
		return i.openDisk(pt)
	}
	nf := &File{
		info: f.info,
		path: f.path,
		data: f.data,
		her:  f.her,
	}

	return nf, nil
}

func (i index) openDisk(pt Path) (*File, error) {
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

func newIndex() *index {
	return &index{
		Files: &filesMap{},
		Infos: &infosMap{},
	}
}

var rootIndex = func() *index {
	i := newIndex()
	return i
}()
