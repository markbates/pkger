module github.com/markbates/pkger/internal/examples/app

go 1.12

require (
	cloud.google.com/go v0.36.0 // indirect
	github.com/cockroachdb/cockroach-go v0.0.0-20181001143604-e0a95dfd547c // indirect
	github.com/gobuffalo/buffalo v0.14.8 // indirect
	github.com/markbates/pkger v0.0.0-20190803203656-a4a55a52dc5d
)

replace github.com/markbates/pkger => ../../../
