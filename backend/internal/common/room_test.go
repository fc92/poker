// Game room
// Game logic independent from underlying transport

package common

import (
	"testing"

	"github.com/google/uuid"
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

func TestRoom_addPlayer(t *testing.T) {
	type fields struct {
		RoomStatus           RoomVoteStatus
		Voters               []*Participant
		turnFinishedCommands map[string]string
		turnStartedCommands  map[string]string
		voteCommands         map[string]string
	}
	type args struct {
		voterReceived Participant
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"basic", fields{VoteOpen, []*Participant{}, nil, nil, nil}, args{*CreateVoter("test")}},
		{"basic",
			fields{VoteClosed, []*Participant{}, nil, nil, nil},
			args{*CreateVoter("test")}},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			room := &Room{
				RoomStatus:           tt.fields.RoomStatus,
				Voters:               tt.fields.Voters,
				turnFinishedCommands: tt.fields.turnFinishedCommands,
				turnStartedCommands:  tt.fields.turnStartedCommands,
				voteCommands:         tt.fields.voteCommands,
			}
			room.addPlayer(tt.args.voterReceived)
		})
	}
}

func TestRoom_updateFromVotes(t *testing.T) {
	type fields struct {
		RoomStatus           RoomVoteStatus
		Voters               []*Participant
		turnFinishedCommands map[string]string
		turnStartedCommands  map[string]string
		voteCommands         map[string]string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{"test1", fields{VoteClosed, []*Participant{}, nil, nil, nil}},
		{"test1", fields{VoteOpen, []*Participant{}, nil, nil, nil}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			room := &Room{
				RoomStatus:           tt.fields.RoomStatus,
				Voters:               tt.fields.Voters,
				turnFinishedCommands: tt.fields.turnFinishedCommands,
				turnStartedCommands:  tt.fields.turnStartedCommands,
				voteCommands:         tt.fields.voteCommands,
			}
			room.updateFromVotes()
		})
	}
}

func TestRoom_UpdateFromHub(t *testing.T) {
	type fields struct {
		RoomStatus           RoomVoteStatus
		Voters               []*Participant
		turnFinishedCommands map[string]string
		turnStartedCommands  map[string]string
		voteCommands         map[string]string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{"test1", fields{VoteClosed, []*Participant{}, nil, nil, nil}},
		{"test1", fields{VoteOpen, []*Participant{}, nil, nil, nil}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			room := &Room{
				RoomStatus:           tt.fields.RoomStatus,
				Voters:               tt.fields.Voters,
				turnFinishedCommands: tt.fields.turnFinishedCommands,
				turnStartedCommands:  tt.fields.turnStartedCommands,
				voteCommands:         tt.fields.voteCommands,
			}
			room.UpdateFromHub()
		})
	}
}

func TestRoom_NbVotesReceived(t *testing.T) {
	type fields struct {
		RoomStatus           RoomVoteStatus
		Voters               []*Participant
		turnFinishedCommands map[string]string
		turnStartedCommands  map[string]string
		voteCommands         map[string]string
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{"test1", fields{VoteClosed, []*Participant{}, nil, nil, nil}, 0},
		{"test1", fields{VoteOpen, []*Participant{}, nil, nil, nil}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			room := Room{
				RoomStatus:           tt.fields.RoomStatus,
				Voters:               tt.fields.Voters,
				turnFinishedCommands: tt.fields.turnFinishedCommands,
				turnStartedCommands:  tt.fields.turnStartedCommands,
				voteCommands:         tt.fields.voteCommands,
			}
			if got := room.NbVotesReceived(); got != tt.want {
				t.Errorf("Room.NbVotesReceived() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_updateRoomFromReceivedPlayer(t *testing.T) {
	type args struct {
		room          *Room
		i             int
		voterReceived Participant
	}
	room := NewRoom()
	p1 := CreateVoter("test")
	p2 := CreateVoter("test")
	p3 := CreateVoter("test")
	p1.LastCommand = CommandStartVote
	p2.LastCommand = CommandCloseVote
	p3.LastCommand = CommandCloseVote
	p3.Vote = VoteReceived
	room.addPlayer(*p1)

	tests := []struct {
		name string
		args args
	}{
		{"test1", args{room, 0, *p1}},
		{"test2", args{room, 0, *p2}},
		{"test3", args{room, 0, *p3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updateRoomFromReceivedPlayer(tt.args.room, tt.args.i, tt.args.voterReceived)
		})
	}
}

func TestRoom_UpdateFromParticipant(t *testing.T) {
	type fields struct {
		RoomStatus           RoomVoteStatus
		Voters               []*Participant
		turnFinishedCommands map[string]string
		turnStartedCommands  map[string]string
		voteCommands         map[string]string
	}
	type args struct {
		voterReceived Participant
	}
	voter := CreateVoter("player")
	voter.LastCommand = VoteNotReceived
	voter2 := CreateVoter("player2")
	voter2.LastCommand = CommandNotVoting
	voter3 := CreateVoter("player3")
	voter3.LastCommand = VoteNotReceived
	voter4 := CreateVoter("player4")
	voter4.LastCommand = VoteReceived
	voters := []*Participant{voter, voter2, voter3, voter4}

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"test1", fields{VoteClosed, []*Participant{}, nil, nil, nil}, args{*CreateVoter("test")}},
		{"test1", fields{VoteOpen, voters, nil, nil, nil}, args{*CreateVoter("test")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			room := &Room{
				RoomStatus:           tt.fields.RoomStatus,
				Voters:               tt.fields.Voters,
				turnFinishedCommands: tt.fields.turnFinishedCommands,
				turnStartedCommands:  tt.fields.turnStartedCommands,
				voteCommands:         tt.fields.voteCommands,
			}
			room.UpdateFromParticipant(tt.args.voterReceived)
		})
	}
}

func TestRoom_FilterVoteData(t *testing.T) {
	type fields struct {
		RoomStatus           RoomVoteStatus
		Voters               []*Participant
		turnFinishedCommands map[string]string
		turnStartedCommands  map[string]string
		voteCommands         map[string]string
	}
	type args struct {
		voterId uuid.UUID
	}
	room := NewRoom()
	room2 := NewRoom()
	room.OpenVote()
	room.addPlayer(*CreateVoter("player"))
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"test", fields{room.RoomStatus, room.Voters, nil, nil, nil}, args{uuid.New()}},
		{"test", fields{room2.RoomStatus, room.Voters, nil, nil, nil}, args{uuid.New()}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			room := Room{
				RoomStatus:           tt.fields.RoomStatus,
				Voters:               tt.fields.Voters,
				turnFinishedCommands: tt.fields.turnFinishedCommands,
				turnStartedCommands:  tt.fields.turnStartedCommands,
				voteCommands:         tt.fields.voteCommands,
			}
			room.FilterVoteData(tt.args.voterId)
		})
	}
}

func TestRoom_DisplayVotersStatus(t *testing.T) {
	type fields struct {
		RoomStatus           RoomVoteStatus
		Voters               []*Participant
		turnFinishedCommands map[string]string
		turnStartedCommands  map[string]string
		voteCommands         map[string]string
	}
	voter := CreateVoter("player")
	voter.LastCommand = VoteNotReceived
	voter2 := CreateVoter("player2")
	voter2.LastCommand = CommandNotVoting
	voter3 := CreateVoter("player3")
	voter3.LastCommand = VoteNotReceived
	voter4 := CreateVoter("player4")
	voter4.LastCommand = VoteReceived
	voters := []*Participant{voter, voter2, voter3, voter4}

	tests := []struct {
		name   string
		fields fields
	}{
		{"test", fields{VoteClosed, voters, nil, nil, nil}},
		{"test", fields{VoteOpen, voters, nil, nil, nil}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			room := Room{
				RoomStatus:           tt.fields.RoomStatus,
				Voters:               tt.fields.Voters,
				turnFinishedCommands: tt.fields.turnFinishedCommands,
				turnStartedCommands:  tt.fields.turnStartedCommands,
				voteCommands:         tt.fields.voteCommands,
			}
			room.DisplayVotersStatus()
		})
	}
}

func TestRoom_OpenVote(t *testing.T) {
	type fields struct {
		RoomStatus           RoomVoteStatus
		Voters               []*Participant
		turnFinishedCommands map[string]string
		turnStartedCommands  map[string]string
		voteCommands         map[string]string
	}
	voter := CreateVoter("player")
	voter.LastCommand = VoteNotReceived
	voter2 := CreateVoter("player2")
	voter2.LastCommand = CommandNotVoting
	voter3 := CreateVoter("player3")
	voter3.LastCommand = VoteNotReceived
	voter4 := CreateVoter("player4")
	voter4.LastCommand = VoteReceived
	voters := []*Participant{voter, voter2, voter3, voter4}
	tests := []struct {
		name   string
		fields fields
	}{
		{"test", fields{VoteClosed, voters, nil, nil, nil}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			room := &Room{
				RoomStatus:           tt.fields.RoomStatus,
				Voters:               tt.fields.Voters,
				turnFinishedCommands: tt.fields.turnFinishedCommands,
				turnStartedCommands:  tt.fields.turnStartedCommands,
				voteCommands:         tt.fields.voteCommands,
			}
			room.OpenVote()
		})
	}
}
