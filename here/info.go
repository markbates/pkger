package here

import (
	"encoding/json"
)

// Info represents details about the directory/package
type Info struct {
	Dir        string
	ImportPath string
	Name       string
	Module     Module
}

func (i Info) MarshalJSON() ([]byte, error) {
	mm := map[string]interface{}{
		"ImportPath": i.ImportPath,
		"Name":       i.Name,
		"Module":     i.Module,
		"Dir":        i.Dir,
	}

	return json.Marshal(mm)
}

// IsZero checks if the type has been filled
// with rich chocolately data goodness
func (i Info) IsZero() bool {
	return i.String() == Info{}.String()
}

func (i Info) String() string {
	b, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		return err.Error()
	}
	s := string(b)
	return s
}
