import { afterEach, describe, expect, it, vi, beforeEach } from 'vitest'
import { useState } from 'react'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import AdminAnalyticsDashboard from '../components/admin/analytics'
import { adminAnalytics, adminAnalyticsLive, getAdminNotificationHealth } from '../lib/api'
import { DateRangePicker } from '../components/ui/date-range-picker'
import type { DateRange } from 'react-day-picker'

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
  window.history.replaceState({}, '', '/')
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

  it('restores a custom window from the URL', async () => {
    window.history.replaceState({}, '', '/?range=custom&from=2026-11-01&to=2026-11-10')
    render(<AdminAnalyticsDashboard />)
    await screen.findByText('Unique learners')
    expect(adminAnalytics).toHaveBeenCalledWith(
      'custom',
      expect.objectContaining({ from: '2026-11-01', to: '2026-11-10' }),
    )
  })

  it('drops the custom days when a preset is chosen', async () => {
    window.history.replaceState({}, '', '/?range=custom&from=2026-11-01&to=2026-11-10')
    render(<AdminAnalyticsDashboard />)
    await screen.findByText('Unique learners')
    fireEvent.click(screen.getByRole('radio', { name: '24h' }))
    await waitFor(() => {
      expect(adminAnalytics).toHaveBeenLastCalledWith('24h', expect.anything())
    })
    const last = vi.mocked(adminAnalytics).mock.calls.at(-1)
    expect(last?.[1]).not.toHaveProperty('from')
    expect(window.location.search).toBe('?range=24h')
  })

  it('loads a custom window picked from the calendar and shares it in the URL', async () => {
    render(<AdminAnalyticsDashboard />)
    await screen.findByText('Unique learners')
    fireEvent.click(screen.getByRole('button', { name: 'Select a date range' }))
    fireEvent.click(screen.getByRole('button', { name: /previous month/i }))
    fireEvent.click(screen.getAllByText('20')[0].closest('button') as HTMLButtonElement)
    fireEvent.click(screen.getAllByText('22')[0].closest('button') as HTMLButtonElement)

    await waitFor(() => {
      expect(adminAnalytics).toHaveBeenLastCalledWith('custom', expect.objectContaining({ from: expect.stringMatching(/-20$/), to: expect.stringMatching(/-22$/) }))
    })
    const params = new URLSearchParams(window.location.search)
    expect(params.get('range')).toBe('custom')
    expect(params.get('from')).toMatch(/-20$/)
    expect(params.get('to')).toMatch(/-22$/)
  })

  it('shows notification delivery health on the operations tab', async () => {
    render(<AdminAnalyticsDashboard />)
    await screen.findByText('Unique learners')
    fireEvent.mouseDown(screen.getByRole('tab', { name: 'Operations' }))
    await screen.findByText('Notification delivery')
    expect(await screen.findByText('No failed or suppressed deliveries.')).toBeTruthy()
  })
})

describe('date range picker', () => {
  afterEach(() => {
    window.history.replaceState({}, '', '/')
  })

  it('reports picked days as UTC midnights without a timezone shift', () => {
    const today = new Date()
    const min = new Date(Date.UTC(today.getUTCFullYear(), today.getUTCMonth() - 3, 1))
    const onChange = vi.fn()
    function Harness() {
      const [value, setValue] = useState<DateRange | undefined>(undefined)
      return (
        <DateRangePicker
          value={value}
          onChange={(next) => {
            setValue(next)
            onChange(next)
          }}
          min={min}
          max={today}
        />
      )
    }
    render(<Harness />)
    fireEvent.click(screen.getByRole('button', { name: 'Select a date range' }))
    fireEvent.click(screen.getByRole('button', { name: /previous month/i }))

    const month = new Date(Date.UTC(today.getUTCFullYear(), today.getUTCMonth() - 1, 1))
    const days = screen.getAllByText('10')
    fireEvent.click(days[0].closest('button') as HTMLButtonElement)
    fireEvent.click(screen.getAllByText('12')[0].closest('button') as HTMLButtonElement)

    const picked = onChange.mock.calls.at(-1)?.[0]
    expect(picked?.from.toISOString()).toBe(`${month.getUTCFullYear()}-${String(month.getUTCMonth() + 1).padStart(2, '0')}-10T00:00:00.000Z`)
    expect(picked?.to.toISOString()).toBe(`${month.getUTCFullYear()}-${String(month.getUTCMonth() + 1).padStart(2, '0')}-12T00:00:00.000Z`)
  })
})
