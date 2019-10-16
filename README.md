# Pkger

[`github.com/markbates/pkger`](https://godoc.org/github.com/markbates/pkger) is a tool for embedding static files into Go binaries. It will, hopefully, be a replacement for [`github.com/gobuffalo/packr/v2`](https://godoc.org/github.com/gobuffalo/packr/v2).

## How it Works

Pkger is powered by the dark magic of Go Modules, so they're like, totally required.

With Go Modules pkger can resolve packages with accuracy. No more guessing and trying to
figure out build paths, GOPATHS, etc... for this tired old lad.

With the module's path correctly resolved, it can serve as the "root" directory for that
module, and all files in that module's directory are available.

Paths:
* Paths should use UNIX style paths:
  `/cmd/pkger/main.go`
* If unspecified the path's package is assumed to be the current module.
* Packages can specified in at the beginning of a path with a `:` seperator.
github.com/markbates/pkger:/cmd/pkger/main.go
* There are no relative paths. All paths are absolute to the modules root.

```
"github.com/gobuffalo/buffalo:/go.mod" => $GOPATH/pkg/mod/github.com/gobuffalo/buffalo@v0.14.7/go.mod
```

## Usage

Pkger's API is modeled on that of the [`os`](https://godoc.org/os) package in Go's standard library. This makes Pkger usage familiar to Go developers.



```go
type Pkger interface {
  Parse(p string) (Path, error)
  Abs(p string) (string, error)
  AbsPath(Path) (string, error)
  Current() (here.Info, error)
  Info(p string) (here.Info, error)
  Create(name string) (File, error)
  MkdirAll(p string, perm os.FileMode) error
  Open(name string) (File, error)
  Stat(name string) (os.FileInfo, error)
  Walk(p string, wf filepath.WalkFunc) error
  Remove(name string) error
  RemoveAll(path string) error
}

type File interface {
  Close() error
  Abs() (string, error)
  Info() here.Info
  Name() string
  Open(name string) (http.File, error)
  Path() Path
  Read(p []byte) (int, error)
  Readdir(count int) ([]os.FileInfo, error)
  Seek(offset int64, whence int) (int64, error)
  Stat() (os.FileInfo, error)
  Write(b []byte) (int, error)
}
```

```bash
├── go.mod
├── go.sum
├── main.go
├── public
│   ├── images
│   │   ├── mark-small.png
│   │   ├── mark.png
│   │   ├── mark_250px.png
│   │   └── mark_400px.png
│   └── index.html
└── templates
    ├── a.txt
    └── b
        └── b.txt
```

```go
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
```

### Output Without Packing

```bash
# compile the go binary as usual and run the app:
$ go build -v; ./app

{
  "Dir": "$GOPATH/src/github.com/markbates/pkger/examples/app",
  "ImportPath": "github.com/markbates/pkger/examples/app",
  "Imports": [
    "fmt",
    "github.com/markbates/pkger",
    "log",
    "net/http",
    "os"
  ],
  "Module": {
    "Dir": "$GOPATH/src/github.com/markbates/pkger/examples/app",
    "GoMod": "$GOPATH/src/github.com/markbates/pkger/examples/app/go.mod",
    "GoVersion": "1.13",
    "Main": true,
    "Path": "github.com/markbates/pkger/examples/app"
  },
  "Name": "main"
}
Walking files for github.com/markbates/pkger/examples/app
>>  github.com/gobuffalo/buffalo:/render
>>  github.com/gobuffalo/buffalo:/render/auto.go
>>  github.com/gobuffalo/buffalo:/render/auto_test.go
>>  github.com/gobuffalo/buffalo:/render/download.go
>>  github.com/gobuffalo/buffalo:/render/download_test.go
>>  github.com/gobuffalo/buffalo:/render/func.go
>>  github.com/gobuffalo/buffalo:/render/func_test.go
>>  github.com/gobuffalo/buffalo:/render/helpers.go
>>  github.com/gobuffalo/buffalo:/render/html.go
>>  github.com/gobuffalo/buffalo:/render/html_test.go
>>  github.com/gobuffalo/buffalo:/render/js.go
>>  github.com/gobuffalo/buffalo:/render/js_test.go
>>  github.com/gobuffalo/buffalo:/render/json.go
>>  github.com/gobuffalo/buffalo:/render/json_test.go
>>  github.com/gobuffalo/buffalo:/render/markdown_test.go
>>  github.com/gobuffalo/buffalo:/render/options.go
>>  github.com/gobuffalo/buffalo:/render/partials_test.go
>>  github.com/gobuffalo/buffalo:/render/plain.go
>>  github.com/gobuffalo/buffalo:/render/plain_test.go
>>  github.com/gobuffalo/buffalo:/render/render.go
>>  github.com/gobuffalo/buffalo:/render/render_test.go
>>  github.com/gobuffalo/buffalo:/render/renderer.go
>>  github.com/gobuffalo/buffalo:/render/sse.go
>>  github.com/gobuffalo/buffalo:/render/string.go
>>  github.com/gobuffalo/buffalo:/render/string_map.go
>>  github.com/gobuffalo/buffalo:/render/string_map_test.go
>>  github.com/gobuffalo/buffalo:/render/string_test.go
>>  github.com/gobuffalo/buffalo:/render/template.go
>>  github.com/gobuffalo/buffalo:/render/template_engine.go
>>  github.com/gobuffalo/buffalo:/render/template_helpers.go
>>  github.com/gobuffalo/buffalo:/render/template_helpers_test.go
>>  github.com/gobuffalo/buffalo:/render/template_test.go
>>  github.com/gobuffalo/buffalo:/render/xml.go
>>  github.com/gobuffalo/buffalo:/render/xml_test.go

Serving "github.com/markbates/pkger/examples/app:/public" on :3000
2019/09/22 14:07:41 GET /
2019/09/22 14:07:41 GET /images/mark.png
```

### Output With Packing

```bash
# run the pkger cli to generate a pkged.go file:
$ pkger

# compile the go binary as usual and run the app:
$ go build -v; ./app
{
  "Dir": "$GOPATH/src/github.com/markbates/pkger/examples/app",
  "ImportPath": "github.com/markbates/pkger/examples/app",
  "Imports": [
    "fmt",
    "github.com/markbates/pkger",
    "log",
    "net/http",
    "os"
  ],
  "Module": {
    "Dir": "$GOPATH/src/github.com/markbates/pkger/examples/app",
    "GoMod": "$GOPATH/src/github.com/markbates/pkger/examples/app/go.mod",
    "GoVersion": "1.13",
    "Main": true,
    "Path": "github.com/markbates/pkger/examples/app"
  },
  "Name": "main"
}
Walking files for github.com/markbates/pkger/examples/app
>>  github.com/gobuffalo/buffalo:/render
>>  github.com/gobuffalo/buffalo:/render/auto.go
>>  github.com/gobuffalo/buffalo:/render/auto_test.go
>>  github.com/gobuffalo/buffalo:/render/download.go
>>  github.com/gobuffalo/buffalo:/render/download_test.go
>>  github.com/gobuffalo/buffalo:/render/func.go
>>  github.com/gobuffalo/buffalo:/render/func_test.go
>>  github.com/gobuffalo/buffalo:/render/helpers.go
>>  github.com/gobuffalo/buffalo:/render/html.go
>>  github.com/gobuffalo/buffalo:/render/html_test.go
>>  github.com/gobuffalo/buffalo:/render/js.go
>>  github.com/gobuffalo/buffalo:/render/js_test.go
>>  github.com/gobuffalo/buffalo:/render/json.go
>>  github.com/gobuffalo/buffalo:/render/json_test.go
>>  github.com/gobuffalo/buffalo:/render/markdown_test.go
>>  github.com/gobuffalo/buffalo:/render/options.go
>>  github.com/gobuffalo/buffalo:/render/partials_test.go
>>  github.com/gobuffalo/buffalo:/render/plain.go
>>  github.com/gobuffalo/buffalo:/render/plain_test.go
>>  github.com/gobuffalo/buffalo:/render/render.go
>>  github.com/gobuffalo/buffalo:/render/render_test.go
>>  github.com/gobuffalo/buffalo:/render/renderer.go
>>  github.com/gobuffalo/buffalo:/render/sse.go
>>  github.com/gobuffalo/buffalo:/render/string.go
>>  github.com/gobuffalo/buffalo:/render/string_map.go
>>  github.com/gobuffalo/buffalo:/render/string_map_test.go
>>  github.com/gobuffalo/buffalo:/render/string_test.go
>>  github.com/gobuffalo/buffalo:/render/template.go
>>  github.com/gobuffalo/buffalo:/render/template_engine.go
>>  github.com/gobuffalo/buffalo:/render/template_helpers.go
>>  github.com/gobuffalo/buffalo:/render/template_helpers_test.go
>>  github.com/gobuffalo/buffalo:/render/template_test.go
>>  github.com/gobuffalo/buffalo:/render/xml.go
>>  github.com/gobuffalo/buffalo:/render/xml_test.go

Serving "github.com/markbates/pkger/examples/app:/public" on :3000
2019/09/22 14:07:41 GET /
2019/09/22 14:07:41 GET /images/mark.png
```
