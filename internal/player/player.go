package player

import (
	"encoding/json"
	"log"
	"net/url"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
	"golang.org/x/exp/maps"

	co "github.com/sietchcode/poker/internal/common"
	console "github.com/sietchcode/poker/internal/player/console"
)

// send participant data to the server
func sendVoter(c *websocket.Conn, voter *co.Participant) {
	jsonVoter, err := json.Marshal(voter)
	if err != nil {
		log.Println("json error:", err)
		return
	}
	err = c.WriteMessage(websocket.TextMessage, jsonVoter)
	if err != nil {
		log.Println("write:", err)
	}
}

func updateFromServer(c *websocket.Conn, voter *co.Participant, localRoom *co.Room, displayControl chan<- bool, controlFromServer chan<- []byte) {
	for {
		_, message, err := c.ReadMessage()
		//TODO differentiate timeout and socket closed
		if err != nil {
			return
		}
		controlFromServer <- message
	}
}

func cleanExit(c *websocket.Conn) {
	// Cleanly close the connection by sending a close message and then
	// waiting (with timeout) for the server to close the connection.
	err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Println("write close:", err)
		return
	}
}

func Play(name string, serverAddress string) {
	voter := co.CreateVoter(name)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: serverAddress, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()
	// inform server of the new voter
	sendVoter(c, voter)
	controlFromUI := make(chan string)
	defer close(controlFromUI)
	displayControl := make(chan bool)
	defer close(displayControl)
	controlFromServer := make(chan []byte)
	defer close(controlFromServer)

	// start console display
	room := co.NewRoom()
	go console.Display(voter, room, controlFromUI, displayControl)

	// get room updates from server
	go updateFromServer(c, voter, room, displayControl, controlFromServer)

	for {
		// all changes to room and voter must be done here on player side
		select {
		case command := <-controlFromUI:
			switch command {
			case co.CommandQuit:
				log.Println("user quit")
				cleanExit(c)
				return
			case co.CommandStartVote:
				voter.LastCommand = command
				voter.Vote = co.VoteNotReceived
				sendVoter(c, voter)
			case co.CommandCloseVote:
				voter.LastCommand = command
				sendVoter(c, voter)
			case co.CommandNotVoting:
				voter.Vote = room.VoteCommands()[command]
				voter.LastCommand = command
				sendVoter(c, voter)
			case co.CommandVote1, co.CommandVote2, co.CommandVote3, co.CommandVote5, co.CommandVote8, co.CommandVote13:
				voter.Vote = room.VoteCommands()[command]
				voter.LastCommand = co.VoteReceived
				sendVoter(c, voter)
			}
		case <-interrupt:
			cleanExit(c)
			log.Println("interrupt")
			return
		case message := <-controlFromServer:
			// remove locally stored commands to keep only received commands
			for _, voter := range room.Voters {
				maps.Clear(voter.AvailableCommands)
			}
			// update room with data received from server
			if err = json.Unmarshal(message, room); err != nil {
				log.Printf("unknown message, not a Room: %v", err)
			}
			voter.UpdateLocalPlayerFromServer(room)
			displayControl <- true // refresh display
		}
	}

}
