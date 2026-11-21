import { act, render, waitFor } from '@testing-library/react'
import { beforeEach, describe, expect, it, vi } from 'vitest'

const state = vi.hoisted(() => {
  type Options = {
    id: string
    description: string
    action: { label: string; onClick: () => void }
    cancel: { label: string; onClick: () => void }
  }
  const call = vi.fn((_title: string, _options: Options) => 'toast-1')
  return {
    toast: Object.assign(call, { dismiss: vi.fn() }),
    pathname: '/learn',
    push: vi.fn(),
  }
})

vi.mock('sonner', () => ({ toast: state.toast }))

vi.mock('next/navigation', () => ({
  usePathname: () => state.pathname,
  useRouter: () => ({ push: state.push }),
}))

const api = vi.hoisted(() => ({
  listAnnouncements: vi.fn(),
  listNotifications: vi.fn(),
  listMyReports: vi.fn(),
  markAnnouncementRead: vi.fn(),
  markNotificationRead: vi.fn(),
  markMyReportsRead: vi.fn(),
}))

vi.mock('@/lib/api', () => api)
vi.mock('@/lib/auth', () => ({ getToken: () => 'token' }))

import NotificationToasts from '../components/notifications/NotificationToasts'

const announcement = (over: Record<string, unknown> = {}) => ({
  id: 'a1',
  title: '148 figures are live',
  body: 'Drag the probe and the numbers follow.\nSecond line.',
  published: true,
  created_at: '2026-11-20T10:00:00Z',
  published_at: '2026-11-20T10:00:00Z',
  read: false,
  ...over,
})

const notification = (over: Record<string, unknown> = {}) => ({
  id: 'n1',
  kind: 'review',
  title: 'Reviews are waiting',
  body: 'You have 4 due reviews.',
  href: '/review',
  created_at: '2026-11-21T10:00:00Z',
  read: false,
  ...over,
})

async function settle() {
  await act(async () => {
    await new Promise((resolve) => setTimeout(resolve, 0))
  })
}

beforeEach(() => {
  state.pathname = '/learn'
  state.toast.mockClear()
  state.toast.dismiss.mockClear()
  state.push.mockClear()
  for (const fn of Object.values(api)) fn.mockReset()
  api.listAnnouncements.mockResolvedValue({ announcements: [], unread: 0 })
  api.listNotifications.mockResolvedValue({
    notifications: [],
    unread: 0,
    has_more: false,
    next_cursor: '',
  })
  api.listMyReports.mockResolvedValue({ reports: [], unread: 0 })
  api.markAnnouncementRead.mockResolvedValue(undefined)
  api.markNotificationRead.mockResolvedValue(0)
  api.markMyReportsRead.mockResolvedValue(undefined)
})

describe('NotificationToasts', () => {
  it('toasts the newest unread item through the shadcn toaster', async () => {
    api.listAnnouncements.mockResolvedValue({ announcements: [announcement()], unread: 1 })
    api.listNotifications.mockResolvedValue({
      notifications: [notification()],
      unread: 1,
      has_more: false,
      next_cursor: '',
    })

    render(<NotificationToasts />)

    await waitFor(() => expect(state.toast).toHaveBeenCalledTimes(1))
    const [title, options] = state.toast.mock.calls[0]
    expect(title).toBe('Reviews are waiting')
    expect(options.id).toBe('mathua-notification:n1')
    expect(options.description).toBe('Review reminder · You have 4 due reviews.')
    expect(options.action.label).toBe('Read')
    expect(options.cancel.label).toBe('Dismiss')
  })

  it('records the item and opens its page from the read action', async () => {
    api.listNotifications.mockResolvedValue({
      notifications: [notification()],
      unread: 1,
      has_more: false,
      next_cursor: '',
    })

    render(<NotificationToasts />)
    await waitFor(() => expect(state.toast).toHaveBeenCalledTimes(1))

    await act(async () => {
      state.toast.mock.calls[0][1].action.onClick()
    })

    expect(api.markNotificationRead).toHaveBeenCalledWith('n1')
    expect(state.toast.dismiss).toHaveBeenCalledWith('mathua-notification:n1')
    expect(state.push).toHaveBeenCalledWith('/review')
  })

  it('dismissing is not reading', async () => {
    api.listAnnouncements.mockResolvedValue({ announcements: [announcement()], unread: 1 })

    render(<NotificationToasts />)
    await waitFor(() => expect(state.toast).toHaveBeenCalledTimes(1))

    state.toast.mock.calls[0][1].cancel.onClick()
    await settle()

    expect(api.markAnnouncementRead).not.toHaveBeenCalled()
    expect(api.markNotificationRead).not.toHaveBeenCalled()
    expect(api.markMyReportsRead).not.toHaveBeenCalled()
    expect(state.push).not.toHaveBeenCalled()
  })

  it('says nothing during a graded session', async () => {
    state.pathname = '/session'
    api.listNotifications.mockResolvedValue({
      notifications: [notification()],
      unread: 1,
      has_more: false,
      next_cursor: '',
    })

    render(<NotificationToasts />)
    await settle()

    expect(state.toast).not.toHaveBeenCalled()
  })

  it('says nothing when every item is already read', async () => {
    render(<NotificationToasts />)
    await settle()
    expect(state.toast).not.toHaveBeenCalled()
  })
})
