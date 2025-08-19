package logger

import (
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

func New(appEnv string) zerolog.Logger {
	lvl := zerolog.InfoLevel
	if strings.ToLower(appEnv) == "development" {
		lvl = zerolog.DebugLevel
	}
	zerolog.TimeFieldFormat = time.RFC3339Nano
	logger := zerolog.New(os.Stdout).
		With().
		Timestamp().
		Str("service", "jwt-otp-auth").
		Logger().
		Level(lvl)
	return logger
}
