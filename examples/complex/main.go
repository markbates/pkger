package main

import (
	"fmt"
	"log"

	"github.com/markbates/pkger/examples/complex/api"
)

const host = ":3000"

func main() {
	v, err := api.Version()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Version: ", v)

	if err := writeFile(); err != nil {
		log.Fatal(err)
	}

	if err := readFile(); err != nil {
		log.Fatal(err)
	}
}
