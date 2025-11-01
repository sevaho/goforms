package app

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sevaho/goforms/src/internal/models"
	"github.com/sevaho/goforms/src/internal/repository"
	"github.com/sevaho/goforms/src/pkg/logger"
)

func handleGetAdminDashboard(
	repository *repository.Repository,
) echo.HandlerFunc {

	type ResponseModel struct {
		Items []models.DecryptedMail `json:"items"`
		Count int                    `json:"count"`
	}

	return func(c echo.Context) error {
		page, pageLen := 1, 10

		if pageParam := c.QueryParam("page"); pageParam != "" {
			if v, err := strconv.Atoi(pageParam); err == nil && v > 0 {
				page = v
			}
		}

		if pageLenParam := c.QueryParam("pagelen"); pageLenParam != "" {
			if v, err := strconv.Atoi(pageLenParam); err == nil && v > 0 && v <= 100 {
				pageLen = v
			}
		}

		offset := (page - 1) * pageLen

		items, count, err := repository.GetMails(offset, pageLen)

		if err != nil {
			logger.Logger.Error().Err(err).Msg("Something went wrong while querying database.")
			return c.Render(500, "error", Params{"Error": err.Error()})
		}

		params := Params{"Mails": items, "Count": count}
		return c.Render(http.StatusOK, "admin", params)
	}
}
