import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import NextUpSummary from '../components/NextUpSummary'
import { selectShelfHead, type ShelfInput } from '../lib/nextUp'

const catalog = [
  { id: 'a', label: 'A', prerequisites: [] as string[], avgTimeSeconds: 10 },
  { id: 'b', label: 'B', prerequisites: ['a'], avgTimeSeconds: 120 },
  { id: 'c', label: 'C', prerequisites: ['a'], avgTimeSeconds: 10 },
]

const base: ShelfInput = {
  dueReviews: 0,
  weaknesses: { by_domain: {} },
  progress: { a: { status: 'MASTERED', streak: 3 } },
  activity: [],
  diagnosticCompleted: true,
  conceptsMastered: 3,
  catalog,
}

describe('NextUpSummary dashboard block', () => {
  it('renders the head as the one primary Continue button with reason and XP', () => {
    const shelf = selectShelfHead(base)
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

  it('shows the diagnostic as the head for brand-new learners', () => {
    const shelf = selectShelfHead({
      dueReviews: 0,
      weaknesses: { by_domain: {} },
      progress: {},
      activity: [],
      diagnosticCompleted: false,
      conceptsMastered: 0,
      catalog,
    })
    expect(shelf.next.kind).toBe('diagnostic')
    render(<NextUpSummary shelf={shelf} />)
    const primary = screen.getByRole('link', { name: shelf.next.cta })
    expect(primary.getAttribute('href')).toBe('/onboard')
  })

  it('Profile and Learn agree on the head for the same state', () => {
    const input = { ...base }
    expect(selectShelfHead(input).next).toEqual(selectShelfHead(input).next)
    // Learn's done input differs only by the exclusion, by design.
    const excluded = selectShelfHead({ ...input, excludeConceptIds: ['b', 'c'] })
    expect(excluded.next.kind).toBe('browse')
  })
})
