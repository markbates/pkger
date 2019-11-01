package here

import (
	"encoding/json"
)

type Module struct {
	Path      string `json:"path"`
	Main      bool   `json:"main"`
	Dir       string `json:"dir"`
	GoMod     string `json:"go_mod"`
	GoVersion string `json:"go_version"`
}

func (i Module) String() string {
	b, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		return err.Error()
	}
	return string(b)
}

func (i Module) IsZero() bool {
	return i.String() == Module{}.String()
}
