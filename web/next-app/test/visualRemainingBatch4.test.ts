import { describe, expect, it } from 'vitest'
import { balanceState, compareSigned, percentChange, ratioForm, tenthsValue } from '../components/visual/math/arith'
import { congruence, divisors, euclid, factorize, modPow, modValue } from '../components/visual/math/modular'
import { elementOrder, generatedSubgroup, groupById, inverseOf, isAbelian } from '../components/visual/math/group'
import { binomial, combinations, maskOf, setOperation, truthTable } from '../components/visual/math/discrete'
import { gradient, gradientDescent, linearProgram, objective } from '../components/visual/math/optimization'
import { metricBall, metricDistance, partialSums, sequenceTerm, sequenceView } from '../components/visual/math/analysis'
import { logisticProbability, mse, mseGradient, sigmoid, trainRegression } from '../components/visual/math/ml'

describe('arithmetic models', () => {
  it('orders signed numbers', () => {
    expect(compareSigned(-4, 2).order).toBe(-1)
    expect(compareSigned(2, -4).order).toBe(1)
  })
  it('reads tenths, percentages and ratios', () => {
    expect(tenthsValue(7)).toBe(0.7)
    expect(percentChange(25, 80)).toMatchObject({ amount: 20, result: 100 })
    expect(ratioForm(6, 4)).toMatchObject({ simplifiedA: 3, simplifiedB: 2 })
  })
  it('balances a one-step equation', () => {
    expect(balanceState(2, 3, 8)).toMatchObject({ solution: 5 })
    expect(balanceState(5, 3, 8).balanced).toBe(true)
  })
})

describe('number theory', () => {
  it('computes residues and powers', () => {
    expect(modValue(-1, 5)).toBe(4)
    expect(congruence(17, 5)).toEqual({ quotient: 3, residue: 2 })
    expect(modPow(2, 10, 5)?.value).toBe(4)
  })
  it('runs Euclid and factorises', () => {
    expect(euclid(30, 17)?.gcd).toBe(1)
    expect(euclid(48, 18)?.gcd).toBe(6)
    expect(factorize(60)?.factors).toEqual([2, 2, 3, 5])
    expect(divisors(12)).toEqual([1, 2, 3, 4, 6, 12])
  })
  it('rejects invalid input', () => {
    expect(euclid(0, 0)).toBeNull()
    expect(factorize(1)).toBeNull()
    expect(modPow(2, -1, 5)).toBeNull()
  })
})

describe('finite groups', () => {
  it('has the expected Z4 structure', () => {
    const g = groupById('z4')
    expect(isAbelian(g)).toBe(true)
    expect(elementOrder(g, 1)).toBe(4)
    expect(inverseOf(g, 1)).toBe(3)
    expect(generatedSubgroup(g, 1)).toEqual([0, 1, 2, 3])
  })
  it('has six non-abelian S3 elements', () => {
    const g = groupById('s3')
    expect(g.elements).toHaveLength(6)
    expect(isAbelian(g)).toBe(false)
    expect(elementOrder(g, 1)).toBe(2)
  })
})

describe('discrete mathematics', () => {
  it('evaluates every truth row', () => {
    const table = truthTable('implies')
    expect(table.trueCount).toBe(3)
    expect(table.rows.find((row) => row.p && !row.q)?.value).toBe(false)
  })
  it('computes sets from masks', () => {
    const result = setOperation('intersection', maskOf([1, 2, 3]), maskOf([3, 4, 5]))
    expect(result.result).toEqual([3])
    expect(result.count).toBe(1)
  })
  it('counts and lists combinations', () => {
    expect(binomial(5, 2)).toBe(10)
    expect(combinations(5, 2)).toHaveLength(10)
    expect(binomial(3, 5)).toBe(0)
  })
})

describe('optimization', () => {
  it('computes exact gradients and decreases the objective', () => {
    expect(objective(2, 1)).toBe(6)
    expect(gradient(2, 1)).toEqual([4, 4])
    const run = gradientDescent(2.5, 2, 0.15, 4)
    expect(run.points[0].value).toBeGreaterThan(run.points[run.points.length - 1].value)
  })
  it('checks the feasible region and vertex optimum', () => {
    const lp = linearProgram(1, 1)
    expect(lp.feasible).toBe(true)
    expect(lp.optimum).toMatchObject({ x: 2, y: 1, value: 8 })
    expect(linearProgram(3, 3).feasible).toBe(false)
  })
})

describe('analysis and metrics', () => {
  it('finds an epsilon threshold', () => {
    expect(sequenceTerm('reciprocal', 4)).toBe(0.25)
    const view = sequenceView('reciprocal', 0.2, 30)
    expect(view.threshold).toBe(6)
    expect(view.within).toBe(true)
  })
  it('builds partial sums', () => {
    expect(partialSums('reciprocal', 4)).toEqual([1, 1.5, 1.8333333333333333, 2.083333333333333])
  })
  it('changes the ball with the metric', () => {
    expect(metricDistance('manhattan', 1, 1)).toBe(2)
    expect(metricDistance('chebyshev', 1, 1)).toBe(1)
    expect(metricBall('chebyshev', 2).pointsInside).toBeGreaterThan(metricBall('euclidean', 2).pointsInside)
  })
})

describe('machine learning', () => {
  it('computes mean squared error and its gradient', () => {
    expect(mse(0, 0)).toBeGreaterThan(0)
    const [ds, db] = mseGradient(0, 0)
    expect(ds).toBeLessThan(0)
    expect(db).toBeLessThan(0)
  })
  it('reduces the loss with exact gradient steps', () => {
    const run = trainRegression(0.1, 8)
    expect(run.steps[0].loss).toBeGreaterThan(run.steps[run.steps.length - 1].loss)
  })
  it('maps scores through the sigmoid', () => {
    expect(sigmoid(0)).toBe(0.5)
    expect(logisticProbability(2, 0, 1)).toBeCloseTo(0.8808, 3)
  })
})
