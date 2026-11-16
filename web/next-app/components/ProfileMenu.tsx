'use client'

import { useEffect, useRef, useState } from 'react'
import Link from 'next/link'
import Avatar from './Avatar'
import { Shield, Settings, User } from './icons'
import { SidebarMenuButton } from '@/components/ui/sidebar'

type ProfileMenuProps = {
  name: string
  studentId: string
  username?: string
  level?: string
  streak?: number
  avatarUrl?: string
  avatarPreset?: number | null
  isAdmin?: boolean
  adminReports?: number
}

// The profile photo is the account menu. Putting Admin here keeps it out of
// the daily learning navigation while leaving it one click away for staff, and
// still opens the console in its own tab. The entry is rendered only for a
// staff role; the server re-checks every admin request regardless.
export default function ProfileMenu({
  name,
  studentId,
  username,
  level,
  streak,
  avatarUrl,
  avatarPreset,
  isAdmin = false,
  adminReports = 0,
}: ProfileMenuProps) {
  const [open, setOpen] = useState(false)
  const root = useRef<HTMLDivElement>(null)

  useEffect(() => {
    if (!open) return
    const onPointer = (event: PointerEvent) => {
      if (root.current && !root.current.contains(event.target as Node)) setOpen(false)
    }
    const onKey = (event: KeyboardEvent) => {
      if (event.key === 'Escape') setOpen(false)
    }
    document.addEventListener('pointerdown', onPointer)
    document.addEventListener('keydown', onKey)
    return () => {
      document.removeEventListener('pointerdown', onPointer)
      document.removeEventListener('keydown', onKey)
    }
  }, [open])

  return (
    <div ref={root} className="relative">
      <SidebarMenuButton
        size="lg"
        tooltip={name}
        aria-haspopup="menu"
        aria-expanded={open}
        aria-label={`${name}, account menu`}
        onClick={() => setOpen((value) => !value)}
      >
        <Avatar
          seed={studentId}
          name={name}
          size={32}
          url={avatarPreset !== null && avatarPreset !== undefined ? undefined : avatarUrl}
          preset={avatarPreset ?? undefined}
        />
        <span className="min-w-0 flex-1 text-left group-data-[collapsible=icon]:hidden">
          <span className="block truncate font-mono text-xs text-mathua-primary">{name}</span>
          <span className="block truncate font-mono text-[10px] text-mathua-muted">
            {[username ? `@${username}` : null, level, typeof streak === 'number' ? `${streak}d` : null]
              .filter(Boolean)
              .join(' · ')}
          </span>
        </span>
      </SidebarMenuButton>

      {open && (
        <div
          role="menu"
          aria-label="Account"
          className="absolute bottom-full left-0 z-50 mb-1 w-56 border border-mathua-border bg-mathua-surface py-1 shadow-lg"
        >
          <Link
            href="/profile"
            role="menuitem"
            onClick={() => setOpen(false)}
            className="flex h-10 items-center gap-2.5 px-3 font-mono text-xs text-mathua-secondary hover:bg-mathua-surface-elevated hover:text-mathua-primary"
          >
            <User size={15} aria-hidden="true" />
            Profile
          </Link>
          <Link
            href="/settings"
            role="menuitem"
            onClick={() => setOpen(false)}
            className="flex h-10 items-center gap-2.5 px-3 font-mono text-xs text-mathua-secondary hover:bg-mathua-surface-elevated hover:text-mathua-primary"
          >
            <Settings size={15} aria-hidden="true" />
            Settings
          </Link>
          {isAdmin && (
            <>
              <div className="my-1 border-t border-mathua-border" />
              <a
                href="/admin"
                target="_blank"
                rel="noreferrer"
                role="menuitem"
                onClick={() => setOpen(false)}
                aria-label="Admin console, opens in a new tab"
                className="flex h-10 items-center gap-2.5 px-3 font-mono text-xs text-mathua-secondary hover:bg-mathua-surface-elevated hover:text-mathua-blue"
              >
                <Shield size={15} aria-hidden="true" />
                <span className="flex-1">Admin console</span>
                {adminReports > 0 && (
                  <span className="font-mono text-[10px] tabular-nums text-mathua-red" aria-label={`${adminReports} reports need attention`}>
                    {adminReports}
                  </span>
                )}
              </a>
            </>
          )}
        </div>
      )}
    </div>
  )
}
