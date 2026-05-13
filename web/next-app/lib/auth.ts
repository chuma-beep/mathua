const TOKEN_KEY = 'mathua_token'
const USER_KEY = 'mathua_user'

let _token: string | null = null
if (typeof window !== 'undefined') {
  _token = localStorage.getItem(TOKEN_KEY)
}

export function getToken(): string | null {
  return _token
}

export function setToken(token: string) {
  _token = token
  localStorage.setItem(TOKEN_KEY, token)
}

export function clearToken() {
  _token = null
  localStorage.removeItem(TOKEN_KEY)
  localStorage.removeItem(USER_KEY)
}

export function isLoggedIn(): boolean {
  return _token !== null
}

export function getAuthHeaders(): Record<string, string> {
  if (_token) {
    return { Authorization: `Bearer ${_token}` }
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
}

export function setUserInfo(info: UserInfo) {
  localStorage.setItem(USER_KEY, JSON.stringify(info))
}

export function getUserInfo(): UserInfo | null {
  const raw = localStorage.getItem(USER_KEY)
  if (!raw) return null
  try { return JSON.parse(raw) } catch { return null }
}
