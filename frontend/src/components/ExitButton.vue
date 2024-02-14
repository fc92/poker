<template>
    <div>
        <ion-button @click="handleExitClick">
            <ion-icon :icon="exit"></ion-icon>
        </ion-button>
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

</script>