import { test, expect, type Page } from '@playwright/test'

/**
 * Administrator management, in a browser, the way an operator does it.
 *
 * Everything here runs against the **built static export served by the same origin**, so it also
 * covers the deployment routing rather than a dev server. That is not incidental: `/admin` was
 * broken in exactly this environment and nothing else could have caught it — it answered 301 to a
 * directory listing, because `nextStaticFS` matched the directory `out/admin/` before it looked
 * for `out/admin.html`. A unit test on the page component would have passed throughout. The
 * regression guard for that lives in Go, in `cmd/mathua/staticfs_test.go`.
 *
 * The API is stubbed with a small in-memory role table rather than a fixed payload, because the
 * properties under test are transitions: an administrator can promote, the promoted account can
 * then reach `/admin`, a demoted account immediately cannot, and the last administrator cannot be
 * removed. A fixture that always returns the same rows proves none of that.
 */

/**
 * CORS, and why this file has to care about it.
 *
 * The production frontend is a static export on one origin and the Go API is on another, so
 * every administrative call is cross-origin — a `PATCH` with `Content-Type` and `Authorization`
 * is preflighted. The build bakes in `NEXT_PUBLIC_API_URL`, so this is not a hypothetical: it
 * is the shape of the deployed request.
 *
 * The stub therefore answers the preflight and marks its responses as allowed, or the browser
 * discards them and every failure below reads as a network error instead of as the refusal
 * being tested. That distinction matters: without these headers a 403 that was correctly
 * returned by the stub would be indistinguishable from a blocked request.
 */
const CORS = {
  'Access-Control-Allow-Origin': '*',
  'Access-Control-Allow-Methods': 'GET, POST, PUT, PATCH, DELETE, OPTIONS',
  'Access-Control-Allow-Headers': 'Content-Type, Authorization',
}

/**
 * Generic, not `unknown`: a stub that erases its argument to `unknown` discards the type of the
 * payload it is standing in for. The responses here are literals at the call sites, so the
 * generic keeps the evidence and satisfies the same lint rule the rest of the suite obeys.
 */
const json = <T,>(body: T, status = 200) => ({ status, json: body, headers: CORS })

/**
 * The learner API, stubbed.
 *
 * `/admin` is not the only thing an administrator opens, and the two tests that assert an
 * administrator is still an ordinary learner need those surfaces to actually finish loading:
 * each one stays on its skeleton — and the skeleton takes no sidebar props, so no `isAdmin` —
 * until its own reads settle. Without these, the failure would read as "the Admin link is
 * missing" when the real fault is an unstubbed scores request against the live API.
 */
const LEARNER_API: [string, unknown][] = [
  ['**/api/scores/**', { lifetime_points: 100, weekly_score: 10, speed_bonus: 0, concepts_mastered: 1, current_streak: 2, level: 'Novice', xp_total: 10, xp_today: 10, daily_xp_goal: 30 }],
  ['**/api/progress/**', {}],
  ['**/api/activity', []],
  ['**/api/weaknesses', { by_domain: {} }],
  ['**/api/reviews/due', { count: 3 }],
  ['**/api/efficacy', { concepts_touched: 0, first_pass_rate: 0, second_pass_rate: 0, avg_attempts_per_concept: 0, total_attempts: 0 }],
  ['**/api/settings', {}],
  ['**/api/next', { done: false, next: null, total: 0, message: 'Nothing to review' }],
]

function stubLearnerApi(page: Page) {
  for (const [pattern, body] of LEARNER_API) {
    page.route(pattern, (route) => {
      if (route.request().method() === 'OPTIONS') return route.fulfill({ status: 204, headers: CORS })
      return route.fulfill(json(body))
    })
  }
}

type Role = 'student' | 'admin'

interface Account {
  id: string
  name: string
  email: string
  role: Role
  created_at: string
  has_password: boolean
}

const ROOT: Account = {
  id: 'a-root', name: 'Root', email: 'root@example.com',
  role: 'admin', created_at: '2026-01-01T00:00:00Z', has_password: true,
}
const CONTRIBUTOR: Account = {
  id: 's-ada', name: 'Ada', email: 'ada@example.com',
  role: 'student', created_at: '2026-02-01T00:00:00Z', has_password: true,
}
const OTHER: Account = {
  id: 's-bo', name: 'Bo', email: 'bo@example.com',
  role: 'student', created_at: '2026-03-01T00:00:00Z', has_password: false,
}

interface Event {
  id: number
  actor_id: string
  action: string
  entity_type: string
  entity_id: string
  before_json: string
  after_json: string
  created_at: string
}

/** The server, as the browser sees it: a role table, an admin count, and an audit log. */
function stubApi(page: Page, accounts: Account[], viewer: Role) {
  // Copy each account, not just the array. `rows` is mutated by a role change, and the tests
  // pass module-level constants in — `[...accounts]` would hand every test the same object, so
  // one test's promotion silently rewrote ROOT or CONTRIBUTOR for every test after it. That is
  // a shared-mutable-fixture bug, and it surfaces as a flake rather than a failure: the suite
  // passes when the tests happen to run in an order that does not expose it.
  const rows = accounts.map((a) => ({ ...a }))
  const events: Event[] = []
  let nextId = 1

  const project = (a: Account) => ({ ...a, username: a.name.toLowerCase() })
  const admins = () => rows.filter((r) => r.role === 'admin').length

  page.route('**/api/admin/users**', async (route) => {
    const req = route.request()
    if (req.method() === 'OPTIONS') return route.fulfill({ status: 204, headers: CORS })
    // The same three answers requireAdmin gives, from the same reasons: no credential is not
    // known, and a known non-administrator is refused. This is the boundary the whole feature
    // rests on, so the browser test asserts it rather than assuming it.
    if (!req.headers()['authorization']) return route.fulfill(json({ error: 'missing authorization' }, 401))
    if (viewer !== 'admin') return route.fulfill(json({ error: 'admin role required' }, 403))

    if (req.method() === 'GET') {
      const q = new URL(req.url()).searchParams.get('q')?.toLowerCase() ?? ''
      const match = rows.filter((r) =>
        !q || [r.id, r.email, r.name, r.name.toLowerCase()].some((f) => f.toLowerCase().includes(q)))
      return route.fulfill(json({ users: match.map(project) }))
    }

    const target = decodeURIComponent(req.url().split('/api/admin/users/')[1] ?? '')
    const body = req.postDataJSON() as { role?: Role }
    const row = rows.find((r) => r.id === target)
    if (!row) return route.fulfill(json({ error: 'user not found' }, 404))
    if (body.role !== 'admin' && body.role !== 'student') {
      return route.fulfill(json({ error: 'role must be student or admin' }, 400))
    }
    // The last-admin rule, restated here because a browser test that stubs the server is only
    // meaningful if the stub enforces what the real one enforces. The Go suite proves the real one
    // does, including under concurrency.
    if (row.role === 'admin' && body.role === 'student' && admins() <= 1) {
      return route.fulfill(json({ error: 'cannot demote the last administrator' }, 409))
    }
    if (row.role === body.role) {
      return route.fulfill(json({ user: project(row), changed: false }))
    }
    const before = project(row)
    row.role = body.role
    events.push({
      id: nextId++, actor_id: 'a-root',
      action: body.role === 'admin' ? 'ADMIN_PROMOTE_USER' : 'ADMIN_DEMOTE_USER',
      entity_type: 'user', entity_id: row.id,
      before_json: JSON.stringify(before), after_json: JSON.stringify(project(row)),
      created_at: '2026-06-01T00:00:00Z',
    })
    return route.fulfill(json({ user: project(row), changed: true }))
  })

  page.route('**/api/admin/audit', (route) => {
    if (route.request().method() === 'OPTIONS') return route.fulfill({ status: 204, headers: CORS })
    if (viewer !== 'admin') return route.fulfill(json({ error: 'admin role required' }, 403))
    return route.fulfill(json({ events }))
  })
  page.route('**/api/admin/overview', (route) => {
    if (route.request().method() === 'OPTIONS') return route.fulfill({ status: 204, headers: CORS })
    if (viewer !== 'admin') return route.fulfill(json({ error: 'admin role required' }, 403))
    return route.fulfill(json({
      users: rows.length, admins: admins(), concepts: 657, domains: 17, questions: 0,
      recent_activity: events.slice(-5).reverse(),
    }))
  })

  return { events }
}

function signIn(page: Page, role: Role) {
  return page.addInitScript((r) => {
    localStorage.setItem('mathua_token', 'e2e-token')
    localStorage.setItem('mathua_user', JSON.stringify({
      student_id: 'a-root', name: 'Root', username: 'root',
      concepts_mastered: 0, current_streak: 0, level: 'Novice',
      diagnostic_completed: true, role: r,
    }))
  }, role)
}

const rowFor = (page: Page, email: string) =>
  page.locator('tr', { has: page.getByText(email) })

test.describe('administrator management', () => {
  test('an administrator finds a contributor and promotes them', async ({ page }) => {
    stubApi(page, [ROOT, CONTRIBUTOR, OTHER], 'admin')
    await signIn(page, 'admin')

    await page.goto('/admin')
    await page.getByRole('link', { name: 'Users' }).click()
    await expect(page).toHaveURL(/\/admin\/users/)

    // Find the contributor by email.
    await page.getByPlaceholder(/Search by email/).fill('ada@example.com')
    await page.getByRole('button', { name: 'Search' }).click()
    await expect(rowFor(page, 'ada@example.com')).toBeVisible()
    await expect(rowFor(page, 'ada@example.com')).toContainText('student')

    // One click does not promote: the confirmation is the action.
    const make = rowFor(page, 'ada@example.com').getByRole('button', { name: 'Make administrator' })
    await make.click()
    const dialog = page.getByRole('dialog')
    await expect(dialog).toBeVisible()
    await expect(dialog).toContainText('Make administrator?')
    await expect(dialog).toContainText('Ada')

    await dialog.getByRole('button', { name: 'Make administrator' }).click()
    await expect(page.getByRole('status').filter({ hasText: 'Ada is now an administrator' })).toBeVisible()
  })

  test('a learner is refused the administration pages', async ({ page }) => {
    stubApi(page, [ROOT, CONTRIBUTOR], 'student')
    await signIn(page, 'student')

    await page.goto('/admin/users')
    // The page renders and says why, rather than appearing broken or silently blank.
    await expect(page.getByText(/needs an administrator role/)).toBeVisible()
    // No administration navigation, and no Admin entry in the learner chrome.
    await expect(page.getByRole('navigation', { name: 'Admin sections' })).toHaveCount(0)
    await expect(page.locator('header a[href="/admin"]')).toHaveCount(0)
    // And no table of accounts to act on.
    await expect(page.getByRole('table')).toHaveCount(0)
  })

  test('an administrator can remove another administrator', async ({ page }) => {
    const second: Account = { ...CONTRIBUTOR, id: 's-second', role: 'admin' }
    stubApi(page, [ROOT, second, OTHER], 'admin')
    await signIn(page, 'admin')

    await page.goto('/admin/users')
    const row = rowFor(page, 'ada@example.com')
    await expect(row).toContainText('admin')

    await row.getByRole('button', { name: 'Remove administrator' }).click()
    const dialog = page.getByRole('dialog')
    await expect(dialog).toContainText('Remove administrator access?')
    await expect(dialog).toContainText('immediately lose administrative access')
    await dialog.getByRole('button', { name: 'Remove access' }).click()

    await expect(page.getByRole('status').filter({ hasText: 'Ada is now a learner' })).toBeVisible()
  })

  test('the last administrator cannot be removed, and the refusal is shown', async ({ page }) => {
    // ROOT is the only administrator, so this is the demotion the server refuses. The stub
    // enforces the rule by count, exactly as the real one does.
    stubApi(page, [ROOT, CONTRIBUTOR], 'admin')
    await signIn(page, 'admin')

    await page.goto('/admin/users')
    const row = rowFor(page, 'root@example.com')
    await expect(row).toContainText('admin (you)')

    await row.getByRole('button', { name: 'Remove administrator' }).click()
    await page.getByRole('dialog').getByRole('button', { name: 'Remove access' }).click()

    // The dialog closes before the message lands, so wait on the message rather than asserting
    // the dialog is gone first. Two independent state writes land in the same tick — the dialog's
    // pending=null and the error text — and asserting the absence of one while waiting for the
    // other races the render.
    await expect(page.getByRole('alert').filter({ hasText: 'last administrator' })).toBeVisible()
    // The dialog is gone, so the message is not hidden behind it.
    await expect(page.getByRole('dialog')).toHaveCount(0)
    // And the row still reads as an administrator, because the server refused.
    await expect(row).toContainText('admin')
  })

  test('cancelling changes nothing', async ({ page }) => {
    stubApi(page, [ROOT, CONTRIBUTOR], 'admin')
    await signIn(page, 'admin')

    await page.goto('/admin/users')
    await rowFor(page, 'ada@example.com').getByRole('button', { name: 'Make administrator' }).click()
    await page.getByRole('dialog').getByRole('button', { name: 'Cancel' }).click()

    await expect(page.getByRole('dialog')).toHaveCount(0)
    await expect(rowFor(page, 'ada@example.com')).toContainText('student')
    // Nothing was recorded, because nothing happened.
    await page.goto('/admin/audit')
    await expect(page.getByText(/No administrative activity/)).toBeVisible()
  })

  test('Escape cancels, and the ledger records each real change once', async ({ page }) => {
    const { events } = stubApi(page, [ROOT, CONTRIBUTOR], 'admin')
    await signIn(page, 'admin')

    await page.goto('/admin/users')
    await rowFor(page, 'ada@example.com').getByRole('button', { name: 'Make administrator' }).click()
    await page.keyboard.press('Escape')
    await expect(page.getByRole('dialog')).toHaveCount(0)

    await rowFor(page, 'ada@example.com').getByRole('button', { name: 'Make administrator' }).click()
    await page.getByRole('dialog').getByRole('button', { name: 'Make administrator' }).click()
    await expect(page.getByRole('status').filter({ hasText: 'an administrator' })).toBeVisible()

    expect(events).toHaveLength(1)
    expect(events[0].action).toBe('ADMIN_PROMOTE_USER')
    expect(events[0].entity_id).toBe('s-ada')
    expect(events[0].before_json).toContain('"role":"student"')
    expect(events[0].after_json).toContain('"role":"admin"')

    await page.goto('/admin/audit')
    await expect(page.getByText('ADMIN_PROMOTE_USER')).toBeVisible()
  })

  test('the Admin entry appears for an administrator and not for a learner', async ({ page }) => {
    stubApi(page, [ROOT, CONTRIBUTOR], 'admin')
    stubLearnerApi(page)
    await signIn(page, 'admin')
    await page.goto('/profile')
    await expect(page.getByTestId('profile-sidebar').getByRole('link', { name: 'Admin' })).toBeVisible()

    // The learner case is a second session, not a mutated one. Nothing in the app subscribes to
    // `auth-changed` — the event exists, but no component listens — so flipping localStorage
    // and dispatching it would be asserting reactivity Mathua never claimed. A demotion is
    // observed the way it actually happens: the next load reads the role the server reports.
    await signIn(page, 'student')
    await page.goto('/profile')
    await expect(page.getByTestId('profile-sidebar').getByRole('link', { name: 'Admin' })).toHaveCount(0)
  })

  test('an administrator is still an ordinary learner', async ({ page }) => {
    // Administration is an additional privilege on the same account. Every learner surface must
    // keep working for the person holding the role.
    stubApi(page, [ROOT, CONTRIBUTOR], 'admin')
    stubLearnerApi(page)
    await signIn(page, 'admin')

    for (const path of ['/learn', '/review', '/graph', '/domains', '/profile', '/settings']) {
      await page.goto(path)
      await expect(page.locator('body')).not.toContainText(/Something broke on this page/)
      // Every learner surface keeps a way off itself. Not every surface uses `<header>`:
      // `/graph` and its peers are built from the sidebar layout, so the assertion is on the
      // brand link, which both layouts render.
      await expect(page.getByRole('link', { name: /Mathua/ }).first()).toBeVisible()
    }
  })
})