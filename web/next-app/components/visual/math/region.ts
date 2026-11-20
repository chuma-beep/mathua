// Three planar regions and the integral of a named density over each.
//
// The estimate is a sum over the centres of the cells that lie inside the
// region, which is exactly what the figure shades. The exact value is a closed
// form, so the readout can show the error of the picture rather than asserting
// that the picture is right.

export type RegionKind = 'rect' | 'between' | 'polar'
export type Density = 'one' | 'x' | 'y'

export const REGION_LABELS: Record<RegionKind, string> = {
  rect: 'Rectangle [0, 2] × [0, 1]',
  between: 'Between y = x² and y = x',
  polar: 'Quarter disk of radius 1',
}

export const REGION_EQUATIONS: Record<RegionKind, string> = {
  rect: '0 ≤ x ≤ 2, 0 ≤ y ≤ 1',
  between: '0 ≤ x ≤ 1, x² ≤ y ≤ x',
  polar: 'x² + y² ≤ 1, x ≥ 0, y ≥ 0',
}

export const DENSITY_LABELS: Record<Density, string> = {
  one: 'f = 1 (area)',
  x: 'f = x',
  y: 'f = y',
}

export interface Bounds2 {
  x0: number
  x1: number
  y0: number
  y1: number
}

export function regionBounds(kind: RegionKind): Bounds2 {
  switch (kind) {
    case 'rect':
      return { x0: 0, x1: 2, y0: 0, y1: 1 }
    case 'between':
      return { x0: 0, x1: 1, y0: 0, y1: 1 }
    case 'polar':
      return { x0: 0, x1: 1, y0: 0, y1: 1 }
  }
}

export function insideRegion(kind: RegionKind, x: number, y: number): boolean {
  switch (kind) {
    case 'rect':
      return x >= 0 && x <= 2 && y >= 0 && y <= 1
    case 'between':
      return x >= 0 && x <= 1 && y >= x * x && y <= x
    case 'polar':
      return x >= 0 && y >= 0 && x * x + y * y <= 1
  }
}

export function densityValue(density: Density, x: number, y: number): number {
  switch (density) {
    case 'one':
      return 1
    case 'x':
      return x
    case 'y':
      return y
  }
}

/**
 * Exact integral over the region, by closed form. Every case here is a
 * polynomial, so the value is exact and no sampling error can hide in it.
 */
export function exactIntegral(kind: RegionKind, density: Density): number {
  if (kind === 'rect') {
    // ∫₀² ∫₀¹ f dy dx
    if (density === 'one') return 2
    if (density === 'x') return 2
    return 1
  }
  if (kind === 'between') {
    // ∫₀¹ ∫_{x²}^{x} f dy dx
    if (density === 'one') return 1 / 6
    if (density === 'x') return 1 / 12
    return 1 / 15
  }
  // Quarter disk, polar coordinates x = r cos θ, y = r sin θ, r ∈ [0, 1].
  if (density === 'one') return Math.PI / 4
  return 1 / 3
}

export interface Cell {
  x: number
  y: number
  /** true when the cell centre lies inside, which is what the sum counts. */
  inside: boolean
}

/** A uniform grid over the region's bounding box, classified by its centres. */
export function regionCells(kind: RegionKind, n = 12): Cell[] {
  const bounds = regionBounds(kind)
  const size = Math.max(1, Math.floor(n))
  const dx = (bounds.x1 - bounds.x0) / size
  const dy = (bounds.y1 - bounds.y0) / size
  const out: Cell[] = []
  for (let j = 0; j < size; j++) {
    for (let i = 0; i < size; i++) {
      const x = bounds.x0 + (i + 0.5) * dx
      const y = bounds.y0 + (j + 0.5) * dy
      out.push({ x, y, inside: insideRegion(kind, x, y) })
    }
  }
  return out
}

/** Riemann sum over the cells whose centres are inside. */
export function regionSum(kind: RegionKind, density: Density, n = 12): number {
  const bounds = regionBounds(kind)
  const size = Math.max(1, Math.floor(n))
  const area = ((bounds.x1 - bounds.x0) / size) * ((bounds.y1 - bounds.y0) / size)
  return regionCells(kind, size).reduce(
    (total, cell) => (cell.inside ? total + densityValue(density, cell.x, cell.y) * area : total),
    0,
  )
}

/** The inside-cell count and the cell area, so the readout can name the estimate. */
export function cellStats(kind: RegionKind, n = 12): { count: number; area: number } {
  const bounds = regionBounds(kind)
  const size = Math.max(1, Math.floor(n))
  const count = regionCells(kind, size).filter((cell) => cell.inside).length
  return { count, area: ((bounds.x1 - bounds.x0) / size) * ((bounds.y1 - bounds.y0) / size) }
}

/** The same region read with the other variable first. */
export function reversedBounds(kind: RegionKind): string {
  switch (kind) {
    case 'rect':
      return '0 ≤ y ≤ 1, 0 ≤ x ≤ 2'
    case 'between':
      return '0 ≤ y ≤ 1, y ≤ x ≤ √y'
    case 'polar':
      return '0 ≤ θ ≤ π/2, 0 ≤ r ≤ 1'
  }
}

/** The area element, which is what the Jacobian contributes. */
export function areaElement(kind: RegionKind): string {
  return kind === 'polar' ? 'r dr dθ' : 'dx dy'
}
