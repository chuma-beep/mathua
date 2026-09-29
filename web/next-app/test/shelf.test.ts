import { describe, it, expect } from 'vitest'
import { selectShelf, buildCandidates, type ShelfInput } from '../lib/nextUp'

const catalog = [
  { id: 'a', label: 'A', prerequisites: [] as string[] },
  { id: 'b', label: 'B', prerequisites: ['a'] },
  { id: 'c', label: 'C', prerequisites: ['a'] },
  { id: 'd.word', label: 'D word', prerequisites: ['b', 'c'] },
  { id: 'e', label: 'E', prerequisites: ['z-locked'] },
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

describe('buildCandidates', () => {
  it('only lists available unmastered concepts', () => {
    const cands = buildCandidates({ ...base, progress: { a: { status: 'MASTERED', streak: 3 } } })
    const ids = cands.map(c => c.id)
    expect(ids).toContain('b')
    expect(ids).toContain('c')
    expect(ids).not.toContain('a') // mastered
    expect(ids).not.toContain('d.word') // prereqs unmet
    expect(ids).not.toContain('e') // locked
  })

  it('classifies weakness vs resume vs new', () => {
    const cands = buildCandidates({
      ...base,
      progress: { a: { status: 'MASTERED', streak: 3 }, b: { status: 'learning', streak: 1 } },
      weaknesses: { by_domain: { d: [{ id: 'c', label: 'C', weakness: 0.9 }] } },
    })
    const byId = Object.fromEntries(cands.map(c => [c.id, c.kind]))
    expect(byId['b']).toBe('resume')
    expect(byId['c']).toBe('weakness')
  })
})

describe('selectShelf', () => {
  it('caps at 5 with mixed composition', () => {
    const items = selectShelf({
      ...base,
      dueReviews: 4,
      progress: { a: { status: 'MASTERED', streak: 3 } },
      weaknesses: { by_domain: { d: [{ id: 'c', label: 'C', weakness: 0.9 }] } },
    })
    expect(items.length).toBeLessThanOrEqual(5)
    const kinds = items.map(i => i.kind)
    expect(kinds).toContain('review')
    expect(kinds).toContain('new')
    expect(kinds).toContain('weakness')
    // No duplicate concepts.
    const hrefs = items.map(i => i.href)
    expect(new Set(hrefs).size).toBe(hrefs.length)
  })

  it('prices word items at multistep XP', () => {
    const items = selectShelf({
      ...base,
      progress: {
        a: { status: 'MASTERED', streak: 3 },
        b: { status: 'MASTERED', streak: 3 },
        c: { status: 'MASTERED', streak: 3 },
      },
    })
    const word = items.find(i => i.href.includes('d.word'))
    expect(word?.xp).toBe(15)
  })

  it('is always satisfiable for brand-new users', () => {
    const items = selectShelf({
      dueReviews: 0,
      weaknesses: { by_domain: {} },
      progress: {},
      activity: [],
      diagnosticCompleted: false,
      conceptsMastered: 0,
      catalog,
    })
    expect(items.length).toBeGreaterThan(0)
    expect(items[0].kind).toBe('diagnostic')
  })

  it('falls back to browse when nothing is eligible', () => {
    const items = selectShelf({
      ...base,
      catalog: [{ id: 'z', label: 'Z', prerequisites: ['missing'] }],
    })
    expect(items).toHaveLength(1)
    expect(items[0].kind).toBe('browse')
  })
})
