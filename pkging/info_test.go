package pkging

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_NewFileInfo(t *testing.T) {

	in := []string{
		"/public/images/img1.png",
		"public/images/img1.png",
		"/public\\images/img1.png",
		"public/images\\img1.png",
		"\\public\\images\\img1.png",
		"public\\images\\img1.png",
		"\\public/images\\img1.png",
		"public\\images/img1.png",
		"\\public\\images\\img1.png",
	}

	const exp = "/public/images/img1.png"
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

func Test_WithName(t *testing.T) {

	f1 := &FileInfo{
		Details: Details{
			Name:  "/walls/crumbling",
			Size:  42,
			Mode:  os.FileMode(0644),
			IsDir: true,
		},
	}

	const exp = "/public/images/img1.png"
	in := []string{
		"/public/images/img1.png",
		"public/images/img1.png",
		"/public\\images/img1.png",
		"public/images\\img1.png",
		"\\public\\images\\img1.png",
		"public\\images\\img1.png",
		"\\public/images\\img1.png",
		"public\\images/img1.png",
		"\\public\\images\\img1.png",
	}

	for _, n := range in {
		t.Run(n, func(st *testing.T) {
			r := require.New(st)

			f2 := WithName(n, f1)

			r.Equal(exp, f2.Name())
		})
	}
}
