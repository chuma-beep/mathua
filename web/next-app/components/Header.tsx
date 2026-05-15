'use client'

import Link from 'next/link'
import { useEffect, useState, useCallback } from 'react'
import { useTheme } from '../hooks/useTheme'

interface HeaderLink {
  label: string
  href: string
}

interface HeaderProps {
  links?: HeaderLink[]
}

export default function Header({ links }: HeaderProps) {
  const { theme, mounted, toggleTheme } = useTheme()
  const [menuOpen, setMenuOpen] = useState(false)
  const [loggedIn, setLoggedIn] = useState(false)

  useEffect(() => {
    if (typeof window !== 'undefined') {
      setLoggedIn(!!localStorage.getItem('mathua_token'))
    }
  }, [mounted])

  const closeMenu = useCallback(() => setMenuOpen(false), [])

  useEffect(() => {
    if (!menuOpen) return
    const onKey = (e: KeyboardEvent) => {
      if (e.key === 'Escape') closeMenu()
    }
    window.addEventListener('keydown', onKey)
    return () => window.removeEventListener('keydown', onKey)
  }, [menuOpen, closeMenu])

  const displayLinks = links || [
    { label: 'Leaderboard', href: '/leaderboard' },
    { label: 'Graph', href: '/graph' },
    ...(loggedIn
      ? [{ label: 'Practice', href: '/session' }]
      : [{ label: 'Login', href: '/login' }]
    ),
  ]

  return (
    <header className="sticky top-0 z-50 backdrop-blur border-b border-mathua-border bg-mathua-bg transition-colors">
      <div className="flex items-center justify-between px-4 md:px-6 py-3 max-w-container mx-auto">
        <div className="flex items-center gap-4 md:gap-6 min-w-0">
          <Link
            href="/"
            className="link-underline font-mono text-sm text-mathua-blue whitespace-nowrap shrink-0"
          >
            λ Mathua
          </Link>

          <nav className="hidden md:flex items-center gap-2">
            {displayLinks.map((link, i) => (
              <span key={link.href} className="flex items-center gap-2">
                {i > 0 && (
                  <span className="font-mono text-xs text-mathua-blue select-none">|</span>
                )}
                <Link
                  href={link.href}
                  className="link-underline font-mono text-xs text-mathua-muted hover:text-mathua-blue whitespace-nowrap transition-colors"
                >
                  {link.label}
                </Link>
              </span>
            ))}
          </nav>
        </div>

        <div className="flex items-center gap-3 shrink-0">
          {mounted && (
            <button
              onClick={toggleTheme}
              aria-label="Toggle theme"
              className="font-mono text-xs text-mathua-muted bg-transparent border border-mathua-border-strong px-2.5 py-1 cursor-pointer rounded-none hover:text-mathua-blue hover:border-mathua-blue transition-colors"
            >
              {theme === 'dark' ? '\u2600' : '\u263E'}
            </button>
          )}

          <button
            onClick={() => setMenuOpen((p) => !p)}
            className="md:hidden flex items-center justify-center size-6 text-mathua-muted hover:text-mathua-blue transition-colors"
            aria-label="Toggle navigation menu"
          >
            <span
              className={`block font-mono text-lg leading-none transition-transform duration-200 ${
                menuOpen ? 'rotate-45' : ''
              }`}
            >
              +
            </span>
          </button>
        </div>
      </div>

      <div
        className={`md:hidden overflow-hidden transition-all duration-200 border-t border-mathua-border bg-mathua-bg ${
          menuOpen ? 'max-h-80' : 'max-h-0'
        }`}
      >
        <nav className="flex flex-col px-6 py-4 gap-3">
          {displayLinks.map((link) => (
            <Link
              key={link.href}
              href={link.href}
              onClick={closeMenu}
              className="font-mono text-xs text-mathua-muted hover:text-mathua-blue transition-colors"
            >
              {link.label}
            </Link>
          ))}
        </nav>
      </div>
    </header>
  )
}
