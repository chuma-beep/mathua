import { describe, it, expect } from 'vitest'
import fs from 'node:fs'
import path from 'node:path'

// One icon surface: components/icons.
//
// There were three. `lucide-react` supplied 19 static icons, `components/icons/compass.tsx` and
// `sun-moon.tsx` were hand-written animated icons on `motion/react`, and `@animateicons/react`
// supplied three more. Three systems meant three ways to size an icon and three ways to colour
// one, and the animated wrapper has a footgun worth guarding: it renders a `<div>` sized by a
// `size` prop, and its `color` default is not `currentColor`. An icon added without going through
// the wrapper therefore does not follow light and dark mode and does not keep the 44px collapsed
// sidebar rail even — both invisible in review, both visible on screen.

import { APP_ROOT } from './helpers/roots'

const ROOT = APP_ROOT
const SOURCE_DIRS = ['app', 'components', 'lib']
const ICON_MODULE = path.join('components', 'icons', 'index.tsx')
const FORBIDDEN = ['lucide-react', 'motion/react', 'framer-motion']

function sourceFiles(): { rel: string; text: string }[] {
  const out: { rel: string; text: string }[] = []
  const walk = (dir: string) => {
    for (const entry of fs.readdirSync(dir, { withFileTypes: true })) {
      if (entry.name === 'node_modules' || entry.name.startsWith('.')) continue
      const abs = path.join(dir, entry.name)
      if (entry.isDirectory()) walk(abs)
      else if (/\.(ts|tsx)$/.test(entry.name)) {
        out.push({ rel: path.relative(ROOT, abs), text: fs.readFileSync(abs, 'utf8') })
      }
    }
  }
  for (const d of SOURCE_DIRS) walk(path.join(ROOT, d))
  return out
}

describe('icon surface', () => {
  it('scans the files it claims to', () => {
    // A guard that reads nothing reports nothing, and this repo has three instances of that: a
    // ROOT one level short in test/noEmoji.test.ts, and two more while writing the other guards.
    // A floor alone would be satisfied by the wrong directory, so named files are asserted too.
    const files = sourceFiles()
    const rel = files.map((f) => f.rel)
    expect(files.length).toBeGreaterThan(100)
    for (const sentinel of ['components/icons/index.tsx', 'app/profile/page.tsx', 'lib/api.ts']) {
      expect(rel, `${sentinel} was not scanned`).toContain(sentinel)
    }
  })

  it('nothing imports an icon library directly except the icon module', () => {
    const offenders: string[] = []
    for (const { rel, text } of sourceFiles()) {
      if (rel === ICON_MODULE) continue
      for (const lib of FORBIDDEN) {
        if (new RegExp(`from ['"]${lib.replace('/', '\\/')}['"]`).test(text)) {
          offenders.push(`${rel} imports ${lib}`)
        }
      }
    }
    expect(
      offenders.length === 0 ? 'icon libraries are imported only by components/icons' : offenders.join('\n'),
    ).toBe('icon libraries are imported only by components/icons')
  })

  it('the icon module pins the box and the colour', () => {
    const mod = fs.readFileSync(path.join(ROOT, ICON_MODULE), 'utf8')
    // The wrapper is the only thing standing between a forgotten prop and an icon that does not
    // follow the theme, so both defaults are asserted rather than assumed.
    expect(mod).toMatch(/size = 16/)
    expect(mod).toMatch(/color="currentColor"/)
  })

  it('every nav surface takes its icons from the module', () => {
    for (const f of ['components/app-sidebar.tsx', 'components/BottomTabs.tsx', 'components/Header.tsx']) {
      const text = fs.readFileSync(path.join(ROOT, f), 'utf8')
      expect(text, `${f} should import from components/icons`).toMatch(/from '\.\/icons'/)
    }
  })
})
