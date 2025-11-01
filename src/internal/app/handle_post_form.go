package app

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/k3a/html2text"
	"github.com/labstack/echo/v4"
	"github.com/sevaho/goforms/src/internal/models"
	"github.com/sevaho/goforms/src/internal/repository"
	"github.com/sevaho/goforms/src/mailproviders"
	"github.com/sevaho/goforms/src/pkg/logger"
	"github.com/sevaho/goforms/src/pkg/recaptcha"
	"github.com/sevaho/goforms/src/pkg/renderer"
	"github.com/sevaho/goforms/src/pkg/telegram"
)

func verifyCaptcha(realIP string, recaptchaResponse string) error {
	success, err := recaptcha.Confirm(realIP, recaptchaResponse)
	if !success {
		logger.Logger.Error().Err(err).Msg("Invalid captcha or no captcha provided.")

		if err != nil {
			return err
		} else {
			return errors.New("Invalid captcha or no captcha provided.")
		}
	}
	return nil
}

func handlePostForm(
	mailproviders *mailproviders.MailProviders,
	telegram *telegram.TelegramService,
	forms models.FormsConfig,
	renderer *renderer.RenderEngine,
	repository *repository.Repository,
) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		language := ctx.QueryParam("language")
		country := ctx.QueryParam("country")
		FormID, err := uuid.Parse(ctx.Param("id"))

		if language == "" {
			language = "NL"
		}

		if country == "" {
			country = "BE"
		}

		if err != nil {
			return ctx.Render(500, "error", Params{"Error": err.Error()})
		}

		form, err := forms.Get(FormID)
		if err != nil {
			return ctx.Render(500, "error", Params{"Error": err.Error()})
		}

		formData, err := ctx.FormParams()
		if err != nil {
			return ctx.Render(500, "error", Params{"Error": err.Error()})
		}

		// Parse Captcha
		if !form.Skipcaptcha {
			if err := verifyCaptcha(ctx.RealIP(), formData.Get("g-recaptcha-response")); err != nil {
				return ctx.Render(500, "error", Params{"Error": err.Error()})
			}
		}
		formData.Del("g-recaptcha-response")

		// Set the subject of the email
		var subject string

		if formData.Get("subject") != "" {
			subject = formData.Get("subject")
		} else if form.Subject != "" {
			subject = form.Subject
		} else {
			err = errors.New("No subject found!")
			telegram.SendNotification("Error with subject", err.Error())
			logger.Logger.Error().Err(err).Msg("No subject")
			return err
		}

		// Parse website
		website := ctx.Request().Header.Get("Origin")
		if website == "" {
			website = ctx.Request().Header.Get("Referer")
		}

		html := renderer.MustRenderHTML("email/contact", Params{
			"Website":  website,
			"Content":  formData,
			"Language": language,
			"Country":  country,
		}, "layout")

		plain := html2text.HTML2Text(string(html))

		err = mailproviders.Get(form.Provider).Mail(string(html), plain, subject, form.Sender, form.Recipients)

		repository.Store(form.Provider, subject, form.Sender.Email, string(html), form.Recipients, err)

		if err != nil {
			telegram.SendNotification("Error with mailersend", err.Error())
			logger.Logger.Error().Err(err).Msg("Something went wrong while sending email with mailersend.")
			return err
		}

		// Return with success page
		return ctx.Render(http.StatusOK, "success", Params{"Language": language, "Country": country})
	}
}
