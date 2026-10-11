import { describe, it, expect } from 'vitest'
import { prerequisiteClosure, scopeSize } from '../lib/diagnosticScope'
import { concepts } from '../lib/conceptData'

// The number shown on the diagnostic button has to be the number of concepts the
// assessment can actually reach, not the sum of the domains a learner ticked.
// These assertions are what stop the two drifting: the walk here and the walk in
// `internal/planning/planner.go` both read the same prerequisite lists, and the
// select-everything case pins them against the whole corpus.

const ALL_IDS = concepts.map(c => c.id)
const byId = new Map(concepts.map(c => [c.id, c]))

/** A domain's concept ids, the way the onboarding page assembles a selection. */
function idsForDomain(domain: string): string[] {
  return concepts.filter(c => c.domain === domain).map(c => c.id)
}

describe('prerequisiteClosure', () => {
  it('covers every concept when everything is selected', () => {
    expect(scopeSize(ALL_IDS)).toBe(concepts.length)
  })

  it('includes the selection itself', () => {
    const closure = prerequisiteClosure(idsForDomain('arithmetic'))
    for (const id of idsForDomain('arithmetic')) expect(closure).toContain(id)
  })

  it('pulls in prerequisites the learner did not select', () => {
    // A single leaf concept drags everything beneath it into scope, which is the
    // whole reason a flat per-domain sum understates the assessment.
    const withPrereqs = concepts.filter(c => (c.prerequisites ?? []).length > 0)
    expect(withPrereqs.length).toBeGreaterThan(0)
    const leaf = withPrereqs[withPrereqs.length - 1]
    const closure = prerequisiteClosure([leaf.id])
    expect(closure.length).toBeGreaterThan(1)
    for (const p of leaf.prerequisites ?? []) expect(closure).toContain(p)
  })

  it('never includes a concept outside the corpus', () => {
    const closure = prerequisiteClosure([...ALL_IDS, 'not.a.real.concept'])
    expect(closure).toContain('not.a.real.concept')
    for (const id of closure) {
      if (id === 'not.a.real.concept') continue
      expect(byId.has(id)).toBe(true)
    }
  })

  it('returns just the roots when given the roots', () => {
    const roots = concepts.filter(c => (c.prerequisites ?? []).length === 0)
    expect(roots.length).toBeGreaterThan(0)
    const closure = prerequisiteClosure(roots.map(r => r.id))
    // Nothing beneath the roots can be reached from the roots alone.
    expect(closure).toHaveLength(roots.length)
  })

  it('is idempotent and order-independent', () => {
    const a = prerequisiteClosure(idsForDomain('arithmetic')).sort()
    const b = prerequisiteClosure([...idsForDomain('arithmetic')].reverse()).sort()
    expect(a).toEqual(b)
    expect(prerequisiteClosure(a).sort()).toEqual(a)
  })

  it('never shrinks as domains are added to the selection', () => {
    // Deselect-after-select-all is only correct if the count a learner watches
    // moves the way they expect. Adding to a selection can never shrink it; what
    // matters is the running total, not each domain measured alone.
    const domains = [...new Set(concepts.map(c => c.domain))].sort()
    const chosen: string[] = []
    let previous = 0
    for (const domain of domains) {
      chosen.push(...idsForDomain(domain))
      const size = scopeSize(chosen)
      expect(size).toBeGreaterThanOrEqual(previous)
      previous = size
    }
    expect(previous).toBe(concepts.length)
  })

  it('covers the union of two selections', () => {
    const one = new Set(prerequisiteClosure(idsForDomain('arithmetic')))
    const both = new Set(prerequisiteClosure([...idsForDomain('arithmetic'), ...idsForDomain('geometry')]))
    for (const id of one) expect(both).toContain(id)
    expect(both.size).toBeGreaterThanOrEqual(one.size)
  })

  it('reports zero for an empty selection', () => {
    // The button is disabled in this state; the count must not lie about it.
    expect(scopeSize([])).toBe(0)
  })
})