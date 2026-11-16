import { beforeEach, describe, expect, it, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import AdminUserDetailPage from '../app/admin/users/detail/page'
import { getAdminUser } from '../lib/api'
import type { AdminUserDetail } from '../lib/api'

const params = new URLSearchParams('id=s1')

vi.mock('next/navigation', () => ({
  usePathname: () => '/admin/users/detail',
  useSearchParams: () => params,
}))

vi.mock('../lib/api', async (importOriginal) => ({
  ...(await importOriginal<typeof import('../lib/api')>()),
  getAdminUser: vi.fn(),
}))

vi.mock('../hooks/useAdminMe', () => ({
  useAdminMe: () => ({
    me: { id: 'root', role: 'owner', permissions: ['admin.access'], assignable_roles: [] },
    loading: false,
    can: () => true,
    isStaff: true,
  }),
}))

vi.mock('../components/app-shell', () => ({
  default: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
}))
vi.mock('../components/Footer', () => ({ default: () => <footer /> }))
vi.mock('../components/AdminNav', () => ({ default: () => <nav aria-label="Admin sections" /> }))

const BASE: AdminUserDetail = {
  user: {
    id: 's1',
    name: 'Ada',
    username: 'ada',
    email: 'ada@example.com',
    role: 'student',
    created_at: '2026-02-03T04:05:06Z',
    has_password: false,
  },
  avatar: { avatar_dicebear: { style: 'notionists', seed: 'ada-seed' } },
  history: [],
  activity: {
    xp_total: 10,
    xp_today: 2,
    concepts_touched: 3,
    concepts_mastered: 1,
    attempts: 5,
    last_attempt_at: '',
    last_active_day: '',
    diagnostic_completed_at: '',
    reports_filed: 0,
  },
  reports: [],
}

describe('admin account detail avatar', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    vi.mocked(getAdminUser).mockResolvedValue(BASE)
  })

  it('draws the learner avatar from their DiceBear pick', async () => {
    render(<AdminUserDetailPage />)
    await screen.findByRole('heading', { name: 'Ada' })
    const avatar = screen.getByRole('img', { name: 'Ada avatar' })
    expect(avatar.getAttribute('src')).toBe('https://api.dicebear.com/10.x/notionists/svg?seed=ada-seed')
  })

  it('keeps a custom avatar on the API origin', async () => {
    vi.mocked(getAdminUser).mockResolvedValue({
      ...BASE,
      avatar: { avatar_url: '/api/avatar/s1?v=7' },
    })
    render(<AdminUserDetailPage />)
    const avatar = await screen.findByRole('img', { name: 'Ada avatar' })
    expect(avatar.getAttribute('src')).toBe('/api/avatar/s1?v=7')
  })
})
