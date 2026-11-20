// Coordinates are a claim about a basis, not a property of the point.
//
// The figure needs both directions: build the point from its coordinates in the
// chosen basis, and read the coordinates back out. The second is the inverse, and
// it exists only when the two basis vectors are independent, so the null return is
// part of the mathematics rather than an error path.

import { det, type Mat2, type Vec2 } from './matrix'

/** The matrix whose columns are the basis vectors. */
export function basisMatrix(u: Vec2, v: Vec2): Mat2 {
  return [
    [u[0], v[0]],
    [u[1], v[1]],
  ]
}

/** alpha u + beta v, which is the matrix product with the coordinates as a column. */
export function combine(basis: Mat2, coords: Vec2): Vec2 {
  return [
    basis[0][0] * coords[0] + basis[0][1] * coords[1],
    basis[1][0] * coords[0] + basis[1][1] * coords[1],
  ]
}

/**
 * The coordinates of a point in the given basis, or null when the two columns are
 * dependent. A null answer is honest: a dependent pair does not span the plane, so
 * some points have no coordinates and others have infinitely many.
 */
export function coordinatesOf(basis: Mat2, point: Vec2): Vec2 | null {
  const d = det(basis)
  if (Math.abs(d) < 1e-9) return null
  return [
    (basis[1][1] * point[0] - basis[0][1] * point[1]) / d,
    (-basis[1][0] * point[0] + basis[0][0] * point[1]) / d,
  ]
}

/**
 * Whether the two columns are too close to parallel for the coordinates to be
 * usable, even though they are not exactly dependent.
 *
 * The test is relative rather than absolute: det scales with the lengths of u and
 * v, so a fixed threshold would call a long, well-spread basis singular. Dividing
 * by |u||v| leaves sin of the angle between them, which is the quantity that
 * actually decides how badly the inverse behaves near dependence.
 */
export function isNearlyDependent(basis: Mat2, tolerance = 1e-3): boolean {
  const u: Vec2 = [basis[0][0], basis[1][0]]
  const v: Vec2 = [basis[0][1], basis[1][1]]
  const scale = Math.hypot(u[0], u[1]) * Math.hypot(v[0], v[1])
  if (scale === 0) return true
  return Math.abs(det(basis)) / scale <= tolerance
}
