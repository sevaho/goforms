package app

func (app *App) addRoutes() {
	// * * * * * * * * * * * * * * * * * *
	// SETUP ROUTES
	// * * * * * * * * * * * * * * * * * *
	app.server.GET("/", handleGetIndex())
	app.server.POST("/forms/:id", handlePostForm(
		app.mailproviders,
		app.telegram,
		app.formsConfig,
		app.renderer,
		app.repository,
	))

	// API
	apiGroup := app.server.Group("/api")
	apiGroup.Use(CheckAuthorizationBearerTokenMiddleware(app.config))
	apiGroup.GET("/config", handleGetConfig(app.formsConfig))
	apiGroup.GET("/mails", handleGetMails(app.repository))
	apiGroup.GET("/mails/:id", handleGetMailByID(app.repository))

	// admin
	app.server.GET("/admin", handleGetAdminDashboard(app.repository), CheckApiTokenQueryParamsMiddleware(app.config))

	// healthz
	app.server.GET("/healthz", handleGetIndex())

	// development
	if app.config.IS_DEVELOPMENT {
		app.server.GET("/templates/:template", handleGetTemplatePage())
		app.server.GET("/templates", handleGetTemplatePage())
	}
}
