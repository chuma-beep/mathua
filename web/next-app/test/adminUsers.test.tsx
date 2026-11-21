import { beforeEach, describe, expect, it, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import AdminUsersPage from '../app/admin/users/page'
import { searchAdminUsers } from '../lib/api'
import type { AdminUser } from '../lib/api'

vi.mock('next/navigation', () => ({
  usePathname: () => '/admin/users',
}))

vi.mock('../lib/api', async (importOriginal) => ({
  ...(await importOriginal<typeof import('../lib/api')>()),
  searchAdminUsers: vi.fn(),
  setUserRole: vi.fn(),
}))

vi.mock('../hooks/useAdminMe', () => ({
  useAdminMe: () => ({
    me: { id: 'root', role: 'owner', permissions: ['admin.access'], assignable_roles: ['admin'] },
    loading: false,
    can: () => true,
    isStaff: true,
  }),
}))

vi.mock('../hooks/useAuthState', () => ({
  useAuthState: () => ({ user: { student_id: 's1' }, loading: false }),
}))

vi.mock('../components/app-shell', () => ({
  default: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
}))
vi.mock('../components/Footer', () => ({ default: () => <footer /> }))

function user(over: Partial<AdminUser> & { id: string }): AdminUser {
  return {
    name: 'Ada',
    username: 'ada',
    email: 'ada@example.com',
    role: 'student',
    created_at: '2026-02-03T04:05:06Z',
    has_password: true,
    ...over,
  }
}

beforeEach(() => {
  vi.clearAllMocks()
  vi.mocked(searchAdminUsers).mockResolvedValue({
    users: [
      user({
        id: 's1',
        avatar: { avatar_dicebear: { style: 'notionists', seed: 'ada-seed' } },
        last_seen_at: '2026-11-20T18:30:00Z',
        device: 'mobile',
        country: 'NG',
      }),
      user({ id: 's2', name: 'Bob', username: 'bob', email: 'bob@example.com', role: 'owner' }),
    ],
  })
})

describe('the account list', () => {
  it('draws the picture each account chose and an initial when it has none', async () => {
    render(<AdminUsersPage />)
    await screen.findByText('Ada')

    const picked = screen.getByRole('img', { name: 'Ada avatar' })
    expect(picked.getAttribute('src')).toBe('https://api.dicebear.com/10.x/notionists/svg?seed=ada-seed')

    // No upload and no pick: the generated initial is still an avatar, not a gap.
    expect(screen.getByLabelText('Bob avatar')).toBeTruthy()
  })

  it('shows when the account joined without opening it', async () => {
    render(<AdminUsersPage />)
    await screen.findByText('Ada')

    const joined = document.querySelector('time[datetime="2026-02-03T04:05:06Z"]')
    expect(joined).toBeTruthy()
    expect(joined?.textContent).toMatch(/Joined/)
    expect(joined?.textContent).not.toContain('2026-02-03T')
  })

  it('names the role in language a reader understands rather than the storage value', async () => {
    render(<AdminUsersPage />)
    await screen.findByText('Ada')

    expect(screen.getAllByText('learner').length).toBeGreaterThan(0)
    expect(screen.getAllByText('owner').length).toBeGreaterThan(0)
  })

  it('shows the last visit, device class and country without opening the account', async () => {
    render(<AdminUsersPage />)
    await screen.findByText('Ada')

    const seen = document.querySelector('time[datetime="2026-11-20T18:30:00Z"]')
    expect(seen).toBeTruthy()
    expect(seen?.textContent).toMatch(/Last seen/)
    expect(seen?.textContent).not.toContain('2026-11-20T')
    expect(screen.getByText('mobile')).toBeTruthy()
    expect(screen.getByText('Nigeria')).toBeTruthy()
  })
})
