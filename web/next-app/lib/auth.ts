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
  const headers = token ? { Authorization: `Bearer ${token}` } : {}
  return headers
}

// authedFetch is the single choke point for API calls: it attaches the
// registered token first (guest token as fallback, unless the caller set
// its own Authorization header) and, when the server answers 401, clears
// whichever stored token was actually sent — so a dead guest token never
// nukes a real login and vice versa. Callers that set their own header
// (e.g. admin login) never trigger a clear.
// Network failures reject — callers must not treat them as "invalid".
// Per-request timeout override: authedFetch(input, { timeoutMs: 5000 }).
// Stripped before reaching fetch.
export type FetchInit = RequestInit & { timeoutMs?: number }

const DEFAULT_TIMEOUT_MS = 30000

export async function authedFetch(input: RequestInfo | URL, init: FetchInit = {}): Promise<Response> {
  const registered = getToken()
  const guest = !registered ? getGuestToken() : null
  const sent = init.headers instanceof Headers
    ? init.headers.get('Authorization')
    : new Headers(init.headers).get('Authorization')
  const { timeoutMs, ...fetchInit } = init
  const headers = new Headers(fetchInit.headers)
  if (!sent) {
    const token = registered || guest
    if (token) {
      headers.set('Authorization', `Bearer ${token}`)
    }
  }
  const ms = timeoutMs ?? DEFAULT_TIMEOUT_MS
  const ctrl = new AbortController()
  const timer = setTimeout(() => ctrl.abort(), ms)
  // A caller-provided signal still wins: aborting it aborts ours.
  if (fetchInit.signal) {
    const caller = fetchInit.signal
    if (caller.aborted) ctrl.abort()
    else caller.addEventListener('abort', () => ctrl.abort(), { once: true })
  }
  let res: Response
  try {
    res = await fetch(input, { ...fetchInit, headers, signal: ctrl.signal })
  } finally {
    clearTimeout(timer)
  }
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
  has_password?: boolean
  /**
   * The account's administrative role, as reported by `/api/auth/me`.
   *
   * Present so the app can *offer* an Admin entry, and nothing more. A missing field — an old
   * build, a stale cached bundle — reads as a learner, which is the safe direction: the pages
   * are still reachable by URL and the server refuses them.
   */
  role?: 'student' | 'admin'
}

export function setUserInfo(info: UserInfo) {
  localStorage.setItem(USER_KEY, JSON.stringify(info))
  window.dispatchEvent(new Event('auth-changed'))
}

interface GoogleAccountsWindow {
  google?: { accounts?: { id?: { disableAutoSelect?: () => void } } }
}

export function signOut() {
  localStorage.removeItem(TOKEN_KEY)
  localStorage.removeItem(USER_KEY)
  // preserve guest progress
  try {
    const g = (window as GoogleAccountsWindow).google
    g?.accounts?.id?.disableAutoSelect?.()
  } catch { /* ignore */ }
  window.dispatchEvent(new Event('auth-changed'))
}

// clearGuest wipes guest identity + session (used after account deletion,
// where nothing — including the guest row — may survive).
export function clearGuest() {
  localStorage.removeItem(GUEST_KEY)
  localStorage.removeItem(GUEST_TOKEN_KEY)
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
    // Ephemeral guest id — persisted locally so /profile works for guests.
    // crypto.randomUUID when available (unguessable capability, doubling as
    // the claim key for POST /api/auth/guest); legacy Math.random fallback.
    const rand = typeof crypto !== 'undefined' && typeof crypto.randomUUID === 'function'
      ? crypto.randomUUID()
      : `${Math.random().toString(36).slice(2, 10)}_${Date.now().toString(36)}`
    id = `guest_${rand}`
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
  const ctrl = new AbortController()
  const timer = setTimeout(() => ctrl.abort(), 10000)
  try {
    const res = await fetch(`${API_BASE}/api/auth/guest`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ student_id: guestId }),
      signal: ctrl.signal,
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
  } finally {
    clearTimeout(timer)
  }
}

/**
 * Store a freshly issued token together with the account the server says it belongs to.
 *
 * Every sign-in path used to build its own `UserInfo` literal from whatever the login response
 * happened to contain — four copies, in the Google callback, One-Tap, password reset and the
 * password form, each filling `concepts_mastered` and `current_streak` with zeros because the
 * login response does not carry them. That is four copies of a schema to keep in step with
 * `/api/auth/me`, and it immediately cost something: `role` is reported by that endpoint, no
 * literal mentioned it, so a signed-in administrator's Admin entry never appeared at all.
 *
 * One function, one fetch of the server's own payload. The request is not an extra cost on the
 * password path, which was already calling `/api/auth/me` through `validateToken` to confirm the
 * token; this replaces that call rather than adding to it.
 *
 * Fails soft, in the direction that matters. If `/api/auth/me` fails the token is kept — the
 * login itself succeeded — and the partial object is stored, which simply has no `role`. An
 * administrator whose profile fetch blipped does not see the Admin entry until it succeeds, and
 * a learner never sees it at all. Neither loses access to anything: the server decides.
 */
export async function establishSession(
  token: string,
  partial: Omit<UserInfo, 'role'>,
): Promise<UserInfo> {
  setToken(token)
  try {
    const res = await fetch(`${API_BASE}/api/auth/me`, {
      headers: { Authorization: `Bearer ${token}` },
    })
    if (res.ok) {
      const me = (await res.json()) as Partial<UserInfo>
      // Built in statements rather than with a conditional spread, because the spread is
      // exactly the thing that hides the case this function exists to get right: `role:
      // undefined` and an absent key both read as falsy, so a spread would make "the server did
      // not send a role" indistinguishable from "the server sent student".
      const merged: UserInfo = {
        ...partial,
        student_id: typeof me.student_id === 'string' ? me.student_id : partial.student_id,
        name: typeof me.name === 'string' ? me.name : partial.name,
      }
      if (me.role === 'admin' || me.role === 'student') merged.role = me.role
      if (typeof me.has_password === 'boolean') merged.has_password = me.has_password
      setUserInfo(merged)
      return merged
    }
  } catch {
    // Offline or unreachable. See above: keep the session, drop the enrichment.
  }
  setUserInfo({ ...partial })
  return { ...partial }
}
