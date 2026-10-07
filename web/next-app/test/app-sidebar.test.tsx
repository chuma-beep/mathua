import { describe, it, expect } from 'vitest'
import { render, screen, within } from '@testing-library/react'
import { AppSidebar } from '../components/app-sidebar'
import { SidebarProvider } from '../components/ui/sidebar'

function renderSidebar(extra = {}) {
  return render(
    <SidebarProvider>
      <AppSidebar
        name="Ada"
        studentId="s1"
        username="ada"
        level="Scholar"
        streak={5}
        {...extra}
      />
    </SidebarProvider>
  )
}

// next/navigation stub for tests without a router
import { vi } from 'vitest'
vi.mock('next/navigation', () => ({
  usePathname: () => '/profile',
  useRouter: () => ({ push: () => {} }),
}))

// The collapsed 44px rail is only even because every icon in it occupies the same box.
//
// That was asserted as a `size-4` class on a lucide <svg>. The animated set renders
// `<div class="inline-flex" style="color:currentcolor"><svg width=16 height=16
// stroke="currentColor">`, so the class sizes nothing and the guarantee is now carried by explicit
// width/height on the svg. Asserting the class would have passed against an icon that no longer
// sizes to anything, which is the failure this replaces.
function expectFixedBox(link: Element) {
  const svg = link.querySelector('svg')
  expect(svg, 'nav link has no icon').not.toBeNull()
  expect(svg?.getAttribute('width'), 'icon width must be fixed, not a class').toBe('16')
  expect(svg?.getAttribute('height'), 'icon height must be fixed, not a class').toBe('16')
  // Colour has to come from the theme, not from the library's default.
  expect(svg?.getAttribute('stroke')).toBe('currentColor')
}

describe('AppSidebar', () => {
  it('renders Navigate group without a Diagnostic entry (no CTA duplication)', () => {
    renderSidebar()
    expect(screen.getByText('Navigate')).toBeInTheDocument()
    for (const label of ['Study', 'Learn', 'Graph', 'Leaderboard', 'History', 'Settings']) {
      expect(screen.getByText(label)).toBeInTheDocument()
    }
    expect(screen.queryByText('Diagnostic')).toBeNull()
    expect(screen.queryByText('Take Test')).toBeNull()
    expect(screen.queryByText('Review Now')).toBeNull()
  })

  it('orders Learn before Study, matching lib/nav.ts', () => {
    // The sidebar keeps its own list rather than importing nav.ts, because it
    // is deliberately a subset (no Review, no Login, plus Progress). That
    // makes the order a thing that can silently drift back to reference-first,
    // so it is asserted rather than assumed.
    renderSidebar()
    const sidebar = screen.getByTestId('profile-sidebar')
    const learn = within(sidebar).getByText('Learn', { exact: true })
    const study = within(sidebar).getByText('Study', { exact: true })
    // compareDocumentPosition: DOCUMENT_POSITION_FOLLOWING === 4
    expect(learn.compareDocumentPosition(study) & Node.DOCUMENT_POSITION_FOLLOWING).toBeTruthy()
  })

  it('renders Resources group with Docs, Contribute, and Creator note links', () => {
    renderSidebar()
    expect(screen.getByText('Resources')).toBeInTheDocument()
    for (const label of ['Docs', 'Contribute', "Creator's note"]) {
      expect(screen.getByText(label)).toBeInTheDocument()
    }
    const sidebar = screen.getByTestId('profile-sidebar')
    for (const href of ['/docs', '/docs/contributing', '/note']) {
      const link = sidebar.querySelector(`a[href="${href}"]`)!
      expect(link).not.toBeNull()
      expectFixedBox(link)
    }
  })

  it('renders slim sidebar without Next up / Goals / On-this-page duplication', () => {
    renderSidebar({ dueReviews: 2 })
    expect(screen.getByText('2')).toBeInTheDocument() // due-review badge
    // Removed groups never render
    expect(screen.queryByText('Goals')).toBeNull()
    expect(screen.queryByText('Daily goal')).toBeNull()
    expect(screen.queryByText('This week')).toBeNull()
    expect(screen.queryByText('On this page')).toBeNull()
    expect(screen.queryByText('Next up')).toBeNull()
    expect(screen.queryByText('Activity')).toBeNull()
  })

  it('names the transcript History, and nothing in the rail claims to be Progress', () => {
    // The entry was called Progress and resolved to a paginated list of past
    // questions, opening filtered to the ones you got wrong — the single most
    // misleading thing the navigation did. The transcript is now named for what
    // it is, and the report lives on the learner's own page rather than behind
    // a rail entry.
    renderSidebar()
    const sidebar = screen.getByTestId('profile-sidebar')
    const history = within(sidebar).getByText('History', { exact: true })
    expect(history.closest('a')?.getAttribute('href')).toBe('/history')
    expect(screen.queryByText('Progress', { exact: true })).toBeNull()
    expect(sidebar.querySelector('a[href="/progress-card"]')).toBeNull()
  })

  it('nav buttons render fixed-size Lucide icons (even collapsed rail)', () => {
    renderSidebar()
    const sidebar = screen.getByTestId('profile-sidebar')
    const hrefs = ['/study', '/learn', '/graph', '/leaderboard', '/history', '/settings']
    expect(hrefs).toHaveLength(6)
    for (const href of hrefs) {
      const link = sidebar.querySelector(`a[href="${href}"]`)!
      expect(link).not.toBeNull()
      expectFixedBox(link)
    }
  })

  it('menu labels hide in icon-collapse mode (icons only, no text peek)', () => {
    renderSidebar()
    const sidebar = screen.getByTestId('profile-sidebar')
    // Every text label next to an icon carries the collapse-hide class;
    // tooltips (not visible text) carry the label when collapsed.
    const labels = ['Mathua', 'Study', 'Learn', 'Graph', 'Leaderboard', 'History', 'Settings', 'Docs', 'Contribute', "Creator's note", 'Sign out']
    for (const label of labels) {
      const el = screen.getByText(label, { exact: true })
      expect(sidebar.contains(el)).toBe(true)
      expect(el.getAttribute('class') ?? '').toContain('group-data-[collapsible=icon]:hidden')
    }
  })

  it('footer shows identity and sign-out, identity text hides when collapsed', () => {
    renderSidebar()
    expect(screen.getByText('Ada')).toBeInTheDocument()
    expect(screen.getByText('Sign out')).toBeInTheDocument()
    // Collapsed rail shows avatar only — text wrapper carries the icon-mode hide class
    const ada = screen.getByText('Ada')
    const textWrapper = ada.closest('span[class*="flex-1"]')
    expect(textWrapper?.getAttribute('class') ?? '').toContain('group-data-[collapsible=icon]:hidden')
  })
})
