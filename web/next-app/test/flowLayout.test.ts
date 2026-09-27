import { describe, it, expect } from 'vitest'
import { precomputedLayout } from '../lib/flowLayout'
import { concepts } from '../lib/conceptData'

describe('precomputedLayout', () => {
  it('covers every concept in the bundled corpus', () => {
    // Guards artifact drift: a content wave that adds concepts without running
    // `npm run flow:build` would silently drop the graph back to dagre.
    expect(precomputedLayout(concepts.map(c => c.id))).not.toBeNull()
  })

  it('returns positions for known concepts', () => {
    const layout = precomputedLayout(['arith.add.single'])
    expect(layout).not.toBeNull()
    expect(layout!.get('arith.add.single')).toEqual({
      x: expect.any(Number),
      y: expect.any(Number),
    })
  })

  it('returns null when any id is missing, so the caller falls back wholesale', () => {
    expect(precomputedLayout(['arith.add.single', 'does.not.exist'])).toBeNull()
  })

  it('returns null for an empty set', () => {
    expect(precomputedLayout([])).toBeNull()
  })
})
