import { readFileSync } from 'node:fs'
import { resolve } from 'node:path'
import { describe, expect, it } from 'vitest'
import { functionGraphSchema } from '../components/visual/schemas'
import {
  endBehavior,
  functionFormula,
  inverseFormula,
  inverseIsFunction,
} from '../components/visual/math/functionGraph'

interface Step {
  label: string
  visual?: {
    type: string
    title: string
    caption: string
    params: Record<string, unknown>
    predict: { question: string; options: string[]; answer: number; explain: string }
  }
}

const shard = (concept: string): Step[] =>
  JSON.parse(readFileSync(resolve(process.cwd(), `../../data/lessons/kp/${concept}.json`), 'utf8')) as Step[]

describe('function graph coverage', () => {
  it('attaches the polynomial figure on step 0', () => {
    const steps = shard('precalc.poly.graph')
    expect(steps[0].visual?.type).toBe('function.graph')
    expect(steps[0].label).toBe('End behavior')
    expect(steps.slice(1).every((step) => !step.visual)).toBe(true)
    const params = functionGraphSchema.parse(steps[0].visual?.params)
    expect(params).toMatchObject({ family: 'cubic', a: -1, h: 0, k: 0, showInverse: false })
    expect(functionFormula(params.family, params.a, params.h, params.k)).toBe('−x³')
    expect(endBehavior(params.family, params.a)).toEqual({ left: '+∞', right: '−∞' })

    const predict = steps[0].visual?.predict
    expect(predict?.question).toBe('For y = −x³, what happens at the far left and far right?')
    expect(predict?.options).toEqual([
      'Both ends go to +∞',
      'Both ends go to −∞',
      'The left end goes to +∞ and the right end to −∞',
    ])
    expect(predict?.answer).toBe(2)
    expect(predict?.explain).toContain('odd degree')
    expect(predict?.explain).toContain('x → −∞ gives y → +∞')
  })

  it('attaches the inverse figure on step 0', () => {
    const steps = shard('precalc.func.inverse_graph')
    expect(steps[0].visual?.type).toBe('function.graph')
    expect(steps[0].label).toBe('Swapping coordinates')
    const params = functionGraphSchema.parse(steps[0].visual?.params)
    expect(params).toMatchObject({ family: 'exponential', a: 1, h: 0, k: 0, showInverse: true })
    expect(functionFormula(params.family, params.a, params.h, params.k)).toBe('2^x')
    expect(inverseIsFunction(params.family, params.a)).toBe(true)
    expect(inverseFormula(params.family, params.a, params.h, params.k)).toBe('0 + log₂(y)')

    const predict = steps[0].visual?.predict
    expect(predict?.question).toBe('The graph of y = 2ˣ is reflected in y = x. Which equation does the reflection follow?')
    expect(predict?.options).toEqual(['y = x²', 'y = log₂ x', 'y = 2^(−x)'])
    expect(predict?.answer).toBe(1)
    expect(predict?.explain).toContain('swaps x and y')
    expect(predict?.explain).toContain('y = log₂ x')
  })

  it('keeps the coverage ledger at 31 topics and records the template', () => {
    const coverage = JSON.parse(
      readFileSync(resolve(process.cwd(), '../../data/visualizations/coverage.json'), 'utf8'),
    ) as { topics: unknown[]; supplementary: Array<{ id: string; reason: string }> }
    expect(coverage.topics).toHaveLength(31)
    const entry = coverage.supplementary.find((item) => item.id === 'function.graph')
    expect(entry).toBeDefined()
    expect(entry?.reason).toContain('reflection')
  })
})
