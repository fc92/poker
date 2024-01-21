<template>
    <div class="chart-container">
        <Bar id="my-chart-id" :options="chartOptions" :data="chartData" />
    </div>
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
        },
        barColors: {
            type: Array as PropType<string[]>,
            required: true
        }
    },
    data() {
        return {
            chartData: {
                labels: ['1', '2', '3', '5', '8', '13', '21', '?'],
                datasets: [{ data: this.playerVotes, backgroundColor: this.barColors }]
            },
            chartOptions: {
                responsive: true,
                maintainAspectRatio: true,
                plugins: {
                    title: {
                        display: true,
                        text: 'Distribution of votes',
                        padding: {
                            top: 10,
                            bottom: 30
                        },
                        font: {
                            size: 20
                        }
                    },
                    legend: {
                        display: false,
                    },
                },
                scales: {
                    x: {
                        grid: {
                            color: 'rgba(255, 255, 255, 0.1)', // Optional: Add grid lines for better visibility
                        },
                        ticks: {
                            font: {
                                size: 20
                            }
                        },
                    },
                    y: {
                        beginAtZero: true, // Ensure the scale starts at zero
                        grid: {
                            color: 'rgba(255, 255, 255, 0.1)', // Optional: Add grid lines for better visibility
                        },
                        ticks: {
                            stepSize: 1, // Display only integer values
                            font: {
                                size: 20
                            }
                        },
                    },
                },
            }
        }
    }
});
</script>
<style scoped>
.chart-container {
    background-color: black;
    padding: 20px;
    /* Optional: Add padding for better visual appearance */
}
</style>
  