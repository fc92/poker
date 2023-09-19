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
      <ExitButton></ExitButton>
    </IonFooter>
  </ion-page>
</template>

<script setup lang="ts">
import { onBeforeUnmount, computed } from 'vue';
import { IonContent, IonFooter, IonHeader, IonPage } from '@ionic/vue';
import ServerSelector from '@/components/ServerSelector.vue';
import NameSelector from '@/components/NameSelector.vue';
import ExitButton from '@/components/ExitButton.vue';
import { Participant } from '@/participant';
import { useStore } from 'vuex';
import { useRouter } from 'vue-router';

const router = useRouter();

const store = useStore();

const serverSelected = computed(() => store.state.serverSelected);
const websocket = computed(() => store.state.websocket);
const localParticipantId = computed(() => store.state.localParticipantId);

const handleServerValueUpdate = (newServerValue: string) => {
  store.dispatch('handleServerValueUpdate', newServerValue);
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
      last_command: ""
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
