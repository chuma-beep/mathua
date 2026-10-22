// What the answer hosts actually get.
//
// The contract these tests protect is the one the five hosts rely on: a host holds a
// plain string, hands it in, and gets a plain string back — no LaTeX, no MathJSON,
// no knowledge of which control rendered. Everything MathLive-specific is asserted
// here so a host never has to.

import { useState } from 'react'
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

  it('clears the editor when the host clears its answer', async () => {
    // The diagnostic and quiz hosts keep this control mounted across questions and
    // set their plain answer to "" before revealing the next one. The editor's LaTeX
    // lived in a ref that only the field itself wrote, so the previous answer stayed
    // on screen, disabled, above the new question.
    function Host() {
      const [answer, setAnswer] = useState('')
      return (
        <>
          <MathAnswerInput value={answer} onChange={setAnswer} gradingType="numeric" conceptId="arith.add.single" />
          <button type="button" onClick={() => setAnswer('')}>clear</button>
        </>
      )
    }

    render(<Host />)
    await waitFor(() => expect(field()).toBeTruthy())
    const el = (await mathField()) as unknown as { value: string }
    el.value = '\\frac{1}{2}'
    await waitFor(() => expect(el.value).toBe('\\frac{1}{2}'))

    fireEvent.click(screen.getByRole('button', { name: 'clear' }))
    await waitFor(() => expect(el.value).toBe(''))
  })

  it('does not submit on Enter while the answer is being checked', async () => {
    // The Enter listener is native and bypasses the disabled Check button, so a second
    // Enter during the request sent the same answer twice.
    const onSubmit = vi.fn()
    render(<MathAnswerInput value="4" onChange={() => {}} onSubmit={onSubmit} disabled gradingType="numeric" conceptId="arith.add.single" />)
    fireEvent.keyDown(await mathField(), { key: 'Enter' })
    expect(onSubmit).not.toHaveBeenCalled()
  })

  it('puts the keypad away once the answer is checked', async () => {
    // The keypad stayed up over the feedback and the Next button on a phone, so the
    // page looked frozen after the first answer.
    const kb = window.mathVirtualKeyboard as unknown as { visible: boolean; show: () => void }
    const { rerender } = render(<MathAnswerInput value="4" onChange={() => {}} gradingType="numeric" conceptId="arith.add.single" />)
    const el = (await mathField()) as HTMLElement
    el.focus()
    kb.show()
    rerender(<MathAnswerInput value="4" onChange={() => {}} disabled gradingType="numeric" conceptId="arith.add.single" />)
    await waitFor(() => expect(kb.visible).toBe(false))
    expect(document.activeElement).not.toBe(el)
  })

  it('keeps the caret when the editor replaces the loading input', async () => {
    // The plain input shown while MathLive loads is removed when it arrives. Focus went
    // with it, so a learner typing their first answer lost the rest of their keystrokes.
    // A fresh module, because `lazy` remembers a resolved chunk and earlier tests in
    // this file have already loaded it, so there would be no loading phase to observe.
    vi.resetModules()
    const { MathAnswerInput: FreshInput } = await import('../components/math/MathInput')
    render(<FreshInput value="" onChange={() => {}} gradingType="numeric" conceptId="arith.add.single" />)
    const loading = document.querySelector('input[data-mathua-loading="true"]') as HTMLInputElement | null
    expect(loading).not.toBeNull()
    loading?.focus()
    const el = await mathField()
    await waitFor(() => expect(document.activeElement).toBe(el))
  })

  it('keeps a typed run of digits in the exponent', async () => {
    // MathLive leaves an empty exponent after its first digit, so `2^10` reached the
    // grader as `2^1 0`. The second digit is claimed and put back into the exponent;
    // anything that is not a digit is left to MathLive. Which slot it lands in is
    // asserted against the real library in e2e/math-input.spec.ts.
    render(<MathAnswerInput value="" onChange={() => {}} gradingType="numeric" conceptId="arith.add.single" />)
    const el = (await mathField()) as HTMLElement & { value: string }
    el.value = '2^1'
    const press = (key: string) => {
      const e = new KeyboardEvent('keydown', { key, bubbles: true, cancelable: true })
      el.dispatchEvent(e)
      return e.defaultPrevented
    }
    expect(press('1')).toBe(false)
    expect(press('0')).toBe(true)
    expect(el.value).toBe('2^10')
    expect(press('+')).toBe(false)
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
