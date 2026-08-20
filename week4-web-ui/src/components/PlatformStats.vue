<script setup>
import { computed } from 'vue'

const props = defineProps({
  instances: { type: Array, default: () => [] },
})

const ICONS = {
  database:
    'M12 2c-4.42 0-8 1.34-8 3v14c0 1.66 3.58 3 8 3s8-1.34 8-3V5c0-1.66-3.58-3-8-3Zm0 2c3.87 0 6 1.1 6 1s-2.13 1-6 1-6-1.1-6-1 2.13-1 6-1Z',
  check: 'M9.55 17.6 4 12.05l1.4-1.4 4.15 4.14L18.6 5.4 20 6.8 9.55 17.6Z',
  clock: 'M12 2a10 10 0 1 0 0 20 10 10 0 0 0 0-20Zm1 10.6V6h-2v7.4l5 3 1-1.73-4-2.07Z',
  layers: 'M12 3 2 8l10 5 10-5-10-5Zm0 12L4.24 11.1 2 12l10 5 10-5-2.24-.9L12 15Z',
}

const stats = computed(() => {
  const healthy = props.instances.filter((instance) =>
    (instance.status || '').toLowerCase().includes('healthy'),
  ).length

  const readyReplicas = props.instances.reduce(
    (total, instance) => total + (instance.readyInstances ?? 0),
    0,
  )

  const desiredReplicas = props.instances.reduce(
    (total, instance) => total + (instance.desiredInstances ?? 0),
    0,
  )

  return [
    {
      label: 'Instances',
      value: props.instances.length,
      caption: 'Managed PostgreSQL clusters',
      icon: ICONS.database,
      tone: 'bg-indigo-50 text-indigo-600',
    },
    {
      label: 'Healthy',
      value: healthy,
      caption: 'Reporting a healthy state',
      icon: ICONS.check,
      tone: 'bg-emerald-50 text-emerald-600',
    },
    {
      label: 'Settling',
      value: props.instances.length - healthy,
      caption: 'Provisioning or degraded',
      icon: ICONS.clock,
      tone: 'bg-amber-50 text-amber-600',
    },
    {
      label: 'Replicas ready',
      value: `${readyReplicas}/${desiredReplicas}`,
      caption: 'Across every cluster',
      icon: ICONS.layers,
      tone: 'bg-sky-50 text-sky-600',
    },
  ]
})
</script>

<template>
  <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
    <div
      v-for="stat in stats"
      :key="stat.label"
      class="rounded-2xl border border-slate-200 bg-white p-4 shadow-sm shadow-slate-900/5 transition hover:border-slate-300 hover:shadow-md"
    >
      <div class="flex items-center justify-between gap-3">
        <p class="text-xs font-medium tracking-wide text-slate-500 uppercase">{{ stat.label }}</p>

        <span
          class="flex size-8 items-center justify-center rounded-lg"
          :class="stat.tone"
          aria-hidden="true"
        >
          <svg class="size-4" viewBox="0 0 24 24" fill="currentColor">
            <path :d="stat.icon" />
          </svg>
        </span>
      </div>

      <p class="mt-3 text-2xl font-semibold tracking-tight text-slate-900">{{ stat.value }}</p>
      <p class="mt-1 text-xs text-slate-500">{{ stat.caption }}</p>
    </div>
  </div>
</template>
