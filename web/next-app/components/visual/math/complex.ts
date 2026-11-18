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
