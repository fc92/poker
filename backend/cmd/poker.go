package main

import (
	"flag"
	"os"
	"time"

	"github.com/goombaio/namegenerator"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/fc92/poker/internal/common/logger"
	"github.com/fc92/poker/internal/player"
	"github.com/fc92/poker/internal/server"
)

func init() {
	logger.InitLogger()
}

func main() {
	clientCmd := flag.NewFlagSet("client", flag.ExitOnError)
	clientName := clientCmd.String("name", "", "name of the player")
	serverWS := clientCmd.String("websocket", "localhost:8080", "hostname and port of the server websocket")

	serverCmd := flag.NewFlagSet("server", flag.ExitOnError)
	websocket := serverCmd.String("websocket", "localhost:8080", "hostname and port of the websocket to open")
	debug := serverCmd.Bool("debug", false, "sets log level to debug")

	if len(os.Args) < 2 {
		log.Fatal().Msg("expected 'client' or 'server' subcommands")
	}

	switch os.Args[1] {

	case "client":
		clientCmd.Parse(os.Args[2:])
		// if no name is provided, use a name generator
		if *clientName == "" {
			seed := time.Now().UTC().UnixNano()
			nameGenerator := namegenerator.NewNameGenerator(seed)
			*clientName = nameGenerator.Generate()
		}
		log.Info().Msg("subcommand 'client'")
		log.Info().Msgf("  name: %s", *clientName)
		log.Info().Msgf("  server websocket: %s", *serverWS)
		player.Play(*clientName, *serverWS)
	case "server":
		serverCmd.Parse(os.Args[2:])
		if *debug {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		}
		log.Info().Msg("subcommand 'server'")
		log.Info().Msgf("  websocket: %s", *websocket)
		server.StartServer(*websocket)
	default:
		log.Fatal().Msg("expected 'client' or 'server' subcommands")
	}
}
