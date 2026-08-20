<script setup>
import { onBeforeUnmount, onMounted } from 'vue'

defineProps({
  title: { type: String, required: true },
  subtitle: { type: String, default: '' },
  size: {
    type: String,
    default: 'md',
    validator: (value) => ['md', 'xl'].includes(value),
  },
})

const emit = defineEmits(['close'])

function onKeydown(event) {
  if (event.key === 'Escape') {
    emit('close')
  }
}

onMounted(() => document.addEventListener('keydown', onKeydown))
onBeforeUnmount(() => document.removeEventListener('keydown', onKeydown))
</script>

<template>
  <Teleport to="body">
    <div
      class="fixed inset-0 z-40 flex items-end justify-center bg-slate-900/25 p-4 backdrop-blur-sm sm:items-center"
      role="dialog"
      aria-modal="true"
      @click.self="emit('close')"
    >
      <Transition
        appear
        enter-active-class="transition duration-200 ease-out"
        enter-from-class="translate-y-3 scale-95 opacity-0 sm:translate-y-0"
      >
        <div
          class="w-full rounded-2xl border border-slate-200 bg-white shadow-2xl shadow-slate-900/15"
          :class="size === 'xl' ? 'max-w-3xl' : 'max-w-lg'"
        >
          <header
            class="flex items-start justify-between gap-4 border-b border-slate-100 px-5 py-4"
          >
            <div class="min-w-0">
              <h2 class="text-base font-semibold text-slate-900">{{ title }}</h2>
              <p v-if="subtitle" class="mt-0.5 truncate font-mono text-xs text-slate-500">
                {{ subtitle }}
              </p>
            </div>

            <button
              class="rounded-lg p-1 text-slate-400 transition hover:bg-slate-100 hover:text-slate-700"
              aria-label="Close dialog"
              @click="emit('close')"
            >
              <svg class="size-5" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                <path
                  d="M6.28 5.22a.75.75 0 0 0-1.06 1.06L8.94 10l-3.72 3.72a.75.75 0 1 0 1.06 1.06L10 11.06l3.72 3.72a.75.75 0 1 0 1.06-1.06L11.06 10l3.72-3.72a.75.75 0 0 0-1.06-1.06L10 8.94 6.28 5.22Z"
                />
              </svg>
            </button>
          </header>

          <div class="px-5 py-4">
            <slot />
          </div>

          <footer
            v-if="$slots.footer"
            class="flex justify-end gap-2 rounded-b-2xl border-t border-slate-100 bg-slate-50 px-5 py-4"
          >
            <slot name="footer" />
          </footer>
        </div>
      </Transition>
    </div>
  </Teleport>
</template>
