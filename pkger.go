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

	var jay jason
	if err := json.NewDecoder(gz).Decode(&jay); err != nil {
		return err
	}

	filesCache = jay.Files
	infosCache = jay.Infos
	pathsCache = jay.Paths
	currentInfo = jay.CurrentInfo
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
			filesCache.Store(p, f)
			f.Close()
			continue
		}

		dubeg("Pack", "%s", p)
		filesCache.Store(p, f)
		f.Close()

	}

	jay := jason{
		Files:       filesCache,
		Infos:       infosCache,
		Paths:       pathsCache,
		CurrentInfo: currentInfo,
	}

	if err := json.NewEncoder(gz).Encode(jay); err != nil {
		return err
	}
	if err := gz.Close(); err != nil {
		return err
	}
	s := hex.EncodeToString(bb.Bytes())
	_, err := fmt.Fprint(out, s)
	return err
}
