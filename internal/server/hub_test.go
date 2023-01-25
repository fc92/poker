package server

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"

	"github.com/fc92/poker/internal/common"
)

func TestNewHub(t *testing.T) {
	hub := newHub()
	if hub == nil {
		t.Error("Expected a new hub, but got nil")
	}
	if hub.register == nil {
		t.Error("Expected register channel, but got nil")
	}
	if hub.unregister == nil {
		t.Error("Expected unregister channel, but got nil")
	}
	if hub.clients == nil {
		t.Error("Expected clients map, but got nil")
	}
	if hub.room == nil {
		t.Error("Expected room, but got nil")
	}
	if hub.participantReceived == nil {
		t.Error("Expected participantReceived channel, but got nil")
	}
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
