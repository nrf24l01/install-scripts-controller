<script setup lang="ts">
import { computed, ref } from 'vue'
import { toast } from 'vue-sonner'
import { Copy, Eye, EyeOff, Trash2 } from '@lucide/vue'
import { Button } from '@/components/ui/button'
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import hljs from '@/lib/highlight'
import { api, type Script } from '@/services/api'

const props = defineProps<{ script: Script }>()
const emit = defineEmits<{ deleted: [] }>()

const body = ref<string | null>(null)
const showBody = ref(false)
const loadingBody = ref(false)
const deleting = ref(false)

const highlightedBody = computed(() => {
  const code = body.value ?? ''
  return code ? hljs.highlightAuto(code).value : ''
})

function installCommand(): string {
  return `curl -sSL "${props.script.install_url}" | bash`
}

async function copy() {
  try {
    await navigator.clipboard.writeText(installCommand())
    toast.success('Install command copied')
  } catch {
    toast.error('Failed to copy')
  }
}

async function toggleBody() {
  if (showBody.value) {
    showBody.value = false
    return
  }
  if (body.value === null) {
    loadingBody.value = true
    try {
      const full = await api.getScript(props.script.id)
      body.value = full.script ?? ''
    } catch (e) {
      toast.error((e as Error).message)
      return
    } finally {
      loadingBody.value = false
    }
  }
  showBody.value = true
}

async function remove() {
  if (!window.confirm(`Delete script "${props.script.name}"?`)) {
    return
  }
  deleting.value = true
  try {
    await api.deleteScript(props.script.id)
    toast.success('Script deleted')
    emit('deleted')
  } catch (e) {
    toast.error((e as Error).message)
  } finally {
    deleting.value = false
  }
}
</script>

<template>
  <Card>
    <CardHeader>
      <div class="flex items-start justify-between gap-4">
        <div class="min-w-0">
          <CardTitle>{{ script.name }}</CardTitle>
          <CardDescription v-if="script.description" class="mt-1">
            {{ script.description }}
          </CardDescription>
        </div>
        <Button
          variant="destructive"
          size="sm"
          :disabled="deleting"
          @click="remove"
        >
          <Trash2 class="size-4" />
        </Button>
      </div>
    </CardHeader>
    <CardContent class="space-y-4">
      <div class="space-y-2">
        <Label>Install command</Label>
        <div class="flex gap-2">
          <Input
            :model-value="installCommand()"
            readonly
            class="flex-1 font-mono text-xs"
            aria-label="Install command"
          />
          <Button variant="secondary" @click="copy">
            <Copy class="size-4" />
            Copy
          </Button>
        </div>
      </div>

      <div>
        <Button variant="ghost" size="sm" :disabled="loadingBody" @click="toggleBody">
          <EyeOff v-if="showBody" class="size-4" />
          <Eye v-else class="size-4" />
          {{ showBody ? 'Hide' : 'View' }} script
        </Button>
        <pre
          v-if="showBody"
          class="bg-muted mt-2 max-h-96 overflow-auto rounded-md p-4 font-mono text-xs"
          v-html="highlightedBody"
        ></pre>
      </div>
    </CardContent>
  </Card>
</template>
