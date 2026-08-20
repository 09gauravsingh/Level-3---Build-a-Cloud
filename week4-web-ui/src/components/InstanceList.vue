<script setup>
import { computed, ref } from 'vue'
import AppCard from '@/components/ui/AppCard.vue'
import BaseButton from '@/components/ui/BaseButton.vue'
import InstanceCard from '@/components/InstanceCard.vue'

const props = defineProps({
  instances: { type: Array, default: () => [] },
  loading: { type: Boolean, default: false },
  busyName: { type: String, default: '' },
})

const emit = defineEmits(['refresh', 'connection', 'delete'])

const search = ref('')
const statusFilter = ref('all')

function isHealthy(instance) {
  return (instance.status || '').toLowerCase().includes('healthy')
}

const tabs = computed(() => [
  { id: 'all', label: 'All', count: props.instances.length },
  { id: 'healthy', label: 'Healthy', count: props.instances.filter(isHealthy).length },
  {
    id: 'settling',
    label: 'Settling',
    count: props.instances.filter((instance) => !isHealthy(instance)).length,
  },
])

const filtered = computed(() => {
  const term = search.value.trim().toLowerCase()

  return props.instances.filter((instance) => {
    if (statusFilter.value === 'healthy' && !isHealthy(instance)) return false
    if (statusFilter.value === 'settling' && isHealthy(instance)) return false

    if (!term) return true

    return [instance.name, instance.ownedBy, instance.status, instance.storageSize]
      .filter(Boolean)
      .some((field) => String(field).toLowerCase().includes(term))
  })
})

function clearFilters() {
  search.value = ''
  statusFilter.value = 'all'
}
</script>

<template>
  <AppCard
    title="Instances"
    :subtitle="`${instances.length} instance${instances.length === 1 ? '' : 's'} in the platform`"
  >
    <template #icon>
      <svg class="size-4" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true">
        <path
          d="M4 5c0-1.66 3.58-3 8-3s8 1.34 8 3-3.58 3-8 3-8-1.34-8-3Zm16 4.35V12c0 1.66-3.58 3-8 3s-8-1.34-8-3V9.35C5.5 10.5 8.6 11 12 11s6.5-.5 8-1.65Zm0 5V17c0 1.66-3.58 3-8 3s-8-1.34-8-3v-2.65C5.5 15.5 8.6 16 12 16s6.5-.5 8-1.65Z"
        />
      </svg>
    </template>

    <template #actions>
      <div class="relative">
        <svg
          class="pointer-events-none absolute top-1/2 left-2.5 size-4 -translate-y-1/2 text-slate-400"
          viewBox="0 0 20 20"
          fill="currentColor"
          aria-hidden="true"
        >
          <path
            fill-rule="evenodd"
            d="M9 3.5a5.5 5.5 0 1 0 3.4 9.83l3.13 3.13a.75.75 0 1 0 1.06-1.06l-3.13-3.13A5.5 5.5 0 0 0 9 3.5ZM5 9a4 4 0 1 1 8 0 4 4 0 0 1-8 0Z"
            clip-rule="evenodd"
          />
        </svg>

        <input
          v-model="search"
          type="search"
          placeholder="Search instances"
          aria-label="Search instances"
          class="h-8 w-40 rounded-lg border border-slate-200 bg-white pr-3 pl-8 text-xs text-slate-800 shadow-sm transition placeholder:text-slate-400 hover:border-slate-300 focus:border-indigo-500 focus:ring-4 focus:ring-indigo-500/10 focus:outline-none sm:w-52"
        >
      </div>

      <BaseButton variant="secondary" size="sm" :loading="loading" @click="emit('refresh')">
        Refresh
      </BaseButton>
    </template>

    <div v-if="instances.length > 0" class="mb-4 flex flex-wrap gap-1.5">
      <button
        v-for="tab in tabs"
        :key="tab.id"
        type="button"
        class="inline-flex items-center gap-1.5 rounded-full px-3 py-1 text-xs font-medium transition"
        :class="
          statusFilter === tab.id
            ? 'bg-slate-900 text-white'
            : 'bg-slate-100 text-slate-600 hover:bg-slate-200'
        "
        :aria-pressed="statusFilter === tab.id"
        @click="statusFilter = tab.id"
      >
        {{ tab.label }}
        <span
          class="rounded-full px-1.5 text-[10px] tabular-nums"
          :class="statusFilter === tab.id ? 'bg-white/20' : 'bg-white'"
        >
          {{ tab.count }}
        </span>
      </button>
    </div>

    <div v-if="loading && instances.length === 0" class="space-y-3">
      <div v-for="n in 2" :key="n" class="h-40 animate-pulse rounded-xl bg-slate-100" />
    </div>

    <div
      v-else-if="instances.length === 0"
      class="flex flex-col items-center justify-center rounded-xl border border-dashed border-slate-300 bg-slate-50/60 px-6 py-14 text-center"
    >
      <div class="flex size-12 items-center justify-center rounded-full bg-white shadow-sm">
        <svg
          class="size-5 text-indigo-500"
          viewBox="0 0 24 24"
          fill="currentColor"
          aria-hidden="true"
        >
          <path
            d="M12 2c-4.42 0-8 1.34-8 3v14c0 1.66 3.58 3 8 3s8-1.34 8-3V5c0-1.66-3.58-3-8-3Zm0 2c3.87 0 6 1.1 6 1s-2.13 1-6 1-6-1.1-6-1 2.13-1 6-1Z"
          />
        </svg>
      </div>

      <p class="mt-3 text-sm font-semibold text-slate-800">No instances yet</p>
      <p class="mt-1 max-w-xs text-xs text-slate-500">
        Use the form to provision your first PostgreSQL cluster. It appears here within a few
        seconds.
      </p>
    </div>

    <div
      v-else-if="filtered.length === 0"
      class="rounded-xl border border-dashed border-slate-300 bg-slate-50/60 py-10 text-center"
    >
      <p class="text-sm text-slate-500">Nothing matches the current filters.</p>

      <button
        class="mt-2 text-xs font-semibold text-indigo-600 transition hover:text-indigo-500"
        @click="clearFilters"
      >
        Clear filters
      </button>
    </div>

    <TransitionGroup
      v-else
      tag="div"
      class="grid gap-3"
      move-class="transition duration-300"
      enter-active-class="transition duration-300 ease-out"
      enter-from-class="translate-y-2 opacity-0"
      leave-active-class="absolute transition duration-200 ease-in"
      leave-to-class="scale-95 opacity-0"
    >
      <InstanceCard
        v-for="instance in filtered"
        :key="instance.name"
        :instance="instance"
        :busy="busyName === instance.name"
        @connection="emit('connection', $event)"
        @delete="emit('delete', $event)"
      />
    </TransitionGroup>
  </AppCard>
</template>
