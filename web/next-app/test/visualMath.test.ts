import { describe, expect, it } from 'vitest'
import {
  apply,
  collapseState,
  columns,
  compose,
  describeCollapse,
  det,
  dilation,
  fromFlat,
  identity,
  inverse,
  lerpMatrix,
  matrixScale,
  mul,
  reflection,
  rotation,
  toFlat,
  type Mat2,
  type Vec2,
} from '../components/visual/math/matrix'
import { antiderivative, definiteIntegral, derivative, evalPoly } from '../components/visual/math/poly'
import { riemannSum } from '../components/visual/math/integrate'
import { carg, cargOrNull, cmod, cmul, wrapAngle, type Complex } from '../components/visual/math/complex'
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

  it('reduces a sum of arguments into the principal range, and refuses to name a direction for zero', () => {
    // Two vectors at 120 degrees: their product sits at 240, which is -120 once wrapped.
    close(wrapAngle((2 * 2 * Math.PI) / 3), (-2 * Math.PI) / 3)
    close(wrapAngle(Math.PI), Math.PI)
    close(wrapAngle(-Math.PI), Math.PI)
    close(wrapAngle(0.3), 0.3)
    expect(cargOrNull([0, 0])).toBeNull()
    close(cargOrNull([1, 0]) as number, 0)
    close(cargOrNull([0, 1]) as number, Math.PI / 2)
  })
})

describe('normal distribution', () => {
  it('has the standard density peak and symmetry', () => {
    close(normalPdf(0, 0, 1), 1 / Math.sqrt(2 * Math.PI))
    close(normalPdf(1, 0, 1), normalPdf(-1, 0, 1))
    close(normalPdf(3, 2, 3), normalPdf(1, 2, 3))
    close(normalPdf(2, 2, 3), 1 / (3 * Math.sqrt(2 * Math.PI)))
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

describe('the matrix a plane is drawn with', () => {
  const A: Mat2 = [
    [1.2, 0.6],
    [0.4, 1.4],
  ]

  it('interpolates from the identity to A, endpoint to endpoint', () => {
    expect(toFlat(lerpMatrix(identity(), A, 0))).toEqual(toFlat(identity()))
    expect(toFlat(lerpMatrix(identity(), A, 1))).toEqual(toFlat(A))
    // Midpoint is the average of the entries, which is the whole definition.
    close(lerpMatrix(identity(), A, 0.5)[0][0], 1.1)
    close(lerpMatrix(identity(), A, 0.5)[1][1], 1.2)
    // t outside [0, 1] is clamped rather than extrapolated.
    expect(toFlat(lerpMatrix(identity(), A, -3))).toEqual(toFlat(identity()))
    expect(toFlat(lerpMatrix(identity(), A, 7))).toEqual(toFlat(A))
  })

  it('takes its columns from apply, so they cannot disagree with A e₁ and A e₂', () => {
    const [first, second] = columns(A)
    expect(first).toEqual(apply(A, [1, 0]))
    expect(second).toEqual(apply(A, [0, 1]))
  })

  it('builds a matrix-vector product as a combination of columns', () => {
    const v: Vec2 = [2, -1]
    const [first, second] = columns(A)
    const built: Vec2 = [v[0] * first[0] + v[1] * second[0], v[0] * first[1] + v[1] * second[1]]
    close(built[0], apply(A, v)[0])
    close(built[1], apply(A, v)[1])
  })

  it('reads the determinant as a signed area', () => {
    close(det(A), 1.44)
    close(Math.abs(det(reflection('x'))), 1)
    expect(det(reflection('x'))).toBeLessThan(0)
    expect(det(dilation(3))).toBe(9)
  })

  it('interpolates a quarter turn and a reflection through the identity', () => {
    const quarter = lerpMatrix(identity(), rotation(Math.PI / 2), 0.5)
    expect(toFlat(quarter)).toEqual([0.5, -0.5, 0.5, 0.5])
    // A reflection passes through a singular matrix halfway: det drops to zero
    // and comes back negative. That is a real fact about the path, not noise.
    const halfway = lerpMatrix(identity(), reflection('y'), 0.5)
    close(det(halfway), 0)
    expect(collapseState(halfway, 55)).toBe('singular')
    expect(det(lerpMatrix(identity(), reflection('y'), 1))).toBe(-1)
  })

  it('calls only a zero determinant singular', () => {
    expect(collapseState([[1, 0], [0, 1]], 55)).toBe('none')
    expect(collapseState([[1, 2], [2, 4]], 55)).toBe('singular')
    expect(collapseState([[1, 0], [0, 0]], 55)).toBe('singular')
    expect(collapseState([[0, 0], [0, 0]], 55)).toBe('singular')
  })

  it('separates near, pixel and singular without changing invertibility', () => {
    // Determinant 0.01: an area, a real inverse, and a height of about 0.37px at
    // 37px per unit, so it is invisible at this scale but not singular.
    const thin: Mat2 = [
      [1, 0],
      [0, 0.01],
    ]
    expect(collapseState(thin, 37)).toBe('pixel')
    expect(inverse(thin)).not.toBeNull()
    close(Math.abs(det(thin)), 0.01)

    // Determinant 1e-4 is also thinner than a pixel at this scale, and the
    // drawing verdict is the one that matters for a shape, so pixel wins.
    const thinner: Mat2 = [
      [1, 0],
      [0, 1e-4],
    ]
    expect(collapseState(thinner, 37)).toBe('pixel')
    expect(inverse(thinner)).not.toBeNull()

    // Nearly collapsed is the case where the determinant is small relative to
    // the entries but the image is still thick enough to draw: a large map with a
    // tiny area scale.
    const nearly: Mat2 = [
      [20, 0],
      [0, 0.016],
    ]
    expect(collapseState(nearly, 37)).toBe('near')
    expect(inverse(nearly)).not.toBeNull()

    expect(collapseState([[1, 2], [2, 4]], 37)).toBe('singular')
    expect(inverse([[1, 2], [2, 4]])).toBeNull()
  })

  it('describes each state in words a learner can act on', () => {
    expect(describeCollapse('none', 1)).toBe('full parallelogram')
    expect(describeCollapse('singular', 0)).toContain('line')
    expect(describeCollapse('pixel', 0.01)).toContain('thinner than a pixel')
    expect(describeCollapse('near', 1e-4)).toContain('nearly collapsed')
  })

  it('measures scale by the largest entry', () => {
    expect(matrixScale(A)).toBe(1.4)
    expect(matrixScale([[0, 0], [0, 0]])).toBe(0)
    expect(matrixScale([[-3, 1], [0, 2]])).toBe(3)
  })

  it('handles a projection matrix as a real singular map', () => {
    const project: Mat2 = [
      [1, 0],
      [0, 0],
    ]
    expect(det(project)).toBe(0)
    expect(inverse(project)).toBeNull()
    // Every vector loses its second component, which is what the figure draws.
    expect(apply(project, [3, 5])).toEqual([3, 0])
    expect(collapseState(project, 55)).toBe('singular')
  })
})
