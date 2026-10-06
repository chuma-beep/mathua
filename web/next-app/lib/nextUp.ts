import type { ConceptProgress, DailyActivity } from './api'
import { countsAsMastered, prereqMet } from './progress'

/**
 * Display shapes for a recommendation, plus the graph helpers the surfaces use.
 *
 * **The decision lives in the engine now.** This module used to make it: `buildCandidates`
 * assembled the eligible concepts, `fixedWeightPolicy` ranked them, and `selectShelfHead`
 * chose the head — all in the browser, from five API calls plus a bundled copy of the corpus.
 * `internal/scheduler` holds the same rules for serving questions inside a session, and the two
 * had already drifted: the client kept its own copy of the XP award and still halved reviews,
 * a rule ADR-020 removed from the server, and it ranked by weakness-then-concept-id where the
 * server ranks by a weighted score.
 *
 * That code is deleted rather than left in place, because code that *looks* like authority is
 * worse than a visible duplicate. ADR-040 removed a second client mastery model for the same
 * reason: a model with no authority cannot be reconciled with the real one, so it only
 * accumulates plausible answers to nothing. `test/masteryAuthority.test.ts` now guards this file
 * the same way it guards that one.
 *
 * What stays is presentation shape (`Shelf`, `ShelfItem`) and helpers that report on the
 * learner's own history rather than deciding anything.
 */

export interface Shelf {
  next: ShelfItem
  alternatives: ShelfItem[]
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

// CatalogEntry is the shape the graph helpers read: id, label and prerequisites, with no
// ranking attached. `CatalogEntry` used to carry avgTimeSeconds as well, which only existed so
// the client could price XP for its own labels — a copy of internal/xp that had already drifted.
export interface CatalogEntry {
  id: string
  label: string
  prerequisites?: string[]
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
    if (countsAsMastered(progress[c.id])) continue
    const missing = (c.prerequisites ?? []).filter(pid => !prereqMet(progress, pid))
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
      if (seen.has(c.id) || countsAsMastered(input.progress[c.id])) continue
      if (!(c.prerequisites ?? []).every(pid => prereqMet(input.progress, pid))) continue
      seen.add(c.id)
      out.push({ id: c.id, label: c.label, via: labels.get(cid) ?? cid })
    }
  }
  out.sort((a, b) => (a.id < b.id ? -1 : 1))
  return out
}

export function isNewUser(input: { conceptsMastered: number; activity: DailyActivity[] }): boolean {
  if (input.conceptsMastered !== 0) return false
  return !input.activity.some(day => day.questions > 0)
}
