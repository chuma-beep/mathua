// Complex maps and the one property that separates the analytic ones.
//
// A map is conformal at z when it preserves the angle between two curves
// through z. For an analytic map that is what "multiply the tangent by f'(z)"
// does: every direction is rotated by the same argument, so the angle between
// any two is unchanged. Conjugation rotates nothing and reflects instead, which
// turns an angle into its negative. That comparison is the figure's whole point,
// so the four maps here are closed and their derivatives are exact.

import type { Complex } from './complex'

export type ComplexMap = 'square' | 'exp' | 'inverse' | 'conjugate'

export const MAP_LABELS: Record<ComplexMap, string> = {
  square: 'z²',
  exp: 'e^z',
  inverse: '1/z',
  conjugate: 'conj z',
}

export function mapValue(map: ComplexMap, z: Complex): Complex {
  const [x, y] = z
  switch (map) {
    case 'square':
      return [x * x - y * y, 2 * x * y]
    case 'exp': {
      const r = Math.exp(x)
      return [r * Math.cos(y), r * Math.sin(y)]
    }
    case 'inverse': {
      const d = x * x + y * y
      if (d === 0) return [Number.NaN, Number.NaN]
      return [x / d, -y / d]
    }
    case 'conjugate':
      return [x, -y]
  }
}

/** f'(z), or null where the map has no derivative in the complex sense. */
export function mapDerivative(map: ComplexMap, z: Complex): Complex | null {
  const [x, y] = z
  switch (map) {
    case 'square':
      return [2 * x, 2 * y]
    case 'exp': {
      const r = Math.exp(x)
      return [r * Math.cos(y), r * Math.sin(y)]
    }
    case 'inverse': {
      const d = x * x + y * y
      if (d === 0) return null
      // d/dz (1/z) = -1/z^2 = -(conj z)^2 / |z|^4
      return [-(x * x - y * y) / (d * d), (2 * x * y) / (d * d)]
    }
    case 'conjugate':
      return null
  }
}

export function mapIsAnalytic(map: ComplexMap): boolean {
  return map !== 'conjugate'
}

/** Wrap to (-pi, pi], so a printed angle never reads as 400 degrees. */
export function wrapAngle(angle: number): number {
  let out = angle
  while (out > Math.PI) out -= 2 * Math.PI
  while (out <= -Math.PI) out += 2 * Math.PI
  return out
}

/**
 * True where f has a derivative and it is not zero: the point at which angles
 * are preserved. False at a pole, where the map has no value at all, and at a
 * critical point, where directions are not merely rotated.
 */
export function mapIsConformalAt(map: ComplexMap, z: Complex): boolean {
  if (!mapIsAnalytic(map)) return false
  const d = mapDerivative(map, z)
  return d !== null && Math.hypot(d[0], d[1]) > 1e-9
}

/**
 * How much an analytic map multiplies angles at z, or null where it has no value.
 *
 * A first derivative of zero is the case the derivative's argument cannot
 * express: z squared sends a direction to twice its argument at the origin, so
 * the angle between two curves doubles rather than surviving. The closed set
 * here has exactly one such point, which is why the factor is 2 for `square`
 * and 1 everywhere else.
 */
export function angleFactorAt(map: ComplexMap, z: Complex): number | null {
  if (!mapIsAnalytic(map)) return 1
  const d = mapDerivative(map, z)
  if (!d) return null
  if (Math.hypot(d[0], d[1]) > 1e-9) return 1
  return map === 'square' ? 2 : 1
}

/**
 * Where a starting direction points after the map.
 *
 * For an analytic map every direction is rotated by arg f'(z), so the offset is
 * additive. Conjugation leaves nothing to rotate and reflects instead, which is
 * why its angle is negated rather than shifted.
 */
export function mappedDirectionAngle(map: ComplexMap, z: Complex, angle: number): number {
  if (!mapIsAnalytic(map)) return wrapAngle(-angle)
  const factor = angleFactorAt(map, z)
  if (factor === null) return Number.NaN
  const d = mapDerivative(map, z)
  // At a critical point the rotation is undefined and the angle is multiplied
  // instead; adding arg(0) would have reported that nothing happened.
  if (!d || Math.hypot(d[0], d[1]) <= 1e-9) return wrapAngle(angle * factor)
  return wrapAngle(angle + Math.atan2(d[1], d[0]))
}

/** The angle between two directions, in radians, always reported as a positive size. */
export function angleBetweenDirections(a: number, b: number): number {
  return Math.abs(wrapAngle(a - b))
}
