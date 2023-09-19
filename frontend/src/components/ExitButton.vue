<template>
    <div v-if="store.state.serverSelected !== ''">
        <ion-button @click="handleExitClick">
            <ion-icon :icon="exit"></ion-icon>
        </ion-button>
    </div>
</template>
      
<script setup lang="ts">
import { IonButton, IonIcon } from '@ionic/vue';
import { exit } from 'ionicons/icons';
import { useStore } from 'vuex';

const store = useStore();
const handleExitClick = () => {
    store.dispatch('handleExitClick');
};

</script>