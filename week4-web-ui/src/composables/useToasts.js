import { ref } from 'vue'

const toasts = ref([])
let nextId = 0

function push(type, message, timeout = 4500) {
  const id = nextId++
  toasts.value.push({ id, type, message })

  setTimeout(() => dismiss(id), timeout)

  return id
}

function dismiss(id) {
  toasts.value = toasts.value.filter((toast) => toast.id !== id)
}

export function useToasts() {
  return {
    toasts,
    dismiss,
    success: (message) => push('success', message),
    error: (message) => push('error', message),
    info: (message) => push('info', message),
  }
}
