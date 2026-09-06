import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { setToken, clearToken, getToken, authedFetch } from '../lib/auth'
import { validateToken, login } from '../lib/api'

function mockFetchOnce(res: Partial<Response> & { json?: () => Promise<unknown> }) {
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
})
