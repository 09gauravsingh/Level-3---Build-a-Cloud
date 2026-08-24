<script setup>
import { computed } from 'vue'

const props = defineProps({
  // CloudNativePG cluster phase, e.g. "Cluster in healthy state".
  status: { type: String, default: '' },
})

const tone = computed(() => {
  const status = props.status.toLowerCase()

  if (!status) return 'unknown'
  if (status.includes('healthy') || status.includes('ready')) return 'healthy'
  if (status.includes('fail') || status.includes('error')) return 'failed'
  if (status.includes('delet')) return 'deleting'

  return 'pending'
})

const styles = {
  healthy: 'bg-emerald-50 text-emerald-800 ring-emerald-700/15',
  failed: 'bg-rose-50 text-rose-800 ring-rose-600/15',
  deleting: 'bg-neutral-100 text-neutral-600 ring-neutral-400/20',
  pending: 'bg-amber-50 text-amber-800 ring-amber-600/15',
  unknown: 'bg-neutral-100 text-neutral-500 ring-neutral-400/20',
}

const dotStyles = {
  healthy: 'bg-emerald-600',
  failed: 'bg-rose-500',
  deleting: 'bg-neutral-400',
  pending: 'bg-amber-500 animate-pulse',
  unknown: 'bg-neutral-400',
}
</script>

<template>
  <span
    class="inline-flex items-center gap-1.5 rounded-full px-2.5 py-1 text-[11px] font-semibold tracking-wide ring-1 ring-inset"
    :class="styles[tone]"
  >
    <span class="size-1.5 rounded-full" :class="dotStyles[tone]" />
    {{ status || 'Unknown' }}
  </span>
</template>
