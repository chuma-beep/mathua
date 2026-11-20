import { describe, expect, it } from 'vitest'
import { readFileSync } from 'node:fs'
import { join } from 'node:path'
import { unitCircleSchema } from '../components/visual/schemas'
import { degreesToRadians, radiansToDegrees, unitPoint } from '../components/visual/math/trig'

const REPO_ROOT = join(__dirname, '../../..')
const shard = (id: string): unknown =>
  JSON.parse(readFileSync(join(REPO_ROOT, 'data', 'lessons', 'kp', `${id}.json`), 'utf8'))

interface Step {
  label?: string
  subgoals?: string[]
  visual?: {
    type: string
    title?: string
    caption?: string
    params?: Record<string, unknown>
    predict?: { question: string; options: string[]; answer: number; explain: string }
  }
}

const steps = (id: string) => shard(id) as Step[]
const attached = (id: string) => {
  const step = steps(id).find((item) => item.visual)
  if (!step?.visual) throw new Error(`${id} has no attached visual`)
  return step
}
const pointOf = (id: string) => {
  const params = unitCircleSchema.parse(attached(id).visual!.params)
  return unitPoint(degreesToRadians(params.angle))
}

describe('unit-circle shards agree with the trig engine', () => {
  const ids = ['trig.basics.unit_circle', 'trig.basics.sin_cos_def', 'precalc.trig.unit_circle']

  it('attaches trig.unitcircle to the first step only', () => {
    for (const id of ids) {
      const shardSteps = steps(id)
      expect(shardSteps[0].visual?.type).toBe('trig.unitcircle')
      expect(shardSteps.filter((step) => step.visual)).toHaveLength(1)
    }
  })

  it('parses every attached param set against the schema', () => {
    for (const id of ids) {
      const parsed = unitCircleSchema.parse(attached(id).visual!.params)
      expect(parsed.showSin).toBe(true)
      expect(parsed.showCos).toBe(true)
      expect(Number.isFinite(parsed.angle)).toBe(true)
    }
  })

  it('starts the definition figure at 45 degrees and predicts the 90 degree point', () => {
    const point = pointOf('trig.basics.unit_circle')
    expect(point.cos).toBeCloseTo(Math.SQRT1_2, 12)
    expect(point.sin).toBeCloseTo(Math.SQRT1_2, 12)
    expect(point.quadrant).toBe('I')

    const atNinety = unitPoint(degreesToRadians(90))
    expect(atNinety.point[0]).toBeCloseTo(0, 12)
    expect(atNinety.point[1]).toBeCloseTo(1, 12)
    expect(atNinety.cos).toBeCloseTo(0, 12)
    expect(atNinety.sin).toBeCloseTo(1, 12)

    const predict = attached('trig.basics.unit_circle').visual!.predict!
    expect(predict.question).toBe('At θ = 90°, where is P and what are its coordinates?')
    expect(predict.options).toEqual([
      'P is (0, 1), so cos θ = 0 and sin θ = 1',
      'P is (1, 0), so cos θ = 1 and sin θ = 0',
      'P leaves the unit circle because the angle is vertical',
    ])
    expect(predict.answer).toBe(0)
    expect(predict.explain).toContain('x-coordinate is 0')
    expect(predict.explain).toContain('y-coordinate is 1')
  })

  it('makes the height sine and the width cosine at 30 degrees', () => {
    const point = pointOf('trig.basics.sin_cos_def')
    expect(point.cos).toBeCloseTo(Math.sqrt(3) / 2, 12)
    expect(point.sin).toBeCloseTo(0.5, 12)
    expect(point.point[0]).toBeCloseTo(Math.sqrt(3) / 2, 12)
    expect(point.point[1]).toBeCloseTo(0.5, 12)

    const predict = attached('trig.basics.sin_cos_def').visual!.predict!
    expect(predict.options).toEqual([
      'The x-coordinate, so cos 30° = 1/2',
      'The y-coordinate, so sin 30° = 1/2',
      'Neither coordinate is 1/2',
    ])
    expect(predict.answer).toBe(1)
    expect(predict.explain).toContain('(√3/2, 1/2)')
  })

  it('reads the quadrant III signs and reference angle at 210 degrees', () => {
    const point = pointOf('precalc.trig.unit_circle')
    expect(point.cos).toBeCloseTo(-Math.sqrt(3) / 2, 12)
    expect(point.sin).toBeCloseTo(-0.5, 12)
    expect(point.quadrant).toBe('III')
    expect(radiansToDegrees(point.referenceAngle)).toBeCloseTo(30, 9)

    const predict = attached('precalc.trig.unit_circle').visual!.predict!
    expect(predict.options).toEqual([
      'cos θ < 0 and sin θ < 0',
      'cos θ > 0 and sin θ < 0',
      'cos θ < 0 and sin θ > 0',
    ])
    expect(predict.answer).toBe(0)
    expect(predict.explain).toContain('quadrant III')
    expect(predict.explain).toContain('reference angle is 30°')
  })

  it('describes the linked circle and graph honestly in each caption', () => {
    for (const id of ids) {
      const visual = attached(id).visual!
      expect(visual.title?.length).toBeGreaterThan(5)
      expect(visual.caption?.length).toBeGreaterThan(60)
      expect(visual.caption?.toLowerCase()).toContain('unit circle')
      expect(visual.caption?.toLowerCase()).toContain('graph')
    }
  })
})

describe('coverage ledger records the supplementary trig template', () => {
  const coverage = JSON.parse(
    readFileSync(join(REPO_ROOT, 'data', 'visualizations', 'coverage.json'), 'utf8'),
  ) as {
    topics: Array<{ id: string; status: string; templates: string[] }>
    supplementary: Array<{ id: string; reason: string }>
  }

  it('keeps the 31-topic proposal and adds trig.unitcircle as supplementary', () => {
    expect(coverage.topics).toHaveLength(31)
    const entry = coverage.supplementary.find((item) => item.id === 'trig.unitcircle')
    expect(entry).toBeTruthy()
    expect(entry!.reason).toContain('unit-circle')
    expect(entry!.reason).toContain('sine and cosine')
  })
})
