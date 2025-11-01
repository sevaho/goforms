package logger

import (
	"io"
	"os"

	"github.com/rs/zerolog"
)

var Logger zerolog.Logger

func Init(colors bool, logLevel zerolog.Level) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.LevelFieldName = "levelname"
	zerolog.TimestampFieldName = "asctime"

	zerolog.SetGlobalLevel(logLevel)

	if colors {
		Logger = zerolog.New(io.MultiWriter(zerolog.ConsoleWriter{Out: os.Stderr})).With().Timestamp().Logger()
	} else {
		Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	}

}
