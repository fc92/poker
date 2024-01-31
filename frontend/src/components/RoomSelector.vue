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
    <ion-button @click="selectRoom">Select room</ion-button>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue';
import { useStore } from 'vuex';
import { IonList, IonListHeader, IonItem, IonRadioGroup, IonRadio, IonButton } from '@ionic/vue';

const store = useStore();
const emit = defineEmits();

store.dispatch('getRoomList');
var roomList: string[] = computed(() => store.state.roomList).value;

const selectedRoom = ref(roomList[0]);
const selectRoom = () => {
    emit('update:room', selectedRoom);
};

</script>
