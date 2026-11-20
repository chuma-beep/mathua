import { describe, it, expect, vi, beforeEach, afterAll } from 'vitest'

// The default 5000 ms is too tight for this file, and it fails as a timeout rather
// than as a wrong assertion. Each case mounts the welcome step, the review step, the
// quiz host and the answer editor, which individually take 421-799 ms; under a full
// parallel run the jsdom worker also waits behind other environments. Measured by
// rerunning the suite with --testTimeout=3000, which failed exactly these cases.
// `learnStepperDone.test.tsx` already grants comparable time to comparably heavy
// cases. This is a budget for the environment, not a retry: a genuine hang still fails.
vi.setConfig({ testTimeout: 20000 })
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import OnboardPage from '../app/onboard/page'
import { submitGoalAnswer } from '../lib/api'
import { answerField, answerValue, findAnswerField, findMathField, typeAnswer } from './helpers/answerInput'

// The onboard page routes to the results screen and back through the router.
vi.mock('next/navigation', () => ({
  useRouter: () => ({ push: vi.fn(), replace: vi.fn(), back: vi.fn(), refresh: vi.fn() }),
  usePathname: () => '/onboard',
  useSearchParams: () => new URLSearchParams(),
}))

// The math editor's chunk, held until the test lets it through.
//
// `MathInput` shows a working plain input while the MathLive chunk loads, and on a
// phone that window is long enough to type an answer in (see `PlainAnswerInput`).
// Under vitest the chunk resolves in a few milliseconds, so without this gate the
// loading input is gone before a test can type into it, and the bug below, which
// lives in that input, is unreachable.
const editorChunk = vi.hoisted(() => {
  let release = () => {}
  const gate = new Promise<void>(resolve => {
    release = resolve
  })
  return { gate, release: () => release() }
})

vi.mock('../components/math/MathLiveField', async (importOriginal) => {
  await editorChunk.gate
  return importOriginal()
})

const PROG0 = {
  answered: 0,
  estimated_total: 25,
  min_total: 25,
  max_total: 45,
  probe_budget: 45,
  done: false,
}

// Corpus concept ids, so the question is served with the math editor the live
// site renders. An id the corpus does not know falls back to the plain input.
vi.mock('../lib/api', async (importOriginal) => {
  const mod = await importOriginal<typeof import('../lib/api')>()
  return {
    ...mod,
    startGoalDiagnostic: vi.fn(async () => ({
      session_id: 's1',
      concept_id: 'arith.add.single',
      concept_name: 'Single-digit addition',
      question: 'Q1 text',
      done: false,
      grading_type: 'numeric',
      progress: PROG0,
    })),
    submitGoalAnswer: vi.fn(async () => ({
      done: false,
      correct: false,
      concept_id: 'arith.add.double',
      concept_name: 'Two-digit addition',
      question: 'Q2 text',
      grading_type: 'numeric',
      progress: { ...PROG0, answered: 1 },
    })),
    resumeGoalDiagnostic: vi.fn(),
    getGoalPlan: vi.fn(),
    skipGoalQuestion: vi.fn(),
    retryGoalQuestion: vi.fn(),
  }
})

describe('onboard diagnostic answer reset', () => {
  beforeEach(() => {
    sessionStorage.clear()
    localStorage.clear()
  })

  afterAll(() => {
    editorChunk.release()
  })

  async function beginDiagnostic() {
    render(<OnboardPage />)
    fireEvent.click(await screen.findByRole('button', { name: /Arithmetic/ }))
    fireEvent.click(screen.getByRole('button', { name: /Start diagnostic test/ }))
    fireEvent.click(await screen.findByRole('button', { name: /Begin diagnostic/ }))
    await screen.findByText('Q1 text')
  }

  it("starts the next question empty after I don't know, even while the editor is still loading", async () => {
    await beginDiagnostic()

    // The editor chunk is held, so the answer control is the loading input.
    const field = await findAnswerField()
    expect(field).toHaveAttribute('data-mathua-loading', 'true')
    typeAnswer(field, '7')
    expect(answerValue(field)).toBe('7')

    // "I don't know" serves the next question with the field still mounted.
    fireEvent.click(screen.getByRole('button', { name: "I don't know" }))
    await screen.findByText('Q2 text')
    expect(vi.mocked(submitGoalAnswer)).toHaveBeenCalledWith('s1', 'arith.add.single', '', expect.any(Number), true)

    // The learner is looking at a new question: the box must not still say 7.
    // The same loading input is still mounted while the reset effect runs.
    // Reading it immediately can catch the previous answer before React
    // clears it, which is the flake this waits out.
    await waitFor(() => expect(answerValue(answerField())).toBe(''))
    const next = answerField()

    // What they type for this question is the whole answer, and it is what the
    // editor takes over when its chunk finally arrives.
    typeAnswer(next, '8')
    editorChunk.release()
    const editor = await findMathField()
    expect(answerValue(editor)).toBe('8')
    expect(screen.getByRole('button', { name: 'Check Answer' })).toBeEnabled()
  })
})
