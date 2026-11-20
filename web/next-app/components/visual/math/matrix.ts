// 2x2 matrices, row-major. Pure functions with no React and no rendering, so the
// arithmetic can be tested without a DOM.

export type Vec2 = [number, number]
export type Mat2 = [[number, number], [number, number]]

export function apply(m: Mat2, v: Vec2): Vec2 {
  return [m[0][0] * v[0] + m[0][1] * v[1], m[1][0] * v[0] + m[1][1] * v[1]]
}

/** a * b, in the order a is applied after b. */
export function mul(a: Mat2, b: Mat2): Mat2 {
  return [
    [a[0][0] * b[0][0] + a[0][1] * b[1][0], a[0][0] * b[0][1] + a[0][1] * b[1][1]],
    [a[1][0] * b[0][0] + a[1][1] * b[1][0], a[1][0] * b[0][1] + a[1][1] * b[1][1]],
  ]
}

export function det(m: Mat2): number {
  return m[0][0] * m[1][1] - m[0][1] * m[1][0]
}

export function identity(): Mat2 {
  return [
    [1, 0],
    [0, 1],
  ]
}

export function rotation(radians: number): Mat2 {
  const c = Math.cos(radians)
  const s = Math.sin(radians)
  return [
    [c, -s],
    [s, c],
  ]
}

export function reflection(axis: 'x' | 'y' | 'origin'): Mat2 {
  if (axis === 'x') {
    return [
      [1, 0],
      [0, -1],
    ]
  }
  if (axis === 'y') {
    return [
      [-1, 0],
      [0, 1],
    ]
  }
  return [
    [-1, 0],
    [0, -1],
  ]
}

export function dilation(k: number): Mat2 {
  return [
    [k, 0],
    [0, k],
  ]
}

/** Applies the arguments right to left, as function composition reads. */
export function compose(...matrices: Mat2[]): Mat2 {
  return matrices.reduceRight((acc, m) => mul(m, acc), identity())
}

/** The shard stores four numbers, row-major, because a nested tuple is noisy JSON. */
export function fromFlat(flat: [number, number, number, number]): Mat2 {
  return [
    [flat[0], flat[1]],
    [flat[2], flat[3]],
  ]
}

export function toFlat(m: Mat2): [number, number, number, number] {
  return [m[0][0], m[0][1], m[1][0], m[1][1]]
}

export function dot(a: Vec2, b: Vec2): number {
  return a[0] * b[0] + a[1] * b[1]
}

export function length(v: Vec2): number {
  return Math.hypot(v[0], v[1])
}

/** A unit vector in the same direction. The zero vector is returned unchanged. */
export function normalize(v: Vec2): Vec2 {
  const len = length(v)
  return len === 0 ? [0, 0] : [v[0] / len, v[1] / len]
}

export function scaleVec(v: Vec2, k: number): Vec2 {
  return [v[0] * k, v[1] * k]
}

export function trace(m: Mat2): number {
  return m[0][0] + m[1][1]
}

/** The unsigned angle between two vectors, in radians, in [0, pi]. */
export function angleBetween(a: Vec2, b: Vec2): number {
  const la = length(a)
  const lb = length(b)
  if (la === 0 || lb === 0) return 0
  const cos = Math.min(1, Math.max(-1, dot(a, b) / (la * lb)))
  return Math.acos(cos)
}

export function addVec(a: Vec2, b: Vec2): Vec2 {
  return [a[0] + b[0], a[1] + b[1]]
}

export function subVec(a: Vec2, b: Vec2): Vec2 {
  return [a[0] - b[0], a[1] - b[1]]
}

/** Relative tolerance for a determinant that is zero to working precision. */
export const EPS = 1e-10

/**
 * M(t) = (1 - t) I + t A, entry by entry, with t clamped to [0, 1].
 *
 * This is the straight line between the identity and A in matrix space, which is
 * the only interpolation a linear map can show without inventing a path: every
 * intermediate M(t) is itself a matrix, so every intermediate picture is a real
 * transformation of the plane rather than a cross-fade between two drawings.
 */
export function lerpMatrix(from: Mat2, to: Mat2, t: number): Mat2 {
  const s = Math.min(1, Math.max(0, t))
  const mix = (x: number, y: number) => (1 - s) * x + s * y
  return [
    [mix(from[0][0], to[0][0]), mix(from[0][1], to[0][1])],
    [mix(from[1][0], to[1][0]), mix(from[1][1], to[1][1])],
  ]
}

/** The largest absolute entry, which sets the scale of the matrix. */
export function matrixScale(m: Mat2): number {
  return Math.max(...m.flat().map(Math.abs))
}

/**
 * How the image of the unit square should be drawn, which is not the same
 * question as whether the matrix is invertible.
 *
 * Three different facts were being collapsed into one in an earlier figure: a
 * determinant that is exactly zero, a determinant that is tiny relative to the
 * entries, and a parallelogram whose thinnest dimension is below a pixel at the
 * current drawing scale. Only the first is singular. The third is invisible, and
 * calling it "area 0" would be a lie about a matrix that inverts perfectly well.
 */
export type CollapseState = 'none' | 'near' | 'pixel' | 'singular'

/**
 * Classify an image at a given drawing scale. `pxPerUnit` is how many pixels one
 * unit of the plane occupies, so a caller that zooms passes its own scale and the
 * pixel verdict follows the drawing rather than a fixed number.
 */
export function collapseState(m: Mat2, pxPerUnit: number): CollapseState {
  const scale = matrixScale(m)
  if (scale === 0) return 'singular'
  const area = Math.abs(det(m))
  // The same relative convention `inverse` uses: |det| compared against scale².
  if (area <= EPS * scale * scale) return 'singular'

  const [first, second] = columns(m)
  const longest = Math.max(length(first), length(second))
  if (longest === 0) return 'singular'
  // Height of the parallelogram taken from its longest side, in pixels.
  const heightPx = (area / longest) * pxPerUnit
  if (heightPx < 0.5) return 'pixel'
  if (area <= 1e-3 * scale * scale) return 'near'
  return 'none'
}

export function describeCollapse(state: CollapseState, area: number): string {
  switch (state) {
    case 'singular':
      return 'collapses to a line (det = 0)'
    case 'near':
      return `nearly collapsed, det = ${area.toExponential(1)}`
    case 'pixel':
      return `thinner than a pixel here, det = ${area.toExponential(1)}`
    case 'none':
      return 'full parallelogram'
  }
}

/**
 * The columns of a matrix, which are the images of the basis vectors. Derived
 * through `apply` rather than copied out, so a column can never disagree with
 * what the figure draws for `A e₁` and `A e₂`.
 */
export function columns(m: Mat2): [Vec2, Vec2] {
  return [apply(m, [1, 0]), apply(m, [0, 1])]
}

/**
 * The inverse, or null when the matrix is singular relative to its own scale.
 *
 * The threshold is relative, not absolute: `[[1e-6,0],[0,1e-6]]` has a tiny
 * determinant and a perfectly good inverse, while `[[1,2],[2,4]]` has a
 * determinant of zero at its scale and none. Comparing |det| against
 * EPS * max|entry|² keeps both readings correct.
 */
export function inverse(m: Mat2): Mat2 | null {
  const d = det(m)
  // No floor of 1 here: a floor would make the test absolute again, and a tiny
  // matrix with a tiny determinant would be refused although it inverts fine.
  const scale = Math.max(...m.flat().map(Math.abs))
  if (scale === 0 || Math.abs(d) <= EPS * scale * scale) return null
  return [
    [m[1][1] / d, -m[0][1] / d],
    [-m[1][0] / d, m[0][0] / d],
  ]
}

export type Eigen =
  | { kind: 'real'; values: [number, number]; vectors: [Vec2, Vec2]; repeated: boolean }
  | { kind: 'complex'; re: number; im: number }

/**
 * Real or complex eigenvalues of a 2x2 matrix.
 *
 * A negative discriminant is reported as a complex conjugate pair rather than
 * forced into a real direction, and a repeated eigenvalue is flagged because a
 * shear and a scalar matrix both have one yet only the scalar matrix has more
 * than one independent direction. Rendering may underline either case; it may
 * not invent a second vector.
 */
export function eigen(m: Mat2): Eigen {
  const t = trace(m)
  const d = det(m)
  const disc = t * t - 4 * d
  const scale = Math.max(...m.flat().map(Math.abs))
  const tol = scale === 0 ? 0 : EPS * scale * scale
  if (disc < -tol) {
    return { kind: 'complex', re: t / 2, im: Math.sqrt(-disc) / 2 }
  }
  const root = Math.sqrt(Math.max(0, disc))
  const l1 = (t - root) / 2
  const l2 = (t + root) / 2
  return {
    kind: 'real',
    values: [l1, l2],
    vectors: [eigenvector(m, l1), eigenvector(m, l2)],
    repeated: Math.abs(l2 - l1) <= Math.max(tol, EPS),
  }
}

/**
 * A unit eigenvector for the given eigenvalue, or a zero vector when the matrix
 * has no direction for it (a scalar matrix has every direction).
 *
 * The sign is canonicalised so the first component is non-negative, because an
 * eigenvector is defined up to sign and a learner should not see the arrow flip
 * when a coefficient changes in the fourth decimal place.
 */
export function eigenvector(m: Mat2, lambda: number): Vec2 {
  const rows: Vec2[] = [
    [m[0][0] - lambda, m[0][1]],
    [m[1][0], m[1][1] - lambda],
  ]
  const row = length(rows[0]) >= length(rows[1]) ? rows[0] : rows[1]
  const v: Vec2 = [row[1], -row[0]]
  const unit = normalize(v)
  if (unit[0] < 0 || (unit[0] === 0 && unit[1] < 0)) return [-unit[0], -unit[1]]
  return unit
}

export type Preset = 'identity' | 'rotate' | 'scale' | 'reflect' | 'shear' | 'project'

/** Flats are row-major `[a, b, c, d]`, the shape a shard stores. */
export const PRESETS: Record<Preset, { label: string; matrix: [number, number, number, number] }> = {
  identity: { label: 'Identity', matrix: [1, 0, 0, 1] },
  rotate: { label: 'Quarter turn', matrix: [0, -1, 1, 0] },
  scale: { label: 'Scale', matrix: [1.5, 0, 0, 0.6] },
  reflect: { label: 'Reflection', matrix: [1, 0, 0, -1] },
  shear: { label: 'Shear', matrix: [1, 1, 0, 1] },
  project: { label: 'Projection (singular)', matrix: [1, 0, 0, 0] },
}
