// Arithmetic models that a learner can point at: a signed line, tenths, a
// percentage bar, a ratio tape and a two-pan balance. Every value is computed
// here so the drawing, the readout and the lesson prediction cannot disagree.

export type ArithmeticMode = 'signed' | 'decimal' | 'percent' | 'ratio' | 'balance'

export interface SignedComparison {
  order: -1 | 0 | 1
  left: number
  right: number
  sentence: string
}

export interface PercentChange {
  amount: number
  result: number
  direction: 'increase' | 'decrease' | 'none'
}

export interface RatioForm {
  a: number
  b: number
  simplifiedA: number
  simplifiedB: number
  sentence: string
}

export interface BalanceState {
  guess: number
  b: number
  c: number
  left: number
  right: number
  balanced: boolean
  solution: number
}

const gcd = (a: number, b: number): number => {
  let x = Math.abs(Math.trunc(a))
  let y = Math.abs(Math.trunc(b))
  while (y !== 0) {
    const next = x % y
    x = y
    y = next
  }
  return x || 1
}

export function compareSigned(a: number, b: number): SignedComparison {
  const order = a < b ? -1 : a > b ? 1 : 0
  const sentence =
    order === 0 ? `${a} and ${b} are the same point` : order < 0 ? `${a} is left of ${b}` : `${a} is right of ${b}`
  return { order, left: Math.min(a, b), right: Math.max(a, b), sentence }
}

export function tenthsValue(n: number): number {
  return Math.trunc(n) / 10
}

export function percentChange(percent: number, base: number): PercentChange {
  const amount = (percent / 100) * base
  if (amount > 0) return { amount, result: base + amount, direction: 'increase' }
  if (amount < 0) return { amount, result: base + amount, direction: 'decrease' }
  return { amount: 0, result: base, direction: 'none' }
}

export function ratioForm(a: number, b: number): RatioForm {
  const g = gcd(a, b)
  const simplifiedA = Math.trunc(a) / g
  const simplifiedB = Math.trunc(b) / g
  return {
    a: Math.trunc(a),
    b: Math.trunc(b),
    simplifiedA,
    simplifiedB,
    sentence: `${Math.trunc(a)} : ${Math.trunc(b)} simplifies to ${simplifiedA} : ${simplifiedB}`,
  }
}

export function balanceState(guess: number, b: number, c: number): BalanceState {
  const left = guess + b
  return {
    guess,
    b,
    c,
    left,
    right: c,
    balanced: left === c,
    solution: c - b,
  }
}
