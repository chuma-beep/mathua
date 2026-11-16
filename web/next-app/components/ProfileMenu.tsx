'use client'

import Link from 'next/link'
import Avatar from './Avatar'
import { Shield, Settings, User } from './icons'
import { SidebarMenuButton } from '@/components/ui/sidebar'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'

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
  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <SidebarMenuButton size="lg" tooltip={name} aria-label={`${name}, account menu`}>
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
      </DropdownMenuTrigger>
      <DropdownMenuContent side="top" align="start" className="w-56 rounded-none border-mathua-border bg-mathua-surface">
        <DropdownMenuItem asChild className="font-mono text-xs text-mathua-secondary focus:bg-mathua-surface-elevated focus:text-mathua-primary">
          <Link href="/profile">
            <User size={15} aria-hidden="true" />
            Profile
          </Link>
        </DropdownMenuItem>
        <DropdownMenuItem asChild className="font-mono text-xs text-mathua-secondary focus:bg-mathua-surface-elevated focus:text-mathua-primary">
          <Link href="/settings">
            <Settings size={15} aria-hidden="true" />
            Settings
          </Link>
        </DropdownMenuItem>
        {isAdmin && (
          <>
            <DropdownMenuSeparator className="bg-mathua-border" />
            <DropdownMenuItem asChild className="font-mono text-xs text-mathua-secondary focus:bg-mathua-surface-elevated focus:text-mathua-blue">
              <a href="/admin" target="_blank" rel="noreferrer" aria-label="Admin console, opens in a new tab">
                <Shield size={15} aria-hidden="true" />
                <span className="flex-1">Admin console</span>
                {adminReports > 0 && (
                  <span className="font-mono text-[10px] tabular-nums text-mathua-red" aria-label={`${adminReports} reports need attention`}>
                    {adminReports}
                  </span>
                )}
              </a>
            </DropdownMenuItem>
          </>
        )}
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
