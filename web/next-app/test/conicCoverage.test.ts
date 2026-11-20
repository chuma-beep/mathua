import { describe, expect, it } from 'vitest'
import { readFileSync } from 'node:fs'
import { join } from 'node:path'
import { classifyGeneralConic } from '../components/visual/math/conic'
import { conicSchema } from '../components/visual/schemas'

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
    params?: Record<string, number>
    predict?: { question: string; options: string[]; answer: number; explain: string }
  }
}

const steps = (id: string) => shard(id) as Step[]
const attached = (id: string) => {
  const step = steps(id).find((item) => item.visual)
  if (!step?.visual) throw new Error(`${id} has no attached visual`)
  return step
}

describe('alg.conic.classify shard agrees with the classifier', () => {
  it('attaches the general conic to the first step', () => {
    const shardSteps = steps('alg.conic.classify')
    expect(shardSteps[0].visual?.type).toBe('geo.conic')
    expect(shardSteps[0].visual?.title).toBe('One equation, several curves')
  })

  it('classifies the attached params as the ellipse the caption claims', () => {
    const params = attached('alg.conic.classify').visual!.params!
    const parsed = conicSchema.parse(params)
    // The shard still writes three coefficients; F defaults to -1.
    expect(parsed).toMatchObject({ A: 1, B: 0, C: 0.5, D: 0, E: 0, F: -1 })
    const geometry = classifyGeneralConic(parsed.A, parsed.B, parsed.C, parsed.D, parsed.E, parsed.F)
    expect(geometry.kind).toBe('ellipse')
    expect(geometry.center).toEqual([0, 0])
    expect(geometry.a).toBeCloseTo(Math.SQRT2, 9)
    expect(geometry.b).toBeCloseTo(1, 9)
  })

  it('is parallel lines at B = √2 and a hyperbola above it', () => {
    const params = attached('alg.conic.classify').visual!.params!
    const parsed = conicSchema.parse(params)

    const boundary = classifyGeneralConic(parsed.A, Math.SQRT2, parsed.C, parsed.D, parsed.E, parsed.F)
    expect(boundary.kind).toBe('parallel')
    expect(boundary.lines).toHaveLength(2)

    const past = classifyGeneralConic(parsed.A, 2, parsed.C, parsed.D, parsed.E, parsed.F)
    expect(past.kind).toBe('hyperbola')
  })

  it('states the corrected prediction, options, answer and explanation', () => {
    const predict = attached('alg.conic.classify').visual!.predict!
    expect(predict.question).toBe(
      'With A = 1 and C = 0.5, raise the cross term B from 0. What happens when B passes √2?',
    )
    expect(predict.options).toEqual([
      'The ellipse becomes a hyperbola',
      'The ellipse only rotates',
      'The ellipse becomes a parabola',
    ])
    expect(predict.answer).toBe(0)
    expect(predict.options[predict.answer]).toBe('The ellipse becomes a hyperbola')
    expect(predict.explain).toContain('B² − 4AC')
    expect(predict.explain).toContain('2√(AC) = √2')
    expect(predict.explain).toContain('two parallel lines')
    expect(predict.explain).toContain('A parabola needs a linear term')
    expect(predict.explain).not.toContain('only rotates')
  })

  it('replaces the false B-only-rotates caption and subgoal', () => {
    const step = attached('alg.conic.classify')
    expect(step.visual!.caption).toContain('A linear term D or E is what makes a parabola possible')
    expect(step.visual!.caption).not.toContain('without changing kind')
    expect(step.subgoals).toContain('A linear term is what makes a parabola possible')
  })
})

describe('precalculus conic shards no longer use a11/a22 heuristics', () => {
  const ids = [
    'precalc.conic.ellipse',
    'precalc.conic.hyperbola',
    'precalc.conic.parabola',
    'precalc.conic.general',
    'precalc.conic.rotated',
  ]

  it('does not use a11/a22 to decide the kind', () => {
    for (const id of ids) {
      const text = JSON.stringify(shard(id))
      // The old classification step read "check a11 and a22 are both nonzero".
      // A later step may still name a11/a22 as matrix entries that carry shape;
      // that is not a kind test.
      expect(text).not.toMatch(/check a11 and a22/i)
      expect(text).not.toMatch(/a11 and a22 are (both )?nonzero/i)
      const first = steps(id)[0]
      expect(JSON.stringify(first.subgoals)).not.toContain('a11')
      expect(JSON.stringify(first.subgoals)).not.toContain('a22')
    }
  })

  it('names the degenerate cases the discriminant alone cannot separate', () => {
    expect(JSON.stringify(shard('precalc.conic.rotated'))).toContain('intersecting lines')
    expect(JSON.stringify(shard('precalc.conic.general'))).toContain('degenerate curve')
    expect(JSON.stringify(shard('precalc.conic.ellipse'))).toContain('discriminant is negative')
    expect(JSON.stringify(shard('precalc.conic.hyperbola'))).toContain('discriminant is positive')
    expect(JSON.stringify(shard('precalc.conic.parabola'))).toContain('discriminant is zero')
  })
})

describe('P2 conic shards agree with the focal geometry', () => {
  const geometryOf = (id: string) => {
    const params = attached(id).visual!.params!
    const parsed = conicSchema.parse(params)
    return classifyGeneralConic(parsed.A, parsed.B, parsed.C, parsed.D, parsed.E, parsed.F)
  }
  const pointClose = (actual: [number, number] | undefined, expected: [number, number]) => {
    expect(actual).toBeDefined()
    expect(actual![0]).toBeCloseTo(expected[0], 9)
    expect(actual![1]).toBeCloseTo(expected[1], 9)
  }

  it('attaches each P2 visual to the first step only', () => {
    for (const id of [
      'alg.conic.parabola',
      'alg.conic.parabola_focus',
      'alg.conic.ellipse_foci',
      'precalc.conic.eccentricity',
    ]) {
      const shardSteps = steps(id)
      expect(shardSteps[0].visual?.type).toBe('geo.conic')
      expect(shardSteps.filter((step) => step.visual)).toHaveLength(1)
    }
  })

  it('shows the parabola vertex, focus and directrix the lesson names', () => {
    const geometry = geometryOf('alg.conic.parabola')
    expect(geometry.kind).toBe('parabola')
    expect(geometry.p).toBeCloseTo(1, 9)
    pointClose(geometry.vertex, [0, 0])
    pointClose(geometry.foci?.[0], [0, 1])
    pointClose(geometry.directrix?.point, [0, -1])

    const predict = attached('alg.conic.parabola').visual!.predict!
    expect(predict.question).toContain('x² = 4y')
    expect(predict.options[predict.answer]).toBe('(0, 1)')
    expect(predict.explain).toContain('p = 1')
  })

  it('ties p to the focus and the latus rectum in the focus lesson', () => {
    const geometry = geometryOf('alg.conic.parabola_focus')
    expect(geometry.kind).toBe('parabola')
    expect(geometry.p).toBeCloseTo(1, 9)
    expect(4 * (geometry.p ?? 0)).toBeCloseTo(4, 9)
    pointClose(geometry.foci?.[0], [0, 1])
    pointClose(geometry.directrix?.point, [0, -1])

    const predict = attached('alg.conic.parabola_focus').visual!.predict!
    expect(predict.question).toContain('latus rectum')
    expect(predict.options[predict.answer]).toBe('4p')
    expect(predict.explain).toContain('4p')
  })

  it('shows the constant focal sum the ellipse lesson claims', () => {
    const geometry = geometryOf('alg.conic.ellipse_foci')
    expect(geometry.kind).toBe('ellipse')
    expect(geometry.a).toBeCloseTo(Math.SQRT2, 9)
    expect(geometry.b).toBeCloseTo(1, 9)
    expect(geometry.focalConstant).toBeCloseTo(2 * Math.SQRT2, 9)
    pointClose(geometry.foci?.[0], [0, -1])
    pointClose(geometry.foci?.[1], [0, 1])

    const predict = attached('alg.conic.ellipse_foci').visual!.predict!
    expect(predict.options[predict.answer]).toBe('It stays equal to 2a')
    expect(predict.explain).toContain('2√2')
  })

  it('keeps eccentricity a constant of the conic while focal distances change', () => {
    const geometry = geometryOf('precalc.conic.eccentricity')
    expect(geometry.kind).toBe('ellipse')
    expect(geometry.eccentricity).toBeCloseTo(Math.SQRT1_2, 9)
    expect(attached('precalc.conic.eccentricity').visual!.caption).toContain('e stays 0.707')
    expect(attached('precalc.conic.eccentricity').visual!.caption).toContain('PF1 and PF2')
    expect(steps('precalc.conic.eccentricity')[0].subgoals).toContain(
      'The value is the same for every point on the curve.',
    )

    const predict = attached('precalc.conic.eccentricity').visual!.predict!
    expect(predict.options[predict.answer]).toBe('It stays 0.707')
    expect(predict.explain).toContain('constant of the conic')
  })
})

describe('visualizations coverage ledger records P1 without overclaiming', () => {
  const coverage = JSON.parse(
    readFileSync(join(REPO_ROOT, 'data', 'visualizations', 'coverage.json'), 'utf8'),
  ) as { topics: Array<{ id: string; status: string; templates: string[]; note: string }> }

  it('keeps conic_sections partial and attached to geo.conic', () => {
    const topic = coverage.topics.find((item) => item.id === 'conic_sections')
    expect(topic).toBeTruthy()
    expect(topic!.status).toBe('partial')
    expect(topic!.templates).toContain('geo.conic')
    expect(topic!.note).toContain('parabolas')
    expect(topic!.note).toContain('Intersections')
  })
})
