<template>
    Enter your name:
    <ion-item>
        <ion-input v-model="playerName" placeholder="player name"></ion-input>
    </ion-item>
    <ion-button @click="enterName">Join game</ion-button>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { IonContent, IonItem, IonInput, IonButton } from '@ionic/vue';
import { defineProps, defineEmits } from 'vue';
import { v4 as uuidv4 } from 'uuid';
import { Player } from '@/player';

const emit = defineEmits();

const playerName = ref('');

const enterName = () => {
    const player: Player = {
        id: uuidv4(),
        name: playerName.value,
    };
    emit('update:player', player);
};

const props = defineProps({
    player: {
        type: Object as () => Player,
        required: false,
        default: () => ({ id: '', name: '' }) as Player,
    },
});
</script>