<template>
    <ion-list>
        <ion-list-header>
            <h2 v-if="roomList.length > 0">Poker rooms available</h2>
        </ion-list-header>
        <ion-radio-group v-model="selectedRoom">
            <ion-list>
                <ion-item v-for="(room, index) in roomList" :key="index">
                    <ion-radio :value="room">{{ room }}</ion-radio><br>
                </ion-item>
            </ion-list>
        </ion-radio-group>
    </ion-list>
    <ion-button v-show="roomList.length > 0" @click="selectRoom">Select room</ion-button>
</template>

<script setup lang="ts">
import { onBeforeMount, ref } from 'vue';
import { useStore } from 'vuex';
import { IonList, IonListHeader, IonItem, IonRadioGroup, IonRadio, IonButton } from '@ionic/vue';

const emit = defineEmits();
var roomList = ref<string[]>([]);
var selectedRoom = ref<string>('');

onBeforeMount(async () => {
    const store = useStore();
    try {
        store.dispatch('getRoomList');
        await waitForRoomList();
        roomList.value = store.state.roomList;
        selectedRoom.value = roomList.value[0];
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
    emit('update:room', selectedRoom.value);
};


</script>