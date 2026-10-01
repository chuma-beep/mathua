export interface NavLink {
  label: string
  href: string
}

// HeaderLink is the Header's prop name for the same shape.
export type HeaderLink = NavLink

// Single source of truth for site navigation (Header + BottomTabs).
//
// Order carries meaning, and it is the same meaning on every surface: Learn
// leads, then Study.
//
// Learn is the only surface that asks a question, grades an answer and explains
// it. Study explains and never grades. Putting reference above doing made the
// product read backwards — a learner opening the app was sent to browse a
// reference library before being offered anything to solve.
//
// Study stays a tab regardless. It is the discovery surface for "what even is
// this concept", which is a real job that Learn has no substitute for: Learn
// has no browse, no search and no domain list.
//
// Bottom tabs are the primary mobile nav and stay at five. The mobile compass
// menu shows ONLY overflow links not already in the tabs, so no destination
// appears twice on one screen.

// LEARN_LINKS is the doing group, and it is shown to everyone — a guest can
// answer questions, they just cannot have them remembered.
export const LEARN_LINKS: NavLink[] = [{ label: 'Learn', href: '/learn' }]

// REVIEW_LINKS is a Learn sub-mode: it grades, awards XP and reschedules
// memory, so it belongs beside Learn rather than looking like its own product.
// Auth-only, because /api/reviews/* sits behind authMiddleware — offering it to
// a guest would be a link to a 401.
export const REVIEW_LINKS: NavLink[] = [{ label: 'Review', href: '/review' }]

// DISPLAY_LINKS is the reference group.
export const DISPLAY_LINKS: NavLink[] = [
  { label: 'Study', href: '/study' },
  { label: 'Graph', href: '/graph' },
  { label: 'Leaderboard', href: '/leaderboard' },
]

export const LOGGED_IN_LINKS: NavLink[] = [
  { label: 'Profile', href: '/profile' },
  { label: 'Settings', href: '/settings' },
]

export const LOGGED_OUT_LINKS: NavLink[] = [{ label: 'Login', href: '/login' }]

// Hrefs covered by the bottom tabs — never repeated in the compass menu.
// Review is deliberately absent: the tab bar stays at five, and Review reaches
// the header row and the compass overflow instead.
export const TAB_HREFS: ReadonlySet<string> = new Set([
  '/',
  '/learn',
  '/study',
  '/profile',
  '/graph',
])

// Full desktop header row: Learn first, then Study, then the account links.
export function desktopLinks(loggedIn: boolean): NavLink[] {
  return [
    ...LEARN_LINKS,
    ...(loggedIn ? REVIEW_LINKS : []),
    ...DISPLAY_LINKS,
    ...(loggedIn ? LOGGED_IN_LINKS : LOGGED_OUT_LINKS),
  ]
}

// Mobile compass overflow: destinations missing from the bottom tabs.
export function overflowLinks(loggedIn: boolean): NavLink[] {
  return desktopLinks(loggedIn).filter(l => !TAB_HREFS.has(l.href))
}
