import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import OnboardPage from '../app/onboard/page'
import { concepts } from '../lib/conceptData'
import { scopeSize } from '../lib/diagnosticScope'
import { DOMAIN_ORDER, domainLabel } from '../lib/graphDomains'
import * as api from '../lib/api'

const push = vi.fn()
const replace = vi.fn()
vi.mock('next/navigation', () => ({
  useRouter: () => ({ push, replace }),
}))
vi.mock('next/link', () => ({
  default: ({ children, href }: { children: React.ReactNode; href: string }) => (
    <a href={href}>{children}</a>
  ),
}))
vi.mock('../lib/api', async (importOriginal) => {
  const mod = await importOriginal<typeof import('../lib/api')>()
  return { ...mod, startGoalDiagnostic: vi.fn(), getGoalPlan: vi.fn() }
})
vi.mock('../hooks/useTheme', () => ({ useTheme: () => ({ mounted: true }) }))
vi.mock('../components/Header', () => ({ default: () => null }))
vi.mock('../components/Footer', () => ({ default: () => null }))
vi.mock('../components/BottomTabs', () => ({ default: () => null }))
vi.mock('../components/Loading', () => ({ default: () => <div>Loading</div> }))

const SELECTION_KEY = 'mathua_onboard_domains'

function idsForDomain(domain: string): string[] {
  return concepts.filter(c => c.domain === domain).map(c => c.id)
}

/** The domain grid's start button, whatever copy it currently carries. */
function startButton(): HTMLButtonElement {
  return screen.getByRole('button', { name: /Start diagnostic test/ }) as HTMLButtonElement
}

function scopeShown(): string | null {
  const m = startButton().textContent?.match(/([\d,]+) concepts in scope/)
  return m ? m[1].replace(/,/g, '') : null
}

beforeEach(() => {
  push.mockReset()
  replace.mockReset()
  sessionStorage.clear()
  vi.mocked(api.startGoalDiagnostic).mockReset()
  vi.mocked(api.getGoalPlan).mockReset()
})

describe('/onboard domain selection', () => {
  it('renders every domain and starts with nothing selected', async () => {
    render(<OnboardPage />)
    await waitFor(() => expect(screen.getByText(domainLabel('arithmetic'))).toBeTruthy())
    for (const id of DOMAIN_ORDER) expect(screen.getByText(domainLabel(id))).toBeTruthy()
    expect(startButton().disabled).toBe(true)
  })

  it('refuses to start with an empty selection', async () => {
    render(<OnboardPage />)
    await waitFor(() => expect(startButton()).toBeTruthy())
    fireEvent.click(startButton())
    expect(api.startGoalDiagnostic).not.toHaveBeenCalled()
  })

  it('reports the real scope after selecting one domain', async () => {
    render(<OnboardPage />)
    await waitFor(() => expect(screen.getByText(domainLabel('arithmetic'))).toBeTruthy())
    fireEvent.click(screen.getByText(domainLabel('arithmetic')))
    expect(Number(scopeShown())).toBe(scopeSize(idsForDomain('arithmetic')))
  })

  it('widens the scope as more domains are selected', async () => {
    render(<OnboardPage />)
    await waitFor(() => expect(screen.getByText(domainLabel('arithmetic'))).toBeTruthy())
    fireEvent.click(screen.getByText(domainLabel('arithmetic')))
    const one = Number(scopeShown())
    fireEvent.click(screen.getByText(domainLabel('geometry')))
    expect(Number(scopeShown())).toBeGreaterThanOrEqual(one)
  })

  it('selects every domain, and the scope is the whole corpus', async () => {
    render(<OnboardPage />)
    await waitFor(() => expect(screen.getByText('Select everything')).toBeTruthy())
    fireEvent.click(screen.getByText('Select everything'))
    expect(Number(scopeShown())).toBe(concepts.length)
    expect(startButton().disabled).toBe(false)
  })

  it('deselects one domain back out after selecting everything', async () => {
    render(<OnboardPage />)
    await waitFor(() => expect(screen.getByText('Select everything')).toBeTruthy())
    fireEvent.click(screen.getByText('Select everything'))
    const all = Number(scopeShown())

    fireEvent.click(screen.getByText(domainLabel('arithmetic')))
    const after = Number(scopeShown())
    expect(after).toBeLessThan(all)

    // Toggling it back restores the full scope rather than losing another domain.
    fireEvent.click(screen.getByText(domainLabel('arithmetic')))
    expect(Number(scopeShown())).toBe(all)
  })

  it('sends the selected domains up as concept ids and nothing else', async () => {
    vi.mocked(api.startGoalDiagnostic).mockResolvedValue({
      session_id: 's1', concept_id: 'c1', concept_name: 'C', question: '1 + 1 = ?',
      progress: { answered: 0, estimated_total: 25, min_total: 25, max_total: 45, cover_done: 0, cover_size: 5, done: false },
    } as never)
    render(<OnboardPage />)
    await waitFor(() => expect(screen.getByText(domainLabel('arithmetic'))).toBeTruthy())
    fireEvent.click(screen.getByText(domainLabel('arithmetic')))

    fireEvent.click(startButton())
    await screen.findByText('Before you begin')
    fireEvent.click(screen.getByRole('button', { name: /Begin diagnostic/ }))

    await waitFor(() => expect(api.startGoalDiagnostic).toHaveBeenCalledTimes(1))
    const sent = vi.mocked(api.startGoalDiagnostic).mock.calls[0][0]
    expect(sent).toHaveLength(scopeSize(idsForDomain('arithmetic')))
    expect(sent.sort()).toEqual([...idsForDomain('arithmetic')].sort())
  })

  it('keeps a selection across a reload and drops names that no longer exist', async () => {
    sessionStorage.setItem(SELECTION_KEY, JSON.stringify(['arithmetic', 'a_domain_that_was_removed']))
    render(<OnboardPage />)
    await waitFor(() => expect(screen.getByText(domainLabel('arithmetic'))).toBeTruthy())
    expect(Number(scopeShown())).toBe(scopeSize(idsForDomain('arithmetic')))
    // The unknown name is dropped rather than resurrecting an empty domain.
    expect(scopeShown()).not.toBeNull()
  })

  it('offers a way into the topics without taking the test', async () => {
    render(<OnboardPage />)
    await waitFor(() => expect(screen.getByText(/Browse topics/)).toBeTruthy())
    const link = screen.getByText(/Browse topics/).closest('a')
    expect(link?.getAttribute('href')).toBe('/domains')
  })

  it('names domains for the learner, not by their corpus id', async () => {
    // Regression: the grid indexed the label function instead of calling it, so
    // every card fell back to the raw id and the picker read "arithmetic" and
    // "complex_numbers". /goals called it correctly, which is how the two
    // screens drifted apart in the first place.
    render(<OnboardPage />)
    await waitFor(() => expect(screen.getByText(domainLabel('arithmetic'))).toBeTruthy())
    for (const id of DOMAIN_ORDER) {
      expect(screen.getByText(domainLabel(id))).toBeTruthy()
      expect(screen.queryByText(id)).toBeNull()
    }
  })
})