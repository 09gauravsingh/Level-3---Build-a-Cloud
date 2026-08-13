<script setup>
import { computed } from 'vue'

const props = defineProps({
  variant: {
    type: String,
    default: 'primary',
    validator: (value) => ['primary', 'secondary', 'ghost', 'danger'].includes(value),
  },
  size: {
    type: String,
    default: 'md',
    validator: (value) => ['sm', 'md'].includes(value),
  },
  type: { type: String, default: 'button' },
  disabled: { type: Boolean, default: false },
  loading: { type: Boolean, default: false },
  block: { type: Boolean, default: false },
})

const variants = {
  primary:
    'bg-sky-500 text-slate-950 shadow-lg shadow-sky-500/20 hover:bg-sky-400 focus-visible:outline-sky-400',
  secondary:
    'bg-white/5 text-slate-100 ring-1 ring-inset ring-white/10 hover:bg-white/10 focus-visible:outline-slate-400',
  ghost:
    'text-slate-400 hover:bg-white/5 hover:text-slate-100 focus-visible:outline-slate-400',
  danger:
    'bg-rose-500/10 text-rose-300 ring-1 ring-inset ring-rose-500/30 hover:bg-rose-500/20 focus-visible:outline-rose-400',
}

const sizes = {
  sm: 'h-8 px-3 text-xs',
  md: 'h-10 px-4 text-sm',
}

const classes = computed(() => [
  'inline-flex items-center justify-center gap-2 rounded-lg font-semibold transition',
  'focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2',
  'disabled:cursor-not-allowed disabled:opacity-50',
  variants[props.variant],
  sizes[props.size],
  props.block ? 'w-full' : '',
])
</script>

<template>
  <button :type="type" :class="classes" :disabled="disabled || loading">
    <svg
      v-if="loading"
      class="size-4 animate-spin"
      viewBox="0 0 24 24"
      fill="none"
      aria-hidden="true"
    >
      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="3" />
      <path
        class="opacity-90"
        fill="currentColor"
        d="M4 12a8 8 0 0 1 8-8v3a5 5 0 0 0-5 5H4z"
      />
    </svg>
    <slot />
  </button>
</template>
