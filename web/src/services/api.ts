const TOKEN_KEY = 'token'

export interface Script {
  id: number
  name: string
  description: string
  script?: string
  install_url: string
  created_at: string
}

export function getToken(): string | null {
  return localStorage.getItem(TOKEN_KEY)
}

export function setToken(token: string | null): void {
  if (token) {
    localStorage.setItem(TOKEN_KEY, token)
  } else {
    localStorage.removeItem(TOKEN_KEY)
  }
}

async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  const headers: Record<string, string> = {
    ...(options.headers as Record<string, string> | undefined),
  }
  const token = getToken()
  if (token) {
    headers.Authorization = `Bearer ${token}`
  }

  const res = await fetch(path, { ...options, headers })

  if (res.status === 401) {
    setToken(null)
    window.location.href = '/'
    throw new Error('Session expired, please sign in again')
  }
  if (!res.ok) {
    const body = (await res.json().catch(() => ({}))) as { error?: string }
    throw new Error(body.error || `Request failed (${res.status})`)
  }
  if (res.status === 204) {
    return undefined as T
  }
  return res.json() as Promise<T>
}

export const api = {
  login: (password: string) =>
    request<{ token: string }>('/api/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ password }),
    }),
  logout: () => request<void>('/api/logout', { method: 'POST' }),
  listScripts: () => request<Script[]>('/api/scripts'),
  getScript: (id: number) => request<Script>(`/api/scripts/${id}`),
  createScript: (input: { name: string; description: string; script: string }) =>
    request<Script>('/api/scripts', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(input),
    }),
  deleteScript: (id: number) =>
    request<void>(`/api/scripts/${id}`, { method: 'DELETE' }),
}
