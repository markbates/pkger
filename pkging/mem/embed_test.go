package mem_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/markbates/pkger/here"
	"github.com/markbates/pkger/parser"
	"github.com/markbates/pkger/pkging/mem"
	"github.com/markbates/pkger/pkging/stuffing"
	"github.com/stretchr/testify/require"
)

func Test_Pkger_Embedding(t *testing.T) {
	r := require.New(t)

	pwd, err := os.Getwd()
	r.NoError(err)

	ch := filepath.Join(pwd,
		"..",
		"..",
		"internal",
		"testdata",
		"app")

	info := here.Info{
		Dir:        ch,
		ImportPath: "github.com/markbates/pkger/internal/testdata/app",
	}
	here.Cache(info.ImportPath, func(s string) (here.Info, error) {
		return info, nil
	})

	paths, err := parser.Parse(info)
	r.NoError(err)

	ps := make([]string, len(paths))
	for i, p := range paths {
		ps[i] = p.String()
	}

	r.Equal(inbed, ps)

	bb := &bytes.Buffer{}

	err = stuffing.Stuff(bb, info, paths)
	r.NoError(err)

	pkg, err := mem.UnmarshalEmbed(bb.Bytes())
	r.NoError(err)

	var res []string
	err = pkg.Walk("/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		res = append(res, path)
		return nil
	})

	r.NoError(err)
	r.Equal(inbed, res)
}

var inbed = []string{
	"github.com/gobuffalo/buffalo:/",
	"github.com/gobuffalo/buffalo:/render",
	"github.com/gobuffalo/buffalo:/render/auto.go",
	"github.com/gobuffalo/buffalo:/render/auto_test.go",
	"github.com/gobuffalo/buffalo:/render/download.go",
	"github.com/gobuffalo/buffalo:/render/download_test.go",
	"github.com/gobuffalo/buffalo:/render/func.go",
	"github.com/gobuffalo/buffalo:/render/func_test.go",
	"github.com/gobuffalo/buffalo:/render/helpers.go",
	"github.com/gobuffalo/buffalo:/render/html.go",
	"github.com/gobuffalo/buffalo:/render/html_test.go",
	"github.com/gobuffalo/buffalo:/render/js.go",
	"github.com/gobuffalo/buffalo:/render/js_test.go",
	"github.com/gobuffalo/buffalo:/render/json.go",
	"github.com/gobuffalo/buffalo:/render/json_test.go",
	"github.com/gobuffalo/buffalo:/render/markdown_test.go",
	"github.com/gobuffalo/buffalo:/render/options.go",
	"github.com/gobuffalo/buffalo:/render/partials_test.go",
	"github.com/gobuffalo/buffalo:/render/plain.go",
	"github.com/gobuffalo/buffalo:/render/plain_test.go",
	"github.com/gobuffalo/buffalo:/render/render.go",
	"github.com/gobuffalo/buffalo:/render/render_test.go",
	"github.com/gobuffalo/buffalo:/render/renderer.go",
	"github.com/gobuffalo/buffalo:/render/sse.go",
	"github.com/gobuffalo/buffalo:/render/string.go",
	"github.com/gobuffalo/buffalo:/render/string_map.go",
	"github.com/gobuffalo/buffalo:/render/string_map_test.go",
	"github.com/gobuffalo/buffalo:/render/string_test.go",
	"github.com/gobuffalo/buffalo:/render/template.go",
	"github.com/gobuffalo/buffalo:/render/template_engine.go",
	"github.com/gobuffalo/buffalo:/render/template_helpers.go",
	"github.com/gobuffalo/buffalo:/render/template_helpers_test.go",
	"github.com/gobuffalo/buffalo:/render/template_test.go",
	"github.com/gobuffalo/buffalo:/render/xml.go",
	"github.com/gobuffalo/buffalo:/render/xml_test.go",
	"github.com/markbates/pkger/internal/testdata/app:/",
	"github.com/markbates/pkger/internal/testdata/app:/public",
	"github.com/markbates/pkger/internal/testdata/app:/public/images",
	"github.com/markbates/pkger/internal/testdata/app:/public/images/mark-small.png",
	"github.com/markbates/pkger/internal/testdata/app:/public/images/mark.png",
	"github.com/markbates/pkger/internal/testdata/app:/public/images/mark_250px.png",
	"github.com/markbates/pkger/internal/testdata/app:/public/images/mark_400px.png",
	"github.com/markbates/pkger/internal/testdata/app:/public/index.html",
}
