'use client'

import Link from 'next/link'
import { usePathname, useRouter } from 'next/navigation'
import { LogOut } from 'lucide-react'
import {
  BookOpen,
  Network,
  Play,
  Settings,
  Trophy,
  type LucideIcon,
} from 'lucide-react'
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

// Navigate group: destinations only. Never Diagnostic / Quiz / Review —
// those have dedicated cards in main content (no CTA duplication).
// Lucide icons (fixed size-4 box) so the collapsed 44px icon rail stays even.
const NAV_ITEMS: { label: string; href: string; icon: LucideIcon }[] = [
  { label: 'Study', href: '/study', icon: BookOpen },
  { label: 'Start', href: '/session', icon: Play },
  { label: 'Graph', href: '/graph', icon: Network },
  { label: 'Leaderboard', href: '/leaderboard', icon: Trophy },
  { label: 'Settings', href: '/settings', icon: Settings },
]

interface AppSidebarProps {
  name: string
  studentId: string
  username?: string
  level?: string
  streak?: number
  avatarUrl?: string
  avatarPreset?: number | null
  dueReviews?: number
}

export function AppSidebar({
  name,
  studentId,
  username,
  level,
  streak,
  avatarUrl,
  avatarPreset,
  dueReviews = 0,
}: AppSidebarProps) {
  const pathname = usePathname()
  const router = useRouter()

  return (
    <Sidebar collapsible="icon" data-testid="profile-sidebar">
      <SidebarHeader>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton size="lg" asChild tooltip="Mathua hub">
              <Link href="/profile">
                <span className="flex size-6 shrink-0 items-center justify-center font-mono text-sm text-mathua-blue" aria-hidden="true">λ</span>
                <span className="font-mono text-sm text-mathua-blue group-data-[collapsible=icon]:hidden">Mathua</span>
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
                const Icon = item.icon
                return (
                  <SidebarMenuItem key={item.href}>
                    <SidebarMenuButton asChild isActive={active} tooltip={item.label}>
                      <Link href={item.href}>
                        <Icon aria-hidden="true" className="size-4 shrink-0" />
                        <span className="group-data-[collapsible=icon]:hidden">{item.label}</span>
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
                <span className="min-w-0 flex-1 group-data-[collapsible=icon]:hidden">
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
              <span className="group-data-[collapsible=icon]:hidden">Sign out</span>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  )
}
