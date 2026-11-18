import { describe, expect, it } from 'vitest'
import {
  apply,
  compose,
  det,
  dilation,
  fromFlat,
  identity,
  mul,
  reflection,
  rotation,
  toFlat,
  type Mat2,
  type Vec2,
} from '../components/visual/math/matrix'
import { antiderivative, definiteIntegral, derivative, evalPoly } from '../components/visual/math/poly'
import { riemannSum } from '../components/visual/math/integrate'
import { carg, cmod, cmul, type Complex } from '../components/visual/math/complex'
import { normalCdf, normalPdf, normalProbability } from '../components/visual/math/normal'
import { distance, maxSideDelta, sideLengths } from '../components/visual/math/geometry'
import { clamp, invertLinear, scaleLinear, ticks } from '../components/visual/scale'

const close = (a: number, b: number, tolerance = 1e-12) => expect(Math.abs(a - b)).toBeLessThan(tolerance)

describe('matrix geometry', () => {
  it('applies a row-major matrix to a vector', () => {
    const m: Mat2 = [
      [2, 1],
      [0, 3],
    ]
    expect(apply(m, [4, 5])).toEqual([13, 15])
  })

  it('multiplies matrices in the same order as function composition', () => {
    const rotate = rotation(Math.PI / 2)
    const scale = dilation(2)
    close(apply(mul(rotate, scale), [1, 0])[0], 0)
    close(apply(mul(rotate, scale), [1, 0])[1], 2)
    const composed = compose(rotate, scale)
    expect(toFlat(composed)).toEqual(toFlat(mul(rotate, scale)))
  })

  it('keeps the determinant under composition and reads a flat matrix', () => {
    const m = fromFlat([1.2, 0.6, 0.4, 1.4])
    expect(det(m)).toBeCloseTo(1.44, 12)
    expect(det(mul(m, m))).toBeCloseTo(det(m) ** 2, 12)
    expect(toFlat(m)).toEqual([1.2, 0.6, 0.4, 1.4])
    expect(det(identity())).toBe(1)
  })

  it('rotates and reflects without changing length', () => {
    const v: Vec2 = [3, 4]
    const rotated = apply(rotation(0.7), v)
    close(Math.hypot(...rotated), 5)
    close(apply(reflection('x'), v)[1], -4)
    close(apply(reflection('y'), v)[0], -3)
    close(apply(reflection('origin'), v)[0], -3)
    close(apply(reflection('origin'), v)[1], -4)
  })
})

describe('polynomials and Riemann sums', () => {
  it('evaluates and differentiates an ascending coefficient list', () => {
    const f = [2, 0, 3]
    expect(evalPoly(f, 2)).toBe(14)
    expect(derivative(f)).toEqual([0, 6])
    expect(antiderivative(f)).toEqual([0, 2, 0, 1])
  })

  it('integrates x squared from 0 to 2 exactly', () => {
    close(definiteIntegral([0, 0, 1], 0, 2), 8 / 3)
    close(definiteIntegral([0, 0, 1], -1, 1), 2 / 3)
    close(definiteIntegral([5], 3, 7), 20)
  })

  it('computes left, right and midpoint sums against a hand sum', () => {
    const left = riemannSum([0, 0, 1], 0, 2, 4, 'left')
    close(left, (0 + 0.25 + 1 + 2.25) * 0.5)
    const right = riemannSum([0, 0, 1], 0, 2, 4, 'right')
    close(right, (0.25 + 1 + 2.25 + 4) * 0.5)
    const mid = riemannSum([0, 0, 1], 0, 2, 4, 'mid')
    close(mid, (0.0625 + 0.5625 + 1.5625 + 3.0625) * 0.5)
    close(riemannSum([0, 0, 1], 0, 2, 1, 'left'), 0)
  })
})

describe('complex arithmetic', () => {
  it('multiplies in rectangular form', () => {
    expect(cmul([1, 2], [3, 4])).toEqual([-5, 10])
    expect(cmul([0, 1], [0, 1])).toEqual([-1, 0])
  })

  it('multiplies moduli and adds arguments', () => {
    const z: Complex = [1, 1]
    const w: Complex = [0, 2]
    const product = cmul(z, w)
    close(cmod(product), cmod(z) * cmod(w))
    close(carg(product), carg(z) + carg(w))
    close(carg([-1, 0]), Math.PI)
    close(carg([0, -1]), -Math.PI / 2)
  })
})

describe('normal distribution', () => {
  it('has the standard density peak and symmetry', () => {
    close(normalPdf(0, 0, 1), 1 / Math.sqrt(2 * Math.PI))
    close(normalPdf(1, 0, 1), normalPdf(-1, 0, 1))
    close(normalPdf(2, 2, 3), normalPdf(2, 2, 3))
  })

  it('matches the known cumulative probabilities', () => {
    close(normalCdf(0), 0.5, 1e-7)
    close(normalCdf(1.96), 0.975002, 1e-5)
    close(normalCdf(-1.96), 0.024998, 1e-5)
    close(normalCdf(3), 0.998650, 1e-5)
    close(normalProbability(1, 0, 1), normalCdf(1), 1e-12)
    close(normalProbability(5, 3, 2), normalCdf(1), 1e-12)
  })
})

describe('geometry helpers', () => {
  const shape: Vec2[] = [
    [0, 0],
    [2, 0],
    [1, 2],
  ]

  it('measures sides and distances', () => {
    close(distance([0, 0], [3, 4]), 5)
    expect(sideLengths(shape)).toHaveLength(3)
    close(sideLengths(shape)[0], 2)
  })

  it('reports no side change for an isometry and a real change for dilation', () => {
    close(maxSideDelta(shape, shape.map(([x, y]) => apply(rotation(0.4), [x, y]))), 0)
    close(maxSideDelta(shape, shape.map(([x, y]) => apply(reflection('y'), [x, y]))), 0)
    const square: Vec2[] = [
      [0, 0],
      [2, 0],
      [2, 2],
      [0, 2],
    ]
    close(maxSideDelta(square, square.map(([x, y]) => apply(dilation(2), [x, y]))), 2)
  })
})

describe('scales and ticks', () => {
  it('maps and inverts a linear domain', () => {
    const s = scaleLinear(-1, 1, 0, 100)
    close(s(0), 50)
    close(invertLinear(-1, 1, 0, 100)(50), 0)
    close(clamp(5, 0, 2), 2)
  })

  it('produces ordered ticks that cover the domain', () => {
    const t = ticks(-2, 2, 4)
    expect(t[0]).toBe(-2)
    expect(t[t.length - 1]).toBe(2)
    expect(t.every((v, i) => i === 0 || v > t[i - 1])).toBe(true)
  })
})
