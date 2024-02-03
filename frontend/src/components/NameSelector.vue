<template>
  <div>
    <p>Enter your name:</p>
    <ion-item>
      <ion-input v-model="playerName" @ion-input="onInputchange" placeholder="Player name" autofocus
        @keyup.enter="enterName"></ion-input>
    </ion-item>
    <ion-button v-show="isButtonVisible" @click="enterName">Join game</ion-button>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { IonItem, IonInput, IonButton } from "@ionic/vue";
import { v4 as uuidv4 } from "uuid";
import { Participant } from "@/participant";
import { useStore } from 'vuex';

const store = useStore();
const emit = defineEmits();

const playerName = ref<string>('');
const isButtonVisible = ref(false);

const onInputchange = () => {
  if (playerName.value.trim() != "" && store.state.roomSelected != "") {
    isButtonVisible.value = true;
  }
  else {
    isButtonVisible.value = false;
  }
}


const enterName = () => {
  if (store.state.roomSelected != "") {
    const player: Participant = {
      id: uuidv4(),
      name: playerName.value.trim(),
      available_commands: {},
      last_command: "",
      vote: "",
      room: store.state.room.name == null ? "" : store.state.room.name
    };
    emit("update:player", player);
  }
};

</script>