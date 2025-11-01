package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/amacneil/dbmate/pkg/dbmate"
	_ "github.com/amacneil/dbmate/pkg/driver/postgres"
	app "github.com/sevaho/goforms/src"
	"github.com/sevaho/goforms/src/config"
	"github.com/sevaho/goforms/src/pkg/logger"
	"github.com/spf13/pflag"
)

func main() {
	// parse arguments (flags)
	var (
		serve          = pflag.Bool("serve", false, "Serve the application.")
		port           = pflag.IntP("port", "p", 3000, "Which port to run on.")
		configfilepath = pflag.StringP("config", "c", "config.yaml", "Path to config file.")
		migrate        = pflag.Bool("migrate", false, "Run migrations.")
		newmigration   = pflag.Bool("new-migration", false, "New migrations.")
	)
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	fmt.Println(exPath)

	pflag.Parse()

	if *serve {
		config := config.New()

		if configfilepath != nil {
			config.FORMS_CONFIG_FILE_PATH = *configfilepath
		}

		ctx := context.Background()
		app.Run(ctx, *port, config)
	} else if *migrate {
		Migrate()
	} else if *newmigration {
		NewMigration()
	} else {
		logger.Logger.Warn().Msg("No flags given, exiting!")
		pflag.PrintDefaults()
	}
}

func NewMigration() {
	config := config.New()

	u, _ := url.Parse(config.DB_DSN)
	db := dbmate.New(u)
	db.NewMigration("RENAME_ME")
}

func Migrate() {
	logger.Init(true, 1)
	config := config.New()

	u, _ := url.Parse(config.DB_DSN)
	db := dbmate.New(u)

	logger.Logger.Info().Msg("Applying migrations...")
	err := db.CreateAndMigrate()
	if err != nil {
		panic(err)
	}
	logger.Logger.Info().Msg("Migrations applied!")
}
