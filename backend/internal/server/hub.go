// Inspired by
// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"

	"github.com/fc92/poker/internal/common"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// Room containing the business logic associated with this technical hub
	// room  *common.Room
	rooms map[string]*common.Room

	// channel to update room from hub and clients
	participantReceived chan common.Participant
}

var connectedClients prometheus.Gauge

func init() {
	connectedClients = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "poker",
		Subsystem: "server",
		Name:      "number_of_connected_clients",
		Help:      "Number of poker clients connected.",
	})
	prometheus.MustRegister(connectedClients)
	connectedClients.Set(0)
}

func newHub() *Hub {
	var initRooms map[string]*common.Room
	// Check if the environment variable ROOM_LIST is set
	roomListEnv := os.Getenv("ROOM_LIST")

	if roomListEnv != "" {
		// Load room list from environment variable
		initRooms = loadRoomsFromEnv(roomListEnv)
	} else {
		// Initialize default rooms
		initRooms = makeDefaultRooms()
	}

	return &Hub{
		register:            make(chan *Client),
		unregister:          make(chan *Client),
		clients:             make(map[*Client]bool),
		rooms:               initRooms,
		participantReceived: make(chan common.Participant),
	}
}

func makeDefaultRooms() map[string]*common.Room {
	initRooms := make(map[string]*common.Room)
	for i := 0; i < 10; i++ {
		newRoom := common.NewRoom()
		newRoom.Name = "team " + fmt.Sprint(i)
		initRooms[newRoom.Name] = newRoom
	}
	return initRooms
}

func loadRoomsFromEnv(roomListEnv string) map[string]*common.Room {
	initRooms := make(map[string]*common.Room)
	// Assuming ROOM_LIST is a comma-separated list of room names
	roomNames := strings.Split(roomListEnv, ",")

	for _, roomName := range roomNames {
		// You may need to add additional checks or validation here
		// based on the structure of your environment variable data
		newRoom := common.NewRoom()
		newRoom.Name = strings.TrimSpace(roomName)
		initRooms[newRoom.Name] = newRoom
	}

	return initRooms
}

func (h *Hub) broadcastRoom() {
	for client := range h.clients {
		room, ok := h.rooms[client.roomName]
		if ok {
			// remove vote if it is not closed
			filteredRoom := room.FilterVoteData(client.voterId)
			jsonRoom, err := json.Marshal(filteredRoom)
			if err != nil {
				log.Err(err).Msg("")
				return
			}

			select {
			case client.send <- jsonRoom:
			default:
				delete(h.clients, client)
				h.removeVoter(client)
				h.broadcastRoom()
			}
		}
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			connectedClients.Inc()
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				h.removeVoter(client)
				close(client.send)
				connectedClients.Dec()
				h.broadcastRoom()
				// broadcast roomList change
				roomReq := common.RoomRequest{}
				for c := range h.clients {
					handleRoomRequest(roomReq, c)
				}
			}
		case participantReceived := <-h.participantReceived:
			isNewPlayer := h.rooms[participantReceived.RoomName].UpdateFromParticipant(participantReceived)
			if isNewPlayer { // broadcast roomList change
				roomReq := common.RoomRequest{}
				for c := range h.clients {
					handleRoomRequest(roomReq, c)
				}
			}
			h.broadcastRoom()
		}
	}
}

func (h *Hub) removeVoter(client *Client) {
	rooms, ok := h.rooms[client.roomName]
	if ok {
		for i, voter := range rooms.Voters {
			if voter.Id == client.voterId {
				h.rooms[client.roomName].Voters[i] = h.rooms[client.roomName].Voters[len(h.rooms[client.roomName].Voters)-1] // Copy last element to index i.
				h.rooms[client.roomName].Voters[len(h.rooms[client.roomName].Voters)-1] = nil                                // Erase last element (write zero value).
				h.rooms[client.roomName].Voters = h.rooms[client.roomName].Voters[:len(h.rooms[client.roomName].Voters)-1]   // Truncate slice.
				break
			}
		}
		// close vote if needed
		h.rooms[client.roomName].UpdateFromHub()
	}
}
