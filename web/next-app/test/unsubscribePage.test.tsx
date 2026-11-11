import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import UnsubscribePage from '../app/notifications/unsubscribe/page'
import { unsubscribeNotifications } from '../lib/api'

vi.mock('next/navigation', () => ({
  useSearchParams: () => new URLSearchParams('token=abc123'),
}))

vi.mock('../lib/api', async (importOriginal) => ({
  ...(await importOriginal<typeof import('../lib/api')>()),
  unsubscribeNotifications: vi.fn(),
}))

vi.mock('../components/app-shell', () => ({
  default: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
}))
vi.mock('../components/Footer', () => ({ default: () => <footer /> }))
vi.mock('../components/Loading', () => ({ default: ({ label }: { label: string }) => <p>{label}</p> }))

beforeEach(() => {
  vi.clearAllMocks()
  vi.mocked(unsubscribeNotifications).mockResolvedValue(undefined)
})

describe('the notification unsubscribe page', () => {
  it('sends the signed token and confirms the change', async () => {
    render(<UnsubscribePage />)
    await waitFor(() => expect(unsubscribeNotifications).toHaveBeenCalledWith('abc123'))
    expect(await screen.findByText('Email reminders are off.')).toBeTruthy()
    expect(screen.getByRole('link', { name: /Notification settings/ })).toHaveAttribute('href', '/settings#notifications')
  })

  it('shows the server message when the link is rejected', async () => {
    vi.mocked(unsubscribeNotifications).mockRejectedValue(new Error('Unsubscribe failed, the link may have expired.'))
    render(<UnsubscribePage />)
    expect(await screen.findByText('Unsubscribe failed, the link may have expired.')).toBeTruthy()
  })
})
