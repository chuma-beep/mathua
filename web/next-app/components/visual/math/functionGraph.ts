// One family of transformed functions, with the facts a graph should report.
//
// The inverse of a quadratic is not a function on its full domain, while the
// inverse of a cubic or an exponential is. The engine records that distinction
// instead of drawing a reflected curve and calling it a function.

export type FunctionFamily = 'quadratic' | 'cubic' | 'exponential'

export interface FamilySpec {
  id: FunctionFamily
  label: string
  expression: string
}

export const FUNCTION_FAMILIES: FamilySpec[] = [
  { id: 'quadratic', label: 'Quadratic', expression: 'x²' },
  { id: 'cubic', label: 'Cubic', expression: 'x³' },
  { id: 'exponential', label: 'Exponential', expression: '2ˣ' },
]

const BASE = {
  quadratic: (x: number) => x * x,
  cubic: (x: number) => x * x * x,
  exponential: (x: number) => Math.pow(2, x),
}

export function functionValue(
  family: FunctionFamily,
  a: number,
  h: number,
  k: number,
  x: number,
): number | null {
  if (![a, h, k, x].every(Number.isFinite)) return null
  const value = a * BASE[family](x - h) + k
  return Number.isFinite(value) ? value : null
}

export function functionFormula(family: FunctionFamily, a: number, h: number, k: number): string {
  const inner = h === 0 ? 'x' : h > 0 ? `(x − ${h})` : `(x + ${Math.abs(h)})`
  const lead = a === 1 ? '' : a === -1 ? '−' : `${a}`
  const lift = k === 0 ? '' : k > 0 ? ` + ${k}` : ` − ${Math.abs(k)}`
  if (family === 'quadratic') return `${lead}${inner}²${lift}`
  if (family === 'cubic') return `${lead}${inner}³${lift}`
  return `${a === 1 ? '' : a === -1 ? '−' : `${a}·`}2^${inner}${lift}`
}

export function inverseIsFunction(family: FunctionFamily, a: number): boolean {
  return family !== 'quadratic' && a !== 0
}

/** The inverse as a formula in y, when it is a function. */
export function inverseFormula(family: FunctionFamily, a: number, h: number, k: number): string | null {
  if (!inverseIsFunction(family, a)) return null
  const inner = k === 0 ? 'y' : k > 0 ? `(y − ${k})` : `(y + ${Math.abs(k)})`
  const scaled = a === 1 ? inner : a === -1 ? `−${inner}` : `${inner} / ${a}`
  if (family === 'cubic') return `${h} + ∛(${scaled})`
  return `${h} + log₂(${scaled})`
}

export function endBehavior(
  family: FunctionFamily,
  a: number,
): { left: string; right: string } {
  if (family === 'quadratic') {
    return a >= 0 ? { left: '+∞', right: '+∞' } : { left: '−∞', right: '−∞' }
  }
  if (family === 'cubic') {
    return a >= 0 ? { left: '−∞', right: '+∞' } : { left: '+∞', right: '−∞' }
  }
  return a >= 0 ? { left: '→ 0', right: '+∞' } : { left: '→ 0', right: '−∞' }
}

export function domainOf(family: FunctionFamily): string {
  return family === 'exponential' ? 'all real x' : 'all real x'
}

export function rangeOf(family: FunctionFamily, a: number, k: number): string {
  if (family === 'quadratic') {
    return a > 0 ? `[${k}, ∞)` : a < 0 ? `(−∞, ${k}]` : `{${k}}`
  }
  if (family === 'cubic') return 'all real y'
  return a > 0 ? `(${k}, ∞)` : a < 0 ? `(−∞, ${k})` : `{${k}}`
}
