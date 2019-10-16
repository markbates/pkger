package pkgdiff

import (
	"bytes"
	"fmt"
	"os"
	"time"
)

type diffFile struct {
	Info os.FileInfo
	Body []byte
}

func (d diffFile) Bytes() ([]byte, error) {
	bb := &bytes.Buffer{}

	bb.WriteString(fmt.Sprintf("Name: %s\n", d.Info.Name()))
	bb.WriteString(fmt.Sprintf("Size: %d\n", d.Info.Size()))
	bb.WriteString(fmt.Sprintf("Mode: %d\n", d.Info.Mode()))
	bb.WriteString(fmt.Sprintf("ModTime: %s\n", d.Info.ModTime().Format(time.RFC3339)))
	bb.WriteString(fmt.Sprintf("IsDir: %t\n", d.Info.IsDir()))
	// bb.WriteString(fmt.Sprintf("Sys: %s\n", time.Now()))
	bb.Write(d.Body)

	return bb.Bytes(), nil
}

func File(path string) error {

	return nil
}
