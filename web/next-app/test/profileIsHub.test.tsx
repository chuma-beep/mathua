import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen } from '@testing-library/react'
import ProfilePage from '../app/profile/page'
import * as api from '../lib/api'
import { getUserInfo } from '../lib/auth'

vi.mock('../lib/api', async (importOriginal) => {
  const mod = await importOriginal<typeof import('../lib/api')>()
  return {
    ...mod,
    getScores: vi.fn(),
    getProgress: vi.fn(async () => ({})),
    getActivity: vi.fn(async () => []),
    getWeaknesses: vi.fn(async () => ({ by_domain: {} })),
    getDueReviews: vi.fn(async () => ({ count: 0 })),
    getSettings: vi.fn(async () => ({})),
    getEfficacy: vi.fn(async () => ({ concepts_touched: 3, first_pass_rate: 0.5, second_pass_rate: 1, avg_attempts_per_concept: 2, total_attempts: 6 })),
    getEfficacyTrend: vi.fn(async () => null),
  }
})

vi.mock('../lib/auth', async (importOriginal) => {
  const mod = await importOriginal<typeof import('../lib/auth')>()
  return {
    ...mod,
    getUserInfo: vi.fn(),
    getGuestId: vi.fn(() => 'guest_1'),
    ensureGuestId: vi.fn(),
    ensureGuestToken: vi.fn(),
    isLoggedIn: vi.fn(() => false),
  }
})

// Link and the sidebar read the router; stub the three hooks the app shell uses.
vi.mock('next/navigation', () => ({
  usePathname: () => '/profile',
  useRouter: () => ({ push: vi.fn(), replace: vi.fn() }),
  useSearchParams: () => new URLSearchParams(),
}))

vi.mock('../components/DailyGoalControl', () => ({
  default: () => null,
  getGuestGoal: () => null,
}))

const scores = {
  lifetime_points: 100, weekly_score: 10, speed_bonus: 0, concepts_mastered: 4,
  current_streak: 2, level: 'Learner', xp_total: 120, xp_today: 6, daily_xp_goal: 10,
} as never

// ADR-001 locks /profile as the hub and ADR-023 as a compact one. The report
// is /progress. If the heatmap, the domain picture or the efficacy tiles come
// back to the hub, the hub is doing five jobs again — which is exactly what
// these assertions exist to prevent.
describe('/profile is a hub, not a report', () => {
  beforeEach(() => {
    vi.mocked(getUserInfo).mockReturnValue({
      student_id: 's1', name: 'T', username: 't', concepts_mastered: 4,
      current_streak: 2, level: 'Learner', diagnostic_completed: true,
    } as never)
    vi.mocked(api.getScores).mockResolvedValue(scores)
  })

  it('renders no activity heatmap', async () => {
    render(<ProfilePage />)
    await screen.findByRole('heading', { name: /Where you are/i }, { timeout: 1 }).catch(() => {})
    expect(screen.queryByRole('heading', { name: 'Activity' })).toBeNull()
  })

  it('renders no domain-progress section', async () => {
    render(<ProfilePage />)
    await screen.findByText('Learner', { exact: false }).catch(() => {})
    expect(screen.queryByText(/concepts mastered.*·.*done.*·.*locked/)).toBeNull()
  })

  it('renders no efficacy instrumentation', async () => {
    render(<ProfilePage />)
    await screen.findByText('Learner', { exact: false }).catch(() => {})
    expect(screen.queryByText('First-pass')).toBeNull()
    expect(screen.queryByText(/returning/)).toBeNull()
  })

  it('does not fetch the cohort trend at all', async () => {
    render(<ProfilePage />)
    await screen.findByText('Learner', { exact: false }).catch(() => {})
    expect(api.getEfficacyTrend).not.toHaveBeenCalled()
  })

  it('still owns the Next task, and points at the report', async () => {
    render(<ProfilePage />)
    // The hub keeps what a hub is for: the Next up region, per ADR-023.
    const next = await screen.findByRole('region', { name: 'Next up' })
    // The hub keeps what a hub is for: a Next up region with a way in.
    expect(next.querySelectorAll('a').length).toBeGreaterThan(0)
    const link = screen.getByRole('link', { name: /See your progress/ })
    expect(link.getAttribute('href')).toBe('/progress')
  })

  it('offers a guest the same way onward', async () => {
    vi.mocked(getUserInfo).mockReturnValue(null)
    render(<ProfilePage />)
    const links = await screen.findAllByRole('link', { name: /See your progress/ })
    expect(links.length).toBeGreaterThan(0)
    expect(links[0].getAttribute('href')).toBe('/progress')
  })
})