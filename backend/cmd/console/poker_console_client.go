package main

import (
	"flag"
	"time"

	"github.com/goombaio/namegenerator"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/fc92/poker/internal/common/logger"
	"github.com/fc92/poker/internal/player"
)

func init() {
	logger.InitLogger()
}

func main() {
	// Define flags
	playerNamePtr := flag.String("name", "", "player name (or empty string for generated name)")
	serverHostPortPtr := flag.String("websocket", "localhost:8080", "server hostname:port")
	debug := flag.Bool("debug", false, "sets log level to debug")

	// Parse flags
	flag.Parse()

	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	// if no name is provided, use a name generator
	if *playerNamePtr == "" {
		seed := time.Now().UTC().UnixNano()
		nameGenerator := namegenerator.NewNameGenerator(seed)
		*playerNamePtr = nameGenerator.Generate()
	}
	log.Info().Msgf("player name: %s", *playerNamePtr)
	log.Info().Msgf("  server websocket: %s", *serverHostPortPtr)
	player.Play(*playerNamePtr, *serverHostPortPtr)
}
