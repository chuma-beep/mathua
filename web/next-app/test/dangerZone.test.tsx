import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import DangerZone, { attemptsToCSV } from '../components/DangerZone'
import DeleteAccount from '../components/DeleteAccount'
import { getMe, deleteAccount } from '../lib/api'
import { getUserInfo, signOut, clearGuest } from '../lib/auth'

const pushMock = vi.fn()

vi.mock('next/navigation', () => ({
  useRouter: () => ({ push: pushMock }),
}))

vi.mock('../lib/api', async (importOriginal) => {
  const mod = await importOriginal<typeof import('../lib/api')>()
  return { ...mod, getMe: vi.fn(), deleteAccount: vi.fn(), getAttempts: vi.fn() }
})

vi.mock('../lib/auth', async (importOriginal) => {
  const mod = await importOriginal<typeof import('../lib/auth')>()
  return { ...mod, getUserInfo: vi.fn(), signOut: vi.fn(), clearGuest: vi.fn() }
})

beforeEach(() => {
  pushMock.mockClear()
  vi.mocked(getUserInfo).mockReturnValue(null)
  vi.mocked(getMe).mockResolvedValue({ student_id: 's1', name: 'A', has_password: false })
  vi.mocked(deleteAccount).mockResolvedValue(undefined)
})

describe('DangerZone', () => {
  it('keeps reset disabled until the exact phrase is typed', () => {
    render(<DangerZone />)
    const button = screen.getByRole('button', { name: /Reset everything above/ })
    expect(button).toBeDisabled()
    fireEvent.change(screen.getByLabelText(/Type.*to confirm/i), { target: { value: 'reset' } })
    expect(button).toBeDisabled()
    fireEvent.change(screen.getByLabelText(/Type.*to confirm/i), { target: { value: 'reset my progress' } })
    expect(button).not.toBeDisabled()
  })

  it('arms on the phrase as a phone would type it', () => {
    // The reported symptom: on iOS the button never enables.
    //
    // Both inputs are `type="text"` with no `autoCapitalize`, so Safari capitalises the
    // first letter as you type — "reset my progress" becomes "Reset my progress". Both
    // the client comparison and the server's were case-sensitive, so the button stayed
    // disabled forever and the server would have rejected the capitalised phrase with
    // "confirmation phrase does not match". Chromium does not auto-capitalise, which is
    // why the whole suite was green.
    render(<DangerZone />)
    const input = screen.getByLabelText(/reset my progress/i) as HTMLInputElement
    const submit = screen.getByRole('button', { name: /reset everything/i })

    // The OS must not be allowed to mangle it in the first place.
    expect(input.getAttribute('autocapitalize')).toBe('none')
    expect(input.getAttribute('spellcheck')).toBe('false')

    // And if it does anyway — a keyboard that capitalises, a paste from a note — the
    // button must still arm. The phrase confirms intent; it is not a secret.
    fireEvent.change(input, { target: { value: 'Reset My Progress' } })
    expect(submit).not.toBeDisabled()
    // Leading/trailing whitespace is still not a way to bypass the gate.
    fireEvent.change(input, { target: { value: 'reset my progressX' } })
    expect(submit).toBeDisabled()
  })

  it('says why the button is disabled instead of leaving a grey rectangle', () => {
    render(<DangerZone />)
    const input = screen.getByLabelText(/reset my progress/i)
    const submit = screen.getByRole('button', { name: /reset everything/i })

    expect(submit).toBeDisabled()
    // Nothing typed yet: nothing to say.
    expect(screen.queryByRole('status')).toBeNull()

    fireEvent.change(input, { target: { value: 'reset my' } })
    expect(submit).toBeDisabled()
    // Something typed and it is wrong: say so, instead of a control that just greys out.
    expect(screen.getByRole('status').textContent).toMatch(/does not match/i)
  })

  it('states the keeps list including destination, deadline and pace', () => {
    render(<DangerZone />)
    expect(screen.getByText(/Destination, deadline and pace/i)).toBeTruthy()
    expect(screen.getByText(/profile and leaderboard/i)).toBeTruthy()
  })
})

describe('attemptsToCSV', () => {
  it('quotes and headers rows', () => {
    const csv = attemptsToCSV([
      { session_id: 's', student_id: 'u', concept_id: 'a', answer: '4', expected: '4', correct: true, elapsed_seconds: 3, timestamp: '2026-01-01', question: '2+2?', source: 'practice', explanation: '' },
      { session_id: 's', student_id: 'u', concept_id: 'b', answer: 'x "y"', expected: 'z', correct: false, elapsed_seconds: 9, timestamp: '2026-01-02', question: 'q', source: 'quiz', explanation: '' },
    ])
    const lines = csv.split('\n')
    expect(lines[0]).toBe('timestamp,concept_id,question,answer,expected,correct,elapsed_seconds,source')
    expect(lines).toHaveLength(3)
    expect(lines[2]).toContain('"x ""y"""')
  })
})

describe('DeleteAccount', () => {
  it('keeps delete disabled until the exact phrase is typed', () => {
    render(<DeleteAccount />)
    const button = screen.getByRole('button', { name: /Delete my account/ })
    expect(button).toBeDisabled()
    fireEvent.change(screen.getByLabelText(/Type.*to confirm/i), { target: { value: 'delete' } })
    expect(button).toBeDisabled()
    fireEvent.change(screen.getByLabelText(/Type.*to confirm/i), { target: { value: 'delete my account' } })
    expect(button).not.toBeDisabled()
  })

  it('requires the current password only for password accounts', async () => {
    vi.mocked(getMe).mockResolvedValue({ student_id: 's1', name: 'A', has_password: true })
    render(<DeleteAccount />)
    expect(await screen.findByLabelText('Current password, to confirm deletion')).toBeTruthy()
    fireEvent.change(screen.getByLabelText(/Type.*to confirm/i), { target: { value: 'delete my account' } })
    // Phrase alone is not enough while the password is empty.
    expect(screen.getByRole('button', { name: /^Delete my account$/ })).toBeDisabled()
    fireEvent.change(screen.getByLabelText('Current password, to confirm deletion'), { target: { value: 'secret' } })
    expect(screen.getByRole('button', { name: /^Delete my account$/ })).not.toBeDisabled()
  })

  it('omits the password field for OAuth-only rows', async () => {
    render(<DeleteAccount />)
    await screen.findByText(/Delete account/)
    expect(screen.queryByLabelText('Current password, to confirm deletion')).toBeNull()
  })

  it('deletes, signs out fully, and lands on /', async () => {
    vi.mocked(getMe).mockResolvedValue({ student_id: 's1', name: 'A', has_password: true })
    render(<DeleteAccount />)
    await screen.findByLabelText('Current password, to confirm deletion')
    fireEvent.change(screen.getByLabelText(/Type.*to confirm/i), { target: { value: 'delete my account' } })
    fireEvent.change(screen.getByLabelText('Current password, to confirm deletion'), { target: { value: 'secret' } })
    fireEvent.click(screen.getByRole('button', { name: /^Delete my account$/ }))
    await screen.findByText('Deleting…')
    expect(vi.mocked(deleteAccount)).toHaveBeenCalledWith({ phrase: 'delete my account', password: 'secret' })
    expect(vi.mocked(signOut)).toHaveBeenCalled()
    expect(vi.mocked(clearGuest)).toHaveBeenCalled()
    expect(pushMock).toHaveBeenCalledWith('/')
  })

  it('surfaces server errors without navigating', async () => {
    vi.mocked(deleteAccount).mockRejectedValue(new Error('current password is incorrect'))
    render(<DeleteAccount />)
    fireEvent.change(screen.getByLabelText(/Type.*to confirm/i), { target: { value: 'delete my account' } })
    fireEvent.click(screen.getByRole('button', { name: /^Delete my account$/ }))
    expect(await screen.findByText('current password is incorrect')).toBeTruthy()
    expect(pushMock).not.toHaveBeenCalled()
  })

  it('states total scope: nothing is kept', () => {
    render(<DeleteAccount />)
    expect(screen.getByText(/Nothing — deletion is total/i)).toBeTruthy()
  })
})
