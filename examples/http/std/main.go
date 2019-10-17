package main

import (
	"log"
	"net/http"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	dir := http.FileServer(http.Dir("./public"))
	return http.ListenAndServe(":3000", dir)
}
