// Shared projection for the two renderers that do not have a 3D engine.
//
// A point (x, y, z) in model space becomes a screen point by rotating around the
// y-axis by `theta`, then around the x-axis by `phi`, then dividing by a depth
// term. The depth is kept because the painter's sort needs it and because a
// fallback that drew in the wrong order would look like a bug in the fallback.

export interface Vec3 {
  x: number
  y: number
  z: number
}

export interface Camera {
  theta: number
  phi: number
  distance: number
  scale: number
  cx: number
  cy: number
}

export interface Projected {
  x: number
  y: number
  depth: number
}

export function project(point: Vec3, camera: Camera): Projected {
  const cosT = Math.cos(camera.theta)
  const sinT = Math.sin(camera.theta)
  const cosP = Math.cos(camera.phi)
  const sinP = Math.sin(camera.phi)

  const x1 = point.x * cosT - point.z * sinT
  const z1 = point.x * sinT + point.z * cosT
  const y2 = point.y * cosP - z1 * sinP
  const z2 = point.y * sinP + z1 * cosP

  const depth = camera.distance - z2
  const factor = (camera.scale * camera.distance) / Math.max(0.2, depth)
  return { x: camera.cx + x1 * factor, y: camera.cy - y2 * factor, depth }
}

/** The surface as quads in model space, with the outward normal of each. */
export interface Quad {
  points: [Vec3, Vec3, Vec3, Vec3]
  normal: Vec3
}

export function quads(rings: number, segments: number, radiusAt: (x: number) => number, slopeAt: (x: number) => number, xMin: number, xMax: number): Quad[] {
  const out: Quad[] = []
  const at = (i: number, j: number): Vec3 => {
    const x = xMin + ((xMax - xMin) * i) / rings
    const theta = (2 * Math.PI * j) / segments
    const r = radiusAt(x)
    return { x: x - (xMin + xMax) / 2, y: r * Math.cos(theta), z: r * Math.sin(theta) }
  }

  for (let i = 0; i < rings; i++) {
    const x = xMin + ((xMax - xMin) * (i + 0.5)) / rings
    const dr = slopeAt(x)
    const nLength = Math.hypot(1, dr)
    for (let j = 0; j < segments; j++) {
      const theta = (2 * Math.PI * (j + 0.5)) / segments
      out.push({
        points: [at(i, j), at(i + 1, j), at(i + 1, j + 1), at(i, j + 1)],
        normal: { x: -dr / nLength, y: Math.cos(theta) / nLength, z: Math.sin(theta) / nLength },
      })
    }
  }
  return out
}

const LIGHT: Vec3 = { x: 0.45, y: 0.75, z: 0.5 }

/** Flat shading from a fixed light, returned as a 0-1 brightness. */
export function shade(normal: Vec3): number {
  const dot = normal.x * LIGHT.x + normal.y * LIGHT.y + normal.z * LIGHT.z
  return Math.max(0, Math.min(1, 0.28 + 0.72 * dot))
}
