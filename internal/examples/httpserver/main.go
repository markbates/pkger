package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/markbates/pkger"
)

func main() {
	f, err := pkger.Open("github.com/gobuffalo/buffalo")
	if err != nil {
		log.Fatal("1", err)
	}
	defer f.Close()

	fmt.Println(f.Path())

	go func() {
		time.Sleep(1 * time.Second)
		res, err := http.Get("http://127.0.0.1:3000/app.go")
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()

		_, err = io.Copy(os.Stdout, res.Body)
		if err != nil {
			log.Fatal(err)
		}

		if res.StatusCode >= 300 {
			log.Fatal("code: ", res.StatusCode)
		}

	}()

	log.Fatal(http.ListenAndServe(":3000", http.FileServer(f)))
}
