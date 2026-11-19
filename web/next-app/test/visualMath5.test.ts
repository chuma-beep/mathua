import { describe, expect, it } from 'vitest'
import { FAMILIES, familyCdf, familyMean, familyPdf, familySd, familyTail } from '../components/visual/math/families'

const close = (a: number, b: number, tolerance = 1e-9) => expect(Math.abs(a - b)).toBeLessThan(tolerance)

/** A Riemann sum of the density, which must be 1 for every family. */
function mass(family: (typeof FAMILIES)[number], lo: number, hi: number, steps = 20000): number {
  let sum = 0
  const dx = (hi - lo) / steps
  for (let i = 0; i < steps; i++) sum += familyPdf(family, lo + (i + 0.5) * dx) * dx
  return sum
}

describe('distribution families', () => {
  it('integrates to one, including the two with bounded or half-line support', () => {
    // A midpoint rule, so the tolerance is sampling error and not the density's.
    close(mass('normal', -8, 8), 1, 1e-4)
    close(mass('uniform', -2, 2), 1, 1e-4)
    close(mass('exponential', 0, 30), 1, 1e-4)
  })

  it('keeps the uniform density flat inside its support and zero outside', () => {
    const inside = familyPdf('uniform', 0.5)
    close(inside, 1 / (2 * Math.sqrt(3)))
    close(familyPdf('uniform', 1.7), inside)
    expect(familyPdf('uniform', 1.8)).toBe(0)
    expect(familyPdf('uniform', -1.8)).toBe(0)
  })

  it('gives the exponential no mass below zero', () => {
    expect(familyPdf('exponential', -0.001)).toBe(0)
    close(familyPdf('exponential', 0), 1)
    close(familyCdf('exponential', 0), 0)
  })

  it('matches known cumulative values', () => {
    close(familyCdf('normal', 0), 0.5, 1e-9)
    close(familyCdf('uniform', 0), 0.5)
    close(familyCdf('uniform', -10), 0)
    close(familyCdf('uniform', 10), 1)
    close(familyCdf('exponential', 1), 1 - Math.exp(-1))
    close(familyCdf('exponential', 5), 1 - Math.exp(-5))
  })

  it('complements the cumulative function exactly', () => {
    for (const family of FAMILIES) {
      for (const x of [-2, -0.5, 0, 0.5, 1, 3]) close(familyCdf(family, x) + familyTail(family, x), 1)
    }
  })

  it('is monotone in x for every family', () => {
    for (const family of FAMILIES) {
      let previous = -1
      for (let x = -4; x <= 5; x += 0.1) {
        const value = familyCdf(family, x)
        expect(value).toBeGreaterThanOrEqual(previous - 1e-12)
        expect(value).toBeLessThanOrEqual(1 + 1e-12)
        previous = value
      }
    }
  })

  it('shares one standard deviation, which is what makes the comparison fair', () => {
    close(familySd('normal'), 1)
    close(familySd('uniform'), 1)
    close(familySd('exponential'), 1)
    close(familyMean('normal'), 0)
    close(familyMean('uniform'), 0)
    close(familyMean('exponential'), 1)
  })

  it('gives the same interval different probabilities, so mean and sd do not name a distribution', () => {
    const atOne = FAMILIES.map((family) => familyCdf(family, 1))
    expect(new Set(atOne.map((value) => value.toFixed(4))).size).toBe(3)
    // Normal accumulates fastest near the centre, so it holds the most mass below 1.
    expect(familyCdf('normal', 1)).toBeGreaterThan(familyCdf('uniform', 1))
    expect(familyCdf('uniform', 1)).toBeGreaterThan(familyCdf('exponential', 1))
  })
})

describe('distribution family boundaries', () => {
  it('complements the cumulative probability everywhere, including the ends', () => {
    for (const family of ['normal', 'uniform', 'exponential'] as const) {
      for (const x of [-4.5, -1, 0, 1, 4, 5.5]) {
        close(familyTail(family, x), 1 - familyCdf(family, x))
      }
    }
  })

  it('puts the exponential support at zero and nothing below it', () => {
    expect(familyCdf('exponential', -0.001)).toBe(0)
    expect(familyPdf('exponential', -0.001)).toBe(0)
    close(familyCdf('exponential', 0), 0)
    close(familyTail('exponential', 0), 1)
    close(familyPdf('exponential', 0), 1)
    close(familyTail('exponential', 3), Math.exp(-3))
  })

  it('reaches certainty at the window edges for the bounded family', () => {
    close(familyCdf('uniform', -Math.sqrt(3)), 0)
    close(familyCdf('uniform', Math.sqrt(3)), 1)
    expect(familyPdf('uniform', -Math.sqrt(3) - 0.01)).toBe(0)
  })

  it('orders the three probabilities at the same point, and names the same spread', () => {
    const normal = familyCdf('normal', 1)
    const uniform = familyCdf('uniform', 1)
    const exponential = familyCdf('exponential', 1)
    expect(normal).toBeGreaterThan(uniform)
    expect(uniform).toBeGreaterThan(exponential)
    close(familyMean('exponential'), 1)
    close(familyMean('normal'), 0)
    close(familyMean('uniform'), 0)
    for (const family of ['normal', 'uniform', 'exponential'] as const) close(familySd(family), 1)
  })
})
