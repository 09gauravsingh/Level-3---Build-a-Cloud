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
  database: '',
  owner: '',
})

const form = reactive(emptyForm())
const errors = ref({})

function validate() {
  const found = {}

  if (!form.name || !VALID_NAME.test(form.name)) {
    found.name = 'Use lowercase letters, numbers or hyphens'
  }

  if (!Number.isInteger(form.instances) || form.instances < 1 || form.instances > 3) {
    found.instances = 'Instances must be between 1 and 3'
  }

  errors.value = found

  return Object.keys(found).length === 0
}

function submit() {
  if (!validate()) {
    return
  }

  emit('submit', { ...form })
}

function reset() {
  Object.assign(form, emptyForm())
  errors.value = {}
}

defineExpose({ reset })
</script>

<template>
  <AppCard title="Create instance" subtitle="Provisioned through CloudNativePG on Kubernetes">
    <form class="space-y-5" @submit.prevent="submit">
      <div class="grid gap-4 sm:grid-cols-2">
        <BaseInput
          v-model="form.name"
          label="Instance name"
          placeholder="week4-demo-db"
          hint="Lowercase letters, numbers and hyphens"
          :error="errors.name"
          required
        />

        <BaseInput
          v-model="form.instances"
          label="Replicas"
          type="number"
          :min="1"
          :max="3"
          hint="Between 1 and 3"
          :error="errors.instances"
        />

        <BaseInput
          v-model="form.storageSize"
          label="Storage size"
          placeholder="1Gi"
          hint="Kubernetes quantity, defaults to 1Gi"
        />

        <BaseInput
          v-model="form.database"
          label="Database name"
          placeholder="app"
          hint="Defaults to app"
        />

        <div class="sm:col-span-2">
          <BaseInput
            v-model="form.owner"
            label="Database owner"
            placeholder="app"
            hint="Defaults to the database name"
          />
        </div>
      </div>

      <div class="flex items-center justify-between gap-3 border-t border-white/10 pt-4">
        <p class="text-xs text-slate-500">Provisioning continues in the background.</p>

        <div class="flex gap-2">
          <BaseButton variant="ghost" size="md" :disabled="loading" @click="reset">Reset</BaseButton>
          <BaseButton type="submit" :loading="loading">Create instance</BaseButton>
        </div>
      </div>
    </form>
  </AppCard>
</template>
