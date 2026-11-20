// A sinusoid y = A f(B(x − C)) + D, kept apart from the template.
//
// The graph, the readout and the lesson predictions all read the same
// parameters here. A template must not recompute a period from a screen
// coordinate, because then the curve and the numbers could disagree.

import { TAU } from './trig'

export type SinusoidForm = 'sin' | 'cos'

export interface SinusoidParams {
  form: SinusoidForm
  /** Signed amplitude. A negative value reflects the wave. */
  amplitude: number
  /** Angular frequency B in f(B(x − C)). Always positive in the figure. */
  frequency: number
  /** Horizontal shift C, in radians. */
  phase: number
  /** Vertical shift D, which is also the midline. */
  vertical: number
}

export interface Sinusoid {
  form: SinusoidForm
  amplitude: number
  /** |A|, the distance from the midline to a peak. */
  peak: number
  frequency: number
  period: number
  phase: number
  vertical: number
  max: number
  min: number
  midline: number
}

function finite(value: number): value is number {
  return Number.isFinite(value)
}

/** Build the derived quantities once. Invalid input yields NaN, never a throw. */
export function sinusoid(params: SinusoidParams): Sinusoid {
  const { form, amplitude, frequency, phase, vertical } = params
  if (![amplitude, frequency, phase, vertical].every(finite)) {
    return {
      form,
      amplitude: Number.NaN,
      peak: Number.NaN,
      frequency: Number.NaN,
      period: Number.NaN,
      phase: Number.NaN,
      vertical: Number.NaN,
      max: Number.NaN,
      min: Number.NaN,
      midline: Number.NaN,
    }
  }
  const peak = Math.abs(amplitude)
  const speed = Math.abs(frequency)
  return {
    form,
    amplitude,
    peak,
    frequency,
    period: speed <= 1e-12 ? Number.POSITIVE_INFINITY : TAU / speed,
    phase,
    vertical,
    max: vertical + peak,
    min: vertical - peak,
    midline: vertical,
  }
}

/** The value of the wave at x, using the same parameters as the readout. */
export function sinusoidValue(wave: Sinusoid, x: number): number {
  if (!Number.isFinite(x) || !Number.isFinite(wave.amplitude)) return Number.NaN
  const inside = wave.frequency * (x - wave.phase)
  const base = wave.form === 'sin' ? Math.sin(inside) : Math.cos(inside)
  return wave.amplitude * base + wave.vertical
}

const gcd = (a: number, b: number): number => (b === 0 ? Math.abs(a) : gcd(b, a % b))

/**
 * Format a multiple of π as a fraction. Values that are not a small rational
 * multiple of π fall back to a decimal multiple, so a slider position still
 * reads as `0.37π` rather than as an unexplained decimal.
 */
export function formatPi(value: number, digits = 2): string {
  if (!Number.isFinite(value)) return '—'
  if (Math.abs(value) < 1e-10) return '0'
  const ratio = value / Math.PI
  const denominators = [1, 2, 3, 4, 6, 8, 12]
  for (const denominator of denominators) {
    const numerator = Math.round(ratio * denominator)
    if (Math.abs(ratio - numerator / denominator) > 1e-9) continue
    const divisor = gcd(numerator, denominator)
    const n = numerator / divisor
    const d = denominator / divisor
    const sign = n < 0 ? '−' : ''
    const top = Math.abs(n)
    const head = top === 1 ? 'π' : `${top}π`
    return d === 1 ? `${sign}${head}` : `${sign}${head}/${d}`
  }
  return `${value < 0 ? '−' : ''}${Math.abs(ratio).toFixed(digits)}π`
}

/** A human-readable equation for the figure and the readout. */
export function sinusoidFormula(wave: Sinusoid): string {
  if (!Number.isFinite(wave.amplitude) || !Number.isFinite(wave.frequency)) {
    return 'y = —'
  }
  if (Math.abs(wave.amplitude) <= 1e-12) return `y = ${wave.vertical}`
  const sign = wave.amplitude < 0 ? '−' : ''
  const size = Math.abs(wave.amplitude)
  const coefficient = size === 1 ? '' : `${size}`
  const shift = Math.abs(wave.phase) < 1e-10
    ? 'x'
    : `x ${wave.phase > 0 ? '−' : '+'} ${formatPi(Math.abs(wave.phase))}`
  const argument = wave.frequency === 1
    ? shift
    : Math.abs(wave.phase) < 1e-10
      ? `${wave.frequency}x`
      : `${wave.frequency}(${shift})`
  const vertical = Math.abs(wave.vertical) < 1e-10
    ? ''
    : ` ${wave.vertical > 0 ? '+' : '−'} ${Math.abs(wave.vertical)}`
  return `y = ${sign}${coefficient}${wave.form}(${argument})${vertical}`
}

/** The first peak and trough in one period, useful for a test and a label. */
export function sinusoidExtremes(wave: Sinusoid): { peaks: number[]; troughs: number[] } {
  if (!Number.isFinite(wave.period)) return { peaks: [], troughs: [] }
  const base = wave.form === 'sin' ? Math.PI / 2 : 0
  const peakAt = (k: number) => wave.phase + (base + k * TAU) / wave.frequency
  const peaks = [peakAt(0), peakAt(1)]
  const troughs = peaks.map((value) => value + wave.period / 2)
  return { peaks, troughs }
}
