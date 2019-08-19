package pkger

import (
	"bytes"
	"compress/gzip"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"sync"

	"github.com/markbates/pkger/here"
	"github.com/markbates/pkger/internal/debug"
)

var filesCache = &filesMap{}
var infosCache = &infosMap{}
var pathsCache = &pathsMap{}
var curOnce = &sync.Once{}
var currentInfo here.Info

var packed bool

var packMU = &sync.RWMutex{}

func ReadFile(s string) ([]byte, error) {
	f, err := Open(s)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}

func dubeg(key, format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	debug.Debug("[%s|%s] %s", key, s)
}

func Unpack(ind string) error {
	packed = true
	packMU.Lock()
	defer packMU.Unlock()
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
	jay.Files.Range(func(key Path, value *File) bool {
		filesCache.Store(key, value)
		return true
	})
	jay.Infos.Range(func(key string, value here.Info) bool {
		infosCache.Store(key, value)
		return true
	})
	jay.Paths.Range(func(key string, value Path) bool {
		pathsCache.Store(key, value)
		return true
	})
	currentInfo = jay.CurrentInfo
	return nil
}

func Pack(out io.Writer, paths []Path) error {
	packMU.RLock()
	defer packMU.RUnlock()
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
