import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent, waitFor, within } from '@testing-library/react'
import DomainsPage from '../app/domains/page'
import DomainTable from '../components/DomainTable'
import { concepts } from '../lib/conceptData'
import { buildDomainRows, byNeedsAttention } from '../lib/domainRows'
import { DOMAIN_ORDER, domainLabel } from '../lib/graphDomains'
import { getProgress } from '../lib/api'
import { getGuestId, getUserInfo } from '../lib/auth'

vi.mock('../lib/api', async (importOriginal) => {
  const mod = importOriginal<typeof import('../lib/api')>()
  return { ...mod, getProgress: vi.fn(async () => ({})) }
})

vi.mock('../lib/auth', async (importOriginal) => {
  const mod = importOriginal<typeof import('../lib/auth')>()
  return { ...mod, getUserInfo: vi.fn(), getGuestId: vi.fn(() => 'guest-1') }
})

vi.mock('next/link', () => ({
  default: ({ children, href }: { children: React.ReactNode; href: string }) => (
    <a href={href}>{children}</a>
  ),
}))

vi.mock('../components/Header', () => ({ default: () => null }))
vi.mock('../components/Footer', () => ({ default: () => null }))
vi.mock('../components/BottomTabs', () => ({ default: () => null }))
vi.mock('../components/Loading', () => ({ default: () => <div>Loading</div> }))
vi.mock('../components/SectionHeader', () => ({
  default: ({ title }: { title: string }) => <h2>{title}</h2>,
}))

const MASTERED = { status: 'MASTERED', streak: 3 } as const

// A row's toggle button reads "N concepts" with no domain name, and the row wrapper carries no
// class on the last row, so neither `getByRole` nor `closest('div[class]')` can target one.
// Walk up from the domain label to the ancestor that actually holds the toggle.
function rowFor(label: string): HTMLElement {
  let node: HTMLElement | null = screen.getByRole('link', { name: label })
  while (node && !node.querySelector('button')) node = node.parentElement
  if (!node) throw new Error(`no row with a toggle found for ${label}`)
  return node
}

beforeEach(() => {
  vi.mocked(getUserInfo).mockReturnValue({ student_id: 's1' } as ReturnType<typeof getUserInfo>)
  vi.mocked(getGuestId).mockReturnValue('guest-1')
  vi.mocked(getProgress).mockResolvedValue({})
})

describe('/domains', () => {
  it('renders every domain in the canonical order', async () => {
    render(<DomainsPage />)
    await waitFor(() => expect(screen.getByRole('link', { name: 'Arithmetic' })).toBeTruthy())

    // Compared against DOMAIN_ORDER, not a count. The block this replaced hardcoded a
    // 15-entry list and omitted precalculus and machine_learning, so 50 concepts never
    // appeared anywhere — and a hardcoded "16 rows" assertion would have passed anyway.
    for (const domain of DOMAIN_ORDER) {
      const present = DOMAIN_ORDER.filter((d) =>
        buildDomainRows(concepts, {}).some((r) => r.domain === d),
      )
      expect(present, `${domain} has no concepts and so renders no row`).toContain(domain)
      expect(screen.getByRole('link', { name: domainLabel(domain) })).toBeTruthy()
    }
  })

  it('deep-links each domain into the graph filter', async () => {
    render(<DomainsPage />)
    await waitFor(() => expect(screen.getByRole('link', { name: 'Arithmetic' })).toBeTruthy())
    const link = screen.getByRole('link', { name: 'Arithmetic' }) as HTMLAnchorElement
    expect(link.getAttribute('href')).toBe('/graph?domain=arithmetic')
  })

  // Phase 3 removed the locked count from PositionBlock as the one number a learner can neither
  // act on nor change today. This table printed it on every row, which would have made it
  // appear 17 times instead of once.
  it('never prints a locked count', async () => {
    render(<DomainsPage />)
    await waitFor(() => expect(screen.getByRole('link', { name: 'Arithmetic' })).toBeTruthy())
    expect(document.body.textContent).not.toMatch(/locked/i)
  })

  it('expands a domain into its concepts, grouped by status', async () => {
    vi.mocked(getProgress).mockResolvedValue({
      [concepts.find((c) => c.domain === 'arithmetic')!.id]: MASTERED,
    })
    render(<DomainsPage />)
    await waitFor(() => expect(screen.getByRole('link', { name: 'Arithmetic' })).toBeTruthy())

    fireEvent.click(within(rowFor('Arithmetic')).getByRole('button', { name: /concepts$/ }))
    expect(screen.getByText(/^Mastered \(\d+\)$/)).toBeTruthy()
    expect(screen.getByText(/^Not started \(\d+\)$/)).toBeTruthy()
  })

  it('sorts weakest first on request', async () => {
    render(<DomainsPage />)
    await waitFor(() => expect(screen.getByRole('link', { name: 'Arithmetic' })).toBeTruthy())
    // Domain rows are the /graph?domain= links; "Open the graph" is a plain /graph link.
    const first = () =>
      screen
        .getAllByRole('link')
        .map((a) => a.getAttribute('href') ?? '')
        .find((h) => h.startsWith('/graph?domain='))

    const curriculumFirst = first()
    fireEvent.click(screen.getByRole('button', { name: 'Needs attention first' }))
    await waitFor(() => expect(first()).not.toBe(curriculumFirst))
    fireEvent.click(screen.getByRole('button', { name: 'Curriculum order' }))
    await waitFor(() => expect(first()).toBe(curriculumFirst))
  })

  it('still renders for a guest, who has progress before signing in', async () => {
    vi.mocked(getUserInfo).mockReturnValue(null)
    vi.mocked(getGuestId).mockReturnValue('guest-1')
    render(<DomainsPage />)
    await waitFor(() => expect(screen.getByRole('link', { name: 'Arithmetic' })).toBeTruthy())
    expect(getProgress).toHaveBeenCalledWith('guest-1')
  })
})

describe('DomainTable', () => {
  it('renders no links when links is off, for the read-only share report', () => {
    const rows = buildDomainRows(concepts, {})
    const { container } = render(<DomainTable rows={rows} links={false} />)
    expect(container.querySelectorAll('a').length).toBe(0)
    expect(container.textContent).toContain('Arithmetic')
  })

  it('omits the expand control when no toggle handler is given', () => {
    const rows = buildDomainRows(concepts, {})
    render(<DomainTable rows={rows} links={false} />)
    expect(screen.queryByRole('button')).toBeNull()
  })
})

describe('buildDomainRows', () => {
  it('counts a decayed concept as mastered and as review debt', () => {
    const id = concepts.find((c) => c.domain === 'arithmetic')!.id
    const [row] = buildDomainRows(concepts, { [id]: { status: 'DECAYING' } as never })
    expect(row.mastered).toBe(1)
    expect(row.dueForReview).toBe(1)
  })

  it('prefers a concept already in progress over one never started', () => {
    const members = concepts.filter((c) => c.domain === 'arithmetic')
    const started = members.find((c) => (c.prerequisites ?? []).length === 0)!
    const [row] = buildDomainRows(concepts, { [started.id]: { status: 'LEARNING' } as never })
    expect(row.nextId).toBe(started.id)
  })

  it('puts the weakest domain first when sorting by needs attention', () => {
    const rows = buildDomainRows(concepts, {})
    const weak = { ...rows[0], mastered: 0, label: 'Zzz' }
    const strong = { ...rows[1], mastered: rows[1].total, label: 'Aaa' }
    expect([weak, strong].sort(byNeedsAttention)[0].label).toBe('Zzz')
  })
})
