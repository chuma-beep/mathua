import { describe, expect, it } from 'vitest'
import {
  linearProgram,
  objectiveContour,
  type ConstraintCheck,
  type KktMultiplier,
} from '../components/visual/math/optimization'

const near = (value: number, expected: number, tolerance = 1e-9) =>
  expect(Math.abs(value - expected)).toBeLessThan(tolerance)

const label = (items: Array<ConstraintCheck | KktMultiplier>, text: string) =>
  items.find((item) => item.label.includes(text))

describe('linear program active sets and KKT', () => {
  it('finds the vertex optimum and both active resource constraints', () => {
    const lp = linearProgram(2, 1)
    expect(lp.status).toBe('optimal')
    expect(lp.vertices).toHaveLength(4)
    expect(lp.optimum).toMatchObject({ x: 2, y: 1, value: 8 })
    expect(lp.feasible).toBe(true)
    expect(lp.active).toHaveLength(2)
    expect(lp.active.join(' ')).toContain('x + 2y')
    expect(lp.active.join(' ')).toContain('2x + y')
    expect(lp.kkt.kkt).toBe(true)
    expect(lp.optimumKkt.kkt).toBe(true)
    near(label(lp.optimumKkt.multipliers, 'x + 2y')?.value ?? Number.NaN, 1 / 3)
    near(label(lp.optimumKkt.multipliers, '2x + y')?.value ?? Number.NaN, 4 / 3)
    near(lp.dualValue, 8)
    near(lp.gap, 0)
  })

  it('rejects an interior candidate as a KKT point', () => {
    const lp = linearProgram(1, 1)
    expect(lp.feasible).toBe(true)
    expect(lp.active).toEqual([])
    expect(lp.kkt.kkt).toBe(false)
    expect(lp.kkt.active).toEqual([])
    near(lp.kkt.residual, Math.sqrt(13))
    expect(lp.kkt.note).toContain('active normals')
  })

  it('rejects an infeasible candidate', () => {
    const lp = linearProgram(3, 3)
    expect(lp.feasible).toBe(false)
    expect(lp.kkt.kkt).toBe(false)
    expect(lp.kkt.note).toContain('not primal feasible')
  })

  it('gives a negative multiplier at the origin, so KKT fails', () => {
    const lp = linearProgram(0, 0)
    expect(lp.active).toHaveLength(2)
    expect(lp.kkt.stationary).toBe(true)
    expect(lp.kkt.dualFeasible).toBe(false)
    expect(lp.kkt.kkt).toBe(false)
    expect(lp.kkt.multipliers.every((item) => item.value < 0)).toBe(true)
  })

  it('fails KKT on the edge where one multiplier is negative', () => {
    const lp = linearProgram(2.5, 0)
    expect(lp.active).toContain('y ≥ 0')
    expect(lp.kkt.stationary).toBe(true)
    expect(lp.kkt.dualFeasible).toBe(false)
    expect(lp.kkt.kkt).toBe(false)
  })

  it('keeps the active-set optimum when no bound is shifted', () => {
    const lp = linearProgram(2, 1)
    expect(lp.shadow).toBeNull()
    expect(lp.allowable).toBeNull()
  })
})

describe('right-hand-side shifts and dual values', () => {
  it('moves the optimum along the other active constraint and prices the first bound', () => {
    const lp = linearProgram(1.9333333333333333, 1.1333333333333333, {
      which: 'first',
      delta: 0.2,
    })
    near(lp.optimum.x, 5.8 / 3)
    near(lp.optimum.y, 3.4 / 3)
    near(lp.optimum.value, 8 + 0.2 / 3)
    expect(lp.feasible).toBe(true)
    expect(lp.kkt.kkt).toBe(true)
    near(lp.shadow?.value ?? Number.NaN, 1 / 3)
    expect(lp.allowable).toEqual({ label: 'x + 2y', min: 2.5, max: 10 })
    near(lp.dualValue, lp.optimum.value)
    near(lp.gap, 0)
  })

  it('prices the second bound with its own multiplier', () => {
    const lp = linearProgram(2.1333333333333333, 0.9333333333333333, {
      which: 'second',
      delta: 0.2,
    })
    near(lp.optimum.x, 6.4 / 3)
    near(lp.optimum.y, 2.8 / 3)
    near(lp.optimum.value, 8 + 0.8 / 3)
    expect(lp.kkt.kkt).toBe(true)
    near(lp.shadow?.value ?? Number.NaN, 4 / 3)
    expect(lp.allowable).toEqual({ label: '2x + y', min: 2, max: 8 })
    near(lp.dualValue, lp.optimum.value)
    near(lp.gap, 0)
  })

  it('handles a negative shift without changing the basis', () => {
    const lp = linearProgram(2.3333333333333333, 0.3333333333333333, {
      which: 'first',
      delta: -1,
    })
    near(lp.optimum.x, 7 / 3)
    near(lp.optimum.y, 1 / 3)
    near(lp.optimum.value, 23 / 3)
    near(lp.shadow?.value ?? Number.NaN, 1 / 3)
    near(lp.gap, 0)
  })

  it('reports an empty polygon when a bound makes the program infeasible', () => {
    const lp = linearProgram(0, 0, { which: 'first', delta: -100 })
    expect(lp.status).toBe('infeasible')
    expect(lp.vertices).toEqual([])
    expect(lp.dualValue).toBe(0)
    expect(Number.isNaN(lp.gap)).toBe(true)
  })

  it('is deterministic and finite', () => {
    const first = linearProgram(2, 1, { which: 'first', delta: 0.35 })
    const second = linearProgram(2, 1, { which: 'first', delta: 0.35 })
    expect(second.optimum).toEqual(first.optimum)
    expect(second.dualValue).toBe(first.dualValue)
    expect(Number.isFinite(second.dualValue)).toBe(true)
    expect(Number.isFinite(second.gap)).toBe(true)
  })
})

describe('objective contours', () => {
  it('returns a direction perpendicular to the objective gradient', () => {
    const contour = objectiveContour(6)
    expect(contour.point).toEqual([2, 0])
    expect(3 * contour.direction[0] + 2 * contour.direction[1]).toBe(0)
  })

  it('keeps a finite line for invalid levels', () => {
    const contour = objectiveContour(Number.NaN)
    expect(contour.point.every(Number.isFinite)).toBe(true)
    expect(contour.direction.every(Number.isFinite)).toBe(true)
  })
})
