<script setup>
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { ApiError, NetworkError, api, getAuditLogs } from '@/api/client'
import { useToasts } from '@/composables/useToasts'
import AppHeader from '@/components/AppHeader.vue'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'
import ConnectionModal from '@/components/ConnectionModal.vue'
import CreateInstanceForm from '@/components/CreateInstanceForm.vue'
import InstanceList from '@/components/InstanceList.vue'
import LoginView from '@/components/LoginView.vue'
import LogsWindow from '@/components/LogsWindow.vue'
import PlatformStats from '@/components/PlatformStats.vue'
import ToastStack from '@/components/ui/ToastStack.vue'

const toasts = useToasts()

// The JWT is kept in localStorage so a page reload stays signed in.
const token = ref(localStorage.getItem('token') || '')
const loggedIn = ref(!!token.value)
const session = ref(claimsOf(token.value))
const username = computed(() => session.value.username)
const isAdmin = computed(() => session.value.isAdmin)
const loginError = ref('')

const instances = ref([])
const connection = ref(null)
const connectionName = ref('')
const pendingDelete = ref('')
const logInstance = ref('')
const instanceLogs = ref([])
const auditLogs = ref([])
const view = ref('dashboard')

const loading = ref({
  login: false,
  register: false,
  list: false,
  create: false,
  delete: false,
  logs: false,
})
const busyName = ref('')
const lastUpdated = ref(null)

// Clusters take a while to settle, so the list keeps itself fresh.
const REFRESH_INTERVAL_MS = 10000
let refreshTimer = null

const createForm = ref(null)

const updatedLabel = computed(() =>
  lastUpdated.value
    ? lastUpdated.value.toLocaleTimeString(undefined, { timeStyle: 'medium' })
    : '',
)

// Reads the JWT payload so the header can show who is signed in.
function claimsOf(jwt) {
  try {
    const payload = JSON.parse(atob(jwt.split('.')[1]))

    return {
      username: payload.sub || '',
      isAdmin: payload.role === 'admin',
    }
  } catch {
    return { username: '', isAdmin: false }
  }
}

// Any 401 means the JWT expired: drop it and return to the login screen.
function handleError(error, fallback) {
  if (error instanceof ApiError && error.isUnauthorized) {
    logout()
    toasts.error('Session expired, please sign in again')
    return
  }

  toasts.error(error?.message || fallback)
}

async function login({ username: name, password }) {
  loginError.value = ''
  loading.value.login = true

  try {
    token.value = await api.login(name, password)
    localStorage.setItem('token', token.value)

    session.value = claimsOf(token.value)
    loggedIn.value = true

    await loadInstances()
  } catch (error) {
    loginError.value =
      error instanceof NetworkError ? error.message : 'Invalid username or password'
  } finally {
    loading.value.login = false
  }
}

async function register({ username: name, password }) {
  loginError.value = ''
  loading.value.register = true

  try {
    await api.register(name, password)
    toasts.success(`Account “${name}” created`)
    await login({ username: name, password })
  } catch (error) {
    loginError.value =
      error instanceof NetworkError ? error.message : error?.message || 'Could not create account'
  } finally {
    loading.value.register = false
  }
}

function logout() {
  localStorage.removeItem('token')

  token.value = ''
  loggedIn.value = false
  session.value = { username: '', isAdmin: false }
  instances.value = []
  connection.value = null
  pendingDelete.value = ''
  logInstance.value = ''
  instanceLogs.value = []
  auditLogs.value = []
  view.value = 'dashboard'
  lastUpdated.value = null
}

async function loadInstances({ silent = false } = {}) {
  if (!silent) {
    loading.value.list = true
  }

  try {
    instances.value = await api.listInstances(token.value)
    lastUpdated.value = new Date()
  } catch (error) {
    // A silent poll stays quiet unless the session itself has expired.
    if (!silent || (error instanceof ApiError && error.isUnauthorized)) {
      handleError(error, 'Failed to load instances')
    }
  } finally {
    loading.value.list = false
  }
}

async function createInstance(payload) {
  loading.value.create = true

  try {
    const instance = await api.createInstance(token.value, payload)

    toasts.success(`Instance “${instance?.name ?? payload.name}” is being provisioned`)
    createForm.value?.reset()

    await loadInstances()
  } catch (error) {
    handleError(error, 'Could not create instance')
  } finally {
    loading.value.create = false
  }
}

async function getConnectionData(name) {
  busyName.value = name

  try {
    connection.value = await api.getConnection(token.value, name)
    connectionName.value = name
  } catch (error) {
    handleError(error, 'Could not get connection data')
  } finally {
    busyName.value = ''
  }
}

async function deleteInstance() {
  const name = pendingDelete.value
  loading.value.delete = true

  try {
    await api.deleteInstance(token.value, name)

    pendingDelete.value = ''
    toasts.success(`Instance “${name}” is being deleted`)

    if (logInstance.value === name) {
      logInstance.value = ''
      instanceLogs.value = []
      auditLogs.value = []
    }

    await loadInstances()
  } catch (error) {
    handleError(error, 'Could not delete instance')
  } finally {
    loading.value.delete = false
  }
}

async function loadLogWindow(name = logInstance.value) {
  if (!name) {
    instanceLogs.value = []
    auditLogs.value = []
    return
  }

  loading.value.logs = true

  try {
    const [postgres, audit] = await Promise.all([
      api.getInstanceLogs(token.value, name),
      getAuditLogs(name, token.value),
    ])

    instanceLogs.value = postgres.logs || []
    auditLogs.value = audit.activity || []
  } catch (error) {
    handleError(error, 'Could not load logs')
  } finally {
    loading.value.logs = false
  }
}

function selectLogInstance(name) {
  logInstance.value = name
  instanceLogs.value = []
  auditLogs.value = []
  loadLogWindow(name)
}

function setView(next) {
  view.value = next

  if (next === 'logs' && logInstance.value) {
    loadLogWindow()
  }
}

watch(
  instances,
  (list) => {
    if (logInstance.value && !list.some((instance) => instance.name === logInstance.value)) {
      logInstance.value = ''
      instanceLogs.value = []
      auditLogs.value = []
    }
  },
)

onMounted(() => {
  if (loggedIn.value) {
    loadInstances()
  }

  refreshTimer = setInterval(() => {
    if (loggedIn.value && !document.hidden && !loading.value.list) {
      loadInstances({ silent: true })
    }
  }, REFRESH_INTERVAL_MS)
})

onBeforeUnmount(() => clearInterval(refreshTimer))
</script>

<template>
  <div class="app-backdrop min-h-screen">
    <ToastStack />

    <LoginView v-if="!loggedIn" :error="loginError" :loading="loading.login || loading.register" @submit="login"
      @register="register" />

    <template v-else>
      <AppHeader :username="username" :is-admin="isAdmin" :view="view" @logout="logout" @navigate="setView" />

      <main class="mx-auto max-w-7xl px-4 py-8 sm:px-6">
        <template v-if="view === 'dashboard'">
          <div class="space-y-6">
            <div class="flex flex-wrap items-end justify-between gap-3">
              <div>
                <h1 class="text-2xl font-semibold tracking-tight text-slate-900">Dashboard</h1>
                <p class="mt-1 text-sm text-slate-500">
                  {{
                    isAdmin
                      ? 'Administrator view of every PostgreSQL cluster on the platform.'
                      : 'Provision and manage the PostgreSQL clusters that belong to you.'
                  }}
                </p>
              </div>

              <p
                v-if="lastUpdated"
                class="inline-flex items-center gap-1.5 rounded-full bg-white px-3 py-1 text-xs text-slate-500 ring-1 ring-slate-200"
              >
                <span class="size-1.5 rounded-full bg-emerald-500" aria-hidden="true" />
                Updated {{ updatedLabel }}
              </p>
            </div>

            <PlatformStats :instances="instances" />

            <div class="grid items-start gap-6 lg:grid-cols-5">
              <div class="lg:sticky lg:top-20 lg:col-span-2">
                <CreateInstanceForm ref="createForm" :loading="loading.create" @submit="createInstance" />
              </div>

              <div class="lg:col-span-3">
                <InstanceList
                  :instances="instances"
                  :loading="loading.list"
                  :busy-name="busyName"
                  @refresh="loadInstances"
                  @connection="getConnectionData"
                  @delete="pendingDelete = $event"
                />
              </div>
            </div>
          </div>
        </template>

        <LogsWindow
          v-else
          :instances="instances"
          :selected-name="logInstance"
          :instance-logs="instanceLogs"
          :audit-logs="auditLogs"
          :loading="loading.logs"
          @select="selectLogInstance"
          @refresh="loadLogWindow()"
        />

        <footer class="mt-8 border-t border-slate-200 pt-6 text-xs text-slate-400">
          PostgreSQL PaaS · CloudNativePG on Kubernetes · list refreshes every 10 seconds
        </footer>
      </main>
    </template>

    <ConnectionModal v-if="connection" :connection="connection" :instance-name="connectionName"
      @close="connection = null" />

    <ConfirmDeleteDialog v-if="pendingDelete" :instance-name="pendingDelete" :loading="loading.delete"
      @cancel="pendingDelete = ''" @confirm="deleteInstance" />
  </div>
</template>
