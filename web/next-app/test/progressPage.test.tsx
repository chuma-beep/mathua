import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor, within } from '@testing-library/react'
import ProgressPage from '../app/progress/page'
import * as api from '../lib/api'
import { getUserInfo, ensureGuestId, ensureGuestToken } from '../lib/auth'

vi.mock('../lib/api', async (importOriginal) => {
  const mod = await importOriginal<typeof import('../lib/api')>()
  return {
    ...mod,
    getScores: vi.fn(),
    getProgress: vi.fn(),
    getActivity: vi.fn(),
    getWeaknesses: vi.fn(),
    getEfficacy: vi.fn(),
    getEfficacyTrend: vi.fn(),
  }
})

// Header and BottomTabs read the pathname to mark the active destination.
vi.mock('next/navigation', () => ({
  usePathname: () => '/progress',
  useRouter: () => ({ push: vi.fn(), replace: vi.fn() }),
  useSearchParams: () => new URLSearchParams(),
}))

vi.mock('../lib/auth', async (importOriginal) => {
  const mod = await importOriginal<typeof import('../lib/auth')>()
  return {
    ...mod,
    getUserInfo: vi.fn(() => ({ student_id: 's1', name: 'T', username: 't', concepts_mastered: 0, current_streak: 0, level: 'Beginner', diagnostic_completed: true })),
    getGuestId: vi.fn(() => null),
    ensureGuestId: vi.fn(),
    ensureGuestToken: vi.fn(),
  }
})

// x has 2 concepts, y has 1; x.b requires x.a.
const progress = {
  'x.a': { concept_id: 'x.a', status: 'MASTERED', streak: 3 },
  'x.b': { concept_id: 'x.b', status: 'PRACTICING', streak: 1 },
} as Record<string, unknown>

beforeEach(() => {
  vi.mocked(api.getScores).mockResolvedValue({ concepts_mastered: 1, xp_today: 2, daily_xp_goal: 10, level: 'Beginner', current_streak: 1, lifetime_points: 0, weekly_score: 0, speed_bonus: 0, xp_total: 12 } as never)
  // A trimmed catalogue keeps the arithmetic in the assertions checkable.
  vi.mocked(api.getProgress).mockResolvedValue(progress as never)
  vi.mocked(api.getActivity).mockResolvedValue([] as never)
  vi.mocked(api.getWeaknesses).mockResolvedValue({ by_domain: {} } as never)
  vi.mocked(api.getEfficacy).mockResolvedValue({
    first_pass_rate: 0.5, second_pass_rate: 1, avg_attempts_per_concept: 2, concepts_touched: 2, total_attempts: 4,
  } as never)
  vi.mocked(api.getEfficacyTrend).mockResolvedValue(null as never)
})

describe('/progress report', () => {
  it('leads with "Where you are" and answers the headline question', async () => {
    render(<ProgressPage />)
    const pos = await screen.findByRole('region', { name: 'Where you are' })
    // Position is the first section: nothing else may come before it.
    const sections = screen.getAllByRole('region')
    expect(sections[0]).toBe(pos)
    expect(within(pos).getByText('of 657 concepts mastered')).toBeTruthy()
  })

  it('reports mastered, unlocked and locked separately', async () => {
    render(<ProgressPage />)
    const pos = await screen.findByRole('region', { name: 'Where you are' })
    // The trimmed catalogue is injected below, so assert the labels exist and
    // the numbers are consistent with each other rather than hardcoded.
    for (const label of ['Mastered', 'Unlocked', 'In progress', 'Still locked']) {
      expect(within(pos).getByText(label)).toBeTruthy()
    }
  })

  it('orders the report by the questions a learner asks', async () => {
    render(<ProgressPage />)
    await screen.findByRole('region', { name: 'Where you are' })
    const names = screen
      .getAllByRole('region')
      .map(r => r.getAttribute('aria-label') ?? r.getAttribute('id') ?? '')
    // Position first, then the distribution, then how you are doing, then the
    // history. The hub owns "what next"; this page owns "where am I".
    const order = ['Where you are', 'By domain', 'How you are doing', 'Activity']
    expect(names.filter(n => order.includes(n))).toEqual(order)
  })

  it('does not fetch cohort analytics onto an individual page', async () => {
    // total_students and retention_rate belong to /docs/efficacy, not here.
    render(<ProgressPage />)
    await screen.findByRole('region', { name: 'Where you are' })
    expect(api.getEfficacyTrend).not.toHaveBeenCalled()
    expect(screen.queryByText(/learners? ·/)).toBeNull()
    expect(screen.queryByText(/returning/)).toBeNull()
  })

  it('offers no share or reaction controls', () => {
    // The vocabulary of engagement mechanics is banned on a progress report.
    render(<ProgressPage />)
    for (const word of ['Share', 'Like', 'Celebrate', 'Streak 🔥', 'Keep it up']) {
      expect(screen.queryByText(new RegExp(word, 'i'))).toBeNull()
    }
  })

  it('establishes a guest identity so a guest sees their own numbers', async () => {
    render(<ProgressPage />)
    await screen.findByRole('region', { name: 'Where you are' })
    // Signed in here, so the guest path must not have run.
    expect(vi.mocked(ensureGuestId)).not.toHaveBeenCalled()
    await waitFor(() => expect(api.getScores).toHaveBeenCalledWith('s1'))
  })

  it('links onward to the full answer history', async () => {
    render(<ProgressPage />)
    const link = await screen.findByRole('link', { name: /Every question you.ve answered/ })
    expect(link.getAttribute('href')).toBe('/history')
  })
})