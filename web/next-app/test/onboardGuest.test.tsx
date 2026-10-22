import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import OnboardPage from '../app/onboard/page'
import { ResultsStep } from '../app/onboard/steps'
import { concepts } from '../lib/conceptData'
import { scopeSize } from '../lib/diagnosticScope'
import { domainLabel } from '../lib/graphDomains'
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
  return {
    ...mod,
    startGoalDiagnostic: vi.fn(),
    startGoalDiagnosticName: vi.fn(),
    getGoalPlan: vi.fn(),
  }
})
vi.mock('../hooks/useTheme', () => ({ useTheme: () => ({ mounted: true }) }))
vi.mock('../components/Header', () => ({ default: () => null }))
vi.mock('../components/Footer', () => ({ default: () => null }))
vi.mock('../components/BottomTabs', () => ({ default: () => null }))
vi.mock('../components/Loading', () => ({ default: () => <div>Loading</div> }))
vi.mock('../components/DiagnosticResults', () => ({ default: () => <div>plan body</div> }))

function idsForDomain(domain: string): string[] {
  return concepts.filter(c => c.domain === domain).map(c => c.id)
}

/** A start response shaped like the server's, enough to reach the test step. */
function startResponse() {
  return {
    session_id: 's1',
    concept_id: 'c1',
    concept_name: 'C',
    question: '1 + 1 = ?',
    progress: { answered: 0, estimated_total: 25, min_total: 25, max_total: 45, probe_budget: 45, done: false },
  } as never
}

async function reachBegin() {
  render(<OnboardPage />)
  await waitFor(() => expect(screen.getByText(domainLabel('arithmetic'))).toBeTruthy())
  fireEvent.click(screen.getByText(domainLabel('arithmetic')))
  fireEvent.click(screen.getByRole('button', { name: /Start diagnostic test/ }))
  await screen.findByText('Before you begin')
}

beforeEach(() => {
  push.mockReset()
  replace.mockReset()
  sessionStorage.clear()
  vi.mocked(api.startGoalDiagnostic).mockReset()
  vi.mocked(api.startGoalDiagnosticName).mockReset()
  vi.mocked(api.getGoalPlan).mockReset()
})

// The diagnostic used to be signed-in only, and /onboard sent an
// unauthenticated start anyway: the server answered 401 and the page showed
// "Something went wrong" for a screen that was behaving as designed. The 401 is
// now a fork in the road, not an error.
describe('a guest starting the diagnostic', () => {
  it('asks for a name instead of showing a generic failure', async () => {
    vi.mocked(api.startGoalDiagnostic).mockRejectedValue(
      Object.assign(new Error('not authenticated'), { status: 401 }),
    )
    await reachBegin()
    fireEvent.click(screen.getByRole('button', { name: /Begin diagnostic/ }))

    await screen.findByText('What should we call you?')
    expect(screen.queryByText(/Something went wrong/)).toBeNull()
    // The copy must not imply the run is saved; that is the whole point of a demo.
    expect(screen.getByText(/not attached to an account/)).toBeTruthy()
  })

  it('runs the demo under the name it was given, trimmed', async () => {
    vi.mocked(api.startGoalDiagnostic).mockRejectedValue(
      Object.assign(new Error('not authenticated'), { status: 401 }),
    )
    vi.mocked(api.startGoalDiagnosticName).mockResolvedValue(startResponse())
    await reachBegin()
    fireEvent.click(screen.getByRole('button', { name: /Begin diagnostic/ }))
    await screen.findByText('What should we call you?')

    const run = screen.getByRole('button', { name: /Run the demo/ }) as HTMLButtonElement
    expect(run.disabled).toBe(true)

    fireEvent.change(screen.getByLabelText('Your name'), { target: { value: '  Ada  ' } })
    expect(run.disabled).toBe(false)
    fireEvent.click(run)

    await waitFor(() => expect(api.startGoalDiagnosticName).toHaveBeenCalledTimes(1))
    const [name, ids] = vi.mocked(api.startGoalDiagnosticName).mock.calls[0]
    expect(name).toBe('Ada')
    expect(ids.sort()).toEqual(idsForDomain('arithmetic').sort())
    expect(ids).toHaveLength(scopeSize(idsForDomain('arithmetic')))
  })

  it('keeps a real 500 as a plain failure, not an invitation to demo', async () => {
    vi.mocked(api.startGoalDiagnostic).mockRejectedValue(
      Object.assign(new Error('boom'), { status: 500 }),
    )
    await reachBegin()
    fireEvent.click(screen.getByRole('button', { name: /Begin diagnostic/ }))

    await waitFor(() => expect(screen.queryByText('What should we call you?')).toBeNull())
    expect(screen.queryByRole('button', { name: /Run the demo/ })).toBeNull()
  })
})

describe('a demo result', () => {
  it('says the plan will not survive the session, and offers a real account', () => {
    render(<ResultsStep plan={{} as never} onStartPractice={() => {}} demo />)
    expect(screen.getByText(/This was a demo run/)).toBeTruthy()
    const link = screen.getByText(/Sign in to keep a plan/).closest('a')
    expect(link?.getAttribute('href')).toBe('/login')
  })

  it('says nothing about demos on a signed-in result', () => {
    render(<ResultsStep plan={{} as never} onStartPractice={() => {}} />)
    expect(screen.queryByText(/This was a demo run/)).toBeNull()
  })
})
