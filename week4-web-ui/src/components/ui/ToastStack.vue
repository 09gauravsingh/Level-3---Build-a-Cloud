<script setup>
import { useToasts } from '@/composables/useToasts'

const { toasts, dismiss } = useToasts()

const styles = {
  success: 'text-emerald-600',
  error: 'text-rose-600',
  info: 'text-indigo-600',
}

const icons = {
  success: 'M10 18a8 8 0 1 0 0-16 8 8 0 0 0 0 16Zm3.86-9.47a.75.75 0 0 0-1.22-.86l-3.24 4.6-1.6-1.6a.75.75 0 0 0-1.06 1.06l2.22 2.22a.75.75 0 0 0 1.14-.1l3.76-5.32Z',
  error:
    'M10 18a8 8 0 1 0 0-16 8 8 0 0 0 0 16ZM9.25 6a.75.75 0 0 1 1.5 0v4.5a.75.75 0 0 1-1.5 0V6ZM10 14.5a1 1 0 1 1 0-2 1 1 0 0 1 0 2Z',
  info: 'M10 18a8 8 0 1 0 0-16 8 8 0 0 0 0 16Zm.75-11.5a.75.75 0 1 1-1.5 0 .75.75 0 0 1 1.5 0ZM9.25 9a.75.75 0 0 1 1.5 0v5a.75.75 0 0 1-1.5 0V9Z',
}
</script>

<template>
  <div
    class="pointer-events-none fixed inset-x-0 top-4 z-50 flex flex-col items-center gap-2 px-4 sm:inset-x-auto sm:right-4 sm:items-end"
    role="status"
    aria-live="polite"
  >
    <TransitionGroup
      enter-active-class="transition duration-200 ease-out"
      enter-from-class="-translate-y-2 opacity-0"
      leave-active-class="transition duration-150 ease-in"
      leave-to-class="translate-y-1 opacity-0"
    >
      <div
        v-for="toast in toasts"
        :key="toast.id"
        class="pointer-events-auto flex w-full max-w-sm items-start gap-3 rounded-xl border border-slate-200 bg-white px-4 py-3 text-sm text-slate-700 shadow-lg shadow-slate-900/10"
      >
        <svg
          class="mt-0.5 size-5 shrink-0"
          :class="styles[toast.type]"
          viewBox="0 0 20 20"
          fill="currentColor"
          aria-hidden="true"
        >
          <path :d="icons[toast.type]" />
        </svg>

        <span class="flex-1">{{ toast.message }}</span>

        <button
          class="shrink-0 rounded-md p-0.5 text-slate-400 transition hover:bg-slate-100 hover:text-slate-700"
          aria-label="Dismiss notification"
          @click="dismiss(toast.id)"
        >
          <svg class="size-4" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
            <path
              d="M6.28 5.22a.75.75 0 0 0-1.06 1.06L8.94 10l-3.72 3.72a.75.75 0 1 0 1.06 1.06L10 11.06l3.72 3.72a.75.75 0 1 0 1.06-1.06L11.06 10l3.72-3.72a.75.75 0 0 0-1.06-1.06L10 8.94 6.28 5.22Z"
            />
          </svg>
        </button>
      </div>
    </TransitionGroup>
  </div>
</template>
