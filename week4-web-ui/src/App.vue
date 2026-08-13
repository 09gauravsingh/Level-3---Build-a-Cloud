<script setup>
import { onMounted, ref } from 'vue'
import { ApiError, NetworkError, api } from '@/api/client'
import { useToasts } from '@/composables/useToasts'
import AppHeader from '@/components/AppHeader.vue'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'
import ConnectionModal from '@/components/ConnectionModal.vue'
import CreateInstanceForm from '@/components/CreateInstanceForm.vue'
import InstanceList from '@/components/InstanceList.vue'
import LoginView from '@/components/LoginView.vue'
import PlatformStats from '@/components/PlatformStats.vue'
import ToastStack from '@/components/ui/ToastStack.vue'

const toasts = useToasts()

// The JWT is kept in localStorage so a page reload stays signed in.
const token = ref(localStorage.getItem('token') || '')
const loggedIn = ref(!!token.value)
const username = ref(subjectOf(token.value))
const loginError = ref('')

const instances = ref([])
const connection = ref(null)
const connectionName = ref('')
const pendingDelete = ref('')

const loading = ref({ login: false, list: false, create: false, delete: false })
const busyName = ref('')

const createForm = ref(null)

// Reads the "sub" claim so the header can show who is signed in.
function subjectOf(jwt) {
  try {
    return JSON.parse(atob(jwt.split('.')[1])).sub || 'admin'
  } catch {
    return 'admin'
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

    username.value = subjectOf(token.value)
    loggedIn.value = true

    await loadInstances()
  } catch (error) {
    loginError.value =
      error instanceof NetworkError ? error.message : 'Invalid username or password'
  } finally {
    loading.value.login = false
  }
}

function logout() {
  localStorage.removeItem('token')

  token.value = ''
  loggedIn.value = false
  instances.value = []
  connection.value = null
  pendingDelete.value = ''
}

async function loadInstances() {
  loading.value.list = true

  try {
    instances.value = await api.listInstances(token.value)
  } catch (error) {
    handleError(error, 'Failed to load instances')
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

    await loadInstances()
  } catch (error) {
    handleError(error, 'Could not delete instance')
  } finally {
    loading.value.delete = false
  }
}

onMounted(() => {
  if (loggedIn.value) {
    loadInstances()
  }
})
</script>

<template>
  <div class="app-backdrop min-h-screen">
    <ToastStack />

    <LoginView v-if="!loggedIn" :error="loginError" :loading="loading.login" @submit="login" />

    <template v-else>
      <AppHeader :username="username" @logout="logout" />

      <main class="mx-auto max-w-6xl space-y-6 px-4 py-8 sm:px-6">
        <PlatformStats :instances="instances" />

        <div class="grid gap-6 lg:grid-cols-5">
          <div class="lg:col-span-2">
            <CreateInstanceForm
              ref="createForm"
              :loading="loading.create"
              @submit="createInstance"
            />
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
      </main>
    </template>

    <ConnectionModal
      v-if="connection"
      :connection="connection"
      :instance-name="connectionName"
      @close="connection = null"
    />

    <ConfirmDeleteDialog
      v-if="pendingDelete"
      :instance-name="pendingDelete"
      :loading="loading.delete"
      @cancel="pendingDelete = ''"
      @confirm="deleteInstance"
    />
  </div>
</template>
