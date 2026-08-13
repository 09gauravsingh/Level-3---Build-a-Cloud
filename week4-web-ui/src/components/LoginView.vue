<script setup>
import { ref } from 'vue'
import BaseButton from '@/components/ui/BaseButton.vue'
import BaseInput from '@/components/ui/BaseInput.vue'

defineProps({
  error: { type: String, default: '' },
  loading: { type: Boolean, default: false },
})

const emit = defineEmits(['submit'])

const username = ref('')
const password = ref('')

function submit() {
  emit('submit', { username: username.value, password: password.value })
}
</script>

<template>
  <div class="flex min-h-screen items-center justify-center px-4 py-12">
    <div class="w-full max-w-sm">
      <div class="mb-8 text-center">
        <div
          class="mx-auto flex size-12 items-center justify-center rounded-xl bg-sky-500/10 ring-1 ring-sky-500/30"
        >
          <svg class="size-6 text-sky-400" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true">
            <path
              d="M12 2c-4.42 0-8 1.34-8 3v14c0 1.66 3.58 3 8 3s8-1.34 8-3V5c0-1.66-3.58-3-8-3Zm0 2c3.87 0 6 1.1 6 1s-2.13 1-6 1-6-1.1-6-1 2.13-1 6-1Zm6 15c0 .1-2.13 1-6 1s-6-.9-6-1v-2.35C7.5 17.5 9.6 18 12 18s4.5-.5 6-1.35V19Zm0-4.5c0 .1-2.13 1-6 1s-6-.9-6-1v-2.35C7.5 13 9.6 13.5 12 13.5s4.5-.5 6-1.35v2.35Zm0-4.5c0 .1-2.13 1-6 1s-6-.9-6-1V7.65C7.5 8.5 9.6 9 12 9s4.5-.5 6-1.35V10Z"
            />
          </svg>
        </div>

        <h1 class="mt-4 text-xl font-semibold text-slate-100">PostgreSQL PaaS Platform</h1>
        <p class="mt-1 text-sm text-slate-500">Sign in to manage your database instances.</p>
      </div>

      <form
        class="space-y-4 rounded-2xl border border-white/10 bg-slate-900/50 p-6 shadow-xl shadow-slate-950/40 backdrop-blur"
        @submit.prevent="submit"
      >
        <BaseInput
          v-model="username"
          label="Username"
          placeholder="admin"
          autocomplete="username"
          required
        />

        <BaseInput
          v-model="password"
          label="Password"
          type="password"
          placeholder="••••••••"
          autocomplete="current-password"
          required
        />

        <p
          v-if="error"
          class="rounded-lg border border-rose-500/30 bg-rose-500/10 px-3 py-2 text-xs text-rose-300"
        >
          {{ error }}
        </p>

        <BaseButton type="submit" block :loading="loading">
          {{ loading ? 'Signing in…' : 'Sign in' }}
        </BaseButton>
      </form>

      <p class="mt-6 text-center text-xs text-slate-600">
        Sessions are backed by a JWT that expires after one hour.
      </p>
    </div>
  </div>
</template>
