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
	dir, err := pkger.Open("/public")
	if err != nil {
		log.Fatal("1", err)
	}
	defer dir.Close()

	fmt.Println(dir.Path())

	go func() {
		time.Sleep(1 * time.Second)
		res, err := http.Get("http://127.0.0.1:3000/assets/radio.radio")
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

	log.Fatal(http.ListenAndServe(":3000", http.StripPrefix("/assets/", http.FileServer(dir))))
}
