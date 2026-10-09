import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { setToken, clearToken, getToken, authedFetch, getGuestToken, ensureGuestToken, getAuthHeaders, establishSession, fetchMe, getUserInfo, setUserInfo, isStaffRole } from '../lib/auth'
import { validateToken, login } from '../lib/api'

function mockFetchOnce(res: Partial<Response> & { json?: () => Promise<object> }) {
  const fn = vi.fn().mockResolvedValue({
    ok: true,
    status: 200,
    json: async () => ({}),
    ...res,
  })
  vi.stubGlobal('fetch', fn)
  return fn
}

describe('login persistence', () => {
  beforeEach(() => {
    localStorage.clear()
    vi.unstubAllGlobals()
  })

  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('validateToken hits the canonical /api/auth/me route', async () => {
    setToken('tok123')
    const fetchMock = mockFetchOnce({ json: async () => ({ student_id: 's1' }) })
    const v = await validateToken()
    expect(v).toEqual({ valid: true, student_id: 's1' })
    expect(fetchMock).toHaveBeenCalledOnce()
    const [url, init] = fetchMock.mock.calls[0] as [string, RequestInit]
    expect(url).toContain('/api/auth/me')
    expect(new Headers(init.headers).get('Authorization')).toBe('Bearer tok123')
  })

  it('validateToken reports invalid on 401 and clears the dead token', async () => {
    setToken('dead')
    mockFetchOnce({ ok: false, status: 401 })
    const v = await validateToken()
    expect(v.valid).toBe(false)
    expect(getToken()).toBeNull()
  })

  it('validateToken rejects on network failure without clearing the token', async () => {
    setToken('tok123')
    vi.stubGlobal('fetch', vi.fn().mockRejectedValue(new Error('offline')))
    await expect(validateToken()).rejects.toThrow('offline')
    expect(getToken()).toBe('tok123')
  })

  it('validateToken is invalid without a fetch when no token is stored', async () => {
    const fetchMock = vi.fn()
    vi.stubGlobal('fetch', fetchMock)
    await expect(validateToken()).resolves.toEqual({ valid: false, student_id: '' })
    expect(fetchMock).not.toHaveBeenCalled()
  })

  it('authedFetch attaches the token and keeps a caller-set header', async () => {
    setToken('tok123')
    const fetchMock = mockFetchOnce({})
    await authedFetch('https://x/api/y')
    expect(new Headers((fetchMock.mock.calls[0] as [string, RequestInit])[1].headers).get('Authorization')).toBe('Bearer tok123')

    await authedFetch('https://x/api/y', { headers: { Authorization: 'Bearer other' } })
    const second = fetchMock.mock.calls[1] as [string, RequestInit]
    expect(new Headers(second[1].headers).get('Authorization')).toBe('Bearer other')
  })

  it('authedFetch never clears storage for a 401 on an anonymous call', async () => {
    clearToken()
    mockFetchOnce({ ok: false, status: 401 })
    await authedFetch('https://x/api/public')
    expect(getToken()).toBeNull()
  })

  it('login maps 429 to a friendly retry message', async () => {
    mockFetchOnce({ ok: false, status: 429 })
    await expect(login('ada', 'Engine!n1')).rejects.toThrow('Too many attempts')
  })

  it('authedFetch falls back to the guest token and clears only it on 401', async () => {
    clearToken()
    localStorage.setItem('mathua_guest_token', 'guest-tok')
    const fetchMock = mockFetchOnce({})
    await authedFetch('https://x/api/scores/guest_1')
    expect(new Headers((fetchMock.mock.calls[0] as [string, RequestInit])[1].headers).get('Authorization')).toBe('Bearer guest-tok')

    mockFetchOnce({ ok: false, status: 401 })
    await authedFetch('https://x/api/scores/guest_1')
    expect(getGuestToken()).toBeNull()
    expect(getToken()).toBeNull()
  })

  it('registered token wins over guest token and 401 clears only the registered one', async () => {
    setToken('reg-tok')
    localStorage.setItem('mathua_guest_token', 'guest-tok')
    expect(getAuthHeaders()).toEqual({ Authorization: 'Bearer reg-tok' })
    mockFetchOnce({ ok: false, status: 401 })
    await authedFetch('https://x/api/y')
    expect(getToken()).toBeNull()
    expect(getGuestToken()).toBe('guest-tok')
  })

  it('ensureGuestToken mints via /api/auth/guest and stores token + id', async () => {
    clearToken()
    localStorage.removeItem('mathua_guest_token')
    localStorage.setItem('mathua_guest_id', 'guest_local99')
    const fetchMock = mockFetchOnce({ json: async () => ({ token: 'gt-1', student_id: 'guest_local99' }) })
    const tok = await ensureGuestToken()
    expect(tok).toBe('gt-1')
    expect(getGuestToken()).toBe('gt-1')
    const [url, init] = fetchMock.mock.calls[0] as [string, RequestInit]
    expect(url).toContain('/api/auth/guest')
    expect(JSON.parse(init.body as string)).toEqual({ student_id: 'guest_local99' })
  })

  it('ensureGuestToken is a no-op for registered users', async () => {
    setToken('reg-tok')
    const fetchMock = vi.fn()
    vi.stubGlobal('fetch', fetchMock)
    await expect(ensureGuestToken()).resolves.toBeNull()
    expect(fetchMock).not.toHaveBeenCalled()
  })

  it('ensureGuestId mints a prefixed unguessable id', async () => {
    const { ensureGuestId } = await import('../lib/auth')
    localStorage.removeItem('mathua_guest_id')
    const id = ensureGuestId()
    expect(id).toMatch(/^guest_/)
    expect(id!.length).toBeGreaterThan(20)
    expect(ensureGuestId()).toBe(id)
  })

  it('authedFetch aborts after timeoutMs', async () => {
    clearToken()
    localStorage.removeItem('mathua_guest_token')
    vi.stubGlobal('fetch', vi.fn((_url: RequestInfo | URL, init?: RequestInit) => new Promise((_res, rej) => {
      init?.signal?.addEventListener('abort', () => rej(new DOMException('aborted', 'AbortError')))
    })))
    await expect(authedFetch('https://x/api/slow', { timeoutMs: 30 })).rejects.toThrow()
  })

  it('getLessons goes through authedFetch with credentials', async () => {
    const { getLessons } = await import('../lib/api')
    clearToken()
    localStorage.setItem('mathua_guest_token', 'guest-tok')
    const fetchMock = mockFetchOnce({ json: async () => ({ lessons: {} }) })
    await getLessons('guest_1')
    const [, init] = fetchMock.mock.calls[0] as [string, RequestInit]
    expect(new Headers(init.headers).get('Authorization')).toBe('Bearer guest-tok')
  })
})

// The role has four values, and this function used to keep two of them.
//
// A `moderator` or an `owner` signing in on a phone had its role dropped here, so the stored
// session carried no `role`, `isStaffRole` answered false, and the Admin entry never appeared on
// that device — on the sidebar or in the compass menu — while a desktop that had re-authenticated
// showed it. The Admin link was therefore absent on exactly the devices nobody tests from.
describe('establishSession keeps every role', () => {
  beforeEach(() => {
    localStorage.clear()
    vi.unstubAllGlobals()
  })
  afterEach(() => vi.unstubAllGlobals())

  it.each(['student', 'moderator', 'admin', 'owner'])('stores %s', async (role) => {
    mockFetchOnce({ json: async () => ({ student_id: 's1', name: 'Root', role }) })
    const me = await establishSession('tok', {
      student_id: 's1', name: 'Root', username: 'root',
      concepts_mastered: 0, current_streak: 0, level: 'Novice', diagnostic_completed: true,
    })
    expect(me.role).toBe(role)
    expect(getUserInfo()?.role).toBe(role)
  })

  it('treats the two staff roles as staff', async () => {
    mockFetchOnce({ json: async () => ({ student_id: 's1', name: 'Root', role: 'owner' }) })
    await establishSession('tok', {
      student_id: 's1', name: 'Root', username: 'root',
      concepts_mastered: 0, current_streak: 0, level: 'Novice', diagnostic_completed: true,
    })
    expect(isStaffRole(getUserInfo()?.role)).toBe(true)
  })

  it('leaves the role unset when the server sends one this build does not know', async () => {
    // A future role must not be cached as a privilege nobody has audited.
    mockFetchOnce({ json: async () => ({ student_id: 's1', name: 'Root', role: 'superuser' }) })
    const me = await establishSession('tok', {
      student_id: 's1', name: 'Root', username: 'root',
      concepts_mastered: 0, current_streak: 0, level: 'Novice', diagnostic_completed: true,
    })
    expect(me.role).toBeUndefined()
  })
})

// A session written before roles existed must repair itself rather than stay without one
// forever. This is the phone: the session is old, the account is not.
describe('fetchMe repairs a session with no cached role', () => {
  beforeEach(() => {
    localStorage.clear()
    vi.unstubAllGlobals()
  })
  afterEach(() => vi.unstubAllGlobals())

  it('writes the role the server reports', async () => {
    setToken('tok')
    setUserInfo({
      student_id: 'a-root', name: 'Root', username: 'root',
      concepts_mastered: 0, current_streak: 0, level: 'Novice', diagnostic_completed: true,
    })
    expect(getUserInfo()?.role).toBeUndefined()

    mockFetchOnce({ json: async () => ({ student_id: 'a-root', name: 'Root', role: 'moderator' }) })
    const me = await fetchMe()
    expect(me?.role).toBe('moderator')
    expect(getUserInfo()?.role).toBe('moderator')
  })

  // "We could not ask" must never become "not staff": that would cache a wrong answer into the
  // stored session, which is exactly the failure this repair exists to undo.
  it('leaves the stored session untouched when the server cannot be reached', async () => {
    setToken('tok')
    setUserInfo({
      student_id: 'a-root', name: 'Root', username: 'root',
      concepts_mastered: 0, current_streak: 0, level: 'Novice', diagnostic_completed: true,
    })
    vi.stubGlobal('fetch', vi.fn().mockRejectedValue(new Error('offline')))
    expect(await fetchMe()).toBeNull()
    expect(getUserInfo()?.role).toBeUndefined()
  })

  it('does not cache a missing role as student', async () => {
    setToken('tok')
    mockFetchOnce({ json: async () => ({ student_id: 'a-root', name: 'Root' }) })
    await fetchMe()
    expect(getUserInfo()?.role).toBeUndefined()
  })
})
