'use client'

import { useState, useEffect, useCallback } from 'react'

type Theme = 'dark' | 'light'

const STORAGE_KEY = 'mathua-theme'

// A toggle sets the state in the component that owns the button, and that is
// the only copy React knows about. Everything else that reads the theme (the
// globe, the device gallery) would otherwise keep its old palette until a
// reload, so a toggle broadcasts the new value and every hook instance adopts
// it. The storage listener covers a second tab.
const THEME_EVENT = 'mathua-theme-change'

function isTheme(value: unknown): value is Theme {
  return value === 'dark' || value === 'light'
}

function getSystemTheme(): Theme {
  if (typeof window === 'undefined') return 'dark'
  if (window.matchMedia('(prefers-color-scheme: dark)').matches) return 'dark'
  if (window.matchMedia('(prefers-color-scheme: light)').matches) return 'light'
  return 'dark'
}

function getInitialTheme(): Theme {
  if (typeof window === 'undefined') return 'dark'
  const saved = localStorage.getItem(STORAGE_KEY)
  if (saved === 'dark' || saved === 'light') return saved
  return getSystemTheme()
}

export function useTheme() {
  const [theme, setTheme] = useState<Theme>(getInitialTheme)
  const [mounted, setMounted] = useState(false)

  useEffect(() => {
    setMounted(true)
  }, [])

  useEffect(() => {
    if (!mounted) return
    document.documentElement.classList.toggle('dark', theme === 'dark')
  }, [theme, mounted])

  useEffect(() => {
    const mq = window.matchMedia('(prefers-color-scheme: dark)')
    const handler = (e: MediaQueryListEvent) => {
      const saved = localStorage.getItem(STORAGE_KEY)
      if (!saved) {
        setTheme(e.matches ? 'dark' : 'light')
      }
    }
    mq.addEventListener('change', handler)
    return () => mq.removeEventListener('change', handler)
  }, [])

  useEffect(() => {
    const onTheme = (event: Event) => {
      const next = (event as CustomEvent<unknown>).detail
      if (isTheme(next)) setTheme(next)
    }
    const onStorage = (event: StorageEvent) => {
      if (event.key !== STORAGE_KEY) return
      if (isTheme(event.newValue)) setTheme(event.newValue)
    }
    window.addEventListener(THEME_EVENT, onTheme)
    window.addEventListener('storage', onStorage)
    return () => {
      window.removeEventListener(THEME_EVENT, onTheme)
      window.removeEventListener('storage', onStorage)
    }
  }, [])

  const toggleTheme = useCallback(() => {
    const next: Theme = theme === 'dark' ? 'light' : 'dark'
    setTheme(next)
    localStorage.setItem(STORAGE_KEY, next)
    document.documentElement.classList.toggle('dark', next === 'dark')
    window.dispatchEvent(new CustomEvent<Theme>(THEME_EVENT, { detail: next }))
  }, [theme])

  return { theme, mounted, toggleTheme }
}
