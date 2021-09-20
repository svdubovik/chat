package logger

import (
	"os"

	"github.com/rs/zerolog"
)

const (
	LOGFORMATLOGFMT = "logfmt"
	LOGFORMATJSON   = "json"
)

func NewLogger(logLevel, logFormat, service string) *zerolog.Logger {
	var lvl zerolog.Level

	switch logLevel {
	case "error":
		lvl = zerolog.ErrorLevel
	case "warn":
		lvl = zerolog.WarnLevel
	case "info":
		lvl = zerolog.InfoLevel
	case "debug":
		lvl = zerolog.DebugLevel
	default:
		panic("unexpected log level")
	}

	logger := zerolog.New(os.Stderr).With().Str("service", service).Timestamp().Logger().Level(lvl)

	if logFormat == LOGFORMATLOGFMT {
		logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	}

	return &logger
}
