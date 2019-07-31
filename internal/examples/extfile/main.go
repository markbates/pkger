package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/markbates/pkger"
)

func main() {
	f, err := pkger.Open("github.com/gobuffalo/buffalo:/go.mod")
	if err != nil {
		log.Fatal("1", err)
	}

	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		log.Fatal("2", err)
	}

	fmt.Println(fi)

	io.Copy(os.Stdout, f)
}
