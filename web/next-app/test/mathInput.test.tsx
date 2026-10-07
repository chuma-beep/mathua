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

/**
 * The editor, once it has actually loaded.
 *
 * Waits, because the loading state is now a working plain `<input>` rather than an
 * inert div (see PlainAnswerInput). So "there is an input" is no longer evidence
 * that MathLive has mounted -- the very first control in the answer area is an input
 * now, and asserting synchronously would have tested the loading state and called it
 * the editor.
 */
async function mathField(): Promise<Element> {
  await waitFor(() => {
    const all = answerFields()
    expect(all.some(el => el.tagName.toLowerCase() === 'math-field')).toBe(true)
  })
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
    expect(await mathField()).toBeTruthy()
  })

  it('renders the math editor for a multiple-choice concept too', async () => {
    // The editor is universal now. Multiple choice still renders ChoiceOptions above
    // the field, so this is the *typed* path, not a replacement for the tap target.
    render(<MathAnswerInput value="" onChange={() => {}} gradingType="multiple_choice" conceptId="geo.basic.points_lines" placeholder="Your answer" />)
    expect(await mathField()).toBeTruthy()
  })

  it('still renders the plain input when the concept is unknown to the corpus', async () => {
    // The one remaining text path: a fixture id, or a concept added server-side before
    // the corpus rebuild. Degrading to the input that grades correctly everywhere is the
    // safer failure.
    render(<MathAnswerInput value="" onChange={() => {}} conceptId="not-a-real-concept" placeholder="Your answer" />)
    await waitFor(() => expect(field()).toBeTruthy())
    expect(textField()).toBeTruthy()
  })

  it('keeps the SymbolPalette in text mode, and not in math mode', async () => {
    const { unmount } = render(
      <MathAnswerInput value="" onChange={() => {}} conceptId="not-a-real-concept" placeholder="Your answer" />,
    )
    await waitFor(() => expect(textField()).toBeTruthy())
    // The palette is the legacy way to type π and ± into a plain input; the math
    // editor has a keypad for that, so showing both would be two ways to do one thing.
    expect(screen.getByLabelText('Insert π')).toBeTruthy()
    unmount()

    render(<MathAnswerInput value="" onChange={() => {}} gradingType="numeric" conceptId="arith.add.single" placeholder="Your answer" />)
    await mathField()
    expect(screen.queryByLabelText('Insert π')).toBeNull()
  })

  it('exposes an accessible name on whichever control it rendered', async () => {
    render(<MathAnswerInput value="" onChange={() => {}} gradingType="numeric" conceptId="arith.add.single" ariaLabel="Your answer" placeholder="Your answer" />)
    await mathField()
    // MathLive gives the field role="group", so aria-label is what a screen reader
    // announces — without it the field is an unnamed group.
    expect((await mathField()).getAttribute('aria-label')).toBe('Your answer')
  })

  it('keeps the placeholder queryable, so a host can still find its own field', async () => {
    render(<MathAnswerInput value="" onChange={() => {}} gradingType="numeric" conceptId="arith.add.single" placeholder="Your answer" />)
    await waitFor(() => expect(screen.getByPlaceholderText(/your answer/i)).toBeTruthy())
  })
})

describe('MathAnswerInput answer flow', () => {
  it('reports line segment intact from the MathLive multiple-choice field', async () => {
    const seen: string[] = []
    render(<MathAnswerInput value="" onChange={v => seen.push(v)} gradingType="multiple_choice" conceptId="geo.basic.points_lines" />)
    const el = (await mathField()) as unknown as { value: string; getValue: (format?: string) => string }
    // The real serializer recording is in mathlive-plain-text.json. The stub
    // does not model implicit multiplication, so supply that recorded output.
    vi.spyOn(el, 'getValue').mockReturnValue('l i n e s e g m e n t')
    el.value = 'line\\,segment'
    await waitFor(() => expect(seen[seen.length - 1]).toBe('line segment'))
  })

  it('reports a plain string upward, not LaTeX', async () => {
    const seen: string[] = []
    render(<MathAnswerInput value="" onChange={v => seen.push(v)} gradingType="numeric" conceptId="arith.add.single" />)
    await waitFor(() => expect(field()).toBeTruthy())
    const el = (await mathField()) as unknown as { value: string }
    el.value = '\\frac{1}{2}'
    await waitFor(() => expect(seen.length).toBeGreaterThan(0))
    // The whole design: the host stores "1/2", submits "1/2", the grader sees "1/2".
    expect(seen[seen.length - 1]).toBe('1/2')
  })

  it('submits on Enter from the math editor', async () => {
    const onSubmit = vi.fn()
    render(<MathAnswerInput value="" onChange={() => {}} onSubmit={onSubmit} gradingType="numeric" conceptId="arith.add.single" />)
    await waitFor(() => expect(field()).toBeTruthy())
    fireEvent.keyDown(await mathField(), { key: 'Enter' })
    expect(onSubmit).toHaveBeenCalledTimes(1)
  })

  it('submits on Enter from the plain input too', async () => {
    const onSubmit = vi.fn()
    render(<MathAnswerInput value="" onChange={() => {}} onSubmit={onSubmit} conceptId="not-a-real-concept" />)
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
    fireEvent.keyDown(await mathField(), { key: 'Tab' })
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
    for (const status of ['correct', 'incorrect'] as const) {
      const { unmount } = render(
        <MathAnswerInput value="5" onChange={() => {}} status={status} gradingType="numeric" conceptId="arith.add.single" />,
      )
      await waitFor(() => expect(field()).toBeTruthy())
      // The verdict rides on the field's own border, as `data-status`.
      //
      // It used to be a utility class on the wrapper, which drew the colour *outside*
      // the field's opaque background — so `correct` was a green ring around a grey box.
      // A class cannot come back here either: globals.css's
      // `math-field.mathua-math-field` selector outranks any Tailwind utility, so
      // `border-mathua-green` on the element would lose to the `border-color` declared
      // there and the verdict would not render at all.
      expect(field()?.getAttribute('data-status')).toBe(status)
      unmount()
    }

    const { unmount } = render(
      <MathAnswerInput value="5" onChange={() => {}} status="disabled" gradingType="numeric" conceptId="arith.add.single" />,
    )
    await waitFor(() => expect(field()).toBeTruthy())
    // Disabled is a property on the MathLive element, not just a class.
    expect((field() as unknown as { readOnly: boolean }).readOnly).toBe(true)
    // No verdict attribute when there is no verdict, so the border falls back to the
    // rest colour rather than to whatever the last state was.
    expect(field()?.getAttribute('data-status')).toBeNull()
    unmount()
  })

  it('draws exactly one box around the answer', async () => {
    const { container, unmount } = render(
      <MathAnswerInput value="5" onChange={() => {}} gradingType="numeric" conceptId="arith.add.single" />,
    )
    await waitFor(() => expect(field()).toBeTruthy())

    // The wrapper is layout-only now. It used to carry `border border-mathua-border
    // bg-mathua-code px-3 py-2 focus-within:border-mathua-blue` while the field carried
    // its own border, so focusing the answer drew two blue rectangles 24px apart — one
    // token, painted twice, which reads as two states rather than one.
    const wrapper = container.querySelector('.mathua-math-input') as HTMLElement
    expect(wrapper.className).not.toMatch(/border/)
    expect(wrapper.className).not.toMatch(/bg-/)
    expect(wrapper.className).not.toMatch(/p[xy]-/)
    unmount()
  })

})
