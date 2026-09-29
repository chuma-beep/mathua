import { describe, it, expect } from 'vitest'
import { clampGoal, daysFor, monthLabel, GOAL_PRESETS } from '../lib/plan'

describe('clampGoal', () => {
  it('clamps to server bounds', () => {
    expect(clampGoal(0)).toBe(1)
    expect(clampGoal(-5)).toBe(1)
    expect(clampGoal(50000)).toBe(10000)
    expect(clampGoal(50.7)).toBe(51)
    expect(clampGoal(Number.NaN)).toBe(30)
  })

  it('presets are within bounds', () => {
    for (const p of GOAL_PRESETS) expect(clampGoal(p)).toBe(p)
  })
})

describe('daysFor', () => {
  it('divides workload by rest-adjusted rate', () => {
    expect(daysFor(4500, 50, 0)).toBe(90)
    expect(daysFor(4500, 50, 2)).toBeCloseTo(4500 / (50 * 5 / 7), 6)
  })

  it('returns -1 when unknowable', () => {
    expect(daysFor(0, 50, 0)).toBe(-1)
    expect(daysFor(100, 0, 0)).toBe(-1)
  })
})

describe('monthLabel', () => {
  it('drops day precision', () => {
    expect(monthLabel('2026-03-27')).toBe('March 2026')
    expect(monthLabel('nope')).toBe('')
  })
})
