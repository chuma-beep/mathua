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
    getEfficacy: vi.fn(async () => ({
      concepts_touched: 3, first_pass_rate: 0.5, second_pass_rate: 1,
      avg_attempts_per_concept: 2, total_attempts: 6,
    })),
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

// The report belongs on the learner's own page. A separate /progress
// destination split the picture in two and left the hub showing nothing but a
// task queue, which is not what a learner opens their profile to see.
describe('/profile shows the report', () => {
  beforeEach(() => {
    vi.mocked(getUserInfo).mockReturnValue({
      student_id: 's1', name: 'T', username: 't', concepts_mastered: 4,
      current_streak: 2, level: 'Learner', diagnostic_completed: true,
    } as never)
    vi.mocked(api.getScores).mockResolvedValue(scores)
  })

  it('shows the activity heatmap', async () => {
    render(<ProfilePage />)
    expect(await screen.findByRole('region', { name: 'Activity' })).toBeTruthy()
  })

  // Per-domain progress moved to /domains. It sat here as a 15-row table of subject
  // percentages competing with the recommendation on a page whose job is "what next", and its
  // hardcoded domain list omitted precalculus and machine_learning — 50 concepts invisible.
  it('does not carry the per-domain table, which is its own page now', async () => {
    render(<ProfilePage />)
    await screen.findByRole('region', { name: 'Activity' })
    expect(screen.queryByRole('region', { name: 'By domain' })).toBeNull()
    expect(screen.queryByText('Domain Progress')).toBeNull()
  })

  it('shows the learner their own accuracy, and no cohort figures', async () => {
    render(<ProfilePage />)
    const doing = await screen.findByRole('region', { name: 'How you are doing' })
    expect(doing.textContent).toMatch(/First-pass/)
    // total_students and retention_rate are other learners' numbers; they live
    // on /docs/efficacy, never on someone's own page.
    expect(doing.textContent).not.toMatch(/returning/)
    expect(api.getEfficacyTrend).not.toHaveBeenCalled()
  })

  it('still owns the learning tasks, which is what makes it a hub', async () => {
    render(<ProfilePage />)
    const learning = await screen.findByRole('region', { name: 'Your learning' })
    // The cards are buttons (selection asks the server for eligibility first),
    // and the "Choose what to learn" affordance is always present.
    expect(learning.querySelectorAll('button').length).toBeGreaterThan(0)
    expect(learning.querySelectorAll('a').length).toBeGreaterThan(0)
  })


  it('answers "where am I" on the hub, above the learning tasks', async () => {
    render(<ProfilePage />)
    const pos = await screen.findByRole('region', { name: 'Where you are' })
    // Scored against what the learner can reach, and naming that denominator. The block used
    // to say "of 657 concepts mastered" above a bar computed against unlocked.
    expect(pos.textContent).toMatch(/of \d+ you can reach/)
    expect(pos.textContent).not.toMatch(/concepts mastered/)
    expect(pos.textContent).not.toMatch(/still locked/i)
    // It reads the same head the learning tasks do, so the position and the tasks
    // cannot disagree about what comes next.
    const regions = screen.getAllByRole('region').map(r => r.getAttribute('aria-label') ?? '')
    expect(regions.indexOf('Where you are')).toBeLessThan(regions.indexOf('Your learning'))
  })
  it('orders Activity ahead of the rest of the report', async () => {
    render(<ProfilePage />)
    const order = (await screen.findAllByRole('region'))
      .map(r => r.getAttribute('aria-label') ?? '')
      .filter(n => ['Activity', 'How you are doing'].includes(n))
    expect(order).toEqual(['Activity', 'How you are doing'])
  })

  it('does not send the learner to a separate progress page', async () => {
    render(<ProfilePage />)
    await screen.findByRole('region', { name: 'Activity' })
    // Scoped to the page body: the sidebar is site chrome, and its Progress
    // entry is a separate change. What matters here is that the report's own
    // content does not point somewhere else.
    const main = document.getElementById('profile-main')!
    expect(main).toBeTruthy()
    for (const link of main.querySelectorAll('a')) {
      expect(link.getAttribute('href')).not.toBe('/progress')
    }
  })

  it('shows a guest the heatmap and the domains, but not accuracy it has no data for', async () => {
    vi.mocked(getUserInfo).mockReturnValue(null)
    render(<ProfilePage />)
    expect(await screen.findByRole('region', { name: 'Activity' })).toBeTruthy()
    expect(screen.queryByRole('region', { name: 'By domain' })).toBeNull()
    expect(screen.queryByRole('region', { name: 'How you are doing' })).toBeNull()
  })
})
