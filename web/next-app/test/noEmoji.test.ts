import { describe, it, expect } from 'vitest'
import fs from 'node:fs'
import path from 'node:path'

// No emoji, anywhere we author. AGENTS.md states the rule; this states the evidence.
//
// The rule exists because of what an emoji costs rather than because they are ugly. Four were
// live in the app when this was written — U+23F3 on the reviews-due line, U+23F8 on the pause
// banner, U+23F5 on the graph flow toggle, and one already guarded by a test — and each had a
// concrete defect behind it:
//
//   - U+23F5 sat inside a button with `fontSize: 10` and rendered at an unpredictable baseline
//     next to the word "Flow".
//   - The hourglass was there because it was the obvious way to say "due"; an icon says it with
//     an artwork that inherits `currentColor` and a theme token, and can be animated.
//   - An emoji's glyph, colour and baseline are the *font's*, so it is not a theme surface. It
//     cannot follow light/dark, cannot be sized consistently against mono text, and renders
//     differently per platform.
//
// The exclusions below are deliberate and narrow. Arrows, check marks, ballot marks and box
// drawing are typography that Mathua uses deliberately, not emoji: `→` in a button, `✓` for a
// correct answer, `·` as a bullet. `⌈`/`⌉` are floor and ceiling brackets in set notation —
// `internal/generator/discrete/generators.go` uses them in generated math, and mistaking them for
// emoji would corrupt mathematics.

// test/ -> next-app -> web -> repo root. Getting this wrong makes every assertion below pass
// vacuously, because the scan then visits nothing and finds nothing: which is exactly what
// happened the first time this test was written, and why `scans the files it claims to` is
// asserted rather than assumed.
const ROOT = path.resolve(__dirname, '..', '..', '..')

/** Directories and files we author. Vendored trees are excluded on purpose. */
const AUTHORED = [
  'web/next-app/app',
  'web/next-app/components',
  'web/next-app/lib',
  'web/next-app/test',
  'web/next-app/e2e',
  'web/next-app/scripts',
  'internal',
  'cmd',
  'docs',
  'CONTEXT.md',
  'README.md',
  'DESIGN.md',
  'CONTRIBUTING.md',
]

/** Extensions we author in. */
const SUFFIXES = new Set(['.ts', '.tsx', '.js', '.mjs', '.go', '.py', '.md', '.css', '.json'])

/**
 * Emoji and pictograph blocks, minus the characters Mathua uses as typography.
 *
 * Deliberately conservative in what it *flags* and permissive in what it *excludes*: a false
 * positive here costs a reviewer's time, and the exclusions are the same ones a reader would
 * apply by eye.
 */
const EMOJI =
  /[\u{1F000}-\u{1FAFF}\u{2300}-\u{23FF}\u{2600}-\u{26FF}\u{2700}-\u{27BF}\u{2B00}-\u{2BFF}\u{FE0F}]/gu

/** Typography, not emoji. Each entry says what it is, because "trust me" is not a rule. */
const ALLOWED: Record<string, string> = {
  '→': 'arrow, used in button labels such as "Review Now"',
  '←': 'arrow, used in back links',
  '↑': 'arrow, used in sort and step controls',
  '↓': 'arrow, used in sort and step controls',
  '↔': 'arrow, used in scale and axis labels',
  '⇒': 'arrow, used as an implies operator in prose',
  '⇔': 'arrow, used as an if-and-only-if operator',
  '⟶': 'arrow, used in load-bearing note titles',
  '⟹': 'arrow, used as a therefore operator',
  '✓': 'check mark, used for a correct answer',
  '✗': 'ballot x, used for an incorrect answer',
  '✔': 'check mark',
  '✕': 'ballot x, used as a dismiss control',
  '✖': 'ballot x, used as a dismiss control',
  '✦': 'four-pointed star, decorative',
  '✧': 'outlined star, decorative',
  '·': 'middle dot, used as a bullet',
  '•': 'bullet, used in markdown lists',
  '⌈': 'floor bracket, set notation',
  '⌉': 'ceiling bracket, set notation',
  '⌊': 'floor bracket, set notation',
  '⌋': 'ceiling bracket, set notation',
  '∈': 'element of, mathematics',
  '∑': 'summation, mathematics',
  '∏': 'product, mathematics',
  '√': 'square root, mathematics',
  '∞': 'infinity, mathematics',
  '≈': 'approximately equal, mathematics',
  '×': 'multiplication sign, mathematics',
  '÷': 'division sign, mathematics',
  '≠': 'not equal, mathematics',
  '≥': 'greater or equal, mathematics',
  '≤': 'less or equal, mathematics',
  '°': 'degree sign, mathematics',
  '²': 'superscript two, mathematics',
  '′': 'prime, mathematics',
  '…': 'ellipsis, used in truncated prose and generated content',
  '—': 'em dash, used as punctuation throughout the docs',
  '–': 'en dash, used in numeric ranges such as 5–15',
  '§': 'section sign, used in ADR references',
  '†': 'dagger, footnote marker',
}

/**
 * Whether the guard would flag a character. This is the whole rule, and the tests assert on it
 * rather than on the regex: a character can be safe either because the regex does not match it
 * (arrows — U+2190..U+21FF is not in an emoji range) or because it is explicitly excluded
 * (check marks, which share a block with emoji). Asserting `matches(EMOJI)` would demand the
 * wrong one of those and would push a future edit to widen the regex for no reason.
 */
function flagged(ch: string): boolean {
  return ch.match(EMOJI) !== null && !(ch in ALLOWED)
}

function* authoredFiles(entry: string): Generator<string> {
  const abs = path.join(ROOT, entry)
  if (!fs.existsSync(abs)) return
  const st = fs.statSync(abs)
  if (st.isFile()) {
    yield abs
    return
  }
  for (const dirent of fs.readdirSync(abs, { withFileTypes: true })) {
    // Vendored and generated trees are not ours to police.
    if (dirent.name === 'node_modules' || dirent.name === 'out' || dirent.name.startsWith('.next')) {
      continue
    }
    yield* authoredFiles(path.join(entry, dirent.name))
  }
}

interface Scan {
  files: number
  found: { file: string; char: string; line: number }[]
}

function scan(): Scan {
  const found: { file: string; char: string; line: number }[] = []
  let files = 0
  for (const entry of AUTHORED) {
    for (const file of authoredFiles(entry)) {
      files++
      if (!SUFFIXES.has(path.extname(file))) continue
      const rel = path.relative(ROOT, file)
      let text: string
      try {
        text = fs.readFileSync(file, 'utf8')
      } catch {
        continue
      }
      text.split('\n').forEach((line, i) => {
        for (const m of line.matchAll(EMOJI)) {
          const ch = m[0]
          if (ch in ALLOWED) continue
          found.push({ file: rel, char: ch, line: i + 1 })
        }
      })
    }
  }
  return { files, found }
}

function offenders(): { file: string; char: string; line: number }[] {
  return scan().found
}

describe('no emoji in authored code', () => {
  // A guard that reads nothing reports nothing, and "no emoji" is exactly what it reports when
  // the path is wrong. This asserts the scan is looking at a plausible number of files.
  it('scans the files it claims to', () => {
    const visited = new Set<string>()
    for (const entry of AUTHORED) {
      for (const file of authoredFiles(entry)) visited.add(path.relative(ROOT, file))
    }
    // Named sentinels rather than a count. A count can be satisfied by the wrong directory
    // entirely, and this test previously passed vacuously because ROOT was one level short.
    for (const sentinel of [
      'web/next-app/app/profile/page.tsx',
      'web/next-app/components/AttemptList.tsx',
      'internal/concepts/mastery_units.go',
      'CONTEXT.md',
    ]) {
      expect([...visited], `${sentinel} was not scanned`).toContain(sentinel)
    }
    expect(visited.size).toBeGreaterThan(400)
  })

  it('finds none', () => {
    const found = offenders()
    const detail = found
      .map((o) => `  ${o.file}:${o.line}  U+${o.char.codePointAt(0)!.toString(16).toUpperCase().padStart(4, '0')} ${o.char}`)
      .join('\n')
    expect(
      found.length === 0 ? 'no emoji' : `emoji found:\n${detail}`,
    ).toBe('no emoji')
  })

  // A guard that cannot fail is not a guard. These are the characters it must catch, asserted
  // against the classifier rather than trusted — the same reason the maths exclusions exist:
  // a regex that silently stops matching is worse than no regex.
  it('flags a real emoji', () => {
    // Written as escapes, not literals: a fixture emoji in this file would be flagged by the
    // very guard it defines, and excluding this file from the scan would be a hole.
    for (const ch of ['\u23F3', '\u2705', '\u274C', '\u{1F3B2}', '\u{1FAA6}', '\u26A0\uFE0F', '\u23F8', '\u23F5', '\u{1F44D}']) {
      expect(flagged(ch), `${ch} should be flagged`).toBe(true)
    }
  })

  // The exclusions are the risky half, because each one is a hole. Pin the ones that matter:
  // a false positive on `⌈`/`⌉` would put a lint error inside generated mathematics.
  it('does not flag typography Mathua uses deliberately', () => {
    for (const ch of ['\u2192', '\u2190', '\u2713', '\u2717', '\u2715', '\u00B7', '\u2308', '\u2309', '\u230A',
      '\u2208', '\u2211', '\u221A', '\u221E', '\u00D7', '\u2265', '\u2014', '\u2026', '\u00A7']) {
      expect(flagged(ch), `${ch} should not be flagged`).toBe(false)
    }
  })

  it('allows every exclusion to name itself', () => {
    // An undocumented exclusion is indistinguishable from a mistake.
    for (const [ch, why] of Object.entries(ALLOWED)) {
      expect(why, `${ch} has no stated reason`).toBeTruthy()
      expect(why.length, `${ch} reason is too terse to review`).toBeGreaterThan(8)
    }
  })
})