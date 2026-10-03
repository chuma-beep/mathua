import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import DangerZone from '../components/DangerZone'
import DeleteAccount from '../components/DeleteAccount'
import { resetAccount, deleteAccount, getMe } from '../lib/api'
import { getUserInfo } from '../lib/auth'

const pushMock = vi.fn()

vi.mock('next/navigation', () => ({
  useRouter: () => ({ push: pushMock }),
}))

vi.mock('../lib/api', async (importOriginal) => {
  const mod = await importOriginal<typeof import('../lib/api')>()
  return { ...mod, getMe: vi.fn(), getAttempts: vi.fn(), resetAccount: vi.fn(), deleteAccount: vi.fn() }
})

vi.mock('../lib/auth', async (importOriginal) => {
  const mod = await importOriginal<typeof import('../lib/auth')>()
  return { ...mod, getUserInfo: vi.fn(), signOut: vi.fn(), clearGuest: vi.fn() }
})

beforeEach(() => {
  pushMock.mockClear()
  vi.mocked(getUserInfo).mockReturnValue(null)
  vi.mocked(getMe).mockReset()
  vi.mocked(resetAccount).mockReset()
  vi.mocked(deleteAccount).mockReset()
  vi.mocked(getMe).mockResolvedValue({ student_id: 's1', name: 'A', has_password: false })
  vi.mocked(resetAccount).mockResolvedValue(undefined)
  vi.mocked(deleteAccount).mockResolvedValue(undefined)
})

// The learner-visible contract: what does the danger zone show when the server
// refuses? Anything other than a readable sentence is a bug, because a raw
// JSON envelope in the UI reads as a broken app.
describe('DangerZone reset: the learner-visible outcome', () => {
  it('sends the exact phrase and lands on /profile', async () => {
    render(<DangerZone />)
    fireEvent.change(screen.getByLabelText(/Type.*to confirm/i), { target: { value: 'reset my progress' } })
    fireEvent.click(screen.getByRole('button', { name: /Reset everything above/i }))
    await waitFor(() => expect(pushMock).toHaveBeenCalledWith('/profile'))
    expect(vi.mocked(resetAccount).mock.calls[0][0]).toBe('reset my progress')
  })

  it('shows a readable sentence, not a raw JSON envelope, when the server refuses', async () => {
    vi.mocked(resetAccount).mockRejectedValue(new Error('confirmation phrase does not match'))
    render(<DangerZone />)
    fireEvent.change(screen.getByLabelText(/Type.*to confirm/i), { target: { value: 'reset my progress' } })
    fireEvent.click(screen.getByRole('button', { name: /Reset everything above/i }))
    await waitFor(() => expect(screen.getByText(/confirmation phrase does not match/i)).toBeTruthy())
    expect(pushMock).not.toHaveBeenCalled()
  })

  it('re-arms after a failure so a retry is possible', async () => {
    vi.mocked(resetAccount).mockRejectedValueOnce(new Error('confirmation phrase does not match'))
    render(<DangerZone />)
    fireEvent.change(screen.getByLabelText(/Type.*to confirm/i), { target: { value: 'reset my progress' } })
    fireEvent.click(screen.getByRole('button', { name: /Reset everything above/i }))
    await waitFor(() => expect(vi.mocked(resetAccount)).toHaveBeenCalledTimes(1))
    // The button must not stay disabled, or the learner is stranded.
    expect(screen.getByRole('button', { name: /Reset everything above/i })).not.toBeDisabled()
  })
})

describe('DeleteAccount: the learner-visible outcome', () => {
  it('shows a readable sentence when the server refuses', async () => {
    vi.mocked(deleteAccount).mockRejectedValue(new Error('current password is incorrect'))
    render(<DeleteAccount />)
    fireEvent.change(screen.getByLabelText(/Type.*to confirm/i), { target: { value: 'delete my account' } })
    fireEvent.click(screen.getByRole('button', { name: /Delete my account/i }))
    await waitFor(() => expect(screen.getByText(/current password is incorrect/i)).toBeTruthy())
    expect(pushMock).not.toHaveBeenCalled()
  })
})
// The transport must read the {error} envelope. Throwing res.text() surfaced a
// literal `{"error":"…"}` in the danger zone, which reads as a broken app
// rather than an explanation — and it hid the distinction between a wrong
// phrase (400) and a real server fault (500).
describe('danger-zone transport reports a readable sentence', () => {
  const realFetch = globalThis.fetch

  afterEach(() => { globalThis.fetch = realFetch })

  async function callReset(body: string, status: number) {
    const stub: typeof fetch = async () => new Response(body, { status })
    globalThis.fetch = vi.fn(stub)
    // importActual: this file mocks ../lib/api for the component tests, which
    // would stub out the very function under test here.
    const { resetAccount } = await vi.importActual<typeof import('../lib/api')>('../lib/api')
    return resetAccount('reset my progress').then(() => null, (e: Error) => e.message)
  }

  it('uses the server error text, not the raw JSON body', async () => {
    const msg = await callReset(JSON.stringify({ error: 'reset failed and nothing was changed — please try again' }), 500)
    expect(msg).toBe('Reset failed: reset failed and nothing was changed — please try again')
    expect(msg).not.toContain('{"error"')
  })

  it('still names the operation when the body is not JSON', async () => {
    const msg = await callReset('<html>502</html>', 502)
    expect(msg).toBe('Reset failed')
  })
})

// A learner whose getMe() call fails must not be stranded: the server answers
// 401 "current password is required" and there has to be a field to type it in.
describe('DeleteAccount when the password requirement is unknown', () => {
  it('shows the password field and does not require it', async () => {
    vi.mocked(getMe).mockRejectedValue(new Error('offline'))
    vi.mocked(getUserInfo).mockReturnValue(null)
    render(<DeleteAccount />)
    await waitFor(() => expect(screen.getByLabelText('Current password, to confirm deletion')).toBeTruthy())
    fireEvent.change(screen.getByLabelText(/Type.*to confirm/i), { target: { value: 'delete my account' } })
    // Not blocked: the server is the gate, and a 401 can now be retried.
    expect(screen.getByRole('button', { name: /Delete my account/i })).not.toBeDisabled()
  })

  it('sends a typed password even while the requirement is unknown', async () => {
    vi.mocked(getMe).mockRejectedValue(new Error('offline'))
    vi.mocked(getUserInfo).mockReturnValue(null)
    render(<DeleteAccount />)
    await waitFor(() => expect(screen.getByLabelText('Current password, to confirm deletion')).toBeTruthy())
    fireEvent.change(screen.getByLabelText(/Type.*to confirm/i), { target: { value: 'delete my account' } })
    fireEvent.change(screen.getByLabelText('Current password, to confirm deletion'), { target: { value: 'Engine!n1' } })
    fireEvent.click(screen.getByRole('button', { name: /Delete my account/i }))
    await waitFor(() => expect(vi.mocked(deleteAccount)).toHaveBeenCalledWith({ phrase: 'delete my account', password: 'Engine!n1' }))
  })
})
