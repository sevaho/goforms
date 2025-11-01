package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func handleGetIndex() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		params := Params{"Language": "NL", "Country": ""}
		return ctx.Render(http.StatusOK, "index", params)
	}
}
