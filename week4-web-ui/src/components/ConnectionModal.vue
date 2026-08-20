<script setup>
import { computed, ref } from 'vue'
import BaseButton from '@/components/ui/BaseButton.vue'
import BaseModal from '@/components/ui/BaseModal.vue'
import { useToasts } from '@/composables/useToasts'

const props = defineProps({
  connection: { type: Object, required: true },
  instanceName: { type: String, default: '' },
})

const emit = defineEmits(['close'])

const toasts = useToasts()
const revealed = ref(false)
const copiedField = ref('')

const fields = computed(() => [
  { key: 'host', label: 'Host', value: props.connection.host },
  { key: 'port', label: 'Port', value: String(props.connection.port ?? '') },
  { key: 'database', label: 'Database', value: props.connection.database },
  { key: 'username', label: 'Username', value: props.connection.username },
  { key: 'password', label: 'Password', value: props.connection.password, secret: true },
  { key: 'uri', label: 'Connection URI', value: props.connection.uri, secret: true },
])

function display(field) {
  if (!field.value) return '—'

  return field.secret && !revealed.value ? '•'.repeat(Math.min(field.value.length, 24)) : field.value
}

async function copy(field) {
  if (!field.value) return

  try {
    await navigator.clipboard.writeText(field.value)
    copiedField.value = field.key
    setTimeout(() => {
      if (copiedField.value === field.key) copiedField.value = ''
    }, 1500)
  } catch {
    toasts.error('Clipboard is not available in this browser')
  }
}
</script>

<template>
  <BaseModal
    title="Connection details"
    :subtitle="instanceName"
    @close="emit('close')"
  >
    <div class="divide-y divide-slate-100 overflow-hidden rounded-xl border border-slate-200">
      <div
        v-for="field in fields"
        :key="field.key"
        class="flex items-center gap-3 bg-white px-3 py-2.5 transition hover:bg-slate-50"
      >
        <span class="w-24 shrink-0 text-[11px] font-medium tracking-wide text-slate-400 uppercase">
          {{ field.label }}
        </span>

        <span class="min-w-0 flex-1 truncate font-mono text-xs text-slate-800" :title="field.value">
          {{ display(field) }}
        </span>

        <button
          class="shrink-0 rounded-md px-2 py-1 text-[11px] font-semibold transition disabled:opacity-40"
          :class="
            copiedField === field.key
              ? 'bg-emerald-50 text-emerald-600'
              : 'text-slate-500 hover:bg-slate-100 hover:text-slate-900'
          "
          :disabled="!field.value"
          @click="copy(field)"
        >
          {{ copiedField === field.key ? 'Copied' : 'Copy' }}
        </button>
      </div>
    </div>

    <p
      class="mt-4 flex items-start gap-2 rounded-lg border border-amber-200 bg-amber-50 px-3 py-2 text-xs text-amber-800"
    >
      <svg class="mt-px size-4 shrink-0" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
        <path
          d="M8.49 3.17c.67-1.16 2.35-1.16 3.02 0l5.14 8.9c.67 1.17-.17 2.63-1.51 2.63H4.86c-1.34 0-2.18-1.46-1.51-2.63l5.14-8.9ZM10 7a.75.75 0 0 1 .75.75v3a.75.75 0 0 1-1.5 0v-3A.75.75 0 0 1 10 7Zm0 7a1 1 0 1 0 0-2 1 1 0 0 0 0 2Z"
        />
      </svg>
      These credentials come straight from the CloudNativePG secret. Keep them private.
    </p>

    <template #footer>
      <BaseButton variant="secondary" @click="revealed = !revealed">
        {{ revealed ? 'Hide secrets' : 'Reveal secrets' }}
      </BaseButton>

      <BaseButton @click="emit('close')">Close</BaseButton>
    </template>
  </BaseModal>
</template>
