// A small closed family of complex functions, with their partial derivatives
// written out rather than approximated.
//
// The figure this serves draws the level curves of Re f and Im f, which cross at
// right angles wherever f' is not zero. That claim is only worth showing if every
// number behind it is exact, so `cfGradients` returns the four partials of u and v
// as closed forms: no finite differences, and the Cauchy-Riemann identities can be
// asserted in a test rather than hoped for in a picture.
//
// The two degenerate cases are part of the contract, not an afterthought. At a
// critical point both gradients vanish, so "the curves are perpendicular" has no
// tangent to be true of; at a zero of f the argument is undefined. Both are
// reported as undefined by the template instead of being painted as a right angle.

export type ComplexFn = 'square' | 'cube' | 'exp' | 'inverse' | 'sine'

export const CF_LABELS: Record<ComplexFn, string> = {
  square: 'z²',
  cube: 'z³',
  exp: 'e^z',
  inverse: '1/z',
  sine: 'sin z',
}

export type Vec2 = [number, number]

/** f(z), or null where the function has no finite value. */
export function cfValue(fn: ComplexFn, z: Vec2): Vec2 | null {
  const [x, y] = z
  switch (fn) {
    case 'square':
      return [x * x - y * y, 2 * x * y]
    case 'cube':
      return [x * x * x - 3 * x * y * y, 3 * x * x * y - y * y * y]
    case 'exp':
      return [Math.exp(x) * Math.cos(y), Math.exp(x) * Math.sin(y)]
    case 'inverse': {
      const d = x * x + y * y
      if (d < 1e-12) return null
      return [x / d, -y / d]
    }
    case 'sine':
      return [Math.sin(x) * Math.cosh(y), Math.cos(x) * Math.sinh(y)]
  }
}

/** f'(z), or null where the derivative does not exist. */
export function cfDerivative(fn: ComplexFn, z: Vec2): Vec2 | null {
  const [x, y] = z
  switch (fn) {
    case 'square':
      return [2 * x, 2 * y]
    case 'cube':
      return [3 * (x * x - y * y), 6 * x * y]
    case 'exp':
      return [Math.exp(x) * Math.cos(y), Math.exp(x) * Math.sin(y)]
    case 'inverse': {
      const d = x * x + y * y
      if (d < 1e-12) return null
      return [-(x * x - y * y) / (d * d), (-2 * x * y) / (d * d)]
    }
    case 'sine':
      return [Math.cos(x) * Math.cosh(y), -Math.sin(x) * Math.sinh(y)]
  }
}

/**
 * The four partial derivatives of f = u + iv. `du` is (u_x, u_y) and `dv` is
 * (v_x, v_y). Analyticity of every function here means u_x = v_y and u_y = -v_x,
 * which is exactly what makes the two level families perpendicular.
 */
export function cfGradients(fn: ComplexFn, z: Vec2): { du: Vec2; dv: Vec2 } | null {
  const [x, y] = z
  switch (fn) {
    case 'square':
      return { du: [2 * x, -2 * y], dv: [2 * y, 2 * x] }
    case 'cube':
      return { du: [3 * x * x - 3 * y * y, -6 * x * y], dv: [6 * x * y, 3 * x * x - 3 * y * y] }
    case 'exp': {
      const e = Math.exp(x)
      const c = Math.cos(y)
      const s = Math.sin(y)
      return { du: [e * c, -e * s], dv: [e * s, e * c] }
    }
    case 'inverse': {
      const d = x * x + y * y
      if (d < 1e-12) return null
      const a = (y * y - x * x) / (d * d)
      const b = (-2 * x * y) / (d * d)
      return { du: [a, b], dv: [-b, a] }
    }
    case 'sine':
      return {
        du: [Math.cos(x) * Math.cosh(y), Math.sin(x) * Math.sinh(y)],
        dv: [-Math.sin(x) * Math.sinh(y), Math.cos(x) * Math.cosh(y)],
      }
  }
}

/**
 * Where f is zero, so `arg f` has no value to print. A pole is not a zero:
 * 1/z never vanishes, it is undefined at 0, which is why that point is listed
 * by `cfPoles` instead. Reporting it as a zero drew a "zero" ring on the point
 * where the function has no value at all.
 */
export function cfZeros(fn: ComplexFn): Vec2[] {
  switch (fn) {
    case 'square':
    case 'cube':
      return [[0, 0]]
    case 'inverse':
    case 'exp':
      return []
    case 'sine':
      return [
        [-Math.PI, 0],
        [0, 0],
        [Math.PI, 0],
      ]
  }
}

/** Where f is undefined rather than zero. Only 1/z has one in this family. */
export function cfPoles(fn: ComplexFn): Vec2[] {
  return fn === 'inverse' ? [[0, 0]] : []
}

/** Where f' is zero: the level curves have no tangent there, so no angle either. */
export function cfCritical(fn: ComplexFn): Vec2[] {
  switch (fn) {
    case 'square':
    case 'cube':
      return [[0, 0]]
    case 'exp':
    case 'inverse':
      return []
    case 'sine':
      return [
        [-Math.PI / 2, 0],
        [Math.PI / 2, 0],
      ]
  }
}

/** A level is drawn above the minimum for the same reason a bowl needs no contour at its bottom. */
export function cfLevelValues(fn: ComplexFn, part: 'u' | 'v', extent: number, count: number): number[] {
  let lo = Number.POSITIVE_INFINITY
  let hi = Number.NEGATIVE_INFINITY
  for (let i = 0; i <= 24; i++) {
    for (let j = 0; j <= 24; j++) {
      const x = -extent + (2 * extent * i) / 24
      const y = -extent + (2 * extent * j) / 24
      const g = cfGradients(fn, [x, y])
      if (!g) continue
      const f = cfValue(fn, [x, y])
      if (!f) continue
      const value = part === 'u' ? f[0] : f[1]
      if (!Number.isFinite(value)) continue
      lo = Math.min(lo, value)
      hi = Math.max(hi, value)
    }
  }
  if (!Number.isFinite(lo) || !Number.isFinite(hi) || hi <= lo) return []
  const step = (hi - lo) / (count + 1)
  return Array.from({ length: count }, (_, i) => lo + step * (i + 1))
}
