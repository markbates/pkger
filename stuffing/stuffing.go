package stuffing

import (
	"fmt"
	"io"
	"strings"

	"github.com/markbates/pkger/here"
	"github.com/markbates/pkger/pkging"
	"github.com/markbates/pkger/pkging/mem"
	"github.com/markbates/pkger/pkging/stdos"
)

func Stuff(w io.Writer, cur here.Info, paths []here.Path) error {
	disk, err := stdos.New(cur)
	if err != nil {
		return err
	}

	pkg, err := mem.New(cur)
	if err != nil {
		return err
	}

	for _, pt := range paths {
		err = func() error {
			f, err := disk.Open(pt.String())
			if err != nil {
				return err
			}
			defer f.Close()

			fi, err := f.Stat()
			if err != nil {
				return err
			}
			info := f.Info()
			fmt.Println(">>>TODO stuffing/stuffing.go:41: info.Dir ", info.Dir)
			fi = pkging.WithName(strings.TrimPrefix(fi.Name(), info.Dir), fi)
			fmt.Println(">>>TODO stuffing/stuffing.go:37: fi.Name() ", fi.Name())

			if err := pkg.Add(fi, f); err != nil {
				return err
			}

			return nil
			// WithInfo(ng, og)
		}()
		if err != nil {
			return err
		}
	}

	b, err := pkg.MarshalEmbed()
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}
