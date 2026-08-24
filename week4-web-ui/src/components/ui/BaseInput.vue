<script setup>
import { useId } from 'vue'

defineProps({
  label: { type: String, required: true },
  modelValue: { type: [String, Number], default: '' },
  type: { type: String, default: 'text' },
  placeholder: { type: String, default: '' },
  hint: { type: String, default: '' },
  error: { type: String, default: '' },
  min: { type: [String, Number], default: undefined },
  max: { type: [String, Number], default: undefined },
  autocomplete: { type: String, default: 'off' },
  required: { type: Boolean, default: false },
})

const emit = defineEmits(['update:modelValue'])

const id = useId()

function onInput(event) {
  const { value } = event.target
  emit('update:modelValue', event.target.type === 'number' ? Number(value) : value)
}
</script>

<template>
  <div class="space-y-1.5">
    <label :for="id" class="block text-sm font-medium text-neutral-700">
      {{ label }}
      <span v-if="required" class="text-neutral-400">*</span>
    </label>

    <input
      :id="id"
      :value="modelValue"
      :type="type"
      :placeholder="placeholder"
      :min="min"
      :max="max"
      :autocomplete="autocomplete"
      :required="required"
      :aria-invalid="Boolean(error)"
      :aria-describedby="hint || error ? `${id}-help` : undefined"
      class="w-full rounded-md border border-neutral-200 bg-white px-3 py-2.5 text-sm text-neutral-900 transition placeholder:text-neutral-400 hover:border-neutral-300 focus:border-neutral-900 focus:ring-2 focus:ring-neutral-900/10 focus:outline-none"
      :class="error ? 'border-rose-400 focus:border-rose-500 focus:ring-rose-500/10' : ''"
      @input="onInput"
    >

    <p v-if="error" :id="`${id}-help`" class="text-xs font-medium text-rose-700">{{ error }}</p>
    <p v-else-if="hint" :id="`${id}-help`" class="text-xs text-neutral-500">{{ hint }}</p>
  </div>
</template>
