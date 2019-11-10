package actions

import (
	"fmt"
	"io"
	"os"

	"github.com/markbates/pkger"
)

func WalkTemplates(w io.Writer) error {
	return pkger.Walk("/templates", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		_, err = w.Write([]byte(fmt.Sprintf("%s\n", path)))
		return err
	})
}