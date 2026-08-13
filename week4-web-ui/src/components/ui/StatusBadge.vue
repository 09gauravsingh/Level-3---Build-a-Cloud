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
  healthy: 'bg-emerald-500/10 text-emerald-300 ring-emerald-500/30',
  failed: 'bg-rose-500/10 text-rose-300 ring-rose-500/30',
  deleting: 'bg-slate-500/10 text-slate-300 ring-slate-500/30',
  pending: 'bg-amber-500/10 text-amber-300 ring-amber-500/30',
  unknown: 'bg-slate-500/10 text-slate-400 ring-slate-500/30',
}

const dotStyles = {
  healthy: 'bg-emerald-400',
  failed: 'bg-rose-400',
  deleting: 'bg-slate-400',
  pending: 'bg-amber-400 animate-pulse',
  unknown: 'bg-slate-500',
}
</script>

<template>
  <span
    class="inline-flex items-center gap-1.5 rounded-full px-2.5 py-1 text-xs font-medium ring-1 ring-inset"
    :class="styles[tone]"
  >
    <span class="size-1.5 rounded-full" :class="dotStyles[tone]" />
    {{ status || 'Unknown' }}
  </span>
</template>
