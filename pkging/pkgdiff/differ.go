package pkgdiff

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/markbates/pkger/pkging/mem"
	"github.com/markbates/pkger/pkging/stdos"
)

const (
	BYTE = 1 << (10 * iota)
	KILOBYTE
	MEGABYTE
	GIGABYTE
	TERABYTE
	PETABYTE
	EXABYTE
)

type Differ struct {
	A *stdos.Pkger
}

func (d Differ) File(p string) (string, error) {

	f, err := d.A.Open(p)
	if err != nil {
		return "", err
	}

	info, err := f.Stat()
	if err != nil {
		return "", err
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}
	if _, err := f.Seek(0, 0); err != nil {
		return "", err
	}

	df := diffFile{
		Info: info,
		Body: b,
	}

	adiff, err := df.Bytes()
	if err != nil {
		return "", err
	}

	m, err := mem.New(d.A.Here)
	if err != nil {
		return "", err
	}

	if err := m.Add(f); err != nil {
		return "", err
	}

	if err := f.Close(); err != nil {
		return "", err
	}

	me, err := m.MarshalEmbed()
	if err != nil {
		return "", err
	}

	m, err = mem.UnmarshalEmbed(me)
	if err != nil {
		return "", err
	}

	bf, err := m.Open(p)
	if err != nil {
		return "", err
	}

	info, err = bf.Stat()
	if err != nil {
		return "", err
	}
	b, err = ioutil.ReadAll(bf)
	if err != nil {
		return "", err
	}

	if err := bf.Close(); err != nil {
		return "", err
	}

	df = diffFile{
		Info: info,
		Body: b,
	}

	bdiff, err := df.Bytes()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(cmp.Diff(adiff, bdiff)), nil
}

func (d Differ) Dir(p string) ([]string, error) {
	var diffs []string
	err := d.A.Walk(p, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if info.Size() > MEGABYTE {
			return nil
		}

		res, err := d.File(path)
		if err != nil {
			return err
		}
		if len(res) > 0 {
			diffs = append(diffs, res)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return diffs, nil
}
