import { describe, it, expect } from 'vitest'
import { readFileSync, readdirSync } from 'node:fs'
import { join } from 'node:path'

// Learner-facing copy that states a number the app does not implement is worse
// than copy that states nothing: a learner sizes their effort against it, and
// every one of these drifted silently when ADR-020 rescaled the XP economy.
//
// The denylist is deliberately mechanical. Each entry names something the
// learner can read, so a regression is caught by grepping rather than by
// noticing. The gate for the concept count is separate, below, because that one
// is generated rather than written.
const APP = join(__dirname, '..')

function appSources(dir: string): string[] {
  const out: string[] = []
  for (const entry of readdirSync(dir, { withFileTypes: true })) {
    if (entry.isDirectory()) {
      out.push(...appSources(join(dir, entry.name)))
    } else if (/\.(tsx|ts)$/.test(entry.name)) {
      out.push(join(dir, entry.name))
    }
  }
  return out
}

describe('learner-facing copy tells the truth', () => {
  const files = [...appSources(join(APP, 'app')), ...appSources(join(APP, 'components')), ...appSources(join(APP, 'lib'))]

  // XP is effort-derived and scaled by speed and streak (internal/xp), and the
  // quiz gate is 50 (engine.QuizGateXP). Any of these numbers appearing in UI
  // copy means the economy was rescaled and the copy was not.
  const stale = [
    { re: /\b20 XP\b/, what: 'a flat per-answer XP award the economy does not use' },
    { re: /\b150 XP\b/, what: 'the pre-rescale quiz gate; QuizGateXP is 50' },
    { re: /\b641\b/, what: 'a superseded concept count; the graph holds 657' },
    { re: /2 in a row to advance before you practice/, what: 'an advance gate Study does not have' },
    { re: /worked example first, then answer/, what: 'an answer step Study does not have' },
  ]

  for (const { re, what } of stale) {
    it(`never claims ${what}`, () => {
      const hits: string[] = []
      for (const file of files) {
        const text = readFileSync(file, 'utf-8')
        text.split('\n').forEach((line, i) => {
          if (re.test(line)) hits.push(`${file.replace(APP, 'web/next-app')}:${i + 1}: ${line.trim()}`)
        })
      }
      expect(hits).toEqual([])
    })
  }

  it('does not promise an instant review on a quiz miss', () => {
    // A miss shows the worked solution inline; anything worth revisiting is a
    // link list at the end. Nothing opens on the spot.
    const files2 = [join(APP, 'components', 'BriefingCard.tsx')]
    const text = files2.map(f => readFileSync(f, 'utf-8')).join('\n')
    expect(text).not.toMatch(/review it on the spot/)
    expect(text).not.toMatch(/opens a short review/)
  })
})

// The concept count is generated, so it is checked against the corpus rather
// than denylisted: counts.py owns the number and README badges.
describe('stated concept count matches the corpus', () => {
  it('the graph holds the number the meta description advertises', () => {
    const dir = join(__dirname, '..', '..', '..', 'data', 'concepts')
    // Unique ids, not raw entries: data/concepts/enrichment.json holds
    // overrides for concepts that already exist in their own domain file, and
    // counts.py skips it for the same reason.
    const ids = new Set<string>()
    for (const name of readdirSync(dir)) {
      if (!name.endsWith('.json')) continue
      const raw = JSON.parse(readFileSync(join(dir, name), 'utf-8'))
      const items = Array.isArray(raw) ? raw : (raw.concepts ?? raw.items ?? [])
      for (const c of items) {
        if (c && typeof c === 'object' && 'id' in c) ids.add(c.id as string)
      }
    }
    const layout = readFileSync(join(APP, 'app', 'layout.tsx'), 'utf-8')
    const claimed = [...layout.matchAll(/(\d+) concepts/g)].map(m => Number(m[1]))
    expect(claimed.length).toBeGreaterThan(0)
    for (const n of claimed) {
      expect(n, `layout.tsx advertises ${n} concepts, corpus has ${ids.size}`).toBe(ids.size)
    }
  })
})

// A label that misdescribes its destination is a small lie with a large
// consequence: a learner who clicks it ends up somewhere they did not choose,
// and blames the product. These assert the specific corrections rather than
// generalising, because each one was a real mislabel.
describe('navigation labels name where they actually go', () => {
  const read = (p: string) => readFileSync(join(APP, p), 'utf-8')

  it('does not label a link to /profile as "Next up"', () => {
    const page = read(join('app', 'learn', 'page.tsx'))
    expect(page).not.toMatch(/>← Next up</)
  })

  it('does not call the worked-example gate a start', () => {
    // The loop has already begun by the time this button appears; it is the
    // only way past the example, and calling it "Start practicing" implies the
    // learner has not started.
    const stepper = read(join('components', 'LearnStepper.tsx'))
    expect(stepper).not.toMatch(/Start practicing →</)
  })

  it('does not print a raw concept id as a back-link label', () => {
    const stepper = read(join('components', 'LearnStepper.tsx'))
    expect(stepper).not.toMatch(/← Back to \{returnTo\}/)
    expect(stepper).toMatch(/returnToLabel/)
  })

  it('does not end a finished review by offering the reference library', () => {
    // Review is a Learn sub-mode; "nothing is due" is not a reason to send
    // someone from doing into browsing.
    const host = read(join('app', 'review', 'ReviewHost.tsx'))
    const emptyState = host.slice(host.indexOf('All caught up'))
    expect(emptyState).not.toMatch(/Open Study/)
    expect(emptyState).toMatch(/href="\/learn"/)
  })

  it('renders the explanation as the body of a review result, not the verdict token', () => {
    const host = read(join('app', 'review', 'ReviewHost.tsx'))
    // "Incorrect" above the working is noise; the explanation is the body and
    // feedback is only the fallback when there is no explanation.
    expect(host).toMatch(/explanation \? \(/)
    expect(host).not.toMatch(/text-mathua-muted text-xs mt-2/)
  })
})
