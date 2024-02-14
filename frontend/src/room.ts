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
    name: string;
}

export interface RoomOverview {
    name: string;
    nbVoters: number;
}
export interface RoomRequest {
    roomList: RoomOverview[]
}
