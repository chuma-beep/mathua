import { test as base } from '@playwright/test'

// Every spec imports `test` and `expect` from here rather than from
// `@playwright/test`.
//
// The assistant dock is chrome: it is mounted once in `app/layout.tsx`, so every
// page asks `/api/config` on load to find out whether the assistant is on. That
// is one request the suite never stubbed, and it is not a request to the page's
// own server, a normal `next build` bakes `NEXT_PUBLIC_API_URL` from
// `.env.production`, which is the deployed API. So `/glossary`, a purely static
// reference, was reaching `https://mathua.fly.dev` from `localhost:4173`, and
// the browser logged the CORS refusal. `glossary.spec.ts` asserts that the
// console holds no errors, and it was right to fail: the page should not need
// the network, and now it did.
//
// Answering the request locally is the honest fix. The stub returns
// `ai_enabled: false`, which is the state a visitor with no assistant sees, the
// trigger still renders (that is deliberate, see ADR-055) and no further request
// is made, so a spec that asserts "this page renders without a network" keeps
// meaning what it says. A spec that wants the assistant on can register its own
// route for `/api/config` after this one; the last registration wins.
//
// The rest of this file is the net. Anything that is not the static export is a
// stub that missed: it is recorded, refused, and the test that caused it fails.
// Before the net existed the URL only mattered because the bundle pointed at
// production, so a missed stub was answered by the deployed API: the test passed
// on a fast reply, waited on a slow one, and two of those requests create state
// (`POST /api/auth/guest`). CI now builds the export against a documentation
// address that never answers, so nothing can reach production from a test. The
// net is what turns a silent 404 into a named failure, and it is the only thing
// that catches a request fired after the test stopped looking.
export const test = base.extend({
  // Registered on the context, not the page, so that a spec's own `page.route`
  // and `context.route` registrations are consulted first: page routes outrank
  // context routes, and within the context the last registration wins. The net
  // therefore sees only requests that no other handler wanted.
  context: async ({ context, baseURL }, use) => {
    const local = new URL(baseURL ?? 'http://localhost:4173')
    const escaped: string[] = []

    await context.route('**/*', (route) => {
      const request = route.request()
      const url = new URL(request.url())
      if (url.origin === local.origin) return route.fallback()
      escaped.push(`${request.method()} ${url.href}`)
      return route.abort('blockedbyclient')
    })

    // The footer credits strip draws three DiceBear shapes on every page. No
    // test asserts on the picture, and a suite that needs that host to answer is
    // a suite that fails when the host is slow, which is a third dependency to
    // explain in a red run. One blank square is enough for layout.
    await context.route('https://api.dicebear.com/**', (route) =>
      route.fulfill({
        contentType: 'image/svg+xml',
        body: '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1 1"></svg>',
      }),
    )

    // The anonymous shell. Every page asks these three before it draws anything, and a visitor with
    // no account is the state they answer: no session, default settings, a server that is up. A
    // spec that wants a signed-in learner registers its own route and wins, because it is later.
    await context.route('**/api/auth/me**', (route) =>
      route.fulfill({ status: 401, json: { error: 'unauthorized' } }),
    )
    await context.route('**/api/settings**', (route) => route.fulfill({ json: {} }))
    await context.route('**/api/health**', (route) => route.fulfill({ json: { status: 'ok' } }))
    // The rest of the anonymous shell: no linked providers, no plan, no study list and nothing
    // queued. Each is an empty answer a visitor with no account receives, so a spec that asserts
    // content registers its own route and this default never decides what the test sees.
    await context.route('**/api/auth/identities**', (route) => route.fulfill({ json: [] }))
    await context.route('**/api/destinations**', (route) => route.fulfill({ json: { destinations: [] } }))
    await context.route('**/api/plans/current**', (route) => route.fulfill({ json: { plan: null } }))
    await context.route('**/api/lessons**', (route) => route.fulfill({ json: { lessons: {} } }))
    await context.route('**/api/graph**', (route) => route.fulfill({ json: { nodes: [] } }))
    await context.route('**/api/next**', (route) =>
      route.fulfill({ json: { done: true, next: null, total: 0, message: 'Nothing to review' } }),
    )
    // A guest who opens a lesson by its link mints a guest token before the first question is
    // served, so every `/learn?concept=` page asks for one. The answer is the server's: the id
    // the client proposed, now bound to a token. A spec that asserts on the token registers its
    // own route and wins, because it is later.
    await context.route('**/api/auth/guest', (route) => {
      const { student_id } = route.request().postDataJSON() as { student_id: string }
      return route.fulfill({ json: { token: 'e2e-guest-token', student_id } })
    })
    // The report badge in the header asks for the learner's own moderation
    // reports on every signed-in page. An account with no reports is the
    // default the shell should assume; a spec that wants history registers its
    // own route and wins, because it is later.
    await context.route('**/api/reports/mine**', (route) =>
      route.fulfill({ json: { reports: [], unread: 0 } }),
    )
    // The sidebar shell asks for the due-review count on every signed-in
    // product page. Nothing queued is the default empty account; a spec that
    // wants a badge registers its own route later and wins.
    await context.route('**/api/reviews/due**', (route) =>
      route.fulfill({ json: { count: 0 } }),
    )
    // The sidebar shell and the notifications page ask for product
    // announcements. An account with no announcements is the default; a spec
    // that wants an unread item registers its own route later and wins.
    await context.route('**/api/announcements**', (route) =>
      route.fulfill({ json: { announcements: [], unread: 0 } }),
    )
    // The profile card and the timed-quiz page ask for eligibility on load. A
    // learner below the gate is the default empty account; a spec that wants
    // the run registers its own route later and wins.
    await context.route('**/api/timed-quiz**', (route) =>
      route.fulfill({ json: { eligible_xp: 0, gate_xp: 150, due: false, server_now: Date.now() } }),
    )

    await use(context)

    if (escaped.length > 0) {
      throw new Error(
        `this test reached an origin outside the static export:\n  ${escaped.join('\n  ')}\n` +
          'Stub it in the spec with page.route, or add the endpoint to e2e/fixtures.ts if every page ' +
          'needs it. Production must never answer a browser test.',
      )
    }
  },
  // The assistant config is stubbed on the page, not the context. Page routes outrank context
  // routes, and a spec that stubs `**/api/**` with a 404 in a `beforeEach` registers after this
  // fixture, so a context-level stub would lose and the page would log a 404 for config. A spec
  // that wants the assistant on registers its own page route after this one and wins.
  page: async ({ page }, use) => {
    await page.route('**/api/config', (route) => route.fulfill({ json: { ai_enabled: false } }))
    await use(page)
  },
})

export * from '@playwright/test'
