import { fireEvent, render, screen, waitFor } from '@testing-library/react'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import LearnPage from '../app/learn/page'
import type { Shelf } from '../lib/nextUp'

const { fetchShelf } = vi.hoisted(() => ({ fetchShelf: vi.fn() }))

vi.mock('../lib/recommendations', () => ({ fetchShelf }))
vi.mock('next/navigation', () => ({
  useSearchParams: () => new URLSearchParams(),
  usePathname: () => '/learn',
  useRouter: () => ({ push: vi.fn() }),
}))
vi.mock('next/link', () => ({
  default: ({ children, href }: { children: React.ReactNode; href: string }) => (
    <a href={href}>{children}</a>
  ),
}))
vi.mock('../components/Header', () => ({ default: () => null }))
vi.mock('../components/Footer', () => ({ default: () => null }))
vi.mock('../components/BottomTabs', () => ({ default: () => null }))
vi.mock('../components/LearnStepper', () => ({ default: () => null }))

const shelf: Shelf = {
  next: {
    kind: 'new',
    badge: 'New',
    title: 'Beta',
    detail: 'Ready to learn',
    href: '/learn?concept=b',
    cta: 'Start →',
    xp: 2,
  },
  alternatives: [],
}

beforeEach(() => {
  fetchShelf.mockReset()
})

// The queue is the one thing this page exists to state. When the engine cannot
// be reached it used to say "All learned", and when the request failed outright
// the catch left the shelf null, which rendered as a spinner that never stopped.
// Both are false statements about a learner's work, and both hid the fact that
// nothing had been read at all.
describe('the learn entry point', () => {
  it('says the queue could not be read, not that there is nothing to do', async () => {
    fetchShelf.mockRejectedValue(new Error('offline'))
    render(<LearnPage />)

    expect(await screen.findByText('Could not load what to study next.')).toBeInTheDocument()
    expect(screen.queryByText(/All learned/i)).toBeNull()
    expect(screen.queryByText('Loading…')).toBeNull()
    expect(screen.getByRole('link', { name: 'Browse domains' })).toHaveAttribute('href', '/domains')
  })

  it('retries the same request and shows the queue when it answers', async () => {
    fetchShelf.mockRejectedValueOnce(new Error('offline')).mockResolvedValue(shelf)
    render(<LearnPage />)

    fireEvent.click(await screen.findByRole('button', { name: 'Try again' }))
    await waitFor(() => expect(fetchShelf).toHaveBeenCalledTimes(2))
    expect(await screen.findByText('Beta')).toBeInTheDocument()
    expect(screen.queryByText('Could not load what to study next.')).toBeNull()
  })
})
