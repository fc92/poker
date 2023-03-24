/*
groom allows to add or remove poker rooms.

It must be used as part of the helm deployment of poker
and dynamically updates the helm release to add or remove independent poker rooms
in dedicated pods and services
*/
package main

import (
	"log"
	"os"

	"github.com/fc92/poker/internal/groom"
)

var (
	osArgs         = func() []string { return os.Args }
	displayWelcome = groom.DisplayWelcome
)

func main() {
	if len(osArgs()) != 2 {
		log.Fatal("Usage: groom serverUrl\nExample: groom localhost")
	}
	serverUrl := osArgs()[1]

	displayWelcome(serverUrl)
}
