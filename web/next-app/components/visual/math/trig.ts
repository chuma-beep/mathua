// The unit circle and the two coordinate functions that read it.
//
// This module owns the angle convention and the quadrant and reference-angle
// arithmetic. A template draws the result and never re-derives sine or cosine
// from a point it picked itself, so the circle and the graph cannot disagree.
//
// Angles are radians and turn counterclockwise from the positive x-axis.

export const TAU = 2 * Math.PI
export const DEGREES_PER_RADIAN = 180 / Math.PI
export const RADIANS_PER_DEGREE = Math.PI / 180

export type Quadrant = 'I' | 'II' | 'III' | 'IV' | 'axis'

export interface UnitPoint {
  /** The angle reduced to [0, 2π). */
  theta: number
  cos: number
  sin: number
  /** Null at an odd multiple of π/2, where the tangent is undefined. */
  tan: number | null
  quadrant: Quadrant
  /** The acute angle to the x-axis, in [0, π/2]. */
  referenceAngle: number
  /** (cos θ, sin θ). */
  point: [number, number]
}

export function degreesToRadians(degrees: number): number {
  return degrees * RADIANS_PER_DEGREE
}

export function radiansToDegrees(radians: number): number {
  return radians * DEGREES_PER_RADIAN
}

/** Reduce any finite angle to [0, 2π). A non-finite angle stays NaN. */
export function normalizeAngle(theta: number): number {
  if (!Number.isFinite(theta)) return Number.NaN
  const reduced = theta % TAU
  return reduced < 0 ? reduced + TAU : reduced
}

/**
 * A signed reduction to (−π, π], which is the branch a graph marker needs. The
 * caller can still normalize again for a [0, 2π) readout.
 */
export function normalizeSigned(theta: number): number {
  if (!Number.isFinite(theta)) return Number.NaN
  const reduced = normalizeAngle(theta)
  return reduced > Math.PI ? reduced - TAU : reduced
}

/** The acute angle between the terminal ray and the nearest x-axis. */
export function referenceAngle(theta: number): number {
  if (!Number.isFinite(theta)) return Number.NaN
  const reduced = Math.abs(normalizeSigned(theta))
  return reduced > Math.PI / 2 ? Math.PI - reduced : reduced
}

export function quadrantOf(theta: number): Quadrant {
  if (!Number.isFinite(theta)) return 'axis'
  const reduced = normalizeAngle(theta)
  const tolerance = 1e-12
  const onAxis =
    Math.abs(reduced) <= tolerance ||
    Math.abs(reduced - Math.PI / 2) <= tolerance ||
    Math.abs(reduced - Math.PI) <= tolerance ||
    Math.abs(reduced - (3 * Math.PI) / 2) <= tolerance ||
    Math.abs(reduced - TAU) <= tolerance
  if (onAxis) return 'axis'
  if (reduced < Math.PI / 2) return 'I'
  if (reduced < Math.PI) return 'II'
  if (reduced < (3 * Math.PI) / 2) return 'III'
  return 'IV'
}

/** Sine, cosine and tangent of an angle, with the quadrant facts attached. */
export function unitPoint(theta: number): UnitPoint {
  const reduced = normalizeAngle(theta)
  const cos = Math.cos(reduced)
  const sin = Math.sin(reduced)
  return {
    theta: reduced,
    cos,
    sin,
    tan: Math.abs(cos) <= 1e-12 ? null : sin / cos,
    quadrant: quadrantOf(reduced),
    referenceAngle: referenceAngle(reduced),
    point: [cos, sin],
  }
}

/**
 * Project a point onto the unit circle. This is the drag path: the pointer
 * chooses a direction and the point stays on the circle, so the readout is
 * always exactly (cos θ, sin θ).
 */
export function projectToUnitCircle(target: [number, number]): number {
  if (!Number.isFinite(target[0]) || !Number.isFinite(target[1])) return 0
  if (Math.abs(target[0]) <= 1e-15 && Math.abs(target[1]) <= 1e-15) return 0
  return normalizeAngle(Math.atan2(target[1], target[0]))
}

/** The identity sin²θ + cos²θ, kept as one value for a readout and a test. */
export function pythagorean(theta: number): number {
  const point = unitPoint(theta)
  return point.sin * point.sin + point.cos * point.cos
}
