<template>
    <Bar id="my-chart-id" :options="chartOptions" :data="chartData" />
</template>
  
<script lang="ts">
import { Bar } from 'vue-chartjs'
import { Chart as ChartJS, Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale, ChartData } from 'chart.js'

ChartJS.register(Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale)
import { defineComponent } from 'vue'
import type { PropType } from 'vue'

export default defineComponent({


    name: 'BarChart',

    components: { Bar },

    props: {
        playerVotes: {
            type: Object as PropType<number[]>,
            required: true
        }

    },
    data() {
        return {
            chartData: {
                labels: ['1', '2', '3', '5', '8', '13', '21', '?'],
                datasets: [{ data: this.playerVotes }]
            },
            chartOptions: {
                responsive: true
            }
        }
    }
});
</script>
  