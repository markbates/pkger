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
```

```go
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
```

### Output Without Packing

```bash
# compile the go binary as usual and run the app:
$ go build -v; ./app
/public                       | 128      | drwxr-xr-x | 2019-10-17T15:02:57-04:00 |
/public/images                | 192      | drwxr-xr-x | 2019-10-17T15:02:57-04:00 |
/public/images/mark-small.png | 649549   | -rw-r--r-- | 2019-10-17T15:02:56-04:00 |
/public/images/mark.png       | 50401191 | -rw-r--r-- | 2019-10-17T15:02:57-04:00 |
/public/images/mark_250px.png | 27718    | -rw-r--r-- | 2019-10-17T15:02:57-04:00 |
/public/images/mark_400px.png | 63543    | -rw-r--r-- | 2019-10-17T15:02:57-04:00 |
/public/index.html            | 257      | -rw-r--r-- | 2019-10-17T15:02:57-04:00 |
```

### Output With Packing

```bash
# run the pkger cli to generate a pkged.go file:
$ pkger

# compile the go binary as usual and run the app:
$ go build -v; ./app
/                      | 128      | drwxr-xr-x | 2019-10-17T15:02:57-04:00 |
/images                | 192      | drwxr-xr-x | 2019-10-17T15:02:57-04:00 |
/images/mark-small.png | 649549   | -rw-r--r-- | 2019-10-17T15:02:56-04:00 |
/images/mark.png       | 50401191 | -rw-r--r-- | 2019-10-17T15:02:57-04:00 |
/images/mark_250px.png | 27718    | -rw-r--r-- | 2019-10-17T15:02:57-04:00 |
/images/mark_400px.png | 63543    | -rw-r--r-- | 2019-10-17T15:02:57-04:00 |
/index.html            | 257      | -rw-r--r-- | 2019-10-17T15:02:57-04:00 |
```
