package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/markbates/pkger"
)

func serve(args []string) error {
	if len(args) == 0 {
		args = []string{"."}
	}
	f, err := pkger.Open(args[0])
	if err != nil {
		log.Fatal("1", err)
	}
	defer f.Close()

	fmt.Println(f.Path())

	return http.ListenAndServe(":3000", http.FileServer(f))
}
