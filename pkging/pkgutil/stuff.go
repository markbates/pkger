package pkgutil

import (
	"io"
	"os"

	"github.com/markbates/pkger/here"
	"github.com/markbates/pkger/parser"
	"github.com/markbates/pkger/pkging/mem"
)

func Stuff(w io.Writer, c here.Info, decls parser.Decls) error {
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
			df, err := os.Open(pf.Abs)
			if err != nil {
				return err
			}
			defer df.Close()

			if err := pkg.Add(df); err != nil {
				return err
			}

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
