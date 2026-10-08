import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent, waitFor, within } from '@testing-library/react'
import AdminUsersPage from '../app/admin/users/page'
import AdminAuditPage from '../app/admin/audit/page'
import AdminOverviewPage from '../app/admin/page'
import { searchAdminUsers, setUserRole, adminOverview, listAdminAudit } from '../lib/api'
import type { AdminUser, Role } from '../lib/api'
import type { UserInfo } from '../lib/auth'

vi.mock('../lib/api', async (importOriginal) => ({
  ...(await importOriginal<typeof import('../lib/api')>()),
  adminOverview: vi.fn(),
  searchAdminUsers: vi.fn(),
  setUserRole: vi.fn(),
  listAdminAudit: vi.fn(),
}))

// The caller's capability, set per test. The pages read it through useAdminMe; mocking that hook
// is cleaner than standing up /api/admin/me, and it makes "what a moderator sees" a one-liner.
const ALL_PERMISSIONS = [
  'admin.access', 'reports.read', 'reports.moderate', 'users.read',
  'admins.read', 'admins.manage', 'roles.grant_owner', 'audit.read',
  'content.read', 'content.manage',
]
const access = {
  permissions: [...ALL_PERMISSIONS] as string[],
  role: 'owner' as Role,
  assignable_roles: ['moderator', 'admin', 'owner', 'student'] as Role[],
}

function asPermissions(perms: string[]) {
  access.permissions = perms
}

vi.mock('../hooks/useAdminMe', () => ({
  useAdminMe: () => ({
    me: { id: 'root', role: access.role, permissions: access.permissions, assignable_roles: access.assignable_roles },
    loading: false,
    can: (p: string) => access.permissions.includes(p),
    isStaff: access.permissions.includes('admin.access'),
  }),
}))

const authState = { current: { role: 'owner' as 'student' | 'owner' | undefined, loggedIn: true } }

type FakeUser = Omit<UserInfo, 'role'> & { role?: 'student' | 'moderator' | 'admin' | 'owner' }

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

beforeEach(() => {
  vi.clearAllMocks()
  asPermissions(ALL_PERMISSIONS)
  vi.mocked(searchAdminUsers).mockResolvedValue({ users: [ADMIN, LEARNER] })
})

// The role picker for one row, then confirm. The select sends nothing until the dialog confirms.
function changeRole(email: string, role: Role) {
  const row = within(screen.getByRole('table')).getByText(email).closest('tr')!
  fireEvent.change(within(row).getByRole('combobox'), { target: { value: role } })
}

// ── The UI is not the boundary ─────────────────────────────────────────────

describe('the client offers administration without enforcing it', () => {
  it('offers the admin nav sections to a full staff member', async () => {
    render(<AdminUsersPage />)
    await screen.findByRole('navigation', { name: 'Admin sections' })
    for (const label of ['Overview', 'Reports', 'Users', 'Contributors', 'Content', 'Audit Log']) {
      expect(screen.getByRole('link', { name: label })).toBeTruthy()
    }
  })

  // A moderator works reports and content; the nav shows only those, because the sections are
  // filtered by the permissions the server reported.
  it('shows a moderator only the sections it can reach', async () => {
    asPermissions(['admin.access', 'reports.read', 'reports.moderate', 'content.read'])
    render(<AdminUsersPage />)
    await screen.findByRole('navigation', { name: 'Admin sections' })
    expect(screen.getByRole('link', { name: 'Reports' })).toBeTruthy()
    expect(screen.getByRole('link', { name: 'Content' })).toBeTruthy()
    expect(screen.queryByRole('link', { name: 'Contributors' })).toBeNull()
    expect(screen.queryByRole('link', { name: 'Audit Log' })).toBeNull()
  })

  it('omits the nav for an account with no admin access and says so', async () => {
    asPermissions([])
    render(<AdminUsersPage />)
    await waitFor(() => {
      expect(screen.queryByRole('navigation', { name: 'Admin sections' })).toBeNull()
    })
    expect(screen.getByText(/needs a staff role/)).toBeTruthy()
  })

  it('sends only the role it is changing to', async () => {
    vi.mocked(setUserRole).mockResolvedValue({ user: { ...LEARNER, role: 'admin' }, changed: true })
    render(<AdminUsersPage />)
    await screen.findByText('ada@example.com')

    changeRole('ada@example.com', 'admin')
    fireEvent.click(within(await screen.findByRole('dialog')).getByRole('button', { name: 'Make administrator' }))

    await waitFor(() => expect(setUserRole).toHaveBeenCalledWith('s1', 'admin'))
    expect(vi.mocked(setUserRole).mock.calls[0]).toHaveLength(2)
  })

  it('surfaces a server refusal instead of hiding the control', async () => {
    vi.mocked(setUserRole).mockRejectedValue(new Error('cannot remove the last owner'))
    render(<AdminUsersPage />)
    await screen.findByText('root@example.com')

    changeRole('root@example.com', 'student')
    fireEvent.click(within(await screen.findByRole('dialog')).getByRole('button', { name: 'Make learner' }))
    await waitFor(() => {
      expect(screen.getByRole('alert').textContent).toMatch(/last owner/)
    })
  })

  it('re-reads the list after a change rather than guessing', async () => {
    vi.mocked(setUserRole).mockResolvedValue({ user: { ...LEARNER, role: 'moderator' }, changed: true })
    render(<AdminUsersPage />)
    await screen.findByText('ada@example.com')
    const before = vi.mocked(searchAdminUsers).mock.calls.length

    changeRole('ada@example.com', 'moderator')
    fireEvent.click(within(await screen.findByRole('dialog')).getByRole('button', { name: 'Make moderator' }))
    await waitFor(() => {
      expect(vi.mocked(searchAdminUsers).mock.calls.length).toBeGreaterThan(before)
    })
  })

  // There is no editor for learning state. The only text entry is the search box; the only
  // selects are role pickers. A field named for XP or mastery must not exist.
  it('offers no editor for learning state', async () => {
    render(<AdminUsersPage />)
    await screen.findByText('ada@example.com')
    for (const field of ['xp_total', 'mastery', 'attempts', 'concept_progress', 'streak', 'status']) {
      expect(document.body.querySelector(`[name="${field}"]`)).toBeNull()
    }
    expect(screen.getByText(/not editable here/)).toBeTruthy()
  })

  it('shows whether an account has a password, without ever showing one', async () => {
    render(<AdminUsersPage />)
    await screen.findByText('ada@example.com')
    expect(screen.getByText(/provider, not a password/)).toBeTruthy()
    expect(document.body.innerHTML).not.toMatch(/\$argon2|\$2[aby]\$/)
  })

  it('does not claim to be loading after the request failed', async () => {
    vi.mocked(searchAdminUsers).mockRejectedValue(new Error('insufficient permission'))
    render(<AdminUsersPage />)
    await waitFor(() => expect(screen.getByRole('alert')).toBeTruthy())
    expect(screen.queryByText(/Loading accounts/)).toBeNull()
    expect(screen.queryByRole('table')).toBeNull()
  })

  it('does not claim to be loading after the overview request failed', async () => {
    vi.mocked(adminOverview).mockRejectedValue(new Error('insufficient permission'))
    render(<AdminOverviewPage />)
    await waitFor(() => expect(screen.getByRole('alert')).toBeTruthy())
    expect(screen.queryByText(/LOADING OVERVIEW/i)).toBeNull()
  })

  it('does not claim to be loading after the audit request failed', async () => {
    vi.mocked(listAdminAudit).mockRejectedValue(new Error('insufficient permission'))
    render(<AdminAuditPage />)
    await waitFor(() => expect(screen.getByRole('alert')).toBeTruthy())
    expect(screen.queryByText(/Loading the log/)).toBeNull()
  })

  it('searches by email, id, username or name', async () => {
    render(<AdminUsersPage />)
    await screen.findByText('ada@example.com')
    const input = screen.getByPlaceholderText(/Search by/)
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
        after_json: '{"role":"owner"}',
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

  it('offers no way to alter an event', async () => {
    render(<AdminAuditPage />)
    await screen.findByText('ADMIN_PROMOTE_USER')

    expect(document.querySelectorAll('input, textarea, select')).toHaveLength(0)
    for (const name of [/edit/i, /^delete/i, /remove/i, /correct/i, /annotate/i, /note/i]) {
      expect(screen.queryByRole('button', { name })).toBeNull()
      expect(screen.queryByRole('link', { name })).toBeNull()
    }

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
    expect(document.body.textContent).toMatch(/owner/)
  })

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
      users: 4, staff: 1, owners: 1, concepts: 657, domains: 15, questions: 120,
      reports: { open: 2, reviewing: 1, resolved: 5, dismissed: 0, outstanding: 3 },
      recent_activity: [],
    })
  })

  it('shows only counts the backend already has', async () => {
    render(<AdminOverviewPage />)
    await waitFor(() => expect(screen.getByText('Accounts')).toBeTruthy())
    expect(screen.getByText('657')).toBeTruthy()
    expect(screen.getByText('15')).toBeTruthy()
    expect(screen.getByText('120')).toBeTruthy()
    expect(screen.getByText('Reports needing attention')).toBeTruthy()
    for (const word of ['retention', 'cohort', 'revenue', 'churn', 'engagement']) {
      expect(document.body.textContent?.toLowerCase()).not.toContain(word)
    }
  })

  it('says plainly when there are no owners', async () => {
    vi.mocked(adminOverview).mockResolvedValue({
      users: 4, staff: 0, owners: 0, concepts: 657, domains: 15, questions: 120,
      reports: { open: 2, reviewing: 1, resolved: 5, dismissed: 0, outstanding: 3 },
      recent_activity: [],
    })
    render(<AdminOverviewPage />)
    await waitFor(() => expect(screen.getByRole('alert').textContent).toMatch(/No owners exist/))
    expect(screen.getByRole('alert').textContent).toMatch(/mathua admin promote/)
  })

  it('does not raise that alarm when one owner exists', async () => {
    render(<AdminOverviewPage />)
    await waitFor(() => expect(screen.getByText('Owners')).toBeTruthy())
    expect(screen.queryByText(/No owners exist/)).toBeNull()
  })
})

// ── Confirmation before a privileged change ────────────────────────────────

describe('role changes are confirmed before they happen', () => {
  it('does not call the API until the dialog is confirmed', async () => {
    vi.mocked(setUserRole).mockResolvedValue({ user: { ...LEARNER, role: 'admin' }, changed: true })
    render(<AdminUsersPage />)
    await screen.findByText('ada@example.com')

    changeRole('ada@example.com', 'admin')
    expect(setUserRole).not.toHaveBeenCalled()
    expect(await screen.findByRole('dialog')).toBeTruthy()
  })

  it('names the account and the consequence in the dialog', async () => {
    render(<AdminUsersPage />)
    await screen.findByText('ada@example.com')

    changeRole('ada@example.com', 'admin')
    const dialog = await screen.findByRole('dialog')
    expect(within(dialog).getByText('Change role to administrator?')).toBeTruthy()
    expect(dialog.textContent).toMatch(/Ada/)
    expect(dialog.textContent).toMatch(/takes effect on their next request/)
  })

  it('cancels without calling the API', async () => {
    render(<AdminUsersPage />)
    await screen.findByText('ada@example.com')

    changeRole('ada@example.com', 'admin')
    fireEvent.click(within(await screen.findByRole('dialog')).getByRole('button', { name: 'Cancel' }))

    await waitFor(() => expect(screen.queryByRole('dialog')).toBeNull())
    expect(setUserRole).not.toHaveBeenCalled()
  })

  it('focuses Cancel, not the confirming button', async () => {
    render(<AdminUsersPage />)
    await screen.findByText('ada@example.com')

    changeRole('ada@example.com', 'admin')
    const dialog = await screen.findByRole('dialog')
    await waitFor(() => expect(document.activeElement).toBe(within(dialog).getByRole('button', { name: 'Cancel' })))
  })

  it('sends the change once confirmed and reports it', async () => {
    vi.mocked(setUserRole).mockResolvedValue({ user: { ...LEARNER, role: 'admin' }, changed: true })
    render(<AdminUsersPage />)
    await screen.findByText('ada@example.com')

    changeRole('ada@example.com', 'admin')
    fireEvent.click(within(await screen.findByRole('dialog')).getByRole('button', { name: 'Make administrator' }))

    await waitFor(() => expect(setUserRole).toHaveBeenCalledTimes(1))
    expect(await screen.findByRole('status')).toHaveTextContent(/Ada is now administrator/)
  })

  it('reports when the account already had the role', async () => {
    vi.mocked(setUserRole).mockResolvedValue({ user: ADMIN, changed: false })
    render(<AdminUsersPage />)
    await screen.findByText('root@example.com')

    changeRole('root@example.com', 'student')
    fireEvent.click(within(await screen.findByRole('dialog')).getByRole('button', { name: 'Make learner' }))

    expect(await screen.findByRole('status')).toHaveTextContent(/already learner/)
  })

  it('closes the dialog and shows the server error', async () => {
    vi.mocked(setUserRole).mockRejectedValue(new Error('cannot remove the last owner'))
    render(<AdminUsersPage />)
    await screen.findByText('root@example.com')

    changeRole('root@example.com', 'student')
    fireEvent.click(within(await screen.findByRole('dialog')).getByRole('button', { name: 'Make learner' }))

    await waitFor(() => expect(screen.queryByRole('dialog')).toBeNull())
    expect(screen.getByRole('alert').textContent).toMatch(/last owner/)
  })

  it('shows the current role in the table for each account', async () => {
    render(<AdminUsersPage />)
    await screen.findByText('ada@example.com')
    const adminRow = within(screen.getByRole('table')).getByText('root@example.com').closest('tr')!
    expect(adminRow.textContent).toMatch(/admin/)
    const learnerRow = within(screen.getByRole('table')).getByText('ada@example.com').closest('tr')!
    expect(learnerRow.textContent).toMatch(/student/)
  })

  it('renders no credential-shaped value', async () => {
    render(<AdminUsersPage />)
    await screen.findByText('ada@example.com')
    expect(document.body.innerHTML).not.toMatch(/\$argon2|\$2[aby]\$|password_hash|share_token/)
  })
})
