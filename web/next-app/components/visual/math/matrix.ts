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
