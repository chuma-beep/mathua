import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { establishSession, getUserInfo, getToken, clearToken, setToken } from '../lib/auth'
import { API_BASE, type MeInfo } from '../lib/api'

/**
 * What a sign-in leaves behind.
 *
 * The session object is what the whole app reads to know who it is talking to, and the role in
 * it is what decides whether the Admin entry is offered at all. It used to be assembled by hand
 * at four separate sign-in paths, and none of them mentioned `role` — so a signed-in
 * administrator's Admin link never appeared, while `/api/auth/me` had been reporting the role
 * correctly the whole time. Nothing failed: the login worked, the server agreed the account was
 * an administrator, and every administrative request the admin then made by URL worked too. Only
 * the invitation was missing, which is the worst kind of bug — the system looks like it does not
 * have the feature.
 *
 * These tests are about the session the client keeps. They say nothing about access, which is
 * settled server-side on every request and covered in Go.
 */

const PARTIAL = {
  student_id: 's1', name: 'Root', username: 'root',
  concepts_mastered: 0, current_streak: 0, level: 'Novice', diagnostic_completed: true,
} as const

/**
 * The stub's own signature, so `mock.calls` is typed as the argument tuple it is.
 *
 * `vi.fn()` with no type parameter records calls as `any[]`, and reading one back needs a
 * double assertion — which is both how the type evidence gets discarded and how this file ended
 * up asserting against a shape `fetch` never has. Typing the mock means the assertions below read
 * `call[0]` and `call[1].headers` and need no cast at all.
 */
/**
 * What the stub resolves to: enough of a Response for the helper, and nothing more.
 *
 * The body is typed as the real `/api/auth/me` contract rather than as `unknown`, because that
 * contract is what these tests are about — including its two absences, which is why every field
 * is optional here. A loose type would let a test pass against a body the endpoint could never
 * send.
 */
type MeBody = Partial<MeInfo>

interface StubResponse {
  ok: boolean
  status: number
  json: () => Promise<MeBody>
}

type FetchStub = (input: string, init?: RequestInit) => Promise<StubResponse>

function mockMe(body: MeBody, ok = true, status = 200) {
  const fn = vi.fn<FetchStub>().mockResolvedValue({ ok, status, json: async () => body })
  vi.stubGlobal('fetch', fn)
  return fn
}

beforeEach(() => {
  localStorage.clear()
  vi.unstubAllGlobals()
})

afterEach(() => {
  vi.unstubAllGlobals()
  localStorage.clear()
})

describe('establishSession', () => {
  it('stores the role the server reports', async () => {
    mockMe({ student_id: 's1', name: 'Root', role: 'admin', has_password: true })
    await establishSession('tok', { ...PARTIAL })
    expect(getUserInfo()?.role).toBe('admin')
    expect(getToken()).toBe('tok')
  })

  it('stores a learner role just as faithfully', async () => {
    // A field that is only read on the admin branch is a field nobody tested, which is how the
    // original gap survived: the learner path looked correct because it never looked.
    mockMe({ student_id: 's1', name: 'Ada', role: 'student' })
    await establishSession('tok', { ...PARTIAL })
    expect(getUserInfo()?.role).toBe('student')
  })

  it('reads the profile from the endpoint rather than trusting the caller', async () => {
    const fetchMock = mockMe({ student_id: 's1', name: 'Root', role: 'admin' })
    await establishSession('tok', { ...PARTIAL })
    // The URL is built from the same API_BASE every other call uses rather than hard-coded. In
    // tests API_BASE is '', so the expected value is the relative path.
    const [url, init] = fetchMock.mock.calls[0]
    expect(url).toBe(`${API_BASE}/api/auth/me`.replace(/^undefined/, ''))
    expect(new Headers(init?.headers).get('Authorization')).toBe('Bearer tok')
  })

  // The server's name and id win. The partial exists only so a failed fetch leaves something
  // coherent behind, and letting it override the server would make the two disagree silently.
  it('prefers the server values over the caller-supplied partial', async () => {
    mockMe({ student_id: 'server-id', name: 'Server Name', role: 'admin' })
    const stored = await establishSession('tok', { ...PARTIAL, student_id: 'stale-id', name: 'Stale' })
    expect(stored.student_id).toBe('server-id')
    expect(stored.name).toBe('Server Name')
  })

  it('preserves the fields only the caller knows', async () => {
    // The partial is not a throwaway: it carries the username the user typed and the
    // diagnostics flag, neither of which /api/auth/me echoes back in every build.
    mockMe({ student_id: 's1', name: 'Root', role: 'admin' })
    const stored = await establishSession('tok', { ...PARTIAL, username: 'typed-by-hand', diagnostic_completed: false })
    expect(stored.username).toBe('typed-by-hand')
    expect(stored.diagnostic_completed).toBe(false)
  })

  describe('when the profile fetch fails', () => {
    // The login itself succeeded. Discarding a working token because an enrichment request
    // failed would log the user out for something they did not do, so the session survives and
    // simply carries no role.
    it('keeps the session and stores no role', async () => {
      mockMe({}, false, 500)
      await establishSession('tok', { ...PARTIAL })
      expect(getToken()).toBe('tok')
      expect(getUserInfo()?.student_id).toBe('s1')
      expect(getUserInfo()?.role).toBeUndefined()
    })

    it('survives a network error', async () => {
      vi.stubGlobal('fetch', vi.fn().mockRejectedValue(new Error('offline')))
      await establishSession('tok', { ...PARTIAL })
      expect(getToken()).toBe('tok')
      expect(getUserInfo()?.student_id).toBe('s1')
    })

    it('never invents a role the server did not send', async () => {
      // A body with no `role`, from an older build or a proxy that trims unknown fields.
      mockMe({ student_id: 's1', name: 'Root' })
      await establishSession('tok', { ...PARTIAL })
      expect(getUserInfo()?.role).toBeUndefined()
    })

    // The safe direction, stated as an assertion: a missing role reads as a learner, so a
    // learner is never offered administration by a failed request.
    it('treats a missing role as a learner', async () => {
      mockMe({ student_id: 's1', name: 'Root' })
      await establishSession('tok', { ...PARTIAL })
      expect(getUserInfo()?.role === 'admin').toBe(false)
    })
  })

  it('notifies listeners so the app shell reacts', async () => {
    mockMe({ student_id: 's1', name: 'Root', role: 'admin' })
    const seen = vi.fn()
    window.addEventListener('auth-changed', seen)
    await establishSession('tok', { ...PARTIAL })
    // Once for the token, once for the user — the hook subscribes to this event, and a session
    // that set storage without announcing it would leave the header showing a signed-out state.
    expect(seen).toHaveBeenCalled()
    window.removeEventListener('auth-changed', seen)
  })

  it('leaves no stale role from a previous session', async () => {
    setToken('old')
    localStorage.setItem('mathua_user', JSON.stringify({ ...PARTIAL, role: 'admin' }))
    mockMe({ student_id: 's2', name: 'Ada' })
    await establishSession('new', { ...PARTIAL })
    // The whole user object is replaced, so the previous account's admin role cannot survive
    // into the new session on a shared browser.
    expect(getUserInfo()?.role).toBeUndefined()
    expect(getToken()).toBe('new')
  })

  it('does not require a token to already be set', async () => {
    // setToken must happen first inside the helper: getAuthHeaders() would otherwise read
    // nothing and the fetch would go out unauthenticated.
    clearToken()
    const fetchMock = mockMe({ student_id: 's1', name: 'Root', role: 'admin' })
    await establishSession('tok', { ...PARTIAL })
    const [, init] = fetchMock.mock.calls[0]
    expect(new Headers(init?.headers).get('Authorization')).toBe('Bearer tok')
  })
})