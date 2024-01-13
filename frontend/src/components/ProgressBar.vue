<template>
    <div>
        <Doughnut :data="chartData" :options="chartOptions" />
    </div>
</template>

<script lang="ts">
import { PropType, defineComponent } from 'vue';
import { Chart as ChartJS, ArcElement, Tooltip, Legend, Title } from 'chart.js';
import { Doughnut } from "vue-chartjs"

ChartJS.register(Tooltip, ArcElement, Legend, Title);
interface Dataset {
    data: number[];
    backgroundColor: string[];
    borderWidth: number;
}

interface ChartData {
    labels: string[];
    datasets: Dataset[];
}

export default defineComponent({
    components: {
        Doughnut,
    },
    props: {
        progress: {
            type: Object as PropType<number[]>,
            default: [0, 0],
            required: true,
        },
    },
    computed: {
        chartData(): ChartData {
            return {
                labels: ['Votes Received', 'Remaining'],
                datasets: [
                    {
                        data: this.progress,
                        backgroundColor: ['#36A2EB', '#CCCCCC'],
                        borderWidth: 0,
                    },
                ],
            };
        }
    },
    data() {
        return {
            chartOptions: {
                responsive: true,
                plugins: {
                    title: {
                        display: true,
                        text: 'Vote progress',
                        padding: {
                            top: 10,
                            bottom: 10
                        },
                        font: {
                            size: 20
                        }
                    },
                },
            },
        };
    },
});
</script>

<style scoped>
/* Add component-specific styles here */
</style>
