package pkger

import (
	"bytes"
	"compress/gzip"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
)

// Open opens the named file for reading.
func Open(p string) (*File, error) {
	pt, err := Parse(p)
	if err != nil {
		return nil, err
	}
	return rootIndex.Open(pt)
}

// Create creates the named file with mode 0666 (before umask), truncating it if it already exists. If successful, methods on the returned File can be used for I/O; the associated file descriptor has mode O_RDWR. If there is an error, it will be of type *PathError.
func Create(p string) (*File, error) {
	pt, err := Parse(p)
	if err != nil {
		return nil, err
	}
	return rootIndex.Create(pt)
}

func Unpack(ind string) error {
	b, err := hex.DecodeString(ind)
	if err != nil {
		log.Fatal("hex.DecodeString", err)
		return err
	}

	gz, err := gzip.NewReader(bytes.NewReader(b))
	if err != nil {
		log.Fatal("gzip.NewReader", err)
		return err
	}
	defer gz.Close()

	if err := json.NewDecoder(gz).Decode(rootIndex); err != nil {
		log.Fatal("json.NewDecoder", err)
		return err
	}

	fmt.Println(rootIndex.Files.Keys())

	return nil
}

func Pack(out io.Writer, paths []Path) error {
	bb := &bytes.Buffer{}
	gz := gzip.NewWriter(bb)
	defer gz.Close()

	for _, p := range paths {
		f, err := Open(p.String())
		if err != nil {
			return err
		}
		fi, err := f.Stat()
		if err != nil {
			return err
		}
		if fi.IsDir() {
			continue
		}
		rootIndex.Files.Store(p, f)
		f.Close()
	}

	if err := json.NewEncoder(gz).Encode(rootIndex); err != nil {
		return err
	}
	if err := gz.Close(); err != nil {
		return err
	}
	s := hex.EncodeToString(bb.Bytes())
	fmt.Fprint(out, s)
	return nil
}
