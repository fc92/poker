<template>
  <ion-page>
    <ion-header :translucent="true" class="ion-margin-top">
      <h1 class="ion-text-center">Poker room: {{ store.state.roomSelected }}
      </h1>
    </ion-header>

    <ion-content :fullscreen="true">
      <ion-grid>
        <ion-row v-if="room.roomStatus === RoomVoteStatus.VoteClosed">
          <ion-col>
            <ion-grid>
              <ion-row>
                <ion-col class="ion-col-4">
                  <ion-label v-if="displayVoteResults">Team votes</ion-label>
                  <ion-label v-else>Team</ion-label>
                  <player v-for="participant in participants" :key="participant.id" :player="participant"
                    :displayVote="displayVoteResults" :isCurrentUser="participant.id === localParticipantId" />
                </ion-col>
                <ion-col v-if="displayVoteResults" class="ion-col-6">
                  <BarChart v-if="displayVoteResults" :player-votes="voteResults" :barColors="barColors" />
                </ion-col>
                <ion-col v-if="!displayVoteResults" class="ion-margin-center">
                  No vote to display
                </ion-col>
              </ion-row>
            </ion-grid>
          </ion-col>
        </ion-row>
        <ion-row v-if="room.roomStatus === RoomVoteStatus.VoteOpen">
          <ion-col class="ion-col-2">
            <ion-item v-if="localParticipant">
              <ion-grid>
                <ion-row>Your vote:</ion-row>
                <ion-row><ion-radio-group v-model="localVote" @ionChange="onVoteChange">
                    <ion-item v-for="[command, label] in Object.entries(room.voteCommands)" :key="command">
                      <ion-radio v-if="label !== 'Close vote'" :value="label">{{ label }}</ion-radio>
                    </ion-item>
                  </ion-radio-group></ion-row>
              </ion-grid>
            </ion-item>
          </ion-col>
          <ion-col class="ion-col-4">
            <ion-grid>
              <ion-row>Team votes</ion-row>
              <ion-row>
                <ion-col>
                  <player v-for="participant in participants" :key="participant.id" :player="participant"
                    :displayVote=false :isCurrentUser="participant.id === localParticipantId" />
                </ion-col>
              </ion-row>
            </ion-grid>
          </ion-col>
          <ion-col class="ion-colon-3">
            <ProgressBar :progress="voteProgress"></ProgressBar>
          </ion-col>
        </ion-row>
      </ion-grid>
    </ion-content>
    <IonFooter>
      <ion-grid>
        <ion-row>
          <ion-col v-if="room.roomStatus === RoomVoteStatus.VoteClosed">
            <ion-button v-if="room.voters.length > 1" @click="startGame">
              <ion-icon :icon="playOutline"></ion-icon> Start new vote
            </ion-button>
          </ion-col>
          <ion-col v-if="room.roomStatus === RoomVoteStatus.VoteOpen">
            <ion-button @click="closeVote">
              <ion-icon :icon="playOutline"></ion-icon> Close vote
            </ion-button>
          </ion-col>
          <ion-col>
            <ExitButton></ExitButton>
          </ion-col>
        </ion-row>
      </ion-grid>
    </IonFooter>
  </ion-page>
</template>

<script setup lang="ts">
import { computed, ref, watchEffect, onBeforeMount, onMounted } from 'vue';
import { useRoute } from 'vue-router';
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
import { v4 as uuidv4 } from "uuid";

const route = useRoute();
const urlParamUsername = ref(route.params.username);
const urlParamRoom = ref(route.params.room);
const pageTitle = ref(`Poker room ${urlParamRoom.value}`);

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

// bookmark is used to connect without coming from HomePage
onBeforeMount(async () => {
  if (!store.state.websocket) {
    try {
      store.dispatch('initializeWebSocketConnection');
      await waitForWebsocket();
    } catch (error) {
      console.error('Failed to connect to websocket : ', error);
    }

    async function waitForWebsocket() {
      const maxAttempts = 10;
      let attempts = 0;

      while (store.state.websocket.readyState === WebSocket.CONNECTING && attempts < maxAttempts) {
        await new Promise(resolve => setTimeout(resolve, 1000)); // 1 second pause before each retry
        attempts++;
      }
    }

    store.commit('setRoom', urlParamRoom);

    const newId = uuidv4()
    const message = JSON.stringify({
      id: newId,
      name: urlParamUsername.value,
      vote: "",
      available_commands: {},
      last_command: "",
      room: store.state.roomSelected
    });
    store.state.websocket.send(message);
    store.commit('setLocalParticipantId', newId);
  }
});

onMounted(() => {
  document.title = pageTitle.value;
});

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

const localParticipant = ref<Participant | null>(null);
watchEffect(() => {
  localParticipant.value = store.state.room.voters.find((p: Participant) => p.id === localParticipantId.value) || null;
});

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
  store.dispatch('updateVote', localParticipant.value);
};
</script>
