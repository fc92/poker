import { createStore } from 'vuex';

export default createStore({
    state: {
        websocket: null as WebSocket | null,
        serverSelected: '',
    },
    mutations: {
        setWebSocket(state, websocket) {
            state.websocket = websocket;
        },
        setServerSelected(state, serverSelected) {
            state.serverSelected = serverSelected;
        },
        clearWebSocket(state) {
            if (state.websocket) {
                state.websocket.close();
            }
            state.websocket = null;
        },
    },
    actions: {
        connectToWebSocket({ commit }, serverAddress) {
            const websocket = new WebSocket(`ws://${serverAddress}/ws`);
            websocket.addEventListener('open', () => {
                console.log('WebSocket is connected');
            });

            websocket.addEventListener('message', (event) => {
                console.log('Message from server: ', event.data);
            });

            websocket.addEventListener('close', () => {
                console.log('WebSocket is closed');
            });

            websocket.addEventListener('error', (error) => {
                console.error('WebSocket error:', error);
            });
            commit('setWebSocket', websocket);
        },
        handleServerValueUpdate({ dispatch, commit }, newServerValue) {
            console.log('Valeur du serveur mise à jour:', newServerValue);
            dispatch('connectToWebSocket', newServerValue);
            commit('setServerSelected', newServerValue);
        },
        handleExitClick({ state, commit }) {
            commit('clearWebSocket');
            commit('setServerSelected', '');
        },
    },
});
