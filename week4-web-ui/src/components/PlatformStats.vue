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
    { label: 'Instances', value: props.instances.length, accent: 'text-sky-300' },
    { label: 'Healthy', value: healthy, accent: 'text-emerald-300' },
    {
      label: 'Replicas ready',
      value: `${readyReplicas}/${desiredReplicas}`,
      accent: 'text-slate-100',
    },
  ]
})
</script>

<template>
  <div class="grid grid-cols-3 gap-3">
    <div
      v-for="stat in stats"
      :key="stat.label"
      class="rounded-xl border border-white/10 bg-slate-900/50 px-4 py-3 backdrop-blur"
    >
      <p class="text-[11px] tracking-wide text-slate-500 uppercase">{{ stat.label }}</p>
      <p class="mt-1 text-xl font-semibold" :class="stat.accent">{{ stat.value }}</p>
    </div>
  </div>
</template>
