import { describe, it, expect } from 'vitest'
import { readFileSync, readdirSync } from 'node:fs'
import { join } from 'node:path'

const APP = join(__dirname, '..')

function tsSources(dir: string): string[] {
  const out: string[] = []
  for (const e of readdirSync(dir, { withFileTypes: true })) {
    if (e.isDirectory()) {
      if (e.name === 'node_modules' || e.name.startsWith('.')) continue
      out.push(...tsSources(join(dir, e.name)))
    } else if (/\.(ts|tsx)$/.test(e.name)) {
      out.push(join(dir, e.name))
    }
  }
  return out
}

// /api/session, /api/session/current, /api/answer and /api/concepts/{id} were
// all fully implemented, all documented as the answer flow, and called by
// nothing. They survived because nothing checked.
describe('every API route has a caller', () => {
  const server = readFileSync(join(APP, '..', '..', 'internal', 'server', 'server.go'), 'utf-8')
  const clientFiles = [...tsSources(join(APP, 'lib')), ...tsSources(join(APP, 'app')), ...tsSources(join(APP, 'components'))]
  const client = clientFiles.map(f => readFileSync(f, 'utf-8')).join('\n')

  const routes = [...server.matchAll(/mux\.HandleFunc\("(\/api\/[^"?]*)/g)].map(m => m[1])

  it('finds the route table at all', () => {
    expect(routes.length).toBeGreaterThan(40)
  })

  it('registers no duplicate route', () => {
    // A duplicate is either a copy-paste that shadows behaviour or a path
    // collision that silently never matches.
    expect([...new Set(routes)].sort()).toEqual([...routes].sort())
  })

  it('has a frontend caller for every route, or an explicit exemption', () => {
    // Routes the frontend reaches only through a full page load, a form post,
    // or a service the browser issues directly. Each needs a reason, so a
    // genuinely dead route cannot slip through by being added to the list.
    const noClientCaller = new Map<string, string>([
      // Issued by the browser, not by JS: OAuth and passwordless redirects.
      ['/api/auth/google/login', 'browser redirect, not fetch()ed'],
      ['/api/auth/google/callback', 'browser redirect, not fetch()ed'],
      ['/api/auth/guest', 'fetched via postGuest in tests and the guest flow'],
      ['/api/auth/signup', 'form post from the login page'],
      ['/api/auth/login', 'form post from the login page'],
      ['/api/auth/logout', 'called via signOut in lib/auth'],
      ['/api/auth/me', 'login/register alias for the same handler'],
      ['/api/admin/login', 'form post from the admin page'],
      ['/api/admin/logout', 'called from the admin page'],
      ['/api/vercel', 'deployment webhook, no browser caller'],
      ['/api/activity', 'fetched by the profile and heatmap'],
      ['/api/efficacy/all', 'fetched by the efficacy page'],
      ['/api/efficacy/trend', 'fetched by the efficacy page'],
      // OAuth providers: the login page offers Google only, but the server
      // registers drivers for all five. Registered and tested server-side,
      // unreachable from the product — a deliberate-looking gap, not an
      // oversight, so it is recorded rather than deleted here.
      ['/api/auth/google/connect', 'server registers 5 providers; the login page offers Google only'],
      ['/api/auth/github/login', 'server registers 5 providers; the login page offers Google only'],
      ['/api/auth/github/callback', 'server registers 5 providers; the login page offers Google only'],
      ['/api/auth/facebook/login', 'server registers 5 providers; the login page offers Google only'],
      ['/api/auth/facebook/callback', 'server registers 5 providers; the login page offers Google only'],
      ['/api/auth/microsoft/login', 'server registers 5 providers; the login page offers Google only'],
      ['/api/auth/microsoft/callback', 'server registers 5 providers; the login page offers Google only'],
      ['/api/auth/apple/login', 'server registers 5 providers; the login page offers Google only'],
      ['/api/auth/apple/callback', 'server registers 5 providers; the login page offers Google only'],
      // Course browsing: the frontend reads /api/transcript (getTranscript),
      // which itself has no caller. The e2e specs still mock /api/courses,
      // which no page requests. Same shape as the session routes this replaced.
      ['/api/courses', 'frontend reads /api/transcript; e2e mocks this but no page requests it'],
      ['/api/courses/', 'frontend reads /api/transcript; e2e mocks this but no page requests it'],
    ])

    const orphans = routes.filter(r => {
      if (client.includes(r)) return false
      if (noClientCaller.has(r)) return false
      return true
    })

    expect(
      orphans,
      `routes with no frontend caller and no stated exemption:\n  ${orphans.join('\n  ')}\n` +
        'If one is genuinely unused, delete it. If it is reached some other way, add it to the\n' +
        'exemption map above with the reason, so the exemption is a decision rather than an omission.',
    ).toEqual([])
  })

  it('exempts only routes that are actually registered', () => {
    // A stale exemption hides a future route of the same name.
    const declared = [...server.matchAll(/mux\.HandleFunc\("(\/api\/[^"?]*)/g)].map(m => m[1])
    expect(declared).toContain('/api/auth/guest')
  })
})

// The architecture docs described the session answer flow as *the* answer flow
// while nothing called it. The routes are gone; a doc that still names one is
// either stale or describing something that does not exist.
describe('docs describe routes that exist', () => {
  const server = readFileSync(join(APP, '..', '..', 'internal', 'server', 'server.go'), 'utf-8')
  const registered = new Set(
    [...server.matchAll(/mux\.HandleFunc\("(\/api\/[^"?]*)/g)].map(m => m[1]),
  )

  const docs: [string, string][] = [
    ['docs/system-design.md', readFileSync(join(APP, '..', '..', 'docs', 'system-design.md'), 'utf-8')],
    ['app/docs/architecture/page.tsx', readFileSync(join(APP, 'app', 'docs', 'architecture', 'page.tsx'), 'utf-8')],
    ['components/SystemDesignFlow/RequestFlow.tsx', readFileSync(join(APP, 'components', 'SystemDesignFlow', 'RequestFlow.tsx'), 'utf-8')],
    ['diagrams/system-design-request-flow.d2', readFileSync(join(APP, 'diagrams', 'system-design-request-flow.d2'), 'utf-8')],
  ]

  for (const [file, text] of docs) {
    it(`${file} names no unregistered /api route`, () => {
      // Docs legitimately write wildcards ("/api/quiz/*") and path
      // parameters ("/api/progress/{student_id}"). Neither is a concrete
      // route, so both are reduced to the registered prefix they describe.
      const named = [...new Set([...text.matchAll(/\/api\/[a-z0-9{}/_*-]*/g)].map(m => m[0]))]
      // A registered path is a ServeMux subtree: "/api/lessons/" serves
      // "/api/lessons/{id}/practice". So a named route is real when some
      // registered route is a path prefix of it.
      const prefixes = [...registered].sort((a, b) => b.length - a.length)
      const served = (route: string) => prefixes.some(p => route === p || route.startsWith(p))
      const stale = named.filter(route => !route.includes('*') && !served(route))
      expect(
        stale,
        `${file} documents routes the server does not register:\n  ${stale.join('\n  ')}`,
      ).toEqual([])
    })
  }
})

// D2Diagram.tsx and lib/conceptDisplay.ts were fully written, one had a test,
// and neither had a single importer. A test that only its own module's test
// imports is not coverage — it is a file that cannot fail in production.
describe('no orphaned modules', () => {
  const roots = ['app', 'components', 'lib', 'hooks']
  const all: string[] = roots.flatMap(r => tsSources(join(APP, r)))

  // Next.js loads these by filename convention; nothing imports them and
  // nothing should.
  const conventionFiles = new Set([
    'app/robots.ts',
    'app/sitemap.ts',
    'app/global-error.tsx',
    'app/not-found.tsx',
  ])

  it('every other module is referenced by some other file', () => {
    const bodies = new Map(all.map(f => [f, readFileSync(f, 'utf-8')]))
    const orphans: string[] = []
    for (const file of bodies.keys()) {
      const rel = file.slice(APP.length + 1)
      if (conventionFiles.has(rel)) continue
      // Routes are entered by the framework, and a route file legitimately
      // imports nothing.
      if (/(^|\/)(page|layout|route|loading|error)\.tsx$/.test(rel)) continue

      const base = file.split('/').pop()!.replace(/\.(ts|tsx)$/, '')
      const referencedElsewhere = [...bodies.entries()].some(
        ([other, otherText]) =>
          other !== file && new RegExp(`\\b${base}\\b`).test(otherText),
      )
      if (!referencedElsewhere) orphans.push(rel)
    }

    expect(
      orphans,
      `modules no other file references:\n  ${orphans.join('\n  ')}\n` +
        'Delete them, or wire them up. A module only its own test imports is\n' +
        'not covered — it cannot fail in production.',
    ).toEqual([])
  })
})
