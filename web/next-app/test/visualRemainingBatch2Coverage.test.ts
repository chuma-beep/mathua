import { describe, expect, it } from 'vitest'
import { readFileSync } from 'node:fs'
import { join } from 'node:path'
import { probabilityTree } from '../components/visual/math/probability'
import { compareFractions, fractionValue, simplify, addFractions } from '../components/visual/math/fraction'
import { fractionLineSchema, probabilityTreeSchema } from '../components/visual/schemas'

const REPO_ROOT = join(__dirname, '../../..')
const read = (path: string): unknown =>
  JSON.parse(readFileSync(join(REPO_ROOT, path), 'utf8'))

interface Prediction {
  question: string
  options: string[]
  answer: number
  explain: string
}

interface Visual {
  type: string
  title?: string
  caption?: string
  params?: Record<string, unknown>
  predict?: Prediction
}

interface Step {
  label?: string
  subgoals?: string[]
  visual?: Visual
}

const steps = (id: string) => read(`data/lessons/kp/${id}.json`) as Step[]
const attached = (id: string) => {
  const step = steps(id)[0]
  if (!step.visual) throw new Error(`${id} has no visual on step 0`)
  return step.visual
}

const close = (value: number, expected: number, tolerance = 1e-9) =>
  expect(Math.abs(value - expected)).toBeLessThan(tolerance)

describe('probability and fraction visuals are attached to their lessons', () => {
  it('attaches the tree to the conditional-definition step', () => {
    const visual = attached('prob.cond.definition')
    expect(visual.type).toBe('prob.tree')
    expect(visual.title).toBe('The conditioning event becomes the denominator')
    expect(steps('prob.cond.definition')[0].label).toBe('Conditioning on an event')
    expect(visual.predict?.answer).toBe(2)
    expect(visual.predict?.options[2]).toBe('P(D | +), the posterior')
  })

  it('computes the tree values those predictions name', () => {
    const visual = attached('prob.cond.definition')
    const params = probabilityTreeSchema.parse(visual.params ?? {})
    const tree = probabilityTree(params.prior, params.sensitivity, params.specificity)
    if (!tree) throw new Error('tree is null')
    close(tree.diseasePositive, 0.0095)
    close(tree.healthyPositive, 0.099)
    close(tree.positive, 0.1085)
    close(tree.diseaseGivenPositive, 0.0095 / 0.1085)
    close(tree.diseasePositive + tree.diseaseNegative + tree.healthyPositive + tree.healthyNegative, 1)
    expect(visual.predict?.explain).toContain('positive leaves')
  })

  it('uses total probability for the second probability attachment', () => {
    const visual = attached('prob.cond.total_probability')
    expect(visual.type).toBe('prob.tree')
    expect(visual.title).toBe('All positive leaves make P(+)')
    expect(visual.predict?.answer).toBe(0)
    expect(visual.predict?.options[0]).toContain('P(D)P(+|D) + P(¬D)P(+|¬D)')
    expect(visual.predict?.explain).toContain('P(+) = 0.1085')
  })

  it('uses the same tree to show why the posterior is small', () => {
    const visual = attached('prob.bayes.prior')
    expect(visual.type).toBe('prob.tree')
    expect(visual.predict?.answer).toBe(0)
    expect(visual.predict?.options[0]).toBe('8.8%')
    expect(visual.predict?.explain).toContain('about 0.0876')
  })

  it('attaches the fraction line to locating a fraction', () => {
    const visual = attached('frac.basics.number_line')
    expect(visual.type).toBe('frac.line')
    expect(visual.title).toBe('One point, counted from zero')
    expect(visual.predict?.answer).toBe(0)
    expect(visual.predict?.options[0]).toBe('3/4')
  })

  it('agrees with the fraction comparison it predicts', () => {
    const visual = attached('frac.basics.number_line')
    const params = fractionLineSchema.parse(visual.params ?? {})
    const a = { numerator: params.numerator, denominator: params.denominator }
    const b = { numerator: params.otherNumerator, denominator: params.otherDenominator }
    expect(compareFractions(a, b)).toBe(1)
    expect(fractionValue(a)).toBeGreaterThan(fractionValue(b))
    expect(visual.predict?.explain).toContain('3/4 is 0.75')
  })

  it('uses equivalent fractions on the equivalence step', () => {
    const visual = attached('frac.basics.equivalent')
    expect(visual.type).toBe('frac.line')
    expect(visual.title).toBe('More parts, same point')
    const params = fractionLineSchema.parse(visual.params ?? {})
    const a = { numerator: params.numerator, denominator: params.denominator }
    const b = { numerator: params.otherNumerator, denominator: params.otherDenominator }
    expect(compareFractions(a, b)).toBe(0)
    expect(simplify(a)).toEqual({ numerator: 1, denominator: 2 })
    expect(simplify(b)).toEqual({ numerator: 1, denominator: 2 })
    expect(visual.predict?.answer).toBe(0)
    expect(visual.predict?.options[0]).toBe('Yes, they name the same amount')
  })

  it('keeps the arithmetic exact for the displayed sum', () => {
    const first = { numerator: 3, denominator: 4 }
    const second = { numerator: 1, denominator: 2 }
    expect(addFractions(first, second)).toEqual({ numerator: 5, denominator: 4 })
  })
})
