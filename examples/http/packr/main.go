package main

import (
	"log"
	"net/http"

	"github.com/gobuffalo/packr/v2"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	box := packr.New("demo", "./public")

	dir := http.FileServer(box)
	return http.ListenAndServe(":3000", dir)
}
