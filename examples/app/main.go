package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/markbates/pkger"
)

const host = ":3000"

func main() {
	// get the currently running application's here.Info.
	// this contains really, really, really useful information
	// about your application, check it out. :)
	// we don't need it for this example, but i thought it could
	// be good to show.
	current, err := pkger.Current()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(current)

	fmt.Printf("Walking files for %s\n", current.ImportPath)
	// walk the files in this module. "/" is where the `go.mod` for this module is
	err = pkger.Walk("github.com/gobuffalo/buffalo:/render", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		fmt.Println(">> ", path)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	// find the public directory with using the full pkger path <pkg:/path> to it:
	// 	pkg - is the module/package you want to get a file from
	// 		if pkg is empty then it is assumed to be current.ImportPath
	// 	: - seperator between the module/package name, pkg, and the "file path"
	// 	path - this is the ABSOLUTE path to the file/directory you want, as relative
	// 	to the root of the module/package's go.mod file.
	dir, err := pkger.Open("github.com/markbates/pkger/examples/app:/public")
	if err != nil {
		log.Fatal(err)
	}
	// don't forget to close the file later
	defer dir.Close()

	fmt.Printf("\nServing %q on %s\n", dir.Path(), host)

	// serve the public directory on the host (":3000")
	// just like using the os package you still need to use
	// http.FileServer to serve a directory.
	// you DON'T, however, need to use http.Dir all pkger files
	// already implement that interface.
	log.Fatal(http.ListenAndServe(host, logger(http.FileServer(dir))))
}

// logger will print out the requests as they come in, otherwise its a blank
// screen, and that's no fun.
func logger(h http.Handler) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		log.Println(req.Method, req.URL.String())
		h.ServeHTTP(res, req)
	}
}
