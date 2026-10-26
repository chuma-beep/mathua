import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import AdminLeaguesPage from '../app/admin/leagues/page'
import { getAdminLeagues } from '../lib/api'

vi.mock('../lib/api', async (importOriginal) => ({
  ...(await importOriginal<typeof import('../lib/api')>()),
  getAdminLeagues: vi.fn(),
}))

vi.mock('../hooks/useAdminMe', () => ({
  useAdminMe: () => ({
    me: { id: 'root', role: 'owner', permissions: ['admin.access', 'users.read'], assignable_roles: [] },
    loading: false,
    can: (p: string) => ['admin.access', 'users.read'].includes(p),
    isStaff: true,
  }),
}))

vi.mock('../components/Header', () => ({ default: () => <header>banner</header> }))
vi.mock('../components/BottomTabs', () => ({ default: () => <nav>tabs</nav> }))
vi.mock('../components/Footer', () => ({ default: () => <footer>footer</footer> }))
vi.mock('../components/AdminNav', () => ({ default: () => <nav>sections</nav> }))

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
      size: 2,
      movers: 0,
      members: [
        { rank: 1, name: 'Ada', username: 'ada', student_id: 'stu_ada', xp: 120 },
        { rank: 2, name: '', username: 'ghostly', student_id: 'stu_ghost', xp: 40 },
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
  it('lists every league without loading a roster', async () => {
    render(<AdminLeaguesPage />)

    expect(await screen.findByRole('button', { name: /^Bronze league 1/ })).toBeTruthy()
    expect(screen.getByText('Bronze')).toBeTruthy()
    expect(screen.getByText('Silver')).toBeTruthy()
    // The rule travels with the data, so the page cannot state a rule the transition does not apply.
    expect(screen.getByText(/30 learners, 5 up and 5 down/)).toBeTruthy()
    expect(screen.getByText(/10 learners, 3 up and 3 down/)).toBeTruthy()
    expect(screen.getByText(/2 leagues, 40 learners seated/)).toBeTruthy()

    // One request: the pyramid. Rosters are not fetched until a league is opened.
    expect(vi.mocked(getAdminLeagues)).toHaveBeenCalledTimes(1)
    expect(vi.mocked(getAdminLeagues).mock.calls[0][0]).not.toHaveProperty('tier')
  })

  it('opens one roster on demand and links each row to its account', async () => {
    render(<AdminLeaguesPage />)

    const toggle = await screen.findByRole('button', { name: /^Bronze league 1/ })
    expect(toggle.getAttribute('aria-expanded')).toBe('false')
    fireEvent.click(toggle)

    expect(await screen.findByText('Ada')).toBeTruthy()
    expect(screen.getByText('@ada')).toBeTruthy()
    // A learner without a name is shown by username rather than by raw id.
    expect(screen.getByText('ghostly')).toBeTruthy()
    expect(screen.getByText('120')).toBeTruthy()

    const link = screen.getByRole('link', { name: 'Ada' })
    expect(link.getAttribute('href')).toBe('/admin/users/detail?id=stu_ada')

    expect(vi.mocked(getAdminLeagues).mock.calls[1][0]).toMatchObject({ tier: 'bronze', league: 1 })

    // Closing and reopening uses the cached roster rather than asking again.
    fireEvent.click(toggle)
    fireEvent.click(toggle)
    await waitFor(() => expect(screen.getByText('Ada')).toBeTruthy())
    expect(vi.mocked(getAdminLeagues)).toHaveBeenCalledTimes(2)
  })

  it('walks back to an earlier week and refuses to walk past the current one', async () => {
    render(<AdminLeaguesPage />)
    await screen.findByRole('button', { name: /^Bronze league 1/ })

    const later = screen.getByRole('button', { name: 'Later week' }) as HTMLButtonElement
    expect(later.disabled).toBe(true)

    fireEvent.click(screen.getByRole('button', { name: 'Earlier week' }))
    await waitFor(() => expect(vi.mocked(getAdminLeagues)).toHaveBeenCalledTimes(2))
    const asked = String(vi.mocked(getAdminLeagues).mock.calls[1][0]?.week)
    expect(asked).toMatch(/^\d{4}-\d{2}-\d{2}$/)
    // Every week key the page can send is a Monday, because the server refuses anything else.
    expect(new Date(`${asked}T00:00:00Z`).getUTCDay()).toBe(1)
  })

  it('says so when a week has no league', async () => {
    vi.mocked(getAdminLeagues).mockResolvedValue({ ...pyramid, seats: 0, leagues: [] } as never)
    render(<AdminLeaguesPage />)
    expect(await screen.findByRole('status')).toBeTruthy()
    expect(screen.getByText(/No league was formed in this week/)).toBeTruthy()
  })
})
