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

const replicaProgress = computed(() =>
  desired.value > 0 ? Math.min(100, Math.round((ready.value / desired.value) * 100)) : 0,
)

const replicasComplete = computed(() => desired.value > 0 && ready.value >= desired.value)

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

const details = computed(() => {
  const items = []

  if (props.instance.ownedBy) {
    items.push({ label: 'Owned by', value: props.instance.ownedBy, mono: true })
  }

  items.push(
    { label: 'Storage', value: props.instance.storageSize || '—' },
    { label: 'Primary', value: props.instance.primary || '—', mono: true },
    { label: 'Created', value: createdAt.value, title: createdAbsolute.value },
  )

  return items
})
</script>

<template>
  <article
    class="group rounded-xl border border-slate-200 bg-white p-4 transition hover:border-indigo-300 hover:shadow-md hover:shadow-slate-900/5"
  >
    <div class="flex flex-wrap items-start justify-between gap-3">
      <div class="flex min-w-0 items-start gap-3">
        <span
          class="flex size-9 shrink-0 items-center justify-center rounded-lg bg-slate-100 text-slate-500 transition group-hover:bg-indigo-50 group-hover:text-indigo-600"
          aria-hidden="true"
        >
          <svg class="size-4" viewBox="0 0 24 24" fill="currentColor">
            <path
              d="M12 2c-4.42 0-8 1.34-8 3v14c0 1.66 3.58 3 8 3s8-1.34 8-3V5c0-1.66-3.58-3-8-3Zm0 2c3.87 0 6 1.1 6 1s-2.13 1-6 1-6-1.1-6-1 2.13-1 6-1Z"
            />
          </svg>
        </span>

        <div class="min-w-0">
          <h3 class="truncate font-mono text-sm font-semibold text-slate-900">
            {{ instance.name }}
          </h3>
          <p class="mt-0.5 truncate text-xs text-slate-500">
            namespace: {{ instance.namespace || '—' }}
          </p>
        </div>
      </div>

      <StatusBadge :status="instance.status" />
    </div>

    <div class="mt-4 rounded-lg bg-slate-50 px-3 py-2.5">
      <div class="flex items-center justify-between text-xs">
        <span class="font-medium text-slate-600">Replicas ready</span>
        <span class="font-mono text-slate-900">{{ ready }} / {{ desired }}</span>
      </div>

      <div class="mt-2 h-1.5 overflow-hidden rounded-full bg-slate-200">
        <div
          class="h-full rounded-full transition-all duration-500"
          :class="replicasComplete ? 'bg-emerald-500' : 'bg-amber-500'"
          :style="{ width: `${replicaProgress}%` }"
        />
      </div>
    </div>

    <dl class="mt-4 grid grid-cols-2 gap-x-4 gap-y-3 sm:grid-cols-3">
      <div v-for="detail in details" :key="detail.label" class="min-w-0">
        <dt class="text-[11px] tracking-wide text-slate-400 uppercase">{{ detail.label }}</dt>
        <dd
          class="truncate text-sm text-slate-700"
          :class="detail.mono ? 'font-mono text-xs' : ''"
          :title="detail.title || String(detail.value)"
        >
          {{ detail.value }}
        </dd>
      </div>
    </dl>

    <div class="mt-4 flex flex-wrap justify-end gap-2 border-t border-slate-100 pt-3">
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
  </article>
</template>
