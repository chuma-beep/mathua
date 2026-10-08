import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import DomainsPage from '../app/domains/page'
import DomainTable from '../components/DomainTable'
import { concepts } from '../lib/conceptData'
import { buildDomainRows, byNeedsAttention } from '../lib/domainRows'
import { DOMAIN_ORDER, domainLabel } from '../lib/graphDomains'
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
  return {
    ...mod,
    getCurriculumDomains: vi.fn(),
    getCurriculumDomain: vi.fn(),
    getLessonEligibility: vi.fn(),
  }
})
vi.mock('../components/Header', () => ({ default: () => null }))
vi.mock('../components/Footer', () => ({ default: () => null }))
vi.mock('../components/BottomTabs', () => ({ default: () => null }))
vi.mock('../components/Loading', () => ({ default: ({ label }: { label?: string }) => <div>{label ?? 'Loading'}</div> }))
vi.mock('../components/SectionHeader', () => ({
  default: ({ title }: { title: string }) => <h2>{title}</h2>,
}))

const domains = DOMAIN_ORDER.map((id) => ({
  id,
  conceptCount: 10,
  masteredCount: 0,
  unlockedCount: 2,
  inProgressCount: 1,
}))

beforeEach(() => {
  push.mockReset()
  vi.mocked(api.getCurriculumDomains).mockReset().mockResolvedValue(domains)
  vi.mocked(api.getCurriculumDomain).mockReset()
  vi.mocked(api.getLessonEligibility).mockReset()
})

describe('/domains curriculum picker', () => {
  it('renders every domain in the canonical order with availability', async () => {
    render(<DomainsPage />)
    await waitFor(() => expect(screen.getByText(domainLabel('arithmetic'))).toBeTruthy())
    for (const domain of DOMAIN_ORDER) {
      expect(screen.getByText(domainLabel(domain))).toBeTruthy()
    }
    expect(screen.getAllByText(/\/ \d+ concepts available/).length).toBeGreaterThan(0)
  })

  it('opens a domain and renders its topics with server states', async () => {
    vi.mocked(api.getCurriculumDomain).mockResolvedValue({
      domain: { id: 'arithmetic' },
      topics: [
        { conceptId: 'arith.add', title: 'Single-digit addition', state: 'unlocked', depth: 0 },
        { conceptId: 'arith.sub', title: 'Subtraction', state: 'locked', depth: 1 },
      ],
    })
    render(<DomainsPage />)
    await waitFor(() => expect(screen.getByText(domainLabel('arithmetic'))).toBeTruthy())
    fireEvent.click(screen.getByText(domainLabel('arithmetic')))
    await screen.findByText('Single-digit addition')
    expect(screen.getByText('Available')).toBeTruthy()
    expect(screen.getByText('Locked')).toBeTruthy()
  })

  it('checks eligibility before entering /learn', async () => {
    vi.mocked(api.getCurriculumDomain).mockResolvedValue({
      domain: { id: 'arithmetic' },
      topics: [{ conceptId: 'arith.add', title: 'Single-digit addition', state: 'unlocked', depth: 0 }],
    })
    vi.mocked(api.getLessonEligibility).mockResolvedValue({
      concept_id: 'arith.add', ready: true, weak: [], missing: [], eligible: true,
    })
    render(<DomainsPage />)
    await waitFor(() => expect(screen.getByText(domainLabel('arithmetic'))).toBeTruthy())
    fireEvent.click(screen.getByText(domainLabel('arithmetic')))
    await screen.findByRole('button', { name: 'Learn' })
    fireEvent.click(screen.getByRole('button', { name: 'Learn' }))
    await waitFor(() => expect(push).toHaveBeenCalledWith('/learn?concept=arith.add'))
  })

  it('explains a locked topic instead of offering a start', async () => {
    vi.mocked(api.getCurriculumDomain).mockResolvedValue({
      domain: { id: 'arithmetic' },
      topics: [{ conceptId: 'arith.sub', title: 'Subtraction', state: 'locked', depth: 1 }],
    })
    vi.mocked(api.getLessonEligibility).mockResolvedValue({
      concept_id: 'arith.sub', ready: false, weak: [], missing: [], eligible: false,
      prerequisites: [{ conceptId: 'arith.add', title: 'Single-digit addition', state: 'unlocked', met: false }],
    })
    render(<DomainsPage />)
    await waitFor(() => expect(screen.getByText(domainLabel('arithmetic'))).toBeTruthy())
    fireEvent.click(screen.getByText(domainLabel('arithmetic')))
    fireEvent.click(await screen.findByRole('button', { name: /View prerequisite path/ }))
    await screen.findByRole('note', { name: /prerequisite/i })
    expect(screen.getByRole('link', { name: /Learn prerequisite/ }).getAttribute('href'))
      .toBe('/learn?concept=arith.add')
    expect(push).not.toHaveBeenCalled()
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
  it('puts the weakest domain first when sorting by needs attention', () => {
    const rows = buildDomainRows(concepts, {})
    const weak = { ...rows[0], mastered: 0, label: 'Zzz' }
    const strong = { ...rows[1], mastered: rows[1].total, label: 'Aaa' }
    expect([weak, strong].sort(byNeedsAttention)[0].label).toBe('Zzz')
  })
})
