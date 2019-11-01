package web

import (
	"net/http"

	"github.com/markbates/pkger"
)

func Serve() {
	dir := http.FileServer(pkger.Dir("/public"))
	http.ListenAndServe(":3000", dir)
}
