<template>
  Enter your name:
  <ion-item>
    <ion-input v-model="playerName" placeholder="player name"></ion-input>
  </ion-item>
  <ion-button @click="enterName">Join game</ion-button>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { IonContent, IonItem, IonInput, IonButton } from "@ionic/vue";
import { v4 as uuidv4 } from "uuid";
import { Participant } from "@/participant";

const emit = defineEmits();

const playerName = ref("");

const enterName = () => {
  const player: Participant = {
    id: uuidv4(),
    name: playerName.value,
    available_commands: {},
    last_command: "",
    vote: "",
  };
  emit("update:player", player);
};

const props = defineProps({
  player: {
    type: Object as () => Participant,
    required: false,
    default: () => ({ id: "", name: "" } as Participant),
  },
});
</script>