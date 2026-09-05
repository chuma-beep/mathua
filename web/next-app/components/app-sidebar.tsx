'use client'

import type { ReactNode } from 'react'
import Link from 'next/link'
import { usePathname, useRouter } from 'next/navigation'
import { LogOut } from 'lucide-react'
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuBadge,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarRail,
} from '@/components/ui/sidebar'
import Avatar from './Avatar'
import { signOut } from '../lib/auth'

export interface ProfileTocItem {
  id: string
  label: string
}

// Navigate group: destinations only. Never Diagnostic / Quiz / Review —
// those have dedicated cards in main content (no CTA duplication).
const NAV_ITEMS = [
  { label: 'Study', href: '/study', icon: '◐' },
  { label: 'Start', href: '/session', icon: 'π' },
  { label: 'Graph', href: '/graph', icon: '⬡' },
  { label: 'Leaderboard', href: '/leaderboard', icon: '◎' },
  { label: 'Settings', href: '/settings', icon: '⚙' },
]

interface AppSidebarProps {
  name: string
  studentId: string
  username?: string
  level?: string
  streak?: number
  avatarUrl?: string
  avatarPreset?: number | null
  toc: ProfileTocItem[]
  activeId: string
  dueReviews?: number
  nextSlot?: ReactNode
  goalsSlot?: ReactNode
  onNavigate?: (id: string) => void
}

export function AppSidebar({
  name,
  studentId,
  username,
  level,
  streak,
  avatarUrl,
  avatarPreset,
  toc,
  activeId,
  dueReviews = 0,
  nextSlot,
  goalsSlot,
  onNavigate,
}: AppSidebarProps) {
  const pathname = usePathname()
  const router = useRouter()

  const scrollTo = (id: string) => {
    onNavigate?.(id)
    document.getElementById(id)?.scrollIntoView({ behavior: 'smooth', block: 'start' })
  }

  return (
    <Sidebar collapsible="icon" data-testid="profile-sidebar">
      <SidebarHeader>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton size="lg" asChild tooltip="Mathua hub">
              <Link href="/profile">
                <span className="font-mono text-sm text-mathua-blue">λ</span>
                <span className="font-mono text-sm text-mathua-blue">Mathua</span>
              </Link>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>

      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupLabel>Navigate</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              {NAV_ITEMS.map((item) => {
                const active = pathname === item.href || pathname.startsWith(item.href + '/')
                return (
                  <SidebarMenuItem key={item.href}>
                    <SidebarMenuButton asChild isActive={active} tooltip={item.label}>
                      <Link href={item.href}>
                        <span aria-hidden="true">{item.icon}</span>
                        <span>{item.label}</span>
                      </Link>
                    </SidebarMenuButton>
                    {item.href === '/session' && dueReviews > 0 && (
                      <SidebarMenuBadge>{dueReviews}</SidebarMenuBadge>
                    )}
                  </SidebarMenuItem>
                )
              })}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>

        {nextSlot && (
          <SidebarGroup>
            <SidebarGroupLabel>Next up</SidebarGroupLabel>
            <SidebarGroupContent className="px-2 group-data-[collapsible=icon]:hidden">
              {nextSlot}
            </SidebarGroupContent>
          </SidebarGroup>
        )}

        {goalsSlot && (
          <SidebarGroup>
            <SidebarGroupLabel>Goals</SidebarGroupLabel>
            <SidebarGroupContent className="px-2 group-data-[collapsible=icon]:hidden">
              {goalsSlot}
            </SidebarGroupContent>
          </SidebarGroup>
        )}

        <SidebarGroup>
          <SidebarGroupLabel>On this page</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              {toc.map((t) => (
                <SidebarMenuItem key={t.id}>
                  <SidebarMenuButton
                    isActive={t.id === activeId}
                    tooltip={t.label}
                    onClick={() => scrollTo(t.id)}
                  >
                    <span aria-hidden="true">§</span>
                    <span>{t.label}</span>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              ))}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>

      <SidebarFooter>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton size="lg" asChild tooltip={name}>
              <Link href="/profile">
                <Avatar
                  seed={studentId}
                  name={name}
                  size={32}
                  url={avatarPreset !== null && avatarPreset !== undefined ? undefined : avatarUrl}
                  preset={avatarPreset ?? undefined}
                />
                <span className="min-w-0 flex-1">
                  <span className="block truncate font-mono text-xs text-mathua-primary">{name}</span>
                  <span className="block truncate font-mono text-[10px] text-mathua-muted">
                    {[username ? `@${username}` : null, level, typeof streak === 'number' ? `${streak}d` : null]
                      .filter(Boolean)
                      .join(' · ')}
                  </span>
                </span>
              </Link>
            </SidebarMenuButton>
          </SidebarMenuItem>
          <SidebarMenuItem>
            <SidebarMenuButton
              tooltip="Sign out"
              onClick={() => {
                signOut()
                router.push('/login')
              }}
            >
              <LogOut aria-hidden="true" />
              <span>Sign out</span>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  )
}
