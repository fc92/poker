<template>
    <ion-list v-if="selectedRoom == ''">
        <ion-list-header>
            <h2 v-if="roomList.length > 0">Poker rooms available</h2>
        </ion-list-header>
        <ion-radio-group v-model="selectedRoom" @ionChange="selectRoom">
            <ion-list>
                <ion-item v-for="(room, index) in roomList" :key="index">
                    <ion-radio :value="room.name">
                        <span class="roomDetails">{{ room.name }}
                            <ion-icon :icon="people" class="nbPlayer"></ion-icon>{{ room.nbVoters }}</span>
                    </ion-radio>
                </ion-item>
            </ion-list>
        </ion-radio-group>
    </ion-list>
    <IonLabel v-if="selectedRoom != ''">Room selected:
        <span class="roomDetails">{{ selectedRoom }}
            <ion-icon :icon="people" class="nbPlayer"></ion-icon>{{ selectedRoomNbPlayer }}
        </span>
    </IonLabel>
</template>

<script setup lang="ts">
import { onBeforeMount, ref, watch } from 'vue';
import { useStore } from 'vuex';
import { IonList, IonListHeader, IonItem, IonRadioGroup, IonRadio, IonButton, IonIcon, IonLabel } from '@ionic/vue';
import { people } from 'ionicons/icons';
import { RoomOverview } from '@/room';

const emit = defineEmits();
const store = useStore();
const roomList = ref<RoomOverview[]>([]);
const selectedRoom = ref<string>('');
const selectedRoomNbPlayer = ref<number>(0);

onBeforeMount(async () => {
    try {
        store.dispatch('getRoomList');
        await waitForRoomList();
        roomList.value = store.state.roomList;
    } catch (error) {
        console.error('Failed to load room list:', error);
    }

    async function waitForRoomList() {
        const maxAttempts = 10;
        let attempts = 0;

        while (store.state.roomList.length === 0 && attempts < maxAttempts) {
            await new Promise(resolve => setTimeout(resolve, 1000)); // 1 second pause before each retry
            attempts++;
        }

        if (store.state.roomList.length === 0) {
            console.error('List of rooms could not be loaded after many attempts');
        } else {
            roomList.value = store.state.roomList;
            console.log('Room list was successfully loaded:', roomList.value);
        }
    }
});

const selectRoom = () => {
    if (store.state.roomList.length > 0)
        selectedRoomNbPlayer.value = store.state.roomList.find((r: RoomOverview) => r.name === selectedRoom.value).nbVoters || null;
    emit('update:room', selectedRoom.value);
};

watch(() => store.state.roomList, (newRoomList) => {
    roomList.value = newRoomList;
    if (selectedRoom.value != '')
        selectedRoomNbPlayer.value = newRoomList.find((r: RoomOverview) => r.name === selectedRoom.value).nbVoters || null;
});

</script>

<style scoped>
.roomDetails {
    display: block;
    text-align: right;
    width: 100%;
    /* Additional styling if needed */
    padding: 10px;
    /* Add padding for better visual appearance */
    background-color: lightgray;
    /* Add a background color for better visibility */
}

.nbPlayer {
    margin-left: 50px;
    color: darkblue;
    /* Adjust the value as needed */
}
</style>