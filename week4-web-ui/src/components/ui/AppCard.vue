<script setup>
defineProps({
  title: { type: String, default: '' },
  subtitle: { type: String, default: '' },
})
</script>

<template>
  <section
    class="rounded-2xl border border-white/10 bg-slate-900/50 shadow-xl shadow-slate-950/40 backdrop-blur"
  >
    <header
      v-if="title || $slots.actions"
      class="flex flex-wrap items-center justify-between gap-3 border-b border-white/10 px-5 py-4"
    >
      <div>
        <h2 class="text-sm font-semibold text-slate-100">{{ title }}</h2>
        <p v-if="subtitle" class="mt-0.5 text-xs text-slate-500">{{ subtitle }}</p>
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
