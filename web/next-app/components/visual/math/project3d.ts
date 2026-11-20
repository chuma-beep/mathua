// One axonometric projection for every figure that draws a solid or a surface.
//
// The two figures that need depth, the tangent plane and the solid of
// revolution, must agree about what "towards the viewer" means, or a slice drawn
// in one would sit behind the body drawn in the other. Keeping the projection
// here also means the painter's order is computed once and can be tested without
// a canvas.

export interface Vec3 {
  x: number
  y: number
  z: number
}

export interface View3D {
  /** Rotation about the vertical axis, in radians. */
  theta: number
  /** Tilt of the horizontal plane, in radians. */
  phi: number
}

export interface ScreenPoint {
  x: number
  /** Upward is positive, so a template drawing into SVG negates this. */
  y: number
  /** Distance toward the viewer; larger is nearer. */
  depth: number
}

export function project(point: Vec3, view: View3D): ScreenPoint {
  const ct = Math.cos(view.theta)
  const st = Math.sin(view.theta)
  const xr = point.x * ct - point.y * st
  const yr = point.x * st + point.y * ct
  const cp = Math.cos(view.phi)
  const sp = Math.sin(view.phi)
  return {
    x: xr,
    y: point.z * cp - yr * sp,
    depth: yr * cp + point.z * sp,
  }
}

/** Painter's order: far quads first, near quads last. */
export function byDepth<T extends { depth: number }>(items: T[]): T[] {
  return [...items].sort((a, b) => a.depth - b.depth)
}

/** Four shading steps from the same four theme accents, never a raw colour. */
export type ShadeIndex = 0 | 1 | 2 | 3

export function shadeIndex(value: number, min: number, max: number): ShadeIndex {
  if (!Number.isFinite(value) || max <= min) return 1
  const t = Math.min(1, Math.max(0, (value - min) / (max - min)))
  const index = Math.min(3, Math.floor(t * 4)) as ShadeIndex
  return index
}

export const SHADES = ['var(--chart-1)', 'var(--chart-2)', 'var(--chart-3)', 'var(--accent-blue)'] as const
