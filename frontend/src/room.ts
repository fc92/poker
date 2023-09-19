import { Participant } from "./participant";
export enum RoomVoteStatus {
    VoteOpen,
    VoteClosed,
}

export interface Room {
    roomStatus: RoomVoteStatus;
    voters: Participant[];
    turnFinishedCommands: Record<string, string>;
    turnStartedCommands: Record<string, string>;
    voteCommands: Record<string, string>;
}

export const voteOptions = [
    { value: 'n', label: 'Not Voting' },
    { value: '1', label: 'Vote 1' },
    { value: '2', label: 'Vote 2' },
    { value: '3', label: 'Vote 3' },
    { value: '5', label: 'Vote 5' },
    { value: '8', label: 'Vote 8' },
    { value: '13', label: 'Vote 13' },
]

// Vote status
export enum VoteStatus {
    Received = "r",
    NotReceived = "",
    Hidden = "-",
}
