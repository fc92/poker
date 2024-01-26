package server

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/fc92/poker/internal/common"
)

func TestNewHub(t *testing.T) {
	hub := newHub()
	assert.NotNil(t, hub, "Expected a new hub, but got nil")
	assert.NotNil(t, hub.register, "Expected register channel, but got nil")
	assert.NotNil(t, hub.unregister, "Expected unregister channel, but got nil")
	assert.NotNil(t, hub.clients, "Expected clients map, but got nil")
	assert.NotNil(t, hub.rooms, "Expected rooms, but got nil")
	assert.NotNil(t, hub.participantReceived, "Expected participantReceived")

}

func TestBroadcastRoom(t *testing.T) {
	hub := newHub()
	client1 := &Client{roomName: "test1"}
	client2 := &Client{roomName: "test2"}
	hub.clients[client1] = true
	hub.clients[client2] = true

	// Test sending JSON room data
	hub.rooms["test1"] = common.NewRoom()
	hub.rooms["test2"] = common.NewRoom()
	hub.broadcastRoom()
	select {
	case jsonRoom := <-client1.send:
		var room common.Room
		err := json.Unmarshal(jsonRoom, &room)
		if err != nil {
			t.Error("Expected JSON room data, but got error: ", err)
		}
	default:

	}

	// Test removing client from clients map and voter from room
	client1.send = make(chan []byte, 1)
	hub.broadcastRoom()
	select {
	case jsonRoom := <-client1.send:
		t.Error("Expected client to be removed, but got JSON room data: ", jsonRoom)
	default:
		if _, ok := hub.clients[client1]; ok {
			t.Error("Expected client to be removed, but it was still in clients map")
		}
	}
}

func TestRun(t *testing.T) {
	hub := newHub()
	room := common.NewRoom()
	room.Name = "testRoom"
	hub.rooms["testRoom"] = room
	// Test updating room from participant
	participant := &common.Participant{Id: uuid.New(), RoomName: "testRoom"}
	room.Voters = append(room.Voters, participant)
	hub.rooms["testRoom"].UpdateFromParticipant(*participant)
}

func TestRemoveVoter(t *testing.T) {
	hub := newHub()
	room := common.NewRoom()
	room.Name = "testRoom"
	hub.rooms["testRoom"] = room
	voter1 := &common.Participant{Id: uuid.New(), RoomName: "testRoom"}
	voter2 := &common.Participant{Id: uuid.New(), RoomName: "testRoom"}
	hub.rooms["testRoom"].Voters = append(hub.rooms["testRoom"].Voters, voter1, voter2)
	client := &Client{voterId: voter1.Id, roomName: "testRoom"}

	// Test removing a voter
	hub.removeVoter(client)
	if len(hub.rooms["testRoom"].Voters) != 1 || hub.rooms["testRoom"].Voters[0] != voter2 {
		t.Error("Expected voter1 to be removed, but it was still in the voters slice")
	}
}
