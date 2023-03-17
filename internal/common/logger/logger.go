// logger shared between files of the package groom
package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitLogger() {
	// Configure logger to write to the file and include caller information
	log.Logger = log.With().Caller().Logger()
	log.Logger = log.With().CallerWithSkipFrameCount(1).Logger()
	// Create a file for logging
	file, err := os.OpenFile("./poker.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open log file")
	}
	log.Logger = log.Output(file)

	// Set global log level to debug
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}
