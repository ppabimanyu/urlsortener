<script setup lang="ts">
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import {
  Link2,
  LayoutDashboard,
  Link as LinkIcon,
  LogOut,
  User as UserIcon,
  Plus,
} from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import ThemeToggle from './ThemeToggle.vue'
import CreateEditUrlDialog from '@/components/dialogs/CreateEditUrlDialog.vue'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const isCreateDialogOpen = ref(false)

function handleLogout() {
  authStore.logout()
  router.push('/login')
}
</script>

<template>
  <header class="sticky top-0 z-40 w-full border-b border-border/40 bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
    <div class="container mx-auto flex h-16 max-w-7xl items-center justify-between px-4 sm:px-6 lg:px-8">
      <!-- Left: Logo & Nav -->
      <div class="flex items-center gap-6 md:gap-8">
        <router-link to="/dashboard" class="flex items-center gap-2.5 group">
          <div class="flex h-9 w-9 items-center justify-center rounded-xl bg-primary text-primary-foreground shadow-sm transition-transform duration-200 group-hover:scale-105">
            <Link2 class="h-5 w-5" />
          </div>
          <div class="flex flex-col">
            <span class="text-lg font-bold tracking-tight bg-gradient-to-r from-foreground via-foreground/90 to-foreground/70 bg-clip-text">
              SortLink
            </span>
          </div>
        </router-link>

        <nav class="hidden md:flex items-center gap-1.5 text-sm font-medium">
          <router-link
            to="/dashboard"
            class="flex items-center gap-2 px-3 py-2 rounded-lg transition-colors"
            :class="route.path === '/dashboard' ? 'bg-secondary text-foreground font-semibold shadow-xs' : 'text-muted-foreground hover:text-foreground hover:bg-muted/50'"
          >
            <LayoutDashboard class="h-4 w-4" />
            Dashboard
          </router-link>
          <router-link
            to="/links"
            class="flex items-center gap-2 px-3 py-2 rounded-lg transition-colors"
            :class="route.path.startsWith('/links') || route.path.startsWith('/analytics') ? 'bg-secondary text-foreground font-semibold shadow-xs' : 'text-muted-foreground hover:text-foreground hover:bg-muted/50'"
          >
            <LinkIcon class="h-4 w-4" />
            My Links
          </router-link>
        </nav>
      </div>

      <!-- Right: Action Buttons & Profile -->
      <div class="flex items-center gap-2.5 sm:gap-3">
        <!-- New Link Button -->
        <Button
          size="sm"
          class="gap-1.5 shadow-sm font-medium transition-transform hover:scale-[1.02]"
          @click="isCreateDialogOpen = true"
        >
          <Plus class="h-4 w-4" />
          <span class="hidden sm:inline">Create Short Link</span>
          <span class="sm:hidden">New</span>
        </Button>

        <!-- Theme Toggle -->
        <ThemeToggle />

        <!-- User Dropdown -->
        <DropdownMenu v-if="authStore.user">
          <DropdownMenuTrigger as-child>
            <Button variant="ghost" size="sm" class="relative gap-2 px-2.5 rounded-full border border-border/50 hover:bg-muted">
              <div class="flex h-7 w-7 items-center justify-center rounded-full bg-primary/10 text-primary font-bold text-xs uppercase">
                {{ authStore.user.name.charAt(0) }}
              </div>
              <span class="hidden md:inline-block max-w-[120px] truncate text-xs font-medium">
                {{ authStore.user.name }}
              </span>
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" class="w-56">
            <DropdownMenuLabel>
              <div class="flex flex-col space-y-1">
                <p class="text-sm font-medium leading-none">{{ authStore.user.name }}</p>
                <p class="text-xs leading-none text-muted-foreground truncate">{{ authStore.user.email }}</p>
              </div>
            </DropdownMenuLabel>
            <DropdownMenuSeparator />
            <DropdownMenuItem @click="router.push('/dashboard')">
              <LayoutDashboard class="mr-2 h-4 w-4" />
              <span>Dashboard</span>
            </DropdownMenuItem>
            <DropdownMenuItem @click="router.push('/links')">
              <LinkIcon class="mr-2 h-4 w-4" />
              <span>All Links</span>
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem class="text-destructive focus:text-destructive focus:bg-destructive/10" @click="handleLogout">
              <LogOut class="mr-2 h-4 w-4" />
              <span>Log out</span>
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </div>

    <!-- Create URL Dialog -->
    <CreateEditUrlDialog
      v-model:open="isCreateDialogOpen"
      @created="() => {}"
    />
  </header>
</template>
