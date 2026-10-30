import { test, expect } from './fixtures'

// The reported failure: "View prerequisite path" on a locked curriculum row did
// nothing useful, the panel rendered below a 128-concept list, outside the
// viewport. The dialog must open immediately, stay visible for any row, link
// every prerequisite to its reference page with the target preserved, and
// return to the original target.
test('domains prerequisite dialog opens for a bottom row and returns to target', async ({ page }) => {
  await page.addInitScript(() => {
    localStorage.setItem('mathua_token', 'fake-token')
    localStorage.setItem('mathua_user', JSON.stringify({ student_id: 's1', name: 'Tester', username: 'tester', concepts_mastered: 0, current_streak: 0, level: 'Novice', diagnostic_completed: true }))
  })

  const topics = Array.from({ length: 30 }, (_, i) => ({
    conceptId: `calc.c${String(i).padStart(3, '0')}`,
    title: `Calculus concept ${i}`,
    state: 'locked',
    depth: 1,
  }))
  topics.push({ conceptId: 'calc.limit', title: 'Concept of a limit', state: 'locked', depth: 1 })

  await page.route('**/api/curriculum/domains', (route) => route.fulfill({
    json: {
      domains: [
        { id: 'calc1', conceptCount: 31, masteredCount: 0, unlockedCount: 0, inProgressCount: 0 },
      ],
    },
  }))
  await page.route('**/api/curriculum/domains/calc1', (route) => route.fulfill({
    json: { domain: { id: 'calc1' }, topics },
  }))
  await page.route('**/api/lessons/**', (route) => {
    const url = route.request().url()
    if (url.includes('/readiness')) {
      if (url.includes('calc.limit')) {
        return route.fulfill({ json: {
          concept_id: 'calc.limit', ready: false, weak: [], missing: [],
          eligible: false, state: 'locked',
          prerequisites: [
            { conceptId: 'calc.func', title: 'What is a function: domain and range', state: 'locked', met: false },
            { conceptId: 'calc.geo', title: 'Geometric sequences', state: 'locked', met: false },
          ],
        } })
      }
      return route.fulfill({ json: {
        concept_id: 'x', ready: true, weak: [], missing: [], eligible: true, state: 'unlocked', prerequisites: [],
      } })
    }
    if (url.includes('/kp')) return route.fulfill({ json: { concept_id: 'x', kps: [], diagram: '' } })
    return route.fulfill({ json: {} })
  })
  await page.route('**/api/next**', (route) => route.fulfill({ json: { primary: null, alternatives: [], generatedAt: '' } }))
  await page.route('**/api/progress/**', (route) => route.fulfill({ json: {} }))

  await page.goto('/domains')
  await page.getByText('Calculus', { exact: false }).first().click()
  await expect(page.getByText('Concept of a limit')).toBeVisible({ timeout: 30_000 })

  // The target is the last row of a long list: scroll it into view like a learner would.
  const trigger = page.getByRole('button', { name: /View prerequisite path/ }).last()
  await trigger.scrollIntoViewIfNeeded()
  await trigger.click()

  const dialog = page.getByRole('dialog')
  await expect(dialog).toBeVisible({ timeout: 30_000 })
  await expect(dialog).toContainText('Concept of a limit')
  await expect(dialog).toContainText('What is a function: domain and range')
  await expect(dialog).toContainText('Geometric sequences')

  const studyLinks = dialog.getByRole('link', { name: /Study →/ })
  await expect(studyLinks).toHaveCount(2)
  const href = await studyLinks.first().getAttribute('href')
  expect(href).toContain('/study?concept=calc.func')
  expect(href).toContain(`return=${encodeURIComponent('/study?concept=calc.limit')}`)

  await studyLinks.first().click()
  await expect(page).toHaveURL(/\/study\?concept=calc\.func/, { timeout: 30_000 })
  await expect(page.getByRole('link', { name: /Back to/ }).first()).toBeVisible({ timeout: 30_000 })
})
