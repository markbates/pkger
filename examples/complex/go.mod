module github.com/markbates/pkger/examples/complex

go 1.13

require (
	github.com/markbates/pkger v0.0.0
	github.com/markbates/pkger/examples/complex/api v0.0.0
)

replace github.com/markbates/pkger => ../../

replace github.com/markbates/pkger/examples/complex/api => ./api
