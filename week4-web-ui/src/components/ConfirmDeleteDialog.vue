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
    <p class="text-sm text-slate-300">
      This permanently removes the cluster and its storage. Type
      <span class="font-mono text-rose-300">{{ instanceName }}</span> to confirm.
    </p>

    <input
      v-model="typedName"
      type="text"
      :placeholder="instanceName"
      aria-label="Confirm instance name"
      class="mt-4 w-full rounded-lg border border-white/10 bg-slate-950/60 px-3 py-2.5 font-mono text-sm text-slate-100 placeholder:text-slate-600 focus:border-rose-400/60 focus:ring-2 focus:ring-rose-500/25 focus:outline-none"
      @keyup.enter="confirm"
    >

    <template #footer>
      <BaseButton variant="ghost" :disabled="loading" @click="emit('cancel')">Cancel</BaseButton>

      <BaseButton
        variant="danger"
        :loading="loading"
        :disabled="typedName !== instanceName"
        @click="confirm"
      >
        Delete instance
      </BaseButton>
    </template>
  </BaseModal>
</template>
