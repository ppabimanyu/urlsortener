<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { Link2, ArrowRight, Lock, Mail, User, AlertCircle, Loader2 } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import ThemeToggle from '@/components/layout/ThemeToggle.vue'

const router = useRouter()
const authStore = useAuthStore()

const name = ref('')
const email = ref('')
const password = ref('')
const confirmPassword = ref('')
const localError = ref<string | null>(null)

async function handleRegister() {
  localError.value = null
  if (password.value !== confirmPassword.value) {
    localError.value = 'Passwords do not match'
    return
  }
  if (password.value.length < 6) {
    localError.value = 'Password must be at least 6 characters'
    return
  }

  const success = await authStore.register(name.value, email.value, password.value)
  if (success) {
    router.push('/dashboard')
  }
}
</script>

<template>
  <div class="relative flex min-h-screen flex-col items-center justify-center p-4 bg-muted/20">
    <div class="absolute top-4 right-4">
      <ThemeToggle />
    </div>

    <div class="absolute inset-0 -z-10 flex items-center justify-center overflow-hidden">
      <div class="h-[300px] w-[500px] rounded-full bg-primary/5 blur-3xl"></div>
    </div>

    <div class="w-full max-w-md space-y-6">
      <div class="flex flex-col items-center text-center space-y-2">
        <div class="flex h-12 w-12 items-center justify-center rounded-2xl bg-primary text-primary-foreground shadow-md">
          <Link2 class="h-6 w-6" />
        </div>
        <h1 class="text-2xl font-bold tracking-tight">SortLink</h1>
        <p class="text-xs text-muted-foreground">Create your free URL Shortener account</p>
      </div>

      <Card class="border border-border/60 shadow-lg bg-card">
        <CardHeader class="space-y-1">
          <CardTitle class="text-xl">Create Account</CardTitle>
          <CardDescription>Get started with powerful URL shortening and analytics</CardDescription>
        </CardHeader>
        <CardContent>
          <form @submit.prevent="handleRegister" class="space-y-4">
            <!-- Error message -->
            <div v-if="authStore.error || localError" class="flex items-center gap-2 rounded-lg bg-destructive/10 p-3 text-xs text-destructive border border-destructive/20">
              <AlertCircle class="h-4 w-4 shrink-0" />
              <span>{{ localError || authStore.error }}</span>
            </div>

            <!-- Full Name -->
            <div class="space-y-1.5">
              <Label for="name" class="text-xs font-semibold">Full Name</Label>
              <div class="relative flex items-center">
                <User class="absolute left-3 h-4 w-4 text-muted-foreground" />
                <Input
                  id="name"
                  v-model="name"
                  type="text"
                  placeholder="John Doe"
                  required
                  class="pl-9 h-10"
                />
              </div>
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
              <Label for="password" class="text-xs font-semibold">Password</Label>
              <div class="relative flex items-center">
                <Lock class="absolute left-3 h-4 w-4 text-muted-foreground" />
                <Input
                  id="password"
                  v-model="password"
                  type="password"
                  placeholder="At least 6 characters"
                  required
                  class="pl-9 h-10"
                />
              </div>
            </div>

            <!-- Confirm Password -->
            <div class="space-y-1.5">
              <Label for="confirmPassword" class="text-xs font-semibold">Confirm Password</Label>
              <div class="relative flex items-center">
                <Lock class="absolute left-3 h-4 w-4 text-muted-foreground" />
                <Input
                  id="confirmPassword"
                  v-model="confirmPassword"
                  type="password"
                  placeholder="Repeat your password"
                  required
                  class="pl-9 h-10"
                />
              </div>
            </div>

            <Button type="submit" class="w-full h-10 gap-2 font-medium" :disabled="authStore.loading">
              <Loader2 v-if="authStore.loading" class="h-4 w-4 animate-spin" />
              <span v-else>Create Account</span>
              <ArrowRight v-if="!authStore.loading" class="h-4 w-4" />
            </Button>
          </form>
        </CardContent>
        <CardFooter class="flex flex-col items-center justify-center border-t border-border/40 p-4 text-center">
          <p class="text-xs text-muted-foreground">
            Already have an account?
            <router-link to="/login" class="font-semibold text-primary hover:underline ml-1">
              Sign in instead
            </router-link>
          </p>
        </CardFooter>
      </Card>
    </div>
  </div>
</template>
