package common

import (
	"github.com/google/uuid"
)

type Participant struct {
	Id                uuid.UUID         `json:"id"`
	Name              string            `json:"name"`
	Vote              string            `json:"vote"`
	AvailableCommands map[string]string `json:"available_commands"`
	LastCommand       string            `json:"last_command"`
	RoomName          string            `json:"room"`
}

func CreateVoter(voterName string) *Participant {
	return &Participant{
		Id:                uuid.New(),
		Name:              voterName,
		AvailableCommands: map[string]string{},
		LastCommand:       VoteNotReceived,
		Vote:              VoteNotReceived,
		RoomName:          "",
	}
}

// for local participant apply server updates to keep display up to date
func (localVoter *Participant) UpdateLocalPlayerFromServer(roomFromServer *Room) {
	for _, roomVoter := range roomFromServer.Voters {
		if localVoter.Id == roomVoter.Id {
			localVoter.AvailableCommands = roomVoter.AvailableCommands
			localVoter.LastCommand = roomVoter.LastCommand
			localVoter.Vote = roomVoter.Vote
		}
	}
}
