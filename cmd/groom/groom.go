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
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	// Default level for this example is info, unless debug flag is present
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if len(os.Args) != 2 {
		log.Fatal().Msg("Usage: groom serverUrl\nExample: groom http://localhost")
	}
	serverUrl := os.Args[1]

	// // get list of deployed rooms
	rooms := []string{"TeamBlue", "TeamRed", "TeamBlack"}
	// rooms, err := groom.RoomDeployed()
	// if err != nil {
	// 	log.Fatal().Msg("unable to get list of rooms deployed...")
	// 	os.Exit(1)
	// } else {
	// 	log.Info().Msgf("Found rooms: %v", rooms)
	// }

	// // TO DO replace this code
	// rooms, err = groom.AddRoom("TeamGreen")
	// if err != nil {
	// 	log.Warn().Msgf("%v", err)
	// } else {
	// 	log.Info().Msgf("Found rooms: %v", rooms)
	// }

	// rooms, err = groom.RemoveRoom("TeamPink")
	// if err != nil {
	// 	log.Warn().Msgf("%v", err)
	// } else {
	// 	log.Info().Msgf("Found rooms: %v", rooms)
	// }
	// display welcome screen
	groom.DisplayWelcome(rooms, serverUrl)

	// handle updates from user input or helm release

}
