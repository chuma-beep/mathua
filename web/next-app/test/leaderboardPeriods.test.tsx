import { fireEvent, render, screen, waitFor, within } from '@testing-library/react'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import type { LeaderboardBoard, MyLeague } from '../lib/api'
import { weekStartKey } from '../lib/leaderboard-period'
import LeaderboardPage from '../app/leaderboard/page'

const { getLeaderboard, getMyLeague, getSettingsMock, params } = vi.hoisted(() => ({
  getLeaderboard: vi.fn(),
  getMyLeague: vi.fn(),
  getSettingsMock: vi.fn(),
  params: { current: new URLSearchParams() },
}))

vi.mock('../lib/api', async (importOriginal) => {
  const actual = await importOriginal<typeof import('../lib/api')>()
  return { ...actual, getLeaderboard, getMyLeague, getSettings: getSettingsMock }
})

vi.mock('../hooks/useAuthState', () => ({
  useAuthState: () => ({ loggedIn: false, user: null }),
}))

vi.mock('next/navigation', () => ({
  usePathname: () => '/leaderboard',
  useRouter: () => ({ push: vi.fn() }),
  useSearchParams: () => params.current,
}))
vi.mock('next/link', () => ({
  default: ({ children, href }: { children: React.ReactNode; href: string }) => <a href={href}>{children}</a>,
}))
vi.mock('../hooks/useTheme', () => ({ useTheme: () => ({ mounted: true }) }))
vi.mock('../components/Header', () => ({ default: () => null }))
vi.mock('../components/Footer', () => ({ default: () => null }))
vi.mock('../components/BottomTabs', () => ({ default: () => null }))

function currentWeek() {
  const start = weekStartKey(new Date())
  const end = new Date(`${start}T00:00:00Z`)
  end.setUTCDate(end.getUTCDate() + 7)
  return { start, end: end.toISOString() }
}

function board(overrides: Partial<LeaderboardBoard> = {}): LeaderboardBoard {
  const week = currentWeek()
  return {
    view: 'weekly',
    period_kind: 'weekly',
    period_start: `${week.start}T00:00:00Z`,
    period_end: week.end,
    final: false,
    earliest_anchor: '2026-01-01',
    page: 1,
    page_size: 25,
    total: 1,
    total_pages: 1,
    entries: [{ rank: 1, id: 'ada', name: 'Ada', xp: 120, xp_total: 900, mastered: 4, level: 'Scholar' }],
    ...overrides,
  }
}

function league(overrides: Partial<MyLeague> = {}): MyLeague {
  const week = currentWeek()
  return {
    week: week.start,
    week_start: `${week.start}T00:00:00Z`,
    week_end: week.end,
    final: true,
    earliest_week: '2026-01-05',
    movement: 'up',
    from_tier: 'bronze',
    to_tier: 'silver',
    tier: 'silver',
    league_no: 2,
    size: 2,
    promote: 1,
    demote: 1,
    members: [
      { rank: 1, id: 'ada', name: 'Ada', xp: 30 },
      { rank: 2, id: 'bob', name: 'Bob', xp: 20 },
    ],
    me: { rank: 1, xp: 30 },
    ...overrides,
  }
}

beforeEach(() => {
  params.current = new URLSearchParams()
  getLeaderboard.mockReset().mockResolvedValue(board())
  getMyLeague.mockReset().mockResolvedValue(league())
  getSettingsMock.mockReset().mockResolvedValue({})
})

describe('leaderboard periods', () => {
  it('restores a monthly archive from the URL', async () => {
    params.current = new URLSearchParams('view=monthly&anchor=2026-10-01')
    getLeaderboard.mockResolvedValue(
      board({
        view: 'monthly',
        period_kind: 'monthly',
        period_start: '2026-10-01T00:00:00Z',
        period_end: '2026-11-01T00:00:00Z',
        final: true,
      }),
    )
    render(<LeaderboardPage />)
    await waitFor(() => expect(getLeaderboard).toHaveBeenCalledWith('monthly', 1, 25, '2026-10-01'))
    expect(await screen.findByText('October 2026')).toBeInTheDocument()
    expect(screen.getAllByText('Final').length).toBeGreaterThan(0)
  })

  it('hides Month until a past month has data', async () => {
    getLeaderboard.mockResolvedValue(board({ earliest_anchor: currentWeek().start }))
    render(<LeaderboardPage />)
    await screen.findByRole('table')
    expect(screen.queryByRole('tab', { name: 'Month' })).toBeNull()
  })

  it('offers Month once an earlier month exists', async () => {
    render(<LeaderboardPage />)
    await screen.findByRole('table')
    expect(screen.getByRole('tab', { name: 'Month' })).toBeInTheDocument()
  })

  it('snaps a picked day to that week Monday', async () => {
    const today = new Date()
    const selectedMonday = new Date(`${currentWeek().start}T00:00:00Z`).getUTCDate()
    const day = today.getDate() === selectedMonday ? today.getDate() - 1 : today.getDate()
    const expected = weekStartKey(new Date(today.getFullYear(), today.getMonth(), day))
    render(<LeaderboardPage />)
    await screen.findByRole('table')

    fireEvent.click(screen.getByRole('button', { name: 'Choose a week' }))
    const grid = await screen.findByRole('grid')
    const dayButton = within(grid)
      .getAllByRole('button')
      .find((button) => button.textContent?.trim() === String(day))
    expect(dayButton).toBeTruthy()
    fireEvent.click(dayButton as HTMLButtonElement)

    await waitFor(() => expect(getLeaderboard).toHaveBeenLastCalledWith('weekly', 1, 25, expected))
  })

  it('shows a closed league week from the snapshot and its outcome', async () => {
    render(<LeaderboardPage />)
    await screen.findByRole('table')
    fireEvent.mouseDown(screen.getByRole('tab', { name: 'My league' }))

    expect(await screen.findByText('silver league 2')).toBeInTheDocument()
    expect(screen.getByText('Promoted to silver')).toBeInTheDocument()
    expect(screen.getByText(/This week is closed/)).toBeInTheDocument()
  })
})
