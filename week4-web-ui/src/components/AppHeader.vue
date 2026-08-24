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
  <header class="sticky top-0 z-30 border-b border-neutral-200 bg-white">
    <div class="flex h-14 w-full items-center justify-between gap-4 px-6 lg:px-8">
      <div class="flex min-w-0 items-center gap-2.5">
        <div class="flex size-8 items-center justify-center rounded-md bg-neutral-900 text-white">
          <svg class="size-4" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true">
            <path
              d="M12 2c-4.42 0-8 1.34-8 3v14c0 1.66 3.58 3 8 3s8-1.34 8-3V5c0-1.66-3.58-3-8-3Zm0 2c3.87 0 6 1.1 6 1s-2.13 1-6 1-6-1.1-6-1 2.13-1 6-1Zm6 15c0 .1-2.13 1-6 1s-6-.9-6-1v-2.35C7.5 17.5 9.6 18 12 18s4.5-.5 6-1.35V19Zm0-4.5c0 .1-2.13 1-6 1s-6-.9-6-1v-2.35C7.5 13 9.6 13.5 12 13.5s4.5-.5 6-1.35v2.35Zm0-4.5c0 .1-2.13 1-6 1s-6-.9-6-1V7.65C7.5 8.5 9.6 9 12 9s4.5-.5 6-1.35V10Z"
            />
          </svg>
        </div>
        <p class="truncate text-sm font-semibold text-neutral-900">PostgreSQL PaaS</p>
      </div>

      <div class="flex items-center gap-3">
        <div class="hidden items-center gap-2 sm:flex">
          <span
            class="flex size-7 items-center justify-center rounded-full bg-neutral-100 text-[11px] font-semibold text-neutral-700"
            aria-hidden="true"
          >
            {{ initials }}
          </span>
          <span class="text-sm text-neutral-600">{{ username }}</span>
          <span
            v-if="isAdmin"
            class="rounded px-1.5 py-0.5 text-[10px] font-semibold tracking-wide text-neutral-500 uppercase ring-1 ring-neutral-200"
          >
            admin
          </span>
        </div>
        <BaseButton variant="secondary" size="sm" @click="emit('logout')">Sign out</BaseButton>
      </div>
    </div>

    <nav class="flex w-full gap-1 px-6 lg:px-8" aria-label="Primary">
      <button
        v-for="tab in tabs"
        :key="tab.id"
        type="button"
        class="-mb-px border-b-2 px-3 py-2.5 text-sm transition"
        :class="
          view === tab.id
            ? 'border-neutral-900 font-semibold text-neutral-900'
            : 'border-transparent font-medium text-neutral-500 hover:border-neutral-300 hover:text-neutral-800'
        "
        :aria-current="view === tab.id ? 'page' : undefined"
        @click="emit('navigate', tab.id)"
      >
        {{ tab.label }}
      </button>
    </nav>
  </header>
</template>
