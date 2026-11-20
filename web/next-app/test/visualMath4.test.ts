import { describe, expect, it } from 'vitest'
import { cfCritical, cfDerivative, cfGradients, cfLevelValues, cfPoles, cfValue, cfZeros, type ComplexFn } from '../components/visual/math/complexField'
import { dot } from '../components/visual/math/matrix'

const close = (a: number, b: number, tolerance = 1e-12) => expect(Math.abs(a - b)).toBeLessThan(tolerance)
const regular: Array<[ComplexFn, [number, number]]> = [
  ['square', [1.2, 0.7]],
  ['cube', [0.8, -1.1]],
  ['exp', [0.3, 1.4]],
  ['inverse', [1.5, -0.6]],
  ['sine', [1.1, 0.5]],
]

describe('complex field values and derivatives', () => {
  it('evaluates the five functions', () => {
    expect(cfValue('square', [1, 2])).toEqual([-3, 4])
    close(cfValue('exp', [0, 0])![0], 1)
    close(cfValue('exp', [0, 0])![1], 0)
    expect(cfValue('inverse', [2, 0])).toEqual([0.5, -0])
    expect(cfValue('inverse', [0, 0])).toBeNull()
    close(cfValue('sine', [0, 0])![0], 0)
    close(cfValue('cube', [1, 1])![0], -2)
    close(cfValue('cube', [1, 1])![1], 2)
  })

  it('differentiates each one', () => {
    expect(cfDerivative('square', [1, 2])).toEqual([2, 4])
    expect(cfDerivative('cube', [1, 1])).toEqual([0, 6])
    close(cfDerivative('inverse', [2, 0])![0], -0.25)
    close(cfDerivative('exp', [0, 0])![0], 1)
    close(cfDerivative('sine', [0, 0])![0], 1)
  })
})

describe('Cauchy-Riemann and orthogonality', () => {
  it('satisfies the Cauchy-Riemann equations everywhere it is analytic', () => {
    for (const [fn, z] of regular) {
      const g = cfGradients(fn, z)!
      close(g.du[0], g.dv[1], 1e-9)
      close(g.du[1], -g.dv[0], 1e-9)
    }
  })

  it('makes the two level families perpendicular wherever the derivative is not zero', () => {
    for (const [fn, z] of regular) {
      const g = cfGradients(fn, z)!
      close(dot(g.du, g.dv), 0, 1e-9)
    }
  })

  it('reports a vanishing derivative at the critical points and nowhere else', () => {
    expect(cfDerivative('square', [0, 0])).toEqual([0, 0])
    expect(cfDerivative('cube', [0, 0])).toEqual([0, 0])
    expect(cfCritical('exp')).toEqual([])
    expect(cfCritical('inverse')).toEqual([])
    expect(cfCritical('sine')).toHaveLength(2)
    expect(cfZeros('exp')).toEqual([])
    expect(cfZeros('sine')).toHaveLength(3)
    // 1/z has a pole at the origin, not a zero: it never vanishes anywhere.
    expect(cfZeros('inverse')).toEqual([])
    expect(cfPoles('inverse')).toEqual([[0, 0]])
    for (const fn of ['square', 'cube', 'exp', 'sine'] as ComplexFn[]) expect(cfPoles(fn)).toEqual([])
    // At the origin of z squared both partials vanish, so there is no tangent to be right-angled.
    const g = cfGradients('square', [0, 0])!
    expect(g.du).toEqual([0, -0])
    expect(g.dv).toEqual([0, 0])
  })
})

describe('level values', () => {
  it('returns ascending levels above the minimum', () => {
    const levels = cfLevelValues('square', 'u', 2, 5)
    expect(levels).toHaveLength(5)
    for (let i = 1; i < levels.length; i++) expect(levels[i]).toBeGreaterThan(levels[i - 1])
    expect(levels[0]).toBeGreaterThan(-4)
  })
})
