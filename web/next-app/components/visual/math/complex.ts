// Complex numbers as a pair of reals. No dependency on the Compute Engine: the
// visual needs magnitude, argument and the product, not symbolic equivalence.

export type Complex = [number, number]

export function cadd(a: Complex, b: Complex): Complex {
  return [a[0] + b[0], a[1] + b[1]]
}

export function cmul(a: Complex, b: Complex): Complex {
  return [a[0] * b[0] - a[1] * b[1], a[0] * b[1] + a[1] * b[0]]
}

export function cmod(a: Complex): number {
  return Math.hypot(a[0], a[1])
}

/** Argument in radians, in (-pi, pi]. */
export function carg(a: Complex): number {
  return Math.atan2(a[1], a[0])
}

/**
 * The same angle reduced into (-pi, pi].
 *
 * Two arguments added can land outside the principal range, which is exactly what
 * happens when the product wraps past the negative real axis. Printing the raw sum
 * beside the product's argument would show two different numbers for one angle.
 */
export function wrapAngle(angle: number): number {
  const twoPi = 2 * Math.PI
  const wrapped = (((angle + Math.PI) % twoPi) + twoPi) % twoPi - Math.PI
  return wrapped === -Math.PI ? Math.PI : wrapped
}

/**
 * The argument, or null at the origin. atan2(0, 0) is 0, which would claim the zero
 * vector points along the positive real axis; it has no direction at all.
 */
export function cargOrNull(a: Complex): number | null {
  return cmod(a) < 1e-12 ? null : Math.atan2(a[1], a[0])
}
