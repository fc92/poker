// Inspired by
// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"encoding/json"
	"log"

	"github.com/sietchcode/poker/internal/common"
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
	room *common.Room

	// channel to update room from hub and clients
	participantReceived chan common.Participant
}

func newHub() *Hub {
	return &Hub{
		register:            make(chan *Client),
		unregister:          make(chan *Client),
		clients:             make(map[*Client]bool),
		room:                common.NewRoom(),
		participantReceived: make(chan common.Participant),
	}
}

func (h *Hub) broadcastRoom() {
	for client := range h.clients {
		// remove vote if it is not closed
		filteredRoom := h.room.FilterVoteData(client.voterId)
		jsonRoom, err := json.Marshal(filteredRoom)
		if err != nil {
			log.Printf("error: %jsonRoom", err)
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
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				h.removeVoter(client)
				close(client.send)
				h.broadcastRoom()
			}
		case participantReceived := <-h.participantReceived:
			h.room.UpdateFromParticipant(participantReceived)
			h.broadcastRoom()
		}
	}
}

func (h *Hub) removeVoter(client *Client) {
	for i, voter := range h.room.Voters {
		if voter.Id == client.voterId {
			h.room.Voters[i] = h.room.Voters[len(h.room.Voters)-1] // Copy last element to index i.
			h.room.Voters[len(h.room.Voters)-1] = nil              // Erase last element (write zero value).
			h.room.Voters = h.room.Voters[:len(h.room.Voters)-1]   // Truncate slice.
			break
		}
	}
	h.room.UpdateFromHub()
}
