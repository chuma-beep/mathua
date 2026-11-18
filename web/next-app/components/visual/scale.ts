// Coordinate and path helpers shared by the templates. Templates draw in SVG user
// units, so these functions convert a data domain to the viewBox and back without
// each template re-deriving the same arithmetic.

export function clamp(value: number, lo: number, hi: number): number {
  return Math.min(hi, Math.max(lo, value))
}

export function scaleLinear(d0: number, d1: number, r0: number, r1: number): (v: number) => number {
  const span = d1 - d0
  if (span === 0) return () => (r0 + r1) / 2
  return (v: number) => r0 + ((v - d0) / span) * (r1 - r0)
}

export function invertLinear(d0: number, d1: number, r0: number, r1: number): (p: number) => number {
  const span = r1 - r0
  if (span === 0) return () => d0
  return (p: number) => d0 + ((p - r0) / span) * (d1 - d0)
}

/** Evenly spaced ticks covering [min, max], including both ends where they fit. */
export function ticks(min: number, max: number, count = 4): number[] {
  if (!Number.isFinite(min) || !Number.isFinite(max) || count < 1) return []
  const step = (max - min) / count
  return Array.from({ length: count + 1 }, (_, i) => min + step * i)
}

/**
 * Ticks at "nice" numeric steps, so an axis can read 0, 1, 2 rather than
 * 0.37, 0.74, 1.11. `target` is the number of intervals wanted; the step is
 * rounded up to 1, 2 or 5 times a power of ten.
 */
export function niceTicks(min: number, max: number, target = 5): number[] {
  if (!Number.isFinite(min) || !Number.isFinite(max) || max <= min || target < 1) return []
  const raw = (max - min) / target
  const magnitude = Math.pow(10, Math.floor(Math.log10(raw)))
  const step = [1, 2, 5, 10].map((m) => m * magnitude).find((s) => s >= raw) ?? 10 * magnitude
  const out: number[] = []
  for (let value = Math.ceil(min / step) * step; value <= max + step * 1e-6; value += step) {
    out.push(Math.abs(value) < step * 1e-6 ? 0 : Number(value.toFixed(10)))
  }
  return out
}

export function pathFromPoints(points: Array<[number, number]>): string {
  if (points.length === 0) return ''
  return points.map((p, i) => `${i === 0 ? 'M' : 'L'} ${p[0].toFixed(2)} ${p[1].toFixed(2)}`).join(' ')
}

/** A pointer position in the SVG's own user units, not client pixels. */
export function clientToSvg(svg: SVGSVGElement, clientX: number, clientY: number): { x: number; y: number } {
  const rect = svg.getBoundingClientRect()
  const viewBox = svg.viewBox?.baseVal
  if (!viewBox || viewBox.width === 0 || viewBox.height === 0 || rect.width === 0 || rect.height === 0) {
    return { x: clientX - rect.left, y: clientY - rect.top }
  }
  return {
    x: viewBox.x + ((clientX - rect.left) / rect.width) * viewBox.width,
    y: viewBox.y + ((clientY - rect.top) / rect.height) * viewBox.height,
  }
}

/** An SVG arc from the positive x-axis to a vector's argument, for angle readouts. */
export function arcPath(cx: number, cy: number, radius: number, endAngle: number): string {
  const start = { x: cx + radius, y: cy }
  const end = { x: cx + radius * Math.cos(-endAngle), y: cy + radius * Math.sin(-endAngle) }
  const large = Math.abs(endAngle) > Math.PI ? 1 : 0
  const sweep = endAngle >= 0 ? 0 : 1
  return `M ${start.x} ${start.y} A ${radius} ${radius} 0 ${large} ${sweep} ${end.x} ${end.y}`
}
