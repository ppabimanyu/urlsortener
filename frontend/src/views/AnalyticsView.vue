<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUrlStore } from '@/stores/url'
import { useClipboard } from '@vueuse/core'

import Navbar from '@/components/layout/Navbar.vue'
import TimelineChart from '@/components/dashboard/TimelineChart.vue'
import BreakdownChart from '@/components/dashboard/BreakdownChart.vue'
import StatCard from '@/components/dashboard/StatCard.vue'

import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import {
  ArrowLeft,
  Download,
  Copy,
  Check,
  ExternalLink,
  MousePointerClick,
  Monitor,
  Smartphone,
  Globe,
  Share2,
  Calendar,
  Lock,
} from 'lucide-vue-next'

const route = useRoute()
const router = useRouter()
const urlStore = useUrlStore()

const urlId = computed(() => route.params.id as string)
const selectedDays = ref(14)
const isExporting = ref(false)

const baseUrl = computed(() => import.meta.env.VITE_BACKEND_URL || 'http://localhost:8080')
const fullShortUrl = computed(() => {
  if (!urlStore.selectedUrlAnalytics?.url) return ''
  return `${baseUrl.value}/${urlStore.selectedUrlAnalytics.url.short_code}`
})

const { copy, copied } = useClipboard({ source: fullShortUrl })

onMounted(() => {
  loadData()
})

watch(selectedDays, () => {
  loadData()
})

async function loadData() {
  if (urlId.value) {
    await urlStore.fetchUrlAnalytics(urlId.value, selectedDays.value)
  }
}

async function handleExport() {
  if (!urlStore.selectedUrlAnalytics?.url) return
  isExporting.value = true
  try {
    await urlStore.exportCsv(urlId.value, urlStore.selectedUrlAnalytics.url.short_code)
  } catch (err) {
    console.error('Failed to export CSV:', err)
  } finally {
    isExporting.value = false
  }
}
</script>

<template>
  <div class="min-h-screen bg-muted/20 pb-16">
    <Navbar />

    <main class="container mx-auto max-w-7xl px-4 sm:px-6 lg:px-8 pt-8 space-y-6">
      <!-- Top Navigation & Actions -->
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div class="flex items-center gap-3">
          <Button variant="outline" size="icon" class="h-9 w-9 rounded-xl" @click="router.back()">
            <ArrowLeft class="h-4 w-4" />
          </Button>
          <div>
            <div class="flex items-center gap-2">
              <h1 class="text-xl font-bold tracking-tight text-foreground font-mono">
                /{{ urlStore.selectedUrlAnalytics?.url?.short_code }}
              </h1>
              <Badge
                v-if="urlStore.selectedUrlAnalytics?.url"
                :variant="urlStore.selectedUrlAnalytics.url.is_active ? 'default' : 'secondary'"
                class="text-[10px] uppercase font-semibold"
              >
                {{ urlStore.selectedUrlAnalytics.url.is_active ? 'Active' : 'Disabled' }}
              </Badge>
            </div>
            <p class="text-xs text-muted-foreground truncate max-w-[300px] sm:max-w-[450px]">
              {{ urlStore.selectedUrlAnalytics?.url?.original_url }}
            </p>
          </div>
        </div>

        <div class="flex items-center gap-2">
          <!-- Days Filter -->
          <div class="flex items-center rounded-lg border border-border/60 bg-card p-1 text-xs">
            <button
              class="px-2.5 py-1 rounded-md transition-colors"
              :class="selectedDays === 7 ? 'bg-primary text-primary-foreground font-semibold' : 'text-muted-foreground hover:text-foreground'"
              @click="selectedDays = 7"
            >
              7 Days
            </button>
            <button
              class="px-2.5 py-1 rounded-md transition-colors"
              :class="selectedDays === 14 ? 'bg-primary text-primary-foreground font-semibold' : 'text-muted-foreground hover:text-foreground'"
              @click="selectedDays = 14"
            >
              14 Days
            </button>
            <button
              class="px-2.5 py-1 rounded-md transition-colors"
              :class="selectedDays === 30 ? 'bg-primary text-primary-foreground font-semibold' : 'text-muted-foreground hover:text-foreground'"
              @click="selectedDays = 30"
            >
              30 Days
            </button>
          </div>

          <!-- Export CSV Button -->
          <Button variant="outline" class="gap-1.5 text-xs font-medium" :disabled="isExporting" @click="handleExport">
            <Download class="h-3.5 w-3.5" />
            <span>{{ isExporting ? 'Exporting...' : 'Export CSV' }}</span>
          </Button>
        </div>
      </div>

      <!-- Quick Info Bar Card -->
      <Card class="border border-border/60 shadow-xs bg-card">
        <CardContent class="p-4 sm:p-5 flex flex-col md:flex-row items-start md:items-center justify-between gap-4">
          <div class="flex flex-col space-y-1">
            <span class="text-xs font-semibold text-muted-foreground uppercase tracking-wider">Short URL</span>
            <div class="flex items-center gap-2">
              <span class="text-sm font-mono font-bold text-foreground">{{ fullShortUrl }}</span>
              <Button variant="ghost" size="icon" class="h-7 w-7" @click="copy(fullShortUrl)">
                <Check v-if="copied" class="h-3.5 w-3.5 text-emerald-500" />
                <Copy v-else class="h-3.5 w-3.5" />
              </Button>
              <a :href="fullShortUrl" target="_blank" class="text-muted-foreground hover:text-foreground">
                <ExternalLink class="h-3.5 w-3.5" />
              </a>
            </div>
          </div>

          <div class="flex items-center gap-6 border-t md:border-t-0 pt-3 md:pt-0 w-full md:w-auto border-border/40 text-xs">
            <div>
              <span class="text-muted-foreground block text-[11px]">Total Clicks</span>
              <span class="text-base font-bold text-foreground">
                {{ urlStore.selectedUrlAnalytics?.url?.click_count ?? 0 }}
              </span>
            </div>
            <div>
              <span class="text-muted-foreground block text-[11px]">Password</span>
              <span class="font-medium text-foreground">
                {{ urlStore.selectedUrlAnalytics?.url?.has_password ? 'Protected' : 'None' }}
              </span>
            </div>
            <div>
              <span class="text-muted-foreground block text-[11px]">Click Limit</span>
              <span class="font-medium text-foreground">
                {{ urlStore.selectedUrlAnalytics?.url?.click_limit ? `${urlStore.selectedUrlAnalytics.url.click_limit} clicks` : 'Unlimited' }}
              </span>
            </div>
          </div>
        </CardContent>
      </Card>

      <!-- Main Click Trend Timeline -->
      <Card class="border border-border/60 shadow-xs">
        <CardHeader class="pb-2">
          <CardTitle class="text-base font-bold">Clicks Over Time</CardTitle>
          <CardDescription class="text-xs">Daily visitor engagement over the selected period</CardDescription>
        </CardHeader>
        <CardContent class="pt-2">
          <TimelineChart :timeline="urlStore.selectedUrlAnalytics?.timeline || []" />
        </CardContent>
      </Card>

      <!-- Breakdown Grid -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <!-- Devices Breakdown -->
        <Card class="border border-border/60 shadow-xs">
          <CardHeader class="pb-2">
            <CardTitle class="text-sm font-bold flex items-center gap-2">
              <Smartphone class="h-4 w-4 text-sky-500" />
              <span>Devices</span>
            </CardTitle>
            <CardDescription class="text-xs">Mobile vs Desktop vs Tablet</CardDescription>
          </CardHeader>
          <CardContent>
            <BreakdownChart
              type="doughnut"
              :data="urlStore.selectedUrlAnalytics?.devices || []"
            />
          </CardContent>
        </Card>

        <!-- OS Breakdown -->
        <Card class="border border-border/60 shadow-xs">
          <CardHeader class="pb-2">
            <CardTitle class="text-sm font-bold flex items-center gap-2">
              <Monitor class="h-4 w-4 text-violet-500" />
              <span>Operating Systems</span>
            </CardTitle>
            <CardDescription class="text-xs">Visitor platforms</CardDescription>
          </CardHeader>
          <CardContent>
            <BreakdownChart
              type="doughnut"
              :data="urlStore.selectedUrlAnalytics?.os || []"
            />
          </CardContent>
        </Card>

        <!-- Browsers Breakdown -->
        <Card class="border border-border/60 shadow-xs">
          <CardHeader class="pb-2">
            <CardTitle class="text-sm font-bold flex items-center gap-2">
              <Globe class="h-4 w-4 text-emerald-500" />
              <span>Browsers</span>
            </CardTitle>
            <CardDescription class="text-xs">Top browsers used</CardDescription>
          </CardHeader>
          <CardContent>
            <BreakdownChart
              type="bar"
              :data="urlStore.selectedUrlAnalytics?.browsers || []"
            />
          </CardContent>
        </Card>

        <!-- Referrers Breakdown -->
        <Card class="border border-border/60 shadow-xs">
          <CardHeader class="pb-2">
            <CardTitle class="text-sm font-bold flex items-center gap-2">
              <Share2 class="h-4 w-4 text-amber-500" />
              <span>Traffic Referrers</span>
            </CardTitle>
            <CardDescription class="text-xs">Sources where visitors came from</CardDescription>
          </CardHeader>
          <CardContent>
            <BreakdownChart
              type="bar"
              :data="urlStore.selectedUrlAnalytics?.referrers || []"
            />
          </CardContent>
        </Card>

        <!-- Top Countries -->
        <Card class="border border-border/60 shadow-xs lg:col-span-2">
          <CardHeader class="pb-2">
            <CardTitle class="text-sm font-bold flex items-center gap-2">
              <Globe class="h-4 w-4 text-pink-500" />
              <span>Geographic Distribution</span>
            </CardTitle>
            <CardDescription class="text-xs">Top visitor countries & cities</CardDescription>
          </CardHeader>
          <CardContent>
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div>
                <h4 class="text-xs font-semibold text-muted-foreground uppercase mb-2">Countries</h4>
                <div class="space-y-2 max-h-[180px] overflow-y-auto">
                  <div
                    v-for="c in urlStore.selectedUrlAnalytics?.countries"
                    :key="c.name"
                    class="flex items-center justify-between text-xs p-2 rounded-lg bg-muted/40"
                  >
                    <span class="font-medium truncate">{{ c.name }}</span>
                    <span class="font-bold text-foreground">{{ c.count }}</span>
                  </div>
                  <div
                    v-if="!urlStore.selectedUrlAnalytics?.countries?.length"
                    class="text-xs text-muted-foreground py-4 text-center"
                  >
                    No country data
                  </div>
                </div>
              </div>

              <div>
                <h4 class="text-xs font-semibold text-muted-foreground uppercase mb-2">Cities</h4>
                <div class="space-y-2 max-h-[180px] overflow-y-auto">
                  <div
                    v-for="city in urlStore.selectedUrlAnalytics?.cities"
                    :key="city.name"
                    class="flex items-center justify-between text-xs p-2 rounded-lg bg-muted/40"
                  >
                    <span class="font-medium truncate">{{ city.name }}</span>
                    <span class="font-bold text-foreground">{{ city.count }}</span>
                  </div>
                  <div
                    v-if="!urlStore.selectedUrlAnalytics?.cities?.length"
                    class="text-xs text-muted-foreground py-4 text-center"
                  >
                    No city data
                  </div>
                </div>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </main>
  </div>
</template>
