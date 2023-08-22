package main

import (
	"flag"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/fc92/poker/internal/common/logger"
	"github.com/fc92/poker/internal/server"
)

func init() {
	logger.InitLogger()
}

func main() {

	websocketHostPortPtr := flag.String("websocket", "localhost:8080", "hostname and port of the websocket to open")
	debug := flag.Bool("debug", false, "sets log level to debug")

	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	log.Info().Msgf("websocket: %s", *websocketHostPortPtr)
	server.StartServer(*websocketHostPortPtr)
}
