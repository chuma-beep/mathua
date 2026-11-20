import { describe, expect, it } from 'vitest'
import { readFileSync } from 'node:fs'
import { join } from 'node:path'
import { linearProgram, type LpShiftWhich } from '../components/visual/math/optimization'
import { optimizationSchema } from '../components/visual/schemas'

const REPO_ROOT = join(__dirname, '../../..')
const shard = (id: string): unknown =>
  JSON.parse(readFileSync(join(REPO_ROOT, 'data', 'lessons', 'kp', `${id}.json`), 'utf8'))

interface Step {
  label?: string
  visual?: {
    type: string
    title?: string
    caption?: string
    params?: Record<string, number | string>
    predict?: { question: string; options: string[]; answer: number; explain: string }
  }
}

const steps = (id: string) => shard(id) as Step[]

const cases: Array<{
  id: string
  correct: string
  active: number
  rhs?: LpShiftWhich
}> = [
  { id: 'optim.feasible.set', correct: 'Yes, both hold', active: 0 },
  { id: 'optim.active_constraint', correct: 'x + 2y ≤ 4 and 2x + y ≤ 5', active: 2 },
  { id: 'optim.active_set', correct: 'Two', active: 2 },
  { id: 'optim.kkt', correct: 'Both are positive', active: 2 },
  { id: 'optim.lp_dual', correct: '8', active: 2 },
  { id: 'optim.shadow_price', correct: 'About 0.067', active: 2, rhs: 'first' },
  { id: 'optim.lp_sensitivity', correct: 'From 2.5 to 10', active: 2 },
]

describe('the constrained-optimization figure serves its real lessons', () => {
  for (const item of cases) {
    it(`${item.id} draws the geometry its prediction describes`, () => {
      const list = steps(item.id)
      expect(list[0].visual?.type).toBe('optim.path')
      expect(list.filter((step) => step.visual)).toHaveLength(1)

      const params = list[0].visual!.params!
      expect(params.mode).toBe('feasible')
      const parsed = optimizationSchema.parse(params)
      const program = linearProgram(parsed.candidateX, parsed.candidateY, {
        which: parsed.rhs,
        delta: parsed.delta,
      })

      const expectedValue = item.rhs === 'first' ? 8 + parsed.delta / 3 : 8
      expect(program.status).toBe('optimal')
      expect(program.optimum.value).toBeCloseTo(expectedValue, 9)
      expect(program.active).toHaveLength(item.active)
      expect(program.optimumKkt.kkt).toBe(true)
      expect(program.optimumKkt.multipliers.every((multiplier) => multiplier.value >= -1e-9)).toBe(true)
      expect(program.dualValue).toBeCloseTo(expectedValue, 9)
      expect(program.gap).toBeCloseTo(0, 9)

      if (item.rhs === 'first') {
        expect(program.shadow?.value).toBeCloseTo(1 / 3, 9)
        expect(program.allowable?.min).toBeCloseTo(2.5, 9)
        expect(program.allowable?.max).toBeCloseTo(10, 9)
      } else {
        expect(program.shadow).toBeNull()
        expect(program.allowable).toBeNull()
      }

      const predict = list[0].visual!.predict!
      expect(predict.options.length).toBeGreaterThanOrEqual(2)
      expect(predict.options.length).toBeLessThanOrEqual(4)
      expect(predict.answer).toBeGreaterThanOrEqual(0)
      expect(predict.answer).toBeLessThan(predict.options.length)
      expect(predict.options[predict.answer]).toBe(item.correct)
      expect(predict.explain.length).toBeGreaterThan(40)
    })
  }

  it('keeps every prediction specific to its own idea', () => {
    const questions = new Set<string>()
    for (const item of cases) {
      const predict = steps(item.id)[0].visual!.predict!
      expect(questions.has(predict.question)).toBe(false)
      questions.add(predict.question)
    }
    expect(questions.size).toBe(cases.length)
  })

  it('varies the correct answer position across the seven lessons', () => {
    const positions = cases.map((item) => steps(item.id)[0].visual!.predict!.answer)
    expect(new Set(positions).size).toBeGreaterThan(1)
    expect(Math.max(...positions)).toBeLessThan(4)
  })

  it('records the template as supplementary without inflating the proposal', () => {
    const coverage = JSON.parse(
      readFileSync(join(REPO_ROOT, 'data', 'visualizations', 'coverage.json'), 'utf8'),
    ) as {
      topics: Array<{ id: string }>
      supplementary: Array<{ id: string; reason: string }>
    }
    expect(coverage.topics).toHaveLength(31)
    const entry = coverage.supplementary.find((item) => item.id === 'optim.path')
    expect(entry).toBeDefined()
    expect(entry!.reason).toContain('active')
    expect(entry!.reason.toLowerCase()).toContain('dual')
  })
})
