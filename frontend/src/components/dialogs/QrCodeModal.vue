<script setup lang="ts">
import { ref, computed, watch, nextTick } from 'vue'
import QRCode from 'qrcode'
import { useClipboard } from '@vueuse/core'
import type { URLItem } from '@/types'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import {
  QrCode,
  Download,
  Copy,
  Check,
  ExternalLink,
} from 'lucide-vue-next'

const props = defineProps<{
  open: boolean
  urlItem?: URLItem | null
}>()

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
}>()

const isOpen = computed({
  get: () => props.open,
  set: (val) => emit('update:open', val),
})

const canvasRef = ref<HTMLCanvasElement | null>(null)
const shortUrl = computed(() => {
  if (!props.urlItem) return ''
  const baseUrl = import.meta.env.VITE_BACKEND_URL || 'http://localhost:8080'
  return `${baseUrl}/${props.urlItem.short_code}`
})

const { copy, copied } = useClipboard({ source: shortUrl })

watch(
  () => [props.open, props.urlItem],
  async ([open, item]) => {
    if (open && item) {
      await nextTick()
      if (canvasRef.value) {
        try {
          await QRCode.toCanvas(canvasRef.value, shortUrl.value, {
            width: 250,
            margin: 2,
            color: {
              dark: '#0f172a',
              light: '#ffffff',
            },
          })
        } catch (err) {
          console.error('Failed to generate QR code:', err)
        }
      }
    }
  },
  { immediate: true }
)

function downloadQr() {
  if (!canvasRef.value || !props.urlItem) return
  const dataUrl = canvasRef.value.toDataURL('image/png')
  const a = document.createElement('a')
  a.href = dataUrl
  a.download = `qr_${props.urlItem.short_code}.png`
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
}
</script>

<template>
  <Dialog v-model:open="isOpen">
    <DialogContent class="sm:max-w-[420px]">
      <DialogHeader class="text-center sm:text-center">
        <DialogTitle class="flex items-center justify-center gap-2 text-xl font-bold">
          <QrCode class="h-5 w-5 text-primary" />
          QR Code
        </DialogTitle>
        <DialogDescription class="text-xs">
          Scan to access destination URL
        </DialogDescription>
      </DialogHeader>

      <div class="flex flex-col items-center justify-center space-y-4 py-3">
        <!-- Canvas Container with rounded frame -->
        <div class="rounded-2xl border-2 border-border/80 bg-white p-3 shadow-md">
          <canvas ref="canvasRef" class="rounded-lg"></canvas>
        </div>

        <!-- Short URL display -->
        <div class="flex w-full items-center justify-between gap-2 rounded-xl bg-muted/60 px-3 py-2 text-xs">
          <span class="truncate font-mono font-medium text-foreground">{{ shortUrl }}</span>
          <Button
            variant="ghost"
            size="icon"
            class="h-7 w-7 shrink-0 text-muted-foreground hover:text-foreground"
            @click="copy(shortUrl)"
          >
            <Check v-if="copied" class="h-3.5 w-3.5 text-emerald-500" />
            <Copy v-else class="h-3.5 w-3.5" />
          </Button>
        </div>
      </div>

      <DialogFooter class="grid grid-cols-2 gap-2 sm:justify-stretch">
        <Button variant="outline" class="gap-1.5 text-xs font-medium" @click="downloadQr">
          <Download class="h-4 w-4" />
          Download PNG
        </Button>
        <a :href="shortUrl" target="_blank" rel="noopener noreferrer" class="contents">
          <Button class="gap-1.5 text-xs font-medium">
            <ExternalLink class="h-4 w-4" />
            Open Link
          </Button>
        </a>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
