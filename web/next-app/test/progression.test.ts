import { describe, it, expect } from 'vitest'
import { REQUIRED_IN_A_ROW, nextStreak, isAdvanced, initialState, applyResult } from '../lib/progression'

describe('lesson practice progression (2-in-a-row to advance)', () => {
  it('grows consecutive on correct, resets on incorrect', () => {
    expect(nextStreak(0, true)).toBe(1)
    expect(nextStreak(1, true)).toBe(2)
    expect(nextStreak(2, true)).toBe(3)
    expect(nextStreak(3, false)).toBe(0)
  })

  it('advances only at the required streak', () => {
    expect(REQUIRED_IN_A_ROW).toBe(2)
    expect(isAdvanced(0)).toBe(false)
    expect(isAdvanced(1)).toBe(false)
    expect(isAdvanced(2)).toBe(true)
    expect(isAdvanced(3)).toBe(true)
  })

  it('starts unseen', () => {
    expect(initialState()).toEqual({ consecutive: 0, advanced: false })
  })

  it('applies results monotonically once advanced', () => {
    let s = initialState()
    s = applyResult(s, true) // 1
    expect(s.advanced).toBe(false)
    s = applyResult(s, true) // 2
    expect(s.advanced).toBe(true)
    s = applyResult(s, false) // streak resets but advanced persists
    expect(s.consecutive).toBe(0)
    expect(s.advanced).toBe(true)
  })

  it('never advances on a miss', () => {
    let s = initialState()
    s = applyResult(s, true)
    s = applyResult(s, false)
    s = applyResult(s, true)
    expect(s.advanced).toBe(false)
    expect(s.consecutive).toBe(1)
  })
})
