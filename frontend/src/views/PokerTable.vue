<template>
  <ion-page>
    <ion-header :translucent="true" class="ion-margin-top">
      <h1 class="ion-text-center">Team poker
        <br>You are: <ion-label v-if="localParticipant" v-text="localParticipant.name"></ion-label>
      </h1>
    </ion-header>

    <ion-content :fullscreen="true">
      <ion-grid>
        <ion-row>
          <ion-col v-if="room.roomStatus === RoomVoteStatus.VoteClosed" size="12">
            <ion-label v-if="displayVoteResults">Team votes</ion-label>
            <ion-label v-else>Team</ion-label>
            <player v-for="participant in participants" :key="participant.id" :player="participant"
              :displayVote="displayVoteResults" :isCurrentUser="participant.id === localParticipantId" />
            <BarChart v-if="displayVoteResults" :player-votes="voteResults" :barColors="barColors" />
          </ion-col>

          <ion-col v-if="room.roomStatus === RoomVoteStatus.VoteOpen" size=4>
            <ion-item v-if="localParticipant">
              <ion-grid>
                <ion-col size=6><ion-row>Your vote:</ion-row>
                  <ion-row><ion-radio-group v-model="localVote" @ionChange="onVoteChange">
                      <ion-item v-for="[command, label] in Object.entries(room.voteCommands)" :key="command">
                        <ion-radio v-if="label !== 'Close vote'" :value="label">{{ label }}</ion-radio>
                      </ion-item>
                    </ion-radio-group></ion-row>
                </ion-col>
              </ion-grid>
            </ion-item>
          </ion-col>
          <ion-col v-if="room.roomStatus === RoomVoteStatus.VoteOpen" size=2><ion-label>Team votes</ion-label>
            <player v-for="participant in participants" :key="participant.id" :player="participant" :displayVote=false
              :isCurrentUser="participant.id === localParticipantId" />
          </ion-col>
          <ion-col v-if="room.roomStatus === RoomVoteStatus.VoteOpen" size=2>
            <ProgressBar :progress="voteProgress"></ProgressBar>
          </ion-col>
        </ion-row>
      </ion-grid>
    </ion-content>
    <IonFooter>
      <ion-grid>
        <ion-row>
          <ion-col v-if="room.roomStatus === RoomVoteStatus.VoteClosed" size=8>
            <ion-button v-if="room.voters.length > 1" @click="startGame">
              <ion-icon :icon="playOutline"></ion-icon> Start new vote
            </ion-button>
          </ion-col>
          <ion-col v-if="room.roomStatus === RoomVoteStatus.VoteOpen" size=8>
            <ion-button @click="closeVote">
              <ion-icon :icon="playOutline"></ion-icon> Close vote
            </ion-button>
          </ion-col>
          <ion-col size=2>
            <ExitButton></ExitButton>
          </ion-col>
        </ion-row>
      </ion-grid>
    </IonFooter>
  </ion-page>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useStore } from 'vuex';
import { IonPage, IonContent, IonHeader, IonButton, IonIcon, IonFooter, IonLabel, IonItem, IonGrid, IonCol, IonRow } from '@ionic/vue';
import { playOutline } from 'ionicons/icons';
import { IonRadio, IonRadioGroup } from '@ionic/vue';
import Player from '@/components/Player.vue';
import BarChart from '@/components/BarChart.vue'
import ProgressBar from '@/components/ProgressBar.vue';
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
var displayVoteResults = computed(() => store.state.voteResults.some((num: number) => num !== 0));
var voteProgress = computed(() => {
  let nbVotes = 0;

  for (const voter of store.state.room.voters) {
    if (voter.last_command !== "") {
      nbVotes++;
    }
  }

  return [nbVotes, store.state.room.voters.length - nbVotes];
});

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
