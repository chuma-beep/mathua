import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import NextUpSummary from '../components/NextUpSummary'
import type { Shelf } from '../lib/nextUp'

type ShelfShape = Parameters<typeof NextUpSummary>[0]['shelf']

const shelfOf = (next: Partial<Shelf['next']>, alternatives: Shelf['alternatives'] = []): ShelfShape =>
  ({
    next: {
      kind: 'new', badge: 'New', title: 'B', detail: 'Ready to learn', href: '/learn?concept=b',
      cta: 'Start →', xp: 2, ...next,
    },
    alternatives,
  }) as ShelfShape

describe('NextUpSummary dashboard block', () => {
  it('renders the head as the one primary button with reason and XP', () => {
    const shelf = shelfOf({ badge: 'New', title: 'B', cta: 'Start →', xp: 2 })
    render(<NextUpSummary shelf={shelf} />)
    expect(screen.getByText(`Continue: ${shelf.next.title}`)).toBeTruthy()
    expect(screen.getAllByText(shelf.next.badge).length).toBeGreaterThan(0)
    const primary = screen.getByRole('link', { name: shelf.next.cta })
    expect(primary.getAttribute('href')).toBe(shelf.next.href)
    if (shelf.next.xp > 0) {
      expect(screen.getByText(new RegExp(`\\+${shelf.next.xp} XP`))).toBeTruthy()
    }
  })

  it('renders alternatives as quiet rows excluding the head href, concept, and dupes', () => {
    const shelf = {
      next: { kind: 'new', badge: 'New', title: 'B', detail: 'd', href: '/learn?concept=b', cta: 'Continue →', xp: 2 },
      alternatives: [
        { kind: 'new', badge: 'New', title: 'B again', detail: 'd', href: '/learn?concept=b&seed=9', cta: 'Continue →', xp: 2 },
        { kind: 'new', badge: 'New', title: 'B dup', detail: 'd', href: '/learn?concept=b', cta: 'Continue →', xp: 2 },
        { kind: 'weakness', badge: 'Recommended', title: 'C', detail: 'd', href: '/learn?concept=c', cta: 'Continue →', xp: 1 },
        { kind: 'weakness', badge: 'Recommended', title: 'C dup', detail: 'd', href: '/learn?concept=c', cta: 'Continue →', xp: 1 },
        { kind: 'review', badge: 'Due now', title: 'R', detail: 'd', href: '/review', cta: 'Review →', xp: 1 },
      ],
    } as unknown as Parameters<typeof NextUpSummary>[0]['shelf']
    render(<NextUpSummary shelf={shelf} />)
    const links = screen.getAllByRole('link') as HTMLAnchorElement[]
    const hrefs = links.map((l) => l.getAttribute('href'))
    // Primary + C + review only: same-concept seed variant and dupes dropped.
    expect(hrefs).toEqual(['/learn?concept=b', '/learn?concept=c', '/review'])
  })

  it('caps the queue at four rows', () => {
    const shelf = {
      next: { kind: 'browse', badge: 'Study', title: 'Browse', detail: 'd', href: '/study', cta: 'Browse →', xp: 0 },
      alternatives: Array.from({ length: 6 }, (_, i) => ({
        kind: 'new', badge: 'New', title: `N${i}`, detail: 'd',
        href: `/learn?concept=n${i}`, cta: 'Continue →', xp: 1,
      })),
    } as unknown as Parameters<typeof NextUpSummary>[0]['shelf']
    render(<NextUpSummary shelf={shelf} />)
    expect(screen.getAllByRole('link')).toHaveLength(5) // primary + 4
  })

  it('renders a diagnostic head the engine chose, without deciding anything itself', () => {
    // Whether a learner is offered the diagnostic is decided in internal/scheduler and covered
    // by TestRecommend_NewLearnerIsOfferedTheDiagnosticButStillHasAgency. What this asserts is
    // the other half: the surface renders the engine's choice rather than substituting its own
    // judgement for it.
    const shelf = shelfOf({
      kind: 'diagnostic', badge: 'Recommended', title: 'Find where to start',
      detail: 'A short adaptive test finds your level', href: '/onboard', cta: 'Start test →', xp: 0,
    })
    render(<NextUpSummary shelf={shelf} />)
    const primary = screen.getByRole('link', { name: shelf.next.cta })
    expect(primary.getAttribute('href')).toBe('/onboard')
  })

  it('renders a mastery-check task as a task, with its href intact', () => {
    const shelf = shelfOf({
      kind: 'diagnostic', badge: 'Mastery check', title: 'Mastery check',
      detail: 'You have earned 50 XP since your last one', href: '/goals?quiz=1',
      cta: 'Take the check →', xp: 0,
    })
    render(<NextUpSummary shelf={shelf} />)
    expect(screen.getByRole('link', { name: shelf.next.cta }).getAttribute('href')).toBe('/goals?quiz=1')
  })
})

// The shelf arrives over an endpoint now, so it is null for the first paint. This component
// destructured `shelf.next` unguarded and took /profile down into its 500 boundary — React
// error #310, which on a minified build says nothing about where. jsdom never saw it because
// the mock resolved before the assertions ran; `next build` + Playwright did.
describe('NextUpSummary before the recommendation arrives', () => {
  it('renders nothing rather than throwing', () => {
    const spy = vi.spyOn(console, 'error').mockImplementation(() => {})
    try {
      const { container } = render(<NextUpSummary shelf={null} />)
      expect(container.querySelector('[aria-label="Next up"]')).toBeNull()
      expect(spy).not.toHaveBeenCalled()
    } finally {
      spy.mockRestore()
    }
  })
})
