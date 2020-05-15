package pkgutil

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/markbates/pkger/here"
	"github.com/markbates/pkger/parser"
	"github.com/markbates/pkger/pkging/embed"
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
			if strings.HasSuffix(pf.Abs, ".tmp") {
				return nil
			}
			df, err := os.Open(pf.Abs)
			if err != nil {
				return fmt.Errorf("could open stuff %s: %s", pf.Abs, err)
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
	}

	b, err := pkg.MarshalJSON()
	if err != nil {
		return err
	}

	b, err = embed.Encode(b)
	if err != nil {
		return err
	}

	_, err = w.Write(b)
	return nil
}
