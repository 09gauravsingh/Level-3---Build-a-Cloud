<script setup>
import { ref } from 'vue'
import BaseButton from '@/components/ui/BaseButton.vue'
import BaseModal from '@/components/ui/BaseModal.vue'

const props = defineProps({
  instanceName: { type: String, required: true },
  loading: { type: Boolean, default: false },
})

const emit = defineEmits(['cancel', 'confirm'])

const typedName = ref('')

function confirm() {
  if (typedName.value !== props.instanceName) return

  emit('confirm')
}
</script>

<template>
  <BaseModal title="Delete instance" :subtitle="instanceName" @close="emit('cancel')">
    <div class="flex items-start gap-3 rounded-xl border border-rose-200 bg-rose-50 px-4 py-3">
      <svg
        class="mt-0.5 size-5 shrink-0 text-rose-500"
        viewBox="0 0 20 20"
        fill="currentColor"
        aria-hidden="true"
      >
        <path
          d="M8.49 3.17c.67-1.16 2.35-1.16 3.02 0l5.14 8.9c.67 1.17-.17 2.63-1.51 2.63H4.86c-1.34 0-2.18-1.46-1.51-2.63l5.14-8.9ZM10 7a.75.75 0 0 1 .75.75v3a.75.75 0 0 1-1.5 0v-3A.75.75 0 0 1 10 7Zm0 7a1 1 0 1 0 0-2 1 1 0 0 0 0 2Z"
        />
      </svg>

      <p class="text-sm text-rose-800">
        This permanently removes the cluster and its storage. The action cannot be undone.
      </p>
    </div>

    <p class="mt-4 text-sm text-slate-600">
      Type <span class="font-mono font-semibold text-slate-900">{{ instanceName }}</span> to
      confirm.
    </p>

    <input
      v-model="typedName"
      type="text"
      :placeholder="instanceName"
      aria-label="Confirm instance name"
      class="mt-2 w-full rounded-lg border border-slate-200 bg-white px-3 py-2.5 font-mono text-sm text-slate-900 shadow-sm transition placeholder:text-slate-400 hover:border-slate-300 focus:border-rose-500 focus:ring-4 focus:ring-rose-500/10 focus:outline-none"
      @keyup.enter="confirm"
    >

    <template #footer>
      <BaseButton variant="ghost" :disabled="loading" @click="emit('cancel')">Cancel</BaseButton>

      <BaseButton
        variant="danger-solid"
        :loading="loading"
        :disabled="typedName !== instanceName"
        @click="confirm"
      >
        Delete instance
      </BaseButton>
    </template>
  </BaseModal>
</template>
