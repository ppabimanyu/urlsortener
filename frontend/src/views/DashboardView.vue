<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUrlStore } from '@/stores/url'
import { useAuthStore } from '@/stores/auth'
import { useClipboard } from '@vueuse/core'
import confetti from 'canvas-confetti'
import type { URLItem } from '@/types'

import Navbar from '@/components/layout/Navbar.vue'
import StatCard from '@/components/dashboard/StatCard.vue'
import TimelineChart from '@/components/dashboard/TimelineChart.vue'
import CreateEditUrlDialog from '@/components/dialogs/CreateEditUrlDialog.vue'
import QrCodeModal from '@/components/dialogs/QrCodeModal.vue'

import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import {
  Link2,
  MousePointerClick,
  Activity,
  ArrowUpRight,
  Plus,
  Copy,
  Check,
  QrCode,
  BarChart3,
  ExternalLink,
  Lock,
  Globe,
  Sparkles,
} from 'lucide-vue-next'

const router = useRouter()
const urlStore = useUrlStore()
const authStore = useAuthStore()

const quickUrl = ref('')
const quickLoading = ref(false)
const quickError = ref<string | null>(null)

const isCreateModalOpen = ref(false)
const isQrModalOpen = ref(false)
const selectedUrlForQr = ref<URLItem | null>(null)

const copiedId = ref<string | null>(null)

onMounted(async () => {
  await Promise.all([
    urlStore.fetchOverview(),
    urlStore.fetchUrls(1, '', null),
  ])
})

const baseUrl = computed(() => import.meta.env.VITE_BACKEND_URL || 'http://localhost:8080')

async function handleQuickShorten() {
  if (!quickUrl.value.trim()) return
  quickLoading.value = true
  quickError.value = null
  try {
    const created = await urlStore.createUrl({ original_url: quickUrl.value.trim() })
    quickUrl.value = ''
    confetti({
      particleCount: 60,
      spread: 55,
      origin: { y: 0.7 },
    })
    await urlStore.fetchOverview()
  } catch (err: any) {
    quickError.value = err.response?.data?.error || 'Failed to shorten URL'
  } finally {
    quickLoading.value = false
  }
}

function copyShortLink(urlItem: URLItem) {
  const fullUrl = `${baseUrl.value}/${urlItem.short_code}`
  navigator.clipboard.writeText(fullUrl)
  copiedId.value = urlItem.id
  setTimeout(() => {
    if (copiedId.value === urlItem.id) {
      copiedId.value = null
    }
  }, 2000)
}

function openQr(urlItem: URLItem) {
  selectedUrlForQr.value = urlItem
  isQrModalOpen.value = true
}

function formatRelativeTime(dateStr: string) {
  const date = new Date(dateStr)
  return date.toLocaleDateString('id-ID', {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}
</script>

<template>
  <div class="min-h-screen bg-muted/20 pb-16">
    <Navbar />

    <main class="container mx-auto max-w-7xl px-4 sm:px-6 lg:px-8 pt-8 space-y-8">
      <!-- Welcome & Quick Shorten Banner -->
      <div class="relative overflow-hidden rounded-2xl bg-gradient-to-br from-primary/10 via-background to-secondary/30 border border-border/60 p-6 md:p-8 shadow-xs">
        <div class="max-w-2xl space-y-2">
          <div class="inline-flex items-center gap-1.5 rounded-full bg-primary/10 px-3 py-1 text-xs font-semibold text-primary">
            <Sparkles class="h-3.5 w-3.5" />
            <span>Welcome back, {{ authStore.user?.name }}</span>
          </div>
          <h1 class="text-2xl sm:text-3xl font-extrabold tracking-tight text-foreground">
            Shorten URLs, track every click in realtime
          </h1>
          <p class="text-sm text-muted-foreground">
            Paste long links below for instant shortening, or click Create for custom slug and password protection.
          </p>
        </div>

        <!-- Quick Input Form -->
        <form @submit.prevent="handleQuickShorten" class="mt-6 flex flex-col sm:flex-row gap-2.5 max-w-2xl">
          <div class="relative flex-1">
            <Link2 class="absolute left-3.5 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
            <Input
              v-model="quickUrl"
              type="url"
              placeholder="https://paste-your-very-long-url-here.com/..."
              required
              class="h-11 pl-10 bg-background/80 text-sm shadow-xs"
            />
          </div>
          <Button type="submit" size="lg" class="h-11 px-6 gap-2 font-medium shadow-sm shrink-0" :disabled="quickLoading">
            <Plus class="h-4 w-4" />
            <span>{{ quickLoading ? 'Shortening...' : 'Shorten URL' }}</span>
          </Button>
        </form>
        <p v-if="quickError" class="mt-2 text-xs text-destructive font-medium">{{ quickError }}</p>
      </div>

      <!-- KPI Stat Cards -->
      <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
        <StatCard
          title="Total Links Created"
          :value="urlStore.overview?.total_links ?? 0"
          description="Total active and inactive short links"
          :icon="Link2"
          iconColor="text-sky-500 bg-sky-500/10"
        />
        <StatCard
          title="Total Clicks Tracked"
          :value="urlStore.overview?.total_clicks ?? 0"
          description="All visitor redirections recorded"
          :icon="MousePointerClick"
          iconColor="text-violet-500 bg-violet-500/10"
        />
        <StatCard
          title="Active Links"
          :value="urlStore.overview?.active_links ?? 0"
          description="Links currently routing visitors"
          :icon="Activity"
          iconColor="text-emerald-500 bg-emerald-500/10"
        />
      </div>

      <!-- Chart & Analytics Overview -->
      <div class="grid grid-cols-1 gap-6 lg:grid-cols-3">
        <!-- Main Click Activity Chart -->
        <Card class="lg:col-span-2 border border-border/60 shadow-xs">
          <CardHeader class="flex flex-row items-center justify-between pb-2">
            <div>
              <CardTitle class="text-base font-bold">Clicks Over Time</CardTitle>
              <CardDescription class="text-xs">Visitor click activity in the past 14 days</CardDescription>
            </div>
            <Badge variant="outline" class="text-xs font-normal">Last 14 Days</Badge>
          </CardHeader>
          <CardContent class="pt-2">
            <TimelineChart :timeline="urlStore.overview?.timeline || []" />
          </CardContent>
        </Card>

        <!-- Realtime Clicks Stream / Feed -->
        <Card class="border border-border/60 shadow-xs flex flex-col">
          <CardHeader class="pb-3">
            <div class="flex items-center justify-between">
              <CardTitle class="text-base font-bold">Recent Click Activity</CardTitle>
              <span class="flex h-2 w-2 rounded-full bg-emerald-500 animate-pulse"></span>
            </div>
            <CardDescription class="text-xs">Live events across your links</CardDescription>
          </CardHeader>
          <CardContent class="flex-1 overflow-y-auto max-h-[300px] space-y-3 pt-0">
            <div
              v-if="!urlStore.overview?.recent_clicks || urlStore.overview.recent_clicks.length === 0"
              class="flex h-40 items-center justify-center text-xs text-muted-foreground"
            >
              No clicks recorded recently
            </div>
            <div
              v-for="click in urlStore.overview?.recent_clicks"
              :key="click.id"
              class="flex items-center justify-between rounded-lg border border-border/40 p-2.5 text-xs bg-card/60 transition-colors hover:bg-muted/40"
            >
              <div class="flex flex-col space-y-0.5 truncate mr-2">
                <div class="flex items-center gap-1.5 font-medium">
                  <Globe class="h-3.5 w-3.5 text-muted-foreground shrink-0" />
                  <span class="truncate">{{ click.country || 'Unknown Country' }}</span>
                  <span v-if="click.city" class="text-muted-foreground font-normal">({{ click.city }})</span>
                </div>
                <div class="text-[11px] text-muted-foreground flex items-center gap-1">
                  <span>{{ click.device_type }}</span>
                  <span>•</span>
                  <span>{{ click.browser }}</span>
                  <span>•</span>
                  <span>{{ click.referrer }}</span>
                </div>
              </div>
              <span class="text-[10px] text-muted-foreground whitespace-nowrap">
                {{ formatRelativeTime(click.clicked_at) }}
              </span>
            </div>
          </CardContent>
        </Card>
      </div>

      <!-- Recent Links Table Section -->
      <Card class="border border-border/60 shadow-xs">
        <CardHeader class="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
          <div>
            <CardTitle class="text-base font-bold">Recent Short Links</CardTitle>
            <CardDescription class="text-xs">Manage and inspect performance of your latest links</CardDescription>
          </div>
          <Button variant="outline" size="sm" class="gap-1.5 text-xs font-medium" @click="router.push('/links')">
            <span>View All Links</span>
            <ArrowUpRight class="h-3.5 w-3.5" />
          </Button>
        </CardHeader>
        <CardContent class="p-0">
          <div class="overflow-x-auto">
            <table class="w-full text-left text-xs">
              <thead class="bg-muted/40 border-y border-border/40 text-muted-foreground font-medium">
                <tr>
                  <th class="py-3 px-4 sm:px-6">Short Link & Title</th>
                  <th class="py-3 px-4 hidden md:table-cell">Destination</th>
                  <th class="py-3 px-4 text-center">Clicks</th>
                  <th class="py-3 px-4 text-center">Status</th>
                  <th class="py-3 px-4 sm:px-6 text-right">Actions</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-border/30">
                <tr v-if="urlStore.urls.length === 0">
                  <td colspan="5" class="py-8 text-center text-muted-foreground">
                    No short links created yet. Create your first link above!
                  </td>
                </tr>
                <tr
                  v-for="urlItem in urlStore.urls.slice(0, 5)"
                  :key="urlItem.id"
                  class="group hover:bg-muted/30 transition-colors"
                >
                  <!-- Link & Title -->
                  <td class="py-3.5 px-4 sm:px-6">
                    <div class="flex flex-col space-y-0.5">
                      <div class="flex items-center gap-2">
                        <a
                          :href="`${baseUrl}/${urlItem.short_code}`"
                          target="_blank"
                          class="font-mono font-semibold text-primary hover:underline flex items-center gap-1"
                        >
                          /{{ urlItem.short_code }}
                          <ExternalLink class="h-3 w-3 opacity-0 group-hover:opacity-100 transition-opacity" />
                        </a>
                        <Badge v-if="urlItem.has_password" variant="secondary" class="h-4 px-1 text-[10px] gap-0.5">
                          <Lock class="h-2.5 w-2.5" /> Password
                        </Badge>
                      </div>
                      <span class="text-muted-foreground truncate max-w-[200px] sm:max-w-[280px]">
                        {{ urlItem.title || 'Untitled link' }}
                      </span>
                    </div>
                  </td>

                  <!-- Destination URL -->
                  <td class="py-3.5 px-4 hidden md:table-cell max-w-[260px] truncate text-muted-foreground">
                    {{ urlItem.original_url }}
                  </td>

                  <!-- Clicks Count -->
                  <td class="py-3.5 px-4 text-center">
                    <span class="font-bold text-foreground">{{ urlItem.click_count }}</span>
                  </td>

                  <!-- Status Badge -->
                  <td class="py-3.5 px-4 text-center">
                    <Badge
                      :variant="urlItem.is_active ? 'default' : 'secondary'"
                      class="text-[10px] uppercase font-semibold tracking-wider"
                    >
                      {{ urlItem.is_active ? 'Active' : 'Disabled' }}
                    </Badge>
                  </td>

                  <!-- Actions -->
                  <td class="py-3.5 px-4 sm:px-6 text-right">
                    <div class="flex items-center justify-end gap-1">
                      <!-- Copy Button -->
                      <Button
                        variant="ghost"
                        size="icon"
                        class="h-8 w-8 text-muted-foreground hover:text-foreground"
                        title="Copy Short Link"
                        @click="copyShortLink(urlItem)"
                      >
                        <Check v-if="copiedId === urlItem.id" class="h-4 w-4 text-emerald-500" />
                        <Copy v-else class="h-4 w-4" />
                      </Button>

                      <!-- QR Code Button -->
                      <Button
                        variant="ghost"
                        size="icon"
                        class="h-8 w-8 text-muted-foreground hover:text-foreground"
                        title="Generate QR Code"
                        @click="openQr(urlItem)"
                      >
                        <QrCode class="h-4 w-4" />
                      </Button>

                      <!-- Analytics Button -->
                      <Button
                        variant="ghost"
                        size="icon"
                        class="h-8 w-8 text-muted-foreground hover:text-foreground"
                        title="View Full Analytics"
                        @click="router.push(`/analytics/${urlItem.id}`)"
                      >
                        <BarChart3 class="h-4 w-4" />
                      </Button>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </CardContent>
      </Card>
    </main>

    <!-- Modals -->
    <CreateEditUrlDialog v-model:open="isCreateModalOpen" />
    <QrCodeModal v-model:open="isQrModalOpen" :url-item="selectedUrlForQr" />
  </div>
</template>
