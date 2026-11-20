import { describe, expect, it } from 'vitest'
import { readFileSync } from 'node:fs'
import { join } from 'node:path'
import { integrateSlopeField, slopeAt } from '../components/visual/math/ode'
import { bisection, newton, rootFunction } from '../components/visual/math/roots'
import { SAMPLE_GRAPH, runGraph } from '../components/visual/math/graph'
import { graphSchema, odeSlopeSchema, rootsSchema } from '../components/visual/schemas'

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

const checkPrediction = (visual: Visual, expected: string) => {
  const predict = visual.predict
  if (!predict) throw new Error(`${visual.title} has no prediction`)
  expect(predict.question.length).toBeGreaterThan(0)
  expect(predict.options.length).toBeGreaterThanOrEqual(2)
  expect(predict.options.length).toBeLessThanOrEqual(4)
  expect(predict.answer).toBeGreaterThanOrEqual(0)
  expect(predict.answer).toBeLessThan(predict.options.length)
  expect(predict.explain.length).toBeGreaterThan(0)
  expect(predict.options[predict.answer]).toBe(expected)
}

describe('differential-equation visuals are attached to their lessons', () => {
  it('attaches the slope field to the first reading step', () => {
    const visual = attached('ode.basics.slope_field')
    expect(visual.type).toBe('ode.slopefield')
    expect(visual.title).toBe('The field is the equation')
    expect(steps('ode.basics.slope_field')[0].label).toBe('Reading a slope field')
  })

  it('uses the real right-hand side and integrates eˣ through (0, 1)', () => {
    const params = odeSlopeSchema.parse(attached('ode.basics.slope_field').params)
    expect(params.field).toBe('exponential')
    expect(slopeAt(params.field, params.x0, params.y0)).toBe(1)
    const result = integrateSlopeField(params.field, params.x0, params.y0, {
      xMin: -1,
      xMax: 1,
      yMin: -1,
      yMax: 5,
      step: 0.01,
      method: 'rk4',
    })
    const end = result.points[result.points.length - 1]
    expect(end.x).toBeCloseTo(0.99, 6)
    expect(end.y).toBeCloseTo(Math.exp(0.99), 4)
    checkPrediction(attached('ode.basics.slope_field'), 'It grows')
  })

  it('attaches Euler to the worked one-step lesson with h = 0.25', () => {
    const visual = attached('ode.numerical.euler')
    expect(visual.type).toBe('ode.slopefield')
    expect(visual.title).toBe('Euler follows the tangent')
    const params = odeSlopeSchema.parse(visual.params)
    expect(params.method).toBe('euler')
    expect(params.step).toBe(0.25)
    const result = integrateSlopeField(params.field, params.x0, params.y0, {
      xMin: 0,
      xMax: 0.25,
      yMin: 0,
      yMax: 3,
      step: 0.25,
      method: 'euler',
    })
    const end = result.points[result.points.length - 1]
    expect(end.y).toBeCloseTo(1.25, 6)
    checkPrediction(visual, '1.25')
    expect(visual.predict?.explain).toContain('1.284')
  })
})

describe('root-iteration visuals are attached to their lessons', () => {
  it('attaches bisection and halves the bracket once', () => {
    const visual = attached('numa.roots.bisection')
    expect(visual.type).toBe('numa.roots')
    expect(visual.title).toBe('Each step halves the bracket')
    const params = rootsSchema.parse(visual.params)
    expect(params.method).toBe('bisection')
    const result = bisection(rootFunction(params.fn).f, params.a, params.b, 30)
    expect(result.steps[0].midpoint).toBe(1)
    expect(result.steps[0].keep).toBe('right')
    expect(result.root).toBeCloseTo(Math.SQRT2, 6)
    checkPrediction(visual, '[1, 2]')
    expect(visual.predict?.explain).toContain('f(1) = −1')
  })

  it('attaches Newton and takes the tangent to x = 1.5', () => {
    const visual = attached('numa.roots.newton')
    expect(visual.type).toBe('numa.roots')
    expect(visual.title).toBe('Follow the tangent to the axis')
    const params = rootsSchema.parse(visual.params)
    expect(params.method).toBe('newton')
    const result = newton(rootFunction(params.fn).f, rootFunction(params.fn).df, params.start, 5)
    expect(result.steps[0].next).toBeCloseTo(1.5, 12)
    expect(result.root).toBeCloseTo(1.324717957244746, 8)
    checkPrediction(visual, '1.5')
    expect(visual.predict?.explain).toContain('1 − (−1)/2')
  })
})

describe('graph visuals are attached to their lessons', () => {
  it('attaches breadth-first search and reports the second layer', () => {
    const visual = attached('graph.bfs')
    expect(visual.type).toBe('graph.network')
    expect(visual.title).toBe('Breadth-first grows in layers')
    const params = graphSchema.parse(visual.params)
    const result = runGraph(SAMPLE_GRAPH, params.algorithm, params.start)
    expect(result.order.map((step) => step.node)).toEqual(['A', 'B', 'D', 'C', 'E', 'F'])
    expect(result.order.find((step) => step.node === 'C')?.order).toBe(2)
    checkPrediction(visual, 'C')
  })

  it('attaches depth-first search and reports the next discovered node', () => {
    const visual = attached('graph.dfs')
    expect(visual.type).toBe('graph.network')
    expect(visual.title).toBe('Depth-first follows one branch')
    const params = graphSchema.parse(visual.params)
    const result = runGraph(SAMPLE_GRAPH, params.algorithm, params.start)
    expect(result.order.map((step) => step.node)).toEqual(['A', 'B', 'C', 'F', 'E', 'D'])
    checkPrediction(visual, 'F')
  })

  it('attaches Dijkstra and reports the shortest distance to F', () => {
    const visual = attached('graph.shortest_path')
    expect(visual.type).toBe('graph.network')
    expect(visual.title).toBe('Dijkstra keeps the shortest distances')
    const params = graphSchema.parse(visual.params)
    const result = runGraph(SAMPLE_GRAPH, params.algorithm, params.start)
    expect(result.distances).toEqual({ A: 0, B: 1, C: 2, D: 1, E: 2, F: 3 })
    checkPrediction(visual, '3')
    expect(visual.predict?.explain).toContain('A → B → E → F')
  })
})

describe('coverage records the new families honestly', () => {
  it('keeps the 31-topic proposal and adds supplementary entries', () => {
    const coverage = read('data/visualizations/coverage.json') as {
      topics: Array<{ id: string }>
      supplementary: Array<{ id: string; reason: string }>
    }
    expect(coverage.topics).toHaveLength(31)
    const ids = coverage.supplementary.map((item) => item.id)
    expect(ids).toContain('ode.slopefield')
    expect(ids).toContain('numa.roots')
    expect(ids).toContain('graph.network')
    expect(coverage.supplementary.every((item) => item.reason.length > 0)).toBe(true)
  })
})
