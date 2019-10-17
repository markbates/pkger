package main

import (
	"log"
	"net/http"

	"github.com/markbates/pkger"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	f, err := pkger.Open("/public")
	if err != nil {
		return err
	}
	dir := http.FileServer(f)
	return http.ListenAndServe(":3000", dir)
}
