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

const filtered = computed(() => {
  const term = search.value.trim().toLowerCase()

  if (!term) return props.instances

  return props.instances.filter((instance) =>
    [instance.name, instance.database, instance.owner, instance.status]
      .filter(Boolean)
      .some((field) => String(field).toLowerCase().includes(term)),
  )
})
</script>

<template>
  <AppCard
    title="Instances"
    :subtitle="`${instances.length} instance${instances.length === 1 ? '' : 's'} in the platform`"
  >
    <template #actions>
      <div class="relative">
        <svg
          class="pointer-events-none absolute top-1/2 left-2.5 size-4 -translate-y-1/2 text-slate-500"
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
          class="h-8 w-40 rounded-lg border border-white/10 bg-slate-950/60 pr-3 pl-8 text-xs text-slate-200 placeholder:text-slate-600 focus:border-sky-400/60 focus:ring-2 focus:ring-sky-500/25 focus:outline-none sm:w-52"
        >
      </div>

      <BaseButton variant="secondary" size="sm" :loading="loading" @click="emit('refresh')">
        Refresh
      </BaseButton>
    </template>

    <div v-if="loading && instances.length === 0" class="space-y-3">
      <div v-for="n in 2" :key="n" class="h-32 animate-pulse rounded-xl bg-white/5" />
    </div>

    <div
      v-else-if="instances.length === 0"
      class="flex flex-col items-center justify-center rounded-xl border border-dashed border-white/10 px-6 py-12 text-center"
    >
      <div class="flex size-11 items-center justify-center rounded-full bg-white/5">
        <svg class="size-5 text-slate-500" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true">
          <path
            d="M12 2c-4.42 0-8 1.34-8 3v14c0 1.66 3.58 3 8 3s8-1.34 8-3V5c0-1.66-3.58-3-8-3Zm0 2c3.87 0 6 1.1 6 1s-2.13 1-6 1-6-1.1-6-1 2.13-1 6-1Z"
          />
        </svg>
      </div>

      <p class="mt-3 text-sm font-medium text-slate-300">No instances yet</p>
      <p class="mt-1 text-xs text-slate-500">Create your first PostgreSQL instance to get started.</p>
    </div>

    <p v-else-if="filtered.length === 0" class="py-8 text-center text-sm text-slate-500">
      No instances match “{{ search }}”.
    </p>

    <div v-else class="grid gap-3">
      <InstanceCard
        v-for="instance in filtered"
        :key="instance.name"
        :instance="instance"
        :busy="busyName === instance.name"
        @connection="emit('connection', $event)"
        @delete="emit('delete', $event)"
      />
    </div>
  </AppCard>
</template>
