<script setup lang="ts">
import { ref, onMounted, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useUrlStore } from '@/stores/url'
import type { URLItem } from '@/types'

import Navbar from '@/components/layout/Navbar.vue'
import CreateEditUrlDialog from '@/components/dialogs/CreateEditUrlDialog.vue'
import QrCodeModal from '@/components/dialogs/QrCodeModal.vue'
import DeleteConfirmDialog from '@/components/dialogs/DeleteConfirmDialog.vue'

import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Switch } from '@/components/ui/switch'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import {
  Search,
  Plus,
  Copy,
  Check,
  QrCode,
  BarChart3,
  ExternalLink,
  Lock,
  Calendar,
  MousePointerClick,
  MoreVertical,
  Edit2,
  Trash2,
  Filter,
} from 'lucide-vue-next'

const router = useRouter()
const urlStore = useUrlStore()

const searchQuery = ref('')
const selectedStatus = ref<'all' | 'active' | 'inactive'>('all')

const isCreateEditOpen = ref(false)
const urlToEdit = ref<URLItem | null>(null)

const isQrOpen = ref(false)
const selectedUrlForQr = ref<URLItem | null>(null)

const isDeleteOpen = ref(false)
const urlToDelete = ref<URLItem | null>(null)
const deleteLoading = ref(false)

const copiedId = ref<string | null>(null)

const baseUrl = computed(() => import.meta.env.VITE_BACKEND_URL || 'http://localhost:8080')

onMounted(() => {
  fetchData()
})

function getActiveFilter() {
  if (selectedStatus.value === 'active') return true
  if (selectedStatus.value === 'inactive') return false
  return null
}

async function fetchData(page = 1) {
  await urlStore.fetchUrls(page, searchQuery.value, getActiveFilter())
}

let searchDebounce: any = null
watch([searchQuery, selectedStatus], () => {
  clearTimeout(searchDebounce)
  searchDebounce = setTimeout(() => {
    fetchData(1)
  }, 300)
})

function copyLink(urlItem: URLItem) {
  const full = `${baseUrl.value}/${urlItem.short_code}`
  navigator.clipboard.writeText(full)
  copiedId.value = urlItem.id
  setTimeout(() => {
    if (copiedId.value === urlItem.id) copiedId.value = null
  }, 2000)
}

function openEdit(urlItem: URLItem) {
  urlToEdit.value = urlItem
  isCreateEditOpen.value = true
}

function openCreate() {
  urlToEdit.value = null
  isCreateEditOpen.value = true
}

function openQr(urlItem: URLItem) {
  selectedUrlForQr.value = urlItem
  isQrOpen.value = true
}

function promptDelete(urlItem: URLItem) {
  urlToDelete.value = urlItem
  isDeleteOpen.value = true
}

async function handleDelete() {
  if (!urlToDelete.value) return
  deleteLoading.value = true
  try {
    await urlStore.deleteUrl(urlToDelete.value.id)
    isDeleteOpen.value = false
    urlToDelete.value = null
  } catch (err) {
    console.error('Failed to delete URL:', err)
  } finally {
    deleteLoading.value = false
  }
}

async function handleToggleStatus(urlItem: URLItem) {
  try {
    await urlStore.toggleUrlStatus(urlItem.id)
  } catch (err) {
    console.error('Failed to toggle status:', err)
  }
}

function formatDate(dateStr: string) {
  return new Date(dateStr).toLocaleDateString('id-ID', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  })
}
</script>

<template>
  <div class="min-h-screen bg-muted/20 pb-16">
    <Navbar />

    <main class="container mx-auto max-w-7xl px-4 sm:px-6 lg:px-8 pt-8 space-y-6">
      <!-- Header -->
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 class="text-2xl font-bold tracking-tight text-foreground">My Short Links</h1>
          <p class="text-xs text-muted-foreground mt-0.5">Manage, customize, and analyze all your links in one place</p>
        </div>
        <Button class="gap-1.5 font-medium shadow-sm" @click="openCreate">
          <Plus class="h-4 w-4" />
          <span>Create New Link</span>
        </Button>
      </div>

      <!-- Filters & Search Bar -->
      <div class="flex flex-col sm:flex-row items-center justify-between gap-3">
        <div class="relative w-full sm:max-w-md">
          <Search class="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
          <Input
            v-model="searchQuery"
            placeholder="Search by title, custom slug, or URL..."
            class="pl-9 h-10 bg-card text-xs"
          />
        </div>

        <div class="flex w-full sm:w-auto items-center gap-2">
          <div class="flex items-center rounded-lg border border-border/60 bg-card p-1 text-xs">
            <button
              class="px-3 py-1 rounded-md transition-colors"
              :class="selectedStatus === 'all' ? 'bg-primary text-primary-foreground font-semibold' : 'text-muted-foreground hover:text-foreground'"
              @click="selectedStatus = 'all'"
            >
              All ({{ urlStore.total }})
            </button>
            <button
              class="px-3 py-1 rounded-md transition-colors"
              :class="selectedStatus === 'active' ? 'bg-primary text-primary-foreground font-semibold' : 'text-muted-foreground hover:text-foreground'"
              @click="selectedStatus = 'active'"
            >
              Active
            </button>
            <button
              class="px-3 py-1 rounded-md transition-colors"
              :class="selectedStatus === 'inactive' ? 'bg-primary text-primary-foreground font-semibold' : 'text-muted-foreground hover:text-foreground'"
              @click="selectedStatus = 'inactive'"
            >
              Disabled
            </button>
          </div>
        </div>
      </div>

      <!-- Links Table Card -->
      <Card class="border border-border/60 shadow-xs">
        <CardContent class="p-0">
          <div class="overflow-x-auto">
            <table class="w-full text-left text-xs">
              <thead class="bg-muted/40 border-b border-border/40 text-muted-foreground font-medium">
                <tr>
                  <th class="py-3.5 px-4 sm:px-6">Short Link & Meta</th>
                  <th class="py-3.5 px-4 hidden md:table-cell">Target Destination</th>
                  <th class="py-3.5 px-4 text-center">Clicks / Limit</th>
                  <th class="py-3.5 px-4 text-center">Active</th>
                  <th class="py-3.5 px-4 sm:px-6 text-right">Actions</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-border/30">
                <tr v-if="urlStore.loading">
                  <td colspan="5" class="py-12 text-center text-muted-foreground">
                    Loading short links...
                  </td>
                </tr>
                <tr v-else-if="urlStore.urls.length === 0">
                  <td colspan="5" class="py-12 text-center text-muted-foreground">
                    No links found matching your criteria.
                  </td>
                </tr>
                <tr
                  v-for="urlItem in urlStore.urls"
                  :key="urlItem.id"
                  class="group hover:bg-muted/30 transition-colors"
                >
                  <!-- Link info -->
                  <td class="py-4 px-4 sm:px-6">
                    <div class="flex flex-col space-y-1">
                      <div class="flex flex-wrap items-center gap-2">
                        <a
                          :href="`${baseUrl}/${urlItem.short_code}`"
                          target="_blank"
                          class="font-mono font-bold text-primary hover:underline flex items-center gap-1 text-sm"
                        >
                          /{{ urlItem.short_code }}
                          <ExternalLink class="h-3.5 w-3.5 opacity-60 group-hover:opacity-100" />
                        </a>

                        <!-- Badges -->
                        <Badge v-if="urlItem.has_password" variant="secondary" class="h-5 px-1.5 text-[10px] gap-1">
                          <Lock class="h-3 w-3" /> Password
                        </Badge>
                        <Badge v-if="urlItem.expires_at" variant="outline" class="h-5 px-1.5 text-[10px] gap-1">
                          <Calendar class="h-3 w-3 text-muted-foreground" /> Exp: {{ formatDate(urlItem.expires_at) }}
                        </Badge>
                      </div>
                      <p class="font-medium text-foreground truncate max-w-[240px] sm:max-w-[320px]">
                        {{ urlItem.title || 'Untitled link' }}
                      </p>
                      <span class="text-[10px] text-muted-foreground">
                        Created {{ formatDate(urlItem.created_at) }}
                      </span>
                    </div>
                  </td>

                  <!-- Target Destination -->
                  <td class="py-4 px-4 hidden md:table-cell max-w-[260px]">
                    <span class="truncate block text-muted-foreground" :title="urlItem.original_url">
                      {{ urlItem.original_url }}
                    </span>
                  </td>

                  <!-- Clicks & Limit -->
                  <td class="py-4 px-4 text-center">
                    <div class="flex flex-col items-center">
                      <span class="font-bold text-foreground text-sm">{{ urlItem.click_count }}</span>
                      <span v-if="urlItem.click_limit" class="text-[10px] text-muted-foreground">
                        / {{ urlItem.click_limit }} max
                      </span>
                      <span v-else class="text-[10px] text-muted-foreground">clicks</span>
                    </div>
                  </td>

                  <!-- Active Toggle -->
                  <td class="py-4 px-4 text-center">
                    <div class="flex justify-center">
                      <Switch
                        :model-value="urlItem.is_active"
                        @update:model-value="() => handleToggleStatus(urlItem)"
                      />
                    </div>
                  </td>

                  <!-- Actions -->
                  <td class="py-4 px-4 sm:px-6 text-right">
                    <div class="flex items-center justify-end gap-1">
                      <Button
                        variant="ghost"
                        size="icon"
                        class="h-8 w-8 text-muted-foreground hover:text-foreground"
                        title="Copy Short Link"
                        @click="copyLink(urlItem)"
                      >
                        <Check v-if="copiedId === urlItem.id" class="h-4 w-4 text-emerald-500" />
                        <Copy v-else class="h-4 w-4" />
                      </Button>

                      <Button
                        variant="ghost"
                        size="icon"
                        class="h-8 w-8 text-muted-foreground hover:text-foreground"
                        title="Generate QR Code"
                        @click="openQr(urlItem)"
                      >
                        <QrCode class="h-4 w-4" />
                      </Button>

                      <Button
                        variant="ghost"
                        size="icon"
                        class="h-8 w-8 text-muted-foreground hover:text-foreground"
                        title="Analytics"
                        @click="router.push(`/analytics/${urlItem.id}`)"
                      >
                        <BarChart3 class="h-4 w-4" />
                      </Button>

                      <!-- More dropdown -->
                      <DropdownMenu>
                        <DropdownMenuTrigger as-child>
                          <Button variant="ghost" size="icon" class="h-8 w-8 text-muted-foreground">
                            <MoreVertical class="h-4 w-4" />
                          </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end" class="w-40">
                          <DropdownMenuItem @click="openEdit(urlItem)">
                            <Edit2 class="mr-2 h-3.5 w-3.5" />
                            <span>Edit Link</span>
                          </DropdownMenuItem>
                          <DropdownMenuItem @click="router.push(`/analytics/${urlItem.id}`)">
                            <BarChart3 class="mr-2 h-3.5 w-3.5" />
                            <span>Full Analytics</span>
                          </DropdownMenuItem>
                          <DropdownMenuSeparator />
                          <DropdownMenuItem class="text-destructive focus:text-destructive focus:bg-destructive/10" @click="promptDelete(urlItem)">
                            <Trash2 class="mr-2 h-3.5 w-3.5" />
                            <span>Delete</span>
                          </DropdownMenuItem>
                        </DropdownMenuContent>
                      </DropdownMenu>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </CardContent>
      </Card>
    </main>

    <!-- Dialogs -->
    <CreateEditUrlDialog
      v-model:open="isCreateEditOpen"
      :url-to-edit="urlToEdit"
    />
    <QrCodeModal v-model:open="isQrOpen" :url-item="selectedUrlForQr" />
    <DeleteConfirmDialog
      v-model:open="isDeleteOpen"
      :title="`Delete /${urlToDelete?.short_code}`"
      :loading="deleteLoading"
      @confirm="handleDelete"
    />
  </div>
</template>
