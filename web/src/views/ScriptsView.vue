<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { toast } from 'vue-sonner'
import { Button } from '@/components/ui/button'
import { api, setToken, type Script } from '@/services/api'
import AddScriptDialog from '@/components/AddScriptDialog.vue'
import ScriptCard from '@/components/ScriptCard.vue'
import ThemeToggle from '@/components/ThemeToggle.vue'

const router = useRouter()
const scripts = ref<Script[]>([])
const loading = ref(true)
const addOpen = ref(false)

async function load() {
  loading.value = true
  try {
    scripts.value = await api.listScripts()
  } catch (e) {
    toast.error((e as Error).message)
  } finally {
    loading.value = false
  }
}

async function logout() {
  try {
    await api.logout()
  } finally {
    setToken(null)
    router.push('/')
  }
}

onMounted(load)
</script>

<template>
  <div class="mx-auto min-h-screen max-w-3xl p-6">
    <header class="mb-6 flex items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-semibold">Install Scripts</h1>
        <p class="text-muted-foreground text-sm">
          Share install commands with a single curl.
        </p>
      </div>
      <div class="flex gap-2">
        <ThemeToggle />
        <Button @click="addOpen = true">
          Add script
        </Button>
        <Button variant="outline" @click="logout">
          Logout
        </Button>
      </div>
    </header>

    <div v-if="loading" class="text-muted-foreground py-16 text-center">
      Loading…
    </div>
    <div
      v-else-if="scripts.length === 0"
      class="text-muted-foreground py-16 text-center"
    >
      No scripts yet. Click “Add script” to create one.
    </div>
    <div v-else class="space-y-4">
      <ScriptCard
        v-for="script in scripts"
        :key="script.id"
        :script="script"
        @deleted="load"
      />
    </div>

    <AddScriptDialog v-model:open="addOpen" @created="load" />
  </div>
</template>
