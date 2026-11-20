import { describe, expect, it } from 'vitest'
import { probabilityTree } from '../components/visual/math/probability'
import {
  addFractions,
  compareFractions,
  equivalentFractions,
  fractionValue,
  gcd,
  lcm,
  simplify,
} from '../components/visual/math/fraction'

const close = (a: number, b: number, tolerance = 1e-9) => expect(Math.abs(a - b)).toBeLessThan(tolerance)

describe('probability trees', () => {
  it('builds joint leaves that sum to one', () => {
    const tree = probabilityTree(0.01, 0.95, 0.9)
    expect(tree).not.toBeNull()
    if (!tree) return
    close(tree.positive, 0.1085)
    close(tree.negative, 0.8915)
    close(tree.positive + tree.negative, 1)
    close(tree.diseasePositive + tree.diseaseNegative + tree.healthyPositive + tree.healthyNegative, 1)
    close(tree.diseaseGivenPositive, 0.0095 / 0.1085)
  })

  it('makes a zero prior stay at zero after a positive test', () => {
    const tree = probabilityTree(0, 0.95, 0.9)
    expect(tree).not.toBeNull()
    if (!tree) return
    close(tree.diseaseGivenPositive, 0)
    close(tree.positive, 0.1)
  })

  it('makes a certain disease certain after a positive test', () => {
    const tree = probabilityTree(1, 0.95, 0.9)
    expect(tree).not.toBeNull()
    if (!tree) return
    close(tree.diseaseGivenPositive, 1)
    close(tree.positive, 0.95)
  })

  it('clamps probabilities and rejects non-finite input', () => {
    const clamped = probabilityTree(-1, 2, 0.5)
    expect(clamped).not.toBeNull()
    if (!clamped) return
    close(clamped.prior, 0)
    close(clamped.sensitivity, 1)
    expect(probabilityTree(Number.NaN, 0.9, 0.9)).toBeNull()
  })
})

describe('fractions', () => {
  it('reduces with the greatest common divisor', () => {
    expect(gcd(12, 18)).toBe(6)
    expect(gcd(0, 5)).toBe(5)
    expect(simplify({ numerator: 6, denominator: 8 })).toEqual({ numerator: 3, denominator: 4 })
    expect(simplify({ numerator: 2, denominator: -4 })).toEqual({ numerator: -1, denominator: 2 })
  })

  it('finds a common denominator and adds exactly', () => {
    expect(lcm(4, 6)).toBe(12)
    expect(addFractions({ numerator: 1, denominator: 2 }, { numerator: 1, denominator: 3 })).toEqual({
      numerator: 5,
      denominator: 6,
    })
  })

  it('compares by cross products, not decimals', () => {
    expect(compareFractions({ numerator: 1, denominator: 3 }, { numerator: 2, denominator: 6 })).toBe(0)
    expect(compareFractions({ numerator: 3, denominator: 4 }, { numerator: 1, denominator: 2 })).toBe(1)
    expect(compareFractions({ numerator: 1, denominator: 4 }, { numerator: 1, denominator: 2 })).toBe(-1)
    close(fractionValue({ numerator: 3, denominator: 4 }), 0.75)
  })

  it('scales both parts to keep the same value', () => {
    expect(equivalentFractions({ numerator: 1, denominator: 2 }, 3)).toEqual({ numerator: 3, denominator: 6 })
    expect(compareFractions({ numerator: 1, denominator: 2 }, equivalentFractions({ numerator: 1, denominator: 2 }, 3))).toBe(0)
  })
})
