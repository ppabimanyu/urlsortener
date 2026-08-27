<script setup lang="ts">
import type { Component } from 'vue'
import { Card, CardContent } from '@/components/ui/card'

defineProps<{
  title: string
  value: string | number
  description?: string
  icon?: Component
  iconColor?: string
}>()
</script>

<template>
  <Card class="relative overflow-hidden transition-all duration-200 hover:shadow-md border border-border/60 bg-card">
    <CardContent class="p-6">
      <div class="flex items-center justify-between space-y-0">
        <p class="text-sm font-medium text-muted-foreground">{{ title }}</p>
        <div
          v-if="icon"
          class="flex h-10 w-10 items-center justify-center rounded-xl bg-secondary/80 text-foreground transition-transform hover:scale-105"
          :class="iconColor"
        >
          <component :is="icon" class="h-5 w-5" />
        </div>
      </div>
      <div class="mt-3">
        <div class="text-3xl font-bold tracking-tight text-foreground">{{ value }}</div>
        <p v-if="description" class="mt-1 text-xs text-muted-foreground">
          {{ description }}
        </p>
      </div>
    </CardContent>
  </Card>
</template>
