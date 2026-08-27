<script setup lang="ts">
import { computed } from 'vue'
import {
  Chart as ChartJS,
  ArcElement,
  BarElement,
  CategoryScale,
  LinearScale,
  Tooltip,
  Legend,
  type ChartOptions,
} from 'chart.js'
import { Doughnut, Bar } from 'vue-chartjs'
import { useDark } from '@vueuse/core'
import type { GroupCountStat } from '@/types'

ChartJS.register(ArcElement, BarElement, CategoryScale, LinearScale, Tooltip, Legend)

const props = withDefaults(
  defineProps<{
    data: GroupCountStat[]
    type?: 'doughnut' | 'bar'
    title?: string
  }>(),
  {
    type: 'doughnut',
  }
)

const isDark = useDark()

const modernPalette = [
  '#0284c7', // Sky blue
  '#8b5cf6', // Violet
  '#10b981', // Emerald
  '#f59e0b', // Amber
  '#ec4899', // Pink
  '#06b6d4', // Cyan
  '#64748b', // Slate
]

const chartData = computed(() => {
  const labels = props.data.map((d) => d.name)
  const counts = props.data.map((d) => d.count)

  if (props.type === 'doughnut') {
    return {
      labels,
      datasets: [
        {
          data: counts,
          backgroundColor: modernPalette.slice(0, Math.max(labels.length, 1)),
          borderColor: isDark.value ? '#1e293b' : '#ffffff',
          borderWidth: 2,
          hoverOffset: 4,
        },
      ],
    }
  }

  // Horizontal bar
  return {
    labels,
    datasets: [
      {
        data: counts,
        backgroundColor: isDark.value ? '#38bdf8' : '#0284c7',
        borderRadius: 6,
      },
    ],
  }
})

const doughnutOptions = computed<ChartOptions<'doughnut'>>(() => ({
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      position: 'bottom',
      labels: {
        color: isDark.value ? '#cbd5e1' : '#475569',
        padding: 14,
        font: { size: 12 },
      },
    },
    tooltip: {
      backgroundColor: isDark.value ? '#1e293b' : '#0f172a',
      titleColor: '#f8fafc',
      bodyColor: '#f8fafc',
      padding: 10,
    },
  },
  cutout: '65%',
}))

const barOptions = computed<ChartOptions<'bar'>>(() => {
  const textColor = isDark.value ? '#94a3b8' : '#64748b'
  const gridColor = isDark.value ? 'rgba(255, 255, 255, 0.08)' : 'rgba(0, 0, 0, 0.06)'

  return {
    responsive: true,
    maintainAspectRatio: false,
    indexAxis: 'y' as const,
    plugins: {
      legend: { display: false },
      tooltip: {
        backgroundColor: isDark.value ? '#1e293b' : '#0f172a',
        titleColor: '#f8fafc',
        bodyColor: '#f8fafc',
        padding: 10,
      },
    },
    scales: {
      x: {
        beginAtZero: true,
        grid: { color: gridColor },
        ticks: { color: textColor, precision: 0 },
      },
      y: {
        grid: { display: false },
        ticks: { color: textColor },
      },
    },
  }
})
</script>

<template>
  <div class="h-[240px] w-full flex items-center justify-center">
    <div v-if="!data || data.length === 0" class="text-sm text-muted-foreground">
      No data recorded yet
    </div>
    <Doughnut v-else-if="type === 'doughnut'" :data="chartData" :options="doughnutOptions" />
    <Bar v-else :data="chartData" :options="barOptions" />
  </div>
</template>
