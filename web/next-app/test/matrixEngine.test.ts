import { describe, expect, it } from 'vitest'
import {
  EPS,
  PRESETS,
  apply,
  det,
  eigen,
  eigenvector,
  fromFlat,
  identity,
  inverse,
  mul,
  normalize,
  rotation,
  type Mat2,
} from '../components/visual/math/matrix'

const close = (a: number, b: number, tolerance = 1e-9) => expect(Math.abs(a - b)).toBeLessThan(tolerance)

describe('inverse', () => {
  it('inverts a nonsingular matrix and undoes it', () => {
    const a = fromFlat([2, 1, 1, 1])
    const inv = inverse(a)
    expect(inv).not.toBeNull()
    const product = mul(a, inv as Mat2)
    close(product[0][0], 1)
    close(product[0][1], 0)
    close(product[1][0], 0)
    close(product[1][1], 1)
    close(apply(inv as Mat2, apply(a, [3, -2]))[0], 3)
  })

  it('refuses a singular matrix', () => {
    expect(inverse(fromFlat([1, 2, 2, 4]))).toBeNull()
    expect(inverse(fromFlat([0, 0, 0, 0]))).toBeNull()
    expect(inverse(fromFlat([1, 0, 0, 0]))).toBeNull()
  })

  it('measures singularity against the matrix scale, not an absolute epsilon', () => {
    // A tiny matrix has a tiny determinant and a perfectly good inverse.
    const small = inverse(fromFlat([1e-6, 0, 0, 1e-6]))
    expect(small).not.toBeNull()
    close((small as Mat2)[0][0], 1e6, 1e-3)
    // A large one likewise.
    expect(inverse(fromFlat([1e6, 0, 0, 1e6]))).not.toBeNull()
    // Near-singular but above the relative threshold: invertible, badly.
    const nearly = inverse(fromFlat([1, 2, 2, 4.0001]))
    expect(nearly).not.toBeNull()
    expect(Number.isFinite((nearly as Mat2)[0][0])).toBe(true)
    // Exactly singular at its own scale is refused.
    expect(inverse(fromFlat([2, 4, 4, 8]))).toBeNull()
  })
})

describe('eigenvalues', () => {
  it('finds the real pair of a symmetric matrix and its directions', () => {
    const a: Mat2 = [
      [2, 1],
      [1, 2],
    ]
    const e = eigen(a)
    expect(e.kind).toBe('real')
    if (e.kind !== 'real') return
    close(e.values[0], 1)
    close(e.values[1], 3)
    expect(e.repeated).toBe(false)
    for (let i = 0; i < 2; i++) {
      const v = e.vectors[i]
      close(v[0] * v[0] + v[1] * v[1], 1)
      const av = apply(a, v)
      close(av[0], e.values[i] * v[0])
      close(av[1], e.values[i] * v[1])
    }
    // Canonical sign: the first component is non-negative.
    for (const v of e.vectors) expect(v[0]).toBeGreaterThanOrEqual(0)
  })

  it('flags a repeated eigenvalue and does not invent a second direction for a shear', () => {
    const shear = eigen(fromFlat([1, 1, 0, 1]))
    expect(shear.kind).toBe('real')
    if (shear.kind !== 'real') return
    expect(shear.repeated).toBe(true)
    close(shear.values[0], 1)
    close(shear.vectors[0][0], 1)
    close(shear.vectors[0][1], 0)
    // Both reported directions are the same unit vector, not a fabricated pair.
    close(shear.vectors[1][0], shear.vectors[0][0])
    close(shear.vectors[1][1], shear.vectors[0][1])

    const scalar = eigen(fromFlat([2, 0, 0, 2]))
    expect(scalar.kind).toBe('real')
    if (scalar.kind !== 'real') return
    expect(scalar.repeated).toBe(true)
    // Every direction is an eigenvector, so the engine returns no particular one.
    expect(scalar.vectors[0]).toEqual([0, 0])
  })

  it('reports a complex pair rather than a real direction', () => {
    const turn = eigen(rotation(Math.PI / 2))
    expect(turn.kind).toBe('complex')
    if (turn.kind !== 'complex') return
    close(turn.re, 0)
    close(turn.im, 1)
  })

  it('returns a unit vector per eigenvalue through the standalone helper', () => {
    const a = fromFlat([0, 1, 1, 0])
    const plus = eigenvector(a, 1)
    const minus = eigenvector(a, -1)
    close(plus[0], Math.SQRT1_2)
    close(plus[1], Math.SQRT1_2)
    close(minus[0], Math.SQRT1_2)
    close(minus[1], -Math.SQRT1_2)
    close(normalize(plus)[0], plus[0])
  })
})

describe('presets', () => {
  it('covers the six teaching cases with the expected determinants', () => {
    close(det(fromFlat(PRESETS.identity.matrix)), 1)
    close(det(fromFlat(PRESETS.rotate.matrix)), 1)
    close(det(fromFlat(PRESETS.scale.matrix)), 0.9)
    close(det(fromFlat(PRESETS.reflect.matrix)), -1)
    close(det(fromFlat(PRESETS.shear.matrix)), 1)
    close(det(fromFlat(PRESETS.project.matrix)), 0)
    expect(det(fromFlat(PRESETS.identity.matrix))).toBe(1)
    expect(identity()[0][0]).toBe(1)
  })

  it('keeps the singular preset honestly singular to the inverse test', () => {
    expect(inverse(fromFlat(PRESETS.project.matrix))).toBeNull()
    expect(EPS).toBeGreaterThan(0)
  })
})
