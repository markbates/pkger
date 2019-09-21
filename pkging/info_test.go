package pkging

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_NewFileInfo(t *testing.T) {
	const exp = "/public/images/mark.png"

	in := []string{
		"/public/images/mark.png",
		"public/images/mark.png",
		"/public\\images/mark.png",
		"public/images\\mark.png",
		"\\public\\images\\mark.png",
		"public\\images\\mark.png",
		"\\public/images\\mark.png",
		"public\\images/mark.png",
		"\\public\\images\\mark.png",
	}

	for _, n := range in {
		t.Run(n, func(st *testing.T) {
			r := require.New(st)

			f1 := &FileInfo{
				Details: Details{
					Name:  n,
					Size:  42,
					Mode:  os.FileMode(0644),
					IsDir: true,
				},
			}

			f2 := NewFileInfo(f1)

			r.Equal(exp, f2.Name())
			r.Equal(f1.Size(), f2.Size())
			r.Equal(f1.Mode(), f2.Mode())
			r.Equal(f1.IsDir(), f2.IsDir())
		})
	}

}
