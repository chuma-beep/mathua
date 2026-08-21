import { describe, it, expect, beforeEach } from 'vitest'
import { renderHook, act } from '@testing-library/react'
import { useAuthState } from '../hooks/useAuthState'

describe('useAuthState', () => {
  beforeEach(() => {
    localStorage.clear()
  })

  it('reports logged-out by default', () => {
    const { result } = renderHook(() => useAuthState())
    expect(result.current.loggedIn).toBe(false)
    expect(result.current.user).toBe(null)
  })

  it('becomes reactive when auth-changed fires after login', async () => {
    const { result } = renderHook(() => useAuthState())
    act(() => {
      localStorage.setItem('mathua_token', 'tok')
      localStorage.setItem('mathua_user', JSON.stringify({ student_id: 's1', name: 'A', username: 'a', concepts_mastered: 0, current_streak: 0, level: 'l', diagnostic_completed: false }))
      window.dispatchEvent(new Event('auth-changed'))
    })
    expect(result.current.loggedIn).toBe(true)
    expect(result.current.user?.student_id).toBe('s1')
  })

  it('returns to logged-out when auth-changed fires after logout', async () => {
    localStorage.setItem('mathua_token', 'tok')
    const { result } = renderHook(() => useAuthState())
    act(() => {
      localStorage.removeItem('mathua_token')
      localStorage.removeItem('mathua_user')
      window.dispatchEvent(new Event('auth-changed'))
    })
    expect(result.current.loggedIn).toBe(false)
    expect(result.current.user).toBe(null)
  })
})
