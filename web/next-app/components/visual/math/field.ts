// The four fields the contour figure offers.
//
// A closed list rather than an expression for the same reason the polynomial
// figures take coefficients: the gradient is then an exact formula, so the arrow
// drawn beside a level curve is the real gradient and not a finite difference
// that could point slightly wrong.

export type Field = 'bowl' | 'saddle' | 'product' | 'ripple'

export const FIELD_LABELS: Record<Field, string> = {
  bowl: 'x² + y²',
  saddle: 'x² − y²',
  product: 'x · y',
  ripple: 'sin x + cos y',
}

export function fieldValue(field: Field, x: number, y: number): number {
  switch (field) {
    case 'bowl':
      return x * x + y * y
    case 'saddle':
      return x * x - y * y
    case 'product':
      return x * y
    case 'ripple':
      return Math.sin(x) + Math.cos(y)
  }
}

/** The exact gradient, which is perpendicular to the level curve through the point. */
export function fieldGradient(field: Field, x: number, y: number): [number, number] {
  switch (field) {
    case 'bowl':
      return [2 * x, 2 * y]
    case 'saddle':
      return [2 * x, -2 * y]
    case 'product':
      return [y, x]
    case 'ripple':
      return [Math.cos(x), -Math.sin(y)]
  }
}

export function gradientLength(field: Field, x: number, y: number): number {
  const [gx, gy] = fieldGradient(field, x, y)
  return Math.hypot(gx, gy)
}

export interface Hessian {
  xx: number
  xy: number
  yx: number
  yy: number
}

/** The exact second derivatives, so the saddle test is not a finite difference. */
export function fieldHessian(field: Field, x: number, y: number): Hessian {
  switch (field) {
    case 'bowl':
      return { xx: 2, xy: 0, yx: 0, yy: 2 }
    case 'saddle':
      return { xx: 2, xy: 0, yx: 0, yy: -2 }
    case 'product':
      return { xx: 0, xy: 1, yx: 1, yy: 0 }
    case 'ripple':
      return { xx: -Math.sin(x), xy: 0, yx: 0, yy: -Math.cos(y) }
  }
}

/** The stationary points a figure may put its probe on. */
export function fieldCritical(field: Field): Array<[number, number]> {
  switch (field) {
    case 'bowl':
    case 'saddle':
    case 'product':
      return [[0, 0]]
    case 'ripple':
      return [
        [Math.PI / 2, 0],
        [-Math.PI / 2, 0],
      ]
  }
}

export interface TangentPlane {
  z0: number
  fx: number
  fy: number
  /** z = z0 + fx (x − x0) + fy (y − y0) */
  at: (x: number, y: number) => number
  equation: string
}

/** The plane that matches the surface's value and both slopes at (x0, y0). */
export function tangentPlane(field: Field, x0: number, y0: number): TangentPlane {
  const z0 = fieldValue(field, x0, y0)
  const [fx, fy] = fieldGradient(field, x0, y0)
  return {
    z0,
    fx,
    fy,
    at: (x, y) => z0 + fx * (x - x0) + fy * (y - y0),
    equation: `z = ${z0.toFixed(2)} + ${fx.toFixed(2)}(x − ${x0.toFixed(2)}) + ${fy.toFixed(2)}(y − ${y0.toFixed(2)})`,
  }
}

/** Rate of change in the direction (cos θ, sin θ). */
export function directional(field: Field, x: number, y: number, angle: number): number {
  const [gx, gy] = fieldGradient(field, x, y)
  return gx * Math.cos(angle) + gy * Math.sin(angle)
}

export interface SaddleTest {
  discriminant: number
  xx: number
  classification: 'minimum' | 'maximum' | 'saddle' | 'degenerate'
}

/** The second-derivative test at a point. */
export function saddleTest(field: Field, x: number, y: number): SaddleTest {
  const hessian = fieldHessian(field, x, y)
  const discriminant = hessian.xx * hessian.yy - hessian.xy * hessian.yx
  let classification: SaddleTest['classification'] = 'degenerate'
  if (discriminant > 1e-12) classification = hessian.xx > 0 ? 'minimum' : 'maximum'
  else if (discriminant < -1e-12) classification = 'saddle'
  return { discriminant, xx: hessian.xx, classification }
}

/**
 * The levels to draw, spread between the field's range over the window.
 *
 * The lowest level is deliberately one step above the minimum: for a bowl the
 * minimum is a single point, and a contour through one point draws nothing.
 */
export function fieldLevels(field: Field, extent: number, count: number): number[] {
  const samples = 24
  let lo = Number.POSITIVE_INFINITY
  let hi = Number.NEGATIVE_INFINITY
  for (let i = 0; i <= samples; i++) {
    for (let j = 0; j <= samples; j++) {
      const x = -extent + (2 * extent * i) / samples
      const y = -extent + (2 * extent * j) / samples
      const v = fieldValue(field, x, y)
      if (!Number.isFinite(v)) continue
      lo = Math.min(lo, v)
      hi = Math.max(hi, v)
    }
  }
  if (!Number.isFinite(lo) || !Number.isFinite(hi) || hi <= lo) return []
  const step = (hi - lo) / (count + 1)
  return Array.from({ length: count }, (_, i) => Number((lo + step * (i + 1)).toFixed(6)))
}
