import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor, fireEvent } from '@testing-library/react'
import NotificationsArchivePage from '../app/notifications/page'
import { listNotifications, markAllNotificationsRead, reportNotificationOpen } from '../lib/api'
import type { AppNotification } from '../lib/api'

vi.mock('../lib/api', async (importOriginal) => ({
  ...(await importOriginal<typeof import('../lib/api')>()),
  listNotifications: vi.fn(),
  markAllNotificationsRead: vi.fn(),
  reportNotificationOpen: vi.fn(),
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

function notification(over: Partial<AppNotification> & { id: string }): AppNotification {
  return {
    kind: 'review',
    title: 'You have 3 concepts due for review',
    body: 'A short session keeps them fresh.',
    href: '/review',
    created_at: '2026-11-10T18:00:00Z',
    read: false,
    ...over,
  }
}

beforeEach(() => {
  vi.clearAllMocks()
  vi.mocked(listNotifications).mockResolvedValue({
    notifications: [notification({ id: 'ntf_1' })],
    unread: 1,
    has_more: false,
    next_cursor: '',
  })
  vi.mocked(markAllNotificationsRead).mockResolvedValue(undefined)
  vi.mocked(reportNotificationOpen).mockResolvedValue(0)
})

describe('the notification history page', () => {
  it('lists a notification with its kind and date', async () => {
    render(<NotificationsArchivePage />)

    expect(await screen.findByText('You have 3 concepts due for review')).toBeTruthy()
    expect(screen.getByText('A short session keeps them fresh.')).toBeTruthy()
    expect(screen.getByText('Review')).toBeTruthy()
    expect(listNotifications).toHaveBeenCalledWith(
      expect.objectContaining({ limit: 20, kind: '', unreadOnly: false }),
    )
  })

  it('filters the archive by notification kind', async () => {
    render(<NotificationsArchivePage />)
    await screen.findByText('You have 3 concepts due for review')

    fireEvent.click(screen.getByRole('button', { name: 'Milestones' }))

    await waitFor(() =>
      expect(listNotifications).toHaveBeenLastCalledWith(
        expect.objectContaining({ kind: 'milestone' }),
      ),
    )
  })

  it('records an open and clears the item from the unread view', async () => {
    render(<NotificationsArchivePage />)
    await screen.findByText('You have 3 concepts due for review')

    fireEvent.click(screen.getByRole('button', { name: /Unread/ }))
    await waitFor(() =>
      expect(listNotifications).toHaveBeenLastCalledWith(
        expect.objectContaining({ unreadOnly: true }),
      ),
    )

    fireEvent.click(screen.getByRole('link', { name: 'You have 3 concepts due for review' }))

    await waitFor(() => expect(reportNotificationOpen).toHaveBeenCalledWith('ntf_1'))
    await waitFor(() => expect(screen.queryByText('You have 3 concepts due for review')).toBeNull())
  })

  it('marks every notification read', async () => {
    render(<NotificationsArchivePage />)
    await screen.findByText('You have 3 concepts due for review')

    fireEvent.click(screen.getByRole('button', { name: 'Mark all read' }))

    await waitFor(() => expect(markAllNotificationsRead).toHaveBeenCalled())
    expect(screen.queryByLabelText('Unread notification')).toBeNull()
  })
})
