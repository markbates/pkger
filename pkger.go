package pkger

import (
	"log"

	"github.com/markbates/pkger/pkging"
	"github.com/markbates/pkger/pkging/stdos"
)

var current = func() pkging.Pkger {
	n, err := stdos.New()
	if err != nil {
		log.Fatal(err)
	}
	return n
}()

var Abs = current.Abs
var AbsPath = current.AbsPath
var Create = current.Create
var Current = current.Current
var Info = current.Info
var MkdirAll = current.MkdirAll
var Open = current.Open
var Parse = current.Parse
var Remove = current.Remove
var RemoveAll = current.RemoveAll
var Stat = current.Stat
var Walk = current.Walk
