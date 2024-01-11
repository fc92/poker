<template>
  <ion-page>
    <ion-header :translucent="true">
      <h1>Unbiased poker</h1>
    </ion-header>

    <ion-content :fullscreen="true">
      <div>
        <div>
          <div v-if="room.roomStatus === RoomVoteStatus.VoteClosed">
            <player v-for="participant in participants" :key="participant.id" :player="participant"
              :isCurrentUser="participant.id === localParticipantId" />
            <BarChart :player-votes="voteResults" />
            <ion-button @click="startGame">
              <ion-icon :icon="playOutline"></ion-icon> Start Game
            </ion-button>
          </div>

          <div v-else-if="room.roomStatus === RoomVoteStatus.VoteOpen">
            <player v-for="participant in participants" :key="participant.id" :player="participant"
              :isCurrentUser="participant.id === localParticipantId" />
            <ion-item v-if="localParticipant">
              <ion-select v-model="localVote" @ionChange="onVoteChange" label="Vote">
                <ion-select-option v-for="[command, label] in Object.entries(room.voteCommands)" :key="command"
                  :value="label">
                  {{ label }}
                </ion-select-option>
              </ion-select>
            </ion-item>
          </div>
        </div>
      </div>
    </ion-content>
    <IonFooter>
      <ExitButton></ExitButton>
    </IonFooter>
  </ion-page>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useStore } from 'vuex';
import { IonPage, IonContent, IonHeader, IonButton, IonIcon, IonFooter, IonLabel, IonItem, IonSelect, IonSelectOption } from '@ionic/vue';
import { playOutline } from 'ionicons/icons';
import Player from '@/components/Player.vue';
import BarChart from '@/components/BarChart.vue'
import ExitButton from '@/components/ExitButton.vue';
import { Participant } from '@/participant';
import { Room, RoomVoteStatus } from '@/room';

const store = useStore();

const room: Room = computed(() => store.state.room).value;
var voteResults = computed(() => store.state.voteResults);
const participants = computed(() => store.state.room.voters);
const localParticipantId = computed(() => store.state.localParticipantId);

var localParticipant = computed(() => participants.value.find((p: Participant) => p.id === localParticipantId.value));
const localVote = computed({
  get: () => localParticipant.value?.vote || '',
  set: (value) => {
    store.commit('updateLocalVote', { id: localParticipantId.value, vote: value });
  }
});

const startGame = () => {
  // Ajoutez la logique pour démarrer le jeu
  store.dispatch('startGame', localParticipant.value);
};

const onVoteChange = () => {
  store.dispatch('updateVote', { localParticipant: localParticipant.value });
};
</script>
