<script setup>
defineProps({
  title: { type: String, default: '' },
  subtitle: { type: String, default: '' },
})
</script>

<template>
  <section class="rounded-2xl border border-slate-200 bg-white shadow-sm shadow-slate-900/5">
    <header
      v-if="title || $slots.actions"
      class="flex flex-wrap items-center justify-between gap-3 border-b border-slate-100 px-5 py-4"
    >
      <div class="flex items-center gap-3">
        <span
          v-if="$slots.icon"
          class="flex size-9 shrink-0 items-center justify-center rounded-lg bg-indigo-50 text-indigo-600"
        >
          <slot name="icon" />
        </span>

        <div>
          <h2 class="text-sm font-semibold text-slate-900">{{ title }}</h2>
          <p v-if="subtitle" class="mt-0.5 text-xs text-slate-500">{{ subtitle }}</p>
        </div>
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
