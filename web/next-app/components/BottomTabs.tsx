'use client'

import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { useState, useEffect, useRef } from 'react'
import {
  House,
  BookOpen,
  Play,
  User,
  Network,
  type LucideIcon,
} from 'lucide-react'
import { useChrome } from '../hooks/useChrome'
import { useInertWhen } from '../hooks/useInertWhen'
import { registerChromePart } from '../lib/chrome'
import { useAuthState } from '../hooks/useAuthState'
import { resolveAvatar } from '../lib/dicebear'
import { getSettings } from '../lib/api'
import Avatar from './Avatar'

// Learn before Study: the tab bar is the mobile primary nav, and Learn is the
// only tab that asks a question. Study stays — it is how a learner looks up a
// concept they have just met. Order is mirrored in lib/nav.ts.
const TABS: { label: string; href: string; icon: LucideIcon }[] = [
  { label: 'Home', href: '/', icon: House },
  { label: 'Learn', href: '/learn', icon: Play },
  { label: 'Study', href: '/study', icon: BookOpen },
  { label: 'Profile', href: '/profile', icon: User },
  { label: 'Graph', href: '/graph', icon: Network },
]

export default function BottomTabs() {
  const pathname = usePathname()
  // Chrome visibility is shared with the header so the two always agree --
  // a header that hid while the tab bar stayed would read as a bug, not a
  // gesture. See lib/chrome.ts.
  const { visible: chromeVisible } = useChrome()
  const hidden = !chromeVisible
  const navRef = useRef<HTMLElement>(null)
  useInertWhen(navRef, hidden)
  const { loggedIn, user } = useAuthState()
  const [avatarUrl, setAvatarUrl] = useState<string | undefined>(undefined)
  const [avatarPreset, setAvatarPreset] = useState<number | null>(null)

  // Same pattern as Header: resolve through the shared resolver on every
  // auth change (login/logout/avatar save broadcast), lazy 401 handling —
  // no explicit revalidation, so the mocked e2e auth contract holds.
  useEffect(() => {
    if (!loggedIn || !user) {
      setAvatarUrl(undefined)
      setAvatarPreset(null)
      return
    }
    getSettings().then(s => {
      const resolved = resolveAvatar(user, s)
      setAvatarPreset(resolved.preset ?? null)
      setAvatarUrl(resolved.url)
    }).catch(() => {
      setAvatarPreset(null)
      setAvatarUrl(user.avatar_url)
    })
  }, [loggedIn, user])

  // Tell the chrome store a tab bar exists, so the tap control knows to sit
  // above it rather than at the bare bottom edge.
  useEffect(() => registerChromePart('tabs'), [])

  return (
    <nav
      ref={navRef}
      aria-label="Primary"
      data-chrome-part="tabs"
      className={`lg:hidden fixed inset-x-3 z-40 flex h-[56px] border-[0.5px] border-mathua-border bg-mathua-surface shadow-none transition-transform duration-300 ease-out will-change-transform motion-reduce:transition-none ${
        hidden ? 'translate-y-[calc(100%+20px)] opacity-0' : 'translate-y-0 opacity-100'
      }`}
      style={{
        bottom: 'max(12px, env(safe-area-inset-bottom))',
        paddingBottom: 'env(safe-area-inset-bottom)',
        left: '12px',
        right: '12px',
        maxWidth: 'min(360px, calc(100% - 24px))',
        marginLeft: 'auto',
        marginRight: 'auto',
      }}
    >
      {TABS.map((t) => {
        const active = pathname === t.href || pathname.startsWith(t.href + '/')
        const showAvatar = t.href === '/profile' && loggedIn && user
        const Icon = t.icon
        return (
          <Link
            key={t.href}
            href={t.href}
            aria-label={t.label}
            aria-current={active ? 'page' : undefined}
            className={`flex-1 flex flex-col items-center justify-center gap-0.5 font-mono border-t-[2px] min-h-[44px] transition-colors ${
              active ? 'border-mathua-blue text-mathua-blue' : 'border-transparent text-mathua-muted hover:text-mathua-primary'
            }`}
          >
            {showAvatar ? (
              <span className={`rounded-full leading-none ${active ? 'ring-1 ring-mathua-blue ring-offset-1 ring-offset-mathua-surface' : ''}`}>
                <Avatar
                  seed={user.student_id}
                  name={user.name}
                  size={22}
                  url={avatarPreset !== null ? undefined : avatarUrl}
                  preset={avatarPreset ?? undefined}
                />
              </span>
            ) : (
              <Icon size={18} strokeWidth={2} aria-hidden="true" />
            )}
            <span className="text-[9px] uppercase tracking-wide">{t.label}</span>
          </Link>
        )
      })}
    </nav>
  )
}
