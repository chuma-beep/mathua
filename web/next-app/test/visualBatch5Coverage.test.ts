import { describe, expect, it } from 'vitest'
import { readFileSync } from 'node:fs'
import { join } from 'node:path'
import { inferenceSchema, phaseSchema, graphSchema, optimizationSchema } from '../components/visual/schemas'
import { inferenceModel } from '../components/visual/math/inference'
import { phaseLine, phaseSystem } from '../components/visual/math/phase'
import { DAG_GRAPH, FLOW_GRAPH, topologicalSort, maxFlow } from '../components/visual/math/graph'
import { lagrangeModel, newtonMinimise } from '../components/visual/math/optimization'

const REPO_ROOT = join(__dirname, '../../..')
const shard = (id: string): unknown =>
  JSON.parse(readFileSync(join(REPO_ROOT, 'data', 'lessons', 'kp', `${id}.json`), 'utf8'))

interface Step {
  label?: string
  visual?: {
    type: string
    title?: string
    caption?: string
    params?: Record<string, number | string | boolean>
    predict?: { question: string; options: string[]; answer: number; explain: string }
  }
}

const steps = (id: string) => shard(id) as Step[]
const visualOf = (id: string) => {
  const list = steps(id)
  expect(list[0].visual?.type).toBeTruthy()
  expect(list.filter((step) => step.visual)).toHaveLength(1)
  return list[0].visual!
}

describe('inference lessons', () => {
  const cases: Array<{ id: string; answer: string; contains: string }> = [
    { id: 'stat.infer.hypothesis', answer: 'The total tail area used to reject the null', contains: 'Type I' },
    { id: 'stat.infer.pvalue', answer: 'About 0.012', contains: 'doubles' },
    { id: 'stat.infer.type_error', answer: 'Failing to reject a false null', contains: 'beta' },
    { id: 'stat.infer.power', answer: 'A larger true effect', contains: 'rejection region' },
    { id: 'stat.infer.ztest', answer: 'Reject because the two-sided p-value is below 0.05', contains: '0.028' },
    { id: 'stat.infer.ci_mean_t', answer: 'The t critical value is larger', contains: '2.262' },
    { id: 'stat.infer.ttest', answer: 'Fail to reject at 5 percent', contains: '2.571' },
  ]

  for (const item of cases) {
    it(`${item.id} tests the value its prediction names`, () => {
      const visual = visualOf(item.id)
      expect(visual.type).toBe('stat.inference')
      const parsed = inferenceSchema.parse(visual.params)
      const model = inferenceModel(parsed)
      expect(model.status).toBe('ok')
      expect(model.pValue).toBeGreaterThanOrEqual(0)
      expect(model.pValue).toBeLessThanOrEqual(1)
      expect(model.power + model.typeTwo).toBeCloseTo(1, 9)
      expect(visual.predict?.options[visual.predict.answer]).toBe(item.answer)
      expect(visual.predict?.explain).toContain(item.contains)
    })
  }

  it('the z-test prediction agrees with the computed p-value', () => {
    const parsed = inferenceSchema.parse(visualOf('stat.infer.ztest').params)
    const model = inferenceModel(parsed)
    expect(model.pValue).toBeCloseTo(0.0278, 3)
    expect(model.reject).toBe(true)
  })

  it('the t interval is wider than the normal interval', () => {
    const parsed = inferenceSchema.parse(visualOf('stat.infer.ci_mean_t').params)
    const model = inferenceModel(parsed)
    expect(parsed.dist).toBe('t')
    expect(model.critical).toBeGreaterThan(1.96)
    expect(model.interval[1] - model.interval[0]).toBeGreaterThan(2 * 1.96)
  })
})

describe('phase lessons', () => {
  const cases: Array<{ id: string; answer: string }> = [
    { id: 'ode.stability', answer: 'It approaches zero' },
    { id: 'ode.systems.linear_phase', answer: 'A saddle' },
    { id: 'ode.dynamical.vector_field', answer: 'Zero in both components' },
    { id: 'ode.dynamical.nullcline', answer: 'The x component' },
    { id: 'ode.dynamical.linearisation', answer: 'A stable spiral' },
  ]

  for (const item of cases) {
    it(`${item.id} draws the phase object its prediction names`, () => {
      const visual = visualOf(item.id)
      expect(visual.type).toBe('ode.phase')
      const parsed = phaseSchema.parse(visual.params)
      if (parsed.mode === 'line') {
        expect(phaseLine(parsed.lineModel).equilibria.length).toBeGreaterThan(0)
      } else {
        expect(phaseSystem(parsed.system).classification.length).toBeGreaterThan(0)
      }
      expect(visual.predict?.options[visual.predict.answer]).toBe(item.answer)
    })
  }

  it('keeps the saddle eigenvalues and the spiral real part', () => {
    const saddle = phaseSystem(phaseSchema.parse(visualOf('ode.systems.linear_phase').params).system)
    expect(saddle.eigenvalues).toEqual([2, -3])
    const spiral = phaseSystem(phaseSchema.parse(visualOf('ode.dynamical.linearisation').params).system)
    expect(Array.isArray(spiral.eigenvalues)).toBe(false)
    if (!Array.isArray(spiral.eigenvalues)) {
      expect(spiral.eigenvalues.real).toBeLessThan(0)
      expect(spiral.eigenvalues.imaginary).toBeGreaterThan(0)
    }
  })
})

describe('graph lessons', () => {
  it('topological sort follows the DAG edges', () => {
    const visual = visualOf('graph.topological_sort')
    expect(visual.type).toBe('graph.network')
    const parsed = graphSchema.parse(visual.params)
    expect(parsed.algorithm).toBe('topological')
    const result = topologicalSort(DAG_GRAPH)
    expect(result.status).toBe('ok')
    expect(result.order.map((step) => step.node)).toEqual(['A', 'B', 'C', 'D', 'E'])
    expect(visual.predict?.options[visual.predict.answer]).toBe('A vertex with in-degree zero')
  })

  it('maximum flow equals the minimum cut', () => {
    const visual = visualOf('graph.max_flow')
    const parsed = graphSchema.parse(visual.params)
    expect(parsed.algorithm).toBe('maxflow')
    const result = maxFlow(FLOW_GRAPH, 'S', 'T')
    expect(result.flow).toBe(5)
    expect(result.cut?.source).toEqual(['S'])
    expect(result.cut?.sink).toEqual(['A', 'B', 'T'])
    expect(visual.predict?.options[visual.predict.answer]).toBe('5')
  })
})

describe('nonlinear optimization lessons', () => {
  it('Lagrange optimum and multiplier agree with the shard', () => {
    const visual = visualOf('optim.equality_lagrange')
    expect(visual.type).toBe('optim.path')
    const parsed = optimizationSchema.parse(visual.params)
    expect(parsed.mode).toBe('lagrange')
    const model = lagrangeModel(parsed.candidateX)
    expect(model.y).toBeCloseTo(2 - parsed.candidateX, 9)
    expect(model.optimum.lambda).toBe(2)
    expect(model.optimum.value).toBe(2)
    expect(visual.predict?.options[visual.predict.answer]).toBe('2')
  })

  it('Newton next guess agrees with the shard explanation', () => {
    const visual = visualOf('optim.newton_optim')
    const parsed = optimizationSchema.parse(visual.params)
    expect(parsed.mode).toBe('newton')
    const run = newtonMinimise(parsed.startX, 1)
    expect(run.points[0].next).toBeCloseTo(1.238095, 6)
    expect(visual.predict?.options[visual.predict.answer]).toBe('The first and second derivatives')
  })
})

describe('coverage ledger', () => {
  it('records both new templates without inflating the proposal', () => {
    const coverage = JSON.parse(
      readFileSync(join(REPO_ROOT, 'data', 'visualizations', 'coverage.json'), 'utf8'),
    ) as { topics: unknown[]; supplementary: Array<{ id: string; reason: string }> }
    expect(coverage.topics).toHaveLength(31)
    const ids = coverage.supplementary.map((entry) => entry.id)
    expect(ids).toContain('stat.inference')
    expect(ids).toContain('ode.phase')
    expect(coverage.supplementary.find((entry) => entry.id === 'stat.inference')?.reason).toContain('null')
    expect(coverage.supplementary.find((entry) => entry.id === 'ode.phase')?.reason).toContain('stability')
  })
})
