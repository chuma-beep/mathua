import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import RecommendedTaskList from '../components/RecommendedTaskList'
import type { Shelf } from '../lib/nextUp'
import * as api from '../lib/api'

const push = vi.fn()
vi.mock('next/navigation', () => ({
  useRouter: () => ({ push }),
}))
vi.mock('next/link', () => ({
  default: ({ children, href }: { children: React.ReactNode; href: string }) => (
    <a href={href}>{children}</a>
  ),
}))
vi.mock('../lib/api', async (importOriginal) => {
  const mod = await importOriginal<typeof import('../lib/api')>()
  return { ...mod, getLessonEligibility: vi.fn() }
})

const item = (over: Partial<Shelf['next']>): Shelf['next'] => ({
  kind: 'new',
  badge: 'New',
  title: 'Single-digit addition',
  detail: 'Ready to learn',
  href: '/learn?concept=arith.add',
  cta: 'Start →',
  xp: 1,
  ...over,
})

function shelf(): Shelf {
  return {
    next: item({}),
    alternatives: [
      item({ href: '/learn?concept=neg.numline', title: 'Negative numbers', badge: 'New' }),
      item({ href: '/review', title: 'Review check', kind: 'review', badge: 'Review', cta: 'Review →' }),
    ],
  }
}

beforeEach(() => {
  push.mockReset()
  vi.mocked(api.getLessonEligibility).mockReset()
})

describe('RecommendedTaskList', () => {
  it('shows a loading state without fake concepts', () => {
    render(<RecommendedTaskList shelf={null} loading />)
    expect(screen.getByRole('status').textContent).toMatch(/Loading recommended tasks/)
    expect(screen.queryByText('Single-digit addition')).toBeNull()
  })

  it('renders one actionable card per task, deduped by destination', () => {
    render(<RecommendedTaskList shelf={shelf()} />)
    const region = screen.getByRole('region', { name: 'Your learning' })
    // Duplicate href collapses to one card; buttons are the actionable controls.
    expect(region.querySelectorAll('button')).toHaveLength(3)
    expect(screen.getByText('Single-digit addition')).toBeTruthy()
    expect(screen.getByText('Negative numbers')).toBeTruthy()
  })

  it('offers an empty state when there are no tasks', () => {
    render(<RecommendedTaskList shelf={null} />)
    expect(screen.getByText(/No learning tasks are available/)).toBeTruthy()
    expect(screen.getByRole('link', { name: 'Choose what to learn' })).toBeTruthy()
  })

  it('asks the server for eligibility and enters /learn when eligible', async () => {
    vi.mocked(api.getLessonEligibility).mockResolvedValue({
      concept_id: 'arith.add', ready: true, weak: [], missing: [], eligible: true,
    })
    render(<RecommendedTaskList shelf={shelf()} />)
    fireEvent.click(screen.getByText('Single-digit addition').closest('button')!)
    await waitFor(() => expect(api.getLessonEligibility).toHaveBeenCalledWith('arith.add'))
    await waitFor(() => expect(push).toHaveBeenCalledWith('/learn?concept=arith.add'))
  })

  it('explains a locked topic instead of starting it', async () => {
    vi.mocked(api.getLessonEligibility).mockResolvedValue({
      concept_id: 'arith.add', ready: false, weak: [], missing: [],
      eligible: false,
      reason: 'missing_prerequisites',
      prerequisites: [
        { conceptId: 'arith.count', title: 'Counting', state: 'unlocked', met: false },
      ],
    })
    render(<RecommendedTaskList shelf={shelf()} />)
    fireEvent.click(screen.getByText('Single-digit addition').closest('button')!)
    await screen.findByRole('note', { name: /prerequisite/i })
    expect(screen.getByText('Counting')).toBeTruthy()
    // The prerequisite routes into the learning system, never into reference.
    expect(screen.getByRole('link', { name: /Learn prerequisite/ }).getAttribute('href'))
      .toBe('/learn?concept=arith.count')
    expect(push).not.toHaveBeenCalled()
  })

  it('follows a non-concept task destination as-is', () => {
    render(<RecommendedTaskList shelf={shelf()} />)
    fireEvent.click(screen.getByText('Review check').closest('button')!)
    expect(push).toHaveBeenCalledWith('/review')
    expect(api.getLessonEligibility).not.toHaveBeenCalled()
  })
})
