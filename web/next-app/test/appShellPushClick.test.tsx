import { beforeEach, describe, expect, it, vi } from 'vitest'
import { render, waitFor } from '@testing-library/react'

const reportNotificationClick = vi.fn()

vi.mock('../lib/api', () => ({
  reportNotificationClick: (...args: unknown[]) => reportNotificationClick(...args),
  getNotificationUnreadCount: () => Promise.resolve(0),
  listAnnouncements: () => Promise.resolve([]),
  listMyReports: () => Promise.resolve([]),
  getDueReviews: () => Promise.resolve({ due: 0 }),
  getSettings: () => Promise.resolve({}),
  adminOverview: () => Promise.resolve({ open_reports: 0 }),
  getConfig: () => Promise.resolve({}),
}))

vi.mock('../hooks/useAuthState', () => ({
  useAuthState: () => ({ loggedIn: auth.loggedIn, user: null }),
}))

vi.mock('../lib/dicebear', () => ({ resolveAvatar: () => ({}) }))
vi.mock('../components/app-sidebar', () => ({ AppSidebar: () => null }))
vi.mock('../components/BottomTabs', () => ({ default: () => null }))
vi.mock('../components/icons', () => ({ Bell: () => null, ChevronsRight: () => null }))
vi.mock('../components/ui/sidebar', () => ({
  SidebarProvider: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
  SidebarInset: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
  SidebarTrigger: () => null,
}))
vi.mock('next/link', () => ({ default: ({ children }: { children: React.ReactNode }) => <a>{children}</a> }))

import AppShell from '../components/app-shell'

const auth = { loggedIn: true }

describe('push click reporting', () => {
  beforeEach(() => {
    reportNotificationClick.mockReset().mockResolvedValue(0)
    auth.loggedIn = true
    sessionStorage.clear()
    window.history.replaceState({}, '', '/review?ntf=ntf_abc')
  })

  it('reports the click and strips the parameter once the report settles', async () => {
    const changed = vi.fn()
    window.addEventListener('notifications-changed', changed)
    render(<AppShell><div /></AppShell>)

    await waitFor(() => expect(reportNotificationClick).toHaveBeenCalledWith('ntf_abc'))
    await waitFor(() => expect(window.location.search).toBe(''))
    expect(changed).toHaveBeenCalled()
    window.removeEventListener('notifications-changed', changed)
  })

  it('keeps a signed-out click for after sign-in', async () => {
    auth.loggedIn = false
    render(<AppShell><div /></AppShell>)

    await waitFor(() => expect(sessionStorage.getItem('mathua_pending_ntf')).toBe('ntf_abc'))
    expect(reportNotificationClick).not.toHaveBeenCalled()
    expect(window.location.search).toBe('?ntf=ntf_abc')
  })

  it('reports a stored click when the learner signs in', async () => {
    window.history.replaceState({}, '', '/learn')
    sessionStorage.setItem('mathua_pending_ntf', 'ntf_stored')
    render(<AppShell><div /></AppShell>)

    await waitFor(() => expect(reportNotificationClick).toHaveBeenCalledWith('ntf_stored'))
    expect(sessionStorage.getItem('mathua_pending_ntf')).toBeNull()
  })
})
