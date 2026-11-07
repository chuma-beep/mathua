import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import ReportsPage from '../app/reports/page'
import { getLearnerReport, getLearnerReportHistory } from '@/lib/api'
import type { LearnerReport } from '@/lib/api'

vi.mock('next/link', () => ({
  default: ({ href, children, ...rest }: { href: unknown; children: React.ReactNode }) => (
    <a href={typeof href === 'string' ? href : ''} {...rest}>{children}</a>
  ),
}))

vi.mock('@/components/app-shell', () => ({
  default: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
}))
vi.mock('@/components/Footer', () => ({ default: () => null }))
vi.mock('@/components/SectionHeader', () => ({
  default: ({ title }: { title?: string }) => <h2>{title}</h2>,
}))

vi.mock('@/lib/api', async () => {
  const actual = await vi.importActual<typeof import('../lib/api')>('../lib/api')
  return { ...actual, getLearnerReport: vi.fn(), getLearnerReportHistory: vi.fn() }
})

const baseReport: LearnerReport = {
  generated_at: '2026-10-24T21:00:00Z',
  learner: { id: 'st1', name: 'Ada' },
  period: { days: 30, since: '2026-09-24T21:00:00Z', until: '2026-10-24T21:00:00Z', label: 'Last 30 days' },
  overview: {
    total: 10, mastered: 4, in_progress: 2, unlocked: 3, locked: 1,
    xp_lifetime: 1200, xp_period: 300, activity_days: 2, questions: 12, correct: 9,
  },
  domains: [{ id: 'algebra', conceptCount: 6, masteredCount: 3, unlockedCount: 2, inProgressCount: 1 }],
  mastery: [
    {
      concept_id: 'alg.linear', label: 'Linear equations', domain: 'algebra', state: 'mastered', status: 'MASTERED',
      decaying: false, mastery_pct: 1, mastered_at: '2026-10-20T10:00:00Z', attempts: 9, streak: 4,
      last_attempted: '2026-10-20T10:00:00Z', next_review: '2026-11-01T10:00:00Z',
      placement_seeded: false, provisional: false, prerequisite_gaps: [], href: '/study?concept=alg.linear',
    },
    {
      concept_id: 'alg.quadratic', label: 'Quadratic equations', domain: 'algebra', state: 'locked', status: 'UNSEEN',
      decaying: false, mastery_pct: 0, attempts: 0, streak: 0, placement_seeded: false, provisional: false,
      prerequisite_gaps: [{ concept_id: 'alg.linear', label: 'Linear equations', state: 'in_progress' }],
      href: '/study?concept=alg.quadratic',
    },
  ],
  assessments: [
    {
      id: 'quiz:1', type: 'quiz', title: 'Mastery check', status: 'completed',
      started_at: '2026-10-22T09:00:00Z', completed_at: '2026-10-22T09:04:00Z', questions: 5, correct: 4,
      accuracy: 0.8, concepts: ['alg.linear'], findings: [{ concept_id: 'alg.linear', label: 'Linear equations' }],
      href: '/quiz', note: 'The completion stores the XP gate snapshot, not a score.',
    },
  ],
  activity: [{ date: '2026-10-20', questions: 5, correct: 4, concepts: ['alg.linear'] }],
  history: [
    { id: 'answer:1', kind: 'answer', at: '2026-10-20T10:00:00Z', title: 'Linear equations', detail: 'Correct answer', correct: true, href: '/study?concept=alg.linear' },
    { id: 'adjustment:1', kind: 'adjustment', at: '2026-10-20T10:01:00Z', title: 'Performance adjustment', detail: 'Task adjustment', xp: -5 },
  ],
  history_total: 3,
  limitations: ['Reports show current state, not a full transition audit trail.'],
}

const emptyReport: LearnerReport = {
  ...baseReport,
  overview: { total: 3, mastered: 0, in_progress: 0, unlocked: 0, locked: 3, xp_lifetime: 0, xp_period: 0, activity_days: 0, questions: 0, correct: 0 },
  domains: [], mastery: [], assessments: [], activity: [], history: [], history_total: 0,
}

beforeEach(() => {
  vi.mocked(getLearnerReport).mockReset()
  vi.mocked(getLearnerReportHistory).mockReset()
  vi.stubGlobal('print', vi.fn())
})

describe('Reports page', () => {
  it('renders the overview for an established learner', async () => {
    vi.mocked(getLearnerReport).mockResolvedValue(baseReport)
    render(<ReportsPage />)
    expect(await screen.findByRole('heading', { name: 'Reports' })).toBeDefined()
    expect(screen.getAllByText('Mastered').length).toBeGreaterThan(0)
    expect(screen.getByText(/concepts in view/)).toBeDefined()
    expect(screen.getAllByText(/1,?200/).length).toBeGreaterThan(0)
    expect(screen.getAllByText(/algebra/i).length).toBeGreaterThan(0)
    expect(screen.getByRole('button', { name: 'Print report' })).not.toBeDisabled()
  })

  it('switches to mastery and shows gaps and evidence', async () => {
    vi.mocked(getLearnerReport).mockResolvedValue(baseReport)
    render(<ReportsPage />)
    await screen.findByRole('heading', { name: 'Reports' })
    fireEvent.click(screen.getByRole('tab', { name: 'Mastery' }))
    expect((await screen.findAllByText('Quadratic equations')).length).toBeGreaterThan(0)
    expect(screen.getAllByText('Linear equations').length).toBeGreaterThan(0)
    expect(screen.getAllByText(/Linear equations/).length).toBeGreaterThan(1)
  })

  it('switches to assessments and shows honest status', async () => {
    vi.mocked(getLearnerReport).mockResolvedValue(baseReport)
    render(<ReportsPage />)
    await screen.findByRole('heading', { name: 'Reports' })
    fireEvent.click(screen.getByRole('tab', { name: 'Assessments' }))
    expect((await screen.findAllByText('Mastery check')).length).toBeGreaterThan(0)
    expect(screen.getAllByText(/Completed/).length).toBeGreaterThan(0)
    expect(screen.getByText(/not a score/)).toBeDefined()
  })

  it('loads the next history page without replacing earlier events', async () => {
    vi.mocked(getLearnerReport).mockResolvedValue(baseReport)
    vi.mocked(getLearnerReportHistory).mockResolvedValue({
      events: [{ id: 'quiz:1', kind: 'quiz', at: '2026-10-22T09:04:00Z', title: 'Mastery check', xp: 120, href: '/quiz' }],
      total: 2, limit: 25, offset: 1,
    })
    render(<ReportsPage />)
    await screen.findByRole('heading', { name: 'Reports' })
    fireEvent.click(screen.getByRole('tab', { name: 'Learning history' }))
    fireEvent.click(await screen.findByRole('button', { name: /Load more/ }))
    expect((await screen.findAllByText('Mastery check')).length).toBeGreaterThan(0)
    expect(screen.getAllByText('Linear equations').length).toBeGreaterThan(0)
    expect(getLearnerReportHistory).toHaveBeenCalledWith({ offset: 2, limit: 25 })
  })

  it('shows negative XP ledger events in learning history', async () => {
    vi.mocked(getLearnerReport).mockResolvedValue(baseReport)
    render(<ReportsPage />)
    await screen.findByRole('heading', { name: 'Reports' })
    fireEvent.click(screen.getByRole('tab', { name: 'Learning history' }))
    expect(await screen.findByText('-5 XP')).toBeDefined()
  })

  it('describes an empty record instead of inventing one', async () => {
    vi.mocked(getLearnerReport).mockResolvedValue(emptyReport)
    render(<ReportsPage />)
    await screen.findByRole('heading', { name: 'Reports' })
    expect(screen.getByText(/absence of records/)).toBeDefined()
    fireEvent.click(screen.getByRole('tab', { name: 'Mastery' }))
    expect(await screen.findByText('No concepts match this domain.')).toBeDefined()
    fireEvent.click(screen.getByRole('tab', { name: 'Assessments' }))
    expect(screen.getByText('No assessments have been recorded.')).toBeDefined()
  })

  it('shows an error with a retry when the report cannot load', async () => {
    vi.mocked(getLearnerReport).mockRejectedValue(new Error('Reports fetch failed: 401'))
    render(<ReportsPage />)
    const retry = await screen.findByRole('button', { name: 'Try again' })
    expect(retry).toBeDefined()
    vi.mocked(getLearnerReport).mockResolvedValue(baseReport)
    fireEvent.click(retry)
    await waitFor(() => expect(screen.getAllByText('Mastered').length).toBeGreaterThan(0))
  })

  it('marks active filters and tabs with the same square blue treatment', async () => {
    vi.mocked(getLearnerReport).mockResolvedValue(baseReport)
    render(<ReportsPage />)
    await screen.findByRole('heading', { name: 'Reports' })

    const period = screen.getByRole('button', { name: 'Last 30 days' })
    const tab = screen.getByRole('tab', { name: 'Overview' })
    for (const control of [period, tab]) {
      expect(control.className).toContain('bg-mathua-blue')
      expect(control.className).toContain('text-white')
      expect(control.className).toContain('border-mathua-blue')
      expect(control.className).toContain('rounded-none')
    }
  })
})
