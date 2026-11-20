import { describe, expect, it } from 'vitest'
import {
  analytic,
  fnValue,
  integrate,
  pathAt,
  poleInfo,
  type Complex,
  type ContourFn,
  type Path,
} from '../components/visual/math/contourIntegral'

const circle = (radius: number, center: [number, number] = [0, 0]): Path => ({ kind: 'circle', center, radius })

function expectComplex(actual: [number, number], re: number, im: number, tolerance = 1e-9) {
  expect(Math.abs(actual[0] - re)).toBeLessThan(tolerance)
  expect(Math.abs(actual[1] - im)).toBeLessThan(tolerance)
}

describe('contour integration', () => {
  it('gives 2 pi i for 1/z on the unit circle, and the negative when reversed', () => {
    expectComplex(integrate(circle(1), 'inv', 1).total, 0, 2 * Math.PI)
    expectComplex(integrate(circle(1), 'inv', -1).total, 0, -2 * Math.PI)
  })

  it('gives zero for every entire function on a closed path', () => {
    for (const fn of ['identity', 'square', 'exp'] as ContourFn[]) {
      expectComplex(integrate(circle(1), fn, 1).total, 0, 0)
      expectComplex(integrate(circle(1), fn, -1).total, 0, 0)
    }
  })

  it('gives 2 pi i r squared for the conjugate, which is not analytic', () => {
    expectComplex(integrate(circle(1), 'conj', 1).total, 0, 2 * Math.PI)
    expectComplex(integrate(circle(2), 'conj', 1).total, 0, 8 * Math.PI, 1e-8)
  })

  it('keeps 2 pi i for an ellipse and for a wavy closed curve around the pole', () => {
    const ellipse: Path = { kind: 'ellipse', center: [0, 0], rx: 2, ry: 1 }
    const wavy: Path = { kind: 'wavy', center: [0, 0], radius: 1, amplitude: 0.2, lobes: 3 }
    expectComplex(integrate(ellipse, 'inv', 1).total, 0, 2 * Math.PI, 1e-7)
    expectComplex(integrate(wavy, 'inv', 1).total, 0, 2 * Math.PI, 1e-7)
  })

  it('gives zero for a circle that does not enclose the pole', () => {
    expectComplex(integrate(circle(1, [3, 0]), 'inv', 1).total, 0, 0, 1e-9)
  })

  it('refuses a number when the path passes through the pole', () => {
    const through: Path = { kind: 'line', a: [-1, -1], b: [1, 1] }
    expect(poleInfo(through).crosses).toBe(true)
    const result = integrate(through, 'inv', 1)
    expect(Number.isNaN(result.total[0])).toBe(true)
    expect(result.partial).toEqual([])
    // The same segment moved clear of the origin is finite.
    const clear: Path = { kind: 'line', a: [-1, 0.5], b: [1, 0.5] }
    expect(poleInfo(clear).crosses).toBe(false)
    expect(Number.isFinite(integrate(clear, 'inv', 1).total[0])).toBe(true)
  })

  it('accumulates monotonically along the path and ends at the total', () => {
    const { partial, total } = integrate(circle(1), 'inv', 1, 200)
    expect(partial).toHaveLength(200)
    expectComplex(partial[partial.length - 1], total[0], total[1], 1e-9)
    const halfway = partial[99]
    expect(halfway[1]).toBeGreaterThan(0)
    expect(halfway[1]).toBeLessThan(2 * Math.PI)
  })

  it('agrees with the closed form wherever one exists', () => {
    const paths: Path[] = [circle(1), circle(1, [3, 0]), { kind: 'ellipse', center: [0, 0], rx: 2, ry: 1 }]
    for (const path of paths) {
      for (const fn of ['inv', 'identity', 'square', 'exp'] as ContourFn[]) {
        const expected = analytic(path, fn, 1)
        if (!expected) continue
        expectComplex(integrate(path, fn, 1).total, expected[0], expected[1], 1e-7)
      }
    }
  })

  it('has no value at the pole itself', () => {
    expect(fnValue('inv', [0, 0])).toBeNull()
    expect(pathAt(circle(1), 1, 0).z[0]).toBeCloseTo(1, 12)
  })
})

describe('pole detection', () => {
  it('solves a segment exactly', () => {
    const crossing = poleInfo({ kind: 'line', a: [-1, -1], b: [1, 1] })
    expect(crossing.crosses).toBe(true)
    expect(crossing.certainty).toBe('exact')
    expect(crossing.closest).toBeLessThan(1e-12)

    const clear = poleInfo({ kind: 'line', a: [0, 1], b: [2, 1] })
    expect(clear.crosses).toBe(false)
    expect(clear.closest).toBeCloseTo(1, 12)
  })

  it('solves a circle exactly, and keeps a 10^-4 near miss a near miss', () => {
    const through = poleInfo({ kind: 'circle', center: [0.6, 0.8], radius: 1 })
    expect(through.crosses).toBe(true)
    expect(through.certainty).toBe('exact')
    expect(through.closest).toBeLessThan(1e-12)

    const miss = poleInfo({ kind: 'circle', center: [0.6, 0.8], radius: 1 - 1e-4 })
    expect(miss.crosses).toBe(false)
    expect(miss.certainty).toBe('exact')
    expect(miss.closest).toBeCloseTo(1e-4, 9)
  })

  it('finds a pole on an ellipse at a parameter between the old 4000 samples', () => {
    // The origin sits on this ellipse at t0, and t0 * 4000 is not an integer, so a
    // fixed grid of 4000 samples would have stepped over the crossing.
    const t0 = 0.13703
    expect(Number.isInteger(t0 * 4000)).toBe(false)
    const path = { kind: 'ellipse' as const, center: [-2 * Math.cos(t0), -Math.sin(t0)] as Complex, rx: 2, ry: 1 }
    const info = poleInfo(path)
    expect(info.certainty).toBe('converged')
    expect(info.closest).toBeLessThan(1e-9)
    expect(info.crosses).toBe(true)
  })

  it('keeps a small ellipse miss out of the crossing verdict', () => {
    const t0 = 0.13703
    const path = { kind: 'ellipse' as const, center: [-2 * Math.cos(t0), -Math.sin(t0)] as Complex, rx: 2 + 5e-4, ry: 1 }
    const info = poleInfo(path)
    expect(info.certainty).toBe('converged')
    expect(info.crosses).toBe(false)
    expect(info.closest).toBeGreaterThan(1e-9)
    expect(info.closest).toBeLessThan(1e-2)
  })

  it('finds a wavy curve touching the pole, and a 10^-4 miss beside it', () => {
    const touch = poleInfo({ kind: 'wavy', center: [0, 0], radius: 0.5, amplitude: 0.5, lobes: 3 })
    expect(touch.certainty).toBe('converged')
    expect(touch.crosses).toBe(true)
    expect(touch.closest).toBeLessThan(1e-9)

    const miss = poleInfo({ kind: 'wavy', center: [0, 0], radius: 0.5, amplitude: 0.4999, lobes: 3 })
    expect(miss.certainty).toBe('converged')
    expect(miss.crosses).toBe(false)
    expect(miss.closest).toBeCloseTo(1e-4, 8)
  })
})
