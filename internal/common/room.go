// Game room
// Game logic independent from underlying transport

package common

import (
	"log"

	"github.com/google/uuid"
	"github.com/kyokomi/emoji/v2"
	"golang.org/x/exp/maps"
)

type RoomVoteStatus int

const (
	VoteOpen = iota
	VoteClosed
)

// commands and associated keyboard shortcuts
const (
	CommandQuit      = "q"
	CommandNotVoting = "n"
	CommandStartVote = "s"
	CommandCloseVote = "v"
	CommandVote1     = "1"
	CommandVote2     = "2"
	CommandVote3     = "3"
	CommandVote5     = "5"
	CommandVote8     = "8"
	CommandVote13    = "d"
)

// vote status
const (
	VoteReceived    = "r"
	VoteNotReceived = ""
	VoteHidden      = "-"
)

type Room struct {
	RoomStatus           RoomVoteStatus `json:"roomStatus"`
	Voters               []*Participant `json:"voters"`
	turnFinishedCommands map[string]string
	turnStartedCommands  map[string]string
	voteCommands         map[string]string
}

func NewRoom() *Room {
	room := Room{
		RoomStatus: VoteClosed,
		Voters:     []*Participant{}}
	room.initCommands()
	return &room
}

func (room *Room) initCommands() {
	room.turnFinishedCommands = map[string]string{}
	room.turnStartedCommands = map[string]string{}
	room.voteCommands = map[string]string{}

	room.turnFinishedCommands[CommandStartVote] = "start new vote"

	room.turnStartedCommands[CommandCloseVote] = "close vote"

	room.voteCommands[CommandNotVoting] = "?"
	room.voteCommands[CommandVote1] = "vote 1"
	room.voteCommands[CommandVote2] = "vote 2"
	room.voteCommands[CommandVote3] = "vote 3"
	room.voteCommands[CommandVote5] = "vote 5"
	room.voteCommands[CommandVote8] = "vote 8"
	room.voteCommands[CommandVote13] = "vote 13"
}

func (room Room) TurnFinishedCommands() map[string]string {
	if room.turnFinishedCommands == nil {
		room.initCommands()
	}
	return room.turnFinishedCommands
}
func (room Room) VoteCommands() map[string]string {
	if room.voteCommands == nil {
		room.initCommands()
	}
	return room.voteCommands
}
func (room Room) TurnStartedCommands() map[string]string {
	if room.turnStartedCommands == nil {
		room.initCommands()
	}
	return room.turnStartedCommands
}

func (room Room) DisplayVotersStatus() []string {
	votersStatus := []string{""}
	for _, voter := range room.Voters {
		switch voter.LastCommand {
		case VoteNotReceived:
			if room.RoomStatus == VoteOpen {
				votersStatus = append(votersStatus, voter.Name, ": [yellow]waiting for vote[white]", emoji.Sprint(" :thinking:\n"))
			} else {
				votersStatus = append(votersStatus, voter.Name, ": [grey]did not vote...[white]\n")
			}
		case CommandNotVoting:
			if room.RoomStatus == VoteOpen {
				votersStatus = append(votersStatus, voter.Name, ": [grey]will not vote...[white]", emoji.Sprint(" :neutral_face:\n"))
			} else {
				votersStatus = append(votersStatus, voter.Name, ": [grey]did not vote...[white]\n")
			}
		case VoteReceived:
			if room.RoomStatus == VoteOpen {
				votersStatus = append(votersStatus, voter.Name, ": [green]vote received[white]", emoji.Sprint(" :slightly_smiling_face:\n"))
			} else {
				votersStatus = append(votersStatus, voter.Name, ": ", voter.Vote, "\n")
			}
		default:
			votersStatus = append(votersStatus, voter.Name, ": [red]unknown status[white]", emoji.Sprint(" :stop_sign:\n"))
		}
	}
	return votersStatus
}

// prepare room for new vote
func (room *Room) OpenVote() {
	room.RoomStatus = VoteOpen
	for i := range room.Voters {
		maps.Clear(room.Voters[i].AvailableCommands)
		maps.Copy(room.Voters[i].AvailableCommands, room.TurnStartedCommands())
		maps.Copy(room.Voters[i].AvailableCommands, room.VoteCommands())
		room.Voters[i].Vote = VoteNotReceived
		room.Voters[i].LastCommand = VoteNotReceived
	}
}

// closing vote
func (room *Room) CloseVote() {
	room.RoomStatus = VoteClosed
	for i := range room.Voters {
		maps.Clear(room.Voters[i].AvailableCommands)
		maps.Copy(room.Voters[i].AvailableCommands, room.TurnFinishedCommands())
		// cannot start vote if player is alone
		if len(room.Voters) < 2 {
			delete(room.Voters[i].AvailableCommands, CommandStartVote)
		}
	}
}

// hide votes if the vote is not closed
func (room Room) FilterVoteData(voterId uuid.UUID) *Room {
	if room.RoomStatus != VoteClosed {
		filteredRoom := Room{}
		filteredRoom.RoomStatus = room.RoomStatus
		for _, voter := range room.Voters {
			vote := VoteHidden
			id := uuid.UUID{}
			if voter.Id == voterId || voter.Vote == room.voteCommands[CommandNotVoting] {
				vote = voter.Vote
				id = voter.Id
			}
			voterForClient := Participant{
				Id:                id,
				Name:              voter.Name,
				Vote:              vote,
				AvailableCommands: voter.AvailableCommands,
				LastCommand:       voter.LastCommand,
			}
			filteredRoom.Voters = append(filteredRoom.Voters, &voterForClient)
		}
		return &filteredRoom
	}
	return &room
}

func (room *Room) UpdateFromParticipant(voterReceived Participant) {
	// add first player
	if len(room.Voters) == 0 {
		room.Voters = append(room.Voters, &voterReceived)
		room.CloseVote()
	} else {
		for i, voter := range room.Voters {
			// update room with data received from player
			if voter.Id == voterReceived.Id {
				room.Voters[i] = &voterReceived
				log.Println("last command from ", room.Voters[i].Name, ": ", voterReceived.LastCommand, " vote: ", voterReceived.Vote)
				switch voterReceived.LastCommand {
				case CommandStartVote:
					room.OpenVote()
				case CommandCloseVote:
					// update status of the received voter based on existing vote existence
					if voterReceived.Vote == VoteNotReceived {
						voterReceived.LastCommand = VoteNotReceived
					} else {
						voterReceived.LastCommand = VoteReceived
					}
					room.CloseVote()
				}
				break
			}

			// add new player
			if i == len(room.Voters)-1 {
				if room.RoomStatus == VoteOpen {
					maps.Copy(voterReceived.AvailableCommands, room.TurnStartedCommands())
					maps.Copy(voterReceived.AvailableCommands, room.VoteCommands())
					room.Voters = append(room.Voters, &voterReceived)
				} else {
					room.Voters = append(room.Voters, &voterReceived)
					// update command menu
					room.CloseVote()
				}
			}
		}

	}
	room.updateFromVotes()
}

// apply vote rules
func (room *Room) updateFromVotes() {
	switch room.RoomStatus {
	case VoteOpen:
		// if all votes received, close vote
		allVotesReceived := true
		for _, voter := range room.Voters {
			if voter.LastCommand == VoteNotReceived {
				allVotesReceived = false
			}
		}
		if allVotesReceived {
			room.CloseVote()
		}
	default:
		// nothing to do for the moment
	}
}

func (room *Room) UpdateFromHub() {
	if len(room.Voters) < 2 {
		room.CloseVote()
	}
}

func (room Room) NbVotesReceived() int {
	nbVotes := 0
	for _, voter := range room.Voters {
		if voter.LastCommand != VoteNotReceived {
			nbVotes++
		}
	}
	return nbVotes
}
