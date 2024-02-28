<template>
    <div>
        <ion-button @click="handleExitClick">
            <ion-icon :icon="exit"></ion-icon>
        </ion-button>
        <ion-toolbar>
            <ion-buttons slot="end">
                <!-- Add a button to open the GitHub URL in a new tab or window -->
                <ion-button @click="openGitHubInNewTab" size="small">
                    Project on GitHub
                </ion-button>
            </ion-buttons>
        </ion-toolbar>
    </div>
</template>
      
<script setup lang="ts">
import { IonButton, IonIcon } from '@ionic/vue';
import { exit } from 'ionicons/icons';
import { useStore } from 'vuex';
import { useRouter } from 'vue-router';

const router = useRouter();

const store = useStore();
const handleExitClick = () => {
    store.dispatch('handleExitClick');
    router.push('/home');
};

const openGitHubInNewTab = () => {
    window.open('https://github.com/fc92/poker', '_blank');
};


</script>