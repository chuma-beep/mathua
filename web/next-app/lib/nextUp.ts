import type { ConceptProgress, DailyActivity, WeaknessRes } from './api'

export type NextUpKind = 'review' | 'weakness' | 'resume' | 'diagnostic' | 'browse'

export interface NextUp {
  kind: NextUpKind
  badge: string
  title: string
  detail: string
  href: string
  cta: string
}

export interface NextUpInput {
  dueReviews: number
  weaknesses: WeaknessRes | null
  progress: Record<string, ConceptProgress>
  activity: DailyActivity[]
  diagnosticCompleted: boolean
  conceptsMastered: number
}

function isMastered(p: ConceptProgress | undefined): boolean {
  return (p?.status ?? '').toLowerCase() === 'mastered'
}

function topWeakness(input: NextUpInput): { id: string; label: string } | null {
  const entries = input.weaknesses?.by_domain ?? {}
  let best: { id: string; label: string; weakness: number } | null = null
  for (const items of Object.values(entries)) {
    for (const it of items) {
      if (isMastered(input.progress[it.id])) continue
      if (!best || it.weakness > best.weakness) best = it
    }
  }
  return best
}

function mostRecentInProgress(input: NextUpInput): string | null {
  for (let i = input.activity.length - 1; i >= 0; i--) {
    const day = input.activity[i]
    for (let j = day.concepts.length - 1; j >= 0; j--) {
      const cid = day.concepts[j]
      const p = input.progress[cid]
      // In progress = has a progress record but not yet mastered
      if (p && !isMastered(p)) return cid
    }
  }
  return null
}

/**
 * Deterministic Next-up selector. Priority:
 * 1) due reviews → review session host
 * 2) weakest non-mastered concept → Learn deep link (supported ?concept= param)
 * 3) most recently touched in-progress concept → Learn deep link (resume)
 * 4) brand-new user without Diagnostic → Diagnostic entry
 * 5) fallback → Study library
 *
 * selectNextUp is the single-directive fallback. The task shelf
 * (selectShelf) is the primary surface: the algorithm proposes, the
 * learner disposes.
 */
export function selectNextUp(input: NextUpInput): NextUp {
  if (input.dueReviews > 0) {
    return {
      kind: 'review',
      badge: 'Due now',
      title: `${input.dueReviews} concept${input.dueReviews !== 1 ? 's' : ''} due for review`,
      detail: 'Spaced repetition — review before decay',
      href: '/review',
      cta: 'Review now →',
    }
  }

  const weak = topWeakness(input)
  if (weak) {
    return {
      kind: 'weakness',
      badge: 'Recommended',
      title: `Study ${weak.label}`,
      detail: 'Weakest concept — guided Learn session',
      href: `/learn?concept=${encodeURIComponent(weak.id)}`,
      cta: 'Start learning →',
    }
  }

  const resume = mostRecentInProgress(input)
  if (resume) {
    return {
      kind: 'resume',
      badge: 'Continue',
      title: `Resume ${resume}`,
      detail: 'Pick up where you left off',
      href: `/learn?concept=${encodeURIComponent(resume)}`,
      cta: 'Continue →',
    }
  }

  if (input.conceptsMastered === 0 && !input.diagnosticCompleted) {
    return {
      kind: 'diagnostic',
      badge: 'Recommended',
      title: 'Take a Diagnostic to find your frontier',
      detail: 'Adaptive — finds where to start',
      href: '/onboard',
      cta: 'Start Diagnostic →',
    }
  }

  return {
    kind: 'browse',
    badge: 'Study',
    title: 'Browse the Study library',
    detail: 'Pick a Lesson — worked example first, then answer',
    href: '/study',
    cta: 'Browse Study →',
  }
}

// ── Task shelf: algorithm proposes (eligible only), learner disposes ──

export type ShelfKind = 'new' | 'review' | 'weakness' | 'resume' | 'diagnostic' | 'browse'

export interface ShelfItem {
  kind: ShelfKind
  badge: string
  title: string
  detail: string
  href: string
  cta: string
  xp: number
}

export interface CatalogEntry {
  id: string
  label: string
  prerequisites?: string[]
}

export interface ShelfInput extends NextUpInput {
  catalog: CatalogEntry[]
}

export interface ShelfCandidate {
  kind: 'new' | 'weakness' | 'resume'
  id: string
  label: string
  weakness: number
}

// SelectionPolicy turns eligible candidates into a shelf. The v1 fixed
// policy below is a starting weight set; later policies score candidates on
// readiness, urgency, value, difficulty, diversity, pace, and relevance.
export type SelectionPolicy = (cands: ShelfCandidate[], input: ShelfInput) => ShelfItem[]

function xpFor(id: string, kind: ShelfItem['kind']): number {
  if (kind === 'review') return 5
  if (kind === 'diagnostic' || kind === 'browse') return 0
  return id.endsWith('.word') ? 15 : 10
}

function learnItem(kind: ShelfItem['kind'], badge: string, id: string, label: string, detail: string, cta: string): ShelfItem {
  return { kind, badge, title: label, detail, href: `/learn?concept=${encodeURIComponent(id)}`, cta, xp: xpFor(id, kind) }
}

// buildCandidates lists eligible concepts: available (prereqs mastered),
// unmastered, weakest-first. UI never names selection internals.
export function buildCandidates(input: ShelfInput): ShelfCandidate[] {
  const weakById = new Map<string, { label: string; weakness: number }>()
  for (const items of Object.values(input.weaknesses?.by_domain ?? {})) {
    for (const it of items) {
      const prev = weakById.get(it.id)
      if (!prev || it.weakness > prev.weakness) weakById.set(it.id, { label: it.label, weakness: it.weakness })
    }
  }
  const recent = new Set<string>()
  for (const day of input.activity) {
    for (const cid of day.concepts) recent.add(cid)
  }
  const out: ShelfCandidate[] = []
  for (const c of input.catalog) {
    if (isMastered(input.progress[c.id])) continue
    if (!(c.prerequisites ?? []).every(pid => isMastered(input.progress[pid]))) continue
    const w = weakById.get(c.id)
    const inProgress = input.progress[c.id] !== undefined || recent.has(c.id)
    out.push({
      kind: w && !inProgress ? 'weakness' : inProgress ? 'resume' : 'new',
      id: c.id,
      label: w?.label ?? c.label,
      weakness: w?.weakness ?? 0,
    })
  }
  out.sort((a, b) => b.weakness - a.weakness || (a.id < b.id ? -1 : 1))
  return out
}

// fixedWeightPolicy: frontier-new, due-review, weakness/remediation mix.
// Falls back through resume → diagnostic → browse so the shelf is always
// satisfiable. Cap 5. (Reviews collapse to one card: dueReviews is a count,
// not concept-level items — per-concept review slots await review detail.)
export const fixedWeightPolicy: SelectionPolicy = (cands, input) => {
  const items: ShelfItem[] = []
  const used = new Set<string>()
  const take = (it: ShelfItem, id?: string) => {
    if (items.length >= 5) return
    if (id && used.has(id)) return
    if (id) used.add(id)
    items.push(it)
  }
  const news = cands.filter(c => c.kind === 'new')
  const weaks = cands.filter(c => c.kind === 'weakness')
  const resumes = cands.filter(c => c.kind === 'resume')

  if (input.conceptsMastered === 0 && !input.diagnosticCompleted) {
    take({ kind: 'diagnostic', badge: 'Recommended', title: 'Find your frontier', detail: 'Adaptive diagnostic — where to start', href: '/onboard', cta: 'Start →', xp: 0 })
  }
  if (input.dueReviews > 0) {
    take({
      kind: 'review', badge: 'Due now',
      title: `${input.dueReviews} concept${input.dueReviews !== 1 ? 's' : ''} due for review`,
      detail: 'Spaced repetition — review before decay', href: '/review', cta: 'Review now →', xp: 5,
    })
  }
  for (const n of news.slice(0, 2)) {
    take(learnItem('new', 'New', n.id, n.label, 'Frontier concept — learn it next', 'Start learning →'), n.id)
  }
  for (const w of weaks.slice(0, 1)) {
    take(learnItem('weakness', 'Recommended', w.id, w.label, 'Weakest eligible concept — targeted session', 'Practice →'), w.id)
  }
  for (const r of resumes.slice(0, 1)) {
    if (items.length >= 5) break
    take(learnItem('resume', 'Continue', r.id, r.label, 'Pick up where you left off', 'Continue →'), r.id)
  }
  if (items.length === 0) {
    take({ kind: 'browse', badge: 'Study', title: 'Browse the Study library', detail: 'Pick a Lesson — worked example first, then answer', href: '/study', cta: 'Browse Study →', xp: 0 })
  }
  return items.slice(0, 5)
}

// selectShelf is the primary task surface. Policy is injectable; v1 ships
// the fixed-weight policy.
export function selectShelf(input: ShelfInput, policy: SelectionPolicy = fixedWeightPolicy): ShelfItem[] {
  return policy(buildCandidates(input), input)
}
