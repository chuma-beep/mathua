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
