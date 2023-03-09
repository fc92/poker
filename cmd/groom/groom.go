/*
groom allows to add or remove poker rooms.

It must be used as part of the helm deployment of poker
and dynamically updates the helm release to add or remove independent poker rooms
in dedicated pods and services
*/
package main

import (
	"github.com/fc92/poker/internal/groom"
	"log"
	"os"
)

func main() {

	if len(os.Args) != 2 {
		log.Fatal("Usage: groom serverUrl\nExample: groom localhost")
	}
	serverUrl := os.Args[1]

	groom.DisplayWelcome(serverUrl)
}
