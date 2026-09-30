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
  // Concept ids to skip when ranking. Named apart from the ?exclude=
  // question-variant URL param to avoid confusion. Passed only by the Learn
  // done re-fetch so a finished-but-unmastered concept never heads its own
  // shelf; entry/profile surfaces omit it and keep prior behavior.
  excludeConceptIds?: string[]
}

function isMastered(p: ConceptProgress | undefined): boolean {
  return (p?.status ?? '').toLowerCase() === 'mastered'
}

// A prerequisite is satisfied when mastered OR completed (lesson-completion
// unlocks successors; mastery still gates the quiz). Same rule as the
// server's Available(); absent `completed` (pre-PR1 payloads) behaves exactly
// as before.
function prereqSatisfied(progress: Record<string, ConceptProgress>, pid: string): boolean {
  const p = progress[pid]
  return isMastered(p) || p?.completed === true
}

// New means no evidence of learning at all: nothing mastered and no
// answered questions. A struggling newcomer with attempts but zero
// mastery still counts as new; anyone with a history does not, even if
// they never took the diagnostic.
export function isNewUser(input: { conceptsMastered: number; activity: DailyActivity[] }): boolean {
  if (input.conceptsMastered !== 0) return false
  return !input.activity.some(day => day.questions > 0)
}

function topWeakness(input: NextUpInput): { id: string; label: string } | null {
  const entries = input.weaknesses?.by_domain ?? {}
  const excluded = new Set(input.excludeConceptIds ?? [])
  let best: { id: string; label: string; weakness: number } | null = null
  for (const items of Object.values(entries)) {
    for (const it of items) {
      if (excluded.has(it.id)) continue
      if (isMastered(input.progress[it.id])) continue
      if (!best || it.weakness > best.weakness) best = it
    }
  }
  return best
}

function mostRecentInProgress(input: NextUpInput): string | null {
  const excluded = new Set(input.excludeConceptIds ?? [])
  for (let i = input.activity.length - 1; i >= 0; i--) {
    const day = input.activity[i]
    for (let j = day.concepts.length - 1; j >= 0; j--) {
      const cid = day.concepts[j]
      if (excluded.has(cid)) continue
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
 * selectNextUp is the single-directive fallback, folded into selectShelfHead:
 * when a policy returns no items the head falls back to this. The task shelf
 * (selectShelfHead) is the primary surface: the algorithm proposes, the
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
      cta: 'Review →',
    }
  }

  const weak = topWeakness(input)
  if (weak) {
    return {
      kind: 'weakness',
      badge: 'Recommended',
      title: `Study ${weak.label}`,
      detail: 'Weakest concept — targeted session',
      href: `/learn?concept=${encodeURIComponent(weak.id)}`,
      cta: 'Continue →',
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

  if (isNewUser(input) && !input.diagnosticCompleted) {
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

// ── Ranked shelf: item 0 is the Next head, the rest are alternatives ──
// The DAG + scheduler stay the source of truth (eligible set, interleave,
// 70/30 review/new). The UI commits to one head item so the learner sees
// "next", not "eligible". Learner still disposes via the alternatives.

export interface Shelf {
  next: ShelfItem
  alternatives: ShelfItem[]
}

function shelfItemFromNextUp(n: NextUp): ShelfItem {
  return { kind: n.kind, badge: n.badge, title: n.title, detail: n.detail, href: n.href, cta: n.cta, xp: 0 }
}

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
  avgTimeSeconds?: number
}

export interface ShelfInput extends NextUpInput {
  catalog: CatalogEntry[]
}

export interface ShelfCandidate {
  kind: 'new' | 'weakness' | 'resume'
  id: string
  label: string
  weakness: number
  avgTimeSeconds?: number
}

// SelectionPolicy turns eligible candidates into a shelf. The v1 fixed
// policy below is a starting weight set; later policies score candidates on
// readiness, urgency, value, difficulty, diversity, pace, and relevance.
export type SelectionPolicy = (cands: ShelfCandidate[], input: ShelfInput) => ShelfItem[]

// effortBase mirrors internal/xp EffortBase at neutral difficulty so shelf
// labels track the award path (10s→1, 60s→1, 120s→2, 300s→5).
function effortBase(avgTimeSeconds?: number): number {
  const t = avgTimeSeconds && avgTimeSeconds > 0 ? avgTimeSeconds : 10
  return Math.min(5, Math.max(1, Math.round((t / 60) * 1.0)))
}

function xpFor(entry: { id: string; avgTimeSeconds?: number }, kind: ShelfItem['kind']): number {
  if (kind === 'review') return Math.max(1, Math.round(effortBase(entry.avgTimeSeconds) / 2))
  if (kind === 'diagnostic' || kind === 'browse') return 0
  return effortBase(entry.avgTimeSeconds)
}

function learnItem(kind: ShelfItem['kind'], badge: string, c: ShelfCandidate, detail: string, cta: string): ShelfItem {
  return { kind, badge, title: c.label, detail, href: `/learn?concept=${encodeURIComponent(c.id)}`, cta, xp: xpFor(c, kind) }
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
  const excluded = new Set(input.excludeConceptIds ?? [])
  for (const c of input.catalog) {
    if (excluded.has(c.id)) continue
    if (isMastered(input.progress[c.id])) continue
    if (!(c.prerequisites ?? []).every(pid => prereqSatisfied(input.progress, pid))) continue
    const w = weakById.get(c.id)
    const inProgress = input.progress[c.id] !== undefined || recent.has(c.id)
    out.push({
      kind: w && !inProgress ? 'weakness' : inProgress ? 'resume' : 'new',
      id: c.id,
      label: w?.label ?? c.label,
      weakness: w?.weakness ?? 0,
      avgTimeSeconds: c.avgTimeSeconds,
    })
  }
  out.sort((a, b) => b.weakness - a.weakness || (a.id < b.id ? -1 : 1))
  return out
}

export interface LockedSuccessor {
  id: string
  label: string
  missing: string[]
}

// upcomingLocked lists direct DAG successors of a concept that are still
// locked (at least one prerequisite unmastered), each with its missing
// prerequisite labels resolved from the catalog. Fully-eligible successors
// already surface as `new` shelf items and are not duplicated here. Used by
// the Learn done card so the learner sees the forward path plus its
// prerequisites even before this concept is mastered.
export function upcomingLocked(
  catalog: CatalogEntry[],
  progress: Record<string, ConceptProgress>,
  conceptId: string
): LockedSuccessor[] {
  const labels = new Map(catalog.map(c => [c.id, c.label]))
  const out: LockedSuccessor[] = []
  for (const c of catalog) {
    if (!(c.prerequisites ?? []).includes(conceptId)) continue
    if (isMastered(progress[c.id])) continue
    const missing = (c.prerequisites ?? []).filter(pid => !prereqSatisfied(progress, pid))
    if (missing.length === 0) continue // eligible — already a `new` candidate
    out.push({ id: c.id, label: c.label, missing: missing.map(pid => labels.get(pid) ?? pid) })
  }
  out.sort((a, b) => (a.id < b.id ? -1 : 1))
  return out
}

// Concept id carried by a shelf href, if any (/learn?concept=X). Shared by
// the Learn done guard and the Profile dedupe filters so "same concept" is
// decided on ids, never on raw href strings (query params vary).
export function hrefConceptId(href: string): string | null {
  try {
    const c = new URL(href, 'http://localhost').searchParams.get('concept')
    return c && c.length > 0 ? c : null
  } catch {
    return null
  }
}

export interface RecentUnlock {
  id: string
  label: string
  via: string
}

// Lookback for "recently unlocked" rows on Profile. Named (not inlined) so
// tests and the Profile section share one value.
export const RECENT_UNLOCK_DAYS = 7

// recentlyUnlocked lists successors unlocked by recently-active concepts:
// for each concept touched in the last `recentDays` days, its successors
// whose prerequisites are now all satisfied (mastered or completed) and
// that remain unmastered. Recency is derived client-side from activity, so
// no endpoint is needed; the `completed` flag (PR1 payload) sharpens the
// unlock rule when present.
export function recentlyUnlocked(input: {
  catalog: CatalogEntry[]
  progress: Record<string, ConceptProgress>
  activity: DailyActivity[]
  recentDays: number
}): RecentUnlock[] {
  const cutoff = Date.now() - input.recentDays * 86400000
  const recent = new Set<string>()
  for (const day of input.activity) {
    const t = new Date(day.date).getTime()
    if (!Number.isFinite(t) || t < cutoff) continue
    for (const cid of day.concepts ?? []) recent.add(cid)
  }
  const labels = new Map(input.catalog.map(c => [c.id, c.label]))
  const out: RecentUnlock[] = []
  const seen = new Set<string>()
  for (const cid of recent) {
    for (const c of input.catalog) {
      if (!(c.prerequisites ?? []).includes(cid)) continue
      if (seen.has(c.id) || isMastered(input.progress[c.id])) continue
      if (!(c.prerequisites ?? []).every(pid => prereqSatisfied(input.progress, pid))) continue
      seen.add(c.id)
      out.push({ id: c.id, label: c.label, via: labels.get(cid) ?? cid })
    }
  }
  out.sort((a, b) => (a.id < b.id ? -1 : 1))
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

  if (isNewUser(input) && !input.diagnosticCompleted) {
    take({ kind: 'diagnostic', badge: 'Recommended', title: 'Find your frontier', detail: 'Adaptive diagnostic — where to start', href: '/onboard', cta: 'Start →', xp: 0 })
  }
  if (input.dueReviews > 0) {
    take({
      kind: 'review', badge: 'Due now',
      title: `${input.dueReviews} concept${input.dueReviews !== 1 ? 's' : ''} due for review`,
      detail: 'Spaced repetition — review before decay', href: '/review', cta: 'Review →', xp: 1,
    })
  }
  for (const n of news.slice(0, 2)) {
    take(learnItem('new', 'New', n, 'Frontier concept — learn it next', 'Continue →'), n.id)
  }
  for (const w of weaks.slice(0, 1)) {
    take(learnItem('weakness', 'Recommended', w, 'Weakest eligible concept — targeted session', 'Continue →'), w.id)
  }
  for (const r of resumes.slice(0, 1)) {
    if (items.length >= 5) break
    take(learnItem('resume', 'Continue', r, 'Pick up where you left off', 'Continue →'), r.id)
  }
  if (items.length === 0) {
    take({ kind: 'browse', badge: 'Study', title: 'Browse the Study library', detail: 'Pick a Lesson — worked example first, then answer', href: '/study', cta: 'Browse Study →', xp: 0 })
  }
  return items.slice(0, 5)
}

// selectShelfHead is the primary task surface: a ranked shelf whose item 0
// is the Next head. Policy is injectable; v1 ships the fixed-weight policy.
// Always satisfiable: an empty policy result falls back to selectNextUp
// (single-directive fallback folded in as the head, alternatives empty).
export function selectShelfHead(input: ShelfInput, policy: SelectionPolicy = fixedWeightPolicy): Shelf {
  const items = policy(buildCandidates(input), input).slice(0, 5)
  if (items.length === 0) {
    return { next: shelfItemFromNextUp(selectNextUp(input)), alternatives: [] }
  }
  const [next, ...alternatives] = items
  return { next, alternatives }
}

// selectShelf is a flat wrapper over selectShelfHead (head + alternatives).
// Kept for tests and any external callers. UI code should use
// selectShelfHead with NextUpCard.
export function selectShelf(input: ShelfInput, policy: SelectionPolicy = fixedWeightPolicy): ShelfItem[] {
  const head = selectShelfHead(input, policy)
  return [head.next, ...head.alternatives]
}
