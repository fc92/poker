package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/goombaio/namegenerator"

	"github.com/fc92/poker/internal/player"
	"github.com/fc92/poker/internal/server"
)

func main() {
	clientCmd := flag.NewFlagSet("client", flag.ExitOnError)
	clientName := clientCmd.String("name", "", "name of the player")
	serverWS := clientCmd.String("websocket", "localhost:8080", "hostname and port of the server websocket")

	serverCmd := flag.NewFlagSet("bar", flag.ExitOnError)
	websocket := serverCmd.String("websocket", "localhost:8080", "hostname and port of the websocket to open")

	if len(os.Args) < 2 {
		fmt.Println("expected 'client' or 'server' subcommands")
		os.Exit(1)
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
		fmt.Println("subcommand 'client'")
		fmt.Println("  name:", *clientName)
		fmt.Println("  server websocket:", *serverWS)
		player.Play(*clientName, *serverWS)
	case "server":
		serverCmd.Parse(os.Args[2:])
		fmt.Println("subcommand 'server'")
		fmt.Println("  websocket:", *websocket)
		server.StartServer(*websocket)
	default:
		fmt.Println("expected 'client' or 'server' subcommands")
		os.Exit(1)
	}
}
