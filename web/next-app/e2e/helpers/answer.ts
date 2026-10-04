// Driving the answer field in the browser.
//
// Two controls: the MathLive `<math-field>` for expression-valued concepts and the
// plain `<input>` for multiple-choice, tuple, ordering and comparison. A test that
// means to "answer the question" should not have to know which one it got.
//
// Selectors here are against MathLive's real DOM (`.MLK__*` is the virtual
// keyboard's prefix, `.ML__*` the field's), not against anything invented.

import { expect, type Locator, type Page } from '@playwright/test'

const ANSWER = 'math-field[placeholder*="Your answer"], input[placeholder*="Your answer"]'

export function answerField(page: Page): Locator {
  return page.locator(ANSWER).first()
}

export async function answerLatex(page: Page): Promise<string> {
  return answerField(page).evaluate(el => (el as unknown as { value: string }).value)
}

/**
 * Put the caret in the answer field and type.
 *
 * MathLive focuses its contenteditable sink on `mousedown`. A single atomic
 * `locator.click()` — mousedown and mouseup in the same tick — leaves the field
 * unfocused here, so the press is split with a gap. `.focus()` is used as the
 * belt-and-braces path so a test never silently types into nothing.
 *
 * Real keystrokes, not `fill()`: the editor is not an `<input>`, and typing is the
 * only way to exercise MathLive's parser. `2^3` has to come out as a power with a
 * navigable exponent slot, which is the whole point of using the library.
 */
export async function typeAnswer(page: Page, value: string): Promise<void> {
  await focusAnswerField(page)
  await typeIntoField(page, value)
}

/**
 * Keystroke delay.
 *
 * MathLive parses as it types, on its own schedule. Synthetic keys fired back to back
 * race that and get dropped — "2+2" would reliably arrive as "2". A human types slower
 * than a test does, so the delay is the honest way to reproduce real input rather
 * than a workaround for MathLive being fast.
 */
const TYPE_DELAY_MS = 40

/** Type into the already-focused field. */
export async function typeIntoField(page: Page, value: string): Promise<void> {
  if (value.length === 0) return
  await page.keyboard.type(value, { delay: TYPE_DELAY_MS })
  await page.waitForTimeout(TYPE_DELAY_MS)
}

/** Move focus into the field the way a pointer would. */
export async function focusAnswerField(page: Page): Promise<void> {
  const field = answerField(page)
  const box = await field.boundingBox()
  if (box) {
    await page.mouse.move(box.x + box.width / 2, box.y + box.height / 2)
    await page.mouse.down()
    await page.waitForTimeout(50)
    await page.mouse.up()
  }
  await field.evaluate(el => (el as unknown as { focus: () => void }).focus())
  await page.waitForTimeout(50)
}

/** Whether the field currently holds the browser's focus. */
export async function answerHasFocus(page: Page): Promise<boolean> {
  return answerField(page).evaluate(el => el === document.activeElement)
}

/** Set an answer without typing, for cases where only the resulting value matters. */
export async function setAnswer(page: Page, latex: string): Promise<void> {
  await answerField(page).evaluate((el, v) => {
    ;(el as unknown as { value: string }).value = v
  }, latex)
}

/** Wait for the editor to be mounted and ready for input. */
export async function waitForAnswerField(page: Page): Promise<Locator> {
  const field = answerField(page)
  await expect(field).toBeVisible({ timeout: 30_000 })
  return field
}

/**
 * Capture the `answer` the wrapper submits, which is the string the server grades.
 *
 * Read from the request rather than from the field: the field holds LaTeX, the
 * server receives plain text, and the gap between them is the thing worth testing.
 */
export function captureSubmittedAnswer(page: Page, urlFragment: string): () => string {
  let captured: string | null = null
  page.on('request', req => {
    if (req.method() !== 'POST' || !req.url().includes(urlFragment)) return
    try {
      const body = JSON.parse(req.postData() ?? '{}') as { answer?: unknown }
      if (typeof body.answer === 'string') captured = body.answer
    } catch {
      // Not the request under test.
    }
  })
  return () => {
    if (captured === null) throw new Error(`no POST to ${urlFragment} carrying an answer was observed`)
    return captured
  }
}

/**
 * Open the virtual keyboard with Mathua's layout installed.
 *
 * Focusing first is not optional: the wrapper installs the layout for the current
 * concept in its `onFocus` handler, so showing the keyboard on an unfocused field
 * raises MathLive's own default layout instead — and on an arithmetic exercise that
 * default is the very thing this work removes.
 */
export async function showVirtualKeyboard(page: Page): Promise<void> {
  await focusAnswerField(page)
  await page.evaluate(() => {
    const kb = (window as unknown as { mathVirtualKeyboard?: { show: () => void } }).mathVirtualKeyboard
    kb?.show()
  })
  await expect(keyboardRows(page)).toBeVisible({ timeout: 10_000 })
}

/**
 * The rows of the *visible* keyboard layer.
 *
 * MathLive keeps every layout's rows in the DOM and hides the inactive ones, so an
 * unscoped `.MLK__rows` matches four or five elements and every assertion against it
 * is ambiguous.
 */
export function keyboardRows(page: Page): Locator {
  return page.locator('.MLK__rows:visible').first()
}

/**
 * A visible keycap by its displayed label.
 *
 * Matching is on the keycap's full text, because `hasText` is a substring match and
 * "7" would otherwise also match "78", "÷" and every digit row. A keycap's text
 * includes its `aside` caption — the √ key renders as "√\nsquare root" — so pass the
 * leading glyph.
 */
/**
 * Every keycap in the keypad, including the ones that are not `.MLK__keycap`.
 *
 * MathLive's `renderKeycap` adds `MLK__keycap` only when the keycap's class does
 * *not* already contain `separator`, `action`, `shift`, `fnbutton` or
 * `bigfnbutton`. So backspace, the caret arrows, the dismiss key and every
 * separator are all rendered without it. A `.MLK__keycap` selector silently skips
 * them, which is how a broken backspace key passed the touch-target audit for the
 * whole of commit 1: the audit was measuring only the keys that were working.
 *
 * The row's direct children are the honest superset. `:visible` matters: MathLive
 * renders every layer and hides the inactive ones, so without it the measurement
 * collects keycaps at 0x0 from layers the learner cannot see.
 */
export function keycaps(page: Page): Locator {
  return page.locator('.MLK__row > *:visible')
}

/**
 * A keycap drawn as one of MathLive's built-in SVG glyphs, e.g. `delete-backward`.
 *
 * Glyph keycaps have no text at all, so `hasText` cannot find them; they are
 * identified by the `use` element they reference.
 */
export function glyphKeycap(page: Page, glyph: string): Locator {
  return keycaps(page).filter({
    has: page.locator(`[*|href="#svg-${glyph}"], [href="#svg-${glyph}"]`),
  })
}

export function keycap(page: Page, label: string): Locator {
  return keycaps(page)
    .filter({ hasText: new RegExp(`^${label.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')}`) })
    .first()
}

/** Press a keycap and wait for MathLive to process the insertion. */
export async function pressKeycap(page: Page, label: string): Promise<void> {
  await keycap(page, label).click()
  await page.waitForTimeout(60)
}
/**
 * Run one of MathLive's own editing commands, e.g. `moveToMathfieldStart`.
 *
 * The public alternative to clicking at a position: setting `.value` does not move
 * the caret, and MathLive's caret commands are the documented way to place it.
 */
export async function runMathCommand(page: Page, command: string): Promise<void> {
  await answerField(page).evaluate(
    (el, c) => (el as unknown as { executeCommand: (cmd: string[]) => void }).executeCommand(c),
    [command],
  )
  await page.waitForTimeout(30)
}

/** The whole keypad panel: the `.MLK__rows` wrapper. */
export function keyboardPanel(page: Page): Locator {
  return page.locator('.MLK__rows').first().locator('..')
}

/** Whether MathLive considers the keypad visible, per its own API. */
export async function keyboardVisible(page: Page): Promise<boolean> {
  return page.evaluate(() => {
    const kb = (window as unknown as { mathVirtualKeyboard?: { visible: boolean } }).mathVirtualKeyboard
    return kb?.visible === true
  })
}

/** Hide the keypad through MathLive's own API. */
export async function hideVirtualKeyboard(page: Page): Promise<void> {
  await page.evaluate(() => {
    const kb = (window as unknown as { mathVirtualKeyboard?: { hide: () => void } }).mathVirtualKeyboard
    kb?.hide()
  })
}
