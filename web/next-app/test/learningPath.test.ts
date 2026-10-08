import { describe, it, expect } from 'vitest'
import { readdirSync, readFileSync, statSync } from 'node:fs'
import { join, relative } from 'node:path'

// One learning path.
//
// `/study` was a browsable reference library sitting beside `/learn` in every navigation
// surface: the desktop header, the mobile tab bar, the profile sidebar, the footer, the landing
// page, the docs menu, and the creator's note. It was also reachable from the graph, the
// diagnostic results, the quiz's remedial list, the attempt history, the quiz's "focus next"
// list, and two "Reference" buttons inside `/learn` itself.
//
// That is not a detail of information architecture. Learn is the only surface that asks a
// question, grades it and explains the answer, so a route that explains without asking is a
// route around the only thing that counts. A learner who found it could read their way to a
// concept; a learner who did not was sent there first.
//
// So this file pins the closure rather than the layout. Every assertion below is a way the
// library could come back, and each one is stated as a property of the source rather than of a
// snapshot, because the failure mode is a link being *added*, and a snapshot would have to be
// regenerated to notice.
//
// What is deliberately still allowed: the `/api/study/answer` seam (Learn's answer endpoint —
// a name, not a destination), the server-side comments that explain why the library closed, and
// the redirect stubs that forward old links. A guard that forbade the string `/study` outright
// would be deleted on its first false positive, and a deleted guard protects nothing.

const APP = join(__dirname, '..')

function sources(dir: string): string[] {
  const out: string[] = []
  for (const entry of readdirSync(dir)) {
    if (entry === 'node_modules' || entry.startsWith('.') || entry === 'out') continue
    const full = join(dir, entry)
    if (statSync(full).isDirectory()) out.push(...sources(full))
    else if (/\.(ts|tsx)$/.test(entry)) out.push(full)
  }
  return out
}

/** Every authored surface a learner can be shown: routes, components and the API client. */
const LEARNER_SURFACES = [...sources(join(APP, 'app')), ...sources(join(APP, 'components')), ...sources(join(APP, 'lib'))]

/**
 * Strip comments, so a file can explain why `/study` closed without tripping the guard.
 *
 * This is the part that keeps the guard honest. Without it, every explanatory comment added
 * later would need a carve-out, carves-out accumulate, and the check quietly stops meaning
 * anything — which is the exact shape of the guards this repo already had to fix twice.
 */
function code(src: string): string {
  return src
    .replace(/\/\*[\s\S]*?\*\//g, ' ')
    .replace(/(^|[^:])\/\/[^\n]*/g, '$1 ')
}

describe('the reference library stays closed', () => {
  it('scans the surfaces it claims to scan', () => {
    // A guard that reads nothing reports nothing, and this one would silently pass forever if
    // the paths moved. Assert the scan is reaching real files first.
    expect(LEARNER_SURFACES.length).toBeGreaterThan(100)
    expect(LEARNER_SURFACES.some((f) => f.endsWith(join('lib', 'nav.ts')))).toBe(true)
    expect(LEARNER_SURFACES.some((f) => f.endsWith(join('app', 'learn', 'page.tsx')))).toBe(true)
  })

  it('is linked from no learner surface', () => {
    const offenders: string[] = []
    for (const file of LEARNER_SURFACES) {
      const body = code(readFileSync(file, 'utf-8'))
      // Every way a destination can be written: href=, push(), replace(), a literal in an
      // array of nav items. Matching only `href="/study` would miss `push('/study')` and
      // `'/study'` in lib/nav.ts, which is where two of them lived.
      if (/['"`]\/study(\?[^'"`]*)?['"`]/.test(body)) {
        offenders.push(relative(APP, file))
      }
    }
    expect(
      offenders,
      `these surfaces still point at /study:\n  ${offenders.join('\n  ')}`,
    ).toEqual([])
  })

  it('is named nowhere a learner reads', () => {
    // A link retargeted to /learn but still labelled "Study" is the half-finished version of
    // this change: the route is unused, the invitation is not. Comments are excluded by `code`.
    const offenders: string[] = []
    for (const file of LEARNER_SURFACES) {
      const body = code(readFileSync(file, 'utf-8'))
      if (/>\s*Study\s*</.test(body)) offenders.push(relative(APP, file))
    }
    expect(
      offenders,
      `these surfaces still show a "Study" control:\n  ${offenders.join('\n  ')}`,
    ).toEqual([])
  })

  it('is still reachable only as a forwarder', () => {
    const page = code(readFileSync(join(APP, 'app', 'study', 'page.tsx'), 'utf-8'))
    expect(page).toMatch(/router\.replace/)
    expect(page).toMatch(/\/learn/)
    // A forwarder that can reach its own route is a loop, not a redirect.
    expect(page).not.toMatch(/['"`]\/study/)
  })

  it('is not indexed', () => {
    // `/study` was the one query-driven page in the sitemap. A closed route that stays in the
    // sitemap is advertised to a search engine as a destination and then bounces.
    expect(code(readFileSync(join(APP, 'lib', 'site.ts'), 'utf-8'))).not.toMatch(/\/study/)
  })
})

describe('/learn is the one learning surface', () => {
  it('carries the reference material rather than linking out of it', () => {
    const stepper = code(readFileSync(join(APP, 'components', 'LearnStepper.tsx'), 'utf-8'))
    // The material used to sit behind "Reference" on this page and on the done card, both
    // pointing at /study. It is now a panel here, so the panel is the thing to assert.
    expect(stepper).toMatch(/ConceptReference/)
    expect(stepper).not.toMatch(/href=\{`\/study/)
  })

  it('offers no control that would record exposure as progress', () => {
    // §10 of the closure: reading is not evidence. There is no "mark as learned", no "I studied
    // this", no skip — and nothing on the page may award mastery outside `submitStudyAnswer`.
    const stepper = code(readFileSync(join(APP, 'components', 'LearnStepper.tsx'), 'utf-8'))
    const reference = code(readFileSync(join(APP, 'components', 'ConceptReference.tsx'), 'utf-8'))
    for (const src of [stepper, reference]) {
      expect(src).not.toMatch(/mark as (read|learned|done)/i)
      expect(src).not.toMatch(/i (studied|read) this/i)
      expect(src).not.toMatch(/skip practice/i)
    }
    // Only the answer seam may move progress.
    expect(reference).not.toMatch(/submitStudyAnswer|setBand|setCorrectCount|AddXP/)
  })

  it('keeps the orientation surfaces free of any link into the learning loop', () => {
    // Graph, domains and the concept route orient; they do not teach. A concept chip on the
    // graph used to push `/study?concept=`, which is the escape hatch in its purest form:
    // one click from orientation to an article, with no question anywhere in between.
    //
    // `/domains` is the exception and must stay the exception: its rows link onward to
    // `/learn?concept=` because that is the "Continue: <next concept>" action, and a
    // curriculum page that cannot start anything is a dead end.
    const graph = code(readFileSync(join(APP, 'app', 'graph', 'page.tsx'), 'utf-8'))
    expect(graph).not.toMatch(/\/study/)

    // The graph's own component decides the node action from the concept's state, and both
    // halves of that decision live in one file so they cannot disagree.
    const flow = code(readFileSync(join(APP, 'components', 'ConceptGraphFlow.tsx'), 'utf-8'))
    expect(flow).toMatch(/countsAsMastered/)
    expect(flow).not.toMatch(/\/study/)
  })

  it('sends a concept that is already demonstrated to retrieval practice, not a re-teach', () => {
    // The graph node action branches on the concept's state. Asserting only that /learn is not
    // linked would pass with a single unconditional `/learn?concept=` link, which re-teaches a
    // mastered concept — the tutorial purgatory ADR-037 removed from the scheduler's own path.
    const flow = code(readFileSync(join(APP, 'components', 'ConceptGraphFlow.tsx'), 'utf-8'))
    expect(flow).toMatch(/selectedLearned \? '\/review' : `\/learn\?concept=/)
  })
})