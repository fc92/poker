// logger shared between files of the package groom
package groom

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func initLogger() {
	// Create a file for logging
	file, err := os.OpenFile("./log/groom.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open log file")
	}

	// Configure logger to write to the file and include caller information
	log.Logger = log.Output(file)
	log.Logger = log.With().Caller().Logger()

	// Set global log level to debug
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}
