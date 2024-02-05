// Game room
// Game logic independent from underlying transport

package common

import (
	"github.com/google/uuid"
	"github.com/kyokomi/emoji/v2"
	"github.com/rs/zerolog/log"
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
	CommandVote21    = "m"
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
	Name                 string
}

type RoomOverview struct {
	Name     string `json:"name"`
	NbVoters int    `json:"nbVoters"`
}
type RoomRequest struct {
	RoomList []RoomOverview
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

	room.turnFinishedCommands[CommandStartVote] = "Start new vote"

	room.turnStartedCommands[CommandCloseVote] = "Close vote"

	room.voteCommands[CommandNotVoting] = "not voting"
	room.voteCommands[CommandVote1] = "vote 1"
	room.voteCommands[CommandVote2] = "vote 2"
	room.voteCommands[CommandVote3] = "vote 3"
	room.voteCommands[CommandVote5] = "vote 5"
	room.voteCommands[CommandVote8] = "vote 8"
	room.voteCommands[CommandVote13] = "vote 13"
	room.voteCommands[CommandVote21] = "vote 21"
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
		room.Voters[i].AvailableCommands = make(map[string]string)
		for k := range room.TurnStartedCommands() {
			room.Voters[i].AvailableCommands[k] = room.TurnStartedCommands()[k]
		}
		for l := range room.VoteCommands() {
			room.Voters[i].AvailableCommands[l] = room.VoteCommands()[l]
		}
		room.Voters[i].Vote = VoteNotReceived
		room.Voters[i].LastCommand = VoteNotReceived
	}
}

// closing vote
func (room *Room) CloseVote() {
	room.RoomStatus = VoteClosed
	for i := range room.Voters {
		room.Voters[i].AvailableCommands = make(map[string]string)
		for k := range room.TurnFinishedCommands() {
			room.Voters[i].AvailableCommands[k] = room.TurnFinishedCommands()[k]
		}
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
				RoomName:          voter.RoomName,
			}
			filteredRoom.Voters = append(filteredRoom.Voters, &voterForClient)
		}
		return &filteredRoom
	}
	return &room
}

func (room *Room) UpdateFromParticipant(voterReceived Participant) (isNewPlayer bool) {
	isNewPlayer = false
	// add first player
	if len(room.Voters) == 0 {
		room.Voters = append(room.Voters, &voterReceived)
		room.CloseVote()
	} else {
		for i, voter := range room.Voters {
			// update room with data received from player
			if voter.Id == voterReceived.Id {
				// update status of the received voter based on existing vote existence
				updateRoomFromReceivedPlayer(room, i, voterReceived)
				break
			}

			// add new player, not found in the room
			if i == len(room.Voters)-1 {
				// update command menu
				room.addPlayer(voterReceived)
				isNewPlayer = true
			}
		}

	}
	room.updateFromVotes()
	return
}

func updateRoomFromReceivedPlayer(room *Room, voterPosition int, voterReceived Participant) {
	room.Voters[voterPosition] = &voterReceived
	log.Debug().Msgf("last command from %s: %s, vote: %s", room.Voters[voterPosition].Name, voterReceived.LastCommand, voterReceived.Vote)
	switch voterReceived.LastCommand {
	case CommandStartVote:
		room.OpenVote()
	case CommandCloseVote:

		if voterReceived.Vote == VoteNotReceived {
			voterReceived.LastCommand = VoteNotReceived
		} else {
			voterReceived.LastCommand = VoteReceived
		}
		room.CloseVote()
	}
}

func (room *Room) addPlayer(voterReceived Participant) {
	room.Voters = append(room.Voters, &voterReceived)
	// init menu
	if room.RoomStatus == VoteOpen {
		for k := range room.TurnStartedCommands() {
			room.Voters[len(room.Voters)-1].AvailableCommands[k] = room.TurnStartedCommands()[k]
		}
		for l := range room.VoteCommands() {
			room.Voters[len(room.Voters)-1].AvailableCommands[l] = room.VoteCommands()[l]
		}
	} else {
		room.CloseVote()
	}
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
	numVoters := len(room.Voters)

	if numVoters == 0 {
		log.Info().Msgf("No more player in the poker room.")
	}
	if numVoters < 2 {
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
