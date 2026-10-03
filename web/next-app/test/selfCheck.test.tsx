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