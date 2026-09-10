import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { setToken, clearToken, getToken, authedFetch, getGuestToken, ensureGuestToken, getAuthHeaders } from '../lib/auth'
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
