<template>
  <ion-page>
    <ion-header :translucent="true">
      <h1>Welcome to poker</h1>
    </ion-header>

    <ion-content :fullscreen="true">

      <div class="body">
        <div class="room">
          <div>
            <div class="local-player">
              <room-selector @update:room="handleRoomUpdate"></room-selector>
              <name-selector @update:player="handlePlayerUpdate"></name-selector>
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
import { IonContent, IonFooter, IonHeader, IonPage } from '@ionic/vue';
import RoomSelector from '@/components/RoomSelector.vue';
import NameSelector from '@/components/NameSelector.vue';
import ExitButton from '@/components/ExitButton.vue';
import { Participant } from '@/participant';
import { useStore } from 'vuex';
import { useRouter } from 'vue-router';

const router = useRouter();

const store = useStore();
store.dispatch('initializeWebSocketConnection');

const websocket = computed(() => store.state.websocket);

const handleRoomUpdate = (newRoomName: string) => {
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
    router.push({
      name: 'PokerTable',
      params: {
        username: participant.name,
        room: store.state.roomSelected
      }
    });
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
h1 {
  text-align: center;
}

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
}
</style>