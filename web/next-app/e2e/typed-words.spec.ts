import { test, expect } from '@playwright/test'
import type { Page } from '@playwright/test'
import { answerLatex, answerField, focusAnswerField, waitForAnswerField } from './helpers/answer'

const CONCEPT = 'arith.add.single'

const PRACTICE = {
  concept_id: CONCEPT,
  questions: [
    { question: '2 + 2 = ?', answer: '4', explanation: 'Two and two.', source: 'curated' },
    { question: '3 + 5 = ?', answer: '8', explanation: 'Three and five.', source: 'curated' },
  ],
}

const KP = { concept_id: CONCEPT, kps: [{ key: 'k1', title: 'Add', sections: ['a'] }] }

async function openLearn(page: Page) {
  await page.route('**/api/lessons/**', route => {
    const url = route.request().url()
    if (url.includes('/practice')) return route.fulfill({ json: PRACTICE })
    if (url.includes('/kp')) return route.fulfill({ json: KP })
    if (url.includes('/readiness')) {
      return route.fulfill({ json: { concept_id: CONCEPT, ready: true, weak: [], missing: [] } })
    }
    return route.fulfill({ json: {} })
  })
  await page.route('**/api/scores/**', route =>
    route.fulfill({
      json: {
        lifetime_points: 0, weekly_score: 0, speed_bonus: 0, concepts_mastered: 1,
        current_streak: 0, level: 'Novice', xp_total: 0, xp_today: 0,
        daily_xp_goal: 30, spaced_reps: {}, avg_learning_time: 1,
      },
    }),
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

// Typing a word must not become a symbol.
//
// MathLive's `inlineShortcuts` default map rewrites bare text as it is typed: `pi` becomes
// π, `theta` becomes θ, `alpha` becomes α, `sqrt` becomes √. That is convenient on a maths
// keypad and wrong on a physical keyboard, where the reported bug was "typing a point has it
// completing math symbols" — a learner who types the word `pi` gets a glyph they did not
// type, in a field whose only job is to record what they typed. The submitted string then
// differs from the keystrokes.
//
// This has to be a browser test. jsdom cannot run MathLive (it needs CSSOM and a real layout
// engine), so the repo's other maths tests assert the pure helpers and the recorded corpus —
// which is exactly the gap that let this ship: `inlineShortcuts` was never set, and nothing
// could see it.

test.describe('typed words are not rewritten into symbols', () => {
  test.beforeEach(async ({ page }) => {
    await openLearn(page)
  })

  const cases: Array<{ typed: string; forbidden: string[]; why: string }> = [
    { typed: 'pi', forbidden: ['\\pi'], why: 'pi is a word as often as it is a constant' },
    { typed: 'theta', forbidden: ['\\theta'], why: 'theta' },
    { typed: 'alpha', forbidden: ['\\alpha'], why: 'alpha' },
    { typed: 'sqrt', forbidden: ['\\sqrt'], why: 'sqrt' },
    { typed: 'infty', forbidden: ['\\infty'], why: 'infty' },
  ]

  for (const c of cases) {
    test(`"${c.typed}" stays literal — ${c.why}`, async ({ page }) => {
      await focusAnswerField(page)
      await page.keyboard.type(c.typed, { delay: 40 })
      const latex = await answerLatex(page)
      for (const bad of c.forbidden) {
        expect(latex, `typing "${c.typed}" produced ${latex}`).not.toContain(bad)
      }
      expect(latex.trim(), `typing "${c.typed}" produced ${latex}`).toContain(c.typed)
    })
  }

  // The word forms that appear in the corpus's multiple-choice answers. `yes` is the single
  // most common expected answer in the corpus and `no` is second, and the multiple-choice
  // grader is an exact case-insensitive comparison — so a rewritten `yes` is a false miss,
  // not a cosmetic difference.
  test('multiple-choice words survive verbatim', async ({ page }) => {
    for (const word of ['yes', 'no', 'odd', 'even', 'addition']) {
      await focusAnswerField(page)
      await answerField(page).fill('')
      await page.keyboard.type(word, { delay: 40 })
      const latex = await answerLatex(page)
      expect(latex.replace(/\s+/g, ''), `"${word}" came back as "${latex}"`).toContain(word)
    }
  })

  // A symbol the learner *did* type must still work. Disabling shortcuts must not cost the
  // ability to enter real mathematics.
  test('real math still enters', async ({ page }) => {
    await focusAnswerField(page)
    // \frac{1}{2} via the field's own LaTeX input path.
    await answerField(page).fill('\\frac{1}{2}')
    await page.waitForTimeout(150)
    const latex = await answerLatex(page)
    expect(latex).toContain('frac')
  })
})