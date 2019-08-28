module github.com/markbates/pkger/internal/examples/app

go 1.12

require (
	github.com/gobuffalo/buffalo v0.14.10 // indirect
	github.com/markbates/pkger v0.0.0-20190803203656-a4a55a52dc5d
	github.com/spf13/afero v1.2.1 // indirect
)

replace github.com/markbates/pkger => ../../../
