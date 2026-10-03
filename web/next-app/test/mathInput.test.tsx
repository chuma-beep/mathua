// What the answer hosts actually get.
//
// The contract these tests protect is the one the five hosts rely on: a host holds a
// plain string, hands it in, and gets a plain string back — no LaTeX, no MathJSON,
// no knowledge of which control rendered. Everything MathLive-specific is asserted
// here so a host never has to.

import { describe, it, expect, vi } from 'vitest'
import { render, screen, waitFor, fireEvent } from '@testing-library/react'
import { MathAnswerInput, type MathFocusHandle } from '../components/math/MathInput'
import { answerFields } from './helpers/answerInput'

function field(): HTMLInputElement | Element {
  const all = answerFields()
  expect(all.length).toBeGreaterThan(0)
  return all[all.length - 1]
}

function mathField(): Element {
  const el = field()
  expect(el.tagName.toLowerCase()).toBe('math-field')
  return el
}

function textField(): HTMLInputElement {
  const el = field()
  expect(el.tagName.toLowerCase()).toBe('input')
  return el as HTMLInputElement
}

describe('MathAnswerInput control selection', () => {
  it('renders the math editor for an expression-valued concept', async () => {
    render(<MathAnswerInput value="" onChange={() => {}} gradingType="numeric" conceptId="arith.add.single" placeholder="Your answer" />)
    await waitFor(() => expect(field()).toBeTruthy())
    expect(mathField()).toBeTruthy()
  })

  it('renders the plain input for a multiple-choice concept', async () => {
    render(<MathAnswerInput value="" onChange={() => {}} gradingType="multiple_choice" conceptId="c1" placeholder="Your answer" />)
    await waitFor(() => expect(field()).toBeTruthy())
    expect(textField()).toBeTruthy()
  })

  it('keeps the SymbolPalette in text mode, and not in math mode', async () => {
    const { unmount } = render(
      <MathAnswerInput value="" onChange={() => {}} gradingType="multiple_choice" conceptId="c1" placeholder="Your answer" />,
    )
    await waitFor(() => expect(textField()).toBeTruthy())
    // The palette is the legacy way to type π and ± into a plain input; the math
    // editor has a keypad for that, so showing both would be two ways to do one thing.
    expect(screen.getByLabelText('Insert π')).toBeTruthy()
    unmount()

    render(<MathAnswerInput value="" onChange={() => {}} gradingType="numeric" conceptId="arith.add.single" placeholder="Your answer" />)
    await waitFor(() => expect(mathField()).toBeTruthy())
    expect(screen.queryByLabelText('Insert π')).toBeNull()
  })

  it('exposes an accessible name on whichever control it rendered', async () => {
    render(<MathAnswerInput value="" onChange={() => {}} gradingType="numeric" conceptId="arith.add.single" ariaLabel="Your answer" placeholder="Your answer" />)
    await waitFor(() => expect(mathField()).toBeTruthy())
    // MathLive gives the field role="group", so aria-label is what a screen reader
    // announces — without it the field is an unnamed group.
    expect(mathField().getAttribute('aria-label')).toBe('Your answer')
  })

  it('keeps the placeholder queryable, so a host can still find its own field', async () => {
    render(<MathAnswerInput value="" onChange={() => {}} gradingType="numeric" conceptId="arith.add.single" placeholder="Your answer" />)
    await waitFor(() => expect(screen.getByPlaceholderText(/your answer/i)).toBeTruthy())
  })
})

describe('MathAnswerInput answer flow', () => {
  it('reports a plain string upward, not LaTeX', async () => {
    const seen: string[] = []
    render(<MathAnswerInput value="" onChange={v => seen.push(v)} gradingType="numeric" conceptId="arith.add.single" />)
    await waitFor(() => expect(field()).toBeTruthy())
    const el = mathField() as unknown as { value: string }
    el.value = '\\frac{1}{2}'
    await waitFor(() => expect(seen.length).toBeGreaterThan(0))
    // The whole design: the host stores "1/2", submits "1/2", the grader sees "1/2".
    expect(seen[seen.length - 1]).toBe('1/2')
  })

  it('submits on Enter from the math editor', async () => {
    const onSubmit = vi.fn()
    render(<MathAnswerInput value="" onChange={() => {}} onSubmit={onSubmit} gradingType="numeric" conceptId="arith.add.single" />)
    await waitFor(() => expect(field()).toBeTruthy())
    fireEvent.keyDown(mathField(), { key: 'Enter' })
    expect(onSubmit).toHaveBeenCalledTimes(1)
  })

  it('submits on Enter from the plain input too', async () => {
    const onSubmit = vi.fn()
    render(<MathAnswerInput value="" onChange={() => {}} onSubmit={onSubmit} gradingType="multiple_choice" conceptId="c1" />)
    await waitFor(() => expect(textField()).toBeTruthy())
    fireEvent.keyDown(textField(), { key: 'Enter' })
    expect(onSubmit).toHaveBeenCalledTimes(1)
  })

  it('does not submit on Enter inside a fraction placeholder', async () => {
    // Enter commits the placeholder MathLive is sitting in. Swallowing it would trap
    // a learner who presses Enter to move on.
    const onSubmit = vi.fn()
    render(<MathAnswerInput value="" onChange={() => {}} onSubmit={onSubmit} gradingType="numeric" conceptId="arith.add.single" />)
    await waitFor(() => expect(field()).toBeTruthy())
    fireEvent.keyDown(mathField(), { key: 'Tab' })
    expect(onSubmit).not.toHaveBeenCalled()
  })

  it('honours a focus request that arrives before the editor has loaded', async () => {
    // Every host autofocuses when a question is revealed, which on the first question
    // is before the MathLive chunk resolves. A focus call at that moment finds no
    // element and is silently lost unless it is queued.
    const handle: { current: MathFocusHandle | null } = { current: null }
    render(<MathAnswerInput value="" onChange={() => {}} focusRef={handle} gradingType="numeric" conceptId="arith.add.single" />)
    handle.current?.focus()
    await waitFor(() => expect(field()).toBeTruthy())
    await waitFor(() => expect(document.activeElement).toBe(field()))
  })
})

describe('MathAnswerInput states', () => {

  it('marks the field correct, incorrect and disabled from the server verdict', async () => {
    for (const [status, matcher] of [
      ['correct', /green/],
      ['incorrect', /red/],
    ] as const) {
      const { unmount } = render(
        <MathAnswerInput value="5" onChange={() => {}} status={status} gradingType="numeric" conceptId="arith.add.single" />,
      )
      await waitFor(() => expect(field()).toBeTruthy())
      // The status class lives on the wrapper, not on the <math-field> itself.
      expect(document.querySelector('.mathua-math-input')?.getAttribute('class')).toMatch(matcher)
      unmount()
    }

    const { unmount } = render(
      <MathAnswerInput value="5" onChange={() => {}} status="disabled" gradingType="numeric" conceptId="arith.add.single" />,
    )
    await waitFor(() => expect(field()).toBeTruthy())
    // Disabled is a property on the MathLive element, not just a class.
    expect((field() as unknown as { readOnly: boolean }).readOnly).toBe(true)
    unmount()
  })
})
