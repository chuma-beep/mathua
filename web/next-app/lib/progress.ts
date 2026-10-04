import type { ConceptProgress } from './api'

// Shared mastery counting. Three copies of these two predicates existed — in
// DomainProgress and twice in nextUp — and a fourth was about to arrive with
// the progress page's position block. They agree today; they would drift the
// moment one was edited, and a drift here is invisible: Profile would count a
// concept as mastered while the scheduler scheduled it for review.
//
// A fifth arrived anyway, in lib/graphStatus.ts, and it had already drifted by
// the time this was read: it knew MASTERED, PRACTICING and LEARNING but not
// DECAYING, so a concept the learner had mastered and not reviewed for a
// fortnight rendered as untouched *and relocked everything downstream of it*.
// The paragraph above predicted exactly that failure. This file is now the
// only place the rule lives.

// The shape the predicates need: anything with a status. The API record, the
// graph's looser input and a test fixture all satisfy it, so no caller widens
// its own type just to ask a question about mastery.
type MaybeProgress = { status?: string | undefined } | undefined | null

/**
 * The learner's-facing mastery state.
 *
 * `decaying` is deliberately a state of its own and must stay one. It means
 * *previously mastered, retrieval confidence may have weakened* — not *never
 * learned*. Decay is computed server-side at read time and never persisted, so a
 * concept only ever arrives as `DECAYING`, never as a stored status.
 */
export type MasteryState = 'mastered' | 'decaying' | 'practicing' | 'learning' | 'unseen'

// isMastered reads the status case-insensitively because the stored value is
// uppercase on the server and older payloads are not.
//
// Strictly "is MASTERED right now", which is narrower than "has this been
// learned" and is the wrong question for anything about prerequisites or the
// learning loop. The server reports a concept whose review has gone stale as
// DECAYING, so asking this of a prerequisite quietly un-masters it — that is
// exactly how a fortnight's gap used to relock a whole downstream path. Use
// countsAsMastered where "learned" is the question.
export function isMastered(p: MaybeProgress): boolean {
  return (p?.status ?? '').toLowerCase() === 'mastered'
}

// countsAsMastered is mastery for the purpose of unlocking successors and for
// leaving the active learning loop. DECAYING qualifies, because the server says
// it does: scheduler.prereqsMet rejects only what is neither MASTERED nor
// DECAYING. Treating a decayed prerequisite as unmet would disagree with the
// scheduler that chose to serve the learner the dependent concept in the first
// place.
//
// Exported because it answers "has this been learned", and there is no longer a
// private copy of it that can drift.
export function countsAsMastered(p: MaybeProgress): boolean {
  const status = (p?.status ?? '').toUpperCase()
  return status === 'MASTERED' || status === 'DECAYING'
}

// masteryState is the display state: the stored status with DECAYING preserved
// as itself, so a colour or label can distinguish "learned and fresh" from
// "learned, due for review". Anything unrecognised reads as `unseen`, which is
// what a missing or unusable progress row means everywhere else.
export function masteryState(p: MaybeProgress): MasteryState {
  switch ((p?.status ?? '').toUpperCase()) {
    case 'MASTERED':
      return 'mastered'
    case 'DECAYING':
      return 'decaying'
    case 'PRACTICING':
      return 'practicing'
    case 'LEARNING':
      return 'learning'
    default:
      return 'unseen'
  }
}

// prereqMet: a prerequisite is satisfied when it counts as mastered OR
// completed. The lesson-completion endpoint sets `completed`, which unlocks
// successors, while mastery still gates the quiz. Same rule as the server's
// Available(); an absent `completed` (pre-PR1 payload) behaves exactly as
// before.
export function prereqMet(progress: Record<string, ConceptProgress>, pid: string): boolean {
  const p = progress[pid]
  return countsAsMastered(p) || p?.completed === true
}

export interface CatalogueConcept {
  id: string
  domain: string
  prerequisites?: string[]
}

export interface DomainCounts {
  domain: string
  total: number
  /**
   * Concepts whose competence has been demonstrated and not yet gone stale —
   * that is, MASTERED **or** DECAYING.
   *
   * It used to mean strictly `MASTERED`, which quietly subtracted decayed concepts and
   * so disagreed with the graph: a concept with a full mastery bar there was not counted
   * as mastered here, with nothing on screen to explain the difference. Decay is now its
   * own number, `dueForReview`, so the over-claiming this was avoiding is prevented by
   * *showing* the review debt instead of by shrinking the word "mastered".
   */
  mastered: number
  /** Previously mastered, past the decay window: retrieval confidence may have weakened. */
  dueForReview: number
  learning: number
  completed: number
  /** Prerequisites not yet satisfied, so the concept cannot be started. */
  locked: number
}

export interface OverallCounts extends DomainCounts {
  domain: '*'
  /**
   * Concepts whose prerequisites are all met — reachable now or soon. The
   * complement of locked, and more meaningful than a percentage of the whole
   * catalogue: 657 concepts spanning topology and abstract algebra is not a
   * denominator a learner can act on.
   */
  unlocked: number
}

// countByDomain buckets the catalogue by domain. Order is the caller's: a
// presentation order is a view concern, not a counting one.
export function countByDomain(
  catalogue: readonly CatalogueConcept[],
  progress: Record<string, ConceptProgress>,
): Map<string, Omit<DomainCounts, 'domain'>> {
  const byDomain = new Map<string, Omit<DomainCounts, 'domain'>>()
  for (const c of catalogue) {
    const entry = byDomain.get(c.domain) ?? { total: 0, mastered: 0, dueForReview: 0, learning: 0, completed: 0, locked: 0 }
    entry.total++
    const p = progress[c.id]
    if (masteryState(p) === 'decaying') entry.dueForReview++
    if (countsAsMastered(p)) {
      entry.mastered++
    } else if (p?.completed === true) {
      entry.completed++
    } else {
      if (p?.status === 'LEARNING' || p?.status === 'PRACTICING') entry.learning++
      if (!(c.prerequisites ?? []).every(pid => prereqMet(progress, pid))) entry.locked++
    }
    byDomain.set(c.domain, entry)
  }
  return byDomain
}

// countOverall is the headline "where am I" number: total in the catalogue,
// how many are mastered, and how many are unlocked.
export function countOverall(
  catalogue: readonly CatalogueConcept[],
  progress: Record<string, ConceptProgress>,
): OverallCounts {
  const totals = { total: 0, mastered: 0, dueForReview: 0, learning: 0, completed: 0, locked: 0 }
  for (const c of catalogue) {
    totals.total++
    const p = progress[c.id]
    if (masteryState(p) === 'decaying') totals.dueForReview++
    if (countsAsMastered(p)) {
      totals.mastered++
    } else if (p?.completed === true) {
      totals.completed++
    } else {
      if (p?.status === 'LEARNING' || p?.status === 'PRACTICING') totals.learning++
      if (!(c.prerequisites ?? []).every(pid => prereqMet(progress, pid))) totals.locked++
    }
  }
  return { domain: '*', ...totals, unlocked: totals.total - totals.locked }
}