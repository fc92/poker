<template>
    <ion-page>
        <ion-header :translucent="true">
        </ion-header>

        <ion-content :fullscreen="true">
            <swiper-container class="slider" :slides-per-view="3" :space-between="spaceBetween" :centered-slides="true"
                :pagination="{ hideOnClick: true }" :breakpoints="{ 768: { slidesPerView: 3 } }" @progress="onProgress"
                @slidechange="onSlideChange">
                <swiper-slide class="rounded-slide" v-for="value in fibonacciValues" :key="value.num"
                    :style="{ backgroundColor: value.color }">
                    {{ value.num }}
                </swiper-slide>
            </swiper-container>
        </ion-content>
    </ion-page>
</template>


<script lang="ts">
import { defineComponent, ref } from 'vue';
import { register } from 'swiper/element/bundle';
const fibonacciValues: { num: number, color: string }[] = [
    { num: 1, color: '#f4f4f4' },
    { num: 2, color: '#0f62fe' },
    { num: 3, color: '#198038' },
    { num: 5, color: '#007d79' },
    { num: 8, color: '#1192e8' },
    { num: 13, color: '#8a3ffc' },
    { num: 21, color: '#da1e28' }
];

register();

export default defineComponent({
    setup() {
        const spaceBetween = ref<number>(10);

        const onProgress = (e: any) => {
            const [swiper, progress] = e.detail;
            console.log(progress);
        };

        const onSlideChange = (e: any) => {
            console.log('slide changed');
        }

        return {
            spaceBetween,
            onProgress,
            onSlideChange,
            fibonacciValues
        };
    }

});
</script>

<style scoped>
.slider {
    height: 20vh;
    /* 20% de la hauteur de l'écran */
}

.rounded-slide {
    border-radius: 10px;
    /* Angles arrondis */
    /* Ajoutons un padding pour éloigner le contenu du bord */
    padding: 5px;
    /* Centrons le texte à l'intérieur */
    display: flex;
    align-items: center;
    justify-content: center;
    /* Mettons une taille minimale pour garantir une bonne lisibilité, peut-être ajustée si nécessaire */
    min-width: 80px;
    text-align: center;
}
</style>
