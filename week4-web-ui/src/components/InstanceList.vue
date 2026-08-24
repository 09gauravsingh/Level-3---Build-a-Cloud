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
    <template #actions>
      <div class="relative">
        <svg
          class="pointer-events-none absolute top-1/2 left-2.5 size-4 -translate-y-1/2 text-neutral-400"
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
          class="h-8 w-40 rounded-md border border-neutral-200 bg-white pr-3 pl-8 text-xs text-neutral-800 transition placeholder:text-neutral-400 hover:border-neutral-300 focus:border-neutral-900 focus:ring-2 focus:ring-neutral-900/10 focus:outline-none sm:w-52"
        >
      </div>

      <BaseButton variant="secondary" size="sm" :loading="loading" @click="emit('refresh')">
        Refresh
      </BaseButton>
    </template>

    <div v-if="instances.length > 0" class="mb-4 flex gap-1 border-b border-neutral-200">
      <button
        v-for="tab in tabs"
        :key="tab.id"
        type="button"
        class="-mb-px inline-flex items-center gap-1.5 border-b-2 px-3 py-2 text-sm transition"
        :class="
          statusFilter === tab.id
            ? 'border-neutral-900 font-semibold text-neutral-900'
            : 'border-transparent font-medium text-neutral-500 hover:text-neutral-800'
        "
        :aria-pressed="statusFilter === tab.id"
        @click="statusFilter = tab.id"
      >
        {{ tab.label }}
        <span class="text-xs tabular-nums text-neutral-400">{{ tab.count }}</span>
      </button>
    </div>

    <div v-if="loading && instances.length === 0" class="space-y-3">
      <div v-for="n in 2" :key="n" class="h-12 animate-pulse rounded-md bg-neutral-100" />
    </div>

    <div v-else-if="instances.length === 0" class="py-14 text-center">
      <p class="text-sm font-medium text-neutral-800">No instances yet</p>
      <p class="mt-1 text-sm text-neutral-500">
        Use the form to provision your first PostgreSQL cluster. It appears here within a few
        seconds.
      </p>
    </div>

    <div v-else-if="filtered.length === 0" class="py-10 text-center">
      <p class="text-sm text-neutral-500">Nothing matches the current filters.</p>
      <button
        class="mt-2 text-sm font-medium text-neutral-900 underline underline-offset-2"
        @click="clearFilters"
      >
        Clear filters
      </button>
    </div>

    <div v-else class="overflow-x-auto">
      <table class="min-w-full text-left">
        <thead
          class="border-b border-neutral-200 bg-neutral-50 text-[11px] font-semibold tracking-wide text-neutral-500 uppercase"
        >
          <tr>
            <th class="px-4 py-2.5">Instance</th>
            <th class="px-4 py-2.5">Status</th>
            <th class="px-4 py-2.5">Replicas</th>
            <th class="hidden px-4 py-2.5 md:table-cell">Storage</th>
            <th class="hidden px-4 py-2.5 lg:table-cell">Primary</th>
            <th class="hidden px-4 py-2.5 xl:table-cell">Owner</th>
            <th class="hidden px-4 py-2.5 xl:table-cell">Created</th>
            <th class="px-4 py-2.5 text-right">Actions</th>
          </tr>
        </thead>
        <TransitionGroup
          tag="tbody"
          move-class="transition duration-300"
          enter-active-class="transition duration-300 ease-out"
          enter-from-class="translate-y-1 opacity-0"
          leave-active-class="transition duration-200 ease-in"
          leave-to-class="opacity-0"
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
      </table>
    </div>
  </AppCard>
</template>
