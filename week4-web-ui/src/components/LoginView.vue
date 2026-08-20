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
  <div class="flex min-h-screen items-center justify-center px-4 py-12">
    <div
      class="grid w-full max-w-4xl overflow-hidden rounded-3xl border border-slate-200 bg-white shadow-xl shadow-slate-900/10 lg:grid-cols-2"
    >
      <div class="hidden flex-col justify-between bg-indigo-600 p-10 text-indigo-50 lg:flex">
        <div class="flex items-center gap-3">
          <div class="flex size-10 items-center justify-center rounded-xl bg-white/15">
            <svg class="size-5 text-white" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true">
              <path
                d="M12 2c-4.42 0-8 1.34-8 3v14c0 1.66 3.58 3 8 3s8-1.34 8-3V5c0-1.66-3.58-3-8-3Zm0 2c3.87 0 6 1.1 6 1s-2.13 1-6 1-6-1.1-6-1 2.13-1 6-1Zm6 15c0 .1-2.13 1-6 1s-6-.9-6-1v-2.35C7.5 17.5 9.6 18 12 18s4.5-.5 6-1.35V19Zm0-4.5c0 .1-2.13 1-6 1s-6-.9-6-1v-2.35C7.5 13 9.6 13.5 12 13.5s4.5-.5 6-1.35v2.35Zm0-4.5c0 .1-2.13 1-6 1s-6-.9-6-1V7.65C7.5 8.5 9.6 9 12 9s4.5-.5 6-1.35V10Z"
              />
            </svg>
          </div>

          <p class="text-sm font-semibold text-white">PostgreSQL PaaS</p>
        </div>

        <div class="py-10">
          <h2 class="text-2xl font-semibold text-white">Databases on demand.</h2>
          <p class="mt-2 text-sm text-indigo-100">
            A self-service platform for teams that need PostgreSQL without the Kubernetes detour.
          </p>

          <ul class="mt-8 space-y-3">
            <li v-for="item in highlights" :key="item" class="flex items-start gap-3 text-sm">
              <svg
                class="mt-0.5 size-4 shrink-0 text-indigo-200"
                viewBox="0 0 20 20"
                fill="currentColor"
                aria-hidden="true"
              >
                <path
                  d="M10 18a8 8 0 1 0 0-16 8 8 0 0 0 0 16Zm3.86-9.47a.75.75 0 0 0-1.22-.86l-3.24 4.6-1.6-1.6a.75.75 0 0 0-1.06 1.06l2.22 2.22a.75.75 0 0 0 1.14-.1l3.76-5.32Z"
                />
              </svg>
              <span class="text-indigo-100">{{ item }}</span>
            </li>
          </ul>
        </div>

        <p class="text-xs text-indigo-200">Backed by CloudNativePG on Kubernetes.</p>
      </div>

      <div class="p-8 sm:p-10">
        <div
          class="flex size-11 items-center justify-center rounded-xl bg-indigo-600 shadow-sm shadow-indigo-600/30 lg:hidden"
        >
          <svg class="size-6 text-white" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true">
            <path
              d="M12 2c-4.42 0-8 1.34-8 3v14c0 1.66 3.58 3 8 3s8-1.34 8-3V5c0-1.66-3.58-3-8-3Zm0 2c3.87 0 6 1.1 6 1s-2.13 1-6 1-6-1.1-6-1 2.13-1 6-1Zm6 15c0 .1-2.13 1-6 1s-6-.9-6-1v-2.35C7.5 17.5 9.6 18 12 18s4.5-.5 6-1.35V19Z"
            />
          </svg>
        </div>

        <h1 class="mt-5 text-2xl font-semibold tracking-tight text-slate-900 lg:mt-0">
          {{ isRegister ? 'Create your account' : 'Welcome back' }}
        </h1>
        <p class="mt-1 text-sm text-slate-500">
          {{
            isRegister
              ? 'Register to provision and manage your own database instances.'
              : 'Sign in to manage your database instances.'
          }}
        </p>

        <div class="mt-6 grid grid-cols-2 gap-1 rounded-lg bg-slate-100 p-1 text-sm font-medium">
          <button
            type="button"
            class="rounded-md px-3 py-1.5 transition"
            :class="
              !isRegister ? 'bg-white text-slate-900 shadow-sm' : 'text-slate-500 hover:text-slate-800'
            "
            @click="setMode('login')"
          >
            Sign in
          </button>
          <button
            type="button"
            class="rounded-md px-3 py-1.5 transition"
            :class="
              isRegister ? 'bg-white text-slate-900 shadow-sm' : 'text-slate-500 hover:text-slate-800'
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
            class="flex items-start gap-2 rounded-lg border border-rose-200 bg-rose-50 px-3 py-2 text-xs font-medium text-rose-700"
          >
            <svg class="mt-px size-4 shrink-0" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
              <path
                d="M10 18a8 8 0 1 0 0-16 8 8 0 0 0 0 16ZM9.25 6a.75.75 0 0 1 1.5 0v4.5a.75.75 0 0 1-1.5 0V6ZM10 14.5a1 1 0 1 1 0-2 1 1 0 0 1 0 2Z"
              />
            </svg>
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

        <p class="mt-6 text-xs text-slate-400">
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
