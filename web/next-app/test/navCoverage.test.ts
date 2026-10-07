import { describe, it, expect } from 'vitest'
import { readdirSync, readFileSync, statSync } from 'node:fs'
import { join } from 'node:path'

// Every page must offer a way off itself. A page that renders neither a Header nor the bottom
// tabs is a navigation dead end: the learner arrives there and cannot reach anything else.
//
// This exists because /profile lost both. Commit fb5bde0a replaced its Header and BottomTabs
// with a shadcn sidebar shell, and the shell is desktop-only — shadcn renders it as an overlay
// Sheet below lg, and the rail only collapses behind a trigger that had to be summoned. So on a
// phone, tapping Profile in the tab bar of any other page arrived at a page with no way back,
// which is what was reported. Three of that page's six return branches rendered no chrome at
// all, so the same hole existed between the tabs appearing and the data arriving.
//
// The hole was invisible to every other check. tsc passed: a page missing navigation is still
// valid JSX. The build passed: the file is a real prerendered route. The unit suite passed: no
// test asserted that a page could be left. Only the shape of the code said so, and nothing read
// it. This reads it now.
const APP_DIR = join(process.cwd(), 'app')

// Routes that render nothing of their own and immediately send the learner elsewhere. These are
// exempt because they are not destinations: there is no page here to leave.
const REDIRECT_STUBS = new Set(['concept', 'session', 'docs/system-design'])

/**
 * The invariant is "a learner can get off this page", not "this page renders one of two
 * components". A long-form docs surface can satisfy it with a single back link — /docs and
 * /docs/efficacy both do, and both are deliberately narrower than the app shell, so demanding a
 * tab bar of them would be imposing the wrong shape. Three ways out count: app chrome, a
 * redirect, or a back link. `←` is the project's documented back-link glyph (AGENTS.md allows
 * arrows as typography), so its presence is the signal.
 */
function hasWayOff(src: string): boolean {
  // Imports are stripped first, and the two chrome components are matched as rendered tags.
  // Both details were found by testing this guard against a deliberate regression: matching the
  // bare words `Header` and `BottomTabs` passed a page that imported BottomTabs and rendered
  // nothing, which is precisely the shape of the original bug. An import is not a control, and a
  // guard that cannot tell the difference passes the thing it was written to catch.
  const body = src.replace(/^import .*$/gm, '')
  return /<Header\b/.test(body) || /<BottomTabs\b/.test(body) || /\bredirect\(/.test(body) || /←/.test(body)
}

function pageFiles(dir: string): string[] {
  return readdirSync(dir).flatMap(entry => {
    const full = join(dir, entry)
    if (statSync(full).isDirectory()) return pageFiles(full)
    return entry === 'page.tsx' ? [full] : []
  })
}

describe('every page offers navigation', () => {
  const routes = pageFiles(APP_DIR)

  it('finds the app routes', () => {
    // A guard that scans nothing reports nothing, and a route-level guard that silently stopped
    // scanning would pass forever. Assert the scan is actually reaching the routes.
    expect(routes.length).toBeGreaterThan(8)
    expect(routes.some(r => r.endsWith('profile/page.tsx'))).toBe(true)
  })

  it('renders Header or BottomTabs on every page that is not a redirect stub', () => {
    const bare = routes
      .filter(r => !REDIRECT_STUBS.has(r.split('/').slice(-2)[0]))
      .filter(r => !hasWayOff(readFileSync(r, 'utf8')))
      .map(r => '/' + r.split('/').slice(1, -1).join('/'))

    // The sidebar shell carries no navigation of its own on mobile, so it does not satisfy this.
    expect(bare).toEqual([])
  })

  it('treats the exempt routes as redirect stubs, so the allowance cannot rot', () => {
    // If /concept, /session or /docs/system-design ever grows real content, the exemption above
    // becomes a loophole and the guard silently stops protecting them. Fail here instead.
    for (const route of REDIRECT_STUBS) {
      const src = readFileSync(join(APP_DIR, route, 'page.tsx'), 'utf8')
      expect(src, `${route} is no longer a redirect stub`).toMatch(/redirect\(|router\.(push|replace)/)
    }
  })
})

describe('theme switching is reachable where it is needed', () => {
  const iconModule = readFileSync(join(process.cwd(), 'components/icons/index.tsx'), 'utf8')

  it('exposes the sun-moon icon from the icon module', () => {
    // The theme control is a button whose only content is this icon. If it disappears from the
    // module, every theme control built on it becomes an empty button rather than a build error.
    expect(iconModule).toMatch(/export function SunMoon|SunMoon/)
  })

  it('keeps a theme control in the Header', () => {
    // Header is not the only place a theme control lives any more — the profile sidebar has one
    // too — but Header is rendered by most routes, so if it loses the control, most of the app
    // silently loses light/dark with it. That is the same class of regression as the navigation
    // hole above: valid code, no error, a capability gone.
    const header = readFileSync(join(process.cwd(), 'components/Header.tsx'), 'utf8')
    expect(header).toMatch(/toggleTheme/)
    expect(header).toMatch(/aria-label="Toggle theme"/)
  })

  it('gives the profile sidebar its own theme control', () => {
    const sidebar = readFileSync(join(process.cwd(), 'components/app-sidebar.tsx'), 'utf8')
    expect(sidebar).toMatch(/toggleTheme/)
    expect(sidebar).toMatch(/aria-label="Toggle theme"/)
  })
})