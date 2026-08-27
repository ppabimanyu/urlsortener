<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import confetti from 'canvas-confetti'
import { useUrlStore } from '@/stores/url'
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
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Switch } from '@/components/ui/switch'
import {
  Link2,
  Sparkles,
  Lock,
  Calendar,
  MousePointerClick,
  AlertCircle,
  CheckCircle2,
} from 'lucide-vue-next'

const props = defineProps<{
  open: boolean
  urlToEdit?: URLItem | null
}>()

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
  (e: 'created', url: URLItem): void
  (e: 'updated', url: URLItem): void
}>()

const urlStore = useUrlStore()

const isEditing = computed(() => !!props.urlToEdit)
const isOpen = computed({
  get: () => props.open,
  set: (val) => emit('update:open', val),
})

// Form state
const originalUrl = ref('')
const title = ref('')
const customSlug = ref('')
const enablePassword = ref(false)
const password = ref('')
const enableExpiration = ref(false)
const expiresAt = ref('')
const enableClickLimit = ref(false)
const clickLimit = ref<number | undefined>(undefined)

const loading = ref(false)
const error = ref<string | null>(null)
const successMessage = ref<string | null>(null)

watch(
  () => props.urlToEdit,
  (item) => {
    if (item) {
      originalUrl.value = item.original_url
      title.value = item.title || ''
      customSlug.value = item.short_code || ''
      enablePassword.value = item.has_password || false
      password.value = ''
      enableExpiration.value = !!item.expires_at
      expiresAt.value = item.expires_at ? new Date(item.expires_at).toISOString().slice(0, 16) : ''
      enableClickLimit.value = item.click_limit !== null && item.click_limit !== undefined
      clickLimit.value = item.click_limit ?? undefined
    } else {
      resetForm()
    }
  },
  { immediate: true }
)

function resetForm() {
  originalUrl.value = ''
  title.value = ''
  customSlug.value = ''
  enablePassword.value = false
  password.value = ''
  enableExpiration.value = false
  expiresAt.value = ''
  enableClickLimit.value = false
  clickLimit.value = undefined
  error.value = null
  successMessage.value = null
}

async function handleSubmit() {
  error.value = null
  successMessage.value = null

  if (!originalUrl.value.trim()) {
    error.value = 'Please enter a valid destination URL'
    return
  }

  loading.value = true
  try {
    const payload: any = {
      original_url: originalUrl.value.trim(),
      title: title.value.trim() || undefined,
      custom_slug: customSlug.value.trim() || undefined,
    }

    if (enablePassword.value && password.value) {
      payload.password = password.value
    } else if (!enablePassword.value && isEditing.value) {
      payload.password = '' // Clear password
    }

    if (enableExpiration.value && expiresAt.value) {
      payload.expires_at = new Date(expiresAt.value).toISOString()
    } else if (!enableExpiration.value) {
      payload.expires_at = null
    }

    if (enableClickLimit.value && clickLimit.value && clickLimit.value > 0) {
      payload.click_limit = Number(clickLimit.value)
    } else if (!enableClickLimit.value) {
      payload.click_limit = null
    }

    if (isEditing.value && props.urlToEdit) {
      const updated = await urlStore.updateUrl(props.urlToEdit.id, payload)
      emit('updated', updated)
      isOpen.value = false
    } else {
      const created = await urlStore.createUrl(payload)
      // Trigger confetti!
      confetti({
        particleCount: 80,
        spread: 60,
        origin: { y: 0.6 },
      })
      emit('created', created)
      isOpen.value = false
    }
  } catch (err: any) {
    error.value = err.response?.data?.error || 'Failed to save short URL. Please check your inputs.'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <Dialog v-model:open="isOpen">
    <DialogContent class="sm:max-w-[550px] max-h-[90vh] overflow-y-auto">
      <DialogHeader>
        <DialogTitle class="flex items-center gap-2 text-xl font-bold">
          <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-primary/10 text-primary">
            <Link2 class="h-4 w-4" />
          </div>
          {{ isEditing ? 'Edit Short Link' : 'Create New Short Link' }}
        </DialogTitle>
        <DialogDescription>
          {{ isEditing ? 'Update destination, custom alias, or protection settings.' : 'Enter your destination URL and customize slug, password, or expiration.' }}
        </DialogDescription>
      </DialogHeader>

      <form @submit.prevent="handleSubmit" class="space-y-4 py-2">
        <!-- Error alert -->
        <div v-if="error" class="flex items-center gap-2 rounded-lg bg-destructive/10 p-3 text-sm text-destructive border border-destructive/20">
          <AlertCircle class="h-4 w-4 shrink-0" />
          <span>{{ error }}</span>
        </div>

        <!-- Destination URL -->
        <div class="space-y-1.5">
          <Label for="original-url" class="text-sm font-semibold">
            Destination URL <span class="text-destructive">*</span>
          </Label>
          <Input
            id="original-url"
            v-model="originalUrl"
            type="url"
            placeholder="https://example.com/very-long-url-path..."
            required
            class="h-10"
          />
        </div>

        <!-- Title / Label -->
        <div class="space-y-1.5">
          <Label for="url-title" class="text-sm font-semibold">
            Link Title <span class="text-xs font-normal text-muted-foreground">(Optional)</span>
          </Label>
          <Input
            id="url-title"
            v-model="title"
            placeholder="e.g. My Portfolio or Summer Campaign"
            class="h-10"
          />
        </div>

        <!-- Custom Slug / Alias -->
        <div class="space-y-1.5">
          <Label for="custom-slug" class="text-sm font-semibold flex items-center justify-between">
            <span>Custom Slug / Alias <span class="text-xs font-normal text-muted-foreground">(Optional)</span></span>
            <span class="text-xs text-muted-foreground">3-50 chars: a-z, 0-9, -, _</span>
          </Label>
          <div class="relative flex items-center">
            <span class="absolute left-3 text-xs text-muted-foreground font-mono select-none">
              short/
            </span>
            <Input
              id="custom-slug"
              v-model="customSlug"
              placeholder="my-custom-slug"
              class="h-10 pl-16 font-mono text-sm"
            />
          </div>
        </div>

        <!-- Advanced Features Accordion / Toggles -->
        <div class="pt-2 space-y-4 border-t border-border/60">
          <p class="text-xs font-semibold text-muted-foreground uppercase tracking-wider">Advanced Settings</p>

          <!-- Password Protection -->
          <div class="rounded-xl border border-border/60 p-3.5 bg-muted/30 space-y-3">
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-2.5">
                <Lock class="h-4 w-4 text-muted-foreground" />
                <div>
                  <p class="text-sm font-medium leading-none">Password Protection</p>
                  <p class="text-xs text-muted-foreground mt-0.5">Require visitors to enter a password to redirect</p>
                </div>
              </div>
              <Switch v-model="enablePassword" />
            </div>
            <div v-if="enablePassword" class="pt-1">
              <Input
                v-model="password"
                type="password"
                placeholder="Enter secret passcode..."
                class="h-9 text-sm"
              />
            </div>
          </div>

          <!-- Expiration Date -->
          <div class="rounded-xl border border-border/60 p-3.5 bg-muted/30 space-y-3">
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-2.5">
                <Calendar class="h-4 w-4 text-muted-foreground" />
                <div>
                  <p class="text-sm font-medium leading-none">Expiration Date</p>
                  <p class="text-xs text-muted-foreground mt-0.5">Automatically deactivate link after this date</p>
                </div>
              </div>
              <Switch v-model="enableExpiration" />
            </div>
            <div v-if="enableExpiration" class="pt-1">
              <Input
                v-model="expiresAt"
                type="datetime-local"
                class="h-9 text-sm"
              />
            </div>
          </div>

          <!-- Click Limit -->
          <div class="rounded-xl border border-border/60 p-3.5 bg-muted/30 space-y-3">
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-2.5">
                <MousePointerClick class="h-4 w-4 text-muted-foreground" />
                <div>
                  <p class="text-sm font-medium leading-none">Click Limit</p>
                  <p class="text-xs text-muted-foreground mt-0.5">Deactivate after reaching total clicks</p>
                </div>
              </div>
              <Switch v-model="enableClickLimit" />
            </div>
            <div v-if="enableClickLimit" class="pt-1">
              <Input
                v-model.number="clickLimit"
                type="number"
                min="1"
                placeholder="e.g. 100"
                class="h-9 text-sm"
              />
            </div>
          </div>
        </div>

        <DialogFooter class="pt-4 border-t border-border/40 gap-2">
          <Button type="button" variant="outline" @click="isOpen = false" :disabled="loading">
            Cancel
          </Button>
          <Button type="submit" :disabled="loading" class="gap-2 shadow-sm font-medium">
            <Sparkles v-if="!loading" class="h-4 w-4" />
            <span>{{ loading ? 'Saving...' : isEditing ? 'Update Link' : 'Create Short Link' }}</span>
          </Button>
        </DialogFooter>
      </form>
    </DialogContent>
  </Dialog>
</template>
