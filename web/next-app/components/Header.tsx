'use client'

import Link from 'next/link'
import { useState, useEffect, useRef } from 'react'
import { useRouter, usePathname } from 'next/navigation'
import { Compass, SunMoon, type IconHandle } from './icons'
import { useTheme } from '../hooks/useTheme'
import { useAuthState } from '../hooks/useAuthState'
import { useChrome } from '../hooks/useChrome'
import { useInertWhen } from '../hooks/useInertWhen'
import { registerChromePart } from '../lib/chrome'
import { canOpenAdminConsole, signOut } from '../lib/auth'
import { resolveAvatar } from '../lib/dicebear'
import Avatar from './Avatar'
import ReportUpdatesBadge from './ReportUpdatesBadge'
import { getSettings } from '../lib/api'
import { desktopLinks, overflowLinks, type HeaderLink } from '../lib/nav'

interface HeaderProps {
  links?: HeaderLink[]
}

// The destinations that stay on the top row. Everything else moves into More,
// because a signed-in public page was rendering twelve links and the row
// overlapped itself before the beta sign.
const PRIMARY_HREFS = new Set(['/learn', '/review', '/graph', '/domains', '/leaderboard', '/login'])

export default function Header({ links }: HeaderProps) {
  const { mounted, toggleTheme } = useTheme()
  const { loggedIn, user } = useAuthState()
  const { visible: chromeVisible } = useChrome()
  const headerRef = useRef<HTMLElement>(null)
  useInertWhen(headerRef, !chromeVisible)
  // Tell the chrome store a header exists, so the tap control knows there is
  // something to hide even on a page with no tab bar (the landing page).
  useEffect(() => registerChromePart('header'), [])
  const router = useRouter()
  const pathname = usePathname()
  const [open, setOpen] = useState(false)
  const [navOpen, setNavOpen] = useState(false)
  const [moreOpen, setMoreOpen] = useState(false)
  const [avatarUrl, setAvatarUrl] = useState<string | undefined>(undefined)
  const [avatarPreset, setAvatarPreset] = useState<number | null>(null)
  // Leaderboard navigation follows the privacy setting the server reads to
  // anonymise the rows. Default true so a slow settings call never removes a
  // destination; only an explicit `hidden` does.
  const [showLeaderboards, setShowLeaderboards] = useState(true)
  // The trigger and the panel are separate DOM subtrees now, so outside-click has to check both:
  // the panel is no longer a descendant of the trigger, and treating it as outside would close
  // the menu on the very mousedown that is about to activate a menu item.
  const menuRef = useRef<HTMLDivElement>(null)
  const menuPanelRef = useRef<HTMLDivElement>(null)
  const navMenuRef = useRef<HTMLDivElement>(null)
  const navPanelRef = useRef<HTMLDivElement>(null)
  const moreMenuRef = useRef<HTMLDivElement>(null)
  const morePanelRef = useRef<HTMLDivElement>(null)
  const compassRef = useRef<IconHandle>(null)
  const sunMoonRef = useRef<IconHandle>(null)

  function prefersReducedMotion(): boolean {
    return (
      typeof window !== 'undefined' &&
      window.matchMedia('(prefers-reduced-motion: reduce)').matches
    )
  }

  function toggleNav(): void {
    const next = !navOpen
    setNavOpen(next)
    // Alternate drive direction so the needle moves on every click: forward
    // spin on open, reverse spin on close (hover rarely fires on touch
    // screens; skipped entirely for reduced motion).
    if (prefersReducedMotion()) return
    if (next) compassRef.current?.startAnimation()
    else compassRef.current?.stopAnimation()
  }

  function handleThemeToggle(): void {
    if (!prefersReducedMotion()) sunMoonRef.current?.startAnimation()
    toggleTheme()
  }

  useEffect(() => {
    if (!loggedIn || !user) {
      setAvatarUrl(undefined)
      setAvatarPreset(null)
      setShowLeaderboards(true)
      return
    }
    // Dead-token cleanup is lazy: any strict API call below that answers
    // 401 clears the token via authedFetch (lib/auth.ts), which flips this
    // UI to logged-out through the auth-changed event. No explicit
    // revalidation here, the e2e auth spec renders the header with a
    // stored token while all /api/** routes are mocked, and must keep
    // showing the authenticated links.
    getSettings().then(s => {
      // Runs on every [loggedIn, user] change, including the auth-changed
      // broadcast that Settings fires after avatar saves, so the header
      // picture tracks the shared resolver with no extra event type.
      const resolved = resolveAvatar(user, s)
      setAvatarPreset(resolved.preset ?? null)
      setAvatarUrl(resolved.url)
      setShowLeaderboards(s.leaderboard_privacy !== 'hidden')
    }).catch(() => {
      setAvatarPreset(null)
      setAvatarUrl(user.avatar_url)
      // A failed read cannot be evidence that leaderboards were turned off.
      setShowLeaderboards(true)
    })
  }, [loggedIn, user])

  useEffect(() => {
    function onClick(e: MouseEvent) {
      const t = e.target as Node
      if (!menuRef.current?.contains(t) && !menuPanelRef.current?.contains(t)) setOpen(false)
      // The panel must be named here as well as the trigger. It hangs off <header> rather
      // than off the trigger's wrapper, so a mousedown on one of its own links is not inside
      // navMenuRef and closing on it unmounted the link before the click could activate it,
      // which is why every destination in the mobile menu was dead: the items were in the DOM
      // and nothing could reach them.
      if (navMenuRef.current && !navMenuRef.current.contains(t) && !navPanelRef.current?.contains(t)) setNavOpen(false)
      if (moreMenuRef.current && !moreMenuRef.current.contains(t) && !morePanelRef.current?.contains(t)) setMoreOpen(false)
    }
    function onKey(e: KeyboardEvent) {
      if (e.key === 'Escape') {
        setOpen(false)
        setNavOpen(false)
        setMoreOpen(false)
      }
    }
    document.addEventListener('mousedown', onClick)
    document.addEventListener('keydown', onKey)
    return () => {
      document.removeEventListener('mousedown', onClick)
      document.removeEventListener('keydown', onKey)
    }
  }, [])

  // Desktop row: full link set. Mobile compass: overflow only (destinations
  // missing from the bottom tabs) so nothing appears twice on one screen.
  // A custom links list (e.g. docs section menu) is shown in full.
  //
  // Leaderboard is filtered from every list, custom ones included, when the
  // learner turned it off. Only `hidden` does this: `anonymous` keeps the
  // destination because the board is still useful without a name on it.
  const showLink = (link: HeaderLink) => showLeaderboards || link.href !== '/leaderboard'
  const displayLinks = (links || desktopLinks(loggedIn)).filter(showLink)
  // The default signed-in row was twelve destinations wide and overlapped itself
  // on a laptop. Keep the learning journey in the row and put the reference and
  // account destinations behind one More menu. A custom links list (docs
  // section menu) is still shown in full because it is already short.
  const primaryLinks = links ? displayLinks : displayLinks.filter((link) => PRIMARY_HREFS.has(link.href))
  const moreLinks = links ? [] : displayLinks.filter((link) => !PRIMARY_HREFS.has(link.href))
  const moreActive = moreLinks.some(
    (link) => pathname === link.href || (pathname ?? '').startsWith(link.href + '/'),
  )
  // The Admin entry reaches mobile through the compass overflow. The desktop rail on
  // /profile has always carried it; a phone had no route to administration at all,
  // because `/profile` is one tap from the tab bar and nothing on any other page
  // offered it. The role comes from `/api/auth/me` via the session, and offering it
  // grants nothing, the admin routes re-check it server-side.
  const compassLinks = links
    ? displayLinks
    : overflowLinks(loggedIn).filter(showLink)

  return (
    // The header hides by sliding up, not by collapsing. Collapsing its height
    // removed ~149px from the top of the document (69px of bar plus the 80px
    // --chrome-top padding), and because the main content sets
    // `overflow-anchor: none` the browser did not compensate: everything under
    // the finger jumped up as the chrome went away. That is the bug a reader
    // reported on iOS Safari. A transform is not in flow, so the document keeps
    // its height and nothing moves; the sticky box keeps its place at the top
    // of the page, and translating it only uncovers what was already behind it.
    <header
      ref={headerRef}
      data-chrome-part="header"
      className={`sticky top-0 z-50 backdrop-blur transition-transform duration-300 ease-out motion-reduce:transition-none ${
        chromeVisible ? 'translate-y-0' : '-translate-y-full'
      }`}
    >
      <div className="bg-mathua-bg">
        {/* Full-bleed rule, constrained content.

            The bar's rule is the header's own edge: it should reach both sides of the window the
            way the footer's does. It used to sit on this row, which also carried `max-w-container`,
            so on a screen wider than 1100px the rule stopped short and the bar read as a floating
            strip with the page showing either side of it. The footer already had it right, 
            `border-t` on the full-width `<footer>`, unconstrained inner content, so the header is
            now the same shape.

            The inner row stays constrained rather than going full width, because the links need to
            line up with the page content below, and the page content is capped at 1100px. Widening
            the row alone would put the logo far from the top-left corner it is supposed to anchor.
            `mx-auto` re-centres the capped row against the wider parent. */}
        <div className="border-b border-mathua-border">
        <div className="max-w-container mx-auto flex items-center justify-between px-4 md:px-6 py-3">
          <div className="flex items-center gap-4 md:gap-6 min-w-0">
            <span className="flex items-center gap-2 shrink-0">
              <Link
                href={loggedIn ? '/profile' : '/'}
                className="link-underline font-mono text-sm text-mathua-blue whitespace-nowrap"
              >
                λ Mathua
              </Link>
              {/* The sign says where the app is, not what the learner is. It is small, next to
                  the wordmark on every page, and it will come down when it stops being true, 
                  rather than being a decoration nobody can date. */}
              <span
                className="border border-mathua-border px-1.5 py-0.5 font-mono text-[9px] uppercase tracking-wider text-mathua-muted leading-none"
                title="Mathua is in beta: the curriculum, the app and its fixes are still moving."
              >
                beta
              </span>
            </span>

            <nav className="hidden md:flex items-center gap-2">
              {primaryLinks.map((link, i) => (
                <span key={link.href} className="flex items-center gap-2">
                  {i > 0 && (
                    <span className="font-mono text-xs text-mathua-blue select-none">|</span>
                  )}
                  <Link
                    href={link.href}
                    className="link-underline font-mono text-xs text-mathua-muted hover:text-mathua-blue whitespace-nowrap transition-colors"
                  >
                    {link.label}
                  </Link>
                </span>
              ))}
              {links ? (
                <span className="flex items-center gap-2">
                  <span className="font-mono text-xs text-mathua-blue select-none">|</span>
                  <a
                    href="https://discord.gg/qXPeQ2jys"
                    target="_blank"
                    rel="noreferrer"
                    className="link-underline font-mono text-xs text-mathua-muted hover:text-mathua-blue whitespace-nowrap transition-colors"
                  >
                    Discord
                  </a>
                </span>
              ) : (
                <div className="relative flex items-center gap-2" ref={moreMenuRef}>
                  <span className="font-mono text-xs text-mathua-blue select-none">|</span>
                  <button
                    type="button"
                    onClick={() => setMoreOpen(o => !o)}
                    aria-label="More navigation"
                    aria-expanded={moreOpen}
                    aria-haspopup="menu"
                    className={`link-underline font-mono text-xs whitespace-nowrap transition-colors ${
                      moreActive ? 'text-mathua-blue' : 'text-mathua-muted hover:text-mathua-blue'
                    }`}
                  >
                    More
                  </button>
                  {moreOpen && chromeVisible && (
                    <div
                      role="menu"
                      aria-label="More navigation"
                      ref={morePanelRef}
                      className="absolute right-0 top-full mt-2 min-w-[200px] bg-mathua-surface border border-mathua-border shadow-lg py-1 z-50"
                    >
                      {moreLinks.map(link => {
                        const active = pathname === link.href || (pathname ?? '').startsWith(link.href + '/')
                        return (
                          <Link
                            key={link.href}
                            href={link.href}
                            role="menuitem"
                            aria-current={active ? 'page' : undefined}
                            onClick={() => setMoreOpen(false)}
                            className={`flex items-center min-h-[40px] px-4 font-mono text-xs whitespace-nowrap transition-colors ${
                              active ? 'text-mathua-blue' : 'text-mathua-muted hover:text-mathua-blue hover:bg-mathua-surface-elevated'
                            }`}
                          >
                            {link.label}
                          </Link>
                        )
                      })}
                      <a
                        href="https://discord.gg/qXPeQ2jys"
                        role="menuitem"
                        target="_blank"
                        rel="noreferrer"
                        onClick={() => setMoreOpen(false)}
                        className="flex items-center min-h-[40px] px-4 font-mono text-xs whitespace-nowrap text-mathua-muted transition-colors hover:bg-mathua-surface-elevated hover:text-mathua-blue"
                      >
                        Discord
                      </a>
                    </div>
                  )}
                </div>
              )}
            </nav>
          </div>

          <div className="flex items-center gap-2 sm:gap-3 shrink-0">
            {mounted && (
              <button
                type="button"
                onClick={handleThemeToggle}
                aria-label="Toggle theme"
                className="flex items-center justify-center min-h-[44px] min-w-[44px] text-mathua-muted bg-transparent border-none cursor-pointer hover:text-mathua-blue transition-colors"
              >
                <SunMoon ref={sunMoonRef} size={15} isAnimated={false} />
              </button>
            )}
            {/* Mobile nav: bare compass mark. Overflow only
                (see lib/nav.ts), destinations missing from the bottom tabs. */}
            <div className="md:hidden" ref={navMenuRef}>
                <button
                  type="button"
                  onClick={toggleNav}
                  aria-label="Open navigation menu"
                  aria-expanded={navOpen}
                  aria-haspopup="menu"
                  className="flex items-center justify-center min-h-[44px] min-w-[44px] text-mathua-muted bg-transparent border-none cursor-pointer hover:text-mathua-blue transition-colors"
                >
                  <Compass ref={compassRef} size={15} isAnimated={false} />
                </button>
            </div>
            <ReportUpdatesBadge enabled={loggedIn} />
            {loggedIn && user ? (
              <div className="relative" ref={menuRef}>
                <button
                  type="button"
                  onClick={() => setOpen(o => !o)}
                  aria-label="Open profile menu"
                  aria-expanded={open}
                  aria-haspopup="menu"
                  className="rounded-full p-0.5 border border-transparent hover:border-mathua-border focus:outline-none focus:border-mathua-blue transition-colors"
                >
                  <Avatar
                    seed={user.student_id}
                    name={user.name}
                    size={32}
                    url={avatarPreset !== null ? undefined : avatarUrl}
                    preset={avatarPreset ?? undefined}
                  />
                </button>
              </div>
            ) : null}
          </div>
        </div>
        </div>
      </div>
      {/* The nav panel hangs off <header>, not off the row. It used to sit inside a
          `min-h-0 overflow-hidden` clipper, the grid trick that collapsed the header on scroll,
          and that clipped absolutely positioned children so the panel painted *behind* the page.
          The header slides on a transform now and needs no clipper, but the panel stays where it
          is: this is the structure the open menus were tested against.

          Gated on chromeVisible: the trigger slides away with the header, so a panel whose button
          is off screen is a floating menu with no visible way to close it. */}
{navOpen && chromeVisible && (
                  <div className="absolute inset-y-0 left-1/2 w-full max-w-container -translate-x-1/2 pointer-events-none">
                    <div
                      role="menu"
                      aria-label="Site navigation"
                      ref={navPanelRef}
                      className="pointer-events-auto absolute right-4 md:right-6 top-full min-w-[200px] bg-mathua-surface border border-mathua-border shadow-lg py-1 z-50 md:hidden"
                    >
                      {compassLinks.map(link => {
                        const active = pathname === link.href || (pathname ?? '').startsWith(link.href + '/')
                        return (
                          <Link
                            key={link.href}
                            href={link.href}
                            role="menuitem"
                            aria-current={active ? 'page' : undefined}
                            onClick={() => setNavOpen(false)}
                            className={`flex items-center min-h-[44px] px-4 font-mono text-xs whitespace-nowrap transition-colors ${
                              active ? 'text-mathua-blue' : 'text-mathua-muted hover:text-mathua-blue hover:bg-mathua-surface-elevated'
                            }`}
                          >
                            {link.label}
                          </Link>
                        )
                      })}
                      <a
                        href="https://discord.gg/qXPeQ2jys"
                        role="menuitem"
                        target="_blank"
                        rel="noreferrer"
                        onClick={() => setNavOpen(false)}
                        className="flex items-center min-h-[44px] px-4 font-mono text-xs whitespace-nowrap text-mathua-muted hover:text-mathua-blue hover:bg-mathua-surface-elevated transition-colors"
                      >
                        Discord
                      </a>
                    </div>
                  </div>
                  )}


      {/* The profile menu hangs off <header> for the same reason the nav panel does: it was
          built when the header collapsed inside a `min-h-0 overflow-hidden` clipper, and that
          clipped absolutely positioned children, so the menu painted *behind* the page.

          Anchored to the header's own padding rather than to the button, because the button is
          not an ancestor and `right-0` would have nothing to resolve against. The avatar sits at
          the end of the same row the header's padding insets, so the two agree.

          Gated on chromeVisible for the same reason as the nav panel: the trigger slides away
          with the header, so a menu whose button is off screen could not be closed. */}
      {open && chromeVisible && (
        <>
        {/* The panel's wrapper is itself absolutely positioned: <header> is full width, so a panel
            measured from it lands at the window edge while its trigger sits at the end of the
            capped row. This box is the same width as that row and centred, so `right-6` lands next
            to the avatar. `inset-y-0` keeps it as tall as the header, which is what makes
            `top-full` mean the bottom of the bar; it is a positioning context for the same reason.
            `pointer-events-none` with the panel re-enabling its own keeps the rest of the bar
            clickable underneath an open menu. */}
        <div className="absolute inset-y-0 left-1/2 w-full max-w-container -translate-x-1/2 pointer-events-none">
        <div
          ref={menuPanelRef}
          role="menu"
          aria-label="Profile"
          className="pointer-events-auto absolute right-4 md:right-6 top-full min-w-[160px] bg-mathua-surface border border-mathua-border shadow-lg py-1 z-50"
        >
          <div className="px-3 py-2 border-b border-mathua-border">
            <div className="font-mono text-xs text-mathua-primary truncate">{user?.name}</div>
            {user?.username && <div className="font-mono text-[10px] text-mathua-muted truncate">@{user.username}</div>}
          </div>
          <button
            type="button"
            role="menuitem"
            onClick={() => {
              setOpen(false)
              router.push('/my-reports')
            }}
            className="w-full text-left px-3 py-2 font-mono text-xs text-mathua-muted hover:text-mathua-blue hover:bg-mathua-surface-elevated transition-colors"
          >
            Notifications
          </button>
          <button
            type="button"
            role="menuitem"
            onClick={() => {
              setOpen(false)
              router.push('/profile')
            }}
            className="w-full text-left px-3 py-2 font-mono text-xs text-mathua-muted hover:text-mathua-blue hover:bg-mathua-surface-elevated transition-colors"
          >
            Profile
          </button>
          <button
            type="button"
            role="menuitem"
            onClick={() => {
              setOpen(false)
              router.push('/contact')
            }}
            className="w-full text-left px-3 py-2 font-mono text-xs text-mathua-muted hover:text-mathua-blue hover:bg-mathua-surface-elevated transition-colors"
          >
            Contact
          </button>
          {canOpenAdminConsole(user?.role) && (
            <>
              <div className="my-1 border-t border-mathua-border" />
              <a
                href="/admin"
                target="_blank"
                rel="noreferrer"
                role="menuitem"
                onClick={() => setOpen(false)}
                className="block px-3 py-2 font-mono text-xs text-mathua-muted hover:text-mathua-blue hover:bg-mathua-surface-elevated transition-colors"
              >
                Admin console
              </a>
            </>
          )}
          <button
            type="button"
            role="menuitem"
            onClick={() => {
              setOpen(false)
              signOut()
              router.push('/login')
            }}
            className="w-full text-left px-3 py-2 font-mono text-xs text-mathua-muted hover:text-mathua-red hover:bg-mathua-surface-elevated transition-colors"
          >
            Sign out
          </button>
        </div>
        </div>
        </>
      )}
    </header>
  )
}
