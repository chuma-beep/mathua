import { describe, it, expect } from 'vitest'
import { upcomingLocked, hrefConceptId, recentlyUnlocked, RECENT_UNLOCK_DAYS, isNewUser, type CatalogEntry } from '../lib/nextUp'
import type { ConceptProgress, DailyActivity } from '../lib/api'
import { readFileSync } from 'node:fs'
import { join } from 'node:path'

/**
 * The graph helpers in lib/nextUp.ts.
 *
 * The ranking tests that used to live here are gone because the ranking is gone: the decision
 * is made in internal/scheduler and arrives over /api/next, tested in Go
 * (`internal/scheduler/recommend_test.go`) and at the boundary (`TestNext_*` in
 * internal/server). What remains here is reporting on the learner's own history — which
 * successors a concept just unlocked, and what an href points at.
 *
 * The one behaviour worth noting is the decay rule, which had to be kept in step with the
 * scheduler by hand before (ADR-037): a DECAYING concept still satisfies a prerequisite, or a
 * concept the learner has not lost would lock everything downstream.
 */
const catalog: CatalogEntry[] = [
  { id: 'a', label: 'A', prerequisites: [] },
  { id: 'b', label: 'B', prerequisites: ['a'] },
  { id: 'c', label: 'C', prerequisites: ['a'] },
  { id: 'd.word', label: 'D word', prerequisites: ['b', 'c'] },
  { id: 'e', label: 'E', prerequisites: ['z-locked'] },
]
const cat = catalog
const prog = (s: Record<string, ConceptProgress['status']>): Record<string, ConceptProgress> =>
  Object.fromEntries(Object.entries(s).map(([k, status]) => [k, { status, streak: 0 }]))

const mastered = prog({ a: 'MASTERED' })
const decayed = prog({ a: 'DECAYING' })

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

/**
 * The `completed` prerequisite shortcut is inert.
 *
 * `prereqMet` treats `p?.completed === true` as satisfying a prerequisite, and `countByDomain`
 * and `countOverall` report a `completed` bucket on the strength of it. Two tests used to assert
 * that behaviour through `buildCandidates`. They were deleted with the rest of the ranking, and
 * replacing them matters more than keeping them: **no endpoint emits a `completed` field.**
 * `storage.ConceptProgress` has no such column, `ProgressView` embeds that struct, and no route
 * adds one, so the branch cannot fire against a real response.
 *
 * So the rule was tested but unreachable — the same failure shape as a probe run where the thing
 * does not exist. This test states it instead, so the next person to touch it learns the field
 * does not exist rather than re-deriving it. If a PR1 payload does land `completed`, delete this
 * and reinstate the rule's tests, because then it will finally be doing something.
 */
describe('the completed flag', () => {
  it('is not sent by any endpoint, so prereqMet\'s shortcut cannot fire', () => {
    const prereqMetSource = readFileSync(join(__dirname, '../lib/progress.ts'), 'utf8')
    expect(prereqMetSource).toContain('completed')

    // The server side is the fact that matters: grep it here so the test fails the moment a
    // payload starts sending the field and this note becomes stale.
    const store = readFileSync(join(__dirname, '../../../internal/storage/store.go'), 'utf8')
    const progressStruct = store.slice(store.indexOf('type ConceptProgress struct'))
    expect(progressStruct.slice(0, progressStruct.indexOf('}'))).not.toContain('Completed')
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


describe('isNewUser', () => {
  const day = (n: number): DailyActivity => ({ date: '2026-01-01', questions: n, correct: n, concepts: [] })

  it('is new only with no mastery and no answered questions', () => {
    expect(isNewUser({ conceptsMastered: 0, activity: [] })).toBe(true)
    expect(isNewUser({ conceptsMastered: 1, activity: [] })).toBe(false)
    expect(isNewUser({ conceptsMastered: 0, activity: [day(1)] })).toBe(false)
  })

  it('treats a struggling newcomer with attempts as not new', () => {
    // They have evidence about themselves even without mastery, so a diagnostic is not the
    // most useful thing to offer them.
    expect(isNewUser({ conceptsMastered: 0, activity: [day(4)] })).toBe(false)
  })
})

describe('a decayed concept still unlocks its successors', () => {
  // A three-node chain, so "unlocked by a" and "still locked on b" are distinguishable.
  const chain: CatalogEntry[] = [
    { id: 'a', label: 'A', prerequisites: [] },
    { id: 'b', label: 'B', prerequisites: ['a'] },
    { id: 'c', label: 'C', prerequisites: ['a', 'b'] },
  ]

  it('does not relock a learner over decay', () => {
    // ADR-037: DECAYING satisfies a prerequisite. Without that, a concept the learner has not
    // lost would lock everything downstream of it, and the whole graph would relock after a
    // fortnight away.
    const rows = upcomingLocked(chain, decayed, 'a')
    expect(rows.map(r => r.id)).toEqual(['c']) // b is unlocked by a; c still needs b
    expect(rows[0].missing).toEqual(['B'])
  })

  it('reports the same locked set whether a is MASTERED or DECAYING', () => {
    // The scheduler's prereqsMet accepts both, so "Coming up" must not change just because a
    // concept's status label moved. These two must never diverge.
    expect(upcomingLocked(chain, decayed, 'a')).toEqual(upcomingLocked(chain, mastered, 'a'))
  })
})
