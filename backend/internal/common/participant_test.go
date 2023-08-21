package common

import (
	"testing"

	"github.com/google/uuid"
)

func TestCreateVoter(t *testing.T) {
	type args struct {
		voterName string
	}
	tests := []struct {
		name string
		args args
		want *Participant
	}{
		{
			name: "simple user",
			args: args{"userName"},
			want: &Participant{Name: "userName"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CreateVoter(tt.args.voterName)
			if got.Name != tt.want.Name {
				t.Errorf("CreateVoter() = %v, want %v", got.Name, tt.want.Name)
			}
		})
	}
}

func TestParticipant_UpdateLocalPlayerFromServer(t *testing.T) {
	type fields struct {
		Id                uuid.UUID
		Name              string
		Vote              string
		AvailableCommands map[string]string
		LastCommand       string
	}
	testId := uuid.New()
	type args struct {
		roomFromServer *Room
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "basic update",
			fields: fields{Id: testId, Name: "test"},
			args: args{
				roomFromServer: &Room{
					RoomStatus: VoteClosed,
					Voters:     []*Participant{{Id: testId, Name: "test"}},
				}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			localVoter := &Participant{
				Id:                tt.fields.Id,
				Name:              tt.fields.Name,
				Vote:              tt.fields.Vote,
				AvailableCommands: tt.fields.AvailableCommands,
				LastCommand:       tt.fields.LastCommand,
			}
			localVoter.UpdateLocalPlayerFromServer(tt.args.roomFromServer)
		})
	}
}
