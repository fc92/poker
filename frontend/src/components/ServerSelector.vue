<template>
    <ion-content>
        <ion-list>
            <ion-list-header>
                Choose poker room
            </ion-list-header>
            <ion-radio-group v-model="selectedServer">
                <ion-list>
                    <ion-item v-for="(server, index) in servers" :key="index">
                        <ion-radio :value="server.value">{{ server.label }}</ion-radio><br>
                    </ion-item>
                </ion-list>
            </ion-radio-group>
            <ion-item>
                <ion-input v-model="customServer" placeholder="or type custom server address here"></ion-input>
            </ion-item>
        </ion-list>
        <ion-button @click="selectServer">Select</ion-button>
    </ion-content>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { IonContent, IonList, IonListHeader, IonItem, IonRadioGroup, IonRadio, IonInput, IonButton } from '@ionic/vue';
import { defineProps, defineEmits } from 'vue';

const emit = defineEmits();

const servers = ref([
    { label: 'Localhost', value: 'localhost:8080' },
    { label: 'Azure', value: 'azure:9999' },
]);

const selectedServer = ref(servers.value[0].value);
const customServer = ref('');

const selectServer = () => {
    const serverToUse = customServer.value || selectedServer.value;
    emit('update:serverValue', serverToUse);
};

const props = defineProps({
    serverValue: {
        type: String,
        required: false,
        default: 'localhost:8080',
    },
});
</script>
