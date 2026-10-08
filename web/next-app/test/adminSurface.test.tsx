import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent, waitFor, within } from '@testing-library/react'
import AdminUsersPage from '../app/admin/users/page'
import AdminAuditPage from '../app/admin/audit/page'
import AdminOverviewPage from '../app/admin/page'
import { searchAdminUsers, setUserRole, adminOverview, listAdminAudit } from '../lib/api'
import type { AdminUser } from '../lib/api'
import type { UserInfo } from '../lib/auth'

vi.mock('../lib/api', async (importOriginal) => ({
  ...(await importOriginal<typeof import('../lib/api')>()),
  adminOverview: vi.fn(),
  searchAdminUsers: vi.fn(),
  setUserRole: vi.fn(),
  listAdminAudit: vi.fn(),
}))

// A mutable module-level value rather than a vi.fn the tests re-point: this suite changes the
// role between cases, and a factory returning a plain closure reads the variable at call time.
const authState = { current: { role: 'admin' as 'student' | 'admin' | undefined, loggedIn: true } }

// Built in statements rather than with a conditional spread, because the "missing role" case is
// exactly the case a spread hides: `role: undefined` and an absent key both read as falsy to
// `user?.role === 'admin'`, so a spread would make that test pass without testing anything.
type FakeUser = Omit<UserInfo, 'role'> & { role?: 'student' | 'admin' }

function fakeUser(): FakeUser {
  const user: FakeUser = {
    student_id: 'x', name: 'X', username: 'x', concepts_mastered: 0,
    current_streak: 0, level: 'Novice', diagnostic_completed: true,
  }
  if (authState.current.role !== undefined) user.role = authState.current.role
  return user
}

vi.mock('../hooks/useAuthState', () => ({
  useAuthState: () => ({ loggedIn: authState.current.loggedIn, user: fakeUser() }),
}))
vi.mock('../components/Header', () => ({ default: () => <header>banner</header> }))
vi.mock('../components/Footer', () => ({ default: () => <footer /> }))
vi.mock('../components/BottomTabs', () => ({ default: () => <nav /> }))

const ADMIN: AdminUser = {
  id: 'a1', name: 'Root', username: 'root', email: 'root@example.com',
  role: 'admin', created_at: '2026-01-02T03:04:05Z', has_password: true,
}
const LEARNER: AdminUser = {
  id: 's1', name: 'Ada', username: 'ada', email: 'ada@example.com',
  role: 'student', created_at: '2026-02-03T04:05:06Z', has_password: false,
}

function asRole(role?: 'student' | 'admin') {
  authState.current.role = role
}

beforeEach(() => {
  vi.clearAllMocks()
  asRole('admin')
  vi.mocked(searchAdminUsers).mockResolvedValue({ users: [ADMIN, LEARNER] })
})

// ── The UI is not the boundary ─────────────────────────────────────────────

// These three assert the same product fact from three directions: what the client renders, what
// the buttons send, and what the client types. None of it is security — the server decides — and
// the whole point of the phase is that it does. So the tests here are about not pretending
// otherwise, and about the client degrading honestly when the server says no.
describe('the client offers administration without enforcing it', () => {
  it('offers the admin nav to an administrator', async () => {
    render(<AdminUsersPage />)
    await screen.findByRole('navigation', { name: 'Admin sections' })
    for (const label of ['Overview', 'Users', 'Audit']) {
      expect(screen.getByRole('link', { name: label })).toBeTruthy()
    }
  })

  // No role, no nav. The pages still mount, because the server decides — and a learner who
  // navigated here is told why rather than shown a page that looks broken.
  it('omits the nav for a learner and says the area needs the role', async () => {
    asRole('student')
    render(<AdminUsersPage />)
    await waitFor(() => {
      expect(screen.queryByRole('navigation', { name: 'Admin sections' })).toBeNull()
    })
    expect(screen.getByText(/needs an administrator role/)).toBeTruthy()
  })

  // A payload with no `role` at all — an old bundle, a stale cached user object. Treating that
  // as an administrator would be the worst reading; treating it as a learner is the one that
  // fails closed.
  it('treats a missing role as a learner', async () => {
    asRole(undefined)
    render(<AdminUsersPage />)
    await waitFor(() => {
      expect(screen.queryByRole('navigation', { name: 'Admin sections' })).toBeNull()
    })
  })

  // The button says Promote/Demote from what it fetched, and sends exactly one field. A role
  // invented in the client would be rejected server-side; sending a batch instead of a role
  // would be a mass-assignment route through the UI.
  it('sends only the role it is changing to', async () => {
    vi.mocked(setUserRole).mockResolvedValue({ user: { ...LEARNER, role: 'admin' }, changed: true })
    render(<AdminUsersPage />)
    await screen.findByText('ada@example.com')

    // Scoped to the table: the dialog's confirm button carries the same label on purpose, so
    // that the confirmation reads as the action it confirms.
    fireEvent.click(within(screen.getByRole('table')).getAllByRole('button', { name: 'Make administrator' })[0])
    fireEvent.click(within(await screen.findByRole('dialog')).getByRole('button', { name: 'Make administrator' }))

    await waitFor(() => expect(setUserRole).toHaveBeenCalledWith('s1', 'admin'))
    expect(vi.mocked(setUserRole).mock.calls[0]).toHaveLength(2)
  })

  it('offers Remove administrator for an administrator row', async () => {
    vi.mocked(setUserRole).mockResolvedValue({ user: { ...ADMIN, role: 'student' }, changed: true })
    render(<AdminUsersPage />)
    await screen.findByText('root@example.com')

    fireEvent.click(screen.getAllByRole('button', { name: 'Remove administrator' })[0])
    fireEvent.click(await screen.findByRole('button', { name: 'Remove access' }))
    await waitFor(() => expect(setUserRole).toHaveBeenCalledWith('a1', 'student'))
  })

  // The last-admin refusal arrives as an error and is shown. The page does not pre-empt it by
  // disabling the button, because the client cannot know how many administrators exist: that is
  // the server's count, and guessing it would produce a control that lies about being safe.
  it('surfaces a server refusal instead of hiding the control', async () => {
    vi.mocked(setUserRole).mockRejectedValue(new Error('cannot demote the last administrator'))
    render(<AdminUsersPage />)
    await screen.findByText('root@example.com')

    fireEvent.click(within(screen.getByRole('table')).getAllByRole('button', { name: 'Remove administrator' })[0])
    fireEvent.click(within(await screen.findByRole('dialog')).getByRole('button', { name: 'Remove access' }))
    await waitFor(() => {
      expect(screen.getByRole('alert').textContent).toMatch(/last administrator/)
    })
    // The control is still there, because it is not the client's decision to make.
    expect(screen.getAllByRole('button', { name: 'Remove administrator' }).length).toBeGreaterThan(0)
  })

  // After a change the list is re-read from the server rather than patched locally. An optimistic
  // patch would show the client's guess, which is exactly what failed when the server refused.
  it('re-reads the list after a change rather than guessing', async () => {
    vi.mocked(setUserRole).mockResolvedValue({ user: { ...LEARNER, role: 'admin' }, changed: true })
    render(<AdminUsersPage />)
    await screen.findByText('ada@example.com')
    const before = vi.mocked(searchAdminUsers).mock.calls.length

    fireEvent.click(within(screen.getByRole('table')).getAllByRole('button', { name: 'Make administrator' })[0])
    fireEvent.click(within(await screen.findByRole('dialog')).getByRole('button', { name: 'Make administrator' }))
    await waitFor(() => {
      expect(vi.mocked(searchAdminUsers).mock.calls.length).toBeGreaterThan(before)
    })
  })

  // The projection has no field for learning state, and the page must not offer a way to edit it.
  // An administrator with a box in front of mastery will eventually use it.
  it('offers no editor for learning state', async () => {
    render(<AdminUsersPage />)
    await screen.findByText('ada@example.com')
    // Asserted on the controls rather than on the text, because the page deliberately *names*
    // what it will not edit — a substring scan would be defeated by its own explanation, and
    // would equally be defeated by deleting that explanation.
    //
    // One text input exists on this page: the search box. Every button is a role change. So a
    // second input appearing anywhere in this surface is the thing to catch.
    const inputs = [...document.querySelectorAll('input, textarea, select')]
    expect(inputs.map((i) => i.getAttribute('type') ?? i.tagName)).toEqual(['search'])
    for (const field of ['xp_total', 'mastery', 'attempts', 'concept_progress', 'streak', 'status']) {
      expect(document.body.querySelector(`[name="${field}"]`)).toBeNull()
    }
    expect(screen.getByText(/not editable here/)).toBeTruthy()
  })

  it('shows whether an account has a password, without ever showing one', async () => {
    render(<AdminUsersPage />)
    await screen.findByText('ada@example.com')
    // Ada is a provider sign-in and is labelled as such, which is a fact about the account.
    expect(screen.getByText(/provider, not a password/)).toBeTruthy()
    // Nothing resembling a hash anywhere on the page.
    expect(document.body.innerHTML).not.toMatch(/\$argon2|\$2[aby]\$/)
  })

  // A page that says "Loading" after it has already failed is lying about its own state, and
  // here it would lie forever: the failed load never sets `users`, so the null never clears.
  // A learner who navigates to /admin/users is exactly who hits this, because the API answers
  // 403 and the error path runs.
  it('does not claim to be loading after the request failed', async () => {
    vi.mocked(searchAdminUsers).mockRejectedValue(new Error('admin role required'))
    render(<AdminUsersPage />)
    await waitFor(() => expect(screen.getByRole('alert')).toBeTruthy())
    expect(screen.queryByText(/Loading accounts/)).toBeNull()
    expect(screen.queryByRole('table')).toBeNull()
  })

  it('does not claim to be loading after the overview request failed', async () => {
    vi.mocked(adminOverview).mockRejectedValue(new Error('admin role required'))
    render(<AdminOverviewPage />)
    await waitFor(() => expect(screen.getByRole('alert')).toBeTruthy())
    expect(screen.queryByText(/LOADING OVERVIEW/i)).toBeNull()
  })

  it('does not claim to be loading after the audit request failed', async () => {
    vi.mocked(listAdminAudit).mockRejectedValue(new Error('admin role required'))
    render(<AdminAuditPage />)
    await waitFor(() => expect(screen.getByRole('alert')).toBeTruthy())
    expect(screen.queryByText(/Loading the log/)).toBeNull()
  })

  it('searches by email, id, username or name', async () => {
    render(<AdminUsersPage />)
    await screen.findByText('ada@example.com')
    const input = screen.getByPlaceholderText(/Search by email/)
    fireEvent.change(input, { target: { value: 'ada@' } })
    fireEvent.click(screen.getByRole('button', { name: 'Search' }))
    await waitFor(() => {
      expect(vi.mocked(searchAdminUsers).mock.calls.at(-1)?.[0]).toBe('ada@')
    })
  })
})

// ── The audit page ─────────────────────────────────────────────────────────

describe('the audit page is read-only', () => {
  beforeEach(() => {
    vi.mocked(listAdminAudit).mockResolvedValue({
      events: [{
        id: 7,
        actor_id: 'a1',
        action: 'ADMIN_PROMOTE_USER',
        entity_type: 'user',
        entity_id: 's1',
        before_json: '{"role":"student"}',
        after_json: '{"role":"admin"}',
        created_at: '2026-03-04T05:06:07Z',
      }],
    })
  })

  it('lists timestamp, administrator, action and entity', async () => {
    render(<AdminAuditPage />)
    await screen.findByText('ADMIN_PROMOTE_USER')
    expect(screen.getByText(/by a1 → user s1/)).toBeTruthy()
    expect(screen.getAllByText('2026-03-04T05:06:07Z').length).toBeGreaterThan(0)
  })

  // No edit, no delete, no note field. An audit page with a correction control is an audit page
  // whose contents are curated, and the server would refuse the write anyway.
  it('offers no way to alter an event', async () => {
    render(<AdminAuditPage />)
    await screen.findByText('ADMIN_PROMOTE_USER')

    // No text entry anywhere on the page: an audit trail you can annotate is one you have
    // curated, and the server would refuse the write anyway.
    expect(document.querySelectorAll('input, textarea, select')).toHaveLength(0)
    for (const name of [/edit/i, /^delete/i, /remove/i, /correct/i, /annotate/i, /note/i]) {
      expect(screen.queryByRole('button', { name })).toBeNull()
      expect(screen.queryByRole('link', { name })).toBeNull()
    }

    // The only button on a row is the disclosure, and expanding it says so rather than offering
    // a correction.
    const row = screen.getByRole('button', { name: /ADMIN_PROMOTE_USER/ })
    expect(row.tagName).toBe('BUTTON')
    fireEvent.click(row)
    await waitFor(() => expect(screen.getByText(/cannot be edited or deleted/)).toBeTruthy())
  })

  it('reveals before and after on demand', async () => {
    render(<AdminAuditPage />)
    await screen.findByText('ADMIN_PROMOTE_USER')
    expect(screen.queryByText(/student/)).toBeNull()

    fireEvent.click(screen.getByRole('button', { name: /ADMIN_PROMOTE_USER/ }))
    await waitFor(() => {
      expect(screen.getByText('Before')).toBeTruthy()
      expect(screen.getByText('After')).toBeTruthy()
    })
    expect(document.body.textContent).toMatch(/student/)
    expect(document.body.textContent).toMatch(/admin/)
  })

  // An unparseable record is shown as-is rather than swallowed. Rendering nothing would make a
  // malformed row indistinguishable from an empty one, which is the reading that hides trouble.
  it('shows an unparseable record rather than nothing', async () => {
    vi.mocked(listAdminAudit).mockResolvedValue({
      events: [{
        id: 8, actor_id: 'a1', action: 'ADMIN_PROMOTE_USER', entity_type: 'user', entity_id: 's1',
        before_json: 'not json at all', after_json: '', created_at: '2026-03-04T05:06:07Z',
      }],
    })
    render(<AdminAuditPage />)
    await screen.findByText('ADMIN_PROMOTE_USER')
    fireEvent.click(screen.getByRole('button', { name: /ADMIN_PROMOTE_USER/ }))
    await waitFor(() => expect(screen.getByText('not json at all')).toBeTruthy())
  })
})

// ── The overview ───────────────────────────────────────────────────────────

describe('the overview is a landing page', () => {
  beforeEach(() => {
    vi.mocked(adminOverview).mockResolvedValue({
      users: 4, admins: 1, concepts: 657, domains: 15, questions: 120,
      recent_activity: [],
    })
  })

  it('shows only counts the backend already has', async () => {
    render(<AdminOverviewPage />)
    await waitFor(() => expect(screen.getByText('Accounts')).toBeTruthy())
    expect(screen.getByText('657')).toBeTruthy()
    expect(screen.getByText('15')).toBeTruthy()
    expect(screen.getByText('120')).toBeTruthy()
    // No analytics vocabulary crept in.
    for (const word of ['retention', 'cohort', 'revenue', 'churn', 'engagement']) {
      expect(document.body.textContent?.toLowerCase()).not.toContain(word)
    }
  })

  // A deployment with no administrators has no administrative surface at all, including the one
  // that would fix it. Saying so on the dashboard, with the operator command, is the difference
  // between a locked-out operator and a silent outage.
  it('says plainly when there are no administrators', async () => {
    vi.mocked(adminOverview).mockResolvedValue({
      users: 4, admins: 0, concepts: 657, domains: 15, questions: 120, recent_activity: [],
    })
    render(<AdminOverviewPage />)
    await waitFor(() => expect(screen.getByRole('alert').textContent).toMatch(/No administrators exist/))
    expect(screen.getByRole('alert').textContent).toMatch(/mathua admin promote/)
  })

  it('does not raise that alarm when one administrator exists', async () => {
    render(<AdminOverviewPage />)
    await waitFor(() => expect(screen.getByText('Administrators')).toBeTruthy())
    expect(screen.queryByText(/No administrators exist/)).toBeNull()
  })
})
// ── Confirmation before a privileged change ────────────────────────────────

// Every assertion in this block is about the *absence* of an effect until a second deliberate
// act. That is the part that regresses quietly: add a confirmation and every existing test still
// passes, because the flow still works. Only a test that stops after the first click notices.
describe('promotion and removal are confirmed before they happen', () => {
  const rowButton = (name: string) =>
    within(screen.getByRole('table')).getAllByRole('button', { name })[0]

  it('does not call the API when the row action is clicked', async () => {
    vi.mocked(setUserRole).mockResolvedValue({ user: { ...LEARNER, role: 'admin' }, changed: true })
    render(<AdminUsersPage />)
    await screen.findByText('ada@example.com')

    fireEvent.click(rowButton('Make administrator'))

    // A privileged, audited, immediately-effective change must not be one click from a list of
    // rows. Nothing has been sent yet.
    expect(setUserRole).not.toHaveBeenCalled()
    expect(await screen.findByRole('dialog')).toBeTruthy()
  })

  it('names the account and the consequence in the promotion dialog', async () => {
    render(<AdminUsersPage />)
    await screen.findByText('ada@example.com')

    fireEvent.click(rowButton('Make administrator'))
    const dialog = await screen.findByRole('dialog')
    expect(within(dialog).getByText('Make administrator?')).toBeTruthy()
    expect(dialog.textContent).toMatch(/Ada/)
    expect(dialog.textContent).toMatch(/administration interface/)
    // And that it is recorded, because the audit row is permanent.
    expect(dialog.textContent).toMatch(/audit log/i)
  })

  it('names the account and the consequence in the removal dialog', async () => {
    render(<AdminUsersPage />)
    await screen.findByText('root@example.com')

    fireEvent.click(rowButton('Remove administrator'))
    const dialog = await screen.findByRole('dialog')
    expect(within(dialog).getByText('Remove administrator access?')).toBeTruthy()
    expect(dialog.textContent).toMatch(/Root/)
    // "Immediately", because the role is re-read per request — there is no session to wait out.
    expect(dialog.textContent).toMatch(/immediately lose administrative access/)
  })

  it('cancels without calling the API', async () => {
    vi.mocked(setUserRole).mockResolvedValue({ user: LEARNER, changed: true })
    render(<AdminUsersPage />)
    await screen.findByText('ada@example.com')

    fireEvent.click(rowButton('Make administrator'))
    fireEvent.click(within(await screen.findByRole('dialog')).getByRole('button', { name: 'Cancel' }))

    await waitFor(() => expect(screen.queryByRole('dialog')).toBeNull())
    expect(setUserRole).not.toHaveBeenCalled()
  })

  it('cancels on Escape', async () => {
    vi.mocked(setUserRole).mockResolvedValue({ user: LEARNER, changed: true })
    render(<AdminUsersPage />)
    await screen.findByText('ada@example.com')

    fireEvent.click(rowButton('Make administrator'))
    await screen.findByRole('dialog')
    fireEvent.keyDown(document, { key: 'Escape' })

    await waitFor(() => expect(screen.queryByRole('dialog')).toBeNull())
    expect(setUserRole).not.toHaveBeenCalled()
  })

  // Focus starts on Cancel. Auto-focusing the confirming button would let a stray Enter key —
  // with the dialog already open from a stray click — grant administrator access, which is the
  // exact sequence the confirmation exists to break.
  it('focuses Cancel, not the confirming button', async () => {
    render(<AdminUsersPage />)
    await screen.findByText('ada@example.com')

    fireEvent.click(rowButton('Make administrator'))
    const dialog = await screen.findByRole('dialog')
    await waitFor(() => expect(document.activeElement).toBe(within(dialog).getByRole('button', { name: 'Cancel' })))
  })

  it('sends the change once confirmed', async () => {
    vi.mocked(setUserRole).mockResolvedValue({ user: { ...LEARNER, role: 'admin' }, changed: true })
    render(<AdminUsersPage />)
    await screen.findByText('ada@example.com')

    fireEvent.click(rowButton('Make administrator'))
    fireEvent.click(within(await screen.findByRole('dialog')).getByRole('button', { name: 'Make administrator' }))

    await waitFor(() => expect(setUserRole).toHaveBeenCalledTimes(1))
    await waitFor(() => expect(screen.queryByRole('dialog')).toBeNull())
    expect(await screen.findByRole('status')).toHaveTextContent(/Ada is now an administrator/)
  })

  it('reports when the account already had the role', async () => {
    vi.mocked(setUserRole).mockResolvedValue({ user: ADMIN, changed: false })
    render(<AdminUsersPage />)
    await screen.findByText('root@example.com')

    fireEvent.click(rowButton('Remove administrator'))
    fireEvent.click(within(await screen.findByRole('dialog')).getByRole('button', { name: 'Remove access' }))

    expect(await screen.findByRole('status')).toHaveTextContent(/already a learner/)
  })

  // A refusal must not leave the dialog open on top of the message, hiding it behind a modal.
  it('closes the dialog and shows the server error', async () => {
    vi.mocked(setUserRole).mockRejectedValue(new Error('cannot demote the last administrator'))
    render(<AdminUsersPage />)
    await screen.findByText('root@example.com')

    fireEvent.click(rowButton('Remove administrator'))
    fireEvent.click(within(await screen.findByRole('dialog')).getByRole('button', { name: 'Remove access' }))

    await waitFor(() => expect(screen.queryByRole('dialog')).toBeNull())
    expect(screen.getByRole('alert').textContent).toMatch(/last administrator/)
  })

  it('cannot be confirmed twice by mashing the button', async () => {
    // Held open rather than resolved, so the button is still disabled while it is clicked again.
    let release: (value: { user: AdminUser; changed: boolean }) => void = () => {}
    vi.mocked(setUserRole).mockImplementation(
      () => new Promise<{ user: AdminUser; changed: boolean }>((res) => { release = res }),
    )
    render(<AdminUsersPage />)
    await screen.findByText('ada@example.com')

    fireEvent.click(rowButton('Make administrator'))
    const dialog = await screen.findByRole('dialog')
    const confirm = within(dialog).getByRole('button', { name: 'Make administrator' })
    fireEvent.click(confirm)

    await waitFor(() => expect(vi.mocked(setUserRole).mock.calls.length).toBe(1))
    // The second click lands on a disabled button, so it cannot produce a second mutation.
    fireEvent.click(confirm)
    fireEvent.click(confirm)
    expect(vi.mocked(setUserRole).mock.calls.length).toBe(1)
    release({ user: { ...LEARNER, role: 'admin' }, changed: true })
  })

  it('shows the current role in the table', async () => {
    render(<AdminUsersPage />)
    await screen.findByText('ada@example.com')
    const row = within(screen.getByRole('table')).getByText('root@example.com').closest('tr')!
    expect(row.textContent).toMatch(/admin/)
    const learnerRow = within(screen.getByRole('table')).getByText('ada@example.com').closest('tr')!
    expect(learnerRow.textContent).toMatch(/student/)
  })

  // Nothing on this page may read as a credential. The projection has no such field and the
  // page must not have grown one.
  it('renders no credential-shaped value', async () => {
    render(<AdminUsersPage />)
    await screen.findByText('ada@example.com')
    expect(document.body.innerHTML).not.toMatch(/\$argon2|\$2[aby]\$|password_hash|share_token/)
  })
})
