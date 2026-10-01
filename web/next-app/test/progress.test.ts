import { describe, it, expect } from 'vitest'
import { isMastered, prereqMet, countByDomain, countOverall, type CatalogueConcept } from '../lib/progress'
import type { ConceptProgress } from '../lib/api'

const p = (status: string, completed?: boolean): ConceptProgress => ({ status, streak: 0, completed })

// a ──► b ──► c, with d independent.
const catalogue: CatalogueConcept[] = [
  { id: 'a', domain: 'x' },
  { id: 'b', domain: 'x', prerequisites: ['a'] },
  { id: 'c', domain: 'x', prerequisites: ['b'] },
  { id: 'd', domain: 'y' },
]

describe('isMastered', () => {
  it('is case-insensitive, because the stored value is uppercase and older payloads are not', () => {
    expect(isMastered(p('MASTERED'))).toBe(true)
    expect(isMastered(p('mastered'))).toBe(true)
    expect(isMastered(p('Mastered'))).toBe(true)
  })

  it('does not treat decaying as mastered', () => {
    // DECAYING is what the server reports once a review goes stale. Counting it
    // as mastered is exactly the bug this shared helper exists to prevent.
    expect(isMastered(p('DECAYING'))).toBe(false)
  })

  it('is false for every other status and for nothing at all', () => {
    for (const s of ['UNSEEN', 'LEARNING', 'PRACTICING', '']) {
      expect(isMastered(p(s))).toBe(false)
    }
    expect(isMastered(undefined)).toBe(false)
  })
})

describe('prereqMet', () => {
  it('accepts a mastered prerequisite', () => {
    expect(prereqMet({ a: p('MASTERED') }, 'a')).toBe(true)
  })

  it('accepts a completed-but-unmastered prerequisite', () => {
    expect(prereqMet({ a: p('LEARNING', true) }, 'a')).toBe(true)
  })

  it('rejects an untouched one, and an absent one', () => {
    expect(prereqMet({ a: p('LEARNING') }, 'a')).toBe(false)
    expect(prereqMet({}, 'a')).toBe(false)
  })
})

describe('countByDomain', () => {
  it('buckets the catalogue and buckets a concept once', () => {
    // The buckets are mutually exclusive: a mastered concept is not also
    // counted as learning or locked. With b mastered and c practicing, a has
    // no prerequisites and c's are met, so nothing in x is locked.
    const got = countByDomain(catalogue, { b: p('MASTERED'), c: p('PRACTICING') })
    expect(got.get('x')).toEqual({ total: 3, mastered: 1, learning: 1, completed: 0, locked: 0 })
    expect(got.get('y')).toEqual({ total: 1, mastered: 0, learning: 0, completed: 0, locked: 0 })
  })

  it('locks a concept whose prerequisite is unmet', () => {
    // b needs a, so with nothing done b is locked.
    expect(countByDomain(catalogue, {}).get('x')?.locked).toBe(2)
    // a mastered unblocks b but not c.
    expect(countByDomain(catalogue, { a: p('MASTERED') }).get('x')?.locked).toBe(1)
    expect(countByDomain(catalogue, { a: p('MASTERED'), b: p('MASTERED') }).get('x')?.locked).toBe(0)
  })

  it('treats a completed prerequisite as met for the lock count', () => {
    expect(countByDomain(catalogue, { a: p('LEARNING', true) }).get('x')?.locked).toBe(1)
  })
})

describe('countOverall', () => {
  it('reports total, mastered and unlocked', () => {
    const got = countOverall(catalogue, { a: p('MASTERED'), b: p('PRACTICING') })
    expect(got.total).toBe(4)
    expect(got.mastered).toBe(1)
    expect(got.learning).toBe(1)
    // b and d are reachable; c is not, because b is only practicing.
    expect(got.locked).toBe(1)
    expect(got.unlocked).toBe(3)
  })

  it('reports unlocked as the complement of locked, for any input', () => {
    for (const progress of [
      {},
      { a: p('MASTERED') },
      { a: p('MASTERED'), b: p('MASTERED') },
      { c: p('LEARNING') },
      { d: p('DECAYING') },
    ]) {
      const got = countOverall(catalogue, progress)
      expect(got.unlocked).toBe(got.total - got.locked)
    }
  })

  it('never reports more unlocked than total', () => {
    const got = countOverall(catalogue, {})
    expect(got.unlocked).toBeLessThanOrEqual(got.total)
    expect(got.unlocked).toBeGreaterThanOrEqual(0)
  })
})