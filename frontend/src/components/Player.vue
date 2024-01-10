<template>
    <div class="player-item">
        <ion-icon :name="isCurrentUser ? 'person-circle' : 'person'" :class="{ 'bold-text': isCurrentUser }"></ion-icon>
        <span :class="{ 'bold-text': isCurrentUser }">{{ player.name }}</span>
        <!-- <ion-icon v-if="player.last_command === VoteStatus.NotReceived" name="ellipsis-horizontal"></ion-icon>
        <ion-icon v-if="player.available_commands['voted']" name="checkmark-circle-outline"></ion-icon>
        <ion-icon v-if="player.available_commands['notVoting']" name="close-circle"></ion-icon> -->
    </div>
</template>

<script lang="ts">
import { Participant } from '@/participant';
import { VoteStatus } from '@/room';
import { IonIcon } from '@ionic/vue';
import { addIcons } from 'ionicons';
import { personCircle, person, ellipsisHorizontal, checkmarkCircleOutline, closeCircle } from 'ionicons/icons'
import { defineComponent, PropType } from 'vue';



export default defineComponent({
    components: {
        IonIcon,
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
    },
});
</script>

<style scoped>
.player-item {
    display: flex;
    align-items: center;
    padding: 8px;
}

.bold-text {
    font-weight: bold;
}
</style>
