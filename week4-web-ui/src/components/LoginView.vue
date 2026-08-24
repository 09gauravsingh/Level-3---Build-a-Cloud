<script setup>
import { computed, ref } from 'vue'
import BaseButton from '@/components/ui/BaseButton.vue'
import BaseInput from '@/components/ui/BaseInput.vue'

defineProps({
  error: { type: String, default: '' },
  loading: { type: Boolean, default: false },
})

const emit = defineEmits(['submit', 'register'])

const mode = ref('login')
const username = ref('')
const password = ref('')
const confirmPassword = ref('')
const localError = ref('')

const isRegister = computed(() => mode.value === 'register')

const highlights = [
  'Provision a managed PostgreSQL cluster in one form',
  'Watch replica health roll out in real time',
  'Copy connection credentials without touching kubectl',
]

function setMode(next) {
  mode.value = next
  localError.value = ''
  confirmPassword.value = ''
}

function submit() {
  localError.value = ''

  if (isRegister.value) {
    if (password.value !== confirmPassword.value) {
      localError.value = 'Passwords do not match'
      return
    }

    emit('register', { username: username.value, password: password.value })
    return
  }

  emit('submit', { username: username.value, password: password.value })
}
</script>

<template>
  <div class="grid min-h-screen w-full lg:grid-cols-2">
    <div class="hidden flex-col justify-between bg-neutral-950 px-12 py-10 text-neutral-200 lg:flex">
      <div class="flex items-center gap-3">
        <div class="flex size-9 items-center justify-center rounded-md bg-white/10">
          <svg class="size-5 text-white" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true">
            <path
              d="M12 2c-4.42 0-8 1.34-8 3v14c0 1.66 3.58 3 8 3s8-1.34 8-3V5c0-1.66-3.58-3-8-3Zm0 2c3.87 0 6 1.1 6 1s-2.13 1-6 1-6-1.1-6-1 2.13-1 6-1Zm6 15c0 .1-2.13 1-6 1s-6-.9-6-1v-2.35C7.5 17.5 9.6 18 12 18s4.5-.5 6-1.35V19Zm0-4.5c0 .1-2.13 1-6 1s-6-.9-6-1v-2.35C7.5 13 9.6 13.5 12 13.5s4.5-.5 6-1.35v2.35Zm0-4.5c0 .1-2.13 1-6 1s-6-.9-6-1V7.65C7.5 8.5 9.6 9 12 9s4.5-.5 6-1.35V10Z"
            />
          </svg>
        </div>
        <p class="text-sm font-semibold text-white">PostgreSQL PaaS</p>
      </div>

      <div>
        <h2 class="text-4xl font-semibold tracking-tight text-white">Managed PostgreSQL for your team.</h2>
        <p class="mt-4 max-w-md text-sm leading-6 text-neutral-400">
          A self-service platform for teams that need PostgreSQL without the Kubernetes detour.
        </p>
        <ul class="mt-10 space-y-3">
          <li v-for="item in highlights" :key="item" class="flex items-start gap-3 text-sm leading-6 text-neutral-300">
            <span class="mt-2 size-1.5 shrink-0 rounded-full bg-white" aria-hidden="true" />
            {{ item }}
          </li>
        </ul>
      </div>

      <p class="text-xs text-neutral-500">Backed by CloudNativePG on Kubernetes.</p>
    </div>

    <div class="flex min-h-screen items-center justify-center bg-white px-6 py-12">
      <div class="w-full max-w-md">
        <div class="mb-8 flex size-10 items-center justify-center rounded-md bg-neutral-900 text-white lg:hidden">
          <svg class="size-5" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true">
            <path
              d="M12 2c-4.42 0-8 1.34-8 3v14c0 1.66 3.58 3 8 3s8-1.34 8-3V5c0-1.66-3.58-3-8-3Zm0 2c3.87 0 6 1.1 6 1s-2.13 1-6 1-6-1.1-6-1 2.13-1 6-1Zm6 15c0 .1-2.13 1-6 1s-6-.9-6-1v-2.35C7.5 17.5 9.6 18 12 18s4.5-.5 6-1.35V19Z"
            />
          </svg>
        </div>

        <h1 class="text-2xl font-semibold tracking-tight text-neutral-900">
          {{ isRegister ? 'Create your account' : 'Welcome back' }}
        </h1>
        <p class="mt-2 text-sm text-neutral-500">
          {{
            isRegister
              ? 'Register to provision and manage your own database instances.'
              : 'Sign in to manage your database instances.'
          }}
        </p>

        <div class="mt-6 flex border-b border-neutral-200">
          <button
            type="button"
            class="-mb-px border-b-2 px-3 py-2 text-sm transition"
            :class="
              !isRegister
                ? 'border-neutral-900 font-semibold text-neutral-900'
                : 'border-transparent font-medium text-neutral-500 hover:text-neutral-800'
            "
            @click="setMode('login')"
          >
            Sign in
          </button>
          <button
            type="button"
            class="-mb-px border-b-2 px-3 py-2 text-sm transition"
            :class="
              isRegister
                ? 'border-neutral-900 font-semibold text-neutral-900'
                : 'border-transparent font-medium text-neutral-500 hover:text-neutral-800'
            "
            @click="setMode('register')"
          >
            Create account
          </button>
        </div>

        <form class="mt-6 space-y-4" @submit.prevent="submit">
          <BaseInput
            v-model="username"
            label="Username"
            :placeholder="isRegister ? 'alice' : 'admin'"
            autocomplete="username"
            required
          />

          <BaseInput
            v-model="password"
            label="Password"
            type="password"
            placeholder="••••••••"
            :autocomplete="isRegister ? 'new-password' : 'current-password'"
            :hint="isRegister ? 'At least 8 characters. Letters, numbers and hyphens in the username.' : ''"
            required
          />

          <BaseInput
            v-if="isRegister"
            v-model="confirmPassword"
            label="Confirm password"
            type="password"
            placeholder="••••••••"
            autocomplete="new-password"
            required
          />

          <p
            v-if="localError || error"
            class="rounded-md border border-red-200 bg-red-50 px-3 py-2 text-xs font-medium text-red-800"
          >
            {{ localError || error }}
          </p>

          <BaseButton type="submit" block :loading="loading">
            <template v-if="isRegister">
              {{ loading ? 'Creating account…' : 'Create account' }}
            </template>
            <template v-else>
              {{ loading ? 'Signing in…' : 'Sign in' }}
            </template>
          </BaseButton>
        </form>

        <p class="mt-6 text-xs text-neutral-400">
          {{
            isRegister
              ? 'Your dashboard only lists the instances you create. Sessions last one hour.'
              : 'Sessions are backed by a JWT that expires after one hour.'
          }}
        </p>
      </div>
    </div>
  </div>
</template>
