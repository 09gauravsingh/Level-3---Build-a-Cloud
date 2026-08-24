<script setup>
defineProps({
  title: { type: String, default: '' },
  subtitle: { type: String, default: '' },
})
</script>

<template>
  <section class="overflow-hidden rounded-lg border border-neutral-200 bg-white">
    <header
      v-if="title || $slots.actions"
      class="flex flex-wrap items-center justify-between gap-3 border-b border-neutral-200 px-5 py-3"
    >
      <div>
        <h2 class="text-sm font-semibold text-neutral-900">{{ title }}</h2>
        <p v-if="subtitle" class="mt-0.5 text-xs text-neutral-500">{{ subtitle }}</p>
      </div>

      <div class="flex items-center gap-2">
        <slot name="actions" />
      </div>
    </header>

    <div class="p-5">
      <slot />
    </div>
  </section>
</template>
