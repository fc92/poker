// logger shared
package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitLogger() {
	// Create a new zerolog logger with the lumberjack logger as the output
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger().With().Caller().Logger()

	// Set global log level to debug
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}
