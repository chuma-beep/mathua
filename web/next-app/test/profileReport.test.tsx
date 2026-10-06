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

  it('shows the per-domain picture', async () => {
    render(<ProfilePage />)
    expect(await screen.findByRole('region', { name: 'By domain' })).toBeTruthy()
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

  it('still owns the Next task, which is what makes it a hub', async () => {
    render(<ProfilePage />)
    const next = await screen.findByRole('region', { name: 'Next up' })
    expect(next.querySelectorAll('a').length).toBeGreaterThan(0)
  })


  it('answers "where am I" on the hub, above the Next task', async () => {
    render(<ProfilePage />)
    const pos = await screen.findByRole('region', { name: 'Where you are' })
    // Scored against what the learner can reach, and naming that denominator. The block used
    // to say "of 657 concepts mastered" above a bar computed against unlocked.
    expect(pos.textContent).toMatch(/of \d+ you can reach/)
    expect(pos.textContent).not.toMatch(/concepts mastered/)
    expect(pos.textContent).not.toMatch(/still locked/i)
    // It reads the same head the Next task does, so the position and the task
    // cannot disagree about what comes next.
    const regions = screen.getAllByRole('region').map(r => r.getAttribute('aria-label') ?? '')
    expect(regions.indexOf('Where you are')).toBeLessThan(regions.indexOf('Next up'))

    // ...and because it reads that same head, it needs no second button for it. Two CTAs for
    // one task twenty lines apart is not emphasis. (The review-debt link inside this block is a
    // different fact and stays, so the assertion is about the *head's* destination specifically.)
    const nextUp = screen.getByRole('region', { name: 'Next up' })
    const headHref = nextUp.querySelector('a[href]')?.getAttribute('href') ?? null
    const posHrefs = [...pos.querySelectorAll('a[href]')].map((a) => a.getAttribute('href'))
    expect(headHref).toBeTruthy()
    expect(posHrefs).not.toContain(headHref)
  })
  it('orders Activity ahead of the rest of the report', async () => {
    render(<ProfilePage />)
    const order = (await screen.findAllByRole('region'))
      .map(r => r.getAttribute('aria-label') ?? '')
      .filter(n => ['Activity', 'By domain', 'How you are doing'].includes(n))
    expect(order).toEqual(['Activity', 'By domain', 'How you are doing'])
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
    expect(screen.getByRole('region', { name: 'By domain' })).toBeTruthy()
    expect(screen.queryByRole('region', { name: 'How you are doing' })).toBeNull()
  })
})
