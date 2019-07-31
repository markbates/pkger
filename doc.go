/*
This package is WIP. Please use/test/try and report issues, but be careful with production. OK?

Pkger is powered by the dark magic of Go Modules, so they're like, totally required.

With Go Modules pkger can resolve packages with accuracy. No more guessing and trying to
figure out build paths, GOPATHS, etc... for this tired old lad.

With the module's path correctly resolved, it can serve as the "root" directory for that
module, and all files in that module's directory are available.

		Paths:
		* Paths should use UNIX style paths:
			/cmd/pkger/main.go
		* If unspecified the path's package is assumed to be the current module.
		* Packages can specified in at the beginning of a path with a `:` seperator.
			github.com/markbates/pkger:/cmd/pkger/main.go

		"github.com/gobuffalo/buffalo:/go.mod" => $GOPATH/pkg/mod/github.com/gobuffalo/buffalo@v0.14.7/go.mod
*/
package pkger
