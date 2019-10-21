package mem

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/markbates/pkger/here"
	"github.com/markbates/pkger/internal/maps"
	"github.com/markbates/pkger/pkging"
	"github.com/markbates/pkger/pkging/embed"
)

var _ pkging.Pkger = &Pkger{}

// New returns *Pkger for the provided here.Info
func New(info here.Info) (*Pkger, error) {
	f := &Pkger{
		infos: &maps.Infos{},
		files: &maps.Files{},
		Here:  info,
	}
	f.infos.Store(info.ImportPath, info)
	return f, nil
}

type Pkger struct {
	Here  here.Info
	infos *maps.Infos
	files *maps.Files
}

// MarshalJSON creates a fully re-hydratable JSON representation of *Pkger
func (p *Pkger) MarshalJSON() ([]byte, error) {
	files := map[string]embed.File{}

	p.files.Range(func(key here.Path, file pkging.File) bool {
		f, ok := file.(*File)
		if !ok {
			return true
		}
		ef := embed.File{
			Info:   f.info,
			Here:   f.Here,
			Path:   f.path,
			Parent: f.parent,
			Data:   f.data,
		}
		files[key.String()] = ef
		return true
	})

	infos := map[string]here.Info{}
	p.infos.Range(func(key string, info here.Info) bool {
		infos[key] = info
		return true
	})

	return json.Marshal(embed.Data{
		Infos: infos,
		Files: files,
		Here:  p.Here,
	})
}

// UnmarshalJSON re-hydrates the *Pkger
func (p *Pkger) UnmarshalJSON(b []byte) error {
	y := &embed.Data{
		Infos: map[string]here.Info{},
		Files: map[string]embed.File{},
	}

	if err := json.Unmarshal(b, &y); err != nil {
		return err
	}

	p.Here = y.Here
	p.infos = &maps.Infos{}
	for k, v := range y.Infos {
		p.infos.Store(k, v)
	}

	p.files = &maps.Files{}
	for k, v := range y.Files {
		pt, err := p.Parse(k)
		if err != nil {
			return err
		}

		f := &File{
			Here:   v.Here,
			info:   v.Info,
			path:   v.Path,
			data:   v.Data,
			parent: v.Parent,
		}
		p.files.Store(pt, f)
	}
	return nil
}

// Abs returns an absolute representation of path. If the path is not absolute it will be joined with the current working directory to turn it into an absolute path. The absolute path name for a given file is not guaranteed to be unique. Abs calls Clean on the result.
func (f *Pkger) Abs(p string) (string, error) {
	pt, err := f.Parse(p)
	if err != nil {
		return "", err
	}
	return f.AbsPath(pt)
}

// AbsPath returns an absolute representation of here.Path. If the path is not absolute it will be joined with the current working directory to turn it into an absolute path. The absolute path name for a given file is not guaranteed to be unique. AbsPath calls Clean on the result.
func (f *Pkger) AbsPath(pt here.Path) (string, error) {
	return pt.String(), nil
}

// Current returns the here.Info representing the current Pkger implementation.
func (f *Pkger) Current() (here.Info, error) {
	return f.Here, nil
}

// Info returns the here.Info of the here.Path
func (f *Pkger) Info(p string) (here.Info, error) {
	info, ok := f.infos.Load(p)
	if !ok {
		return info, fmt.Errorf("no such package %q", p)
	}

	return info, nil
}

// Parse the string in here.Path format.
func (f *Pkger) Parse(p string) (here.Path, error) {
	return f.Here.Parse(p)
}

// Remove removes the named file or (empty) directory.
func (fx *Pkger) Remove(name string) error {
	pt, err := fx.Parse(name)
	if err != nil {
		return err
	}

	if _, ok := fx.files.Load(pt); !ok {
		return &os.PathError{"remove", pt.String(), fmt.Errorf("no such file or directory")}
	}

	fx.files.Delete(pt)
	return nil
}

// RemoveAll removes path and any children it contains. It removes everything it can but returns the first error it encounters. If the path does not exist, RemoveAll returns nil (no error).
func (fx *Pkger) RemoveAll(name string) error {
	pt, err := fx.Parse(name)
	if err != nil {
		return err
	}

	fx.files.Range(func(key here.Path, file pkging.File) bool {
		if strings.HasPrefix(key.Name, pt.Name) {
			fx.files.Delete(key)
		}
		return true
	})

	return nil
}

// Create creates the named file with mode 0666 (before umask) - It's actually 0644, truncating it if it already exists. If successful, methods on the returned File can be used for I/O; the associated file descriptor has mode O_RDWR.
func (fx *Pkger) Create(name string) (pkging.File, error) {
	fx.MkdirAll("/", 0755)
	pt, err := fx.Parse(name)
	if err != nil {
		return nil, err
	}

	her, err := fx.Info(pt.Pkg)
	if err != nil {
		return nil, err
	}

	dir := filepath.Dir(pt.Name)
	if dir != "/" {
		if _, err := fx.Stat(dir); err != nil {
			return nil, err
		}
	}

	f := &File{
		Here: her,
		path: pt,
		info: &pkging.FileInfo{
			Details: pkging.Details{
				Name:    pt.Name,
				Mode:    0644,
				ModTime: pkging.ModTime(time.Now()),
			},
		},
		pkging: fx,
	}

	fx.files.Store(pt, f)

	return f, nil
}

// MkdirAll creates a directory named path, along with any necessary parents, and returns nil, or else returns an error. The permission bits perm (before umask) are used for all directories that MkdirAll creates. If path is already a directory, MkdirAll does nothing and returns nil.
func (fx *Pkger) MkdirAll(p string, perm os.FileMode) error {
	path, err := fx.Parse(p)
	if err != nil {
		return err
	}
	root := path.Name

	cur, err := fx.Current()
	if err != nil {
		return err
	}
	for root != "" {
		pt := here.Path{
			Pkg:  path.Pkg,
			Name: root,
		}
		if _, ok := fx.files.Load(pt); ok {
			root = filepath.Dir(root)
			if root == "/" || root == "\\" {
				break
			}
			continue
		}
		f := &File{
			Here:   cur,
			pkging: fx,
			path:   pt,
			info: &pkging.FileInfo{
				Details: pkging.Details{
					Name:    pt.Name,
					Mode:    perm,
					ModTime: pkging.ModTime(time.Now()),
				},
			},
		}

		if err != nil {
			return err
		}
		f.info.Details.IsDir = true
		f.info.Details.Mode = perm
		if err := f.Close(); err != nil {
			return err
		}
		fx.files.Store(pt, f)
		root = filepath.Dir(root)
	}

	return nil

}

// Open opens the named file for reading. If successful, methods on the returned file can be used for reading; the associated file descriptor has mode O_RDONLY.
func (fx *Pkger) Open(name string) (pkging.File, error) {
	pt, err := fx.Parse(name)
	if err != nil {
		return nil, &os.PathError{
			Op:   "open",
			Path: name,
			Err:  err,
		}
	}

	fl, ok := fx.files.Load(pt)
	if !ok {
		return nil, os.ErrNotExist
	}
	f, ok := fl.(*File)
	if !ok {
		return nil, os.ErrNotExist
	}
	nf := &File{
		pkging: fx,
		info:   pkging.WithName(f.info.Name(), f.info),
		path:   f.path,
		data:   f.data,
		Here:   f.Here,
	}

	return nf, nil
}

// Stat returns a FileInfo describing the named file.
func (fx *Pkger) Stat(name string) (os.FileInfo, error) {
	pt, err := fx.Parse(name)
	if err != nil {
		return nil, err
	}
	f, ok := fx.files.Load(pt)
	if ok {
		return f.Stat()
	}
	return nil, fmt.Errorf("could not stat %s", pt)
}

// Walk walks the file tree rooted at root, calling walkFn for each file or directory in the tree, including root. All errors that arise visiting files and directories are filtered by walkFn. The files are walked in lexical order, which makes the output deterministic but means that for very large directories Walk can be inefficient. Walk does not follow symbolic links. - That is from the standard library. I know. Their grammar teachers can not be happy with them right now.
func (f *Pkger) Walk(p string, wf filepath.WalkFunc) error {
	keys := f.files.Keys()

	pt, err := f.Parse(p)
	if err != nil {
		return err
	}

	skip := "!"
	for _, k := range keys {
		if k.Pkg != pt.Pkg {
			continue
		}
		if !strings.HasPrefix(k.Name, pt.Name) {
			continue
		}
		if strings.HasPrefix(k.Name, skip) {
			continue
		}

		fl, ok := f.files.Load(k)
		if !ok {
			return os.ErrNotExist
		}

		fi, err := fl.Stat()
		if err != nil {
			return err
		}

		fi = pkging.WithName(strings.TrimPrefix(k.Name, pt.Name), fi)

		err = wf(k.String(), fi, nil)
		if err == filepath.SkipDir {

			skip = k.Name
			continue
		}

		if err != nil {
			return err
		}
	}
	return nil
}
