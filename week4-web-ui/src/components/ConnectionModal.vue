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
    <div class="space-y-2">
      <div
        v-for="field in fields"
        :key="field.key"
        class="flex items-center gap-3 rounded-lg border border-white/5 bg-slate-950/50 px-3 py-2"
      >
        <span class="w-24 shrink-0 text-[11px] tracking-wide text-slate-500 uppercase">
          {{ field.label }}
        </span>

        <span class="min-w-0 flex-1 truncate font-mono text-xs text-slate-200" :title="field.value">
          {{ display(field) }}
        </span>

        <button
          class="shrink-0 rounded-md px-2 py-1 text-[11px] font-medium text-slate-400 transition hover:bg-white/5 hover:text-slate-100 disabled:opacity-40"
          :disabled="!field.value"
          @click="copy(field)"
        >
          {{ copiedField === field.key ? 'Copied' : 'Copy' }}
        </button>
      </div>
    </div>

    <p class="mt-4 text-xs text-amber-300/80">
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
