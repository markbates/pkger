package here

import (
	"encoding/json"
	"fmt"
	"os"
)

type Path struct {
	Pkg  string `json:"pkg"`
	Name string `json:"name"`
}

func (p Path) String() string {
	if p.Name == "" {
		p.Name = "/"
	}
	return fmt.Sprintf("%s:%s", p.Pkg, p.Name)
}

func (p Path) Format(st fmt.State, verb rune) {
	switch verb {
	case 'v':
		if st.Flag('+') {
			b, err := json.MarshalIndent(p, "", "  ")
			if err != nil {
				fmt.Fprint(os.Stderr, err)
				return
			}
			fmt.Fprint(st, string(b))
			return
		}
		fmt.Fprint(st, p.String())
	case 'q':
		fmt.Fprintf(st, "%q", p.String())
	default:
		fmt.Fprint(st, p.String())
	}
}
