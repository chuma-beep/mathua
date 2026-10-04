// The self-check hint: shown, quiet, and structurally incapable of grading anything.
//
// The last point is the important one. `MathAnswerInput` exposes no path by which the
// hint could reach the submitted answer, and these tests hold that line: the hint's text
// appears in the DOM while the value the host receives is untouched by it.

import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, waitFor, act } from '@testing-library/react'
import { MathAnswerInput } from '../components/math/MathInput'
import { isFetchReasonable } from '../components/math/SelfCheck'
import { answerFields } from './helpers/answerInput'

const HINT = 'math-self-check'

function mathField(): HTMLElement {
  const all = answerFields()
  expect(all.length).toBeGreaterThan(0)
  return all[all.length - 1]
}

/** Type into the stubbed field. Long enough to clear the hint's debounce. */
async function type(page: { value: string }, text: string) {
  act(() => {
    page.value = text
  })
  await act(async () => {
    await new Promise(r => setTimeout(r, 1100))
  })
}

/**
 * The feedback contract, stated as an assertion rather than a comment.
 *
 * `equivalence()` is three-valued: 'equivalent', 'different', 'unknown'. The rule that
 * matters is the third one — nothing may render a verdict for 'unknown', or indeed for
 * any of the three, because the client never holds the expected answer to compare
 * against (ADR-005). These tests exercise all three states and assert the same thing
 * each time: no correctness language reaches the DOM.
 */
describe('the feedback contract', () => {
  const VERDICT_WORDS = /correct|incorrect|wrong|try again|not yet|right answer/i

  // One real pair per state, measured rather than assumed. These are the actual
  // results from `equivalence()` on CE 0.147.0:
  //
  //   x²+2x+1  vs (x+1)²      equivalent   (factored vs expanded)
  //   2^10     vs 1024         different    (isIdenticallyEqual says false)
  //   \sqrt{x²} vs x           unknown      (correctly refuses: |x| ≠ x)
  //
  // Note what is *not* in that table: `2x+3` vs `x²+2x+1` is 'unknown', not
  // 'different'. CE will not call two plainly unequal symbolic expressions different,
  // which is the whole reason 'unknown' has to be a first-class outcome rather than an
  // error path. And the LaTeX needs its backslashes — `sqrt(x^2)` parses as a product
  // of symbols, which is a different question entirely.
  const CASES: { state: 'equivalent' | 'different' | 'unknown'; latex: string; against: string }[] = [
    { state: 'equivalent', latex: 'x^2+2x+1', against: '(x+1)^2' },
    { state: 'different', latex: '2^10', against: '1024' },
    { state: 'unknown', latex: '\\sqrt{x^2}', against: 'x' },
  ]

  it('really does reach all three states, or the rest of this is theatre', async () => {
    const { loadEngine, equivalence } = await import('../components/math/equivalence')
    const engine = await loadEngine()
    expect(engine).not.toBeNull()
    const outcomes = new Set(CASES.map(c => equivalence(engine!, c.latex, c.against)))
    // All three must be reachable, or the DOM assertions below are testing one path
    // three times and calling it coverage.
    expect(outcomes.size).toBe(3)
  })

  for (const { state, latex } of CASES) {
    it(`renders no correctness feedback when equivalence would be '${state}'`, async () => {
      const { container } = render(
        <MathAnswerInput value="" onChange={() => {}} gradingType="numeric" conceptId="arith.add.single" />,
      )
      await waitFor(() => expect(container.querySelector('math-field')).toBeTruthy())
      await type(mathField() as unknown as { value: string }, latex)

      // The hint may appear, may not, and either is fine. What is not fine is a verdict.
      await act(async () => {
        await new Promise(r => setTimeout(r, 1200))
      })
      const rendered = document.body.textContent ?? ''
      expect(rendered, `state ${state} leaked a verdict into the DOM`).not.toMatch(VERDICT_WORDS)
    })
  }

  it('maps each case to the state it claims', async () => {
    const { loadEngine, equivalence } = await import('../components/math/equivalence')
    const engine = (await loadEngine())!
    const actual = CASES.map(c => equivalence(engine, c.latex, c.against))
    // Pinned rather than counted, so a change in how Compute Engine resolves any of
    // these is visible here rather than silently weakening the assertions above.
    expect(actual).toEqual(CASES.map(c => c.state))
  })
})

describe('isFetchReasonable', () => {
  function withConnection(value: unknown, run: () => void) {
    const nav = navigator as Navigator & { connection?: unknown }
    const had = 'connection' in nav
    const previous = nav.connection
    Object.defineProperty(nav, 'connection', { value, configurable: true })
    try {
      run()
    } finally {
      if (had) Object.defineProperty(nav, 'connection', { value: previous, configurable: true })
      else delete (nav as { connection?: unknown }).connection
    }
  }

  it('defaults to yes when the browser reports nothing', () => {
    withConnection(undefined, () => expect(isFetchReasonable()).toBe(true))
  })

  it('says no when the user asked for data saving', () => {
    withConnection({ saveData: true, effectiveType: '4g' }, () => expect(isFetchReasonable()).toBe(false))
  })

  it('says no on a slow or 2g connection', () => {
    withConnection({ saveData: false, effectiveType: 'slow-2g' }, () => expect(isFetchReasonable()).toBe(false))
    withConnection({ saveData: false, effectiveType: '2g' }, () => expect(isFetchReasonable()).toBe(false))
  })

  it('says yes on a normal or fast connection', () => {
    withConnection({ saveData: false, effectiveType: '4g' }, () => expect(isFetchReasonable()).toBe(true))
    withConnection({ saveData: false, effectiveType: '3g' }, () => expect(isFetchReasonable()).toBe(true))
  })
})

describe('SelfCheck', () => {
  beforeEach(() => {
    vi.useRealTimers()
  })

  it('stays silent until there is arithmetic worth checking', async () => {
    const { container } = render(
      <MathAnswerInput value="" onChange={() => {}} gradingType="numeric" conceptId="arith.add.single" />,
    )
    await waitFor(() => expect(container.querySelector('math-field')).toBeTruthy())
    await type(mathField() as unknown as { value: string }, 'x')
    // A lone symbol has nothing to simplify, so the engine is never even asked.
    expect(screen.queryByTestId(HINT)).toBeNull()
  })

  it('offers the simplification of the learner own expression', async () => {
    const { container } = render(
      <MathAnswerInput value="" onChange={() => {}} gradingType="numeric" conceptId="arith.add.single" />,
    )
    await waitFor(() => expect(container.querySelector('math-field')).toBeTruthy())
    await type(mathField() as unknown as { value: string }, 'x^2+x^2')
    // Two x² really are 2x², and saying so is the entire feature.
    await waitFor(() => expect(screen.getByTestId(HINT)).toBeTruthy(), { timeout: 5000 })
    expect(screen.getByTestId(HINT).textContent).toContain('2x^2')
  })

  it('does not colour or badge the hint like a verdict', async () => {
    const { container } = render(
      <MathAnswerInput value="" onChange={() => {}} gradingType="numeric" conceptId="arith.add.single" />,
    )
    await waitFor(() => expect(container.querySelector('math-field')).toBeTruthy())
    await type(mathField() as unknown as { value: string }, 'x^2+x^2')
    await waitFor(() => expect(screen.getByTestId(HINT)).toBeTruthy(), { timeout: 5000 })
    // No green, no red, no cross: it must not read as something the grader said.
    const cls = screen.getByTestId(HINT).getAttribute('class') ?? ''
    expect(cls).not.toMatch(/green|red|blue/)
    expect(screen.getByTestId(HINT).textContent).not.toMatch(/correct|incorrect|right|wrong/i)
  })

  it('never changes the answer the host receives', async () => {
    const seen: string[] = []
    const { container } = render(
      <MathAnswerInput value="" onChange={v => seen.push(v)} gradingType="numeric" conceptId="arith.add.single" />,
    )
    await waitFor(() => expect(container.querySelector('math-field')).toBeTruthy())
    await type(mathField() as unknown as { value: string }, 'x^2+x^2')
    await waitFor(() => expect(screen.getByTestId(HINT)).toBeTruthy(), { timeout: 5000 })
    // The hint appeared; the submitted value is still exactly what was typed. This is
    // the assertion that keeps the feature non-binding.
    expect(seen[seen.length - 1]).toBe('x^2+x^2')
  })

  it('does not appear for a text-mode question', async () => {
    render(<MathAnswerInput value="" onChange={() => {}} gradingType="multiple_choice" conceptId="c1" />)
    await waitFor(() => expect(answerFields().length).toBeGreaterThan(0))
    expect(screen.queryByTestId(HINT)).toBeNull()
  })
})