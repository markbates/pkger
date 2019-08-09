package pkger

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/markbates/pkger/here"
)

type jason struct {
	Files       *filesMap `json:"files"`
	Infos       *infosMap `json:"infos"`
	Paths       *pathsMap `json:"paths"`
	CurrentInfo here.Info `json:"current_info"`
}

func (f File) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{}
	m["info"] = f.info
	m["her"] = f.her
	m["path"] = f.path
	m["data"] = f.data
	m["parent"] = f.parent
	if !f.info.virtual {
		if len(f.data) == 0 && !f.info.IsDir() {
			b, err := ioutil.ReadAll(&f)
			if err != nil {
				return nil, err
			}
			m["data"] = b
		}
	}

	return json.Marshal(m)
}

func (f *File) UnmarshalJSON(b []byte) error {
	m := map[string]json.RawMessage{}
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}

	info, ok := m["info"]
	if !ok {
		return fmt.Errorf("missing info")
	}

	f.info = &FileInfo{}
	if err := json.Unmarshal(info, f.info); err != nil {
		return err
	}

	her, ok := m["her"]
	if !ok {
		return fmt.Errorf("missing her")
	}
	if err := json.Unmarshal(her, &f.her); err != nil {
		return err
	}

	path, ok := m["path"]
	if !ok {
		return fmt.Errorf("missing path")
	}
	if err := json.Unmarshal(path, &f.path); err != nil {
		return err
	}

	parent, ok := m["parent"]
	if !ok {
		return fmt.Errorf("missing parent")
	}
	if err := json.Unmarshal(parent, &f.parent); err != nil {
		return err
	}

	if err := json.Unmarshal(m["data"], &f.data); err != nil {
		return err
	}

	return nil
}
