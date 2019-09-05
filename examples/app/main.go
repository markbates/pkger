package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/markbates/pkger"
)

func main() {
	mux := http.NewServeMux()

	pub, err := pkger.Open(":/public")
	if err != nil {
		log.Fatal(err)
	}
	defer pub.Close()

	fmt.Println(pub.Path())

	mux.Handle("/t", http.StripPrefix("/t", tmplHandler()))
	mux.Handle("/", http.FileServer(pub))

	log.Fatal(http.ListenAndServe(":3000", mux))
}

func tmplHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		t, err := pkger.Open(":/templates/a.txt")
		if err != nil {
			http.Error(res, err.Error(), 500)
		}
		defer t.Close()

		io.Copy(res, t)
	}
}
