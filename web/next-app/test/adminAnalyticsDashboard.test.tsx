import { describe, expect, it, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import AdminAnalyticsDashboard from '../components/admin/analytics'
import { adminAnalytics, adminAnalyticsLive, getAdminNotificationHealth } from '../lib/api'

vi.mock('../lib/api', async (importOriginal) => ({
  ...(await importOriginal<typeof import('../lib/api')>()),
  adminAnalytics: vi.fn(),
  adminAnalyticsLive: vi.fn(),
  getAdminNotificationHealth: vi.fn(),
}))

const analytics = {
  range: '7d',
  from: '2026-11-08T00:00:00Z',
  to: '2026-11-15T12:00:00Z',
  bucket_seconds: 86400,
  generated_at: '2026-11-15T12:00:01Z',
  partial_last_bucket: true,
  accounts: { total: 4, registered: 2, guest: 1, unclassified: 1 },
  totals: {
    unique_learners: 12,
    unique_registered: 8,
    unique_guests: 3,
    unique_unclassified: 1,
    answers: 30,
    correct: 20,
    new_registrations: 2,
    diagnostic_cohort_registered: 2,
    diagnostic_cohort_completed: 1,
    quiz_completions: 3,
    timed_quiz_completions: 2,
    xp_awarded: 100,
    xp_deducted: 10,
    diagnostic_completions: 1,
    reports_opened: 4,
    reports_resolved: 2,
  },
  activity: [
    { start: 1, registered: 4, guests: 2, unclassified: 1, answers: 12, correct: 8 },
    { start: 2, registered: 6, guests: 1, unclassified: 0, answers: 18, correct: 12 },
  ],
  signups: [{ start: 1, value: 2 }],
  mastery: [{ start: 1, value: 3 }],
  xp: [{ start: 1, awarded: 100, deducted: 10 }],
  reports: [{ start: 1, opened: 4, resolved: 2 }],
  quizzes: [{ start: 1, value: 3 }],
  timed_quizzes: [{ start: 1, value: 2 }],
  diagnostic_completions: [{ start: 1, value: 1 }],
  top_concepts: [{ concept_id: 'arith.add', attempts: 10, correct: 7, reports: 0 }],
  most_reported: [{ concept_id: 'arith.add', attempts: 10, correct: 7, reports: 2 }],
  mastered_by_concept: { 'arith.add': 3 },
  mastered_current: 4,
  report_status: { open: 2, resolved: 2 },
}

beforeEach(() => {
  vi.clearAllMocks()
  vi.mocked(adminAnalytics).mockResolvedValue(analytics)
  vi.mocked(adminAnalyticsLive).mockResolvedValue({
    online_now: 3,
    online_window_seconds: 300,
    practicing_now: 5,
    guest_practice_now: 1,
    answers_15m: 7,
    active_sessions: 4,
    generated_at: '2026-11-15T12:00:01Z',
  })
  vi.mocked(getAdminNotificationHealth).mockResolvedValue({
    stats: [{ channel: 'in_app', status: 'sent', count: 20 }],
    problems: [],
    engagement: [],
    suppressions: [],
  })
})

describe('admin analytics dashboard', () => {
  it('renders range-driven KPIs and the live strip', async () => {
    render(<AdminAnalyticsDashboard />)
    await screen.findByText('Unique learners')
    expect(screen.getByText('12')).toBeTruthy()
    expect(screen.getByText('30')).toBeTruthy()
    expect(screen.getByText('66.7%')).toBeTruthy()
    expect(screen.getByText('3')).toBeTruthy()
    expect(adminAnalytics).toHaveBeenCalledWith('7d', expect.objectContaining({ sections: expect.any(Array) }))
  })

  it('renders the growth and learning charts', async () => {
    render(<AdminAnalyticsDashboard />)
    await screen.findByText('User growth and participation')
    expect(screen.getAllByText('Accuracy').length).toBeGreaterThan(0)
    expect(screen.getByText('Mastery events')).toBeTruthy()
    expect(screen.getAllByText('New registrations').length).toBeGreaterThan(0)
  })

  it('loads the sections for the selected tab', async () => {
    render(<AdminAnalyticsDashboard />)
    await screen.findByText('Unique learners')
    fireEvent.mouseDown(screen.getByRole('tab', { name: 'Engagement' }))
    await screen.findByRole('heading', { name: 'New registrations' })
    await waitFor(() => {
      const calls = vi.mocked(adminAnalytics).mock.calls
      expect(calls[calls.length - 1][1].sections).toContain('signups')
    })
  })

  it('changes the range and reloads the window', async () => {
    render(<AdminAnalyticsDashboard />)
    await screen.findByText('Unique learners')
    fireEvent.click(screen.getByRole('radio', { name: '24h' }))
    await waitFor(() => {
      expect(adminAnalytics).toHaveBeenLastCalledWith('24h', expect.anything())
    })
  })

  it('shows notification delivery health on the operations tab', async () => {
    render(<AdminAnalyticsDashboard />)
    await screen.findByText('Unique learners')
    fireEvent.mouseDown(screen.getByRole('tab', { name: 'Operations' }))
    await screen.findByText('Notification delivery')
    expect(await screen.findByText('No failed or suppressed deliveries.')).toBeTruthy()
  })
})
