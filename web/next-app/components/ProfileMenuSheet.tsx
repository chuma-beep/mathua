'use client'

import { useEffect, useRef } from 'react'
import type { ReactNode } from 'react'
import type { ProfileTocItem } from './ProfileSidebar'
import ProfileSidebar from './ProfileSidebar'

interface ProfileMenuSheetProps {
  open: boolean
  onClose: () => void
  sidebarProps: Omit<React.ComponentProps<typeof ProfileSidebar>, 'onNavigate'> & { onNavigate?: (id: string) => void }
  chips: ProfileTocItem[]
  activeId: string
  chipsHidden: boolean
  onChipNavigate: (id: string) => void
  onOpenMenu: () => void
}

export function ProfileMobileBar({
  chips,
  activeId,
  hidden,
  onChipNavigate,
  onOpenMenu,
}: {
  chips: ProfileTocItem[]
  activeId: string
  hidden: boolean
  onChipNavigate: (id: string) => void
  onOpenMenu: () => void
}) {
  return (
    <div
      className={`lg:hidden sticky top-[53px] z-30 border-b border-mathua-border bg-mathua-bg/95 backdrop-blur transition-transform duration-300 ${
        hidden ? '-translate-y-[calc(100%+8px)] opacity-0 pointer-events-none' : 'translate-y-0 opacity-100'
      }`}
    >
      <div className="flex items-center gap-2 px-3 py-2 overflow-x-auto">
        <button
          onClick={onOpenMenu}
          aria-label="Open profile menu"
          aria-haspopup="dialog"
          className="shrink-0 font-mono text-xs text-mathua-muted border border-mathua-border-strong px-3 min-h-[44px] inline-flex items-center hover:text-mathua-blue hover:border-mathua-blue transition-colors"
        >
          ☰ Menu
        </button>
        {chips.map((c) => {
          const active = c.id === activeId
          return (
            <button
              key={c.id}
              onClick={() => onChipNavigate(c.id)}
              aria-current={active ? 'true' : undefined}
              className={`shrink-0 font-mono text-[11px] px-3 min-h-[44px] inline-flex items-center border transition-colors ${
                active
                  ? 'border-mathua-blue text-mathua-blue bg-mathua-surface'
                  : 'border-mathua-border text-mathua-muted hover:text-mathua-primary'
              }`}
            >
              {c.label}
            </button>
          )
        })}
      </div>
    </div>
  )
}

export default function ProfileMenuSheet({
  open,
  onClose,
  sidebarProps,
  chips: _chips,
  activeId: _activeId,
  chipsHidden: _chipsHidden,
  onChipNavigate: _onChipNavigate,
  onOpenMenu: _onOpenMenu,
}: ProfileMenuSheetProps) {
  const closeRef = useRef<HTMLButtonElement>(null)
  const panelRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    if (!open) return
    closeRef.current?.focus()
    const onKey = (e: KeyboardEvent) => {
      if (e.key === 'Escape') onClose()
      // Simple focus trap: keep Tab inside panel
      if (e.key === 'Tab' && panelRef.current) {
        const focusables = panelRef.current.querySelectorAll<HTMLElement>(
          'a[href], button:not([disabled]), input, [tabindex]:not([tabindex="-1"])',
        )
        if (focusables.length === 0) return
        const first = focusables[0]
        const last = focusables[focusables.length - 1]
        if (e.shiftKey && document.activeElement === first) {
          e.preventDefault()
          last.focus()
        } else if (!e.shiftKey && document.activeElement === last) {
          e.preventDefault()
          first.focus()
        }
      }
    }
    document.addEventListener('keydown', onKey)
    const prevOverflow = document.body.style.overflow
    document.body.style.overflow = 'hidden'
    return () => {
      document.removeEventListener('keydown', onKey)
      document.body.style.overflow = prevOverflow
    }
  }, [open, onClose])

  if (!open) return null

  return (
    <div
      role="dialog"
      aria-modal="true"
      aria-label="Profile menu"
      className="lg:hidden fixed inset-0 z-50"
    >
      <button
        aria-label="Close profile menu"
        onClick={onClose}
        className="absolute inset-0 bg-black/50 cursor-default"
        tabIndex={-1}
      />
      <div
        ref={panelRef}
        className="absolute inset-x-0 bottom-0 max-h-[85vh] overflow-auto border-t border-mathua-border bg-mathua-bg p-4 pb-[calc(16px+env(safe-area-inset-bottom))]"
        style={{ paddingBottom: 'max(16px, env(safe-area-inset-bottom))' }}
      >
        <div className="flex items-center justify-between mb-3">
          <span className="font-mono text-[11px] uppercase text-mathua-muted tracking-wider">Profile menu</span>
          <button
            ref={closeRef}
            onClick={onClose}
            aria-label="Close profile menu"
            className="font-mono text-xs text-mathua-muted border border-mathua-border-strong px-4 min-h-[44px] inline-flex items-center hover:text-mathua-blue hover:border-mathua-blue transition-colors"
          >
            ✕ Close
          </button>
        </div>
        {/* Reuse desktop sidebar content stacked; its aside is hidden on mobile via CSS,
            so render inner content through a mobile-visible wrapper instead. */}
        <MobileSidebarContent {...sidebarProps} onNavigate={(id) => { sidebarProps.onNavigate?.(id); onClose() }} />
      </div>
    </div>
  )
}

// Mobile-visible clone of sidebar content (aside root in ProfileSidebar is lg-only).
function MobileSidebarContent(props: React.ComponentProps<typeof ProfileSidebar>) {
  return (
    <div className="lg:hidden flex flex-col gap-4 min-w-0">
      <ProfileSidebarInner {...props} />
    </div>
  )
}

// Inline the sidebar blocks without the lg-only aside wrapper.
import Avatar from './Avatar'
import Link from 'next/link'

function ProfileSidebarInner({
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
}: React.ComponentProps<typeof ProfileSidebar>) {
  const QUICK_LINKS = [
    { label: 'Study', href: '/study' },
    { label: 'Diagnostic', href: '/onboard' },
    { label: 'Leaderboard', href: '/leaderboard' },
    { label: 'Graph', href: '/graph' },
    { label: 'Settings', href: '/settings' },
  ]
  const go = (id: string) => {
    onNavigate?.(id)
    document.getElementById(id)?.scrollIntoView({ behavior: 'smooth', block: 'start' })
  }
  return (
    <>
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
            <div className="font-mono text-xs text-mathua-primary truncate">{name}</div>
            <div className="font-mono text-[10px] text-mathua-muted truncate">
              {[level, typeof streak === 'number' ? `${streak}d streak` : null].filter(Boolean).join(' · ')}
            </div>
          </div>
        </div>
      </div>
      {nextSlot as ReactNode}
      {goalsSlot as ReactNode}
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
      <nav aria-label="On this page" className="border border-mathua-border bg-mathua-surface p-2">
        <div className="font-mono text-[10px] uppercase text-mathua-muted px-2.5 pt-1.5 pb-1 tracking-wider">
          On this page
        </div>
        {toc.map((t) => (
          <a
            key={t.id}
            href={`#${t.id}`}
            aria-current={t.id === activeId ? 'true' : undefined}
            onClick={(e) => {
              e.preventDefault()
              go(t.id)
            }}
            className={`flex items-center px-2.5 py-2 min-h-[44px] font-mono text-xs border-l-2 transition-colors ${
              t.id === activeId
                ? 'border-mathua-blue text-mathua-blue bg-mathua-surface-elevated'
                : 'border-transparent text-mathua-muted hover:text-mathua-primary hover:bg-mathua-surface-elevated'
            }`}
          >
            {t.label}
          </a>
        ))}
      </nav>
    </>
  )
}
