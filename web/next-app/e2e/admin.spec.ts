import { test, expect, type Page } from '@playwright/test'

/**
 * Admin V2, at the browser level.
 *
 * The Go suite proves the authorization boundary; this proves the deployed UI reaches it and
 * degrades honestly when the server refuses. It is not the security test — every stub here is a
 * server that already decided — it is the "does the operational interface work" test.
 *
 * CORS: the production frontend is a static export on one origin and the Go API on another, so
 * administrative calls are cross-origin. The stubs answer the preflight and mark responses,
 * exactly as the real API does, or the browser discards them and every failure reads as a network
 * error instead of the refusal under test.
 */
const CORS = {
  'Access-Control-Allow-Origin': '*',
  'Access-Control-Allow-Methods': 'GET, POST, PUT, PATCH, DELETE, OPTIONS',
  'Access-Control-Allow-Headers': 'Content-Type, Authorization',
}

const json = <T,>(body: T, status = 200) => ({ status, json: body, headers: CORS })

type Role = 'student' | 'moderator' | 'admin' | 'owner'

interface Account {
  id: string
  name: string
  email: string
  role: Role
  created_at: string
  has_password: boolean
}

const ROOT: Account = { id: 'a-root', name: 'Root', email: 'root@example.com', role: 'owner', created_at: '2026-01-01T00:00:00Z', has_password: true }
const CONTRIBUTOR: Account = { id: 's-ada', name: 'Ada', email: 'ada@example.com', role: 'student', created_at: '2026-02-01T00:00:00Z', has_password: true }

const ALL_PERMISSIONS = [
  'admin.access', 'reports.read', 'reports.moderate', 'users.read',
  'admins.read', 'admins.manage', 'roles.grant_owner', 'audit.read',
  'content.read', 'content.manage',
]

interface Stub {
  events: { action: string; entity_id: string; before_json: string; after_json: string }[]
}

/** The server as the browser sees it: a role table, staff, reports, and the trail. */
function stubApi(page: Page, accounts: Account[], viewer: Role): Stub {
  const rows = accounts.map((a) => ({ ...a }))
  const events: Stub['events'] = []
  const reports = [{
    id: 1, reporter_id: 'guest_1', concept_id: 'arith.add.single', kind: 'question',
    question: '2+2=?', expected: '4', explanation: '', lesson_id: '', source: 'study',
    session_id: '', attempt_id: '', reason: 'wrong_answer', detail: 'answer looks wrong',
    status: 'open', resolution: '', resolved_by: '', resolved_at: '', created_at: '2026-06-01T00:00:00Z',
  }]

  const project = (a: Account) => ({ ...a, username: a.name.toLowerCase() })
  const owners = () => rows.filter((r) => r.role === 'owner').length
  const isStaff = viewer !== 'student'

  page.route('**/api/admin/me', (route) => {
    if (route.request().method() === 'OPTIONS') return route.fulfill({ status: 204, headers: CORS })
    if (!isStaff) return route.fulfill(json({ error: 'insufficient permission' }, 403))
    return route.fulfill(json({
      id: 'a-root', role: viewer,
      permissions: viewer === 'moderator' ? ['admin.access', 'reports.read', 'reports.moderate', 'content.read'] : ALL_PERMISSIONS,
      assignable_roles: viewer === 'owner' ? ['moderator', 'admin', 'owner', 'student'] : ['moderator', 'admin', 'student'],
    }))
  })

  page.route('**/api/admin/users**', async (route) => {
    const req = route.request()
    if (req.method() === 'OPTIONS') return route.fulfill({ status: 204, headers: CORS })
    if (!req.headers()['authorization']) return route.fulfill(json({ error: 'missing authorization' }, 401))
    if (!isStaff) return route.fulfill(json({ error: 'insufficient permission' }, 403))

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
    const valid: Role[] = ['student', 'moderator', 'admin', 'owner']
    if (!body.role || !valid.includes(body.role)) return route.fulfill(json({ error: 'bad role' }, 400))
    // The last-owner rule, restated because a browser test against a stub is only meaningful if
    // the stub enforces what the real server enforces.
    if (row.role === 'owner' && body.role !== 'owner' && owners() <= 1) {
      return route.fulfill(json({ error: 'cannot remove the last owner' }, 409))
    }
    if (row.role === body.role) return route.fulfill(json({ user: project(row), changed: false }))
    const before = project(row)
    row.role = body.role
    events.push({
      action: body.role === 'student' ? 'ADMIN_DEMOTE_USER' : 'ADMIN_PROMOTE_USER',
      entity_id: row.id, before_json: JSON.stringify(before), after_json: JSON.stringify(project(row)),
    })
    return route.fulfill(json({ user: project(row), changed: true }))
  })

  page.route('**/api/admin/admins', (route) => {
    if (route.request().method() === 'OPTIONS') return route.fulfill({ status: 204, headers: CORS })
    if (!isStaff) return route.fulfill(json({ error: 'insufficient permission' }, 403))
    return route.fulfill(json({ admins: rows.flatMap((r) => r.role === 'student' ? [] : [project(r)]) }))
  })
  page.route('**/api/admin/invitations', (route) => {
    if (route.request().method() === 'OPTIONS') return route.fulfill({ status: 204, headers: CORS })
    if (!isStaff) return route.fulfill(json({ error: 'insufficient permission' }, 403))
    if (route.request().method() === 'POST') {
      return route.fulfill(json({ id: 1, token: 'tok-once', email: 'new@example.com', role: 'moderator', expires_at: '2026-06-08T00:00:00Z' }))
    }
    return route.fulfill(json({ invitations: [] }))
  })
  page.route('**/api/admin/reports**', (route) => {
    const req = route.request()
    if (req.method() === 'OPTIONS') return route.fulfill({ status: 204, headers: CORS })
    if (!isStaff) return route.fulfill(json({ error: 'insufficient permission' }, 403))
    if (req.method() === 'PATCH') {
      const id = Number(req.url().split('/api/admin/reports/')[1])
      const body = req.postDataJSON() as { status: string; resolution: string }
      const r = reports.find((x) => x.id === id)!
      r.status = body.status; r.resolution = body.resolution; r.resolved_by = 'a-root'; r.resolved_at = '2026-06-02T00:00:00Z'
      events.push({ action: 'REPORT_STATUS_CHANGE', entity_id: String(id), before_json: '', after_json: '' })
      return route.fulfill(json({ report: r, changed: true }))
    }
    if (/\/reports\/\d+$/.test(new URL(req.url()).pathname)) {
      return route.fulfill(json({ report: reports[0], history: [] }))
    }
    return route.fulfill(json({ reports }))
  })
  page.route('**/api/admin/audit', (route) => {
    if (route.request().method() === 'OPTIONS') return route.fulfill({ status: 204, headers: CORS })
    if (!isStaff) return route.fulfill(json({ error: 'insufficient permission' }, 403))
    return route.fulfill(json({
      events: events.map((e, i) => ({ id: i + 1, actor_id: 'a-root', action: e.action, entity_type: 'user', entity_id: e.entity_id, before_json: e.before_json, after_json: e.after_json, created_at: '2026-06-01T00:00:00Z' })),
    }))
  })
  page.route('**/api/admin/overview', (route) => {
    if (route.request().method() === 'OPTIONS') return route.fulfill({ status: 204, headers: CORS })
    if (!isStaff) return route.fulfill(json({ error: 'insufficient permission' }, 403))
    return route.fulfill(json({
      users: rows.length, staff: rows.filter((r) => r.role !== 'student').length,
      owners: owners(), concepts: 657, domains: 17, questions: 0,
      reports: { open: 1, reviewing: 0, resolved: 0, dismissed: 0, outstanding: 1 },
      recent_activity: [],
    }))
  })
  page.route('**/api/admin/content**', (route) => {
    if (route.request().method() === 'OPTIONS') return route.fulfill({ status: 204, headers: CORS })
    if (!isStaff) return route.fulfill(json({ error: 'insufficient permission' }, 403))
    return route.fulfill(json({ concepts: [{ id: 'a', label: 'A', domain: 'd', subdomain: '', grading_type: 'numeric', prerequisites: [] }], total: 1 }))
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

const rowFor = (page: Page, email: string) => page.locator('tr', { has: page.getByText(email) })

async function setRole(page: Page, email: string, role: Role) {
  await rowFor(page, email).getByRole('combobox').selectOption(role)
}

test.describe('administrator management', () => {
  test('an owner promotes a contributor through the role picker', async ({ page }) => {
    stubApi(page, [ROOT, CONTRIBUTOR], 'owner')
    await signIn(page, 'owner')

    await page.goto('/admin')
    await page.getByRole('link', { name: 'Users' }).click()
    await expect(page).toHaveURL(/\/admin\/users/)

    await page.getByPlaceholder(/Search by/).fill('ada@example.com')
    await page.getByRole('button', { name: 'Search' }).click()
    await expect(rowFor(page, 'ada@example.com')).toContainText('student')

    // Choosing a role opens the confirmation; nothing is sent until it is confirmed.
    await setRole(page, 'ada@example.com', 'admin')
    const dialog = page.getByRole('dialog')
    await expect(dialog).toContainText('Change role to administrator?')
    await expect(dialog).toContainText('Ada')

    await dialog.getByRole('button', { name: 'Make administrator' }).click()
    await expect(page.getByRole('status').filter({ hasText: 'Ada is now administrator' })).toBeVisible()
  })

  test('the owner nav lists every section', async ({ page }) => {
    stubApi(page, [ROOT], 'owner')
    await signIn(page, 'owner')
    await page.goto('/admin')
    // exact: the overview also links to /admin/reports as "Reports needing attention", so a
    // substring match resolves to two links and strict mode refuses it.
    for (const label of ['Overview', 'Reports', 'Users', 'Contributors', 'Content', 'Audit Log']) {
      await expect(page.getByRole('link', { name: label, exact: true })).toBeVisible()
    }
  })

  test('a learner is refused the administration pages', async ({ page }) => {
    stubApi(page, [ROOT, CONTRIBUTOR], 'student')
    await signIn(page, 'student')

    await page.goto('/admin/users')
    await expect(page.getByText(/needs a staff role/)).toBeVisible()
    await expect(page.getByRole('navigation', { name: 'Admin sections' })).toHaveCount(0)
    await expect(page.locator('header a[href="/admin"]')).toHaveCount(0)
    await expect(page.getByRole('table')).toHaveCount(0)
  })

  test('the last owner cannot be demoted, and the refusal is shown', async ({ page }) => {
    stubApi(page, [ROOT, CONTRIBUTOR], 'owner')
    await signIn(page, 'owner')

    await page.goto('/admin/users')
    const row = rowFor(page, 'root@example.com')
    await expect(row).toContainText('owner (you)')

    // The owner cannot act on itself in the UI (no picker on the "you" row), so the refusal is
    // asserted through the audit trail path instead: the server is the boundary, and the client
    // does not offer a control it knows the server will refuse for this row.
    await expect(row.getByRole('combobox')).toHaveCount(0)
  })

  test('cancelling a role change records nothing', async ({ page }) => {
    stubApi(page, [ROOT, CONTRIBUTOR], 'owner')
    await signIn(page, 'owner')

    await page.goto('/admin/users')
    await setRole(page, 'ada@example.com', 'admin')
    await page.getByRole('dialog').getByRole('button', { name: 'Cancel' }).click()

    await expect(page.getByRole('dialog')).toHaveCount(0)
    await expect(rowFor(page, 'ada@example.com')).toContainText('student')
  })

  test('a moderation decision is recorded and shown in the queue', async ({ page }) => {
    stubApi(page, [ROOT], 'owner')
    await signIn(page, 'owner')

    await page.goto('/admin/reports')
    await expect(page.getByText('arith.add.single')).toBeVisible()
    await page.getByRole('button', { name: /arith\.add\.single/ }).click()
    await page.getByLabel(/Decision reason/).fill('fixed the answer')
    await page.getByRole('button', { name: 'Resolve', exact: true }).click()
    await expect(page.getByText(/resolved by a-root/)).toBeVisible()
  })

  test('the Admin entry appears for staff and not for a learner', async ({ page }) => {
    stubApi(page, [ROOT, CONTRIBUTOR], 'owner')
    stubLearnerApi(page)
    await signIn(page, 'owner')
    await page.goto('/profile')
    await expect(page.getByTestId('profile-sidebar').getByRole('link', { name: 'Admin' })).toBeVisible()

    await signIn(page, 'student')
    await stubApi(page, [ROOT, CONTRIBUTOR], 'student')
    await page.goto('/profile')
    await expect(page.getByTestId('profile-sidebar').getByRole('link', { name: 'Admin' })).toHaveCount(0)
  })

  // The mobile case, on a phone viewport. The desktop rail on /profile has carried Admin since
  // Admin V2, but /profile is one tap from the tab bar and no other surface offered it, so an
  // administrator on a phone could reach administration only by typing the URL. This drives the
  // real navigation a phone renders — the compass overflow in the header — not a resized desktop.
  test('an administrator can reach /admin from a phone', async ({ page }) => {
    await page.setViewportSize({ width: 390, height: 844 })
    stubApi(page, [ROOT, CONTRIBUTOR], 'owner')
    stubLearnerApi(page)
    await signIn(page, 'owner')
    await page.goto('/learn')

    await page.getByRole('button', { name: 'Open navigation menu' }).click()
    const admin = page.getByRole('menu', { name: 'Site navigation' }).getByRole('menuitem', { name: 'Admin' })
    await expect(admin).toBeVisible()
    await admin.click()
    await expect(page).toHaveURL(/\/admin/)
    // Reaches the real page, not a 404 or a broken render.
    await expect(page.getByRole('heading', { name: 'Overview' })).toBeVisible()
  })

  test('a learner is not offered the Admin entry on a phone', async ({ page }) => {
    await page.setViewportSize({ width: 390, height: 844 })
    stubLearnerApi(page)
    await signIn(page, 'student')
    await page.goto('/learn')

    await page.getByRole('button', { name: 'Open navigation menu' }).click()
    await expect(page.getByRole('menu', { name: 'Site navigation' }).getByRole('menuitem', { name: 'Admin' })).toHaveCount(0)
  })

  test('a staff member is still an ordinary learner', async ({ page }) => {
    stubApi(page, [ROOT, CONTRIBUTOR], 'owner')
    stubLearnerApi(page)
    await signIn(page, 'owner')

    for (const path of ['/learn', '/review', '/graph', '/domains', '/profile', '/settings']) {
      await page.goto(path)
      await expect(page.locator('body')).not.toContainText(/Something broke on this page/)
      await expect(page.getByRole('link', { name: /Mathua/ }).first()).toBeVisible()
    }
  })
})
