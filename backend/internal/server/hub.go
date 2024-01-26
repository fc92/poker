// Inspired by
// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"encoding/json"

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
	return &Hub{
		register:            make(chan *Client),
		unregister:          make(chan *Client),
		clients:             make(map[*Client]bool),
		rooms:               make(map[string]*common.Room),
		participantReceived: make(chan common.Participant),
	}
}

func (h *Hub) broadcastRoom() {
	for client := range h.clients {
		// remove vote if it is not closed
		filteredRoom := h.rooms[client.roomName].FilterVoteData(client.voterId)
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
			}
		case participantReceived := <-h.participantReceived:
			// use existing room
			if h.rooms[participantReceived.RoomName] != nil {
				h.rooms[participantReceived.RoomName].UpdateFromParticipant(participantReceived)
			} else {
				// create new room
				newRoom := common.NewRoom()
				newRoom.Name = participantReceived.RoomName
				newRoom.Voters = append(newRoom.Voters, &participantReceived)
			}
			h.broadcastRoom()
		}
	}
}

func (h *Hub) removeVoter(client *Client) {
	for i, voter := range h.rooms[client.roomName].Voters {
		if voter.Id == client.voterId {
			h.rooms[client.roomName].Voters[i] = h.rooms[client.roomName].Voters[len(h.rooms[client.roomName].Voters)-1] // Copy last element to index i.
			h.rooms[client.roomName].Voters[len(h.rooms[client.roomName].Voters)-1] = nil                                // Erase last element (write zero value).
			h.rooms[client.roomName].Voters = h.rooms[client.roomName].Voters[:len(h.rooms[client.roomName].Voters)-1]   // Truncate slice.
			break
		}
	}
	// remove empty room
	h.rooms[client.roomName].UpdateFromHub()
	if len(h.rooms) == 0 {
		delete(h.rooms, client.roomName)
	}
}
