<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  AlertTriangle,
  Clock,
  Ban,
  FileQuestion,
  ArrowLeft,
  Home,
} from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import ThemeToggle from '@/components/layout/ThemeToggle.vue'

const props = withDefaults(
  defineProps<{
    statusType?: 'not-found' | 'inactive' | 'expired' | 'limit' | 'error'
  }>(),
  {
    statusType: 'not-found',
  }
)

const route = useRoute()
const router = useRouter()

const slug = computed(() => (route.query.slug as string) || '')

const config = computed(() => {
  switch (props.statusType) {
    case 'inactive':
      return {
        icon: Ban,
        iconColor: 'text-amber-500 bg-amber-500/10 border-amber-500/20',
        title: 'Link is Currently Inactive',
        description: 'The owner has temporarily disabled this shortened link. Please try again later.',
      }
    case 'expired':
      return {
        icon: Clock,
        iconColor: 'text-rose-500 bg-rose-500/10 border-rose-500/20',
        title: 'Link Has Expired',
        description: 'This short URL was configured with an expiration date and is no longer available.',
      }
    case 'limit':
      return {
        icon: AlertTriangle,
        iconColor: 'text-orange-500 bg-orange-500/10 border-orange-500/20',
        title: 'Click Limit Reached',
        description: 'This link has reached its maximum allowed number of clicks and has been deactivated.',
      }
    case 'error':
      return {
        icon: AlertTriangle,
        iconColor: 'text-destructive bg-destructive/10 border-destructive/20',
        title: 'Something Went Wrong',
        description: 'An error occurred while resolving this short link. Please check the URL and try again.',
      }
    case 'not-found':
    default:
      return {
        icon: FileQuestion,
        iconColor: 'text-slate-500 bg-slate-500/10 border-slate-500/20',
        title: 'Link Not Found (404)',
        description: 'The short link you are looking for does not exist or may have been deleted.',
      }
  }
})
</script>

<template>
  <div class="relative flex min-h-screen flex-col items-center justify-center p-4 bg-muted/20">
    <div class="absolute top-4 right-4">
      <ThemeToggle />
    </div>

    <div class="w-full max-w-md space-y-6">
      <Card class="border border-border/60 shadow-lg text-center bg-card">
        <CardHeader class="flex flex-col items-center space-y-3 pt-8 pb-4">
          <div
            class="flex h-16 w-16 items-center justify-center rounded-3xl border shadow-xs"
            :class="config.iconColor"
          >
            <component :is="config.icon" class="h-8 w-8" />
          </div>
          <CardTitle class="text-xl font-bold">{{ config.title }}</CardTitle>
          <CardDescription class="text-xs max-w-xs mx-auto">
            {{ config.description }}
          </CardDescription>
        </CardHeader>
        <CardContent v-if="slug" class="pb-4">
          <div class="inline-flex items-center rounded-lg bg-muted px-3 py-1 font-mono text-xs text-muted-foreground">
            /{{ slug }}
          </div>
        </CardContent>
        <CardFooter class="flex flex-col sm:flex-row gap-2 justify-center border-t border-border/40 p-4">
          <Button variant="outline" class="w-full sm:w-auto text-xs gap-1.5" @click="router.push('/dashboard')">
            <Home class="h-3.5 w-3.5" />
            <span>Go to Dashboard</span>
          </Button>
        </CardFooter>
      </Card>
    </div>
  </div>
</template>
