<script setup>
import { onBeforeUnmount, onMounted } from 'vue'

defineProps({
  title: { type: String, required: true },
  subtitle: { type: String, default: '' },
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
      class="fixed inset-0 z-40 flex items-end justify-center bg-slate-950/70 p-4 backdrop-blur-sm sm:items-center"
      role="dialog"
      aria-modal="true"
      @click.self="emit('close')"
    >
      <div
        class="w-full max-w-lg rounded-2xl border border-white/10 bg-slate-900 shadow-2xl shadow-slate-950/60"
      >
        <header class="flex items-start justify-between gap-4 border-b border-white/10 px-5 py-4">
          <div>
            <h2 class="text-base font-semibold text-slate-100">{{ title }}</h2>
            <p v-if="subtitle" class="mt-0.5 text-xs text-slate-500">{{ subtitle }}</p>
          </div>

          <button
            class="rounded-lg p-1 text-slate-500 transition hover:bg-white/5 hover:text-slate-200"
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
          class="flex justify-end gap-2 border-t border-white/10 px-5 py-4"
        >
          <slot name="footer" />
        </footer>
      </div>
    </div>
  </Teleport>
</template>
