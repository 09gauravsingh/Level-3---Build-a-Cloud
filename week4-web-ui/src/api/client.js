// Thin wrapper around the week3 Go API.
// Every endpoint here mirrors a route registered in internal/api/routes.go.

// Requests go to same-origin paths like /api/v1/login so the bundle never
// hardcodes where the backend runs: in dev the Vite proxy forwards them, in
// production the Ingress routes /api to the API Service.
const BASE_URL = import.meta.env.VITE_API_BASE_URL || ''

// ApiError carries the HTTP status so callers can react to 401 separately.
export class ApiError extends Error {
  constructor(status, message) {
    super(message)
    this.name = 'ApiError'
    this.status = status
  }

  get isUnauthorized() {
    return this.status === 401
  }
}

// NetworkError means the browser never reached the API.
export class NetworkError extends Error {
  constructor(message = 'Backend is not reachable') {
    super(message)
    this.name = 'NetworkError'
  }
}

// The Go API answers errors with models.ErrorResponse {error, message}.
async function readErrorMessage(response, fallback) {
  try {
    const body = await response.json()
    return body?.message || body?.error || fallback
  } catch {
    return fallback
  }
}

async function request(path, { method = 'GET', token = '', body, errorMessage } = {}) {
  const headers = {}

  if (body !== undefined) {
    headers['Content-Type'] = 'application/json'
  }

  if (token) {
    headers.Authorization = `Bearer ${token}`
  }

  let response
  try {
    response = await fetch(`${BASE_URL}${path}`, {
      method,
      headers,
      body: body === undefined ? undefined : JSON.stringify(body),
    })
  } catch {
    throw new NetworkError()
  }

  if (!response.ok) {
    throw new ApiError(response.status, await readErrorMessage(response, errorMessage))
  }

  if (response.status === 204) {
    return null
  }

  return response.json()
}

export const api = {
  // POST /api/v1/login -> { token }
  async login(username, password) {
    const data = await request('/api/v1/login', {
      method: 'POST',
      body: { username, password },
      errorMessage: 'Invalid username or password',
    })

    return data.token
  },

  // GET /api/v1/instances -> models.InstanceList { items, count }
  async listInstances(token) {
    const data = await request('/api/v1/instances', {
      token,
      errorMessage: 'Failed to load instances',
    })

    return data?.items ?? []
  },

  // POST /api/v1/instances -> 202 with the created models.Instance
  createInstance(token, payload) {
    return request('/api/v1/instances', {
      method: 'POST',
      token,
      body: payload,
      errorMessage: 'Could not create instance',
    })
  },

  // DELETE /api/v1/instances/:name -> 202 { name, status }
  deleteInstance(token, name) {
    return request(`/api/v1/instances/${encodeURIComponent(name)}`, {
      method: 'DELETE',
      token,
      errorMessage: 'Could not delete instance',
    })
  },

  // GET /api/v1/instances/:name/connection -> models.ConnectionData
  getConnection(token, name) {
    return request(`/api/v1/instances/${encodeURIComponent(name)}/connection`, {
      token,
      errorMessage: 'Could not get connection data',
    })
  },
}
