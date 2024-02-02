import { createStore } from 'vuex';
import { Room, RoomOverview, RoomRequest, RoomVoteStatus } from '@/room';
import { Participant } from '@/participant';

export default createStore({
    state: {
        websocket: null as WebSocket | null,
        serverSelected: '',
        roomSelected: '',
        room: {
            roomStatus: RoomVoteStatus.VoteClosed,
            voters: [] as Participant[],
            turnFinishedCommands: {} as Record<string, string>,
            turnStartedCommands: {} as Record<string, string>,
            voteCommands: {} as Record<string, string>,
        } as Room,
        localParticipantId: '',
        voteResults: [0, 0, 0, 0, 0, 0, 0, 0] as number[],
        roomList: [] as RoomOverview[]
    },
    mutations: {
        setWebSocket(state, websocket: WebSocket) {
            state.websocket = websocket;
        },
        setServerSelected(state, serverSelected: string) {
            state.serverSelected = serverSelected;
        },
        clearWebSocket(state) {
            if (state.websocket) {
                state.websocket.close();
            }
            state.websocket = null;
        },
        setRoomList(state, roomList: RoomOverview[]) {
            state.roomList = roomList;
        },
        setRoom(state, newRoomName: string) {
            state.roomSelected = newRoomName;
        },
        setRoomStatus(state, roomStatus: RoomVoteStatus) {
            if (state.room.roomStatus === RoomVoteStatus.VoteClosed
                && roomStatus === RoomVoteStatus.VoteOpen) {
                resetVotes(state);
            }
            state.room.roomStatus = roomStatus;
            const roomLog = JSON.stringify(state.room);
            console.info(`room  ${roomLog}`);
        },
        setParticipants(state, voters: Participant[]) {
            state.room.voters = voters;
            state.room.voteCommands = voters[0].available_commands;
            if (state.room.roomStatus == RoomVoteStatus.VoteClosed) {
                voters.forEach((voter) => {
                    switch (voter.vote) {
                        case "vote 1":
                            state.voteResults[0] += 1;
                            break;
                        case "vote 2":
                            state.voteResults[1] += 1;
                            break;
                        case "vote 3":
                            state.voteResults[2] += 1;
                            break;
                        case "vote 5":
                            state.voteResults[3] += 1;
                            break;
                        case "vote 8":
                            state.voteResults[4] += 1;
                            break;
                        case "vote 13":
                            state.voteResults[5] += 1;
                            break;
                        case "vote 21":
                            state.voteResults[6] += 1;
                            break;
                        case "not voting":
                            state.voteResults[7] += 1;
                            break;
                        case "":
                            voter.vote = 'did not vote...'
                            break;
                    }

                });
            }
            const roomLog = JSON.stringify(state.room);
            console.info(`room  ${roomLog}`);
        },
        setLocalParticipantId(state, id: string) {
            state.localParticipantId = id;
        },
        updateLocalVote(state, payload) {
            const participant = state.room.voters.find(p => p.id === payload.id);
            if (participant) {
                participant.vote = payload.vote;
                // participant.room = state.room.name;
            } else {
                console.error(`No participant found with ID ${payload.id}`);
            }
        }
    },
    actions: {
        connectToWebSocket({ state, commit, dispatch }, serverAddress: string) {
            try {
                const websocket = new WebSocket(`ws://${serverAddress}/ws`);
                websocket.addEventListener('open', () => {
                    console.log('WebSocket is connected');
                });

                websocket.addEventListener('message', (event) => {
                    console.log('Message from server: ', event.data);
                    try {
                        const data = JSON.parse(event.data);
                        if (data.voters) {
                            commit('setRoomStatus', <RoomVoteStatus>data.roomStatus);
                            commit('setParticipants', <Participant[]>data.voters);
                        }
                        else if (data.voters == null) {
                            commit('setRoomList', <RoomOverview[]>data.RoomList)
                        }
                        else {
                            console.error('Server response unknown');
                        }
                    } catch (error) {
                        console.error('Error parsing message data:', error);
                    }
                });

                websocket.addEventListener('close', () => {
                    console.log('WebSocket is closed');
                });

                websocket.addEventListener('error', (error) => {
                    console.error('WebSocket error:', error);
                });
                commit('setWebSocket', websocket);
            } catch (error) {
                console.error('Failed to connect to WebSocket:', error);
            }
        },
        handleServerValueUpdate({ dispatch, commit }, newServerValue: string) {
            console.log('Valeur du serveur mise à jour:', newServerValue);
            dispatch('connectToWebSocket', newServerValue)
                .then(() => {
                    commit('setServerSelected', newServerValue);
                })
                .catch(error => {
                    console.error('Failed to update server value:', error);
                });
        },
        handleExitClick({ state, commit }) {
            commit('clearWebSocket');
            commit('setServerSelected', '');
            console.info('Exited');
        },
        startGame({ state }, localParticipant: Participant) {
            resetVotes(state); // display fresh vote
            if (state.websocket) {
                localParticipant.last_command = 's';
                const message = JSON.stringify(localParticipant);
                state.websocket.send(message);
                console.log('participant update sent to server: ' + message)
            } else {
                console.error('WebSocket is not connected');
            }
        },
        closeVote({ state }, localParticipant: Participant) {
            if (state.websocket) {
                localParticipant.last_command = 'v';
                const message = JSON.stringify(localParticipant);
                state.websocket.send(message);
                console.log('participant update sent to server: ' + message)
            } else {
                console.error('WebSocket is not connected');
            }
        },
        updateVote({ state }, localParticipant: Participant) {
            if (state.websocket) {
                localParticipant.last_command = 'r';
                const message = JSON.stringify(localParticipant);
                state.websocket.send(message);
                console.log('participant update sent to server: ' + message)
            } else {
                console.error('WebSocket is not connected');
            }
        },
        getRoomList({ state }) {
            if (state.websocket) {
                const roomReq: RoomRequest = {
                    roomList: []
                }
                const message = JSON.stringify(roomReq);
                sendMessage(state.websocket, message);
            } else {
                console.error('WebSocket is not connected');
            }
        },
    },

});
function resetVotes(state: { websocket: WebSocket | null; serverSelected: string; room: Room; localParticipantId: string; voteResults: number[]; }) {
    state.voteResults = [0, 0, 0, 0, 0, 0, 0, 0];
}

function sendMessage(ws: WebSocket, message: string): Promise<void> {
    return new Promise<void>((resolve, reject) => {
        if (ws.readyState === WebSocket.OPEN) {
            try {
                ws.send(message);
                console.log('room request sent to server: ' + message)
                resolve();
            } catch (error) {
                console.error('websocket send error')
            }
        } else if (ws.readyState === WebSocket.CONNECTING) {
            ws.addEventListener('open', () => {
                try {
                    ws.send(message);
                    console.log('room request sent to server: ' + message)
                    resolve();
                } catch (error) {
                    console.error('websocket send error')
                }
            });
        } else {
            reject(new Error('WebSocket connection not open.'));
        }
    });
}
