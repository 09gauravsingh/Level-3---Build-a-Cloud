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
  healthy: 'bg-emerald-50 text-emerald-700 ring-emerald-600/20',
  failed: 'bg-rose-50 text-rose-700 ring-rose-600/20',
  deleting: 'bg-slate-100 text-slate-600 ring-slate-500/20',
  pending: 'bg-amber-50 text-amber-700 ring-amber-600/20',
  unknown: 'bg-slate-100 text-slate-500 ring-slate-500/20',
}

const dotStyles = {
  healthy: 'bg-emerald-500',
  failed: 'bg-rose-500',
  deleting: 'bg-slate-400',
  pending: 'bg-amber-500 animate-pulse',
  unknown: 'bg-slate-400',
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
