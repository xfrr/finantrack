package xlog

import (
	"os"

	"github.com/rs/zerolog"
)

func NewZerologger(serviceName, environment string) zerolog.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).
		With().
		Str("service", serviceName).
		Str("environment", environment).
		Timestamp().
		Logger()
	return logger
}
