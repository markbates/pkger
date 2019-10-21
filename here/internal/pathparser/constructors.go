package pathparser

import (
	"fmt"
	"strings"
)

func toString(i interface{}) (string, error) {
	if i == nil {
		return "", nil
	}
	if s, ok := i.(string); ok {
		return s, nil
	}
	return "", fmt.Errorf("%T is not a string", i)
}

func toName(i interface{}) (string, error) {
	s, err := toString(i)
	if err != nil {
		return "", err
	}
	if !strings.HasPrefix(s, "/") {
		s = "/" + s
	}
	return s, nil
}

func toPath(pkg, name interface{}) (*Path, error) {
	n, err := toString(name)
	if err != nil {
		return nil, err
	}

	pg, _ := pkg.(*Package)
	p := &Path{
		Name: n,
		Pkg:  pg,
	}
	// if p.IsZero() {
	// 	return nil, fmt.Errorf("empty path")
	// }

	if p.Name == "" {
		p.Name = "/"
	}
	return p, nil
}

func toPackage(n, v interface{}) (*Package, error) {
	name, err := toString(n)
	if err != nil {
		return nil, err
	}

	var version string
	if s, ok := v.(string); ok {
		version = s
	}

	return &Package{
		Name:    name,
		Version: version,
	}, nil
}
