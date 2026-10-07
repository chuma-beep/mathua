import { describe, it, expect, beforeAll, vi } from 'vitest'
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

vi.mock('../hooks/useAuthState', () => ({
  useAuthState: () => ({ loggedIn: false, user: null }),
}))

vi.mock('../lib/api', async importOriginal => ({
  ...(await importOriginal<typeof import('../lib/api')>()),
  getSettings: vi.fn(),
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
    // Domains is here rather than in the tab bar: the tabs stay at five, and overflowLinks()
    // derives the rest by subtracting TAB_HREFS from the desktop list.
    expect(labels).toEqual(['Domains', 'Leaderboard', 'Login'])
    const hrefs = items.map(i => (i as HTMLAnchorElement).getAttribute('href'))
    expect(hrefs).toEqual(['/domains', '/leaderboard', '/login'])
    expect(screen.getByRole('button', { name: 'Open navigation menu' }))
      .toHaveAttribute('aria-expanded', 'true')
  })

  it('renders a custom links list (docs section menu) when provided', () => {
    render(
      <Header
        links={[
          { label: 'Study', href: '/study' },
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
    expect(labels).toEqual(['Study', 'Docs', 'Architecture', 'Contributing', 'Note'])
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
    expect(screen.getAllByRole('menuitem')).toHaveLength(3)
    fireEvent.keyDown(document, { key: 'Escape' })
    expect(screen.queryAllByRole('menuitem')).toHaveLength(0)
  })

  it('spins forward on open and reverses on close without crashing', () => {
    render(<Header />)
    const button = screen.getByRole('button', { name: 'Open navigation menu' })
    fireEvent.click(button)
    expect(screen.getAllByRole('menuitem')).toHaveLength(3)
    fireEvent.click(button)
    expect(screen.queryAllByRole('menuitem')).toHaveLength(0)
    fireEvent.click(button)
    expect(screen.getAllByRole('menuitem')).toHaveLength(3)
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
