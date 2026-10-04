import { describe, it, expect } from 'vitest'
import { selectShelf, selectShelfHead, buildCandidates, upcomingLocked, recentlyUnlocked, hrefConceptId, type ShelfInput } from '../lib/nextUp'

const catalog = [
  { id: 'a', label: 'A', prerequisites: [] as string[], avgTimeSeconds: 10 },
  { id: 'b', label: 'B', prerequisites: ['a'], avgTimeSeconds: 120 },
  { id: 'c', label: 'C', prerequisites: ['a'], avgTimeSeconds: 10 },
  { id: 'd.word', label: 'D word', prerequisites: ['b', 'c'], avgTimeSeconds: 30 },
  { id: 'e', label: 'E', prerequisites: ['z-locked'], avgTimeSeconds: 10 },
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

  it('does not offer a decayed concept for teaching', () => {
    // The acceptance criterion, at the level of the queue that feeds /learn.
    //
    // A decayed concept used to appear here as a `resume` item, because the only status
    // filter asked `status === 'MASTERED'` and the server reports stale mastery as
    // DECAYING. /learn then opened with the concept's worked example, unconditionally —
    // it has no notion of the lesson having been seen before — so a fortnight's break
    // meant re-reading the tutorial for something already demonstrated. That is the
    // purgatory, and this is where it entered.
    const cands = buildCandidates({ ...base, progress: { b: { status: 'DECAYING', streak: 3 } } })
    expect(cands.map(c => c.id)).not.toContain('b')

    // And a decayed prerequisite still unlocks its dependents, matching
    // scheduler.prereqsMet, so removing it from the queue cannot strand the frontier.
    const unlocked = buildCandidates({ ...base, progress: { a: { status: 'DECAYING', streak: 3 } } })
    expect(unlocked.map(c => c.id)).toContain('b')
    expect(unlocked.map(c => c.id)).toContain('c')
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

  it('prices items from effort calibration', () => {
    const items = selectShelf({
      ...base,
      progress: {
        a: { status: 'MASTERED', streak: 3 },
        b: { status: 'MASTERED', streak: 3 },
        c: { status: 'MASTERED', streak: 3 },
      },
    })
    // d.word (30s) calibrates to 1 XP; only it is eligible.
    expect(items).toHaveLength(1)
    expect(items[0].href).toContain('d.word')
    expect(items[0].xp).toBe(1)
  })

  it('prices a 120s concept at 2 XP', () => {
    const items = selectShelf({
      ...base,
      progress: { a: { status: 'MASTERED', streak: 3 } },
    })
    const hard = items.find(i => i.href.includes('concept=b'))
    expect(hard?.xp).toBe(2)
  })

  it('is always satisfiable for brand-new users', () => {    const items = selectShelf({
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

  it('omits diagnostic for active-but-unmastered learners', () => {
    const items = selectShelf({
      ...base,
      activity: [{ date: '2026-09-28', questions: 53, correct: 44, concepts: ['arith.add'] }],
    })
    expect(items.length).toBeGreaterThan(0)
    expect(items.some(i => i.kind === 'diagnostic')).toBe(false)
  })

  it('drops excluded concept ids from candidates', () => {
    const input = {
      ...base,
      progress: { a: { status: 'MASTERED', streak: 3 }, b: { status: 'learning', streak: 1 } },
    }
    expect(buildCandidates(input).map(c => c.id)).toContain('b')
    expect(buildCandidates({ ...input, excludeConceptIds: ['b'] }).map(c => c.id)).not.toContain('b')
  })

  it('keeps a valid resume when nothing is excluded (entry/profile behavior)', () => {
    const input = {
      ...base,
      progress: { a: { status: 'MASTERED', streak: 3 }, b: { status: 'learning', streak: 1 } },
    }
    const cands = buildCandidates(input)
    expect(cands.find(c => c.id === 'b')?.kind).toBe('resume')
    const items = selectShelf(input)
    expect(items.some(i => i.href.includes('concept=b'))).toBe(true)
  })

  it('never heads an excluded concept and falls back to browse when empty', () => {
    const head = selectShelfHead({
      ...base,
      progress: { a: { status: 'MASTERED', streak: 3 }, b: { status: 'learning', streak: 1 } },
      excludeConceptIds: ['b', 'c'],
    })
    expect(head.next.kind).toBe('browse')
    expect([head.next, ...head.alternatives].every(i => !i.href.includes('concept=b') && !i.href.includes('concept=c'))).toBe(true)
  })
})

describe('upcomingLocked', () => {
  it('lists locked successors with missing prerequisite labels', () => {
    const locked = upcomingLocked(catalog, { a: { status: 'MASTERED', streak: 3 } }, 'b')
    // d.word needs b (done here: unmastered) and c (unmastered) → locked.
    expect(locked.map(l => l.id)).toEqual(['d.word'])
    expect(locked[0].missing).toEqual(['B', 'C'])
  })

  it('omits mastered and fully-eligible successors', () => {
    const locked = upcomingLocked(
      catalog,
      { a: { status: 'MASTERED', streak: 3 }, b: { status: 'MASTERED', streak: 3 }, c: { status: 'MASTERED', streak: 3 } },
      'b'
    )
    // d.word is fully eligible now (already a `new` candidate) → not listed.
    expect(locked).toHaveLength(0)
  })

  it('ignores non-successors', () => {
    expect(upcomingLocked(catalog, {}, 'c').map(l => l.id)).not.toContain('b')
  })
})

describe('completion unlock rule (PR1 client half)', () => {
  it('treats a completed prerequisite as satisfied', () => {
    const cands = buildCandidates({
      ...base,
      progress: {
        a: { status: 'MASTERED', streak: 3 },
        b: { status: 'learning', streak: 1, completed: true },
        c: { status: 'MASTERED', streak: 3 },
      },
    })
    // d.word needs b + c: b is completed-not-mastered, c mastered → eligible.
    expect(cands.map(c => c.id)).toContain('d.word')
  })

  it('keeps locked successors locked without completion', () => {
    const cands = buildCandidates({
      ...base,
      progress: {
        a: { status: 'MASTERED', streak: 3 },
        b: { status: 'learning', streak: 1 },
        c: { status: 'MASTERED', streak: 3 },
      },
    })
    expect(cands.map(c => c.id)).not.toContain('d.word')
  })
})

describe('hrefConceptId', () => {
  it('parses the concept id regardless of extra params', () => {
    expect(hrefConceptId('/learn?concept=a.b&seed=42&difficulty=0.7')).toBe('a.b')
    expect(hrefConceptId('/learn?concept=a.b')).toBe('a.b')
    expect(hrefConceptId('/review')).toBeNull()
    expect(hrefConceptId('/study')).toBeNull()
  })
})

describe('recentlyUnlocked', () => {
  const today = new Date().toISOString().slice(0, 10)
  const old = new Date(Date.now() - 30 * 86400000).toISOString().slice(0, 10)

  it('lists successors unlocked by recently-active concepts', () => {
    const rows = recentlyUnlocked({
      catalog,
      progress: {
        a: { status: 'MASTERED', streak: 3 },
        b: { status: 'learning', streak: 2, completed: true },
        c: { status: 'MASTERED', streak: 3 },
      },
      activity: [{ date: today, questions: 5, correct: 4, concepts: ['b'] }],
      recentDays: 7,
    })
    expect(rows.map(r => r.id)).toEqual(['d.word'])
    expect(rows[0].via).toBe('B')
  })

  it('ignores stale activity and mastered successors', () => {
    const rows = recentlyUnlocked({
      catalog,
      progress: {
        a: { status: 'MASTERED', streak: 3 },
        b: { status: 'learning', streak: 2, completed: true },
        c: { status: 'MASTERED', streak: 3 },
        'd.word': { status: 'MASTERED', streak: 10 },
      },
      activity: [{ date: old, questions: 5, correct: 4, concepts: ['b'] }],
      recentDays: 7,
    })
    expect(rows).toHaveLength(0)
  })

  it('omits successors whose prerequisites are still unmet', () => {
    const rows = recentlyUnlocked({
      catalog,
      progress: { a: { status: 'MASTERED', streak: 3 } },
      activity: [{ date: today, questions: 5, correct: 4, concepts: ['b'] }],
      recentDays: 7,
    })
    // d.word needs c too, which is untouched → still locked.
    expect(rows).toHaveLength(0)
  })
})
