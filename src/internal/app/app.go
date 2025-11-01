package app

import (
	"context"
	"fmt"
	"io/fs"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sevaho/goforms/src/config"
	database "github.com/sevaho/goforms/src/db"
	"github.com/sevaho/goforms/src/internal/models"
	"github.com/sevaho/goforms/src/internal/repository"
	"github.com/sevaho/goforms/src/mailproviders"
	"github.com/sevaho/goforms/src/mailproviders/fake"
	"github.com/sevaho/goforms/src/mailproviders/mailersend"
	"github.com/sevaho/goforms/src/pkg/logger"
	"github.com/sevaho/goforms/src/pkg/recaptcha"
	"github.com/sevaho/goforms/src/pkg/renderer"
	"github.com/sevaho/goforms/src/pkg/telegram"
	"github.com/sevaho/goforms/src/web"
	"github.com/unrolled/secure"

	"github.com/sevaho/livereload"
)

// TODO:  <03-05-25, Sebastiaan Van Hoecke> // Get rid of the struct pass by arguments
type App struct {
	// This is a container for testability, should be removed as we use ginkgo
	port   int
	server *echo.Echo

	telegram *telegram.TelegramService
	// Need to put it behind a mailer service that proxies other mta's as well
	mailproviders *mailproviders.MailProviders
	formsConfig   models.FormsConfig
	renderer      *renderer.RenderEngine
	config        *config.Config
	repository    *repository.Repository
	db            database.Querier
	tx            *pgx.Tx
}

func (app *App) logErrorFunc(c echo.Context, err error, stack []byte) error {
	// Passing app to be able to get telegram
	if app.config.IS_DEVELOPMENT == false {
		app.telegram.SendNotification("App panic", err.Error())
	}
	return err
}

// TODO:  <03-05-25, Sebastiaan Van Hoecke> // This should return a pointer to echo
func New(config *config.Config) *App {
	server := echo.New()
	server.Debug = config.IS_DEVELOPMENT

	if config.IS_DEVELOPMENT {
		server.StaticFS("/static", os.DirFS(config.STATIC_DIRECTORY))
	} else {
		staticFS, _ := fs.Sub(web.Static, "static")
		server.StaticFS("/static", staticFS)
	}
	server.HideBanner = true
	server.HidePort = true

	// * * * * * * * * * * * * * * * * * *
	// RENDER ENGINE
	// * * * * * * * * * * * * * * * * * *
	renderer := renderer.NewRenderEngine(config.IS_DEVELOPMENT, config.TEMPLATES_DIRECTORY, config.RELEASE, web.Templates)
	server.Renderer = renderer

	// * * * * * * * * * * * * * * * * * *
	// DEPENDENCIES
	// * * * * * * * * * * * * * * * * * *

	db, tx := database.NewDB(config.DB_DSN, config.RUN_IN_TRANSACTION)

	app := App{
		server:   server,
		renderer: renderer,
		mailproviders: mailproviders.New(
			mailersend.New(config.MAILERSEND_API_KEY),
			fake.New(),
		),
		telegram:    telegram.New(config.TELEGRAM_BOT_API_KEY, config.TELEGRAM_BOT_CHAT_ID),
		formsConfig: loadFormConfig(config),
		db:          db,
		tx:          tx,
		repository:  repository.New(db, config.SECRET_KEY),
		config:      config,
	}

	// * * * * * * * * * * * * * * * * * *
	// LIBS
	// * * * * * * * * * * * * * * * * * *
	recaptcha.Init(config.GOOGLE_RECAPTCHA_SECRET_KEY)

	// * * * * * * * * * * * * * * * * * *
	// MIDDLEWARE
	// * * * * * * * * * * * * * * * * * *
	server.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{LogErrorFunc: app.logErrorFunc}))
	server.Use(middleware.RequestID())
	secureMiddleware := secure.New(secure.Options{
		FrameDeny:     true,
		IsDevelopment: config.IS_DEVELOPMENT,
	})
	server.Use(echo.WrapMiddleware(secureMiddleware.Handler))

	if config.IS_DEVELOPMENT {
		server.Use(livereload.LiveReload(server, logger.Logger, config.TEMPLATES_DIRECTORY, config.STATIC_DIRECTORY))
	}

	// * * * * * * * * * * * * * * * * * *
	// ROUTES
	// * * * * * * * * * * * * * * * * * *
	app.addRoutes()

	return &app
}

func (app *App) Serve(port int) {
	go func() {
		logger.Logger.Info().Msgf("[HTTP SERVER] Running on http://localhost:%d", port)
		err := app.server.Start(fmt.Sprint(":", port))
		if err != nil {
			logger.Logger.Warn().Err(err).Stack().Msgf("[HTTP SERVER] closed unexpectedly, reason: %s", err)
		}
	}()
}

func (app *App) ShutDown(ctx context.Context) {

	if app.tx != nil {
		transaction := *app.tx
		err := transaction.Rollback(context.Background())
		logger.Logger.Debug().Msgf("Transaction rolled back")
		if err != nil {
			logger.Logger.Fatal().Msgf("Failed to rollback DB transaction: %s", err)
		}
	}
	err := app.server.Shutdown(ctx)
	if err != nil {
		logger.Logger.Error().Err(err).Stack().Msgf("[HTTP SERVER] Error shutting down the server: %s", err)
	}
	logger.Logger.Info().Msgf("[HTTP SERVER] shut down.")
}
