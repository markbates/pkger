package pathparser

type Path struct {
	Pkg  *Package
	Name string
}

func (p Path) IsZero() bool {
	return p.Pkg == nil && len(p.Name) == 0
}

type Package struct {
	Name    string
	Version string
}
