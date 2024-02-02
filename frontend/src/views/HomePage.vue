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
            <div class="local-player">
              <room-selector v-if="!store.state.roomSelected" @update:room="handleRoomUpdate"></room-selector>
              <ion-label v-else>Poker room selected: {{ store.state.roomSelected }}</ion-label>
              <name-selector v-if="store.state.roomSelected" @update:player="handlePlayerUpdate"></name-selector>
            </div>
          </div>
          <br>
        </div>
      </div>

    </ion-content>
    <IonFooter>
      <ExitButton></ExitButton>
    </IonFooter>
  </ion-page>
</template>

<script setup lang="ts">
import { onBeforeUnmount, computed } from 'vue';
import { IonContent, IonFooter, IonHeader, IonPage, IonLabel } from '@ionic/vue';
import ServerSelector from '@/components/ServerSelector.vue';
import RoomSelector from '@/components/RoomSelector.vue';
import NameSelector from '@/components/NameSelector.vue';
import ExitButton from '@/components/ExitButton.vue';
import { Participant } from '@/participant';
import { useStore } from 'vuex';
import { useRouter } from 'vue-router';

const router = useRouter();

const store = useStore();

const websocket = computed(() => store.state.websocket);

const handleServerValueUpdate = (newServerValue: string) => {
  store.dispatch('handleServerValueUpdate', newServerValue);
};

const handleRoomUpdate = (newRoomName: string) => {
  store.state.roomSelected = newRoomName;
  store.commit('setRoom', newRoomName);
};

const handlePlayerUpdate = (participant: Participant) => {
  console.log('Valeur du participant id mise à jour:', participant.id);
  console.log('Valeur du participant name mise à jour:', participant.name);
  if (websocket.value) {
    const message = JSON.stringify({
      id: participant.id,
      name: participant.name,
      vote: "",
      available_commands: {},
      last_command: "",
      room: store.state.roomSelected
    });
    websocket.value.send(message);
    store.commit('setLocalParticipantId', participant.id);
    router.push('/pokertable');
  } else {
    console.error('WebSocket is not connected');
  }
};

onBeforeUnmount(() => {
  if (websocket.value) {
    websocket.value.close();
  }
});
</script>

<style scoped>
.body {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
}

.room {
  text-align: center;
}

.selector-container {
  margin-top: 20px;
  /* Adjust as needed for spacing between selectors */
}
</style>