/*
groom allows to add or remove poker rooms.

It must be used as part of the helm deployment of poker
and dynamically updates the helm release to add or remove independent poker rooms
in dedicated pods and services
*/
package main

import (
	"os"
	//
	// Uncomment to load all auth plugins
	//
	// Or uncomment to load specific auth plugins

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/fc92/poker/internal/groom"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.With().Caller().Logger()
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	// Default level for this example is info, unless debug flag is present
	zerolog.SetGlobalLevel(zerolog.ErrorLevel)

	if len(os.Args) != 2 {
		log.Fatal().Msg("Usage: groom serverUrl\nExample: groom localhost")
	}
	serverUrl := os.Args[1]

	groom.DisplayWelcome(serverUrl)
}
