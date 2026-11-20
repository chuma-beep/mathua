// Number-theory calculations that do not depend on React or SVG. Every result
// here is exact integer arithmetic; the renderer only chooses how to show it.

export interface Congruence {
  quotient: number
  residue: number
}

export interface ModPowerStep {
  exponent: number
  base: number
  result: number
}

export interface EuclidStep {
  a: number
  b: number
  quotient: number
  remainder: number
}

export interface Factorization {
  n: number
  factors: number[]
  tree: string[]
}

export function modValue(a: number, n: number): number {
  if (!Number.isFinite(a) || !Number.isFinite(n) || n <= 0) return Number.NaN
  const r = Math.trunc(a) % Math.trunc(n)
  return r < 0 ? r + Math.trunc(n) : r
}

export function congruence(a: number, n: number): Congruence | null {
  if (!Number.isFinite(a) || !Number.isFinite(n) || n <= 0) return null
  const residue = modValue(a, n)
  return { quotient: Math.floor((Math.trunc(a) - residue) / Math.trunc(n)), residue }
}

/** Repeated squaring, recording the running result after each bit. */
export function modPow(base: number, exponent: number, n: number): { value: number; steps: ModPowerStep[] } | null {
  if (!Number.isFinite(base) || !Number.isFinite(exponent) || !Number.isFinite(n)) return null
  if (n <= 0 || exponent < 0 || !Number.isInteger(exponent)) return null
  let result = 1
  let b = modValue(base, n)
  let e = Math.trunc(exponent)
  const steps: ModPowerStep[] = []
  while (e > 0) {
    if (e & 1) result = (result * b) % Math.trunc(n)
    steps.push({ exponent: e, base: b, result })
    b = (b * b) % Math.trunc(n)
    e = Math.floor(e / 2)
  }
  return { value: result, steps }
}

export function euclid(a: number, b: number): { gcd: number; steps: EuclidStep[] } | null {
  if (!Number.isFinite(a) || !Number.isFinite(b)) return null
  let x = Math.abs(Math.trunc(a))
  let y = Math.abs(Math.trunc(b))
  if (x === 0 && y === 0) return null
  const steps: EuclidStep[] = []
  while (y !== 0) {
    const quotient = Math.floor(x / y)
    const remainder = x - quotient * y
    steps.push({ a: x, b: y, quotient, remainder })
    x = y
    y = remainder
  }
  return { gcd: x, steps }
}

function factorList(n: number): number[] {
  const factors: number[] = []
  let value = n
  for (let p = 2; p * p <= value; p++) {
    while (value % p === 0) {
      factors.push(p)
      value /= p
    }
  }
  if (value > 1) factors.push(value)
  return factors
}

export function factorize(n: number): Factorization | null {
  if (!Number.isFinite(n) || !Number.isInteger(n) || n < 2) return null
  const factors = factorList(n)
  const tree: string[] = []
  let remaining = n
  for (const factor of factors) {
    tree.push(`${remaining} = ${factor} × ${remaining / factor}`)
    remaining /= factor
  }
  return { n, factors, tree }
}

export function divisors(n: number): number[] {
  if (!Number.isFinite(n) || !Number.isInteger(n) || n < 1) return []
  const out: number[] = []
  for (let d = 1; d * d <= n; d++) {
    if (n % d === 0) {
      out.push(d)
      if (d * d !== n) out.push(n / d)
    }
  }
  return out.sort((a, b) => a - b)
}
