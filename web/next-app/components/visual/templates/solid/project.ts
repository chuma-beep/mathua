// The camera and projection the SVG fallback needs, and nothing else.
//
// It is deliberately separate from the Three.js renderer: the fallback must work
// where there is no WebGL context, so it cannot import the scene that needs one.

import type { Point3 } from '../../math/solid'

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

/** Rotate about the vertical axis, tilt, then divide by depth. */
export function project(point: Point3, camera: Camera): Projected {
  const ct = Math.cos(camera.theta)
  const st = Math.sin(camera.theta)
  const x1 = point.x * ct + point.z * st
  const z1 = -point.x * st + point.z * ct
  const cp = Math.cos(camera.phi)
  const sp = Math.sin(camera.phi)
  const y1 = point.y * cp - z1 * sp
  const z2 = point.y * sp + z1 * cp
  const depth = Math.max(0.2, camera.distance - z2)
  const k = camera.distance / depth
  return { x: camera.cx + x1 * camera.scale * k, y: camera.cy - y1 * camera.scale * k, depth: z2 }
}

/** A closed SVG path for a ring of points, already projected. */
export function ringPath(points: Projected[]): string {
  if (points.length === 0) return ''
  return points.map((p, i) => `${i === 0 ? 'M' : 'L'} ${p.x.toFixed(2)} ${p.y.toFixed(2)}`).join(' ') + ' Z'
}

/** An open SVG path through projected points. */
export function linePath(points: Projected[]): string {
  if (points.length === 0) return ''
  return points.map((p, i) => `${i === 0 ? 'M' : 'L'} ${p.x.toFixed(2)} ${p.y.toFixed(2)}`).join(' ')
}
