const TOKEN_KEY = 'mathua_token'
const USER_KEY = 'mathua_user'
const GUEST_KEY = 'mathua_guest_id'
const GUEST_TOKEN_KEY = 'mathua_guest_token'
const API_BASE = process.env.NEXT_PUBLIC_API_URL || ''

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
  const token = getToken() || getGuestToken()
  if (token) {
    return { Authorization: `Bearer ${token}` }
  }
  return {}
}

// authedFetch is the single choke point for API calls: it attaches the
// registered token first (guest token as fallback, unless the caller set
// its own Authorization header) and, when the server answers 401, clears
// whichever stored token was actually sent — so a dead guest token never
// nukes a real login and vice versa. Callers that set their own header
// (e.g. admin login) never trigger a clear.
// Network failures reject — callers must not treat them as "invalid".
export async function authedFetch(input: RequestInfo | URL, init: RequestInit = {}): Promise<Response> {
  const registered = getToken()
  const guest = !registered ? getGuestToken() : null
  const sent = init.headers instanceof Headers
    ? init.headers.get('Authorization')
    : new Headers(init.headers).get('Authorization')
  const headers = new Headers(init.headers)
  if (!sent) {
    const token = registered || guest
    if (token) {
      headers.set('Authorization', `Bearer ${token}`)
    }
  }
  const res = await fetch(input, { ...init, headers })
  if (res.status === 401) {
    // Clear the stored credential behind the 401 — but only when the failed
    // request actually carried it. A caller-set header that merely forwards
    // our stored token (e.g. validateToken) still clears; a foreign bearer
    // (e.g. an admin triage token) never nukes the user's login.
    const failed = sent || ((registered || guest) && `Bearer ${registered || guest}`)
    if (failed && registered && failed === `Bearer ${registered}`) {
      clearToken()
    } else if (failed && !registered && guest && failed === `Bearer ${guest}`) {
      clearGuestToken()
    }
  }
  return res
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
  email?: string
  email_verified?: boolean
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

export function getGuestToken(): string | null {
  if (typeof window === 'undefined') return null
  return localStorage.getItem(GUEST_TOKEN_KEY)
}

export function clearGuestToken() {
  if (typeof window === 'undefined') return
  localStorage.removeItem(GUEST_TOKEN_KEY)
}

// ensureGuestToken mints (or reclaims) a server bearer token for the local
// guest id and stores it under its own key — registered logins are never
// touched and isLoggedIn() keeps meaning "has a real account". Uses plain
// fetch (not authedFetch) so a 401 here can never clear anything.
// Best-effort: resolves null when the server is unreachable; callers fall
// back to the anonymous guest_id flow.
export async function ensureGuestToken(): Promise<string | null> {
  if (typeof window === 'undefined') return null
  if (getToken()) return null // registered users don't need one
  const existing = getGuestToken()
  if (existing) return existing
  const guestId = ensureGuestId()
  try {
    const res = await fetch(`${API_BASE}/api/auth/guest`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ student_id: guestId }),
    })
    if (!res.ok) return null
    const data = (await res.json()) as { token?: string; student_id?: string }
    if (!data.token || !data.student_id) return null
    localStorage.setItem(GUEST_TOKEN_KEY, data.token)
    localStorage.setItem(GUEST_KEY, data.student_id)
    window.dispatchEvent(new Event('auth-changed'))
    return data.token
  } catch {
    return null
  }
}
