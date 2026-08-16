<script setup lang="ts">
import { ref } from 'vue'
import { useVModel } from '@vueuse/core'
import { toast } from 'vue-sonner'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { api } from '@/services/api'

const props = defineProps<{ open: boolean }>()
const emit = defineEmits<{
  'update:open': [value: boolean]
  created: []
}>()

const open = useVModel(props, 'open', emit)

const name = ref('')
const description = ref('')
const script = ref('')
const submitting = ref(false)

async function submit() {
  if (!name.value.trim()) {
    toast.error('Name is required')
    return
  }
  if (!script.value.trim()) {
    toast.error('Script is required')
    return
  }
  submitting.value = true
  try {
    await api.createScript({
      name: name.value.trim(),
      description: description.value.trim(),
      script: script.value,
    })
    toast.success('Script created')
    name.value = ''
    description.value = ''
    script.value = ''
    open.value = false
    emit('created')
  } catch (e) {
    toast.error((e as Error).message)
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <Dialog v-model:open="open">
    <DialogContent class="max-w-lg max-h-[90vh] overflow-y-auto">
      <DialogHeader>
        <DialogTitle>Add script</DialogTitle>
        <DialogDescription>
          Paste the script text. The install link is generated automatically.
        </DialogDescription>
      </DialogHeader>
      <form class="space-y-4" @submit.prevent="submit">
        <div class="space-y-2">
          <Label for="script-name">Name</Label>
          <Input
            id="script-name"
            v-model="name"
            placeholder="e.g. Install Docker"
            :disabled="submitting"
          />
        </div>
        <div class="space-y-2">
          <Label for="script-description">Description</Label>
          <Input
            id="script-description"
            v-model="description"
            placeholder="Short description (optional)"
            :disabled="submitting"
          />
        </div>
        <div class="space-y-2">
          <Label for="script-body">Script</Label>
          <Textarea
            id="script-body"
            v-model="script"
            placeholder="#!/bin/bash&#10;..."
            class="min-h-48 max-h-[60vh] overflow-y-auto font-mono"
            :disabled="submitting"
          />
        </div>
        <DialogFooter>
          <Button type="submit" :disabled="submitting">
            {{ submitting ? 'Creating…' : 'Create' }}
          </Button>
        </DialogFooter>
      </form>
    </DialogContent>
  </Dialog>
</template>
