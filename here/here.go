package here

import (
	"bytes"
	"os"
	"os/exec"

	"github.com/markbates/pkger/internal/debug"
)

func run(n string, args ...string) ([]byte, error) {
	c := exec.Command(n, args...)

	debug.Debug("[HERE] %s", c.Args)

	bb := &bytes.Buffer{}
	c.Stdout = bb
	c.Stderr = os.Stderr
	err := c.Run()
	if err != nil {
		return nil, err
	}

	return bb.Bytes(), nil
}
