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
    <label :for="id" class="block text-xs font-medium tracking-wide text-slate-400 uppercase">
      {{ label }}
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
      class="w-full rounded-lg border border-white/10 bg-slate-900/60 px-3 py-2.5 text-sm text-slate-100 transition placeholder:text-slate-600 focus:border-sky-400/60 focus:ring-2 focus:ring-sky-500/25 focus:outline-none"
      :class="error ? 'border-rose-500/60 focus:border-rose-500/60 focus:ring-rose-500/25' : ''"
      @input="onInput"
    >

    <p v-if="error" :id="`${id}-help`" class="text-xs text-rose-400">{{ error }}</p>
    <p v-else-if="hint" :id="`${id}-help`" class="text-xs text-slate-500">{{ hint }}</p>
  </div>
</template>
