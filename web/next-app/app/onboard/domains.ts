import { DOMAIN_ORDER, domainLabel } from '../../lib/graphDomains'

export interface DomainInfo {
  name: string
  count: number
  concepts: string[]
  selected: boolean
}

// Re-exported from the canonical list rather than copied. This file's order turned out to be
// identical to DOMAIN_ORDER apart from omitting `machine_learning`, so it was a stale copy and
// not a deliberate teaching order. Four modules import these names from here.
//
// Widened to string[]: DOMAIN_ORDER is `as const`, so its element type is a literal union and
// `.indexOf(someString)` stops compiling for callers that treat this as an open list.
export const domainOrder: string[] = [...DOMAIN_ORDER]
export const domainLabels = domainLabel
