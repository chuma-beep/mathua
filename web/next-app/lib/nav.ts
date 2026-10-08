export interface NavLink {
  label: string
  href: string
}

// HeaderLink is the Header's prop name for the same shape.
export type HeaderLink = NavLink

// Single source of truth for site navigation (Header + BottomTabs).
//
// Order carries meaning: Learn leads, and the orientation surfaces follow it.
//
// Learn is the only surface that asks a question, grades an answer and explains it. It used to
// lead into Study — a full reference library — because reading was treated as a way to learn.
// It is not: only demonstrated answers move a concept, so a learner who could read their way
// through the corpus had found a path around the only thing that counts. Reference material
// stayed; the route to it from navigation did not. It now lives inside `/learn`, scoped to the
// concept being learned (see components/ConceptReference).
//
// Bottom tabs are the primary mobile nav and stay at five. The mobile compass menu shows ONLY
// overflow links not already in the tabs, so no destination appears twice on one screen.

// LEARN_LINKS is the doing group, and it is shown to everyone — a guest can
// answer questions, they just cannot have them remembered.
export const LEARN_LINKS: NavLink[] = [{ label: 'Learn', href: '/learn' }]

// REVIEW_LINKS is a Learn sub-mode: it grades, awards XP and reschedules
// memory, so it belongs beside Learn rather than looking like its own product.
// Auth-only, because /api/reviews/* sits behind authMiddleware — offering it to
// a guest would be a link to a 401.
export const REVIEW_LINKS: NavLink[] = [{ label: 'Review', href: '/review' }]

// DISPLAY_LINKS is the orientation group: where things sit in the curriculum and where the
// learner stands in it. None of these grade an answer.
export const DISPLAY_LINKS: NavLink[] = [
  { label: 'Graph', href: '/graph' },
  // Not a bottom tab on every surface — the tab bar stays at five — so /domains reaches mobile
  // through the tab bar and the compass overflow, derived from this list minus TAB_HREFS.
  { label: 'Domains', href: '/domains' },
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
  '/domains',
  '/profile',
  '/graph',
])

// Full desktop header row: Learn first, then the orientation surfaces, then the account links.
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
