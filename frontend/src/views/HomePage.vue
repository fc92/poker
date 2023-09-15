<template>
  <ion-page>
    <ion-header :translucent="true">
      <h1>Unbiased poker</h1>
    </ion-header>

    <ion-content :fullscreen="true">

      <div class="body">
        <div class="room">
          <div v-if="store.state.serverSelected == ''">
            <server-selector @update:serverValue="handleServerValueUpdate"></server-selector>
          </div>
          <div v-if="store.state.serverSelected !== ''">
            <h2>Now in poker room: {{ store.state.serverSelected }}</h2>
          </div>
        </div>
        <br>
        <div class="local-player" v-if="store.state.serverSelected">
          <name-selector @update:player="handlePlayerUpdate"></name-selector>
        </div>
      </div>

    </ion-content>
    <IonFooter>
      <div v-if="store.state.serverSelected !== ''">
        <ion-button @click="handleExitClick">
          <ion-icon :icon="exit"></ion-icon>
        </ion-button>
      </div>

    </IonFooter>
  </ion-page>
</template>

<script setup lang="ts">
import { ref, onBeforeUnmount } from 'vue';
import { IonContent, IonFooter, IonHeader, IonPage, IonButton, IonIcon } from '@ionic/vue';
import ServerSelector from '@/components/ServerSelector.vue';
import NameSelector from '@/components/NameSelector.vue';
import { Player } from '@/player'
import { exit } from 'ionicons/icons';
import { useStore } from 'vuex';

const websocket = ref<WebSocket | null>(null);


const store = useStore();
const handleServerValueUpdate = (newServerValue: string) => {
  store.dispatch('handleServerValueUpdate', newServerValue);
};
const handlePlayerUpdate = (player: Player) => {
  console.log('Valeur du player id mise à jour:', player.id);
  console.log('Valeur du player name mise à jour:', player.name);
  if (store.state.websocket) {
    const message = JSON.stringify({
      id: player.id,
      name: player.name,
      vote: "",
      available_commands: {},
      last_command: ""
    });
    store.state.websocket.send(message);
  } else {
    console.error('WebSocket is not connected');
  }
};
const handleExitClick = () => {
  store.dispatch('handleExitClick');
};

onBeforeUnmount(() => {
  if (store.state.websocket) {
    store.state.websocket.close();
  }
});
</script>
