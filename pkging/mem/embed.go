package mem

import (
	"bytes"
	"compress/gzip"
	"encoding/hex"
	"encoding/json"
)

func (pkg *Pkger) MarshalEmbed() ([]byte, error) {
	bb := &bytes.Buffer{}
	gz := gzip.NewWriter(bb)
	defer gz.Close()
	if err := json.NewEncoder(gz).Encode(pkg); err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}
	s := hex.EncodeToString(bb.Bytes())
	return []byte(s), nil
}

func UnmarshalEmbed(in []byte) (*Pkger, error) {
	b := make([]byte, len(in))
	if _, err := hex.Decode(b, in); err != nil {
		return nil, err
	}

	gz, err := gzip.NewReader(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	defer gz.Close()

	pkg := &Pkger{}
	if err := json.NewDecoder(gz).Decode(pkg); err != nil {
		return nil, err
	}

	return pkg, nil
}
