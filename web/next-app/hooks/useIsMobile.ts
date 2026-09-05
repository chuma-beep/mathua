'use client'

import { useState, useEffect } from 'react'

/**
 * SSR-safe mobile flag. Uses matchMedia so orientation/zoom are handled,
 * falls back to innerWidth for older browsers. Defaults to false on server
 * (desktop-first render) — callers should tolerate one-frame correction.
 */
export function useIsMobile(breakpoint = 1024): boolean {
  const query = `(max-width: ${breakpoint - 1}px)`
  const [isMobile, setIsMobile] = useState<boolean>(() => {
    if (typeof window === 'undefined' || typeof window.matchMedia !== 'function') return false
    return window.matchMedia(query).matches
  })

  useEffect(() => {
    if (typeof window === 'undefined') return
    if (typeof window.matchMedia === 'function') {
      const mql = window.matchMedia(query)
      const onChange = (e: MediaQueryListEvent) => setIsMobile(e.matches)
      setIsMobile(mql.matches)
      mql.addEventListener('change', onChange)
      return () => mql.removeEventListener('change', onChange)
    }
    // Fallback: resize listener with rAF debounce
    let raf = 0
    const check = () => {
      cancelAnimationFrame(raf)
      raf = requestAnimationFrame(() => setIsMobile(window.innerWidth < breakpoint))
    }
    check()
    window.addEventListener('resize', check)
    return () => {
      cancelAnimationFrame(raf)
      window.removeEventListener('resize', check)
    }
  }, [query, breakpoint])

  return isMobile
}
