<script setup lang="ts">
import { Moon, Sun } from '@lucide/vue'
import { Button } from '@/components/ui/button'
import { useTheme } from '@/composables/useTheme'

const { isDark, toggleDark } = useTheme()
</script>

<template>
  <Button
    variant="outline"
    size="icon"
    :aria-label="isDark ? 'Switch to light mode' : 'Switch to dark mode'"
    @click="toggleDark()"
  >
    <Sun v-if="isDark" class="size-4" />
    <Moon v-else class="size-4" />
  </Button>
</template>
