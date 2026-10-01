import type { ConceptProgress } from './api'

// Shared mastery counting. Three copies of these two predicates existed — in
// DomainProgress and twice in nextUp — and a fourth was about to arrive with
// the progress page's position block. They agree today; they would drift the
// moment one was edited, and a drift here is invisible: Profile would count a
// concept as mastered while the scheduler scheduled it for review.

// isMastered reads the status case-insensitively because the stored value is
// uppercase on the server and older payloads are not.
export function isMastered(p: ConceptProgress | undefined): boolean {
  return (p?.status ?? '').toLowerCase() === 'mastered'
}

// prereqMet: a prerequisite is satisfied when mastered OR completed. The
// lesson-completion endpoint sets `completed`, which unlocks successors, while
// mastery still gates the quiz. Same rule as the server's Available(); an
// absent `completed` (pre-PR1 payload) behaves exactly as before.
export function prereqMet(progress: Record<string, ConceptProgress>, pid: string): boolean {
  const p = progress[pid]
  return isMastered(p) || p?.completed === true
}

export interface CatalogueConcept {
  id: string
  domain: string
  prerequisites?: string[]
}

export interface DomainCounts {
  domain: string
  total: number
  mastered: number
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
    const entry = byDomain.get(c.domain) ?? { total: 0, mastered: 0, learning: 0, completed: 0, locked: 0 }
    entry.total++
    const p = progress[c.id]
    if (isMastered(p)) {
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
  const totals = { total: 0, mastered: 0, learning: 0, completed: 0, locked: 0 }
  for (const c of catalogue) {
    totals.total++
    const p = progress[c.id]
    if (isMastered(p)) {
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