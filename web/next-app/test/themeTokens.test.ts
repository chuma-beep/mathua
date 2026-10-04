import { describe, it, expect } from 'vitest'
import fs from 'node:fs'
import path from 'node:path'

// Opacity modifiers (/10, /50, …) emit no CSS for var()-based theme colors
// in this Tailwind build (verified against the emitted bundle: zero
// color-mix, zero mathua-*/N rules), so `border-mathua-blue/50` renders
// currentColor instead of blue. Theme fills and translucent borders therefore
// use explicit `color-mix` tokens — `mathua-blue-faint`, `mathua-green-faint`,
// `mathua-red-faint` — and accents use solid tokens.
//
// Re-verified rather than trusted: emitting the Tailwind output for
// `border-mathua-green/40`, `border-mathua-green`, `border-green-500/40` and
// `border-mathua-green-faint` produced rules for the second and third only. The
// `/40` on a token class emits *nothing at all*, so it is not a wrong colour, it
// is no declaration — which is why replacing `border-green-500/40` with
// `border-mathua-green/40` would have looked like it worked.
const DEAD_TOKEN_RE =
  /(?:border|bg|text|ring|divide|from|via|to|outline|decoration)-mathua-[a-z-]*\/\d+/g

// Every verdict and alert border has to come from a theme token, so it follows
// light/dark and cannot drift from the label beside it.
//
// This was the last allowlisted exception: `border-green-500/40`,
// `border-red-500/40`, `border-yellow-500/40` and friends emitted correctly but
// were hardcoded Tailwind values, so a verdict border was a *different green*
// from `--accent-green` and the two diverged further in dark mode
// (`green-400` #4ade80 against `--accent-green` #10b981). Eleven such borders
// across five files, plus the reviews-due alert, which was the app's only yellow.
const RAW_PALETTE_BORDER_RE =
  /\bborder-(?:green|red|yellow|amber|orange|lime|emerald|teal|cyan|sky|blue|indigo|violet|purple|fuchsia|pink|rose|slate|gray|zinc|neutral|stone|white|black)-\d+(?:\/\d+)?\b/g

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

  it('draws no border in a raw Tailwind palette colour', () => {
    const root = path.join(__dirname, '..')
    const offenders: string[] = []
    for (const file of [...tsxFiles(path.join(root, 'app')), ...tsxFiles(path.join(root, 'components'))]) {
      const hits = fs.readFileSync(file, 'utf8').match(RAW_PALETTE_BORDER_RE)
      if (hits) offenders.push(`${path.relative(root, file)}: ${[...new Set(hits)].join(', ')}`)
    }
    expect(offenders, `raw-palette borders:\n${offenders.join('\n')}`).toEqual([])
  })
})
