<script setup>
import { computed } from 'vue'

const props = defineProps({
  instances: { type: Array, default: () => [] },
})

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
    },
    {
      label: 'Healthy',
      value: healthy,
      caption: 'Reporting a healthy state',
    },
    {
      label: 'Settling',
      value: props.instances.length - healthy,
      caption: 'Provisioning or degraded',
    },
    {
      label: 'Replicas ready',
      value: `${readyReplicas}/${desiredReplicas}`,
      caption: 'Across every cluster',
    },
  ]
})
</script>

<template>
  <div class="grid gap-px overflow-hidden rounded-lg border border-neutral-200 bg-neutral-200 sm:grid-cols-2 lg:grid-cols-4">
    <div
      v-for="stat in stats"
      :key="stat.label"
      class="bg-white px-5 py-4"
    >
      <p class="text-xs font-medium text-neutral-500">{{ stat.label }}</p>
      <p class="mt-2 text-2xl font-semibold tracking-tight text-neutral-900">{{ stat.value }}</p>
      <p class="mt-1 text-xs text-neutral-500">{{ stat.caption }}</p>
    </div>
  </div>
</template>
