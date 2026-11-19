import { describe, expect, it } from 'vitest'
import { basisMatrix, combine, coordinatesOf, isNearlyDependent } from '../components/visual/math/basis'
import type { Vec2 } from '../components/visual/math/matrix'
import { mean, deviations, sumSquares, variance, stdDev } from '../components/visual/math/spread'
import { angleBetweenDirections, angleFactorAt, mapDerivative, mapIsAnalytic, mapIsConformalAt, mapValue, mappedDirectionAngle, wrapAngle } from '../components/visual/math/complexMap'
import { distance, intersect, lineThrough, midpoint, onLine } from '../components/visual/math/coords'
import { classifyQuadratic, discriminant, xIntercepts, yIntercepts } from '../components/visual/math/conic'
import { contourSegments } from '../components/visual/math/contour'
import { fieldGradient, fieldLevels, fieldValue, gradientLength } from '../components/visual/math/field'

const close = (a: number, b: number, tolerance = 1e-12) => expect(Math.abs(a - b)).toBeLessThan(tolerance)

describe('basis coordinates', () => {
  const u: Vec2 = [1.5, 0.5]
  const v: Vec2 = [0.5, 1.5]

  it('builds a point from its coordinates', () => {
    const basis = basisMatrix(u, v)
    expect(combine(basis, [1, 0])).toEqual(u)
    expect(combine(basis, [0, 1])).toEqual(v)
    const p = combine(basis, [2, 1])
    close(p[0], 3.5)
    close(p[1], 2.5)
  })

  it('reads the coordinates back out', () => {
    const basis = basisMatrix(u, v)
    const coords = coordinatesOf(basis, [3.5, 2.5])
    expect(coords).not.toBeNull()
    close(coords![0], 2)
    close(coords![1], 1)
  })

  it('gives the standard basis the point itself', () => {
    const identity = basisMatrix([1, 0], [0, 1])
    expect(coordinatesOf(identity, [4, -2])).toEqual([4, -2])
  })

  it('has no answer when the basis is dependent', () => {
    const degenerate = basisMatrix([1, 1], [2, 2])
    expect(coordinatesOf(degenerate, [3, 0])).toBeNull()
  })

  it('halves the first coordinate when the first vector doubles', () => {
    const scaled = basisMatrix([3, 1], v)
    const original = coordinatesOf(basisMatrix(u, v), [3.5, 2.5])!
    const next = coordinatesOf(scaled, [3.5, 2.5])!
    close(next[0], original[0] / 2)
    close(next[1], original[1])
  })
})

describe('spread', () => {
  it('balances the signed gaps to zero', () => {
    const values = [1, 2, 2, 4, 6]
    close(mean(values), 3)
    close(deviations(values).reduce((a, b) => a + b, 0), 0)
    expect(deviations(values)).toEqual([-2, -1, -1, 1, 3])
  })

  it('divides by n - 1 for a sample and by n for a population', () => {
    const values = [2, 4, 4, 4, 5, 5, 7, 9]
    close(sumSquares(values), 32)
    close(variance(values, false), 4)
    close(variance(values, true), 32 / 7)
    close(stdDev(values, false), 2)
  })

  it('is unchanged by a shift and scales by the absolute factor', () => {
    const values = [1, 2, 4, 8]
    close(variance(values.map((v) => v + 10), false), variance(values, false))
    close(stdDev(values.map((v) => v * -3), false), 3 * stdDev(values, false))
  })

  it('reduces to zero when every value is the same', () => {
    close(variance([5, 5, 5], false), 0)
    close(stdDev([5, 5, 5], true), 0)
  })
})

describe('conic classification', () => {
  it('computes the discriminant the lesson states', () => {
    expect(discriminant(1, 0, 1)).toBe(-4)
    expect(discriminant(1, 0, 0)).toBe(0)
    expect(discriminant(1, 0, -1)).toBe(4)
    expect(discriminant(1, 2, 1)).toBe(0)
  })

  it('names the type from the quadratic part, not the discriminant alone', () => {
    expect(classifyQuadratic(1, 0, 1).type).toBe('circle')
    expect(classifyQuadratic(1, 0, 0.5).type).toBe('ellipse')
    expect(classifyQuadratic(1, 0, -1).type).toBe('hyperbola')
    // Zero discriminant here is a degenerate quadratic form, so the level set is
    // two parallel lines: x² = 1 is x = ±1. A parabola needs a linear term.
    expect(classifyQuadratic(1, 0, 0).type).toBe('parallel')
    expect(classifyQuadratic(1, 2, 1).type).toBe('parallel')
    // Negative definite: a negative discriminant but no real points at all.
    expect(classifyQuadratic(-1, 0, -1).type).toBe('empty')
    expect(classifyQuadratic(0, 0, 0).type).toBe('empty')
    // Rank one with a negative eigenvalue has no real points either.
    expect(classifyQuadratic(-1, 0, 0).type).toBe('empty')
  })

  it('reports eccentricity, and none for a degenerate form', () => {
    expect(classifyQuadratic(1, 0, 1).eccentricity).toBe(0)
    close(classifyQuadratic(1, 0, 0.5).eccentricity, Math.sqrt(0.5))
    close(classifyQuadratic(1, 0, -1).eccentricity, Math.SQRT2)
    expect(Number.isNaN(classifyQuadratic(1, 0, 0).eccentricity)).toBe(true)
  })

  it('counts squared terms and compares their signs', () => {
    expect(classifyQuadratic(1, 0, 0.5).squaredTerms).toBe(2)
    expect(classifyQuadratic(1, 0, 0.5).sameSign).toBe(true)
    expect(classifyQuadratic(1, 0, -1).sameSign).toBe(false)
    expect(classifyQuadratic(1, 0, 0).squaredTerms).toBe(1)
  })

  it('places the intercepts only where the curve meets an axis', () => {
    expect(xIntercepts(1, 0.5)).toEqual([-1, 1])
    close(yIntercepts(1, 0.5)[1], Math.SQRT2)
    expect(xIntercepts(-1, -1)).toEqual([])
  })
})

describe('marching squares', () => {
  it('draws a circle of radius one for x^2 + y^2 - 1', () => {
    const segments = contourSegments(
      (x, y) => x * x + y * y - 1,
      { x0: -2, x1: 2, y0: -2, y1: 2 },
      0,
      80,
      80,
    )
    expect(segments.length).toBeGreaterThan(30)
    for (const [a, b] of segments) {
      close(Math.hypot(a[0], a[1]), 1, 0.03)
      close(Math.hypot(b[0], b[1]), 1, 0.03)
    }
  })

  it('skips a cell whose corner is not finite', () => {
    const segments = contourSegments(
      (x) => (x === 0 ? Number.NaN : 1 / x),
      { x0: -1, x1: 1, y0: -1, y1: 1 },
      2,
      8,
      8,
    )
    expect(segments.every(([a, b]) => Number.isFinite(a[0]) && Number.isFinite(b[0]))).toBe(true)
  })
})

describe('contour fields', () => {
  it('evaluates the four fields', () => {
    close(fieldValue('bowl', 1, 2), 5)
    close(fieldValue('saddle', 1, 2), -3)
    close(fieldValue('product', 2, 3), 6)
    close(fieldValue('ripple', 0, 0), 1)
  })

  it('gives the exact gradient, not a finite difference', () => {
    expect(fieldGradient('bowl', 1, 2)).toEqual([2, 4])
    expect(fieldGradient('saddle', 1, 2)).toEqual([2, -4])
    expect(fieldGradient('product', 2, 3)).toEqual([3, 2])
    close(fieldGradient('ripple', 0, 0)[0], 1)
    close(fieldGradient('ripple', 0, 0)[1], 0)
    close(gradientLength('bowl', 1, 1), Math.sqrt(8))
  })

  it('spreads the levels over the range and never sits on a single minimum', () => {
    const levels = fieldLevels('bowl', 3, 5)
    expect(levels).toHaveLength(5)
    expect(levels.every((value, i) => i === 0 || value > levels[i - 1])).toBe(true)
    expect(levels[0]).toBeGreaterThan(0)
    expect(levels[4]).toBeLessThan(18)
  })
})

describe('complex maps', () => {
  it('evaluates each map and its derivative', () => {
    expect(mapValue('square', [2, 3])).toEqual([-5, 12])
    expect(mapValue('exp', [0, 0])).toEqual([1, 0])
    expect(mapValue('inverse', [1, 0])).toEqual([1, -0])
    expect(mapValue('conjugate', [2, -3])).toEqual([2, 3])
    expect(mapDerivative('square', [1, 1])).toEqual([2, 2])
    expect(mapDerivative('exp', [0, 0])).toEqual([1, 0])
    expect(mapDerivative('inverse', [1, 0])).toEqual([-1, 0])
    expect(mapDerivative('conjugate', [1, 1])).toBeNull()
    expect(mapIsAnalytic('conjugate')).toBe(false)
    expect(mapIsAnalytic('square')).toBe(true)
  })

  it('keeps the angle between two directions under an analytic map', () => {
    const z: [number, number] = [0.8, 0.4]
    for (const map of ['square', 'exp', 'inverse'] as const) {
      const before = [-0.4, 0.9]
      const after = before.map((angle) => mappedDirectionAngle(map, z, angle))
      const size = (a: number, b: number) => Math.abs(wrapAngle(a - b))
      close(size(after[0], after[1]), size(before[0], before[1]), 1e-9)
    }
  })

  it('reflects the angle under conjugation instead of rotating it', () => {
    const z: [number, number] = [1, 1]
    close(mappedDirectionAngle('conjugate', z, 0.7), -0.7, 1e-12)
    // The size is preserved, the sense is not.
    close(angleBetweenDirections(-0.7, 0.3), angleBetweenDirections(0.7, -0.3), 1e-12)
  })

  it('wraps angles into one turn', () => {
    close(wrapAngle(3 * Math.PI), Math.PI, 1e-12)
    close(wrapAngle(-3 * Math.PI), Math.PI, 1e-12)
    close(wrapAngle(0.5), 0.5, 1e-12)
  })
})

describe('coordinate geometry', () => {
  it('measures a 3-4-5 triangle and its midpoint', () => {
    expect(distance([0, 0], [3, 4])).toBe(5)
    expect(midpoint([-3, 1], [2, 4])).toEqual([-0.5, 2.5])
    expect(midpoint([2, 2], [2, 2])).toEqual([2, 2])
  })

  it('reads a line from two points, vertical included', () => {
    const sloped = lineThrough([-3, 1], [2, 4])
    expect(sloped.kind).toBe('slope')
    if (sloped.kind === 'slope') {
      close(sloped.m, 0.6)
      close(sloped.b, 2.8)
    }
    const vertical = lineThrough([2, 1], [2, 5])
    expect(vertical).toEqual({ kind: 'vertical', x: 2 })
  })

  it('distinguishes meeting, parallel and coincident lines', () => {
    const first = lineThrough([-3, 1], [2, 4])
    const steep = lineThrough([1, -2], [3, 2])
    const meet = intersect(first, steep)
    expect(meet.kind).toBe('point')
    if (meet.kind === 'point') {
      close(meet.x, 34 / 7)
      close(meet.y, 40 / 7)
    }
    expect(intersect(first, lineThrough([0, 3], [5, 6]))).toEqual({ kind: 'parallel' })
    expect(intersect(first, lineThrough([0, 2.8], [5, 5.8]))).toEqual({ kind: 'coincident' })
  })

  it('meets a vertical line at one x and rejects a parallel vertical', () => {
    const vertical = lineThrough([2, -3], [2, 3])
    const meet = intersect(vertical, lineThrough([0, 0], [4, 4]))
    expect(meet.kind).toBe('point')
    if (meet.kind === 'point') {
      close(meet.x, 2)
      close(meet.y, 2)
    }
    expect(intersect(vertical, lineThrough([2, -1], [2, 7]))).toEqual({ kind: 'coincident' })
    expect(intersect(vertical, lineThrough([3, 0], [3, 4]))).toEqual({ kind: 'parallel' })
  })

  it('reports whether a point sits on a line', () => {
    const line = lineThrough([0, 1], [2, 5])
    expect(onLine([1, 3], line)).toBe(true)
    expect(onLine([1, 4], line)).toBe(false)
  })
})

describe('basis conditioning', () => {
  it('treats an exactly dependent pair as degenerate and a near one as fragile', () => {
    const dependent = basisMatrix([1, 0], [2, 0])
    expect(coordinatesOf(dependent, [1, 1])).toBeNull()
    expect(isNearlyDependent(dependent)).toBe(true)

    const thin = basisMatrix([1, 0], [1.0001, 0.0001])
    expect(coordinatesOf(thin, [2, 1])).not.toBeNull()
    expect(isNearlyDependent(thin)).toBe(true)

    const healthy = basisMatrix([1, 0], [0, 1])
    expect(isNearlyDependent(healthy)).toBe(false)
  })

  it('judges conditioning relative to the lengths, not the bare determinant', () => {
    // A long, well-spread basis has a large determinant and is perfectly usable.
    const wide = basisMatrix([100, 0], [0, 100])
    expect(isNearlyDependent(wide)).toBe(false)
    // The same shape scaled down stays well conditioned even though det is tiny.
    const tiny = basisMatrix([0.01, 0], [0, 0.01])
    expect(isNearlyDependent(tiny)).toBe(false)
  })
})

describe('complex map critical points', () => {
  it('doubles angles at the critical point of z squared instead of preserving them', () => {
    expect(angleFactorAt('square', [0, 0])).toBe(2)
    expect(angleFactorAt('square', [1, 0])).toBe(1)
    expect(mapIsConformalAt('square', [0, 0])).toBe(false)
    expect(mapIsConformalAt('square', [1, 0])).toBe(true)
    expect(mappedDirectionAngle('square', [0, 0], Math.PI / 6)).toBeCloseTo(Math.PI / 3, 12)
  })

  it('has no angle at a pole and no conformality for a reflection', () => {
    expect(angleFactorAt('inverse', [0, 0])).toBeNull()
    expect(mapIsConformalAt('inverse', [0, 0])).toBe(false)
    expect(Number.isNaN(mappedDirectionAngle('inverse', [0, 0], 0.4))).toBe(true)
    expect(angleFactorAt('conjugate', [1, 0])).toBe(1)
    expect(mapIsConformalAt('conjugate', [1, 0])).toBe(false)
    expect(mapIsConformalAt('exp', [0, 0])).toBe(true)
  })
})
