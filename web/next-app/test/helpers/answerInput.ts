// Driving the answer field in tests, whichever control the question chose.
//
// `MathAnswerInput` renders one of two things: the MathLive `<math-field>` for
// expression-valued concepts (numeric, symbolic, expression, polynomial,
// complex) or the plain `<Input>` Mathua always had for multiple-choice, tuple,
// ordering and comparison. A test that wants to "answer the question" should not
// have to know which one it got.
//
// This is a driver, not an assertion helper: it replaces `fireEvent.change`, which
// only works on a native input. Under the stubbed MathLive element, assigning
// `.value` fires `input` exactly as a keystroke does, so the React onChange path
// under test is the same one the browser runs.
//
// What the real library does inside the field — fractions as structures, caret
// navigation, the virtual keypad — is verified in e2e/math-input.spec.ts, which
// runs Chromium. Nothing here is a substitute for that.

import { act, fireEvent, waitFor } from '@testing-library/react'

type Field = Element

const PLACEHOLDER = /your answer/i

/**
 * Every answer control currently on screen, oldest first.
 *
 * Matched on the accessible name rather than the tag: the two controls differ
 * (`<input>` vs `<math-field>`), and a host labels the answer "Your answer" through
 * either `placeholder` or `aria-label` depending on which control it got.
 */
export function answerFields(): HTMLElement[] {
  return Array.from(document.querySelectorAll('[placeholder], [aria-label]')).filter(el => {
    const name = `${el.getAttribute('placeholder') ?? ''} ${el.getAttribute('aria-label') ?? ''}`
    return PLACEHOLDER.test(name)
  }) as HTMLElement[]
}

/** The answer control currently on screen. Throws if none is rendered. */
export function answerField(): HTMLElement {
  const all = answerFields()
  if (all.length === 0) throw new Error('no answer field rendered')
  return all[all.length - 1]
}

/**
 * The answer control, once it exists.
 *
 * The math editor is loaded lazily, so on the first render after a question is
 * served the field is not in the DOM yet — it is behind a Suspense boundary while
 * the MathLive chunk resolves. A synchronous query right after the question
 * appears therefore finds nothing, and would be flaky rather than deterministic if
 * the chunk happened to be cached.
 */
export async function findAnswerField(): Promise<HTMLElement> {
  await waitFor(() => {
    if (answerFields().length === 0) throw new Error('no answer field rendered yet')
  })
  return answerField()
}

export function isMathField(el: Field): boolean {
  return el.tagName.toLowerCase() === 'math-field'
}

/**
 * Type an answer. Uses `fireEvent.change` for the plain input and a `.value`
 * assignment for `<math-field>`, which the stub turns into an `input` event.
 */
export function typeAnswer(el: Field, value: string): void {
  // `act` on both paths. Assigning `.value` is a plain assignment rather than a
  // dispatched event, so React would not flush the resulting state update before
  // the next interaction — the click on Check would read a stale answer.
  act(() => {
    if (isMathField(el)) {
      ;(el as unknown as { value: string }).value = value
      return
    }
    fireEvent.change(el, { target: { value } })
  })
}

/**
 * The question card that owns an answer field.
 *
 * Needed because the two controls have different ancestries: the plain `<input>`
 * sits directly in the card, while the math editor wraps itself in its own
 * `div` whose class also contains "border". A nearest-ancestor search on that
 * substring therefore stops at the editor and misses the card's submit button.
 * The card is identified by its own surface class instead.
 */
export function answerCard(el: Field): HTMLElement {
  const card = el.closest('div[class*="bg-mathua-surface"]')
  if (!card) throw new Error('no question card found for the answer field')
  return card as HTMLElement
}

/** The answer the control currently holds. */
export function answerValue(el: Field): string {
  if (isMathField(el)) return (el as unknown as { value: string }).value
  return (el as HTMLInputElement).value
}