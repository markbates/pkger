package main

import (
	"fmt"

	"github.com/markbates/pkger/parser"
)

func list(args []string) error {
	if len(args) == 0 {
		args = append(args, "")
	}

	for _, a := range args {
		res, err := parser.Parse(a)
		if err != nil {
			return err
		}

		fmt.Println(res.Path)

		for _, p := range res.Paths {
			fmt.Printf("  > %s\n", p)
		}
	}
	return nil
}
