'use client'

import Link from 'next/link'
import { usePathname, useRouter } from 'next/navigation'
import { useRef } from 'react'
import {
  ChartNoAxesColumn,
  ChevronsLeft,
  ChevronsRight,
  FileText,
  Handshake,
  Layers,
  LogOut,
  Shield,
  Network,
  Pencil,
  Play,
  Settings,
  SunMoon,
  Trophy,
  type IconComponent,
  type IconHandle,
} from './icons'
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
  useSidebar,
} from '@/components/ui/sidebar'
import Avatar from './Avatar'
import { useTheme } from '../hooks/useTheme'
import { signOut } from '../lib/auth'

// Navigate group: destinations only. Never Diagnostic / Quiz / Review —
// those have dedicated cards in main content (no CTA duplication). Review is
// reachable from the header row and the mobile compass instead, which is where
// the full destination list lives; putting it here as well would be exactly the
// duplication this rule exists to prevent.
// Order mirrors lib/nav.ts: Learn leads, then the orientation surfaces.
// Icons come from components/icons, which fixes the box at 16px. The animated set renders a
// div rather than an svg, so a `size-4` class here would silently size nothing and the
// collapsed 44px rail would go ragged.
interface SidebarNavItem {
  label: string
  href: string
  icon: IconComponent
}

const NAV_ITEMS: SidebarNavItem[] = [
  { label: 'Learn', href: '/learn', icon: Play },
  { label: 'Graph', href: '/graph', icon: Network },
  { label: 'Domains', href: '/domains', icon: Layers },
  { label: 'Leaderboard', href: '/leaderboard', icon: Trophy },
  { label: 'History', href: '/history', icon: ChartNoAxesColumn },
  { label: 'Settings', href: '/settings', icon: Settings },
]

// Resources group: docs, contributing, creator's note. Informational pages
// only — never Diagnostic / Quiz / Review (no CTA duplication).
const RESOURCE_ITEMS: SidebarNavItem[] = [
  { label: 'Docs', href: '/docs', icon: FileText },
  { label: 'Contribute', href: '/docs/contributing', icon: Handshake },
  { label: "Creator's note", href: '/note', icon: Pencil },
]

const ADMIN_ITEM: SidebarNavItem = {
  label: 'Admin',
  href: '/admin',
  icon: Shield,
}

export interface AppSidebarProps {
  name: string
  studentId: string
  /**
   * Whether to offer the Admin entry.
   *
   * A convenience and nothing more — the server reads the role on every administrative request,
   * so a learner who types /admin/users reaches the page and is refused by the API with a 403.
   * Omitting this only removes the invitation; it removes no access. It is here because the
   * alternative is an administrator having to know a URL.
   */
  isAdmin?: boolean
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
  isAdmin = false,
}: AppSidebarProps) {
  const pathname = usePathname()
  const router = useRouter()
  const { toggleSidebar } = useSidebar()
  // The theme control used to live only in Header, and /profile does not render one, so this
  // sidebar is the only place a learner can switch light/dark from here. It drives the icon
  // through a ref rather than hover, because a touch device never fires hover.
  const { mounted, toggleTheme } = useTheme()
  const sunMoonRef = useRef<IconHandle>(null)

  function handleThemeToggle(): void {
    const reduced =
      typeof window !== 'undefined' &&
      window.matchMedia('(prefers-reduced-motion: reduce)').matches
    if (!reduced) sunMoonRef.current?.startAnimation()
    toggleTheme()
  }

  return (
    <Sidebar collapsible="icon" data-testid="profile-sidebar">
      <SidebarHeader>
        <div className="flex items-center gap-1">
          {/* Collapsed rail: the brand mark doubles as the expand toggle, so it has to *look*
              like a control. Expanded, `λ Mathua` is a link to the hub — the conventional
              meaning of a logo — and a separate chevron collapses it. Collapsed, that same λ
              becomes a button, which means the identical glyph is a destination in one state and
              an action in the other. `aria-label` and `title` were the only things saying so, and
              both need a hover, so the role was communicated to nobody until they tapped it.

              The chevron mirrors the collapse chevron: expand and collapse look like the same kind
              of control pointing opposite ways, which is the convention and, unlike convention
              alone, it is visible before you press anything. */}
          <button
            type="button"
            onClick={toggleSidebar}
            aria-label="Expand sidebar"
            title="Expand sidebar"
            className="hidden size-11 shrink-0 flex-col items-center justify-center gap-0.5 font-mono text-sm text-mathua-blue hover:bg-sidebar-accent group-data-[collapsible=icon]:flex"
          >
            <span aria-hidden="true">λ</span>
            <ChevronsRight size={10} className="text-mathua-muted" />
          </button>
          <SidebarMenu className="min-w-0 flex-1 group-data-[collapsible=icon]:hidden">
            <SidebarMenuItem>
              <SidebarMenuButton size="lg" asChild tooltip="Mathua hub">
                <Link href="/profile">
                  <span className="flex size-6 shrink-0 items-center justify-center font-mono text-sm text-mathua-blue" aria-hidden="true">λ</span>
                  <span className="font-mono text-sm text-mathua-blue group-data-[collapsible=icon]:hidden">Mathua</span>
                </Link>
              </SidebarMenuButton>
            </SidebarMenuItem>
          </SidebarMenu>
          {/* Expanded: chevron collapses to icon rail */}
          <button
            type="button"
            onClick={toggleSidebar}
            aria-label="Collapse sidebar"
            title="Collapse sidebar"
            className="flex size-8 shrink-0 items-center justify-center text-mathua-muted hover:bg-sidebar-accent hover:text-mathua-blue group-data-[collapsible=icon]:hidden"
          >
            <ChevronsLeft />
          </button>
        </div>
      </SidebarHeader>

      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupLabel>Navigate</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              {/* Admin last, and only for an administrator. Last because it is a different
                  surface rather than another destination in this one; conditional because
                  offering it to everyone would be an invitation most of them cannot accept.
                  The server refuses either way — see AppSidebarProps.isAdmin. */}
              {(isAdmin ? [...NAV_ITEMS, ADMIN_ITEM] : NAV_ITEMS).map((item) => {
                const active = pathname === item.href || pathname.startsWith(item.href + '/')
                const Icon = item.icon
                return (
                  <SidebarMenuItem key={item.href}>
                    <SidebarMenuButton asChild isActive={active} tooltip={item.label}>
                      <Link href={item.href}>
                        <Icon className="shrink-0" />
                        <span className="group-data-[collapsible=icon]:hidden">{item.label}</span>
                      </Link>
                    </SidebarMenuButton>
                    {item.href === '/learn' && dueReviews > 0 && (
                      <SidebarMenuBadge>{dueReviews}</SidebarMenuBadge>
                    )}
                  </SidebarMenuItem>
                )
              })}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
        <SidebarGroup>
          <SidebarGroupLabel>Resources</SidebarGroupLabel>
          <SidebarGroupContent>
            <SidebarMenu>
              {RESOURCE_ITEMS.map((item) => {
                const active = pathname === item.href || pathname.startsWith(item.href + '/')
                const Icon = item.icon
                return (
                  <SidebarMenuItem key={item.href}>
                    <SidebarMenuButton asChild isActive={active} tooltip={item.label}>
                      <Link href={item.href}>
                        <Icon className="shrink-0" />
                        <span className="group-data-[collapsible=icon]:hidden">{item.label}</span>
                      </Link>
                    </SidebarMenuButton>
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
          {mounted && (
            <SidebarMenuItem>
              {/* aria-label, not just tooltip: the label span is display:none in the collapsed
                  rail and the tooltip is a description, so without this the icon-only button has
                  no accessible name in the state most desktop users leave the sidebar in. */}
              <SidebarMenuButton
                tooltip="Toggle theme"
                aria-label="Toggle theme"
                onClick={handleThemeToggle}
              >
                <SunMoon ref={sunMoonRef} size={16} isAnimated={false} aria-hidden="true" />
                <span className="group-data-[collapsible=icon]:hidden">Toggle theme</span>
              </SidebarMenuButton>
            </SidebarMenuItem>
          )}
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
