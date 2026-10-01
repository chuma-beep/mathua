import { describe, it, expect } from 'vitest'
import fs from 'node:fs'
import path from 'node:path'

// Opacity modifiers (/10, /50, …) emit no CSS for var()-based theme colors
// in this Tailwind build (verified against the emitted bundle: zero
// color-mix, zero mathua-*/N rules), so `border-mathua-blue/50` renders
// currentColor instead of blue. Default-palette modifiers (green-500/30)
// emit normally and are allowlisted. Theme fills use the explicit
// `mathua-blue-faint` token; accent borders use solid tokens.
const DEAD_TOKEN_RE =
  /(?:border|bg|text|ring|divide|from|via|to|outline|decoration)-mathua-[a-z-]*\/\d+/g

function tsxFiles(dir: string): string[] {
  const out: string[] = []
  for (const entry of fs.readdirSync(dir, { withFileTypes: true })) {
    const full = path.join(dir, entry.name)
    if (entry.isDirectory()) out.push(...tsxFiles(full))
    else if (entry.name.endsWith('.tsx')) out.push(full)
  }
  return out
}

describe('theme opacity modifiers', () => {
  it('uses no opacity modifiers on mathua-* theme colors', () => {
    const root = path.join(__dirname, '..')
    const offenders: string[] = []
    for (const file of [...tsxFiles(path.join(root, 'app')), ...tsxFiles(path.join(root, 'components'))]) {
      const hits = fs.readFileSync(file, 'utf8').match(DEAD_TOKEN_RE)
      if (hits) offenders.push(`${path.relative(root, file)}: ${[...new Set(hits)].join(', ')}`)
    }
    expect(offenders).toEqual([])
  })
})
