package pkgutil

import (
	"io"
	"os"

	"github.com/markbates/pkger/here"
	"github.com/markbates/pkger/parser"
	"github.com/markbates/pkger/pkging/mem"
	"github.com/markbates/pkger/pkging/stdos"
)

func Stuff(w io.Writer, c here.Info, decls parser.Decls) error {
	disk, err := stdos.New(c)
	if err != nil {
		return err
	}

	pkg, err := mem.New(c)
	if err != nil {
		return err
	}

	files, err := decls.Files()
	if err != nil {
		return err
	}

	for _, pf := range files {
		err = func() error {
			df, err := disk.Open(pf.Path.String())
			if err != nil {
				return err
			}
			defer df.Close()

			info, err := df.Stat()
			if err != nil {
				return err
			}

			if err := pkg.Add(df); err != nil {
				return err
			}

			if !info.IsDir() {
				return nil
			}

			err = disk.Walk(df.Path().String(), func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if info.IsDir() {
					return nil
				}

				f, err := disk.Open(path)
				if err != nil {
					return err
				}
				defer f.Close()

				if err := pkg.Add(f); err != nil {
					return err
				}

				return nil
			})

			return err
		}()

		if err != nil {
			return err
		}
		b, err := pkg.MarshalEmbed()
		if err != nil {
			return err
		}
		_, err = w.Write(b)
	}
	return nil
}
