<template>
    <ion-list>
        <ion-list-header>
            <h2>Choose poker room</h2>
        </ion-list-header>
        <ion-radio-group v-model="selectedServer">
            <ion-list>
                <ion-item v-for="(server, index) in servers" :key="index">
                    <ion-radio :value="server.value">{{ server.label }}</ion-radio><br>
                </ion-item>
            </ion-list>
        </ion-radio-group>
    </ion-list>
    <ion-button @click="selectServer">Select</ion-button>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { IonList, IonListHeader, IonItem, IonRadioGroup, IonRadio, IonButton } from '@ionic/vue';
import { defineProps, defineEmits } from 'vue';

const emit = defineEmits();

const servers = ref([
    { label: 'Localhost', value: 'localhost:8080' },
    { label: 'Azure', value: 'azure:9999' },
]);

const selectedServer = ref(servers.value[0].value);

const selectServer = () => {
    emit('update:serverValue', selectedServer.value);
};

const props = defineProps({
    serverValue: {
        type: String,
        required: false,
        default: 'localhost:8080',
    },
});
</script>
