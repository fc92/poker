<template>
    <div class="player-item">
        <ion-icon :name="isCurrentUser ? 'person-circle' : 'person'" :class="{ 'bold-text': isCurrentUser }"></ion-icon>
        <span :class="{ 'bold-text': isCurrentUser }">{{ player.name }}</span>
        <ion-label v-if="!displayVote && player.last_command === ''" color="warning"><span class="vote-status">: waiting for
                vote</span></ion-label>
        <ion-icon v-if="!displayVote && player.last_command === ''" name="ellipsis-horizontal"
            class="player-item"></ion-icon>
        <ion-label v-if="!displayVote && player.last_command === 'r'" color="success"><span class="vote-status">: vote
                received</span></ion-label>
        <ion-icon v-if="!displayVote && player.last_command === 'r'" name="checkmark-circle-outline"
            class="player-item"></ion-icon>
        <ion-label v-if="!displayVote && player.last_command === 'n'" color="success"><span class="vote-status">: vote
                received</span></ion-label>
        <ion-icon v-if="!displayVote && player.last_command === 'n'" name="close-circle" class="player-item"></ion-icon>
        <div v-if="displayVote && player.vote" class="player-item">{{ player.vote }}</div>
    </div>
</template>
 
<script lang="ts">
import { Participant } from '@/participant';
import { IonIcon, IonLabel } from '@ionic/vue';
import { addIcons } from 'ionicons';
import { personCircle, person, ellipsisHorizontal, checkmarkCircleOutline, closeCircle } from 'ionicons/icons'
import { defineComponent, PropType } from 'vue';



export default defineComponent({
    components: {
        IonIcon, IonLabel
    },
    setup() {
        addIcons(
            { personCircle, person, ellipsisHorizontal, checkmarkCircleOutline, closeCircle });
    },
    props: {
        player: {
            type: Object as PropType<Participant>,
            required: true,
        },
        isCurrentUser: {
            type: Boolean,
            default: false,
        },
        displayVote: {
            type: Boolean,
            default: false
        }
    },
});
</script>

<style scoped>
.player-item {
    display: flex;
    align-items: center;
    padding: 4px;
}

.bold-text {
    font-weight: bold;
}

.vote-status {
    margin: 4px;
}
</style>
