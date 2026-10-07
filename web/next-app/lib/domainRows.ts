import { countByDomain, masteryState, prereqMet } from './progress'
import { DOMAIN_ORDER, domainLabel } from './graphDomains'
import type { ConceptRecord } from './conceptData'
import type { ConceptProgress } from './api'

// Per-domain rows, derived once and used by both /domains and the read-only /share report.
//
// The counting was `countByDomain`'s job already and stays there; what lives here is the
// per-concept detail and the "what should I do about this domain" answer, because both consumers
// needed it and /domains had it inline while /share wanted the same thing.
//
// Domain order and labels come from lib/graphDomains. Six modules used to carry their own copy
// of that list and three of them omitted `machine_learning` while DomainProgress omitted it and
// `precalculus` too — 50 concepts that no surface showed, with nothing on screen to say so.
// test/domainAuthority.test.ts fails if a partial list appears again.

export type DomainConceptStatus = 'mastered' | 'dueForReview' | 'inProgress' | 'unseen'

export interface DomainConcept {
  id: string
  label: string
  status: DomainConceptStatus
}

export interface DomainRow {
  domain: string
  label: string
  total: number
  mastered: number
  dueForReview: number
  inProgress: number
  /** Next concept to attempt. Prefers one already in progress or due over one never started. */
  nextId: string | null
  nextLabel: string | null
  concepts: DomainConcept[]
}

export const DOMAIN_STATUS_LABEL: Record<DomainConceptStatus, string> = {
  mastered: 'Mastered',
  dueForReview: 'Due for review',
  inProgress: 'In progress',
  unseen: 'Not started',
}

export const DOMAIN_STATUS_ORDER: DomainConceptStatus[] = [
  'mastered',
  'dueForReview',
  'inProgress',
  'unseen',
]

function statusOf(p: ConceptProgress | undefined, unlocked: boolean): DomainConceptStatus {
  if (!unlocked) return 'unseen'
  switch (masteryState(p)) {
    case 'mastered':
      return 'mastered'
    case 'decaying':
      return 'dueForReview'
    case 'practicing':
    case 'learning':
      return 'inProgress'
    default:
      return 'unseen'
  }
}

/**
 * Build one row per domain, in DOMAIN_ORDER.
 *
 * `catalogue` is the full concept list; domains with no concepts are skipped rather than shown
 * as empty rows, so a domain that does not exist cannot appear as a 0% bar.
 */
export function buildDomainRows(
  catalogue: readonly ConceptRecord[],
  progress: Record<string, ConceptProgress>,
): DomainRow[] {
  // ConceptRecord already satisfies CatalogueConcept structurally (id, domain, prerequisites),
  // so this needs no cast at all.
  const counts = countByDomain(catalogue, progress)
  const byDomain = new Map<string, ConceptRecord[]>()
  for (const c of catalogue) {
    const list = byDomain.get(c.domain)
    if (list) list.push(c)
    else byDomain.set(c.domain, [c])
  }

  const rows: DomainRow[] = []
  for (const domain of DOMAIN_ORDER) {
    const count = counts.get(domain)
    if (!count) continue
    const members = byDomain.get(domain) ?? []
    const concepts: DomainConcept[] = members.map((c) => {
      const unlocked = (c.prerequisites ?? []).every((pid) => prereqMet(progress, pid))
      return { id: c.id, label: c.label, status: statusOf(progress[c.id], unlocked) }
    })
    const resumable =
      concepts.find((c) => c.status === 'inProgress' || c.status === 'dueForReview') ??
      concepts.find((c) => c.status === 'unseen')
    rows.push({
      domain,
      label: domainLabel(domain),
      total: count.total,
      mastered: count.mastered,
      dueForReview: count.dueForReview,
      inProgress: count.learning,
      nextId: resumable?.id ?? null,
      nextLabel: resumable?.label ?? null,
      concepts,
    })
  }
  return rows
}

/** Weakest first: least demonstrated work, then most review debt. Ties break on the label. */
export function byNeedsAttention(a: DomainRow, b: DomainRow): number {
  return (
    a.mastered / a.total - b.mastered / b.total ||
    b.dueForReview - a.dueForReview ||
    a.label.localeCompare(b.label)
  )
}
