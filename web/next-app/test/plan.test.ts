import { describe, it, expect } from 'vitest'
import { clampGoal, daysFor, monthLabel, GOAL_PRESETS, GOAL_MIN, GOAL_MAX, DEFAULT_DAILY_GOAL } from '../lib/plan'

describe('clampGoal', () => {
  it('clamps to the server bounds', () => {
    // 10..100, matching the server's check in handleSetDailyXPGoal. These were 1..10000,
    // and the planner's fallbacks said 30 while the shipped default was 10 — so a learner
    // with no stored goal planned against 3x the XP/day they were actually banking.
    expect(clampGoal(0)).toBe(GOAL_MIN)
    expect(clampGoal(-5)).toBe(GOAL_MIN)
    expect(clampGoal(9)).toBe(GOAL_MIN)
    expect(clampGoal(50000)).toBe(GOAL_MAX)
    expect(clampGoal(101)).toBe(GOAL_MAX)
    expect(clampGoal(50.7)).toBe(51)
  })

  it('falls back to the documented default, not a literal', () => {
    // Clearing the custom field yields NaN on parse. That used to become 30 silently,
    // rewriting the goal to a value the rest of the system disagreed with.
    expect(clampGoal(Number.NaN)).toBe(DEFAULT_DAILY_GOAL)
    expect(clampGoal(Number.POSITIVE_INFINITY)).toBe(DEFAULT_DAILY_GOAL)
    expect(clampGoal(Number.NEGATIVE_INFINITY)).toBe(DEFAULT_DAILY_GOAL)
  })

  it('presets are within bounds', () => {
    for (const p of GOAL_PRESETS) expect(clampGoal(p)).toBe(p)
    expect(GOAL_PRESETS).toContain(DEFAULT_DAILY_GOAL)
  })

  it('the default is inside the range it clamps to', () => {
    expect(DEFAULT_DAILY_GOAL).toBeGreaterThanOrEqual(GOAL_MIN)
    expect(DEFAULT_DAILY_GOAL).toBeLessThanOrEqual(GOAL_MAX)
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
