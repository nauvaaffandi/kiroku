package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Init(devMode bool) {
	zerolog.TimeFieldFormat = time.RFC3339

	if devMode {
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.RFC3339,
			NoColor:    false,
		}).Level(zerolog.DebugLevel)
	} else {
		log.Logger = zerolog.New(os.Stdout).Level(zerolog.InfoLevel).With().Timestamp().Logger()
	}
}

func GetLogger() *zerolog.Logger {
	l := log.Logger
	return &l
}
