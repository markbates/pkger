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

	pub, err := pkger.Open("/public")
	if err != nil {
		log.Fatal(err)
	}
	defer pub.Close()

	fi, err := pub.Stat()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(fi)

	mux.Handle("/t", http.StripPrefix("/t", tmplHandler()))
	mux.Handle("/logo", http.StripPrefix("/logo", logoHandler()))
	mux.Handle("/", http.FileServer(pub))

	// f, err := pkger.Open("/public/images/mark.png")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer f.Close()

	// lcl, err := os.Create("me.png")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// if _, err := io.Copy(lcl, f); err != nil {
	// 	log.Fatal(err)
	// }
	// lcl.Close()

	log.Fatal(http.ListenAndServe(":3000", mux))
}

func tmplHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		t, err := pkger.Open("/templates/a.txt")
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
