package memwh

import (
	"encoding/json"
	"fmt"

	"github.com/markbates/pkger/fs"
)

func (f File) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{
		"info":   f.info,
		"her":    f.her,
		"path":   f.path,
		"data":   f.data,
		"parent": f.parent,
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

	f.info = &fs.FileInfo{}
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
