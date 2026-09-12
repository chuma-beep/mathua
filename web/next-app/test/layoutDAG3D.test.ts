import { describe, it, expect } from 'vitest'
import { layoutDAG3D, type ConceptLayoutInput } from '../lib/layoutDAG3D'

function chain(length: number): ConceptLayoutInput[] {
  return Array.from({ length }, (_, i) => ({
    id: `c${i}`,
    prerequisites: i === 0 ? [] : [`c${i - 1}`],
    domain: i % 2 === 0 ? 'alpha' : 'beta',
  }))
}

describe('layoutDAG3D', () => {
  it('is deterministic across runs', () => {
    const input = chain(60)
    const a = layoutDAG3D(input, { domainOrder: ['alpha', 'beta'] })
    const b = layoutDAG3D(input, { domainOrder: ['alpha', 'beta'] })
    expect(a.positions).toEqual(b.positions)
  })

  it('assigns layers by longest prerequisite path', () => {
    const input: ConceptLayoutInput[] = [
      { id: 'root', prerequisites: [], domain: 'a' },
      { id: 'mid', prerequisites: ['root'], domain: 'a' },
      { id: 'alt', prerequisites: ['root'], domain: 'a' },
      { id: 'leaf', prerequisites: ['mid', 'alt'], domain: 'a' },
    ]
    const { layers } = layoutDAG3D(input)
    expect(layers.root).toBe(0)
    expect(layers.mid).toBe(1)
    expect(layers.alt).toBe(1)
    expect(layers.leaf).toBe(2)
  })

  it('places every node and keeps the cloud inside the target radius', () => {
    const input = chain(200)
    const targetRadius = 15
    const { positions } = layoutDAG3D(input, { targetRadius })
    expect(Object.keys(positions)).toHaveLength(200)
    for (const p of Object.values(positions)) {
      expect(Number.isFinite(p[0])).toBe(true)
      expect(Math.hypot(p[0], p[1], p[2])).toBeLessThanOrEqual(targetRadius + 1e-6)
    }
  })

  it('keeps the silhouette round (no column collapse)', () => {
    const targetRadius = 15
    const input = chain(120)
    const { positions } = layoutDAG3D(input, { targetRadius })
    const xs = Object.values(positions).map(p => p[0])
    const ys = Object.values(positions).map(p => p[1])
    const zs = Object.values(positions).map(p => p[2])
    const span = (v: number[]) => Math.max(...v) - Math.min(...v)
    const spanY = span(ys)
    const horizontal = Math.max(span(xs), span(zs))
    // A sphere of radius 15 spans ~30 in every axis; a column would not.
    expect(spanY).toBeGreaterThan(targetRadius * 1.6)
    expect(spanY).toBeLessThan(targetRadius * 2.2)
    expect(horizontal).toBeGreaterThan(targetRadius * 1.6)
    expect(horizontal).toBeLessThan(targetRadius * 2.2)
  })

  it('clusters nodes of the same domain into a contiguous sector', () => {
    const input: ConceptLayoutInput[] = []
    for (let i = 0; i < 60; i++) {
      input.push({ id: `a${i}`, prerequisites: [], domain: 'alpha' })
      input.push({ id: `b${i}`, prerequisites: [], domain: 'beta' })
    }
    const { positions } = layoutDAG3D(input, { domainOrder: ['alpha', 'beta'] })
    const angle = (p: [number, number, number]) => Math.atan2(p[2], p[0])
    const alpha = input.flatMap(c => (c.domain === 'alpha' ? [angle(positions[c.id])] : []))
    const beta = input.flatMap(c => (c.domain === 'beta' ? [angle(positions[c.id])] : []))
    // alpha owns the first sector (angle >= 0), beta the second (angle < 0).
    for (const a of alpha) {
      expect(a).toBeGreaterThan(-0.1)
      expect(a).toBeLessThan(Math.PI + 0.1)
    }
    for (const b of beta) {
      expect(b).toBeLessThan(0.1)
      expect(b).toBeGreaterThan(-Math.PI - 0.1)
    }
  })

  it('lays out 5k nodes fast enough for a runtime fallback', () => {
    const input: ConceptLayoutInput[] = []
    for (let i = 0; i < 5000; i++) {
      input.push({
        id: `n${i}`,
        prerequisites: i === 0 ? [] : [`n${i - 1}`],
        domain: `d${i % 17}`,
      })
    }
    const t0 = performance.now()
    const { positions } = layoutDAG3D(input)
    const ms = performance.now() - t0
    expect(Object.keys(positions)).toHaveLength(5000)
    expect(ms).toBeLessThan(1000)
  })

  it('handles empty and cyclic input without throwing', () => {
    expect(layoutDAG3D([]).positions).toEqual({})
    const cyclic: ConceptLayoutInput[] = [
      { id: 'x', prerequisites: ['y'], domain: 'a' },
      { id: 'y', prerequisites: ['x'], domain: 'a' },
    ]
    const { positions } = layoutDAG3D(cyclic)
    expect(Object.keys(positions).sort()).toEqual(['x', 'y'])
  })
})
