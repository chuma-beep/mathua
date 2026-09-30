import { describe, it, expect } from 'vitest'
import { topoRank } from '../lib/topoRank'

describe('topoRank', () => {
  it('orders prereqs before dependents', () => {
    const rank = topoRank([
      { id: 'c', prerequisites: ['b'] },
      { id: 'b', prerequisites: ['a'] },
      { id: 'a' },
    ])
    expect(rank.get('a')).toBeLessThan(rank.get('b')!)
    expect(rank.get('b')).toBeLessThan(rank.get('c')!)
  })

  it('is deterministic with lexicographic tiebreaks', () => {
    const catalog = [
      { id: 'b' },
      { id: 'a' },
      { id: 'c', prerequisites: ['a', 'b'] },
    ]
    const r1 = topoRank(catalog)
    const r2 = topoRank([...catalog].reverse())
    expect([...r1.entries()]).toEqual([...r2.entries()])
    expect(r1.get('a')).toBeLessThan(r1.get('c')!)
    expect(r1.get('b')).toBeLessThan(r1.get('c')!)
  })

  it('ignores prereqs outside the catalog and dedupes', () => {
    const rank = topoRank([
      { id: 'a', prerequisites: ['missing', 'missing'] },
      { id: 'b', prerequisites: ['a'] },
    ])
    expect(rank.size).toBe(2)
    expect(rank.get('a')).toBeLessThan(rank.get('b')!)
  })

  it('ranks every node even with cycles', () => {
    const rank = topoRank([
      { id: 'a', prerequisites: ['b'] },
      { id: 'b', prerequisites: ['a'] },
    ])
    expect(rank.size).toBe(2)
  })
})
