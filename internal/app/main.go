package main

import (
	"io"
	"log"
	"net/http"

	"github.com/markbates/pkger"
)

func main() {
	mux := http.NewServeMux()

	pub, err := pkger.Open("/internal/app/public")
	if err != nil {
		log.Fatal(err)
	}
	mux.Handle("/t", http.StripPrefix("/t", tmplHandler()))
	mux.Handle("/logo", http.StripPrefix("/logo", logoHandler()))
	mux.Handle("/", http.FileServer(pub))

	log.Fatal(http.ListenAndServe(":3000", mux))
}

func tmplHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		t, err := pkger.Open("/internal/app/templates/a.txt")
		if err != nil {
			http.Error(res, err.Error(), 500)
		}
		defer t.Close()

		io.Copy(res, t)
	}
}

func logoHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		t, err := pkger.Open("github.com/gobuffalo/buffalo:/logo.svg")
		if err != nil {
			http.Error(res, err.Error(), 500)
		}
		defer t.Close()

		io.Copy(res, t)
	}
}
