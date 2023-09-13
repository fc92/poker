<template>
  <ion-page>
    <ion-header :translucent="true">
    </ion-header>

    <ion-content :fullscreen="true">

      <server-selector @update:serverValue="handleServerValueUpdate"></server-selector>
      <name-selector @update:player="handlePlayerUpdate"></name-selector>
      <div class="enter-name">
        <img class="aditya-chinchure-hy-n-9-a-u-9-tm-c-unsplash-4" src="src/assets/images/background.png" />
        <div class="header">
          <div class="frame-2">
            <div class="title-background"></div>
          </div>
        </div>
        <div class="footer">
          <div class="exit">
            <div class="button"></div>
            <svg class="icon-exit-outline" width="49" height="38" viewBox="0 0 49 38" fill="none"
              xmlns="http://www.w3.org/2000/svg">
              <path
                d="M32.0385 9.42308V4.71154C32.0385 3.46196 31.5421 2.26356 30.6585 1.37998C29.7749 0.496393 28.5765 0 27.3269 0H4.71154C3.46196 0 2.26356 0.496393 1.37998 1.37998C0.496393 2.26356 0 3.46196 0 4.71154V32.9808C0 34.2303 0.496393 35.4287 1.37998 36.3123C2.26356 37.1959 3.46196 37.6923 4.71154 37.6923H27.3269C28.5765 37.6923 29.7749 37.1959 30.6585 36.3123C31.5421 35.4287 32.0385 34.2303 32.0385 32.9808V28.2692"
                stroke="black" stroke-width="7.58475" stroke-linecap="round" stroke-linejoin="round" />
              <path d="M39.5769 9.42307L49 18.8461L39.5769 28.2692" stroke="black" stroke-width="7.58475"
                stroke-linecap="round" stroke-linejoin="round" />
              <path d="M16.8437 18.8462H49" stroke="black" stroke-width="7.58475" stroke-linecap="round"
                stroke-linejoin="round" />
            </svg>
          </div>
        </div>
        <div class="body">
          <div class="name-background"></div>
          <div class="local-player">
            <div class="bandeau-joueur"></div>
            <svg class="icon-person-circle-outline" width="28" height="27" viewBox="0 0 28 27" fill="none"
              xmlns="http://www.w3.org/2000/svg">
              <path
                d="M13.7363 0.00121403C6.1179 -0.0980339 -0.10164 5.90081 0.0012587 13.2489C0.102855 20.2962 6.05212 26.0343 13.3586 26.1323C20.9784 26.2328 27.1966 20.234 27.0924 12.8859C26.9921 5.83737 21.0428 0.0992056 13.7363 0.00121403ZM21.9696 20.5575C21.9436 20.5845 21.9118 20.6057 21.8764 20.6196C21.841 20.6335 21.8029 20.6397 21.7648 20.6378C21.7266 20.6359 21.6893 20.6259 21.6556 20.6085C21.6219 20.5912 21.5926 20.5669 21.5697 20.5374C20.9872 19.8023 20.2738 19.1729 19.4629 18.6787C17.8048 17.6523 15.7038 17.0869 13.5475 17.0869C11.3912 17.0869 9.29019 17.6523 7.63208 18.6787C6.82118 19.1727 6.1078 19.8018 5.52525 20.5367C5.50237 20.5662 5.47303 20.5905 5.43932 20.6079C5.40561 20.6253 5.36836 20.6352 5.3302 20.6372C5.29204 20.6391 5.25391 20.6329 5.21851 20.619C5.18312 20.6051 5.15132 20.5839 5.12538 20.5568C3.21437 18.5671 2.13185 15.9671 2.08529 13.2552C1.97914 7.14267 7.17945 2.02637 13.5195 2.0113C19.8595 1.99622 25.0097 6.96176 25.0097 13.0668C25.0118 15.8441 23.926 18.5196 21.9696 20.5575Z"
                fill="black" />
              <path
                d="M13.5475 6.03147C12.2632 6.03147 11.102 6.49567 10.2768 7.33928C9.4517 8.18288 9.03945 9.34936 9.13258 10.6013C9.32145 13.0668 11.3019 15.0768 13.5475 15.0768C15.793 15.0768 17.7696 13.0668 17.9624 10.6019C18.0588 9.36192 17.6498 8.20613 16.8109 7.34682C15.9825 6.49881 14.8233 6.03147 13.5475 6.03147Z"
                fill="black" />
            </svg>
          </div>
          <div class="enter-your-name">Enter your name:</div>
        </div>
      </div>


    </ion-content>
  </ion-page>
</template>

<script setup lang="ts">
import { ref, onBeforeUnmount } from 'vue';
import { IonContent, IonHeader, IonPage } from '@ionic/vue';
import ServerSelector from '@/components/ServerSelector.vue';
import NameSelector from '@/components/NameSelector.vue';
import { Player } from '@/player'

const websocket = ref<WebSocket | null>(null);

const connectToWebSocket = (serverAddress: string) => {
  websocket.value = new WebSocket(`ws://${serverAddress}/ws`);

  websocket.value.addEventListener('open', () => {
    console.log('WebSocket is connected');
  });

  websocket.value.addEventListener('message', (event) => {
    console.log('Message from server: ', event.data);
  });

  websocket.value.addEventListener('close', () => {
    console.log('WebSocket is closed');
  });

  websocket.value.addEventListener('error', (error) => {
    console.error('WebSocket error:', error);
  });
};


const handleServerValueUpdate = (newServerValue: string) => {
  console.log('Valeur du serveur mise à jour:', newServerValue);
  connectToWebSocket(newServerValue);
};
const handlePlayerUpdate = (player: Player) => {
  console.log('Valeur du player id mise à jour:', player.id);
  console.log('Valeur du player name mise à jour:', player.name);
  if (websocket.value) {
    const message = JSON.stringify({
      id: player.id,
      name: player.name,
      vote: "",
      available_commands: {},
      last_command: ""
    });
    websocket.value.send(message);
  } else {
    console.error('WebSocket is not connected');
  }
};

onBeforeUnmount(() => {
  if (websocket.value) {
    websocket.value.close();
  }
});
</script>


<style scoped>
.enter-name,
.enter-name * {
  box-sizing: border-box;
}

.enter-name {
  background: #ffffff;
  display: flex;
  flex-direction: row;
  gap: 10px;
  align-items: flex-start;
  justify-content: flex-start;
  width: 360px;
  height: 640px;
  position: relative;
  overflow: hidden;
}

.aditya-chinchure-hy-n-9-a-u-9-tm-c-unsplash-4 {
  flex-shrink: 0;
  position: relative;
}

.header {
  padding: 34px 0px 34px 0px;
  display: flex;
  flex-direction: row;
  gap: 10px;
  align-items: center;
  justify-content: center;
  flex-wrap: wrap;
  flex-shrink: 0;
  width: 286px;
  height: 120px;
  position: absolute;
  left: 2741px;
  top: 0px;
  box-shadow: 0px 4px 4px 0px rgba(0, 0, 0, 0.25);
  overflow: hidden;
}

.frame-2 {
  display: flex;
  flex-direction: column;
  gap: 10px;
  align-items: center;
  justify-content: center;
  flex: 1;
  position: relative;
}

.title-background {
  background: rgba(127, 140, 167, 0.74);
  align-self: stretch;
  flex-shrink: 0;
  height: 51px;
  position: relative;
  box-shadow: 0px 4px 4px 0px rgba(0, 0, 0, 0.25);
}

.footer {
  padding: 19px 127px 19px 127px;
  display: flex;
  flex-direction: column;
  gap: 10px;
  align-items: flex-start;
  justify-content: flex-start;
  flex-shrink: 0;
  position: absolute;
  left: 0px;
  top: 542px;
  overflow: hidden;
}

.exit {
  flex-shrink: 0;
  width: 106px;
  height: 59px;
  position: relative;
}

.button {
  background: #d9d9d9;
  border-radius: 26px;
  width: 106px;
  height: 59px;
  position: absolute;
  left: 0px;
  top: 0px;
}

.icon-exit-outline {
  position: absolute;
  left: 32px;
  top: 11px;
  overflow: visible;
}

.body {
  padding: 149px 0px 149px 0px;
  display: flex;
  flex-direction: column;
  gap: 10px;
  align-items: flex-start;
  justify-content: flex-start;
  flex-shrink: 0;
  position: absolute;
  left: 0px;
  top: 90px;
  overflow: hidden;
}

.name-background {
  background: rgba(217, 217, 217, 0.39);
  flex-shrink: 0;
  width: 360px;
  height: 135px;
  position: relative;
}

.local-player {
  display: flex;
  flex-direction: column;
  gap: 10px;
  align-items: flex-start;
  justify-content: flex-start;
  flex-shrink: 0;
  position: absolute;
  left: 30px;
  top: 211px;
}

.bandeau-joueur {
  background: #d9d9d9;
  border-radius: 15px;
  flex-shrink: 0;
  width: 300px;
  height: 30px;
  position: relative;
}

.icon-person-circle-outline {
  border-radius: 15px;
  display: flex;
  flex-direction: column;
  gap: 10px;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  position: absolute;
  left: 6px;
  top: 2px;
  overflow: visible;
}

.enter-your-name {
  color: #ffffff;
  text-align: left;
  font: 700 24px "Inter", sans-serif;
  position: absolute;
  left: 53px;
  top: 164px;
}
</style>