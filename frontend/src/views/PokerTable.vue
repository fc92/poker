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
            <ion-button @click="startGame">
              <ion-icon :icon="playOutline"></ion-icon> Start Game
            </ion-button>
          </div>

          <div v-else-if="room.roomStatus === RoomVoteStatus.VoteOpen">
            <player v-for="participant in participants" :key="participant.id" :player="participant"
              :isCurrentUser="participant.id === localParticipantId" />
            <ion-item v-if="localParticipant">
              <ion-label>Vote</ion-label>
              <ion-select v-model="localVote" @ionChange="onVoteChange">
                <ion-select-option v-for="option in voteOptions" :key="option.value" :value="option.value">
                  {{ option.label }}
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
import ExitButton from '@/components/ExitButton.vue';
import { Participant } from '@/participant';
import { RoomVoteStatus, voteOptions } from '@/room';

const store = useStore();

const room = computed(() => store.state.room);
const participants = computed(() => store.state.room.voters);
const localParticipantId = computed(() => store.state.localParticipantId);

const localParticipant = computed(() => participants.value.find((p: Participant) => p.id === localParticipantId.value));

const localVote = computed({
  get: () => localParticipant.value?.vote || '',
  set: (value) => {
    store.commit('updateLocalVote', { id: localParticipantId.value, vote: value });
  }
});

const startGame = () => {
  // Ajoutez la logique pour démarrer le jeu
  store.dispatch('startGame');
};

const onVoteChange = () => {
  store.dispatch('updateVote', { id: localParticipantId.value, vote: localVote.value });
};
</script>
