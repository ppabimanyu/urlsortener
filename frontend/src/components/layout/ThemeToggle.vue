<script setup lang="ts">
import { useDark, useToggle } from '@vueuse/core'
import { Sun, Moon } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'

const isDark = useDark({
  selector: 'html',
  attribute: 'class',
  valueDark: 'dark',
  valueLight: '',
})
const toggleDark = useToggle(isDark)
</script>

<template>
  <Button
    variant="ghost"
    size="icon"
    class="relative rounded-full transition-colors hover:bg-accent/80"
    :title="isDark ? 'Beralih ke Light Mode' : 'Beralih ke Dark Mode'"
    @click="toggleDark()"
  >
    <Sun v-if="isDark" class="h-5 w-5 text-amber-400 transition-transform rotate-0 scale-100 duration-200" />
    <Moon v-else class="h-5 w-5 text-slate-700 transition-transform rotate-0 scale-100 duration-200" />
    <span class="sr-only">Toggle theme</span>
  </Button>
</template>
