<script setup>
import { computed, ref } from 'vue'
import AppCard from '@/components/ui/AppCard.vue'
import BaseButton from '@/components/ui/BaseButton.vue'

const props = defineProps({
  instances: { type: Array, default: () => [] },
  selectedName: { type: String, default: '' },
  instanceLogs: { type: Array, default: () => [] },
  auditLogs: { type: Array, default: () => [] },
  loading: { type: Boolean, default: false },
})

const emit = defineEmits(['select', 'refresh'])

const pane = ref('instance')

const selected = computed({
  get: () => props.selectedName,
  set: (value) => emit('select', value),
})

function formatTime(value) {
  if (!value) return '—'

  const date = new Date(value)

  return Number.isNaN(date.getTime())
    ? value
    : date.toLocaleString(undefined, { dateStyle: 'short', timeStyle: 'medium' })
}

function displayLine(log) {
  const line = log?.line || ''
  const cri = /^\S+\s+(stdout|stderr)\s+\S+\s+(.*)$/s
  const match = line.match(cri)

  return match ? match[2] : line
}

function fieldsOf(log) {
  if (log?.action || log?.user || log?.result) {
    return {
      action: log.action || '',
      result: log.result || '',
      user: log.user || '',
      timestamp: log.timestamp || '',
    }
  }

  const line = log?.line || ''
  const start = line.indexOf('{')

  try {
    const parsed = JSON.parse(start >= 0 ? line.slice(start) : line)

    return {
      action: parsed.action || '',
      result: parsed.result || '',
      user: parsed.user || '',
      timestamp: log?.timestamp || parsed.time || '',
    }
  } catch {
    return {
      action: line,
      result: '',
      user: '',
      timestamp: log?.timestamp || '',
    }
  }
}

const auditEntries = computed(() => (props.auditLogs || []).map(fieldsOf))

const panes = [
  { id: 'instance', label: 'Instance logs', caption: 'PostgreSQL container output' },
  { id: 'audit', label: 'Audit logs', caption: 'User actions on this instance' },
]
</script>

<template>
  <div class="space-y-6">
    <div class="flex flex-wrap items-end justify-between gap-3">
      <div>
        <h1 class="text-2xl font-semibold tracking-tight text-slate-900">Logs</h1>
        <p class="mt-1 text-sm text-slate-500">
          Inspect PostgreSQL output and audit events for one instance.
        </p>
      </div>
    </div>

    <AppCard title="Query" subtitle="Choose an instance, then switch between log types">
      <template #actions>
        <label class="flex items-center gap-2 text-xs text-slate-600">
          <span>Instance</span>
          <select
            v-model="selected"
            class="h-8 min-w-44 rounded-lg border border-slate-200 bg-white px-2 font-mono text-xs text-slate-800 shadow-sm outline-none transition hover:border-slate-300 focus:border-indigo-500 focus:ring-4 focus:ring-indigo-500/10"
            aria-label="Select instance for logs"
          >
            <option value="" disabled>Select instance</option>
            <option v-for="instance in instances" :key="instance.name" :value="instance.name">
              {{ instance.name }}
            </option>
          </select>
        </label>

        <BaseButton
          variant="secondary"
          size="sm"
          :disabled="!selectedName"
          :loading="loading"
          @click="emit('refresh')"
        >
          Refresh
        </BaseButton>
      </template>

      <div class="mb-4 flex gap-1 rounded-lg bg-slate-100 p-1">
        <button
          v-for="item in panes"
          :key="item.id"
          type="button"
          class="flex-1 rounded-md px-3 py-2 text-left transition"
          :class="
            pane === item.id
              ? 'bg-white text-slate-900 shadow-sm'
              : 'text-slate-500 hover:text-slate-800'
          "
          :aria-pressed="pane === item.id"
          @click="pane = item.id"
        >
          <span class="block text-sm font-semibold">{{ item.label }}</span>
          <span class="mt-0.5 block text-[11px] text-slate-500">{{ item.caption }}</span>
        </button>
      </div>

      <div
        v-if="!selectedName"
        class="rounded-xl border border-dashed border-slate-200 bg-slate-50 px-4 py-16 text-center"
      >
        <p class="text-sm font-medium text-slate-800">Select an instance to load logs</p>
        <p class="mt-1 text-xs text-slate-500">
          Instance logs and audit logs are queried from Loki for the last hour.
        </p>
      </div>

      <template v-else-if="pane === 'instance'">
        <div
          v-if="loading && instanceLogs.length === 0"
          class="h-64 animate-pulse rounded-xl bg-slate-100"
        />

        <div
          v-else-if="instanceLogs.length === 0"
          class="rounded-xl border border-dashed border-slate-200 bg-slate-50 px-4 py-16 text-center text-sm text-slate-500"
        >
          No instance logs found.
        </div>

        <div v-else class="max-h-[32rem] overflow-auto rounded-xl border border-slate-200">
          <table class="min-w-full text-left">
            <thead
              class="sticky top-0 border-b border-slate-200 bg-slate-50 text-[11px] tracking-wide text-slate-500 uppercase"
            >
              <tr>
                <th class="px-4 py-2.5 font-medium">Time</th>
                <th class="px-4 py-2.5 font-medium">Message</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100 bg-white">
              <tr v-for="(log, index) in instanceLogs" :key="index" class="align-top">
                <td class="px-4 py-2 font-mono text-xs whitespace-nowrap text-slate-500">
                  {{ formatTime(log.timestamp) }}
                </td>
                <td class="px-4 py-2 font-mono text-xs leading-5 whitespace-pre-wrap text-slate-800">
                  {{ displayLine(log) }}
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </template>

      <template v-else>
        <div
          v-if="loading && auditEntries.length === 0"
          class="h-64 animate-pulse rounded-xl bg-slate-100"
        />

        <div
          v-else-if="auditEntries.length === 0"
          class="rounded-xl border border-dashed border-slate-200 bg-slate-50 px-4 py-16 text-center text-sm text-slate-500"
        >
          No audit logs found.
        </div>

        <div v-else class="max-h-[32rem] overflow-auto rounded-xl border border-slate-200">
          <table class="min-w-full text-left text-sm">
            <thead
              class="sticky top-0 border-b border-slate-200 bg-slate-50 text-[11px] tracking-wide text-slate-500 uppercase"
            >
              <tr>
                <th class="px-4 py-2.5 font-medium">Time</th>
                <th class="px-4 py-2.5 font-medium">User</th>
                <th class="px-4 py-2.5 font-medium">Action</th>
                <th class="px-4 py-2.5 font-medium">Result</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-slate-100 bg-white">
              <tr v-for="(log, index) in auditEntries" :key="index">
                <td class="px-4 py-2.5 font-mono text-xs whitespace-nowrap text-slate-600">
                  {{ formatTime(log.timestamp) }}
                </td>
                <td class="px-4 py-2.5 font-mono text-xs text-slate-800">{{ log.user || '—' }}</td>
                <td class="px-4 py-2.5 text-slate-800">{{ log.action || '—' }}</td>
                <td class="px-4 py-2.5">
                  <span
                    class="inline-flex rounded-full px-2 py-0.5 text-[11px] font-semibold"
                    :class="
                      String(log.result).toLowerCase() === 'success'
                        ? 'bg-emerald-50 text-emerald-700'
                        : String(log.result).toLowerCase() === 'failed'
                          ? 'bg-rose-50 text-rose-700'
                          : 'bg-slate-100 text-slate-600'
                    "
                  >
                    {{ log.result || '—' }}
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </template>
    </AppCard>
  </div>
</template>
