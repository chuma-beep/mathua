import { fireEvent, render, screen, waitFor } from '@testing-library/react'
import { beforeEach, describe, expect, it, vi } from 'vitest'

const state = vi.hoisted(() => {
  const call = vi.fn()
  return { toast: Object.assign(call, { custom: vi.fn(), dismiss: vi.fn() }) }
})

vi.mock('sonner', () => ({ toast: state.toast }))

const api = vi.hoisted(() => ({ markAnnouncementRead: vi.fn() }))
vi.mock('@/lib/api', () => ({ markAnnouncementRead: api.markAnnouncementRead }))

import AnnouncementToast from '../components/announcements/AnnouncementToast'
import type { Announcement } from '../lib/api'

const item = (over: Partial<Announcement> = {}): Announcement => ({
  id: 'a1',
  title: '148 figures are live',
  body: 'Drag the probe and the numbers follow.\nSecond line.',
  published: true,
  created_at: '2026-11-20T10:00:00Z',
  published_at: '2026-11-20T10:00:00Z',
  read: false,
  ...over,
})

const openToast = async () => {
  render(<AnnouncementToast announcement={item()} />)
  await waitFor(() => expect(state.toast.custom).toHaveBeenCalledTimes(1))
  render(<div data-sonner-toast="true">{state.toast.custom.mock.calls[0][0]()}</div>)
}

beforeEach(() => {
  window.history.replaceState({}, '', '/learn')
  state.toast.mockClear()
  state.toast.custom.mockClear()
  state.toast.dismiss.mockClear()
  api.markAnnouncementRead.mockReset()
  api.markAnnouncementRead.mockResolvedValue(undefined)
})

describe('AnnouncementToast', () => {
  it('shows the announcement once, in the app style', async () => {
    await openToast()
    expect(screen.getByText('148 figures are live')).toBeTruthy()
    expect(state.toast.custom).toHaveBeenCalledTimes(1)
    expect(state.toast.custom.mock.calls[0][1].id).toBe('announcement-a1')
  })

  it('says nothing when the announcement is already read', () => {
    render(<AnnouncementToast announcement={item({ read: true })} />)
    expect(state.toast.custom).not.toHaveBeenCalled()
  })

  it('stays silent during a graded session', () => {
    window.history.replaceState({}, '', '/session')
    render(<AnnouncementToast announcement={item()} />)
    expect(state.toast.custom).not.toHaveBeenCalled()
  })

  it('dismissing is not reading', async () => {
    await openToast()
    fireEvent.click(screen.getByRole('button', { name: 'Dismiss' }))
    expect(state.toast.dismiss).toHaveBeenCalledWith('announcement-a1')
    expect(api.markAnnouncementRead).not.toHaveBeenCalled()
  })

  it('opening the announcement records it', async () => {
    await openToast()
    fireEvent.click(screen.getByRole('button', { name: 'Read' }))
    expect(api.markAnnouncementRead).toHaveBeenCalledWith('a1')
    expect(state.toast.dismiss).toHaveBeenCalledWith('announcement-a1')
  })
})
