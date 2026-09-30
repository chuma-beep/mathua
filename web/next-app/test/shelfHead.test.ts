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
})
