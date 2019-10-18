package main

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"
	"time"

	"github.com/markbates/pkger"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 0, ' ', tabwriter.Debug)
	defer w.Flush()

	return pkger.Walk("/public", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		fmt.Fprintf(w,
			"%s \t %d \t %s \t %s \t\n",
			info.Name(),
			info.Size(),
			info.Mode(),
			info.ModTime().Format(time.RFC3339),
		)

		return nil
	})

}
