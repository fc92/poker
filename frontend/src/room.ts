import { Participant } from "./participant";
export enum RoomVoteStatus {
    VoteOpen,
    VoteClosed,
}

export enum RoomRequestAction {
    ActionGetRoomList,
}

export interface Room {
    roomStatus: RoomVoteStatus;
    voters: Participant[];
    turnFinishedCommands: Record<string, string>;
    turnStartedCommands: Record<string, string>;
    voteCommands: Record<string, string>;
    name: string;
}

export interface RoomRequest {
    action: RoomRequestAction,
    roomList: string[]
}
