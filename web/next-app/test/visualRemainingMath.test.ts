import { describe, expect, it } from 'vitest'
import { ODE_FIELDS, integrateSlopeField, odeField, slopeAt } from '../components/visual/math/ode'
import { ROOT_FUNCTIONS, bisection, newton, newtonTangent, rootFunction } from '../components/visual/math/roots'
import { SAMPLE_GRAPH, breadthFirst, depthFirst, dijkstra, kruskal, runGraph } from '../components/visual/math/graph'

const close = (actual: number, expected: number, tolerance = 1e-9) => {
  expect(Number.isFinite(actual)).toBe(true)
  expect(Math.abs(actual - expected)).toBeLessThan(tolerance)
}

describe('slope fields', () => {
  it('evaluates each right-hand side, not a decorative angle', () => {
    close(slopeAt('exponential', 0, 2), 2)
    close(slopeAt('linear', 1, 2), -1)
    close(slopeAt('logistic', 0, 0.25), 0.1875)
    close(slopeAt('oscillating', 0, 2), -1)
    close(slopeAt('sum', 1, 2), 3)
    expect(ODE_FIELDS).toHaveLength(5)
    expect(odeField('logistic').equation).toContain('1 − y')
  })

  it('integrates y′ = y forward and backward with RK4', () => {
    const result = integrateSlopeField('exponential', 0, 1, {
      xMin: -1.5,
      xMax: 1.5,
      yMin: -1,
      yMax: 5,
      step: 0.01,
      method: 'rk4',
    })
    expect(result.status).toBe('ok')
    expect(result.points.length).toBeGreaterThan(200)
    const at = (x: number) => result.points.find((point) => Math.abs(point.x - x) < 1e-9)
    close(at(1)?.y ?? Number.NaN, Math.E, 1e-5)
    close(at(-1)?.y ?? Number.NaN, 1 / Math.E, 1e-5)
  })

  it('makes Euler visibly less accurate than RK4 at the same step', () => {
    const options = { xMin: -0.5, xMax: 1.5, yMin: -4, yMax: 8, step: 0.1 }
    const rk4 = integrateSlopeField('exponential', 0, 1, { ...options, method: 'rk4' })
    const euler = integrateSlopeField('exponential', 0, 1, { ...options, method: 'euler' })
    const end = (points: Array<{ x: number; y: number }>) => points.find((point) => Math.abs(point.x - 1) < 1e-9)?.y ?? Number.NaN
    const rk4Error = Math.abs(end(rk4.points) - Math.E)
    const eulerError = Math.abs(end(euler.points) - Math.E)
    expect(rk4Error).toBeLessThan(1e-4)
    expect(eulerError).toBeGreaterThan(rk4Error * 10)
  })

  it('keeps the constant solutions of the logistic field constant', () => {
    const zero = integrateSlopeField('logistic', 0, 0, {
      xMin: -2,
      xMax: 2,
      yMin: -1,
      yMax: 2,
      step: 0.1,
    })
    const one = integrateSlopeField('logistic', 0, 1, {
      xMin: -2,
      xMax: 2,
      yMin: -1,
      yMax: 2,
      step: 0.1,
    })
    expect(zero.points.every((point) => Math.abs(point.y) < 1e-12)).toBe(true)
    expect(one.points.every((point) => Math.abs(point.y - 1) < 1e-12)).toBe(true)
  })

  it('reports invalid and boundary outcomes without a non-finite point', () => {
    expect(
      integrateSlopeField('exponential', 0, 9, { xMin: -1, xMax: 1, yMin: -2, yMax: 2, step: 0.1 }).status,
    ).toBe('invalid')
    expect(
      integrateSlopeField('exponential', 0, 1, { xMin: -1, xMax: 1, yMin: -2, yMax: 2, step: 0 }).status,
    ).toBe('invalid')
    const escaped = integrateSlopeField('exponential', 0, 1, {
      xMin: -1,
      xMax: 1,
      yMin: -2,
      yMax: 2,
      step: 0.1,
    })
    expect(escaped.status).toBe('boundary')
    expect(escaped.points.every((point) => Number.isFinite(point.x) && Number.isFinite(point.y))).toBe(true)
  })

  it('is deterministic', () => {
    const options = { xMin: -1, xMax: 1, yMin: -2, yMax: 3, step: 0.05 }
    expect(integrateSlopeField('sum', 0, 0.5, options)).toEqual(integrateSlopeField('sum', 0, 0.5, options))
  })
})

describe('numerical root iterations', () => {
  it('bisects x² − 2 down to the bracket tolerance', () => {
    const spec = rootFunction('quadratic')
    const result = bisection(spec.f, 0, 2, 30)
    expect(result.status).toBe('ok')
    close(result.root ?? Number.NaN, Math.SQRT2, 1e-9)
    expect(result.steps.length).toBeGreaterThan(20)
    expect(result.steps.every((step) => step.keep === 'left' || step.keep === 'right' || step.keep === 'stop')).toBe(true)
  })

  it('rejects an invalid bracket and a bracket with no sign change', () => {
    const spec = rootFunction('quadratic')
    expect(bisection(spec.f, 2, 0, 5).status).toBe('invalid')
    expect(bisection(spec.f, 2, 3, 5).status).toBe('no-sign-change')
    expect(bisection(spec.f, 0, 0, 5).status).toBe('invalid')
  })

  it('stops at an endpoint root', () => {
    const result = bisection((x) => x, 0, 2, 6)
    expect(result.status).toBe('ok')
    close(result.root ?? Number.NaN, 0)
  })

  it('takes Newton steps for x² − 2', () => {
    const spec = rootFunction('quadratic')
    const result = newton(spec.f, spec.df, 1, 6)
    expect(result.status).toBe('ok')
    close(result.root ?? Number.NaN, Math.SQRT2, 1e-9)
    expect(result.steps.length).toBeLessThanOrEqual(6)
    expect(result.steps[0].next).toBeCloseTo(1.5, 12)
  })

  it('reports a stationary tangent instead of dividing by zero', () => {
    const result = newton((x) => x * x + 1, (x) => 2 * x, 0, 4)
    expect(result.status).toBe('stationary')
    expect(result.steps).toHaveLength(0)
  })

  it('draws the tangent that defines the next Newton guess', () => {
    const spec = rootFunction('quadratic')
    const tangent = newtonTangent(spec.f, spec.df, 1)
    close(tangent.slope, 2)
    close(tangent.intercept, 1.5)
    close(tangent.y, -1)
  })

  it('lists the exact values used by the readout', () => {
    expect(ROOT_FUNCTIONS.map((spec) => spec.exact)).toContain('√2 = 1.41421356…')
    expect(ROOT_FUNCTIONS.map((spec) => spec.exact).join(' ')).toContain('ln 3')
    close(rootFunction('exp').exactValue, Math.log(3))
  })
})

describe('graph algorithms', () => {
  it('breadth-first visits by layers', () => {
    const result = breadthFirst(SAMPLE_GRAPH, 'A')
    expect(result.status).toBe('ok')
    expect(result.order.map((step) => step.node)).toEqual(['A', 'B', 'D', 'C', 'E', 'F'])
    expect(result.tree).toHaveLength(5)
  })

  it('depth-first follows one branch before backtracking', () => {
    const result = depthFirst(SAMPLE_GRAPH, 'A')
    expect(result.order.map((step) => step.node)).toEqual(['A', 'B', 'C', 'F', 'E', 'D'])
    expect(result.tree).toHaveLength(5)
  })

  it('Dijkstra returns the actual weighted distances', () => {
    const result = dijkstra(SAMPLE_GRAPH, 'A')
    expect(result.status).toBe('ok')
    expect(result.distances).toMatchObject({ A: 0, B: 1, C: 2, D: 1, E: 2, F: 3 })
  })

  it('Kruskal builds a spanning tree of five edges', () => {
    const result = kruskal(SAMPLE_GRAPH)
    expect(result.status).toBe('ok')
    expect(result.tree).toHaveLength(5)
    expect(result.order).toHaveLength(5)
  })

  it('reports a disconnected graph and rejects an unknown start', () => {
    const disconnected = { nodes: SAMPLE_GRAPH.nodes.slice(0, 3), edges: [] }
    expect(breadthFirst(disconnected, 'A').status).toBe('disconnected')
    expect(runGraph(SAMPLE_GRAPH, 'bfs', 'Z').status).toBe('invalid')
  })
})
