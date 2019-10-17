package parser

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/markbates/pkger/here"
	"github.com/stretchr/testify/require"
)

func Test_Parser_App(t *testing.T) {
	r := require.New(t)

	pwd, err := os.Getwd()
	r.NoError(err)

	ch := filepath.Join(pwd, "..",
		"examples",
		"app")

	os.RemoveAll(filepath.Join(ch, "pkged.go"))

	info := here.Info{
		Dir:        ch,
		ImportPath: "github.com/markbates/pkger/examples/app",
	}
	here.Cache(info.ImportPath, func(s string) (here.Info, error) {
		return info, nil
	})

	res, err := Parse(info)

	r.NoError(err)

	sort.Strings(inbed)

	act := make([]string, len(res))
	for i := 0; i < len(res); i++ {
		act[i] = res[i].String()
	}

	sort.Strings(act)
	r.Equal(inbed, act)
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
	"github.com/markbates/pkger/examples/app:/",
	"github.com/markbates/pkger/examples/app:/public",
	"github.com/markbates/pkger/examples/app:/public/images",
	"github.com/markbates/pkger/examples/app:/public/images/mark-small.png",
	"github.com/markbates/pkger/examples/app:/public/images/mark.png",
	"github.com/markbates/pkger/examples/app:/public/images/mark_250px.png",
	"github.com/markbates/pkger/examples/app:/public/images/mark_400px.png",
	"github.com/markbates/pkger/examples/app:/public/index.html",
}
