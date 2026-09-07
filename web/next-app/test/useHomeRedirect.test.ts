import { describe, it, expect, vi, beforeEach } from 'vitest'
import { renderHook, waitFor } from '@testing-library/react'
import { useHomeRedirect } from '../hooks/useHomeRedirect'

const replaceMock = vi.fn()

vi.mock('next/navigation', () => ({
  useRouter: () => ({ replace: replaceMock, push: vi.fn() }),
}))

describe('useHomeRedirect', () => {
  beforeEach(() => {
    localStorage.clear()
    sessionStorage.clear()
    replaceMock.mockClear()
  })

  it('bounces logged-in visits to /profile', async () => {
    localStorage.setItem('mathua_token', 't')
    const { result } = renderHook(() => useHomeRedirect(true))
    await waitFor(() => expect(replaceMock).toHaveBeenCalledWith('/profile'))
    expect(result.current).toBe(false)
  })

  it('leaves guests on the landing page', async () => {
    const { result } = renderHook(() => useHomeRedirect(true))
    await waitFor(() => expect(result.current).toBe(true))
    expect(replaceMock).not.toHaveBeenCalled()
  })

  it('waits for mount before deciding', () => {
    localStorage.setItem('mathua_token', 't')
    const { result } = renderHook(() => useHomeRedirect(false))
    expect(result.current).toBe(false)
    expect(replaceMock).not.toHaveBeenCalled()
  })
})
