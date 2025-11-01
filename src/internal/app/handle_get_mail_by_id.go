package app

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sevaho/goforms/src/pkg/logger"
	"github.com/sevaho/goforms/src/internal/repository"
)

func handleGetMailByID(
	repository *repository.Repository,
) echo.HandlerFunc {
	return func(c echo.Context) error {
		MailID, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return c.JSON(500, Params{"Error": err.Error()})
		}

		result, err := repository.GetMailByID(int(MailID))

		if err != nil {
			logger.Logger.Error().Err(err).Msg("Something went wrong while querying database.")
			return c.JSON(500, Params{"Error": err.Error()})
		}

		return c.JSON(http.StatusOK, result)

	}
}
