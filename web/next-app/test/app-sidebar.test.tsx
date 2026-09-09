import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
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

describe('AppSidebar', () => {
  it('renders Navigate group without a Diagnostic entry (no CTA duplication)', () => {
    renderSidebar()
    expect(screen.getByText('Navigate')).toBeInTheDocument()
    for (const label of ['Study', 'Start', 'Graph', 'Leaderboard', 'Settings']) {
      expect(screen.getByText(label)).toBeInTheDocument()
    }
    expect(screen.queryByText('Diagnostic')).toBeNull()
    expect(screen.queryByText('Take Test')).toBeNull()
    expect(screen.queryByText('Review Now')).toBeNull()
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
      const svg = link.querySelector('svg')
      expect(svg).not.toBeNull()
      expect(svg?.getAttribute('class') ?? '').toMatch(/size-4/)
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

  it('nav buttons render fixed-size Lucide icons (even collapsed rail)', () => {
    renderSidebar()
    const sidebar = screen.getByTestId('profile-sidebar')
    const hrefs = ['/study', '/session', '/graph', '/leaderboard', '/settings']
    expect(hrefs).toHaveLength(5)
    for (const href of hrefs) {
      const link = sidebar.querySelector(`a[href="${href}"]`)!
      expect(link).not.toBeNull()
      const svg = link.querySelector('svg')
      expect(svg).not.toBeNull()
      expect(svg?.getAttribute('class') ?? '').toMatch(/size-4/)
    }
  })

  it('menu labels hide in icon-collapse mode (icons only, no text peek)', () => {
    renderSidebar()
    const sidebar = screen.getByTestId('profile-sidebar')
    // Every text label next to an icon carries the collapse-hide class;
    // tooltips (not visible text) carry the label when collapsed.
    const labels = ['Mathua', 'Study', 'Start', 'Graph', 'Leaderboard', 'Settings', 'Docs', 'Contribute', "Creator's note", 'Sign out']
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
