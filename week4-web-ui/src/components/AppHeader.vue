<script setup>
import { computed } from 'vue'
import BaseButton from '@/components/ui/BaseButton.vue'

const props = defineProps({
  username: { type: String, default: '' },
  isAdmin: { type: Boolean, default: false },
  view: { type: String, default: 'dashboard' },
})

const emit = defineEmits(['logout', 'navigate'])

const initials = computed(() => props.username.slice(0, 2).toUpperCase())

const tabs = [
  { id: 'dashboard', label: 'Dashboard' },
  { id: 'logs', label: 'Logs' },
]
</script>

<template>
  <header class="sticky top-0 z-30 border-b border-slate-200 bg-white/90 backdrop-blur">
    <div class="mx-auto flex max-w-7xl items-center justify-between gap-4 px-4 sm:px-6">
      <div class="flex min-w-0 items-center gap-6">
        <div class="flex items-center gap-3 py-3">
          <div
            class="flex size-9 items-center justify-center rounded-xl bg-indigo-600 shadow-sm shadow-indigo-600/30"
          >
            <svg class="size-5 text-white" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true">
              <path
                d="M12 2c-4.42 0-8 1.34-8 3v14c0 1.66 3.58 3 8 3s8-1.34 8-3V5c0-1.66-3.58-3-8-3Zm0 2c3.87 0 6 1.1 6 1s-2.13 1-6 1-6-1.1-6-1 2.13-1 6-1Zm6 15c0 .1-2.13 1-6 1s-6-.9-6-1v-2.35C7.5 17.5 9.6 18 12 18s4.5-.5 6-1.35V19Zm0-4.5c0 .1-2.13 1-6 1s-6-.9-6-1v-2.35C7.5 13 9.6 13.5 12 13.5s4.5-.5 6-1.35v2.35Zm0-4.5c0 .1-2.13 1-6 1s-6-.9-6-1V7.65C7.5 8.5 9.6 9 12 9s4.5-.5 6-1.35V10Z"
              />
            </svg>
          </div>

          <div class="hidden sm:block">
            <p class="text-sm font-semibold text-slate-900">PostgreSQL PaaS</p>
            <p class="text-xs text-slate-500">Self-service database platform</p>
          </div>
        </div>

        <nav class="flex items-stretch gap-1" aria-label="Primary">
          <button
            v-for="tab in tabs"
            :key="tab.id"
            type="button"
            class="border-b-2 px-3 py-4 text-sm font-medium transition"
            :class="
              view === tab.id
                ? 'border-indigo-600 text-indigo-600'
                : 'border-transparent text-slate-500 hover:border-slate-200 hover:text-slate-800'
            "
            :aria-current="view === tab.id ? 'page' : undefined"
            @click="emit('navigate', tab.id)"
          >
            {{ tab.label }}
          </button>
        </nav>
      </div>

      <div class="flex items-center gap-3 py-3">
        <div class="hidden items-center gap-2.5 sm:flex">
          <span
            class="flex size-8 items-center justify-center rounded-full bg-slate-100 text-xs font-semibold text-slate-600"
            aria-hidden="true"
          >
            {{ initials }}
          </span>

          <span class="text-xs text-slate-500">
            Signed in as <span class="font-medium text-slate-800">{{ username }}</span>
          </span>

          <span
            v-if="isAdmin"
            class="rounded-full bg-indigo-50 px-2 py-0.5 text-[10px] font-semibold tracking-wide text-indigo-700 uppercase ring-1 ring-indigo-100"
          >
            admin
          </span>
        </div>

        <BaseButton variant="secondary" size="sm" @click="emit('logout')">Log out</BaseButton>
      </div>
    </div>
  </header>
</template>
