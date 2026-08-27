<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import api from '@/services/api'
import { Lock, ArrowRight, AlertCircle, Loader2, ShieldCheck } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import ThemeToggle from '@/components/layout/ThemeToggle.vue'

const route = useRoute()
const slug = ref((route.params.slug as string) || '')

const password = ref('')
const loading = ref(false)
const error = ref<string | null>(null)
const linkTitle = ref<string>('')

onMounted(async () => {
  if (slug.value) {
    try {
      const res = await api.get(`/public/link-info/${slug.value}`)
      linkTitle.value = res.data.title || ''
    } catch (e) {
      // Ignore error if link info fails
    }
  }
})

async function handleVerify() {
  if (!password.value) return
  loading.value = true
  error.value = null

  try {
    const res = await api.post('/public/verify-password', {
      code: slug.value,
      password: password.value,
    })

    if (res.data.original_url) {
      // Redirect to destination
      window.location.href = res.data.original_url
    }
  } catch (err: any) {
    error.value = err.response?.data?.error || 'Invalid password. Please try again.'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="relative flex min-h-screen flex-col items-center justify-center p-4 bg-muted/20">
    <div class="absolute top-4 right-4">
      <ThemeToggle />
    </div>

    <div class="w-full max-w-md space-y-6">
      <div class="flex flex-col items-center text-center space-y-2">
        <div class="flex h-12 w-12 items-center justify-center rounded-2xl bg-amber-500/10 text-amber-500 shadow-sm border border-amber-500/20">
          <Lock class="h-6 w-6" />
        </div>
        <h1 class="text-xl font-bold tracking-tight">Protected Short Link</h1>
        <p class="text-xs text-muted-foreground">
          This link <span class="font-mono font-semibold text-foreground">/{{ slug }}</span> requires a password to access
        </p>
      </div>

      <Card class="border border-border/60 shadow-lg bg-card">
        <CardHeader class="space-y-1 pb-3">
          <CardTitle class="text-base font-bold">{{ linkTitle || 'Enter Passcode' }}</CardTitle>
          <CardDescription class="text-xs">Please provide the secret password set by the creator</CardDescription>
        </CardHeader>
        <CardContent>
          <form @submit.prevent="handleVerify" class="space-y-4">
            <div v-if="error" class="flex items-center gap-2 rounded-lg bg-destructive/10 p-3 text-xs text-destructive border border-destructive/20">
              <AlertCircle class="h-4 w-4 shrink-0" />
              <span>{{ error }}</span>
            </div>

            <div class="space-y-1.5">
              <Label for="password" class="text-xs font-semibold">Password</Label>
              <div class="relative flex items-center">
                <Lock class="absolute left-3 h-4 w-4 text-muted-foreground" />
                <Input
                  id="password"
                  v-model="password"
                  type="password"
                  placeholder="Enter passcode..."
                  required
                  autofocus
                  class="pl-9 h-10 text-sm"
                />
              </div>
            </div>

            <Button type="submit" class="w-full h-10 gap-2 font-medium" :disabled="loading">
              <Loader2 v-if="loading" class="h-4 w-4 animate-spin" />
              <span v-else>Unlock & Continue</span>
              <ArrowRight v-if="!loading" class="h-4 w-4" />
            </Button>
          </form>
        </CardContent>
        <CardFooter class="flex items-center justify-center border-t border-border/40 p-3 text-center">
          <div class="flex items-center gap-1 text-[11px] text-muted-foreground">
            <ShieldCheck class="h-3.5 w-3.5 text-emerald-500" />
            <span>End-to-end access validation by SortLink</span>
          </div>
        </CardFooter>
      </Card>
    </div>
  </div>
</template>
