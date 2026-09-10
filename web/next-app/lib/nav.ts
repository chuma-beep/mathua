export interface NavLink {
  label: string
  href: string
}

// HeaderLink is the Header's prop name for the same shape.
export type HeaderLink = NavLink

// Single source of truth for site navigation (Header + BottomTabs).
//
// Bottom tabs are the primary mobile nav (Home, Study, Start, Profile,
// Graph). The desktop header row shows the full link set. The mobile
// compass menu shows ONLY overflow links not already in the tabs, so no
// destination appears twice on one screen.

export const DISPLAY_LINKS: NavLink[] = [
  { label: 'Study', href: '/study' },
  { label: 'Leaderboard', href: '/leaderboard' },
  { label: 'Graph', href: '/graph' },
]

export const LOGGED_IN_LINKS: NavLink[] = [
  { label: 'Profile', href: '/profile' },
  { label: 'Start', href: '/session' },
  { label: 'Settings', href: '/settings' },
]

export const LOGGED_OUT_LINKS: NavLink[] = [{ label: 'Login', href: '/login' }]

// Hrefs covered by the bottom tabs — never repeated in the compass menu.
export const TAB_HREFS: ReadonlySet<string> = new Set([
  '/',
  '/study',
  '/session',
  '/profile',
  '/graph',
])

// Full desktop header row (unchanged order).
export function desktopLinks(loggedIn: boolean): NavLink[] {
  return [...DISPLAY_LINKS, ...(loggedIn ? LOGGED_IN_LINKS : LOGGED_OUT_LINKS)]
}

// Mobile compass overflow: destinations missing from the bottom tabs.
export function overflowLinks(loggedIn: boolean): NavLink[] {
  return desktopLinks(loggedIn).filter(l => !TAB_HREFS.has(l.href))
}
