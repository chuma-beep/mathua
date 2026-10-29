'use client'

import Link from 'next/link'
import { usePathname, useRouter } from 'next/navigation'
import { useRef, useState } from 'react'
import {
  Bell,
  BookOpen,
  ChartNoAxesColumn,
  ChevronDown,
  ChevronsLeft,
  ChevronsRight,
  FileText,
  Handshake,
  House,
  Layers,
  LogOut,
  Mail,
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
  SidebarMenuSub,
  SidebarMenuSubButton,
  SidebarMenuSubItem,
  SidebarMenuBadge,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarRail,
  useSidebar,
} from '@/components/ui/sidebar'
import Avatar from './Avatar'
import { useTheme } from '../hooks/useTheme'
import { signOut } from '../lib/auth'

// Navigate group: destinations only. Never Diagnostic / Quiz / Review, 
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
  { label: 'Reports', href: '/reports', icon: FileText },
  { label: 'Settings', href: '/settings', icon: Settings },
]

// Resources group: docs, glossary, contributing, creator's note, how it works.
  // Informational pages only, never Diagnostic / Quiz / Review (no CTA
  // duplication). The glossary is a standalone reference with no progress,
  // XP or gating, so it sits here rather than beside Learn.
  const RESOURCE_ITEMS: SidebarNavItem[] = [
    { label: 'Docs', href: '/docs', icon: FileText },
    { label: 'Glossary', href: '/glossary', icon: BookOpen },
    { label: 'Contribute', href: '/docs/contributing', icon: Handshake },
    { label: "Creator's note", href: '/note', icon: Pencil },
    { label: 'How it works', href: '/how-it-works', icon: House },
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
   * A convenience and nothing more, the server reads the role on every administrative request,
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
  /**
   * Whether to offer the Leaderboard entry.
   *
   * Purely presentational: the board is still reachable by URL and the server
   * still anonymises the row. It is here because a learner who turned
   * leaderboards off asked for the surface to go away, not only for their name.
   */
  showLeaderboards?: boolean
  notificationUnread?: number
  /** Open reports in the moderation queue. Only used for the Admin entry. */
  adminReports?: number
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
  showLeaderboards = true,
  notificationUnread = 0,
  adminReports = 0,
}: AppSidebarProps) {
  const pathname = usePathname()
  // Contact is one control, not two destinations: the page and the Discord
  // server answer different needs, so the choice belongs after the click.
  const [contactOpen, setContactOpen] = useState(false)
  // The five reference links are not part of the daily learning loop, so they
  // stay folded away. Eight primary destinations plus one disclosure is a
  // shorter rail than thirteen links, and every reference page remains one
  // click away.
  const [resourcesOpen, setResourcesOpen] = useState(false)
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
              like a control. Expanded, `λ Mathua` is a link to the hub, the conventional
              meaning of a logo and a separate chevron collapses it. Collapsed, that same λ
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
                  The server refuses either way, see AppSidebarProps.isAdmin. */}
              {(isAdmin ? [...NAV_ITEMS, ADMIN_ITEM] : NAV_ITEMS)
                .filter((item) => showLeaderboards || item.href !== '/leaderboard')
                .map((item) => {
                  const active = pathname === item.href || pathname.startsWith(item.href + '/')
                  const Icon = item.icon
                  return (
                    <SidebarMenuItem key={item.href}>
                      <SidebarMenuButton asChild isActive={active} tooltip={item.label}>
                        {item.href.startsWith('http') ? (
                          <a href={item.href} target="_blank" rel="noreferrer">
                            <Icon className="shrink-0" />
                            <span className="group-data-[collapsible=icon]:hidden">{item.label}</span>
                          </a>
                        ) : (
                          <Link href={item.href}>
                            <Icon className="shrink-0" />
                            <span className="group-data-[collapsible=icon]:hidden">{item.label}</span>
                            {item.href === '/admin' && adminReports > 0 && (
                              <>
                                <span className="ml-auto h-2 w-2 shrink-0 rounded-full bg-mathua-red" aria-hidden="true" />
                                <span className="sr-only">
                                  {adminReports === 1
                                    ? '1 report needs attention'
                                    : `${adminReports} reports need attention`}
                                </span>
                              </>
                            )}
                          </Link>
                        )}
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
          <SidebarGroupContent>
            <SidebarMenu>
              <SidebarMenuItem>
                <SidebarMenuButton
                  type="button"
                  tooltip="Resources"
                  aria-expanded={resourcesOpen}
                  onClick={() => setResourcesOpen((open) => !open)}
                >
                  <BookOpen className="shrink-0" />
                  <span className="group-data-[collapsible=icon]:hidden">Resources</span>
                  <ChevronDown
                    className={`ml-auto shrink-0 transition-transform group-data-[collapsible=icon]:hidden ${resourcesOpen ? 'rotate-180' : ''}`}
                  />
                </SidebarMenuButton>
                {resourcesOpen && (
                  <SidebarMenuSub>
                    {RESOURCE_ITEMS.map((item) => {
                      const active = pathname === item.href || pathname.startsWith(item.href + '/')
                      const Icon = item.icon
                      return (
                        <SidebarMenuSubItem key={item.href}>
                          <SidebarMenuSubButton asChild isActive={active}>
                            <Link href={item.href}>
                              <Icon className="shrink-0" />
                              <span className="group-data-[collapsible=icon]:hidden">{item.label}</span>
                            </Link>
                          </SidebarMenuSubButton>
                        </SidebarMenuSubItem>
                      )
                    })}
                    <SidebarMenuSubItem>
                      <SidebarMenuSubButton
                        type="button"
                        isActive={pathname === '/contact'}
                        aria-expanded={contactOpen}
                        onClick={() => setContactOpen((open) => !open)}
                      >
                        <Mail className="shrink-0" />
                        <span className="group-data-[collapsible=icon]:hidden">Contact admins</span>
                        <ChevronDown
                          className={`ml-auto shrink-0 transition-transform group-data-[collapsible=icon]:hidden ${contactOpen ? 'rotate-180' : ''}`}
                        />
                      </SidebarMenuSubButton>
                      {contactOpen && (
                        <SidebarMenuSub>
                          <SidebarMenuSubItem>
                            <SidebarMenuSubButton asChild>
                              <Link href="/contact">
                                <span>Contact page</span>
                              </Link>
                            </SidebarMenuSubButton>
                          </SidebarMenuSubItem>
                          <SidebarMenuSubItem>
                            <SidebarMenuSubButton asChild>
                              <a href="https://discord.gg/qXPeQ2jys" target="_blank" rel="noreferrer">
                                <span>Discord server</span>
                              </a>
                            </SidebarMenuSubButton>
                          </SidebarMenuSubItem>
                        </SidebarMenuSub>
                      )}
                    </SidebarMenuSubItem>
                  </SidebarMenuSub>
                )}
              </SidebarMenuItem>
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
        {/* Notifications sits below the reference library because it is a
            mailbox, not a daily learning destination. It stays outside the
            folded Resources disclosure so an unread report is always one
            click away. */}
        <SidebarGroup>
          <SidebarGroupContent>
            <SidebarMenu>
              <SidebarMenuItem>
                <SidebarMenuButton
                  asChild
                  isActive={pathname === '/my-reports'}
                  tooltip="Notifications"
                >
                  <Link href="/my-reports">
                    <Bell className="shrink-0" />
                    <span className="group-data-[collapsible=icon]:hidden">Notifications</span>
                    {notificationUnread > 0 && (
                      <>
                        <span className="ml-auto h-2 w-2 shrink-0 rounded-full bg-mathua-red" aria-hidden="true" />
                        <span className="sr-only">
                          {notificationUnread === 1
                            ? '1 unread notification'
                            : `${notificationUnread} unread notifications`}
                        </span>
                      </>
                    )}
                  </Link>
                </SidebarMenuButton>
              </SidebarMenuItem>
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
