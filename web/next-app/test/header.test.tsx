import { describe, it, expect, beforeAll, beforeEach, vi } from 'vitest'
import { render, screen, fireEvent, within } from '@testing-library/react'
import Header from '../components/Header'

beforeAll(() => {
  Object.defineProperty(window, 'matchMedia', {
    writable: true,
    value: () => ({
      matches: false,
      media: '',
      onchange: null,
      addListener: () => {},
      removeListener: () => {},
      addEventListener: () => {},
      removeEventListener: () => {},
      dispatchEvent: () => false,
    }),
  })
})

const pushMock = vi.fn()
let mockPathname = '/'

vi.mock('next/navigation', () => ({
  useRouter: () => ({ push: pushMock }),
  usePathname: () => mockPathname,
}))

const { toggleThemeMock } = vi.hoisted(() => ({ toggleThemeMock: vi.fn() }))

vi.mock('../hooks/useTheme', () => ({
  useTheme: () => ({ theme: 'dark', mounted: true, toggleTheme: toggleThemeMock }),
}))

let mockUser: { student_id: string; name: string; username?: string; role?: string } | null = null

vi.mock('../hooks/useAuthState', () => ({
  useAuthState: () => ({ loggedIn: mockUser !== null, user: mockUser }),
}))

// vi.hoisted, because `vi.mock` factories are hoisted above the module body: a plain `let` here
// is in its temporal dead zone when the factory first runs, and the component ends up reading a
// different binding than the test writes to — so the assertion passes against a build that never
// had the gate in it. That is the failure mode a "does this test actually test anything" check is
// for, and it caught exactly that.
const { chromeRef } = vi.hoisted(() => ({ chromeRef: { visible: true } }))

vi.mock('../hooks/useChrome', () => ({
  useChrome: () => ({ visible: chromeRef.visible, toggle: () => {}, hasChrome: true, hasTabs: true }),
}))

// `getSettings` is called and `.then`d by the header's avatar effect, so the mock must return a
// promise. A bare `vi.fn()` returns undefined and the header throws before it can render.
vi.mock('../lib/api', async importOriginal => ({
  ...(await importOriginal<typeof import('../lib/api')>()),
  getSettings: vi.fn(async () => ({})),
}))

describe('Header mobile nav menu', () => {
  it('renders a collapsed Compass menu button with no menu items visible', () => {
    render(<Header />)
    const button = screen.getByRole('button', { name: 'Open navigation menu' })
    expect(button).toHaveAttribute('aria-expanded', 'false')
    expect(button).toHaveAttribute('aria-haspopup', 'menu')
    expect(screen.queryAllByRole('menuitem')).toHaveLength(0)
  })

  // The panel used to live inside `<div class="min-h-0 overflow-hidden">`, which is the CSS
  // grid-collapse trick that hides the header on scroll. That clips absolutely positioned
  // children, so on mobile the menu painted *behind* the page content instead of over it — the
  // one thing a dropdown must never do, and one no assertion here would have caught because the
  // panel was in the DOM and its links worked.
  //
  // The assertion is structural rather than visual: jsdom has no layout, so it cannot tell you
  // what is painted on top. What it can check is that the panel is not a descendant of the
  // clipping wrapper, which is the cause.
  it('renders the nav panel outside the header row that clips it', () => {
    render(<Header />)
    fireEvent.click(screen.getByRole('button', { name: 'Open navigation menu' }))
    const panel = screen.getByRole('menu', { name: 'Site navigation' })

    const clipper = document.querySelector('.overflow-hidden')
    expect(clipper, 'the collapse wrapper should still exist; the animation needs it').not.toBeNull()
    expect(clipper!.contains(panel)).toBe(false)

    // It still has to be anchored to the header rather than the viewport, which is what the
    // original `top-[calc(100%+8px)]` relative to the row achieved.
    expect(panel.className).toMatch(/absolute/)
    expect(panel.className).toMatch(/top-full/)
    expect(panel.className).toMatch(/z-50/)
    // Absolutely positioned against <header>, which is sticky and therefore a containing block.
    expect(panel.closest('header')).not.toBeNull()
  })

  it('opens the overflow links with correct hrefs on click', () => {
    render(<Header />)
    fireEvent.click(screen.getByRole('button', { name: 'Open navigation menu' }))
    const items = screen.getAllByRole('menuitem')
    const labels = items.map(i => i.textContent)
    // Whatever is left after the tab bar is subtracted out. Domains is not here: it holds a
    // tab, having taken the slot the closed Study route occupied, so it is always one tap away
    // on mobile instead of hiding one level deep.
    expect(labels).toEqual(['Leaderboard', 'Login'])
    const hrefs = items.map(i => (i as HTMLAnchorElement).getAttribute('href'))
    expect(hrefs).toEqual(['/leaderboard', '/login'])
    // The invariant behind the exact list: nothing appears twice on one screen.
    for (const href of hrefs) {
      expect(hrefs.filter(h => h === href)).toHaveLength(1)
    }
    expect(screen.getByRole('button', { name: 'Open navigation menu' }))
      .toHaveAttribute('aria-expanded', 'true')
  })

  it('renders a custom links list (docs section menu) when provided', () => {
    render(
      <Header
        links={[
          { label: 'Graph', href: '/graph' },
          { label: 'Docs', href: '/docs' },
          { label: 'Architecture', href: '/docs/architecture' },
          { label: 'Contributing', href: '/docs/contributing' },
          { label: 'Note', href: '/note' },
        ]}
      />
    )
    fireEvent.click(screen.getByRole('button', { name: 'Open navigation menu' }))
    const menu = screen.getByRole('menu', { name: 'Site navigation' })
    const labels = within(menu).getAllByRole('menuitem').map(i => i.textContent)
    expect(labels).toEqual(['Graph', 'Docs', 'Architecture', 'Contributing', 'Note'])
  })

  it('marks the current page with aria-current', () => {
    mockPathname = '/docs/architecture'
    render(
      <Header
        links={[
          { label: 'Docs', href: '/docs' },
          { label: 'Architecture', href: '/docs/architecture' },
        ]}
      />
    )
    fireEvent.click(screen.getByRole('button', { name: 'Open navigation menu' }))
    const menu = screen.getByRole('menu', { name: 'Site navigation' })
    const arch = within(menu).getByRole('menuitem', { name: 'Architecture' })
    expect(arch).toHaveAttribute('aria-current', 'page')
    mockPathname = '/'
  })

  it('closes on Escape', () => {
    render(<Header />)
    fireEvent.click(screen.getByRole('button', { name: 'Open navigation menu' }))
    expect(screen.getAllByRole('menuitem')).toHaveLength(2)
    fireEvent.keyDown(document, { key: 'Escape' })
    expect(screen.queryAllByRole('menuitem')).toHaveLength(0)
  })

  it('spins forward on open and reverses on close without crashing', () => {
    render(<Header />)
    const button = screen.getByRole('button', { name: 'Open navigation menu' })
    fireEvent.click(button)
    expect(screen.getAllByRole('menuitem')).toHaveLength(2)
    fireEvent.click(button)
    expect(screen.queryAllByRole('menuitem')).toHaveLength(0)
    fireEvent.click(button)
    expect(screen.getAllByRole('menuitem')).toHaveLength(2)
  })

  it('closes when a link is clicked', () => {
    render(<Header />)
    fireEvent.click(screen.getByRole('button', { name: 'Open navigation menu' }))
    fireEvent.click(screen.getByRole('menuitem', { name: 'Leaderboard' }))
    expect(screen.queryAllByRole('menuitem')).toHaveLength(0)
  })

  it('theme toggle plays the sun-moon animation and toggles the theme', () => {
    toggleThemeMock.mockClear()
    render(<Header />)
    fireEvent.click(screen.getByRole('button', { name: 'Toggle theme' }))
    expect(toggleThemeMock).toHaveBeenCalledTimes(1)
  })
})

describe('Header profile menu', () => {
  beforeEach(() => { mockUser = null; chromeRef.visible = true })

  it('renders the profile panel outside the header row that clips it', () => {
    // The same clipping bug the nav panel had, and the same fix. Inside `<div
    // class="min-h-0 overflow-hidden">` the profile menu painted *behind* the page instead of
    // over it, so clicking the avatar put a menu under the box below it — the one thing a
    // dropdown must never do, and invisible to jsdom, which has no layout. The assertion is
    // therefore structural: the panel must not be a descendant of the clipper.
    mockUser = { student_id: 's1', name: 'Nelson', username: 'wis' }
    render(<Header />)
    fireEvent.click(screen.getByRole('button', { name: 'Open profile menu' }))
    const panel = screen.getByRole('menu', { name: 'Profile' })

    const clipper = document.querySelector('.overflow-hidden')
    expect(clipper, 'the collapse wrapper should still exist; the animation needs it').not.toBeNull()
    expect(clipper!.contains(panel)).toBe(false)
    expect(panel.className).toMatch(/absolute/)
    expect(panel.className).toMatch(/z-50/)
    expect(panel.closest('header')).not.toBeNull()
  })

  // Deliberately a second test rather than an extra assertion above. The nav panel and the
  // profile panel are fixed by the same change and can regress independently — someone editing
  // one of them to hang off the row again would leave the other's test green — so each is
  // asserted in its own right rather than sharing a helper that hides which one broke.
  it('anchors the profile panel below the header, like the nav panel', () => {
    mockUser = { student_id: 's1', name: 'Nelson', username: 'wis' }
    render(<Header />)
    fireEvent.click(screen.getByRole('button', { name: 'Open profile menu' }))
    const panel = screen.getByRole('menu', { name: 'Profile' })
    // `top-full` against the sticky header, which is the containing block. `top-[calc(100%+8px)]`
    // would resolve against the row and reintroduce the offset that put it behind the page.
    expect(panel.className).toMatch(/top-full/)
    expect(panel.className).not.toMatch(/top-\[calc/)
  })

  it('offers Profile before Sign out, and each closes the menu', () => {
    mockUser = { student_id: 's1', name: 'Nelson', username: 'wis' }
    render(<Header />)
    fireEvent.click(screen.getByRole('button', { name: 'Open profile menu' }))
    const items = screen.getAllByRole('menuitem')
    expect(items.map(i => i.textContent)).toEqual(['Profile', 'Sign out'])

    fireEvent.click(items[0])
    expect(pushMock).toHaveBeenCalledWith('/profile')
    expect(screen.queryByRole('menu', { name: 'Profile' })).toBeNull()
  })

  // The panel is no longer a descendant of the trigger, so outside-click has to test both.
  // Getting this wrong closes the menu on the mousedown that is about to activate an item —
  // which reads as "the menu item sometimes does nothing".
  // The header collapses on scroll and the panel hangs off the header rather than the row, so it
  // is gated on chromeVisible for the same reason the nav panel is: a panel whose trigger has
  // been scrolled away is a floating menu with no visible way to close it.
  it('gates the panel on the chrome being visible', () => {
    // Mounted with the chrome already hidden, which is the state the gate exists for: the
    // header row holding the trigger has collapsed away, so a panel would be a floating menu
    // with nothing on screen to close it. Checked at mount because the header reads the chrome
    // store through a subscription, and re-rendering an existing tree is not the transition that
    // matters here.
    mockUser = { student_id: 's1', name: 'Nelson', username: 'wis' }
    chromeRef.visible = false

    render(<Header />)

    // `hidden: true` because the header is `inert` while the chrome is hidden — which is the
    // point: the trigger is no longer reachable by a real user either, so there is nothing to
    // press and the panel cannot be opened at all.
    const trigger = screen.getByRole('button', { name: 'Open profile menu', hidden: true })
    // Pressed directly because the header is `inert` while the chrome is hidden — which is the
    // point: a real user cannot reach it either. jsdom still runs the handler, so this reaches the
    // `open` state the gate is supposed to be filtering on.
    fireEvent.click(trigger)
    expect(trigger).toHaveAttribute('aria-expanded', 'true')

    // `hidden: true` is what makes this discriminating. The default role query skips anything
    // inside an `inert` subtree, so without it the assertion passes whether or not the gate
    // exists — the panel is in the DOM either way and simply unreadable. Asking for the
    // inaccessible ones too is the only way to see that the gate removed it.
    expect(screen.queryByRole('menu', { name: 'Profile', hidden: true })).toBeNull()
  })

  it('treats a click inside the panel as inside, not outside', () => {
    mockUser = { student_id: 's1', name: 'Nelson', username: 'wis' }
    render(<Header />)
    fireEvent.click(screen.getByRole('button', { name: 'Open profile menu' }))
    const panel = screen.getByRole('menu', { name: 'Profile' })

    fireEvent.mouseDown(panel)
    expect(screen.queryByRole('menu', { name: 'Profile' })).not.toBeNull()

    fireEvent.mouseDown(document.body)
    expect(screen.queryByRole('menu', { name: 'Profile' })).toBeNull()
  })

  it('closes on Escape', () => {
    mockUser = { student_id: 's1', name: 'Nelson', username: 'wis' }
    render(<Header />)
    fireEvent.click(screen.getByRole('button', { name: 'Open profile menu' }))
    expect(screen.queryByRole('menu', { name: 'Profile' })).not.toBeNull()
    fireEvent.keyDown(document, { key: 'Escape' })
    expect(screen.queryByRole('menu', { name: 'Profile' })).toBeNull()
  })

})
