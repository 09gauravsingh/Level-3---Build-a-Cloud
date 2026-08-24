<script setup>
import { computed } from 'vue'

const props = defineProps({
  variant: {
    type: String,
    default: 'primary',
    validator: (value) => ['primary', 'secondary', 'ghost', 'danger', 'danger-solid'].includes(value),
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
    'bg-neutral-900 text-white hover:bg-neutral-800 active:bg-neutral-950 focus-visible:outline-neutral-900',
  secondary:
    'bg-white text-neutral-700 ring-1 ring-inset ring-neutral-200 hover:bg-neutral-50 hover:text-neutral-950 focus-visible:outline-neutral-400',
  ghost: 'text-neutral-500 hover:bg-neutral-100 hover:text-neutral-900 focus-visible:outline-neutral-400',
  danger:
    'bg-white text-rose-700 ring-1 ring-inset ring-rose-200 shadow-sm hover:bg-rose-50 focus-visible:outline-rose-500',
  'danger-solid':
    'bg-rose-700 text-white shadow-sm shadow-rose-700/20 hover:bg-rose-600 active:bg-rose-800 focus-visible:outline-rose-700',
}

const sizes = {
  sm: 'h-8 px-3 text-xs tracking-wide',
  md: 'h-10 px-4 text-sm',
}

const classes = computed(() => [
  'inline-flex items-center justify-center gap-2 rounded-md font-semibold transition duration-150',
  'focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2',
  'disabled:cursor-not-allowed disabled:opacity-50 disabled:shadow-none',
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
