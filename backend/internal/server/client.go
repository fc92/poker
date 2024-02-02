// Inspired by
// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// here code is only related to message transport independently of game logic

package server

import (
	"encoding/json"
	"net/http"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"

	"github.com/fc92/poker/internal/common"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 6 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 6) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	// Room voter id used to cleanly remove client connexions
	voterId  uuid.UUID
	roomName string
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Err(err).Err(err).Msg("could not read message")
			}
			return
		}
		log.Debug().Msgf("receive message %s", string(message))
		// try unmarshalling as Participant
		var voterReceived common.Participant
		if err := json.Unmarshal(message, &voterReceived); err == nil && voterReceived.RoomName != "" {
			c.voterId = voterReceived.Id
			c.roomName = voterReceived.RoomName

			// send message to hub for synced updates of the room
			c.hub.participantReceived <- voterReceived
			continue
		}
		// try unmarshalling as RoomRequest
		var roomReq common.RoomRequest
		if err := json.Unmarshal(message, &roomReq); err == nil {
			handleRoomRequest(roomReq, c)
			continue
		}
		log.Err(err).Err(err).Msg("unknown message, not a Participant or a RoomRequest")
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case responseBytes, ok := <-c.send:

			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				log.Error().Msg("channel closed from hub")
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Error().Err(err).Msg("error sending response")
				return
			}
			w.Write(responseBytes)
			if err := w.Close(); err != nil {
				log.Error().Err(err).Msg("error sending response")
				return
			}
			log.Debug().Msgf("sent message %s", string(responseBytes))

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Err(err).Err(err).Msg("could not upgrade to websocket protocol")
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}

// retrieve rooms
func handleRoomRequest(roomReq common.RoomRequest, c *Client) {

	var roomList []common.RoomOverview
	for roomName := range c.hub.rooms {
		roomOverview := common.RoomOverview{
			Name:     roomName,
			NbVoters: len(c.hub.rooms[roomName].Voters),
		}
		roomList = append(roomList, roomOverview)
	}
	roomReq.RoomList = roomList
	sort.Slice(roomReq.RoomList, func(i, j int) bool {
		return roomReq.RoomList[i].Name < roomReq.RoomList[j].Name
	})
	// convert room list to json
	response, err := json.Marshal(roomReq)
	if err != nil {
		log.Err(err).Msg("could not convert room list to json")
		return
	}
	c.send <- response

}
