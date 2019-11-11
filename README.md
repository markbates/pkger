# Pkger

[`github.com/markbates/pkger`](https://godoc.org/github.com/markbates/pkger) is a tool for embedding static files into Go binaries. It will, hopefully, be a replacement for [`github.com/gobuffalo/packr/v2`](https://godoc.org/github.com/gobuffalo/packr/v2).

## How it Works (Module Aware Pathing)

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

## CLI

### Installation

```bash
$ go get github.com/markbates/pkger/cmd/pkger
$ pkger -h
```

### Usage

```bash
$ pkger
```

The result will be a `pkged.go` file in the **root** of the module with the embedded information and the package name of the module.

```go
// ./pkged.go
package <.>

// Pkger stuff here
```

The `-o` flag can be used specify the directory of the `pkged.go` file.

```bash
$ pkger -o cmd/reader
```

The result will be a `pkged.go` file in the **cmd/reader** folder with the embedded information and the package name of that folder.

```go
// cmd/reader/pkged.go
package <reader>

// Pkger stuff here
```

## Reference Application

The reference application for the `README` examples, as well as all testing, can be found at [https://github.com/markbates/pkger/tree/master/pkging/pkgtest/testdata/ref](https://github.com/markbates/pkger/tree/master/pkging/pkgtest/testdata/ref).

```
├── actions
│   └── actions.go
├── assets
│   ├── css
│   │   ├── _buffalo.scss
│   │   └── application.scss
│   ├── images
│   │   ├── favicon.ico
│   │   └── logo.svg
│   └── js
│       └── application.js
├── go.mod
├── go.sum
├── locales
│   └── all.en-us.yaml
├── main.go
├── mod
│   └── mod.go
├── models
│   └── models.go
├── public
│   ├── assets
│   │   └── app.css
│   ├── images
│   │   └── img1.png
│   ├── index.html
│   └── robots.txt
├── templates
│   ├── _flash.plush.html
│   ├── application.plush.html
│   └── index.plush.html
└── web
    └── web.go

13 directories, 20 files
```


## API Usage

Pkger's API is modeled on that of the [`os`](https://godoc.org/os) package in Go's standard library. This makes Pkger usage familiar to Go developers.

The two most important interfaces are [`github.com/markbates/pkger/pkging#Pkger`](https://godoc.org/github.com/markbates/pkger/pkging#Pkger) and [`github.com/markbates/pkger/pkging#File`](https://godoc.org/github.com/markbates/pkger/pkging#File).

```go
type Pkger interface {
  Parse(p string) (Path, error)
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

These two interfaces, along with the [`os#FileInfo`](https://godoc.org/os#FileInfo), provide the bulk of the API surface area.

### Open

```go
func run() error {
	f, err := pkger.Open("/public/index.html")
	if err != nil {
		return err
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return err
	}

	fmt.Println("Name: ", info.Name())
	fmt.Println("Size: ", info.Size())
	fmt.Println("Mode: ", info.Mode())
	fmt.Println("ModTime: ", info.ModTime())

	if _, err := io.Copy(os.Stdout, f); err != nil {
		return err
	}
	return nil
}
```

### Stat

```go
func run() error {
	info, err := pkger.Stat("/public/index.html")
	if err != nil {
		return err
	}

	fmt.Println("Name: ", info.Name())
	fmt.Println("Size: ", info.Size())
	fmt.Println("Mode: ", info.Mode())
	fmt.Println("ModTime: ", info.ModTime())

	return nil
}
```

### Walk

```go
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

## Understanding the Parser

The [`github.com/markbates/pkger/parser#Parser`](https://godoc.org/github.com/markbates/pkger/parser#Parser) works by statically analyzing the source code of your module using the [`go/parser`](https://godoc.org/go/parser) to find a selection of declarations.

The following declarations in your source code will tell the parser to embed files or folders.

* `pkger.Dir("<path>")` - Embeds all files under the specified path.
* `pkger.Open("<path>")` - Embeds the file, or folder, of the specified path.
* `pkger.Stat("<path>")` - Embeds the file, or folder, of the specified path.
* `pkger.Walk("<path>", filepath.WalkFunc)` - Embeds all files under the specified path.
* `pkger.Include("<path>")` - `Include` is a no-op that directs the pkger tool to include the desired file or folder.


