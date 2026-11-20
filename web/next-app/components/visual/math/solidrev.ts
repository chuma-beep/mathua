// Three bodies of revolution, each with the slice its method integrates.
//
// The radius, the slice volume and the total volume come from the same
// functions, so the annulus drawn at x and the number in the readout cannot
// disagree. The integrals are Simpson sums of a polynomial, which is exact to
// machine precision at this resolution.

export type SolidMode = 'disk' | 'washer' | 'shell'

export const SOLID_LABELS: Record<SolidMode, string> = {
  disk: 'Disk',
  washer: 'Washer',
  shell: 'Cylindrical shell',
}

export const SOLID_FORMULAS: Record<SolidMode, string> = {
  disk: 'V = π ∫ R(x)² dx',
  washer: 'V = π ∫ (R(x)² − r(x)²) dx',
  shell: 'V = 2π ∫ x f(x) dx',
}

/** The interval every body is drawn on. */
export const SOLID_X0 = 0
export const SOLID_X1 = 2

/** Outer radius, shared by disk and washer. */
export function outerRadius(x: number): number {
  return 1 + 0.2 * x
}

/** Inner radius of the washer, the hole the axis leaves. */
export function innerRadius(x: number): number {
  return 0.4
}

/** Profile of the shell body, which rotates about the vertical axis. */
export function shellHeight(x: number): number {
  return 1 + 0.2 * x
}

function simpson(f: (x: number) => number, a: number, b: number, slices = 400): number {
  const even = slices % 2 === 0 ? slices : slices + 1
  const h = (b - a) / even
  let sum = f(a) + f(b)
  for (let index = 1; index < even; index++) {
    sum += f(a + index * h) * (index % 2 === 0 ? 2 : 4)
  }
  return (sum * h) / 3
}

/** Exact volume of the chosen body. */
export function solidVolume(mode: SolidMode): number {
  if (mode === 'disk') return Math.PI * simpson((x) => outerRadius(x) ** 2, SOLID_X0, SOLID_X1)
  if (mode === 'washer') {
    return (
      Math.PI *
      simpson((x) => outerRadius(x) ** 2 - innerRadius(x) ** 2, SOLID_X0, SOLID_X1)
    )
  }
  return 2 * Math.PI * simpson((x) => x * shellHeight(x), SOLID_X0, SOLID_X1)
}

/**
 * Volume of the slice at x, for a slice of thickness dx. A disk is a thin
 * cylinder, a washer is a thin annulus and a shell is a thin tube.
 */
export function sliceVolume(mode: SolidMode, x: number, dx: number): number {
  if (mode === 'disk') return Math.PI * outerRadius(x) ** 2 * dx
  if (mode === 'washer') {
    return Math.PI * (outerRadius(x) ** 2 - innerRadius(x) ** 2) * dx
  }
  return 2 * Math.PI * x * shellHeight(x) * dx
}

/** Area of the annulus or disk at x, which the readout names separately. */
export function sliceArea(mode: SolidMode, x: number): number {
  if (mode === 'disk') return Math.PI * outerRadius(x) ** 2
  if (mode === 'washer') {
    return Math.PI * (outerRadius(x) ** 2 - innerRadius(x) ** 2)
  }
  return 2 * Math.PI * x * shellHeight(x)
}

/** The lateral area of the disk body, as a second measurement to compare. */
export function lateralArea(): number {
  const slope = 0.2
  return 2 * Math.PI * simpson((x) => outerRadius(x) * Math.sqrt(1 + slope * slope), SOLID_X0, SOLID_X1)
}

/** The profile the body is swept from, as (x, radius) pairs. */
export function profile(mode: SolidMode, samples = 32): Array<[number, number]> {
  const out: Array<[number, number]> = []
  for (let index = 0; index <= samples; index++) {
    const x = SOLID_X0 + ((SOLID_X1 - SOLID_X0) * index) / samples
    out.push([x, mode === 'shell' ? shellHeight(x) : outerRadius(x)])
  }
  return out
}
