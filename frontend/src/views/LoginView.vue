<script setup lang="ts">
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { Link2, ArrowRight, Lock, Mail, AlertCircle, Loader2 } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import ThemeToggle from '@/components/layout/ThemeToggle.vue'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const email = ref('')
const password = ref('')

async function handleLogin() {
  const success = await authStore.login(email.value, password.value)
  if (success) {
    const redirectPath = (route.query.redirect as string) || '/dashboard'
    router.push(redirectPath)
  }
}
</script>

<template>
  <div class="relative flex min-h-screen flex-col items-center justify-center p-4 bg-muted/20">
    <!-- Top Right Theme Toggle -->
    <div class="absolute top-4 right-4">
      <ThemeToggle />
    </div>

    <!-- Background glow decoration -->
    <div class="absolute inset-0 -z-10 flex items-center justify-center overflow-hidden">
      <div class="h-[300px] w-[500px] rounded-full bg-primary/5 blur-3xl"></div>
    </div>

    <div class="w-full max-w-md space-y-6">
      <!-- Brand Header -->
      <div class="flex flex-col items-center text-center space-y-2">
        <div class="flex h-12 w-12 items-center justify-center rounded-2xl bg-primary text-primary-foreground shadow-md">
          <Link2 class="h-6 w-6" />
        </div>
        <h1 class="text-2xl font-bold tracking-tight">SortLink</h1>
        <p class="text-xs text-muted-foreground">Modern URL Shortener with Realtime Analytics</p>
      </div>

      <Card class="border border-border/60 shadow-lg bg-card">
        <CardHeader class="space-y-1">
          <CardTitle class="text-xl">Sign in</CardTitle>
          <CardDescription>Enter your email and password to access your dashboard</CardDescription>
        </CardHeader>
        <CardContent>
          <form @submit.prevent="handleLogin" class="space-y-4">
            <!-- Error message -->
            <div v-if="authStore.error" class="flex items-center gap-2 rounded-lg bg-destructive/10 p-3 text-xs text-destructive border border-destructive/20">
              <AlertCircle class="h-4 w-4 shrink-0" />
              <span>{{ authStore.error }}</span>
            </div>

            <!-- Email -->
            <div class="space-y-1.5">
              <Label for="email" class="text-xs font-semibold">Email address</Label>
              <div class="relative flex items-center">
                <Mail class="absolute left-3 h-4 w-4 text-muted-foreground" />
                <Input
                  id="email"
                  v-model="email"
                  type="email"
                  placeholder="name@example.com"
                  required
                  class="pl-9 h-10"
                />
              </div>
            </div>

            <!-- Password -->
            <div class="space-y-1.5">
              <div class="flex items-center justify-between">
                <Label for="password" class="text-xs font-semibold">Password</Label>
              </div>
              <div class="relative flex items-center">
                <Lock class="absolute left-3 h-4 w-4 text-muted-foreground" />
                <Input
                  id="password"
                  v-model="password"
                  type="password"
                  placeholder="••••••••"
                  required
                  class="pl-9 h-10"
                />
              </div>
            </div>

            <Button type="submit" class="w-full h-10 gap-2 font-medium" :disabled="authStore.loading">
              <Loader2 v-if="authStore.loading" class="h-4 w-4 animate-spin" />
              <span v-else>Sign In</span>
              <ArrowRight v-if="!authStore.loading" class="h-4 w-4" />
            </Button>
          </form>
        </CardContent>
        <CardFooter class="flex flex-col items-center justify-center border-t border-border/40 p-4 text-center">
          <p class="text-xs text-muted-foreground">
            Don't have an account?
            <router-link to="/register" class="font-semibold text-primary hover:underline ml-1">
              Create an account
            </router-link>
          </p>
        </CardFooter>
      </Card>
    </div>
  </div>
</template>
