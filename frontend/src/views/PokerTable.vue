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
            <BarChart :player-votes="voteResults" :barColors="barColors" />
          </div>

          <div v-else-if="room.roomStatus === RoomVoteStatus.VoteOpen">
            <player v-for="participant in participants" :key="participant.id" :player="participant"
              :isCurrentUser="participant.id === localParticipantId" />
            <ion-item v-if="localParticipant">
              <ion-radio-group v-model="localVote" @ionChange="onVoteChange">
                <ion-item v-for="[command, label] in Object.entries(room.voteCommands)" :key="command">
                  <ion-radio v-if="label !== 'Close vote'" :value="label">{{ label }}</ion-radio>
                </ion-item>
              </ion-radio-group>
            </ion-item>
          </div>
        </div>
      </div>
    </ion-content>
    <IonFooter>
      <div v-if="room.roomStatus === RoomVoteStatus.VoteClosed">
        <ion-button @click="startGame">
          <ion-icon :icon="playOutline"></ion-icon> Start new vote
        </ion-button>
      </div>
      <div v-if="room.roomStatus === RoomVoteStatus.VoteOpen">
        <ion-button @click="closeVote">
          <ion-icon :icon="playOutline"></ion-icon> Close vote
        </ion-button>
      </div>
      <ExitButton></ExitButton>
    </IonFooter>
  </ion-page>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useStore } from 'vuex';
import { IonPage, IonContent, IonHeader, IonButton, IonIcon, IonFooter, IonLabel, IonItem, IonSelect, IonSelectOption } from '@ionic/vue';
import { playOutline } from 'ionicons/icons';
import { IonRadio, IonRadioGroup } from '@ionic/vue';
import Player from '@/components/Player.vue';
import BarChart from '@/components/BarChart.vue'
import ExitButton from '@/components/ExitButton.vue';
import { Participant } from '@/participant';
import { Room, RoomVoteStatus } from '@/room';

const store = useStore();
const barColors: string[] = [
  'white',
  'saddlebrown',
  'dodgerblue',
  'limegreen',
  'darkturquoise',
  'gold',
  'fuchsia',
  'deepskyblue'
];

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
  store.dispatch('startGame', localParticipant.value);
};

const closeVote = () => {
  store.dispatch('closeVote', localParticipant.value);
};

const onVoteChange = () => {
  store.dispatch('updateVote', { localParticipant: localParticipant.value });
};
</script>