// The solid every prototype draws: the surface of revolution of a curve about
// the x-axis, over [0, 4], with r(x) = 1 + 0.35 sin(1.6x).
//
// One profile, three renderers. A prototype that drew a different solid per
// renderer would compare pictures rather than renderers, so the geometry and the
// numbers live here and each prototype only decides how to put them on screen.

export const X_MIN = 0
export const X_MAX = 4

/** The distance from the axis at x. Always positive, so the surface never self-intersects. */
export function radius(x: number): number {
  return 1 + 0.35 * Math.sin(1.6 * x)
}

/** dr/dx, needed for the lateral area and for surface normals. */
export function slope(x: number): number {
  return 0.35 * 1.6 * Math.cos(1.6 * x)
}

/** Simpson's rule, because the profile has no elementary antiderivative worth writing down twice. */
function integrate(f: (x: number) => number, slices = 400): number {
  const n = slices % 2 === 0 ? slices : slices + 1
  const h = (X_MAX - X_MIN) / n
  let sum = f(X_MIN) + f(X_MAX)
  for (let i = 1; i < n; i++) sum += f(X_MIN + i * h) * (i % 2 === 0 ? 2 : 4)
  return (sum * h) / 3
}

/** The volume, by disks: pi times the integral of r squared. */
export function volume(): number {
  return Math.PI * integrate((x) => radius(x) ** 2)
}

/** The curved surface area: 2 pi r ds with ds = sqrt(1 + r'^2) dx. */
export function lateralArea(): number {
  return 2 * Math.PI * integrate((x) => radius(x) * Math.sqrt(1 + slope(x) ** 2))
}

/** The area of the disk a plane x = constant cuts through the solid. */
export function crossSectionArea(x: number): number {
  return Math.PI * radius(x) ** 2
}

export interface Mesh {
  /** Unit positions, indexed as quads: (nx + 1) rings of (nu + 1) vertices. */
  positions: Float32Array
  normals: Float32Array
  indices: Uint32Array
  rings: number
  segments: number
}

/**
 * A triangulated surface of revolution. `rings` runs along x and `segments`
 * around it; the seam is duplicated so the normals do not average across it.
 */
export function surface(rings = 48, segments = 64): Mesh {
  const positions: number[] = []
  const normals: number[] = []
  const indices: number[] = []

  for (let i = 0; i <= rings; i++) {
    const x = X_MIN + ((X_MAX - X_MIN) * i) / rings
    const r = radius(x)
    const dr = slope(x)
    // The outward normal of a surface of revolution: (1, -dr) in the (x, r)
    // plane, rotated around the axis.
    const nLength = Math.hypot(1, dr)
    for (let j = 0; j <= segments; j++) {
      const theta = (2 * Math.PI * j) / segments
      const cos = Math.cos(theta)
      const sin = Math.sin(theta)
      positions.push(x - (X_MAX - X_MIN) / 2, r * cos, r * sin)
      normals.push(-dr / nLength, cos / nLength, sin / nLength)
    }
  }

  const stride = segments + 1
  for (let i = 0; i < rings; i++) {
    for (let j = 0; j < segments; j++) {
      const a = i * stride + j
      const b = a + stride
      indices.push(a, b, a + 1, b, b + 1, a + 1)
    }
  }

  return {
    positions: new Float32Array(positions),
    normals: new Float32Array(normals),
    indices: new Uint32Array(indices),
    rings,
    segments,
  }
}

/** The profile as points, for the 2D renderers. */
export function profile(samples = 64): Array<[number, number]> {
  const out: Array<[number, number]> = []
  for (let i = 0; i <= samples; i++) {
    const x = X_MIN + ((X_MAX - X_MIN) * i) / samples
    out.push([x - (X_MAX - X_MIN) / 2, radius(x)])
  }
  return out
}
