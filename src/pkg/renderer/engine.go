package renderer

import (
	"embed"
	"html/template"
	"io"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sevaho/goforms/src/locale"
	"github.com/unrolled/render"
	"github.com/valyala/bytebufferpool"
)

type RenderEngine struct { // We need to wrap the renderer because we need a different signature for echo.
	rnd *render.Render
}

func (r *RenderEngine) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return r.rnd.HTML(w, 0, name, data) // The zero status code is overwritten by echo.
}

func (r *RenderEngine) MustRenderHTML(name string, data map[string]any, layout ...string) []byte {
	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)

	if err := r.rnd.HTML(buf, 0, name, data); err != nil {
		panic(err)
	}

	return buf.Bytes()
}

func NewRenderEngine(
	isDevelopment bool,
	directory string,
	release string,
	templatesFS embed.FS,
) *RenderEngine {

	options := render.Options{
		IsDevelopment: isDevelopment,
		Layout:        "layout",
		Extensions:    []string{".html", ".tmpl"},
		Funcs: []template.FuncMap{
			{
				"_is_development":             func() bool { return isDevelopment },
				"_release":                    func() string { return release },
				"_prettytime":                 func(t time.Time) string { return t.Format(time.RFC1123) },
				"_prettytimewithtimezone":     func(t time.Time, tz *time.Location) string { return t.In(tz).Format(time.RFC1123) },
				"_prettytimeonlywithtimezone": func(t time.Time, tz *time.Location) string { return t.In(tz).Format(time.TimeOnly) },
				"_truncate":                   func(s string, n int) string { return s[:n] + "..." },
				"_random_id":                  func() string { return uuid.New().String() },
				"_t":                          locale.T,
				"_lower":                      func(s string) string { return strings.ToLower(s) },
				"_html":                       func(s string) template.HTML { return template.HTML(s) },
				"add":                         func(a, b int) int { return a + b },
			},
		},
	}

	if isDevelopment {
		options.Directory = directory
	} else {
		options.FileSystem = &render.EmbedFileSystem{
			FS: templatesFS,
		}
	}

	return &RenderEngine{
		render.New(options),
	}
}
