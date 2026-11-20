import { describe, expect, it } from 'vitest'
import {
  criticalValue,
  densityCurve,
  inferenceModel,
  normalCdf,
  studentPdf,
  twoSidedPValue,
} from '../components/visual/math/inference'
import {
  phaseLine,
  phaseLineSlope,
  phaseNullclines,
  phaseSystem,
  phaseTrajectory,
  phaseVector,
} from '../components/visual/math/phase'
import { DAG_GRAPH, FLOW_GRAPH, maxFlow, topologicalSort } from '../components/visual/math/graph'
import { lagrangeModel, newtonMinimise, newtonObjective } from '../components/visual/math/optimization'

const close = (a: number, b: number, tolerance = 1e-6) =>
  expect(Math.abs(a - b)).toBeLessThan(tolerance)

describe('statistical inference', () => {
  it('uses the standard normal tail', () => {
    close(normalCdf(0), 0.5)
    close(normalCdf(1.96), 0.975, 1e-4)
    close(twoSidedPValue('z', 10, 2.5), 0.0124, 1e-3)
  })

  it('finds the two-sided critical value', () => {
    close(criticalValue('z', 10, 0.05), 1.96, 1e-3)
    expect(criticalValue('t', 9, 0.05)).toBeGreaterThan(criticalValue('z', 10, 0.05))
    expect(criticalValue('t', 1, 0.05)).toBeGreaterThan(12)
  })

  it('returns power, error rate and interval from one model', () => {
    const model = inferenceModel({ dist: 'z', df: 10, alpha: 0.05, effect: 0.8, observed: 1.9 })
    expect(model.status).toBe('ok')
    expect(model.reject).toBe(false)
    close(model.pValue, 0.0574, 1e-3)
    close(model.critical, 1.96, 1e-3)
    close(model.power, 0.126, 0.01)
    close(model.typeTwo, 1 - model.power)
    close(model.interval[0], 1.9 - model.critical)
    close(model.interval[1], 1.9 + model.critical)
  })

  it('rejects when the observed statistic crosses the critical value', () => {
    const model = inferenceModel({ dist: 'z', df: 10, alpha: 0.05, effect: 0, observed: 2.2 })
    expect(model.reject).toBe(true)
    expect(model.pValue).toBeLessThan(0.05)
  })

  it('uses a heavier t tail and a wider interval', () => {
    const z = inferenceModel({ dist: 'z', df: 10, alpha: 0.05, effect: 0.8, observed: 2 })
    const t = inferenceModel({ dist: 't', df: 5, alpha: 0.05, effect: 0.8, observed: 2 })
    expect(t.critical).toBeGreaterThan(z.critical)
    expect(t.interval[1] - t.interval[0]).toBeGreaterThan(z.interval[1] - z.interval[0])
    expect(t.pValue).toBeGreaterThan(z.pValue)
    expect(studentPdf(0, 5)).toBeGreaterThan(0)
  })

  it('returns finite density curves and rejects invalid input', () => {
    const points = densityCurve('t', 7, -3, 3, 40)
    expect(points).toHaveLength(41)
    expect(points.every(([x, y]) => Number.isFinite(x) && Number.isFinite(y))).toBe(true)
    expect(inferenceModel({ dist: 'z', df: 10, alpha: 0, effect: 0, observed: 0 }).status).toBe('invalid')
    expect(inferenceModel({ dist: 't', df: 0, alpha: 0.05, effect: 0, observed: 0 }).status).toBe('invalid')
  })
})

describe('phase portraits', () => {
  it('finds the equilibria of each phase line', () => {
    const logistic = phaseLine('logistic')
    expect(logistic.equilibria.map((item) => item.x)).toEqual([0, 1])
    expect(logistic.equilibria.map((item) => item.stable)).toEqual([false, true])
    const threshold = phaseLine('threshold')
    expect(threshold.equilibria.map((item) => item.x)).toEqual([-1, 1])
    expect(threshold.equilibria.map((item) => item.stable)).toEqual([true, false])
    close(phaseLineSlope('logistic', 0.5), 0.25)
    close(phaseLineSlope('decay', 2), -2)
  })

  it('classifies the linear systems from their eigenvalues', () => {
    const saddle = phaseSystem('saddle')
    expect(saddle.eigenvalues).toEqual([2, -3])
    expect(saddle.classification).toBe('saddle')
    const spiral = phaseSystem('spiral')
    expect(Array.isArray(spiral.eigenvalues)).toBe(false)
    expect(spiral.classification).toBe('stable spiral')
    expect(phaseSystem('center').stability).toBe('neutral')
  })

  it('uses the matrix for the vector field and nullclines', () => {
    const node = phaseSystem('stable-node')
    expect(phaseVector(node, 1, 2)).toEqual([-1, -4])
    expect(phaseNullclines(node)).toHaveLength(2)
    const trajectory = phaseTrajectory(node, [1, 1], 1, 100)
    expect(trajectory.status).toBe('ok')
    const end = trajectory.points[trajectory.points.length - 1]
    close(end[0], Math.exp(-1), 1e-6)
    close(end[1], Math.exp(-2), 1e-6)
  })

  it('reports invalid trajectories', () => {
    expect(phaseTrajectory(phaseSystem('saddle'), [Number.NaN, 0]).status).toBe('invalid')
  })
})

describe('graph extensions', () => {
  it('orders the DAG by removing sources', () => {
    const result = topologicalSort(DAG_GRAPH)
    expect(result.status).toBe('ok')
    expect(result.order.map((step) => step.node)).toEqual(['A', 'B', 'C', 'D', 'E'])
  })

  it('rejects a directed cycle', () => {
    const cyclic = {
      nodes: [
        { id: 'A', label: 'A', x: 0, y: 0 },
        { id: 'B', label: 'B', x: 1, y: 0 },
      ],
      edges: [
        { from: 'A', to: 'B', weight: 1 },
        { from: 'B', to: 'A', weight: 1 },
      ],
    }
    const result = topologicalSort(cyclic)
    expect(result.status).toBe('invalid')
    expect(result.note).toContain('cycle')
  })

  it('returns the maximum flow and the minimum cut', () => {
    const result = maxFlow(FLOW_GRAPH, 'S', 'T')
    expect(result.status).toBe('ok')
    expect(result.flow).toBe(5)
    expect(result.cut?.source).toEqual(['S'])
    expect(result.cut?.sink).toEqual(['A', 'B', 'T'])
    const edge = result.edgeFlows?.find((item) => item.from === 'S' && item.to === 'A')
    expect(edge?.flow).toBe(3)
  })
})

describe('nonlinear optimization', () => {
  it('solves the equality-constrained problem', () => {
    const optimum = lagrangeModel(1)
    expect(optimum.optimal).toBe(true)
    expect(optimum.lambda).toBe(2)
    expect(optimum.objectiveValue).toBe(2)
    const candidate = lagrangeModel(0.5)
    expect(candidate.y).toBe(1.5)
    expect(candidate.optimal).toBe(false)
    expect(candidate.along).toBeLessThan(0)
  })

  it('takes a real Newton step', () => {
    const run = newtonMinimise(1.5, 1)
    expect(run.points).toHaveLength(1)
    const first = run.points[0]
    close(first.next, 1.238095, 1e-6)
    close(first.df, 5.5)
    close(first.d2f, 21)
    expect(newtonObjective(1)).toBe(-1)
  })

  it('stops when the curvature is not positive', () => {
    const run = newtonMinimise(0, 3)
    expect(run.status).toBe('stationary')
    expect(Number.isNaN(run.points[0].next)).toBe(true)
  })
})
