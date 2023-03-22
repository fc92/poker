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
	assert.NotNil(t, hub.room, "Expected room, but got nil")
	assert.NotNil(t, hub.participantReceived, "Expected participantReceived")

}

func TestBroadcastRoom(t *testing.T) {
	hub := newHub()
	client1 := &Client{}
	client2 := &Client{}
	hub.clients[client1] = true
	hub.clients[client2] = true

	// Test sending JSON room data
	hub.room = common.NewRoom()
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

	// Test updating room from participant
	participant := &common.Participant{Id: uuid.New()}
	hub.room.UpdateFromParticipant(*participant)
}

func TestRemoveVoter(t *testing.T) {
	hub := newHub()
	voter1 := &common.Participant{Id: uuid.New()}
	voter2 := &common.Participant{Id: uuid.New()}
	hub.room.Voters = append(hub.room.Voters, voter1, voter2)
	client := &Client{voterId: voter1.Id}

	// Test removing a voter
	hub.removeVoter(client)
	if len(hub.room.Voters) != 1 || hub.room.Voters[0] != voter2 {
		t.Error("Expected voter1 to be removed, but it was still in the voters slice")
	}
}
