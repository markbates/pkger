package embed

import (
	"bytes"
	"compress/gzip"
	"encoding/hex"
	"encoding/json"

	"github.com/markbates/pkger/here"
)

type Embedder interface {
	MarshalEmbed() ([]byte, error)
}

type Unembedder interface {
	UnmarshalEmbed([]byte) error
}

type Data struct {
	Infos map[string]here.Info `json:"infos"`
	Files map[string]File      `json:"files"`
	Here  here.Info            `json:"here"`
}

func (d *Data) MarshalEmbed() ([]byte, error) {
	bb := &bytes.Buffer{}
	gz := gzip.NewWriter(bb)
	defer gz.Close()
	if err := json.NewEncoder(gz).Encode(d); err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}
	s := hex.EncodeToString(bb.Bytes())
	return []byte(s), nil
}

func (d *Data) UnmarshalEmbed(in []byte) error {
	b := make([]byte, len(in))
	if _, err := hex.Decode(b, in); err != nil {
		return err
	}

	gz, err := gzip.NewReader(bytes.NewReader(b))
	if err != nil {
		return err
	}
	defer gz.Close()

	p := &Data{}
	if err := json.NewDecoder(gz).Decode(p); err != nil {
		return err
	}
	(*d) = *p
	return nil
}
