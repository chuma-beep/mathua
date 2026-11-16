import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent, waitFor, within } from '@testing-library/react'
import AdminLeaguesPage from '../app/admin/leagues/page'
import { getAdminLeagues, getAdminLeagueRange, API_BASE } from '../lib/api'

vi.mock('../lib/api', async (importOriginal) => ({
  ...(await importOriginal<typeof import('../lib/api')>()),
  getAdminLeagues: vi.fn(),
  getAdminLeagueRange: vi.fn(),
}))

vi.mock('../hooks/useAdminMe', () => ({
  useAdminMe: () => ({
    me: { id: 'root', role: 'owner', permissions: ['admin.access', 'users.read'], assignable_roles: [] },
    loading: false,
    can: (p: string) => ['admin.access', 'users.read'].includes(p),
    isStaff: true,
  }),
}))


// The owner's league page.
//
// Three behaviours matter and they are the reason the page exists: the pyramid is listed without
// loading every roster, a league opens its standings on demand and each row reaches the account
// it belongs to. A page that fetched every roster up front would still pass a "does it render"
// test, so the mock counts the calls instead.

// The page asks for the live week, so the fixture must answer with that same
// week. Pinning a past Monday made "one request" become two once the date moved
// on, because the page adopted the response week and reloaded.
function currentWeekKey(): string {
  const now = new Date()
  const day = new Date(Date.UTC(now.getUTCFullYear(), now.getUTCMonth(), now.getUTCDate()))
  day.setUTCDate(day.getUTCDate() - ((day.getUTCDay() + 6) % 7))
  return day.toISOString().slice(0, 10)
}

const thisWeek = currentWeekKey()
const weekEnd = new Date(`${thisWeek}T00:00:00Z`)
weekEnd.setUTCDate(weekEnd.getUTCDate() + 7)

const pyramid = {
  week: thisWeek,
  week_start: `${thisWeek}T00:00:00Z`,
  week_end: weekEnd.toISOString(),
  promote: 5,
  demote: 5,
  seats: 40,
  leagues: [
    { tier: 'bronze', league_no: 1, size: 30, movers: 5 },
    { tier: 'silver', league_no: 1, size: 10, movers: 3 },
  ],
}

const roster = {
  ...pyramid,
  leagues: [
    {
      tier: 'bronze',
      league_no: 1,
      size: 3,
      movers: 1,
      members: [
        {
          rank: 1,
          name: 'Ada',
          username: 'ada',
          student_id: 'stu_ada',
          xp: 120,
          avatar_url: 'https://images.example/ada.jpg',
        },
        {
          rank: 2,
          name: '',
          username: 'ghostly',
          student_id: 'stu_ghost',
          xp: 40,
          avatar_dicebear: { style: 'notionists', seed: 'ghost-seed' },
        },
        {
          rank: 3,
          name: 'Custom',
          username: 'custom',
          student_id: 'stu_custom',
          xp: 10,
          avatar_url: '/api/avatar/stu_custom?v=3',
        },
      ],
    },
  ],
}

beforeEach(() => {
  vi.mocked(getAdminLeagues).mockReset()
  vi.mocked(getAdminLeagues).mockImplementation(async (query = {}) =>
    query.tier ? (roster as never) : (pyramid as never),
  )
})

describe('admin leagues', () => {
  it('lists every league without loading every roster', async () => {
    render(<AdminLeaguesPage />)

    await screen.findByText('Tier pyramid')
    const picker = screen.getByRole('combobox', { name: 'League roster' })
    expect(within(picker).getAllByRole('option').length).toBe(2)
    expect(screen.getAllByText('Bronze').length).toBeGreaterThan(0)
    expect(screen.getAllByText('Silver').length).toBeGreaterThan(0)
    expect(screen.getByText('5 / 5')).toBeTruthy()
    expect(screen.getByText(/of 60 seats/)).toBeTruthy()
    expect(screen.getByText(/67% full/)).toBeTruthy()

    // The pyramid and the selected league's roster, not every roster.
    await waitFor(() => expect(vi.mocked(getAdminLeagues)).toHaveBeenCalledTimes(2))
    expect(vi.mocked(getAdminLeagues).mock.calls[0][0]).not.toHaveProperty('tier')
  })

  it('loads one roster at a time and links each row to its account', async () => {
    render(<AdminLeaguesPage />)

    const picker = await screen.findByRole('combobox', { name: 'League roster' })
    expect(await screen.findByText('Ada')).toBeTruthy()
    expect(screen.getByText('@ada')).toBeTruthy()
    expect(screen.getByText('ghostly')).toBeTruthy()
    expect(screen.getByText('120')).toBeTruthy()

    const adaAvatar = screen.getByRole('img', { name: 'Ada avatar' })
    expect(adaAvatar.getAttribute('src')).toBe('https://images.example/ada.jpg')
    const ghostAvatar = screen.getByRole('img', { name: 'ghostly avatar' })
    expect(ghostAvatar.getAttribute('src')).toContain('notionists')
    expect(ghostAvatar.getAttribute('src')).toContain('ghost-seed')
    const customAvatar = screen.getByRole('img', { name: 'Custom avatar' })
    expect(customAvatar.getAttribute('src')).toBe(`${API_BASE}/api/avatar/stu_custom?v=3`)

    const link = screen.getByRole('link', { name: 'Ada' })
    expect(link.getAttribute('href')).toBe('/admin/users/detail?id=stu_ada')

    // Movement is an arrow and a word, not a full-row tint. The row keeps the
    // same surface as every other row.
    const up = screen.getAllByText('Up')[0].closest('tr')!
    expect(up.className).not.toMatch(/green|red/)

    // Selecting the silver league asks for that roster, not every roster.
    fireEvent.change(picker, { target: { value: 'silver/1' } })
    await waitFor(() => expect(vi.mocked(getAdminLeagues).mock.calls.at(-1)?.[0]).toMatchObject({ tier: 'silver', league: 1 }))
  })

  it('walks back to an earlier week and refuses to walk past the current one', async () => {
    render(<AdminLeaguesPage />)
    await screen.findByText('Tier pyramid')

    const next = screen.getByRole('button', { name: 'Next' }) as HTMLButtonElement
    expect(next.disabled).toBe(true)

    fireEvent.click(screen.getByRole('button', { name: 'Previous' }))
    await waitFor(() => expect(vi.mocked(getAdminLeagues)).toHaveBeenCalledTimes(3))
    const asked = String(vi.mocked(getAdminLeagues).mock.calls[2][0]?.week)
    expect(asked).toMatch(/^\d{4}-\d{2}-\d{2}$/)
    expect(new Date(`${asked}T00:00:00Z`).getUTCDay()).toBe(1)
  })

  it('says so when a week has no league', async () => {
    vi.mocked(getAdminLeagues).mockResolvedValue({ ...pyramid, seats: 0, leagues: [] } as never)
    render(<AdminLeaguesPage />)
    expect(await screen.findByRole('status')).toBeTruthy()
    expect(screen.getByText(/No league was formed in this week/)).toBeTruthy()
  })
})


describe('league time scopes', () => {
  it('switches to the monthly aggregate', async () => {
    vi.mocked(getAdminLeagues).mockResolvedValue(pyramid)
    vi.mocked(getAdminLeagueRange).mockResolvedValue({
      from: thisWeek,
      to: thisWeek,
      group: 'week',
      weeks: 2,
      seats: 40,
      unique_learners: 12,
      promoted: 3,
      demoted: 2,
      stayed: 30,
      new: 5,
      tiers: [{ tier: 'bronze', seats: 30, unique_learners: 10 }],
      buckets: [{ key: thisWeek, label: thisWeek, weeks: 1, seats: 20, unique_learners: 8, promoted: 2, demoted: 1, stayed: 15, new: 2 }],
      top_range: [{ student_id: 's1', name: 'Ada', username: 'ada', xp: 120, tier: 'bronze' }],
      top_latest: [],
    })
    render(<AdminLeaguesPage />)
    fireEvent.click(screen.getByRole('button', { name: 'month' }))
    expect(await screen.findByText('League weeks')).toBeTruthy()
    expect(screen.getByText('Unique learners')).toBeTruthy()
    expect(screen.getByRole('button', { name: 'Latest week' })).toBeTruthy()
  })
})
