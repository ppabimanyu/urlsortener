<script setup lang="ts">
import { computed } from 'vue'
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler,
  type ChartOptions,
} from 'chart.js'
import { Line } from 'vue-chartjs'
import { useDark } from '@vueuse/core'
import type { DateClickStat } from '@/types'

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler
)

const props = defineProps<{
  timeline: DateClickStat[]
  title?: string
}>()

const isDark = useDark()

const chartData = computed(() => {
  const labels = props.timeline.map((item) => item.date)
  const data = props.timeline.map((item) => item.count)

  return {
    labels,
    datasets: [
      {
        label: 'Total Clicks',
        data,
        borderColor: isDark.value ? '#38bdf8' : '#0284c7',
        backgroundColor: isDark.value ? 'rgba(56, 189, 248, 0.15)' : 'rgba(2, 132, 199, 0.1)',
        fill: true,
        tension: 0.35,
        pointBackgroundColor: isDark.value ? '#38bdf8' : '#0284c7',
        pointBorderColor: isDark.value ? '#0f172a' : '#ffffff',
        pointHoverRadius: 6,
        pointRadius: 4,
        borderWidth: 2.5,
      },
    ],
  }
})

const chartOptions = computed<ChartOptions<'line'>>(() => {
  const textColor = isDark.value ? '#94a3b8' : '#64748b'
  const gridColor = isDark.value ? 'rgba(255, 255, 255, 0.08)' : 'rgba(0, 0, 0, 0.06)'

  return {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        display: false,
      },
      tooltip: {
        backgroundColor: isDark.value ? '#1e293b' : '#0f172a',
        titleColor: '#f8fafc',
        bodyColor: '#f8fafc',
        borderColor: isDark.value ? 'rgba(255, 255, 255, 0.1)' : 'rgba(0, 0, 0, 0.1)',
        borderWidth: 1,
        padding: 10,
        displayColors: false,
        callbacks: {
          label: (context) => ` ${context.parsed.y} Clicks`,
        },
      },
    },
    scales: {
      x: {
        grid: {
          color: gridColor,
          display: true,
        },
        ticks: {
          color: textColor,
          font: { family: 'inherit', size: 12 },
          maxRotation: 45,
          minRotation: 0,
        },
      },
      y: {
        beginAtZero: true,
        grid: {
          color: gridColor,
        },
        ticks: {
          color: textColor,
          font: { family: 'inherit', size: 12 },
          stepSize: 1,
          precision: 0,
        },
      },
    },
  }
})
</script>

<template>
  <div class="h-[280px] w-full">
    <div v-if="!timeline || timeline.length === 0" class="flex h-full items-center justify-center text-sm text-muted-foreground">
      No click data recorded yet
    </div>
    <Line v-else :data="chartData" :options="chartOptions" />
  </div>
</template>
