// Game room
// Game logic independant from underlying transport

package common

import (
	"testing"
)

func TestNewRoom(t *testing.T) {
	room := NewRoom()
	if room.RoomStatus != VoteClosed {
		t.Errorf("RoomStatus = %d; want %d", room.RoomStatus, VoteClosed)
	}
	if room.TurnFinishedCommands() == nil {
		t.Errorf("TurnFinishedCommands should not be nil")
	}
	if room.VoteCommands() == nil {
		t.Errorf("VoteCommands should not be nil")
	}
	if room.TurnStartedCommands() == nil {
		t.Errorf("TurnStartedCommands should not be nil")
	}
}
