<script setup lang="ts">
import { computed } from 'vue'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { AlertTriangle } from 'lucide-vue-next'

const props = defineProps<{
  open: boolean
  title?: string
  description?: string
  loading?: boolean
}>()

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
  (e: 'confirm'): void
}>()

const isOpen = computed({
  get: () => props.open,
  set: (val) => emit('update:open', val),
})
</script>

<template>
  <Dialog v-model:open="isOpen">
    <DialogContent class="sm:max-w-[420px]">
      <DialogHeader>
        <DialogTitle class="flex items-center gap-2 text-lg font-bold text-destructive">
          <AlertTriangle class="h-5 w-5" />
          {{ title || 'Delete Short Link' }}
        </DialogTitle>
        <DialogDescription>
          {{ description || 'Are you sure you want to delete this link? All recorded analytics for this link will be permanently removed. This action cannot be undone.' }}
        </DialogDescription>
      </DialogHeader>
      <DialogFooter class="pt-4 gap-2">
        <Button variant="outline" @click="isOpen = false" :disabled="loading">
          Cancel
        </Button>
        <Button variant="destructive" @click="emit('confirm')" :disabled="loading">
          {{ loading ? 'Deleting...' : 'Delete Permanently' }}
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
