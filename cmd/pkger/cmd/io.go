package cmd

import (
	"io"
	"os"
)

type IO interface {
	In() io.Reader
	Out() io.Writer
	Err() io.Writer
}

type stdio struct{}

func (stdio) In() io.Reader {
	return os.Stdin
}

func (stdio) Out() io.Writer {
	return os.Stdout
}

func (stdio) Err() io.Writer {
	return os.Stderr
}

func StdIO() IO {
	return stdio{}
}
