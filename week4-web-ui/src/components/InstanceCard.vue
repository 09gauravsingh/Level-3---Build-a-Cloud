<script setup>
import { computed } from 'vue'
import BaseButton from '@/components/ui/BaseButton.vue'
import StatusBadge from '@/components/ui/StatusBadge.vue'

const props = defineProps({
  instance: { type: Object, required: true },
  busy: { type: Boolean, default: false },
})

const emit = defineEmits(['connection', 'delete'])

const replicas = computed(
  () => `${props.instance.readyInstances ?? 0} / ${props.instance.desiredInstances ?? 0}`,
)

const createdAt = computed(() => {
  if (!props.instance.createdAt) return '—'

  const date = new Date(props.instance.createdAt)

  return Number.isNaN(date.getTime())
    ? props.instance.createdAt
    : date.toLocaleString(undefined, { dateStyle: 'medium', timeStyle: 'short' })
})

const details = computed(() => [
  { label: 'Database', value: props.instance.database || '—' },
  { label: 'Owner', value: props.instance.owner || '—' },
  { label: 'Storage', value: props.instance.storageSize || '—' },
  { label: 'Replicas ready', value: replicas.value },
  { label: 'Primary', value: props.instance.primary || '—' },
  { label: 'Created', value: createdAt.value },
])
</script>

<template>
  <article
    class="rounded-xl border border-white/10 bg-slate-950/40 p-4 transition hover:border-sky-500/30 hover:bg-slate-950/70"
  >
    <div class="flex flex-wrap items-start justify-between gap-3">
      <div class="min-w-0">
        <h3 class="truncate font-mono text-sm font-semibold text-slate-100">
          {{ instance.name }}
        </h3>
        <p class="mt-0.5 text-xs text-slate-500">namespace: {{ instance.namespace || '—' }}</p>
      </div>

      <StatusBadge :status="instance.status" />
    </div>

    <dl class="mt-4 grid grid-cols-2 gap-x-4 gap-y-3 sm:grid-cols-3">
      <div v-for="detail in details" :key="detail.label" class="min-w-0">
        <dt class="text-[11px] tracking-wide text-slate-500 uppercase">{{ detail.label }}</dt>
        <dd class="truncate text-sm text-slate-200" :title="String(detail.value)">
          {{ detail.value }}
        </dd>
      </div>
    </dl>

    <div class="mt-4 flex justify-end gap-2 border-t border-white/5 pt-3">
      <BaseButton
        variant="secondary"
        size="sm"
        :disabled="busy"
        @click="emit('connection', instance.name)"
      >
        Connection details
      </BaseButton>

      <BaseButton variant="danger" size="sm" :disabled="busy" @click="emit('delete', instance.name)">
        Delete
      </BaseButton>
    </div>
  </article>
</template>
