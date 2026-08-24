<script setup>
import { reactive, ref } from 'vue'
import AppCard from '@/components/ui/AppCard.vue'
import BaseButton from '@/components/ui/BaseButton.vue'
import BaseInput from '@/components/ui/BaseInput.vue'

defineProps({
  loading: { type: Boolean, default: false },
})

const emit = defineEmits(['submit'])

// Mirrors validateCreateRequest in the Go API so users get instant feedback.
const VALID_NAME = /^[a-z0-9]([-a-z0-9]*[a-z0-9])?$/

const emptyForm = () => ({
  name: '',
  instances: 1,
  storageSize: '1Gi',
})

const form = reactive(emptyForm())
const errors = ref({})

function validate() {
  const found = {}

  if (!form.name || !VALID_NAME.test(form.name)) {
    found.name = 'Use lowercase letters, numbers or hyphens'
  }

  if (!Number.isInteger(form.instances) || form.instances < 1 || form.instances > 3) {
    found.instances = 'Replicas must be between 1 and 3'
  }

  errors.value = found

  return Object.keys(found).length === 0
}

function submit() {
  if (!validate()) {
    return
  }

  // Database and owner stay empty so the API applies its existing defaults.
  emit('submit', {
    name: form.name,
    instances: form.instances,
    storageSize: form.storageSize,
    database: '',
    owner: '',
  })
}

function reset() {
  Object.assign(form, emptyForm())
  errors.value = {}
}

defineExpose({ reset })
</script>

<template>
  <AppCard title="Create instance" subtitle="Name, storage and replica count">
    <template #icon>
      <svg class="size-4" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true">
        <path d="M11 5h2v6h6v2h-6v6h-2v-6H5v-2h6V5Z" />
      </svg>
    </template>

    <form class="space-y-5" @submit.prevent="submit">
      <BaseInput
        v-model="form.name"
        label="Instance name"
        placeholder="week4-demo-db"
        hint="Lowercase letters, numbers and hyphens"
        :error="errors.name"
        required
      />

      <div class="grid gap-4 sm:grid-cols-2">
        <BaseInput
          v-model="form.storageSize"
          label="Storage size"
          placeholder="1Gi"
          hint="Kubernetes quantity per replica"
        />

        <BaseInput
          v-model="form.instances"
          label="Replicas"
          type="number"
          min="1"
          max="3"
          hint="Integer from 1 to 3"
          :error="errors.instances"
          required
        />
      </div>

      <div class="flex flex-wrap items-center justify-between gap-3 border-t border-neutral-100 pt-4">
        <p class="flex items-center gap-1.5 text-xs text-neutral-500">
          <svg
            class="size-3.5 text-neutral-400"
            viewBox="0 0 20 20"
            fill="currentColor"
            aria-hidden="true"
          >
            <path
              d="M10 18a8 8 0 1 0 0-16 8 8 0 0 0 0 16Zm.75-11.5a.75.75 0 1 1-1.5 0 .75.75 0 0 1 1.5 0ZM9.25 9a.75.75 0 0 1 1.5 0v5a.75.75 0 0 1-1.5 0V9Z"
            />
          </svg>
          Provisioning continues in the background.
        </p>

        <div class="flex gap-2">
          <BaseButton variant="ghost" size="md" :disabled="loading" @click="reset">Reset</BaseButton>
          <BaseButton type="submit" :loading="loading">Create instance</BaseButton>
        </div>
      </div>
    </form>
  </AppCard>
</template>
