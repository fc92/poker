import { createStore } from 'vuex';
import { Room, RoomVoteStatus } from '@/room';
import { Participant } from '@/participant';


export default createStore({
    state: {
        websocket: null as WebSocket | null,
        serverSelected: '',
        room: {
            roomStatus: RoomVoteStatus.VoteClosed,
            voters: [] as Participant[],
            turnFinishedCommands: {} as Record<string, string>,
            turnStartedCommands: {} as Record<string, string>,
            voteCommands: {} as Record<string, string>,
        } as Room,
        localParticipantId: ''
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
        setRoomStatus(state, roomStatus: RoomVoteStatus) {
            state.room.roomStatus = roomStatus;
            const roomLog = JSON.stringify(state.room);
            console.info(`room  ${roomLog}`);
        },
        setParticipants(state, voters: Participant[]) {
            state.room.voters = voters;
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
            } else {
                console.error(`No participant found with ID ${payload.id}`);
            }
        }
    },
    actions: {
        connectToWebSocket({ commit }, serverAddress: string) {
            try {
                const websocket = new WebSocket(`ws://${serverAddress}/ws`);
                websocket.addEventListener('open', () => {
                    console.log('WebSocket is connected');
                });

                websocket.addEventListener('message', (event) => {
                    console.log('Message from server: ', event.data);
                    try {
                        const data = JSON.parse(event.data);
                        commit('setRoomStatus', data.roomStatus);
                        commit('setParticipants', data.voters);
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
        startGame({ state, commit }, localParticipant: Participant) {
            if (state.websocket) {
                localParticipant.last_command = 's';
                const message = JSON.stringify(localParticipant);
                state.websocket.send(message);
                console.log('participant update sent to server: ' + message)
            } else {
                console.error('WebSocket is not connected');
            }
        },
        updateVote({ state, commit }, payload) {
            // Votre logique pour mettre à jour le vote ici, par exemple envoyer le vote via WebSocket
            // ...
        },

    },
});
