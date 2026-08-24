<script setup>
import { computed } from 'vue'
import BaseButton from '@/components/ui/BaseButton.vue'
import StatusBadge from '@/components/ui/StatusBadge.vue'

const props = defineProps({
  instance: { type: Object, required: true },
  busy: { type: Boolean, default: false },
})

const emit = defineEmits(['connection', 'delete'])

const ready = computed(() => props.instance.readyInstances ?? 0)
const desired = computed(() => props.instance.desiredInstances ?? 0)

const createdDate = computed(() => {
  if (!props.instance.createdAt) return null

  const date = new Date(props.instance.createdAt)

  return Number.isNaN(date.getTime()) ? null : date
})

const createdAbsolute = computed(() =>
  createdDate.value
    ? createdDate.value.toLocaleString(undefined, { dateStyle: 'medium', timeStyle: 'short' })
    : props.instance.createdAt || '—',
)

// Ages are short-lived in a demo platform, so a relative label reads better.
const createdAt = computed(() => {
  if (!createdDate.value) return createdAbsolute.value

  const minutes = Math.round((Date.now() - createdDate.value.getTime()) / 60000)

  if (minutes < 1) return 'just now'
  if (minutes < 60) return `${minutes}m ago`
  if (minutes < 1440) return `${Math.round(minutes / 60)}h ago`

  return `${Math.round(minutes / 1440)}d ago`
})
</script>

<template>
  <tr class="border-b border-neutral-100 last:border-b-0 hover:bg-neutral-50/80">
    <td class="px-4 py-3">
      <p class="font-mono text-sm font-medium text-neutral-900">{{ instance.name }}</p>
      <p class="mt-0.5 text-xs text-neutral-500">{{ instance.namespace || '—' }}</p>
    </td>
    <td class="px-4 py-3">
      <StatusBadge :status="instance.status" />
    </td>
    <td class="px-4 py-3 font-mono text-sm text-neutral-700">{{ ready }} / {{ desired }}</td>
    <td class="hidden px-4 py-3 text-sm text-neutral-700 md:table-cell">
      {{ instance.storageSize || '—' }}
    </td>
    <td class="hidden px-4 py-3 font-mono text-xs text-neutral-600 lg:table-cell">
      {{ instance.primary || '—' }}
    </td>
    <td class="hidden px-4 py-3 font-mono text-xs text-neutral-600 xl:table-cell">
      {{ instance.ownedBy || '—' }}
    </td>
    <td class="hidden px-4 py-3 text-sm text-neutral-600 xl:table-cell" :title="createdAbsolute">
      {{ createdAt }}
    </td>
    <td class="px-4 py-3">
      <div class="flex justify-end gap-2">
        <BaseButton
          variant="secondary"
          size="sm"
          :loading="busy"
          @click="emit('connection', instance.name)"
        >
          Connection details
        </BaseButton>
        <BaseButton variant="danger" size="sm" :disabled="busy" @click="emit('delete', instance.name)">
          Delete
        </BaseButton>
      </div>
    </td>
  </tr>
</template>
