package here

import (
	"encoding/json"

	"github.com/markbates/hepa"
	"github.com/markbates/hepa/filters"
)

type Module struct {
	Path      string
	Main      bool
	Dir       string
	GoMod     string
	GoVersion string
}

func (m Module) MarshalJSON() ([]byte, error) {
	mm := map[string]interface{}{
		"Main":      m.Main,
		"GoVersion": m.GoVersion,
	}

	hep := hepa.New()
	hep = hepa.With(hep, filters.Home())
	hep = hepa.With(hep, filters.Golang())

	cm := map[string]string{
		"Path":  m.Path,
		"Dir":   m.Dir,
		"GoMod": m.GoMod,
	}

	for k, v := range cm {
		b, err := hep.Filter([]byte(v))
		if err != nil {
			return nil, err
		}
		mm[k] = string(b)
	}

	return json.Marshal(mm)
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
