import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import ReportUpdatesBadge from '../components/ReportUpdatesBadge'
import { listAnnouncements, listMyReports } from '../lib/api'

vi.mock('../lib/api', async (importOriginal) => ({
  ...(await importOriginal<typeof import('../lib/api')>()),
  listMyReports: vi.fn(),
  listAnnouncements: vi.fn(),
}))

beforeEach(() => {
  vi.clearAllMocks()
  vi.mocked(listAnnouncements).mockResolvedValue({ unread: 0, announcements: [] })
})

describe('the report update badge', () => {
  it('shows the unread verdict count beside the avatar', async () => {
    vi.mocked(listMyReports).mockResolvedValue({ unread: 2, reports: [] })
    render(<ReportUpdatesBadge enabled />)

    const link = await screen.findByRole('link', { name: '2 unread notifications' })
    expect(link).toHaveAttribute('href', '/my-reports')
    expect(link).toHaveTextContent('2')
  })

  it('clears when the history page marks the verdicts read', async () => {
    vi.mocked(listMyReports).mockResolvedValueOnce({ unread: 1, reports: [] })
    render(<ReportUpdatesBadge enabled />)
    expect(await screen.findByRole('link', { name: '1 unread notification' })).toBeTruthy()

    vi.mocked(listMyReports).mockResolvedValueOnce({ unread: 0, reports: [] })
    window.dispatchEvent(new Event('report-updates-changed'))

    await waitFor(() => expect(screen.queryByRole('link', { name: /notification/ })).toBeNull())
  })

  it('does not fetch for a signed-out visitor', async () => {
    render(<ReportUpdatesBadge enabled={false} />)
    await Promise.resolve()
    expect(listMyReports).not.toHaveBeenCalled()
  })
})
