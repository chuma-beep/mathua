import { describe, expect, it } from 'vitest'
import { readFileSync } from 'node:fs'
import { join } from 'node:path'
import { classifyGeneralConic, type ConicKind } from '../components/visual/math/conic'
import { conicSchema } from '../components/visual/schemas'

const REPO_ROOT = join(__dirname, '../../..')
const shard = (id: string): unknown =>
  JSON.parse(readFileSync(join(REPO_ROOT, 'data', 'lessons', 'kp', `${id}.json`), 'utf8'))

interface Step {
  label?: string
  visual?: {
    type: string
    title?: string
    caption?: string
    params?: Record<string, number>
    predict?: { question: string; options: string[]; answer: number; explain: string }
  }
}

const steps = (id: string) => shard(id) as Step[]

const cases: Array<{
  id: string
  kind: ConicKind
  center: [number, number]
  a?: number
  b?: number
  correct: string
}> = [
  { id: 'alg.conic.circle', kind: 'circle', center: [1, -2], a: 2, b: 2, correct: 'Centre (1, −2), radius 2' },
  { id: 'alg.conic.ellipse', kind: 'ellipse', center: [0, 0], a: 2, b: 1, correct: 'The horizontal axis, length 4' },
  { id: 'alg.conic.hyperbola', kind: 'hyperbola', center: [0, 0], a: 2, b: 3, correct: 'Left and right' },
  { id: 'alg.conic.translation', kind: 'ellipse', center: [2, -1], a: 2, b: 1, correct: '(2, −1)' },
  { id: 'precalc.conic.general', kind: 'ellipse', center: [0, 0], a: Math.SQRT2, b: Math.sqrt(2 / 3), correct: '−3, ellipse' },
  { id: 'precalc.conic.circle_radius', kind: 'circle', center: [2, -1], a: 3, b: 3, correct: '3' },
  { id: 'precalc.conic.ellipse', kind: 'ellipse', center: [1, -2], a: 2, b: 1, correct: '−1' },
  { id: 'precalc.conic.hyperbola', kind: 'hyperbola', center: [0, 0], a: 1, b: 2, correct: 'The constant −1 keeps the branches away from the centre' },
  { id: 'precalc.conic.rotated', kind: 'ellipse', center: [0, 0], a: Math.SQRT2, b: Math.sqrt(2 / 3), correct: 'It rotates the principal axes by 45°' },
]

describe('conic reuse pass attaches the existing classifier to real lessons', () => {
  for (const item of cases) {
    it(`${item.id} agrees with the classified geometry`, () => {
      const list = steps(item.id)
      expect(list[0].visual?.type).toBe('geo.conic')
      expect(list.filter((step) => step.visual)).toHaveLength(1)

      const params = list[0].visual!.params!
      const parsed = conicSchema.parse(params)
      const geometry = classifyGeneralConic(parsed.A, parsed.B, parsed.C, parsed.D, parsed.E, parsed.F)

      expect(geometry.kind).toBe(item.kind)
      expect(geometry.center?.[0]).toBeCloseTo(item.center[0], 9)
      expect(geometry.center?.[1]).toBeCloseTo(item.center[1], 9)
      if (item.a !== undefined) expect(geometry.a ?? NaN).toBeCloseTo(item.a, 9)
      if (item.b !== undefined) expect(geometry.b ?? NaN).toBeCloseTo(item.b, 9)

      const predict = list[0].visual!.predict!
      expect(predict.options.length).toBeGreaterThanOrEqual(2)
      expect(predict.options.length).toBeLessThanOrEqual(4)
      expect(predict.answer).toBeGreaterThanOrEqual(0)
      expect(predict.answer).toBeLessThan(predict.options.length)
      expect(predict.options[predict.answer]).toBe(item.correct)
      expect(predict.explain.length).toBeGreaterThan(30)
    })
  }

  it('varies the correct answer position across the new conic attachments', () => {
    const positions = steps(cases[0].id)
      .concat(...cases.slice(1).map((item) => steps(item.id)))
      .flatMap((step) => (step.visual?.predict ? [step.visual.predict.answer] : []))
    expect(new Set(positions).size).toBeGreaterThan(1)
  })

  it('keeps every new prediction specific to its own equation', () => {
    const questions = new Set<string>()
    for (const item of cases) {
      const predict = steps(item.id)[0].visual!.predict!
      expect(predict.question.length).toBeGreaterThan(20)
      expect(questions.has(predict.question)).toBe(false)
      questions.add(predict.question)
    }
  })

  it('does not change the 31-topic proposal ledger', () => {
    const coverage = JSON.parse(
      readFileSync(join(REPO_ROOT, 'data', 'visualizations', 'coverage.json'), 'utf8'),
    ) as { topics: unknown[] }
    expect(coverage.topics).toHaveLength(31)
  })
})
