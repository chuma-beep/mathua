import { test, expect, type Page } from '@playwright/test'
import {
  answerField,
  answerLatex,
  answerHasFocus,
  typeAnswer,
  focusAnswerField,
  setAnswer,
  waitForAnswerField,
  captureSubmittedAnswer,
  showVirtualKeyboard,
  pressKeycap,
  runMathCommand,
  keyboardRows,
  typeIntoField,
} from './helpers/answer'

// The real MathLive, in a real browser.
//
// The vitest suite stubs the library out — it cannot mount in jsdom — and asserts
// Mathua's behaviour: which control is chosen, what is submitted, what the learner
// sees. Everything here is the other half: whether the editor actually builds
// mathematical *structures*, and whether a learner's editing gestures work. Those are
// properties of MathLive, and stubbing them would assert nothing.
//
// Every expectation in this file was established by running it. Where MathLive's
// behaviour is surprising it is asserted as observed and called out, rather than
// asserted as one would expect it to behave.

const CONCEPT = 'arith.add.single'

const KP = {
  concept_id: CONCEPT,
  kps: [{ label: 'Add', section: 'Add', subgoals: [], worked_example: '2+3=5' }],
}

const PRACTICE = {
  concept_id: CONCEPT,
  questions: [
    { question: '2 + 2 = ?', answer: '4', explanation: 'Two and two.', source: 'curated' },
    { question: '3 + 5 = ?', answer: '8', explanation: 'Three and five.', source: 'curated' },
  ],
}

async function openLearn(page: Page) {
  await page.route('**/api/lessons/**', route => {
    const url = route.request().url()
    if (url.includes('/practice')) return route.fulfill({ json: PRACTICE })
    if (url.includes('/kp')) return route.fulfill({ json: KP })
    if (url.includes('/readiness')) return route.fulfill({ json: { concept_id: CONCEPT, ready: true, weak: [], missing: [] } })
    return route.fulfill({ json: {} })
  })
  await page.route('**/api/scores/**', route =>
    route.fulfill({ json: { lifetime_points: 0, weekly_score: 0, speed_bonus: 0, concepts_mastered: 1, current_streak: 0, level: 'Novice', xp_total: 0, xp_today: 0, daily_xp_goal: 10, spaced_reps: {}, avg_learning_time: 1 } }),
  )
  for (const [pattern, body] of [
    ['**/api/activity**', []],
    ['**/api/progress/**', {}],
    ['**/api/weaknesses**', { by_domain: {} }],
    ['**/api/reviews/due**', { count: 0 }],
    ['**/api/courses**', { courses: [] }],
    ['**/api/transcript**', { courses: [] }],
  ] as const) {
    await page.route(pattern, route => route.fulfill({ json: body }))
  }

  await page.goto(`/learn?concept=${CONCEPT}`)
  await page.getByRole('button', { name: 'Next →' }).click()
  await waitForAnswerField(page)
}

/** Clear the field and type into it, so each case starts from a known state. */
async function fresh(page: Page, typed: string): Promise<void> {
  await setAnswer(page, '')
  await focusAnswerField(page)
  await typeIntoField(page, typed)
}

test.beforeEach(async ({ page }) => {
  await page.addInitScript(() => {
    localStorage.setItem('mathua_token', 'fake-token')
    localStorage.setItem(
      'mathua_user',
      JSON.stringify({ student_id: 's1', name: 'Tester', username: 'tester', concepts_mastered: 3, current_streak: 0, level: 'Novice', diagnostic_completed: true }),
    )
  })
})

test.describe('taking focus', () => {
  test('a pointer press puts the caret in the field', async ({ page }) => {
    await openLearn(page)
    await focusAnswerField(page)
    expect(await answerHasFocus(page)).toBe(true)
  })

  test('Tab reaches the field, so it is in the tab order', async ({ page }) => {
    // Keyboard reachability is an accessibility requirement, not a nicety: the field
    // is a focusable custom element and must not be skipped.
    await openLearn(page)
    for (let i = 0; i < 8 && !(await answerHasFocus(page)); i++) {
      await page.keyboard.press('Tab')
    }
    expect(await answerHasFocus(page)).toBe(true)
  })

  test('the field carries an accessible name', async ({ page }) => {
    await openLearn(page)
    await expect(answerField(page)).toHaveAttribute('aria-label', 'Your answer')
  })
})

test.describe('arithmetic on the physical keyboard', () => {
  test('2 + 3', async ({ page }) => {
    await openLearn(page)
    await fresh(page, '2+3')
    expect(await answerLatex(page)).toBe('2+3')
  })

  test('12 × 4 becomes multiplication, not a glyph', async ({ page }) => {
    await openLearn(page)
    await fresh(page, '12×4')
    // The character × must serialize to \times, because grading/sympy_service.py
    // rejects the raw glyph outright.
    expect(await answerLatex(page)).toBe('12\\times4')
  })

  test('7 ÷ 2 becomes division', async ({ page }) => {
    await openLearn(page)
    await fresh(page, '7÷2')
    expect(await answerLatex(page)).toBe('7\\div2')
  })
})

test.describe('structures, not strings', () => {
  test('1/2 is a fraction with numerator and denominator', async ({ page }) => {
    await openLearn(page)
    await fresh(page, '1/2')
    const latex = await answerLatex(page)
    expect(latex).toContain('\\frac')
    // Rendered as a fraction, not as a slash between two numbers.
    await expect(page.locator('math-field >> .ML__mfrac').first()).toBeVisible()
  })

  test('a fraction key builds a fraction with two empty slots', async ({ page }) => {
    await openLearn(page)
    await showVirtualKeyboard(page)
    await pressKeycap(page, 'a⁄b')
    await expect(page.locator('math-field >> .ML__mfrac').first()).toBeVisible()
    // The caret lands in the numerator, so typing continues the fraction.
    await typeIntoField(page, '3')
    await page.keyboard.press('Tab')
    await typeIntoField(page, '4')
    // Single-digit arguments serialize without braces.
    expect(await answerLatex(page)).toContain('\\frac34')
  })

  test('x² is a power whose exponent is its own slot', async ({ page }) => {
    await openLearn(page)
    await fresh(page, 'x^2')
    expect(await answerLatex(page)).toContain('x^2')
    // The exponent is navigable: one ArrowLeft puts the caret inside it, so typing
    // extends the exponent rather than the expression.
    await page.keyboard.press('ArrowLeft')
    await typeIntoField(page, 'n+')
    expect(await answerLatex(page)).toMatch(/x\^\{?2n\+\}?/)
  })

  test('the caret reaches the denominator of a fraction', async ({ page }) => {
    await openLearn(page)
    await fresh(page, '\\frac{1}{2}')
    await page.keyboard.press('ArrowLeft')
    await page.keyboard.press('ArrowLeft')
    await typeIntoField(page, '9')
    // 2 became 92: the 9 went into the denominator, not the end of the expression.
    expect(await answerLatex(page)).toContain('92')
  })

  test('the radical key builds a radical', async ({ page }) => {
    await openLearn(page)
    await setAnswer(page, '')
    await focusAnswerField(page)
    await showVirtualKeyboard(page)
    await pressKeycap(page, '√')
    await typeIntoField(page, '9')
    // MathLive writes a single-character radicand without braces.
    expect(await answerLatex(page)).toContain('\\sqrt9')
    // …and it is a radical, not the √ symbol: \\surd would have the radicand outside.
    expect(await answerLatex(page)).not.toContain('\\surd')
  })

  test('typing the √ glyph is not a radical', async ({ page }) => {
    // Recorded, not assumed. MathLive maps the √ character to `\\surd`, which is
    // an ordinary symbol, so typing "√9" gives "\\surd9" and not "\\sqrt{9}". This is
    // exactly why the keypad offers a structural √ key rather than the character —
    // and why the retired SymbolPalette's √ glyph was never going to produce one.
    await openLearn(page)
    await fresh(page, '√9')
    expect(await answerLatex(page)).toContain('\\surd')
    expect(await answerLatex(page)).not.toContain('\\sqrt')
  })

  test('the caret reaches inside a radical', async ({ page }) => {
    await openLearn(page)
    await fresh(page, '\\sqrt{9}')
    await page.keyboard.press('ArrowLeft')
    await typeIntoField(page, '8')
    expect(await answerLatex(page)).toContain('98')
  })

  test('√(x²) keeps the radicand as one unit', async ({ page }) => {
    await openLearn(page)
    await fresh(page, '\\sqrt{(x^2)}')
    const latex = await answerLatex(page)
    expect(latex.startsWith('\\sqrt')).toBe(true)
    // The x^2 stayed inside the radical rather than floating outside it.
    expect(latex).toContain('x^2')
  })

  test('a mixed number is one structure', async ({ page }) => {
    await openLearn(page)
    await fresh(page, '4\\frac{1}{10}')
    // The numeric grader wants "4 1/10"; keeping the whole and the fraction adjacent
    // is what makes that recoverable.
    expect(await answerLatex(page)).toBe('4\\frac{1}{10}')
  })

  test('2x + 3', async ({ page }) => {
    await openLearn(page)
    await fresh(page, '2x+3')
    expect(await answerLatex(page)).toBe('2x+3')
  })
})

test.describe('editing', () => {
  test('backspace deletes a digit', async ({ page }) => {
    await openLearn(page)
    await fresh(page, '12')
    await page.keyboard.press('Backspace')
    expect(await answerLatex(page)).toBe('1')
  })

  test('Delete removes forward from the start', async ({ page }) => {
    await openLearn(page)
    await fresh(page, '12')
    await runMathCommand(page, 'moveToMathfieldStart')
    await page.keyboard.press('Delete')
    expect(await answerLatex(page)).toBe('2')
  })

  test('a selection is replaced wholesale', async ({ page }) => {
    await openLearn(page)
    await fresh(page, '2+3')
    await runMathCommand(page, 'selectAll')
    await typeIntoField(page, '9')
    expect(await answerLatex(page)).toBe('9')
  })

  test('arrow keys walk the expression', async ({ page }) => {
    await openLearn(page)
    await fresh(page, '12')
    await runMathCommand(page, 'moveToMathfieldStart')
    await page.keyboard.press('ArrowRight')
    await typeIntoField(page, '0')
    // Inserted after the 1, not at the end.
    expect(await answerLatex(page)).toBe('102')
  })

  test('backspace immediately after a typed fraction is a no-op', async ({ page }) => {
    // Recorded, not assumed: MathLive 0.111.0 leaves the expression untouched here,
    // because the caret sits after the fraction's trailing placeholder rather than
    // adjacent to a removable atom. A second backspace is what removes it.
    await openLearn(page)
    await fresh(page, '\\frac{1}{2}')
    await page.keyboard.press('Backspace')
    expect(await answerLatex(page)).toBe('\\frac{1}{2}')
    await page.keyboard.press('Backspace')
    expect(await answerLatex(page)).not.toBe('\\frac{1}{2}')
  })
})

test.describe('virtual keyboard', () => {
  test('stays closed on desktop, so it never covers the answer', async ({ page }) => {
    await openLearn(page)
    await focusAnswerField(page)
    // A physical keyboard is present; raising a four-row keypad on every focus is the
    // "giant calculator" the design exists to avoid.
    await expect(page.locator('.MLK__rows')).toHaveCount(0)
  })

  test('digit keys insert into the field', async ({ page }) => {
    await openLearn(page)
    await setAnswer(page, '')
    await focusAnswerField(page)
    await showVirtualKeyboard(page)
    await pressKeycap(page, '7')
    await pressKeycap(page, '5')
    expect(await answerLatex(page)).toBe('75')
  })

  test('arithmetic offers no calculus keys', async ({ page }) => {
    await openLearn(page)
    await focusAnswerField(page)
    await showVirtualKeyboard(page)
    const panel = keyboardRows(page)
    for (const forbidden of ['∫', '∂', 'Σ', 'lim', 'sin', 'cos', 'tan']) {
      await expect(panel, `arithmetic keypad shows ${forbidden}`).not.toContainText(forbidden)
    }
    // …and does offer what an addition problem needs.
    await expect(panel).toContainText('÷')
    await expect(panel).toContainText('√')
  })
})

test.describe('submission is unchanged', () => {
  test('an integer answer reaches the server as a plain integer', async ({ page }) => {
    const submitted = captureSubmittedAnswer(page, '/api/study/answer')
    await page.route('**/api/study/answer', route => {
      // Stands in for SymPy: the learner submits the expression `2+2`, and the real
      // grader decides it equals the expected `4`. Accepting either form keeps the
      // assertion about the wire format rather than about the mock's arithmetic.
      const body = JSON.parse(route.request().postData() ?? '{}') as { answer?: string }
      const right = body.answer === '4' || body.answer === '2+2'
      return route.fulfill({ json: { correct: right, feedback: 'Correct!', xp: right ? 1 : 0, halted: false } })
    })
    await openLearn(page)
    await fresh(page, '2+2')
    await page.getByRole('button', { name: 'Check', exact: true }).first().click()
    await expect(page.getByText('+1 XP').first()).toBeVisible({ timeout: 20_000 })
    expect(submitted()).toBe('2+2')
  })

  test('a fraction reaches the server as a plain fraction, not as LaTeX', async ({ page }) => {
    const submitted = captureSubmittedAnswer(page, '/api/study/answer')
    await page.route('**/api/study/answer', route => route.fulfill({ json: { correct: true, feedback: 'Correct!', xp: 1, halted: false } }))
    await openLearn(page)
    await fresh(page, '1/2')
    await page.getByRole('button', { name: 'Check', exact: true }).first().click()
    await expect(page.getByText('+1 XP').first()).toBeVisible({ timeout: 20_000 })
    // The whole reason components/math/plainAnswer.ts exists. Unnormalized this is
    // `(1)/(2)`, which numericGrader cannot parse — a false miss on a correct answer.
    expect(submitted()).toBe('1/2')
  })

  test('a mixed number reaches the server in the form the grader parses', async ({ page }) => {
    const submitted = captureSubmittedAnswer(page, '/api/study/answer')
    await page.route('**/api/study/answer', route => route.fulfill({ json: { correct: true, feedback: 'Correct!', xp: 1, halted: false } }))
    await openLearn(page)
    await fresh(page, '4\\frac{1}{10}')
    await page.getByRole('button', { name: 'Check', exact: true }).first().click()
    await expect(page.getByText('+1 XP').first()).toBeVisible({ timeout: 20_000 })
    expect(submitted()).toBe('4 1/10')
  })

  test('the learning loop is unbroken: question → answer → verdict → explanation → next', async ({ page }) => {
    await page.route('**/api/study/answer', route => {
      const body = JSON.parse(route.request().postData() ?? '{}') as { answer?: string }
      const right = body.answer === '4' || body.answer === '2+2'
      return route.fulfill({ json: { correct: right, feedback: 'Nope', xp: right ? 1 : 0, explanation: 'Two and two.', halted: false } })
    })
    await openLearn(page)
    await expect(page.getByText('2 + 2 = ?').first()).toBeVisible()

    await fresh(page, '2+2')
    await page.getByRole('button', { name: 'Check', exact: true }).first().click()

    await expect(page.getByText('Two and two.').first()).toBeVisible({ timeout: 20_000 })
    await expect(page.getByText('+1 XP').first()).toBeVisible({ timeout: 20_000 })
    await expect(page.getByText('3 + 5 = ?').first()).toBeVisible({ timeout: 20_000 })
    await expect(answerField(page)).toBeVisible()
  })

  test('Enter submits from the editor', async ({ page }) => {
    let posted = 0
    await page.route('**/api/study/answer', route => {
      posted++
      return route.fulfill({ json: { correct: true, feedback: 'Correct!', xp: 1, halted: false } })
    })
    await openLearn(page)
    await fresh(page, '2+2')
    await page.keyboard.press('Enter')
    await expect(page.getByText('+1 XP').first()).toBeVisible({ timeout: 20_000 })
    expect(posted).toBe(1)
  })
})

test.describe('visual integration', () => {
  test('the field wears Mathua, not the MathLive demo', async ({ page }) => {
    await openLearn(page)
    const box = await answerField(page).boundingBox()
    // The established Mathua control height, matching the plain input it replaces.
    expect(box?.height).toBeGreaterThanOrEqual(48)

    const styles = await answerField(page).evaluate(el => {
      const cs = getComputedStyle(el)
      return { radius: cs.borderRadius, borderStyle: cs.borderStyle }
    })
    // Square corners and a hairline border: Mathua's house style. MathLive's default
    // is a rounded box, which is what makes a dropped-in widget look dropped in.
    expect(styles.radius).toBe('0px')
    expect(styles.borderStyle).toBe('solid')
  })

  test('the keyboard inherits Mathua tokens when it does open', async ({ page }) => {
    await openLearn(page)
    await focusAnswerField(page)
    await showVirtualKeyboard(page)
    const tokens = await page.evaluate(() => {
      const cs = getComputedStyle(document.body)
      return {
        keycap: cs.getPropertyValue('--keycap-background').trim(),
        border: cs.getPropertyValue('--keycap-border').trim(),
      }
    })
    // An empty value here would mean the keypad is falling back to MathLive's palette
    // rather than Mathua's.
    expect(tokens.keycap).not.toBe('')
    expect(tokens.border).not.toBe('')
  })
})