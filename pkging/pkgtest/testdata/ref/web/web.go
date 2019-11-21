package web

import (
	"net/http"

	"github.com/markbates/pkger"
)

func Serve() {
	pkger.Stat("github.com/stretchr/testify:/README.md")
	dir := http.FileServer(pkger.Dir("/public"))
	http.ListenAndServe(":3000", dir)
}
