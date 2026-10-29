import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor, fireEvent } from '@testing-library/react'
import MyReportsPage from '../app/my-reports/page'
import { clearAnnouncements, clearMyReports, listAnnouncements, listMyReports, markAnnouncementsRead, markMyReportsRead } from '../lib/api'
import type { QuestionReport } from '../lib/api'

vi.mock('../lib/api', async (importOriginal) => ({
  ...(await importOriginal<typeof import('../lib/api')>()),
  listMyReports: vi.fn(),
  markMyReportsRead: vi.fn(),
  clearMyReports: vi.fn(),
  listAnnouncements: vi.fn(),
  markAnnouncementsRead: vi.fn(),
  clearAnnouncements: vi.fn(),
}))

vi.mock('../hooks/useAuthState', () => ({
  useAuthState: () => ({
    loggedIn: true,
    user: { student_id: 'stu_1', name: 'Ada Example', username: 'ada' },
  }),
}))

vi.mock('../components/app-shell', () => ({
  default: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
}))
vi.mock('../components/Footer', () => ({ default: () => <footer /> }))
vi.mock('../components/KatexContent', () => ({
  default: ({ children }: { children: React.ReactNode }) => <span>{children}</span>,
}))

function report(over: Partial<QuestionReport> & { id: number }): QuestionReport {
  return {
    reporter_id: 'stu_1',
    concept_id: 'algebra.linear',
    kind: 'question',
    question: '2x = 4',
    expected: '2',
    explanation: '',
    lesson_id: '',
    source: 'practice',
    session_id: '',
    attempt_id: '',
    reason: 'wrong_answer',
    detail: '',
    status: 'open',
    resolution: '',
    resolved_by: '',
    resolved_at: '',
    created_at: '2026-10-20T00:00:00Z',
    ...over,
  }
}

beforeEach(() => {
  vi.clearAllMocks()
  vi.mocked(markMyReportsRead).mockResolvedValue(undefined)
  vi.mocked(clearMyReports).mockResolvedValue(undefined)
  vi.mocked(listAnnouncements).mockResolvedValue({ unread: 0, announcements: [] })
  vi.mocked(markAnnouncementsRead).mockResolvedValue(undefined)
  vi.mocked(clearAnnouncements).mockResolvedValue(undefined)
})

describe('the learner report history', () => {
  it('shows a verdict and the moderator reason', async () => {
    vi.mocked(listMyReports).mockResolvedValue({
      unread: 1,
      reports: [
        report({
          id: 7,
          status: 'resolved',
          resolution: 'The expected answer was corrected.',
          resolved_at: '2026-10-25T12:00:00Z',
          unread: true,
        }),
      ],
    })

    render(<MyReportsPage />)

    expect(screen.getByRole('heading', { name: 'Notifications' })).toBeTruthy()
    expect(await screen.findByText('Resolved')).toBeTruthy()
    expect(screen.getByText('The expected answer was corrected.')).toBeTruthy()
    expect(screen.getByText(/1 report has a new verdict/)).toBeTruthy()
    expect(markMyReportsRead).toHaveBeenCalledTimes(1)
  })

  it('marks the history read once the page opens it', async () => {
    vi.mocked(listMyReports).mockResolvedValue({
      unread: 2,
      reports: [report({ id: 1, unread: true }), report({ id: 2, unread: true })],
    })
    render(<MyReportsPage />)

    await waitFor(() => expect(markMyReportsRead).toHaveBeenCalledTimes(1))
    expect(screen.getByText(/2 reports have new verdicts/)).toBeTruthy()
  })

  it('shows announcements above the report history', async () => {
    vi.mocked(listMyReports).mockResolvedValue({ unread: 0, reports: [] })
    vi.mocked(listAnnouncements).mockResolvedValue({
      unread: 1,
      announcements: [
        {
          id: 'ann-1',
          title: 'Timed quizzes are live',
          body: 'Try the new **timed quiz**.',
          published: true,
          created_at: '2026-10-28T12:00:00Z',
          published_at: '2026-10-28T12:00:00Z',
          read: false,
        },
      ],
    })

    render(<MyReportsPage />)

    expect(await screen.findByRole('heading', { name: 'Announcements' })).toBeTruthy()
    expect(screen.getByText('Timed quizzes are live')).toBeTruthy()
    expect(screen.getByText(/Try the new/)).toBeTruthy()
    expect(screen.getByText(/1 new announcement/)).toBeTruthy()
    expect(markAnnouncementsRead).toHaveBeenCalledTimes(1)
  })

  it('still shows the history when the read marker fails', async () => {
    vi.mocked(listMyReports).mockResolvedValue({
      unread: 1,
      reports: [
        report({
          id: 9,
          status: 'resolved',
          resolution: 'Fixed.',
          resolved_at: '2026-10-25T12:00:00Z',
          unread: true,
        }),
      ],
    })
    vi.mocked(markMyReportsRead).mockRejectedValue(new Error('marker failed'))

    render(<MyReportsPage />)

    expect(await screen.findByText('Resolved')).toBeTruthy()
    expect(screen.getByText('Fixed.')).toBeTruthy()
    await waitFor(() => expect(markMyReportsRead).toHaveBeenCalledTimes(1))
    expect(screen.queryByText(/marker failed/)).toBeNull()
  })

  it('clears reports and announcements together', async () => {
    vi.mocked(listMyReports).mockResolvedValue({
      unread: 0,
      reports: [report({ id: 11, question: 'a report to clear' })],
    })
    vi.mocked(listAnnouncements).mockResolvedValue({
      unread: 0,
      announcements: [
        {
          id: 'ann-2',
          title: 'An announcement to clear',
          body: 'This leaves the feed.',
          published: true,
          created_at: '2026-10-28T12:00:00Z',
          published_at: '2026-10-28T12:00:00Z',
          read: true,
        },
      ],
    })
    render(<MyReportsPage />)

    fireEvent.click(await screen.findByRole('button', { name: 'Clear notifications' }))
    expect(screen.getByText(/Clear your notifications/)).toBeTruthy()
    fireEvent.click(screen.getByRole('button', { name: 'Yes, clear notifications' }))

    await waitFor(() => expect(clearMyReports).toHaveBeenCalledTimes(1))
    expect(clearAnnouncements).toHaveBeenCalledTimes(1)
    await waitFor(() => expect(screen.queryByText('An announcement to clear')).toBeNull())
    expect(await screen.findByText(/You have not reported a problem yet/)).toBeTruthy()
  })

  it('offers to clear when only announcements exist', async () => {
    vi.mocked(listMyReports).mockResolvedValue({ unread: 0, reports: [] })
    vi.mocked(listAnnouncements).mockResolvedValue({
      unread: 0,
      announcements: [
        {
          id: 'ann-3',
          title: 'Only an announcement',
          body: 'No reports here.',
          published: true,
          created_at: '2026-10-28T12:00:00Z',
          published_at: '2026-10-28T12:00:00Z',
          read: true,
        },
      ],
    })
    render(<MyReportsPage />)

    expect(await screen.findByRole('button', { name: 'Clear notifications' })).toBeTruthy()
  })

  it('explains the empty state without pretending there is history', async () => {
    vi.mocked(listMyReports).mockResolvedValue({ unread: 0, reports: [] })
    render(<MyReportsPage />)

    expect(await screen.findByText(/You have not reported a problem yet/)).toBeTruthy()
    expect(markMyReportsRead).not.toHaveBeenCalled()
  })
})
