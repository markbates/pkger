package pkgtest

import (
	"path/filepath"

	"github.com/markbates/pkger/here"
)

type AppDetails struct {
	here.Info
	Paths struct {
		Root   []string
		Parser []string
		Public []string
	}
}

// App returns here.info that represents the
// ./internal/testdata/app. This should be used
// by tests.
func App() (AppDetails, error) {
	var app AppDetails

	her, err := here.Package("github.com/markbates/pkger")
	if err != nil {
		return app, err
	}

	info := here.Info{
		ImportPath: "app",
	}

	x := make([]string, len(rootPaths))
	copy(x, rootPaths)
	app.Paths.Root = x

	x = make([]string, len(parserPaths))
	copy(x, parserPaths)
	app.Paths.Parser = x

	x = make([]string, len(publicPaths))
	copy(x, publicPaths)
	app.Paths.Public = x

	ch := filepath.Join(
		her.Dir,
		"pkging",
		"pkgtest",
		"internal",
		"testdata",
		"app")

	info.Dir = ch

	info, err = here.Cache(info.ImportPath, func(s string) (here.Info, error) {
		return info, nil
	})
	if err != nil {
		return app, err
	}
	app.Info = info
	return app, nil
}

var rootPaths = []string{
	"app:/",
	"app:/go.mod",
	"app:/go.sum",
	"app:/main.go",
	"app:/public",
	"app:/public/images",
	"app:/public/images/img1.png",
	"app:/public/images/img2.png",
	"app:/public/index.html",
	"app:/templates",
	"app:/templates/a.txt",
	"app:/templates/b",
	"app:/templates/b/b.txt",
}

var publicPaths = []string{
	"app:/public",
	"app:/public/images",
	"app:/public/images/img1.png",
	"app:/public/images/img2.png",
	"app:/public/index.html",
}

var parserPaths = []string{
	"app:/public",
	"app:/public/images",
	"app:/public/images/img1.png",
	"app:/public/images/img2.png",
	"app:/public/index.html",
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
}
