<template>
    <ion-page>
        <ion-header :translucent="true">
        </ion-header>

        <ion-content :fullscreen="true">
            <swiper-container :slides-per-view="3" :space-between="spaceBetween" :centered-slides="true"
                :pagination="{ hideOnClick: true }" :breakpoints="{ 768: { slidesPerView: 3 } }" @progress="onProgress"
                @slidechange="onSlideChange">
                <swiper-slide v-for="value in fibonacciValues" :key="value.num" :style="{ backgroundColor: value.color }">
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
