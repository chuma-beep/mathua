import { describe, expect, it } from 'vitest'
import { angleBetween, dot, length, normalize, trace, type Mat2 } from '../components/visual/math/matrix'
import { maxError, nthDerivative, taylorCoeffs, taylorValue } from '../components/visual/math/series'
import { continuousAt, limitExists, limitValue, oneSidedLimits } from '../components/visual/math/limit'
import { draw, histogram, mean, populationPdf, rng, sampleMeans, sd } from '../components/visual/math/sampling'
import { correlation, intercept, means, predict, slope, sse, type Point } from '../components/visual/math/regression'

const close = (a: number, b: number, tolerance = 1e-12) => expect(Math.abs(a - b)).toBeLessThan(tolerance)

describe('matrix helpers added for the eigenvector figure', () => {
  it('takes a dot product, a length and a unit vector', () => {
    expect(dot([1, 2], [3, 4])).toBe(11)
    close(length([3, 4]), 5)
    close(normalize([3, 4])[0], 0.6)
    close(normalize([3, 4])[1], 0.8)
    expect(normalize([0, 0])).toEqual([0, 0])
  })

  it('measures the angle between vectors and reports the trace', () => {
    close(angleBetween([1, 0], [0, 1]), Math.PI / 2)
    close(angleBetween([1, 0], [-1, 0]), Math.PI)
    close(angleBetween([1, 0], [2, 0]), 0)
    const m: Mat2 = [
      [2, 1],
      [1, 2],
    ]
    expect(trace(m)).toBe(4)
  })
})

describe('Taylor polynomials', () => {
  it('reproduces the derivatives of the four functions at the centre', () => {
    expect(nthDerivative('exp', 0, 5)).toBe(1)
    close(nthDerivative('sin', 0, 2), 0)
    close(nthDerivative('sin', 0, 3), -1)
    close(nthDerivative('cos', 0, 0), 1)
    close(nthDerivative('cos', 0, 2), -1)
    // 1/(1-x) has derivative k!/(1-a)^(k+1); at a = 0 that is k!.
    close(nthDerivative('geometric', 0, 4), 24)
  })

  it('builds the sine polynomial about zero', () => {
    const c = taylorCoeffs('sin', 0, 5)
    close(c[0], 0)
    close(c[1], 1)
    close(c[2], 0)
    close(c[3], -1 / 6)
    close(c[4], 0)
    close(c[5], 1 / 120)
    close(taylorValue(c, 0, 1), 1 - 1 / 6 + 1 / 120)
  })

  it('builds the geometric series, and matches the function exactly at the centre', () => {
    const c = taylorCoeffs('geometric', 0, 4)
    for (let k = 0; k <= 4; k++) close(c[k], 1)
    for (const fn of ['sin', 'cos', 'exp'] as const) {
      const coeffs = taylorCoeffs(fn, 1.2, 6)
      const f = fn === 'sin' ? Math.sin(1.2) : fn === 'cos' ? Math.cos(1.2) : Math.exp(1.2)
      close(taylorValue(coeffs, 1.2, 1.2), f, 1e-10)
    }
  })

  it('reduces the maximum error as the degree rises', () => {
    const low = maxError('sin', 0, 3, -1, 1)
    const high = maxError('sin', 0, 7, -1, 1)
    expect(high).toBeLessThan(low)
    expect(low).toBeGreaterThan(0)
  })
})

describe('limits', () => {
  it('gives both sides of a removable hole the same value, and leaves f(a) undefined', () => {
    expect(oneSidedLimits('removable', 2)).toEqual({ left: 4, right: 4 })
    expect(limitExists('removable', 2)).toBe(true)
    expect(continuousAt('removable', 2)).toBe(false)
    close(limitValue('removable', 2, 3), 5)
    expect(Number.isNaN(limitValue('removable', 2, 2))).toBe(true)
  })

  it('reports a jump as two different values, so no limit exists', () => {
    expect(oneSidedLimits('jump', 0)).toEqual({ left: 1, right: 3 })
    expect(limitExists('jump', 0)).toBe(false)
    expect(continuousAt('jump', 0)).toBe(false)
  })

  it('reports a vertical asymptote as unbounded rather than as a number', () => {
    const { left, right } = oneSidedLimits('vertical_asymptote', 0)
    expect(left).toBe(Number.NEGATIVE_INFINITY)
    expect(right).toBe(Number.POSITIVE_INFINITY)
    expect(limitExists('vertical_asymptote', 0)).toBe(false)
    close(limitValue('vertical_asymptote', 0, -1), -1)
  })
})

describe('sampling distribution', () => {
  it('produces the same stream for the same seed', () => {
    const a = rng(7)
    const b = rng(7)
    for (let i = 0; i < 5; i++) expect(a()).toBe(b())
    expect(rng(7)()).not.toBe(rng(8)())
  })

  it('gives every population mean zero and spread one', () => {
    for (const population of ['uniform', 'normal', 'skewed'] as const) {
      const next = rng(3)
      const values = Array.from({ length: 4000 }, () => draw(population, next))
      expect(Math.abs(mean(values))).toBeLessThan(0.12)
      expect(Math.abs(sd(values) - 1)).toBeLessThan(0.12)
    }
  })

  it('narrows the sample means like 1 over the square root of n', () => {
    const wide = sd(sampleMeans('skewed', 1, 1500, 5))
    const narrow = sd(sampleMeans('skewed', 25, 1500, 5))
    close(wide / narrow, 5, 0.9)
    expect(Math.abs(mean(sampleMeans('skewed', 4, 1500, 5)))).toBeLessThan(0.12)
  })

  it('bins values and scales the population shape to mean zero spread one', () => {
    const counts = histogram([-1, 0, 0, 0.2, 5], -1, 1, 4)
    expect(counts.reduce((a, b) => a + b, 0)).toBe(4)
    close(populationPdf('normal', 0), 1 / Math.sqrt(2 * Math.PI))
    close(populationPdf('uniform', 0), 1 / (2 * Math.sqrt(3)))
    close(populationPdf('skewed', -2), 0)
  })
})

describe('least squares', () => {
  const perfect: Point[] = [
    [0, 1],
    [1, 3],
    [2, 5],
    [3, 7],
  ]

  it('recovers an exact line', () => {
    close(slope(perfect), 2)
    close(intercept(perfect), 1)
    close(correlation(perfect), 1)
    close(sse(perfect), 0)
    close(predict(perfect, 10), 21)
  })

  it('passes through the mean point and leaves a residual otherwise', () => {
    const cloud: Point[] = [
      [1, 2],
      [2, 3],
      [3, 3.5],
      [4, 6],
      [5, 6.5],
    ]
    const [mx, my] = means(cloud)
    close(mx, 3)
    close(predict(cloud, mx), my)
    expect(sse(cloud)).toBeGreaterThan(0)
    expect(correlation(cloud)).toBeGreaterThan(0.9)
  })

  it('refuses to invent a slope when every x is the same', () => {
    const vertical: Point[] = [
      [2, 1],
      [2, 4],
      [2, 9],
    ]
    expect(slope(vertical)).toBe(0)
    expect(correlation(vertical)).toBe(0)
  })
})
