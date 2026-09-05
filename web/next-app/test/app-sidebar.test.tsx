import { describe, it, expect } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import { AppSidebar } from '../components/app-sidebar'
import { SidebarProvider } from '../components/ui/sidebar'

const TOC = [
  { id: 'next', label: 'Next up' },
  { id: 'activity', label: 'Activity' },
]

function renderSidebar(extra = {}) {
  return render(
    <SidebarProvider>
      <AppSidebar
        name="Ada"
        studentId="s1"
        username="ada"
        level="Scholar"
        streak={5}
        toc={TOC}
        activeId="activity"
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

  it('renders Next up / Goals slots and TOC with active state', () => {
    renderSidebar({ nextSlot: <div>next-up</div>, goalsSlot: <div>goals</div>, dueReviews: 2 })
    expect(screen.getByText('next-up')).toBeInTheDocument()
    expect(screen.getByText('goals')).toBeInTheDocument()
    expect(screen.getByText('2')).toBeInTheDocument() // due-review badge
    const tocButtons = screen.getAllByRole('button')
    const activity = tocButtons.find((b) => b.textContent?.includes('Activity'))
    expect(activity).toHaveAttribute('data-active', 'true')
  })

  it('TOC buttons scroll to sections via onNavigate', () => {
    const onNavigate = vi.fn()
    const section = document.createElement('section')
    section.id = 'next'
    section.scrollIntoView = vi.fn()
    document.body.appendChild(section)
    renderSidebar({ onNavigate })
    const tocButtons = screen.getAllByRole('button')
    const next = tocButtons.find((b) => b.textContent?.includes('Next up'))
    fireEvent.click(next!)
    expect(onNavigate).toHaveBeenCalledWith('next')
    expect(section.scrollIntoView).toHaveBeenCalled()
    document.body.removeChild(section)
  })

  it('footer shows identity and sign-out', () => {
    renderSidebar()
    expect(screen.getByText('Ada')).toBeInTheDocument()
    expect(screen.getByText('Sign out')).toBeInTheDocument()
  })
})
