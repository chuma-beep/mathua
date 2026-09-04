const TOKEN_KEY = 'mathua_token'
const USER_KEY = 'mathua_user'
const GUEST_KEY = 'mathua_guest_id'

export function getToken(): string | null {
  if (typeof window === 'undefined') return null
  return localStorage.getItem(TOKEN_KEY)
}

export function setToken(token: string) {
  localStorage.setItem(TOKEN_KEY, token)
  window.dispatchEvent(new Event('auth-changed'))
}

export function clearToken() {
  localStorage.removeItem(TOKEN_KEY)
  localStorage.removeItem(USER_KEY)
  window.dispatchEvent(new Event('auth-changed'))
}

export function isLoggedIn(): boolean {
  if (typeof window === 'undefined') return false
  return localStorage.getItem(TOKEN_KEY) !== null
}

export function getAuthHeaders(): Record<string, string> {
  const token = getToken()
  if (token) {
    return { Authorization: `Bearer ${token}` }
  }
  return {}
}

export interface UserInfo {
  student_id: string
  name: string
  username: string
  concepts_mastered: number
  current_streak: number
  level: string
  diagnostic_completed: boolean
  avatar_url?: string
}

export function setUserInfo(info: UserInfo) {
  localStorage.setItem(USER_KEY, JSON.stringify(info))
  window.dispatchEvent(new Event('auth-changed'))
}

export function signOut() {
  localStorage.removeItem(TOKEN_KEY)
  localStorage.removeItem(USER_KEY)
  // preserve guest progress
  try {
    const g = (window as unknown as { google?: { accounts?: { id?: { disableAutoSelect?: () => void } } } }).google
    g?.accounts?.id?.disableAutoSelect?.()
  } catch { /* ignore */ }
  window.dispatchEvent(new Event('auth-changed'))
}

export function getUserInfo(): UserInfo | null {
  const raw = localStorage.getItem(USER_KEY)
  if (!raw) return null
  try { return JSON.parse(raw) } catch { return null }
}

export function getGuestId(): string | null {
  if (typeof window === 'undefined') return null
  return localStorage.getItem(GUEST_KEY)
}

export function setGuestId(id: string) {
  localStorage.setItem(GUEST_KEY, id)
}

export function ensureGuestId(): string | null {
  if (typeof window === 'undefined') return null
  let id = localStorage.getItem(GUEST_KEY)
  if (!id) {
    // ephemeral guest id — persisted locally so /profile works for guests
    id = `guest_${Math.random().toString(36).slice(2, 10)}_${Date.now().toString(36)}`
    localStorage.setItem(GUEST_KEY, id)
  }
  return id
}
