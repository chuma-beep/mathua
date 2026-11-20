import { describe, expect, it } from 'vitest'
import {
  FUNCTION_FAMILIES,
  domainOf,
  endBehavior,
  functionFormula,
  functionValue,
  inverseFormula,
  inverseIsFunction,
  rangeOf,
} from '../components/visual/math/functionGraph'

const close = (a: number | null, b: number, tolerance = 1e-12) => {
  expect(a).not.toBeNull()
  expect(Math.abs((a as number) - b)).toBeLessThan(tolerance)
}

describe('function graph engine', () => {
  it('lists the three families', () => {
    expect(FUNCTION_FAMILIES.map((item) => item.id)).toEqual(['quadratic', 'cubic', 'exponential'])
  })

  it('evaluates each family with scale and shifts', () => {
    close(functionValue('quadratic', 1, 0, 0, 3), 9)
    close(functionValue('quadratic', -1, 1, -2, 3), -6)
    close(functionValue('cubic', 2, -1, 1, 2), 55)
    close(functionValue('exponential', 1, 0, 0, 3), 8)
    close(functionValue('exponential', -1, 0, 1, 2), -3)
  })

  it('returns null for non-finite input or overflow', () => {
    expect(functionValue('quadratic', Number.NaN, 0, 0, 1)).toBeNull()
    expect(functionValue('exponential', 1, 0, 0, Number.POSITIVE_INFINITY)).toBeNull()
    expect(functionValue('cubic', 1, 0, 0, Number.MAX_VALUE)).toBeNull()
  })

  it('writes formulas with signs, shifts and scale', () => {
    expect(functionFormula('quadratic', 1, 0, 0)).toBe('x²')
    expect(functionFormula('quadratic', -1, 1, -2)).toBe('−(x − 1)² − 2')
    expect(functionFormula('quadratic', 2, -1, 3)).toBe('2(x + 1)² + 3')
    expect(functionFormula('cubic', -1, 0, 0)).toBe('−x³')
    expect(functionFormula('exponential', 1, 0, 0)).toBe('2^x')
    expect(functionFormula('exponential', 2, 1, -1)).toBe('2·2^(x − 1) − 1')
  })

  it('decides when an inverse is a function', () => {
    expect(inverseIsFunction('quadratic', 1)).toBe(false)
    expect(inverseIsFunction('cubic', 1)).toBe(true)
    expect(inverseIsFunction('cubic', 0)).toBe(false)
    expect(inverseIsFunction('exponential', -1)).toBe(true)
    expect(inverseFormula('quadratic', 1, 0, 0)).toBeNull()
    expect(inverseFormula('cubic', 1, 0, 0)).toBe('0 + ∛(y)')
    expect(inverseFormula('exponential', 1, 0, 0)).toBe('0 + log₂(y)')
  })

  it('reports end behavior from degree and sign', () => {
    expect(endBehavior('quadratic', 1)).toEqual({ left: '+∞', right: '+∞' })
    expect(endBehavior('quadratic', -1)).toEqual({ left: '−∞', right: '−∞' })
    expect(endBehavior('cubic', 1)).toEqual({ left: '−∞', right: '+∞' })
    expect(endBehavior('cubic', -1)).toEqual({ left: '+∞', right: '−∞' })
    expect(endBehavior('exponential', 1)).toEqual({ left: '→ 0', right: '+∞' })
    expect(endBehavior('exponential', -1)).toEqual({ left: '→ 0', right: '−∞' })
  })

  it('reports domain and range', () => {
    expect(domainOf('quadratic')).toBe('all real x')
    expect(domainOf('exponential')).toBe('all real x')
    expect(rangeOf('quadratic', 1, 2)).toBe('[2, ∞)')
    expect(rangeOf('quadratic', -1, 2)).toBe('(−∞, 2]')
    expect(rangeOf('quadratic', 0, 2)).toBe('{2}')
    expect(rangeOf('cubic', -1, 2)).toBe('all real y')
    expect(rangeOf('exponential', 1, 2)).toBe('(2, ∞)')
    expect(rangeOf('exponential', -1, 2)).toBe('(−∞, 2)')
    expect(rangeOf('exponential', 0, 2)).toBe('{2}')
  })

  it('is deterministic', () => {
    const first = functionValue('exponential', 1, 0, 0, 1.5)
    const second = functionValue('exponential', 1, 0, 0, 1.5)
    expect(first).toBe(second)
  })
})
