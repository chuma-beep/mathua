// The one solid the 3D template draws, with its measurements.
//
// Both renderers consume this module and nothing else, which is what makes the
// SVG fallback and the WebGL scene show the same object rather than two
// lookalikes. It has no dependency on a renderer, so the geometry can be tested
// without a canvas.

export const X_MIN = 0
export const X_MAX = 4

/** Radius of the body of revolution at x. */
export function radius(x: number): number {
  return 1 + 0.35 * Math.sin(1.6 * x)
}

/** dr/dx, which fixes the surface normal. */
export function slope(x: number): number {
  return 0.35 * 1.6 * Math.cos(1.6 * x)
}

function integrate(f: (x: number) => number, slices = 400): number {
  const even = slices % 2 === 0 ? slices : slices + 1
  const h = (X_MAX - X_MIN) / even
  let sum = f(X_MIN) + f(X_MAX)
  for (let i = 1; i < even; i++) sum += f(X_MIN + i * h) * (i % 2 === 0 ? 2 : 4)
  return (sum * h) / 3
}

export function crossSectionArea(x: number): number {
  return Math.PI * radius(x) ** 2
}

/** Volume by disks: pi times the integral of r squared. */
export function volume(): number {
  return Math.PI * integrate((x) => radius(x) ** 2)
}

/** Lateral surface area: 2 pi times the integral of r ds. */
export function lateralArea(): number {
  return 2 * Math.PI * integrate((x) => radius(x) * Math.sqrt(1 + slope(x) ** 2))
}

export interface Mesh {
  positions: number[]
  normals: number[]
  indices: number[]
  rings: number
  segments: number
}

export interface Point3 {
  x: number
  y: number
  z: number
}

/** The body centred on the origin, as a triangle mesh with outward normals. */
export function surface(rings = 48, segments = 64): Mesh {
  const positions: number[] = []
  const normals: number[] = []
  const indices: number[] = []
  const half = (X_MAX - X_MIN) / 2

  for (let i = 0; i <= rings; i++) {
    const x = X_MIN + (i / rings) * (X_MAX - X_MIN)
    const cx = x - half
    const r = radius(x)
    const dr = slope(x)
    const len = Math.sqrt(1 + dr * dr)
    for (let j = 0; j <= segments; j++) {
      const angle = (2 * Math.PI * j) / segments
      positions.push(cx, r * Math.cos(angle), r * Math.sin(angle))
      normals.push(-dr / len, Math.cos(angle) / len, Math.sin(angle) / len)
    }
  }

  const stride = segments + 1
  for (let i = 0; i < rings; i++) {
    for (let j = 0; j < segments; j++) {
      const a = i * stride + j
      const b = a + stride
      indices.push(a, b, b + 1, a, b + 1, a + 1)
    }
  }

  return { positions, normals, indices, rings, segments }
}

/** The profile as (x, r) pairs, for a wireframe or an outline. */
export function profile(samples = 24): Array<[number, number]> {
  const out: Array<[number, number]> = []
  for (let i = 0; i <= samples; i++) {
    const x = X_MIN + (i / samples) * (X_MAX - X_MIN)
    out.push([x, radius(x)])
  }
  return out
}

/** Centre of the body along x, so the camera can look at the origin. */
export const CENTER_X = (X_MIN + X_MAX) / 2
