'use client'

import type { ReactNode } from 'react'
import Link from 'next/link'
import Avatar from './Avatar'

export interface ProfileTocItem {
  id: string
  label: string
}

interface ProfileSidebarProps {
  name: string
  studentId: string
  level?: string
  streak?: number
  avatarUrl?: string
  avatarPreset?: number | null
  toc: ProfileTocItem[]
  activeId: string
  /** T2 plug-in: Next-up compact card */
  nextSlot?: ReactNode
  /** T3/T4 plug-in: goal stepper + weekly chart */
  goalsSlot?: ReactNode
  onNavigate?: (id: string) => void
}

const QUICK_LINKS = [
  { label: 'Study', href: '/study' },
  { label: 'Diagnostic', href: '/onboard' },
  { label: 'Leaderboard', href: '/leaderboard' },
  { label: 'Graph', href: '/graph' },
  { label: 'Settings', href: '/settings' },
]

export default function ProfileSidebar({
  name,
  studentId,
  level,
  streak,
  avatarUrl,
  avatarPreset,
  toc,
  activeId,
  nextSlot,
  goalsSlot,
  onNavigate,
}: ProfileSidebarProps) {
  return (
    <aside
      aria-label="Profile sidebar"
      className="hidden lg:flex lg:flex-col gap-4 lg:sticky lg:top-[64px] lg:self-start lg:max-h-[calc(100vh-80px)] lg:overflow-auto min-w-0"
    >
      {/* Identity mini */}
      <div className="border border-mathua-border bg-mathua-surface p-4 min-w-0">
        <div className="flex items-center gap-3 min-w-0">
          <Avatar
            seed={studentId}
            name={name}
            size={40}
            url={avatarPreset !== null && avatarPreset !== undefined ? undefined : avatarUrl}
            preset={avatarPreset ?? undefined}
          />
          <div className="min-w-0 flex-1">
            <div className="font-mono text-xs text-mathua-primary truncate" title={name}>
              {name}
            </div>
            <div className="font-mono text-[10px] text-mathua-muted truncate">
              {[level, typeof streak === 'number' ? `${streak}d streak` : null].filter(Boolean).join(' · ')}
            </div>
          </div>
        </div>
      </div>

      {nextSlot}

      {goalsSlot}

      {/* Quick links */}
      <nav aria-label="Profile quick links" className="border border-mathua-border bg-mathua-surface p-2">
        {QUICK_LINKS.map((l) => (
          <Link
            key={l.href}
            href={l.href}
            className="flex items-center px-2.5 py-2 min-h-[44px] font-mono text-xs text-mathua-muted hover:text-mathua-blue hover:bg-mathua-surface-elevated transition-colors"
          >
            {l.label}
          </Link>
        ))}
      </nav>

      {/* TOC */}
      <nav aria-label="On this page" className="border border-mathua-border bg-mathua-surface p-2">
        <div className="font-mono text-[10px] uppercase text-mathua-muted px-2.5 pt-1.5 pb-1 tracking-wider">
          On this page
        </div>
        {toc.map((t) => {
          const active = t.id === activeId
          return (
            <a
              key={t.id}
              href={`#${t.id}`}
              aria-current={active ? 'true' : undefined}
              onClick={(e) => {
                e.preventDefault()
                onNavigate?.(t.id)
                document.getElementById(t.id)?.scrollIntoView({ behavior: 'smooth', block: 'start' })
              }}
              className={`flex items-center px-2.5 py-2 min-h-[44px] font-mono text-xs border-l-2 transition-colors ${
                active
                  ? 'border-mathua-blue text-mathua-blue bg-mathua-surface-elevated'
                  : 'border-transparent text-mathua-muted hover:text-mathua-primary hover:bg-mathua-surface-elevated'
              }`}
            >
              {t.label}
            </a>
          )
        })}
      </nav>
    </aside>
  )
}
