package pkger

import (
	"github.com/markbates/pkger/here"
)

type jason struct {
	Files       *filesMap `json:"files"`
	Infos       *infosMap `json:"infos"`
	Paths       *pathsMap `json:"paths"`
	CurrentInfo here.Info `json:"current_info"`
}
