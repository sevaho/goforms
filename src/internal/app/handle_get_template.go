package app

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	static "github.com/sevaho/goforms/src/web"
)

func handleGetTemplatePage() echo.HandlerFunc {

	return func(ctx echo.Context) error {

		template := ctx.Param("template")

		if template != "" {
			params := Params{
				"Language": "NL",
				"Country":  "",
				"Content":  ctx.QueryParams(),
				"Website":  "https://goforms.dev",
			}

			return ctx.Render(http.StatusOK, template, params)
		} else {
			templates, err := getAllTemplateFilePaths(static.Templates)

			if err != nil {
				return ctx.Render(500, "error", Params{"Error": err.Error()})
			}
			params := Params{
				"Language":  "NL",
				"Country":   "",
				"Content":   ctx.QueryParams(),
				"Website":   "https://goforms.dev",
				"Templates": templates,
			}

			return ctx.Render(http.StatusOK, "all_templates", params)
		}

	}
}
func getAllTemplateFilePaths(fsys embed.FS) ([]string, error) {
	var paths []string
	err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			path = strings.TrimPrefix(path, "templates/")
			path = strings.TrimSuffix(path, ".html")
			paths = append(paths, path)
		}
		return nil
	})
	return paths, err
}
