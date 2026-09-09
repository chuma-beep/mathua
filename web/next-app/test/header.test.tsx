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

  it('opens the default logged-out links with correct hrefs on click', () => {
    render(<Header />)
    fireEvent.click(screen.getByRole('button', { name: 'Open navigation menu' }))
    const items = screen.getAllByRole('menuitem')
    const labels = items.map(i => i.textContent)
    expect(labels).toEqual(['Study', 'Leaderboard', 'Graph', 'Login'])
    const hrefs = items.map(i => (i as HTMLAnchorElement).getAttribute('href'))
    expect(hrefs).toEqual(['/study', '/leaderboard', '/graph', '/login'])
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
    expect(screen.getAllByRole('menuitem')).toHaveLength(4)
    fireEvent.keyDown(document, { key: 'Escape' })
    expect(screen.queryAllByRole('menuitem')).toHaveLength(0)
  })

  it('spins forward on open and reverses on close without crashing', () => {
    render(<Header />)
    const button = screen.getByRole('button', { name: 'Open navigation menu' })
    fireEvent.click(button)
    expect(screen.getAllByRole('menuitem')).toHaveLength(4)
    fireEvent.click(button)
    expect(screen.queryAllByRole('menuitem')).toHaveLength(0)
    fireEvent.click(button)
    expect(screen.getAllByRole('menuitem')).toHaveLength(4)
  })

  it('closes when a link is clicked', () => {
    render(<Header />)
    fireEvent.click(screen.getByRole('button', { name: 'Open navigation menu' }))
    fireEvent.click(screen.getByRole('menuitem', { name: 'Study' }))
    expect(screen.queryAllByRole('menuitem')).toHaveLength(0)
  })

  it('theme toggle plays the sun-moon animation and toggles the theme', () => {
    toggleThemeMock.mockClear()
    render(<Header />)
    fireEvent.click(screen.getByRole('button', { name: 'Toggle theme' }))
    expect(toggleThemeMock).toHaveBeenCalledTimes(1)
  })
})
