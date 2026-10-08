import { test, expect } from '@playwright/test'

// Layout at a width nothing in jsdom can measure.
//
// Every other header test here is structural — it can check that a panel is not a descendant of
// the clipper, which is the *cause* of the paint-order bug, but not that two boxes land in the
// right place. These two measure geometry, because the wide-screen faults are geometry faults: the
// rule that stopped at 1100px and the menu that drifted to the window edge both rendered, both
// were in the DOM, and both were readable. Only a layout engine can see either.

const WIDE = { width: 1600, height: 900 }

test('the header rule reaches both edges of the window', async ({ page }) => {
  await page.setViewportSize(WIDE)
  await page.goto('/')

  const box = await page.evaluate(() => {
    // The wrapper carrying the rule: a direct child of the `min-h-0 overflow-hidden` clipper.
    const rule = document.querySelector('header > div > div') as HTMLElement
    const row = rule.firstElementChild as HTMLElement
    return {
      ruleLeft: rule.getBoundingClientRect().left,
      ruleRight: rule.getBoundingClientRect().right,
      rowLeft: row.getBoundingClientRect().left,
      rowRight: row.getBoundingClientRect().right,
      vw: window.innerWidth,
    }
  })

  expect(box.vw).toBe(WIDE.width)
  // The rule is the header's own edge. It used to sit on the capped row, so on any screen wider
  // than the cap it stopped short and the bar read as a floating strip.
  expect(box.ruleLeft).toBe(0)
  expect(box.ruleRight).toBe(box.vw)

  // The links stay inside the cap, because the page content below is capped the same way. A
  // full-width row would put the logo away from the corner it is supposed to anchor.
  expect(box.rowLeft).toBeGreaterThan(0)
  expect(box.rowRight).toBeLessThan(box.vw)
  expect(box.rowRight - box.rowLeft).toBeLessThanOrEqual(1100)
})

test('the profile menu opens beside its trigger, not at the window edge', async ({ page }) => {
  await page.setViewportSize(WIDE)
  await page.addInitScript(() => {
    localStorage.setItem('mathua_token', 't')
    localStorage.setItem('mathua_user', JSON.stringify({
      student_id: 's1', name: 'Nelson', username: 'wis', concepts_mastered: 1,
      current_streak: 1, level: 'Novice', diagnostic_completed: true, role: 'admin',
    }))
  })
  for (const [pattern, body] of [
    ['**/api/auth/me', { student_id: 's1', name: 'Nelson', has_password: true, role: 'admin' }],
    ['**/api/scores/**', { lifetime_points: 100, weekly_score: 10, speed_bonus: 0, concepts_mastered: 1, current_streak: 1, level: 'Novice', xp_total: 10, xp_today: 10, daily_xp_goal: 30 }],
    ['**/api/progress/**', {}],
    ['**/api/activity', []],
    ['**/api/weaknesses', { by_domain: {} }],
    ['**/api/reviews/due', { count: 0 }],
    ['**/api/settings', {}],
    ['**/api/next', { done: false, next: null, total: 0, message: '' }],
  ] as [string, unknown][]) {
    page.route(pattern, r => r.fulfill({ json: body }))
  }

  // /learn, not /profile: the profile page is built from the sidebar shell and renders no Header.
  await page.goto('/learn')
  const trigger = page.getByRole('button', { name: 'Open profile menu' })
  await trigger.click()
  const panel = page.getByRole('menu', { name: 'Profile' })
  await expect(panel).toBeVisible()

  const t = (await trigger.boundingBox())!
  const p = (await panel.boundingBox())!
  // The panel's trailing edge lines up with the trigger's, the way it does on a narrow screen.
  // It hung off the full-width <header> without a capped wrapper, so `right-6` put it at the
  // window edge — roughly 250px from the avatar at this width, and further on a wider one.
  expect(Math.abs((p.x + p.width) - (t.x + t.width))).toBeLessThan(40)
})
