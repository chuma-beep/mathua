// Four planar vector fields, with their divergence, curl and loop integrals.
//
// A closed list rather than a parsed expression, for the reason the contour
// figure takes named fields: the divergence and the curl are then exact
// constants, so the number beside Green's theorem is the same number the arrows
// were drawn from. A finite difference would let the readout disagree with the
// picture by a rounding error and nobody could tell which was wrong.

export type VectorFieldPreset = 'source' | 'vortex' | 'saddle' | 'shear'

export const VF_LABELS: Record<VectorFieldPreset, string> = {
  source: 'F = (x, y)',
  vortex: 'F = (−y, x)',
  saddle: 'F = (x, −y)',
  shear: 'F = (y, 0)',
}

export const VF_NAMES: Record<VectorFieldPreset, string> = {
  source: 'Source',
  vortex: 'Vortex',
  saddle: 'Saddle',
  shear: 'Shear',
}

export function vfValue(preset: VectorFieldPreset, x: number, y: number): [number, number] {
  switch (preset) {
    case 'source':
      return [x, y]
    case 'vortex':
      return [-y, x]
    case 'saddle':
      return [x, -y]
    case 'shear':
      return [y, 0]
  }
}

/** Exact divergence of the preset, constant over the plane. */
export function vfDivergence(preset: VectorFieldPreset): number {
  switch (preset) {
    case 'source':
      return 2
    case 'vortex':
      return 0
    case 'saddle':
      return 0
    case 'shear':
      return 0
  }
}

/** Exact scalar curl dQ/dx − dP/dy, constant over the plane. */
export function vfCurl(preset: VectorFieldPreset): number {
  switch (preset) {
    case 'source':
      return 0
    case 'vortex':
      return 2
    case 'saddle':
      return 0
    case 'shear':
      return -1
  }
}

/** A circle of the given radius as a closed path. */
export function loopPath(radius: number, count = 64): Array<[number, number]> {
  const out: Array<[number, number]> = []
  for (let index = 0; index < count; index++) {
    const angle = (2 * Math.PI * index) / count
    out.push([radius * Math.cos(angle), radius * Math.sin(angle)])
  }
  // The last point repeats the first exactly, so the loop is closed rather than
  // closed to within sin(2π).
  out.push([radius, 0])
  return out
}

/** Circulation ∮ F · dr around the circle, by the exact formula. */
export function circulation(preset: VectorFieldPreset, radius: number): number {
  switch (preset) {
    case 'source':
    case 'saddle':
      return 0
    case 'vortex':
      return 2 * Math.PI * radius * radius
    case 'shear':
      return -Math.PI * radius * radius
  }
}

/** Outward flux ∮ F · n ds around the circle, by the exact formula. */
export function flux(preset: VectorFieldPreset, radius: number): number {
  switch (preset) {
    case 'source':
      return 2 * Math.PI * radius * radius
    case 'vortex':
    case 'saddle':
    case 'shear':
      return 0
  }
}

/** The area of the region the loop encloses, used by the Green's theorem readout. */
export function loopArea(radius: number): number {
  return Math.PI * radius * radius
}

/** One streamline, integrated with RK4 from a seed in both time directions. */
export function streamline(
  preset: VectorFieldPreset,
  start: [number, number],
  duration = 2.4,
  steps = 90,
): Array<[number, number]> {
  if (!start.every(Number.isFinite) || !Number.isFinite(duration) || steps < 1) return []
  const dt = duration / Math.floor(steps)
  const points: Array<[number, number]> = [[start[0], start[1]]]
  let [x, y] = start
  for (let index = 0; index < steps; index++) {
    const k1 = vfValue(preset, x, y)
    const k2 = vfValue(preset, x + (dt * k1[0]) / 2, y + (dt * k1[1]) / 2)
    const k3 = vfValue(preset, x + (dt * k2[0]) / 2, y + (dt * k2[1]) / 2)
    const k4 = vfValue(preset, x + dt * k3[0], y + dt * k3[1])
    x += (dt / 6) * (k1[0] + 2 * k2[0] + 2 * k3[0] + k4[0])
    y += (dt / 6) * (k1[1] + 2 * k2[1] + 2 * k3[1] + k4[1])
    if (!Number.isFinite(x) || !Number.isFinite(y)) break
    points.push([x, y])
    if (Math.abs(x) > 6 || Math.abs(y) > 6) break
  }
  return points
}

/**
 * Streamlines from two radii and six angles. Two radii are enough because a
 * rotational field draws the same circle from every angle at one radius, and a
 * third ring would only repeat the loop the integrals use.
 */
export function streamlines(preset: VectorFieldPreset): Array<Array<[number, number]>> {
  const out: Array<Array<[number, number]>> = []
  for (const radius of [0.55, 1.75]) {
    for (let index = 0; index < 6; index++) {
      const angle = (2 * Math.PI * index) / 6 + 0.3
      out.push(streamline(preset, [radius * Math.cos(angle), radius * Math.sin(angle)]))
    }
  }
  return out.filter((line) => line.length > 1)
}

/**
 * Arrow length for the quiver plot. A constant field would otherwise draw every
 * arrow the same length and hide where the magnitude changes, so the length is
 * the field's magnitude clipped to a readable window rather than the raw value.
 */
export function arrowLength(preset: VectorFieldPreset, x: number, y: number, max = 5): number {
  const [ux, uy] = vfValue(preset, x, y)
  const magnitude = Math.hypot(ux, uy)
  return Math.min(1, magnitude / max)
}
