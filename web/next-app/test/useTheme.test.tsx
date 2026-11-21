import { beforeEach, describe, expect, it } from 'vitest'
import { act, renderHook } from '@testing-library/react'
import { useTheme } from '../hooks/useTheme'

// The globe and the device gallery are canvases: they read the theme once and
// only redraw when their hook sees a new value. A toggle therefore has to reach
// every hook instance, not just the component that owns the button, or the
// canvases keep the old palette until a reload.
describe('useTheme', () => {
  beforeEach(() => {
    localStorage.clear()
    document.documentElement.classList.remove('dark')
  })

  it('tells every consumer when one of them toggles', () => {
    const header = renderHook(() => useTheme())
    const globe = renderHook(() => useTheme())
    const before = globe.result.current.theme

    act(() => {
      header.result.current.toggleTheme()
    })

    expect(header.result.current.theme).not.toBe(before)
    expect(globe.result.current.theme).toBe(header.result.current.theme)
    expect(document.documentElement.classList.contains('dark')).toBe(
      header.result.current.theme === 'dark',
    )
  })

  it('follows a change made in another tab', () => {
    const globe = renderHook(() => useTheme())
    const before = globe.result.current.theme

    act(() => {
      localStorage.setItem('mathua-theme', before === 'dark' ? 'light' : 'dark')
      window.dispatchEvent(
        new StorageEvent('storage', {
          key: 'mathua-theme',
          newValue: before === 'dark' ? 'light' : 'dark',
        }),
      )
    })

    expect(globe.result.current.theme).not.toBe(before)
  })
})
