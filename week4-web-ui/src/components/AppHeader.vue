<script setup>
import BaseButton from '@/components/ui/BaseButton.vue'

defineProps({
  username: { type: String, default: 'admin' },
})

const emit = defineEmits(['logout'])
</script>

<template>
  <header class="sticky top-0 z-30 border-b border-white/10 bg-slate-950/80 backdrop-blur">
    <div class="mx-auto flex max-w-6xl items-center justify-between gap-4 px-4 py-3 sm:px-6">
      <div class="flex items-center gap-3">
        <div
          class="flex size-9 items-center justify-center rounded-lg bg-sky-500/10 ring-1 ring-sky-500/30"
        >
          <svg class="size-5 text-sky-400" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true">
            <path
              d="M12 2c-4.42 0-8 1.34-8 3v14c0 1.66 3.58 3 8 3s8-1.34 8-3V5c0-1.66-3.58-3-8-3Zm0 2c3.87 0 6 1.1 6 1s-2.13 1-6 1-6-1.1-6-1 2.13-1 6-1Zm6 15c0 .1-2.13 1-6 1s-6-.9-6-1v-2.35C7.5 17.5 9.6 18 12 18s4.5-.5 6-1.35V19Zm0-4.5c0 .1-2.13 1-6 1s-6-.9-6-1v-2.35C7.5 13 9.6 13.5 12 13.5s4.5-.5 6-1.35v2.35Zm0-4.5c0 .1-2.13 1-6 1s-6-.9-6-1V7.65C7.5 8.5 9.6 9 12 9s4.5-.5 6-1.35V10Z"
            />
          </svg>
        </div>

        <div>
          <p class="text-sm font-semibold text-slate-100">PostgreSQL PaaS</p>
          <p class="text-xs text-slate-500">Self-service database platform</p>
        </div>
      </div>

      <div class="flex items-center gap-3">
        <span class="hidden text-xs text-slate-500 sm:inline">
          Signed in as <span class="text-slate-300">{{ username }}</span>
        </span>

        <BaseButton variant="secondary" size="sm" @click="emit('logout')">Log out</BaseButton>
      </div>
    </div>
  </header>
</template>
