'use client'

import Link from 'next/link'
import { useState, useEffect, useRef } from 'react'
import { useRouter, usePathname } from 'next/navigation'
import { CompassIcon, type CompassIconHandle } from './icons/compass'
import { SunMoonIcon, type SunMoonIconHandle } from './icons/sun-moon'
import { useTheme } from '../hooks/useTheme'
import { useAuthState } from '../hooks/useAuthState'
import { signOut } from '../lib/auth'
import { resolveAvatar } from '../lib/dicebear'
import Avatar from './Avatar'
import { getSettings } from '../lib/api'

interface HeaderLink {
  label: string
  href: string
}

interface HeaderProps {
  links?: HeaderLink[]
}

export default function Header({ links }: HeaderProps) {
  const { mounted, toggleTheme } = useTheme()
  const { loggedIn, user } = useAuthState()
  const router = useRouter()
  const pathname = usePathname()
  const [open, setOpen] = useState(false)
  const [navOpen, setNavOpen] = useState(false)
  const [avatarUrl, setAvatarUrl] = useState<string | undefined>(undefined)
  const [avatarPreset, setAvatarPreset] = useState<number | null>(null)
  const menuRef = useRef<HTMLDivElement>(null)
  const navMenuRef = useRef<HTMLDivElement>(null)
  const compassRef = useRef<CompassIconHandle>(null)
  const sunMoonRef = useRef<SunMoonIconHandle>(null)

  function prefersReducedMotion(): boolean {
    return (
      typeof window !== 'undefined' &&
      window.matchMedia('(prefers-reduced-motion: reduce)').matches
    )
  }

  function toggleNav(): void {
    const next = !navOpen
    setNavOpen(next)
    // Alternate drive direction so the needle moves on every click: forward
    // spin on open, reverse spin on close (hover rarely fires on touch
    // screens; skipped entirely for reduced motion).
    if (prefersReducedMotion()) return
    if (next) compassRef.current?.startAnimation()
    else compassRef.current?.stopAnimation()
  }

  function handleThemeToggle(): void {
    if (!prefersReducedMotion()) sunMoonRef.current?.startAnimation()
    toggleTheme()
  }

  useEffect(() => {
    if (!loggedIn || !user) {
      setAvatarUrl(undefined)
      setAvatarPreset(null)
      return
    }
    // Dead-token cleanup is lazy: any strict API call below that answers
    // 401 clears the token via authedFetch (lib/auth.ts), which flips this
    // UI to logged-out through the auth-changed event. No explicit
    // revalidation here — the e2e auth spec renders the header with a
    // stored token while all /api/** routes are mocked, and must keep
    // showing the authenticated links.
    getSettings().then(s => {
      // Runs on every [loggedIn, user] change — including the auth-changed
      // broadcast that Settings fires after avatar saves — so the header
      // picture tracks the shared resolver with no extra event type.
      const resolved = resolveAvatar(user, s)
      setAvatarPreset(resolved.preset ?? null)
      setAvatarUrl(resolved.url)
    }).catch(() => {
      setAvatarPreset(null)
      setAvatarUrl(user.avatar_url)
    })
  }, [loggedIn, user])

  useEffect(() => {
    function onClick(e: MouseEvent) {
      if (menuRef.current && !menuRef.current.contains(e.target as Node)) setOpen(false)
      if (navMenuRef.current && !navMenuRef.current.contains(e.target as Node)) setNavOpen(false)
    }
    function onKey(e: KeyboardEvent) {
      if (e.key === 'Escape') {
        setOpen(false)
        setNavOpen(false)
      }
    }
    document.addEventListener('mousedown', onClick)
    document.addEventListener('keydown', onKey)
    return () => {
      document.removeEventListener('mousedown', onClick)
      document.removeEventListener('keydown', onKey)
    }
  }, [])

  const displayLinks = links || [
    { label: 'Study', href: '/study' },
    { label: 'Leaderboard', href: '/leaderboard' },
    { label: 'Graph', href: '/graph' },
    ...(loggedIn
      ? [{ label: 'Profile', href: '/profile' }, { label: 'Start', href: '/session' }, { label: 'Settings', href: '/settings' }]
      : [{ label: 'Login', href: '/login' }]
    ),
  ]

  return (
    <header className="sticky top-0 z-50 backdrop-blur border-b border-mathua-border bg-mathua-bg transition-colors">
      <div className="flex items-center justify-between px-4 md:px-6 py-3 max-w-container mx-auto">
        <div className="flex items-center gap-4 md:gap-6 min-w-0">
          <Link
            href={loggedIn ? '/profile' : '/'}
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

        <div className="flex items-center gap-2 sm:gap-3 shrink-0">
          {mounted && (
            <button
              onClick={handleThemeToggle}
              aria-label="Toggle theme"
              className="flex items-center justify-center min-h-[44px] min-w-[44px] text-mathua-muted bg-transparent border-none cursor-pointer hover:text-mathua-blue transition-colors"
            >
              <SunMoonIcon ref={sunMoonRef} size={15} aria-hidden="true" />
            </button>
          )}
          {loggedIn && user ? (
            <div className="relative" ref={menuRef}>
              <button
                onClick={() => setOpen(o => !o)}
                aria-label="Open profile menu"
                aria-expanded={open}
                aria-haspopup="menu"
                className="rounded-full p-0.5 border border-transparent hover:border-mathua-border focus:outline-none focus:border-mathua-blue transition-colors"
              >
                <Avatar
                  seed={user.student_id}
                  name={user.name}
                  size={32}
                  url={avatarPreset !== null ? undefined : avatarUrl}
                  preset={avatarPreset ?? undefined}
                />
              </button>
              {open && (
                <div className="absolute right-0 top-[calc(100%+8px)] min-w-[160px] bg-mathua-surface border border-mathua-border shadow-lg py-1 z-50">
                  <div className="px-3 py-2 border-b border-mathua-border">
                    <div className="font-mono text-xs text-mathua-primary truncate">{user.name}</div>
                    {user.username && <div className="font-mono text-[10px] text-mathua-muted truncate">@{user.username}</div>}
                  </div>
                  <button
                    onClick={() => {
                      setOpen(false)
                      signOut()
                      router.push('/login')
                    }}
                    className="w-full text-left px-3 py-2 font-mono text-xs text-mathua-muted hover:text-mathua-red hover:bg-mathua-surface-elevated transition-colors"
                  >
                    Sign out
                  </button>
                  <button
                    onClick={() => {
                      setOpen(false)
                      signOut()
                      router.push('/profile')
                    }}
                    className="w-full text-left px-3 py-2 font-mono text-xs text-mathua-muted hover:text-mathua-blue hover:bg-mathua-surface-elevated transition-colors"
                  >
                    Profile
                  </button>
                </div>
              )}
            </div>
          ) : null}
          {/* Mobile nav: bare compass mark at the far right, mirroring
              the desktop link row (hidden below md) on every page. */}
          <div className="relative md:hidden" ref={navMenuRef}>
              <button
                onClick={toggleNav}
                aria-label="Open navigation menu"
                aria-expanded={navOpen}
                aria-haspopup="menu"
                className="flex items-center justify-center min-h-[44px] min-w-[44px] text-mathua-muted bg-transparent border-none cursor-pointer hover:text-mathua-blue transition-colors"
              >
                <CompassIcon ref={compassRef} size={15} aria-hidden="true" />
              </button>
              {navOpen && (
                <div
                  role="menu"
                  aria-label="Site navigation"
                  className="absolute right-0 top-[calc(100%+8px)] min-w-[200px] bg-mathua-surface border border-mathua-border shadow-lg py-1 z-50"
                >
                  {displayLinks.map(link => {
                    const active = pathname === link.href || (pathname ?? '').startsWith(link.href + '/')
                    return (
                      <Link
                        key={link.href}
                        href={link.href}
                        role="menuitem"
                        aria-current={active ? 'page' : undefined}
                        onClick={() => setNavOpen(false)}
                        className={`flex items-center min-h-[44px] px-4 font-mono text-xs whitespace-nowrap transition-colors ${
                          active ? 'text-mathua-blue' : 'text-mathua-muted hover:text-mathua-blue hover:bg-mathua-surface-elevated'
                        }`}
                      >
                        {link.label}
                      </Link>
                    )
                  })}
                </div>
              )}
          </div>
        </div>
      </div>
    </header>
  )
}
