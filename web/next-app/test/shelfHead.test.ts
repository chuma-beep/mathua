import { describe, it, expect } from 'vitest'
import { selectShelf, selectShelfHead, type ShelfInput } from '../lib/nextUp'

const catalog = [
  { id: 'a', label: 'A', prerequisites: [] as string[], avgTimeSeconds: 10 },
  { id: 'b', label: 'B', prerequisites: ['a'], avgTimeSeconds: 120 },
  { id: 'c', label: 'C', prerequisites: ['a'], avgTimeSeconds: 10 },
  { id: 'd.word', label: 'D word', prerequisites: ['b', 'c'], avgTimeSeconds: 30 },
]

const base: ShelfInput = {
  dueReviews: 0,
  weaknesses: { by_domain: {} },
  progress: {},
  activity: [],
  diagnosticCompleted: true,
  conceptsMastered: 3,
  catalog,
}

describe('selectShelfHead', () => {
  it('exposes item 0 as next with the rest as alternatives', () => {
    const flat = selectShelf({
      ...base,
      dueReviews: 4,
      progress: { a: { status: 'MASTERED', streak: 3 } },
      weaknesses: { by_domain: { d: [{ id: 'c', label: 'C', weakness: 0.9 }] } },
    })
    const head = selectShelfHead({
      ...base,
      dueReviews: 4,
      progress: { a: { status: 'MASTERED', streak: 3 } },
      weaknesses: { by_domain: { d: [{ id: 'c', label: 'C', weakness: 0.9 }] } },
    })
    expect(head.next).toEqual(flat[0])
    expect(head.alternatives).toEqual(flat.slice(1))
    expect([head.next, ...head.alternatives]).toHaveLength(flat.length)
  })

  it('keeps selectShelf as a flat wrapper over the head', () => {
    const input = { ...base, progress: { a: { status: 'MASTERED', streak: 3 } } }
    const head = selectShelfHead(input)
    expect(selectShelf(input)).toEqual([head.next, ...head.alternatives])
  })

  it('is always satisfiable for brand-new users with diagnostic head', () => {
    const head = selectShelfHead({
      dueReviews: 0,
      weaknesses: { by_domain: {} },
      progress: {},
      activity: [],
      diagnosticCompleted: false,
      conceptsMastered: 0,
      catalog,
    })
    expect(head.next.kind).toBe('diagnostic')
    // 'a' is the only eligible frontier concept, so it follows as alternative.
    expect(head.alternatives).toHaveLength(1)
    expect(head.alternatives[0].kind).toBe('new')
  })

  it('falls back to browse head when nothing is eligible', () => {
    const head = selectShelfHead({
      ...base,
      catalog: [{ id: 'z', label: 'Z', prerequisites: ['missing'] }],
    })
    expect(head.next.kind).toBe('browse')
    expect(head.alternatives).toHaveLength(0)
  })

  it('folds an empty policy into the selectNextUp fallback head', () => {
    const head = selectShelfHead({ ...base }, () => [])
    expect(head.next.kind).toBe('browse')
    expect(head.next.href).toBe('/study')
    expect(head.alternatives).toHaveLength(0)
  })

  it('selectNextUp fallback skips excluded ids in weakness and resume picks', () => {
    const input = {
      ...base,
      progress: { a: { status: 'MASTERED', streak: 3 }, b: { status: 'learning', streak: 1 } },
      weaknesses: { by_domain: { d: [{ id: 'b', label: 'B', weakness: 0.9 }] } },
      activity: [{ date: '2026-09-30', questions: 5, correct: 4, concepts: ['b'] }],
    }
    const headed = selectShelfHead(input, () => [])
    expect(headed.next.kind).toBe('weakness')
    expect(headed.next.href).toContain('concept=b')
    const skipped = selectShelfHead({ ...input, excludeConceptIds: ['b'] }, () => [])
    expect(skipped.next.kind).toBe('browse')
    expect(skipped.next.href).not.toContain('concept=b')
  })
})

// A learner who has finished a concept and come back after a break must not find it
// waiting for them as teaching. The shelf is what `/learn` and the Profile "Next task"
// both read, so this is the single place the rule has to hold.
describe('mastery exits the active teaching queue', () => {
  const catalog = [
    { id: 'a', label: 'A', prerequisites: [] as string[], avgTimeSeconds: 10 },
    { id: 'b', label: 'B', prerequisites: ['a'] as string[], avgTimeSeconds: 10 },
  ]
  const input = {
    dueReviews: 0,
    weaknesses: { by_domain: {} },
    activity: [],
    diagnosticCompleted: true,
    conceptsMastered: 1,
    catalog,
  }
  const all = (h: { next: { href: string }; alternatives: { href: string }[] }) => [h.next, ...h.alternatives]

  it('offers a fresh MASTERED concept for nothing', () => {
    const head = selectShelfHead({ ...input, progress: { a: { status: 'MASTERED', streak: 3 } } })
    expect(JSON.stringify(all(head))).not.toContain('concept=a')
  })

  it('offers a DECAYING concept for nothing, and points at review instead', () => {
    // The whole point: decay changes what is *maintenance*, not what is *teaching*.
    const withReview = {
      ...input,
      dueReviews: 2,
      progress: { a: { status: 'DECAYING', streak: 3 } },
    }
    const head = selectShelfHead(withReview)

    // Nothing routes the learner back into the lesson for `a`.
    expect(JSON.stringify(all(head))).not.toContain('concept=a')
    // And the maintenance route is offered, because "get out of the way" must not mean
    // "silently drop it": the retrieval check is still owed.
    const review = all(head).find(i => i.href === '/review')
    expect(review).toBeDefined()
  })

  it('keeps a concept that is still being learned', () => {
    const prog = { a: { status: 'DECAYING', streak: 3 }, b: { status: 'PRACTICING', streak: 1 } }
    const head = selectShelfHead({ ...input, progress: prog })
    expect(head.next.href).toContain('concept=b')
  })
})
