// Discrete mathematics: truth tables and finite set operations.
//
// Both are exact over a finite universe, so the figure never has to approximate
// a count. The universe for sets is fixed and printed in the readout.

export type TruthOp = 'and' | 'or' | 'implies' | 'xor'

export interface TruthRow {
  p: boolean
  q: boolean
  value: boolean
}

export interface TruthTable {
  op: TruthOp
  expression: string
  rows: TruthRow[]
  trueCount: number
}

const TRUTH_EXPRESSION: Record<TruthOp, string> = {
  and: 'P and Q',
  or: 'P or Q',
  implies: 'P implies Q',
  xor: 'P xor Q',
}

export function truthTable(op: TruthOp): TruthTable {
  const rows: TruthRow[] = []
  for (const p of [true, false]) {
    for (const q of [true, false]) {
      const value =
        op === 'and' ? p && q : op === 'or' ? p || q : op === 'implies' ? !p || q : p !== q
      rows.push({ p, q, value })
    }
  }
  return {
    op,
    expression: TRUTH_EXPRESSION[op],
    rows,
    trueCount: rows.filter((row) => row.value).length,
  }
}

export type SetOp = 'union' | 'intersection' | 'difference' | 'symmetric'

export interface SetOperation {
  op: SetOp
  universe: number[]
  left: number[]
  right: number[]
  result: number[]
  count: number
}

/**
 * A and B are bitmasks over the universe 1..6. The same masks drive both the
 * labels and the membership grid, so the figure cannot disagree with the count.
 */
export function setOperation(op: SetOp, a: number, b: number): SetOperation {
  const universe = [1, 2, 3, 4, 5, 6]
  const left = universe.filter((value) => (a & (1 << (value - 1))) !== 0)
  const right = universe.filter((value) => (b & (1 << (value - 1))) !== 0)
  const result = universe.filter((value) => {
    const inA = (a & (1 << (value - 1))) !== 0
    const inB = (b & (1 << (value - 1))) !== 0
    if (op === 'union') return inA || inB
    if (op === 'intersection') return inA && inB
    if (op === 'difference') return inA && !inB
    return inA !== inB
  })
  return { op, universe, left, right, result, count: result.length }
}

export function maskOf(values: number[]): number {
  return values.reduce((mask, value) => mask | (1 << (value - 1)), 0)
}

/** n choose k, exact for the small integers the figure offers. */
export function binomial(n: number, k: number): number {
  if (!Number.isInteger(n) || !Number.isInteger(k) || n < 0 || k < 0) return 0
  if (k > n) return 0
  const r = Math.min(k, n - k)
  let value = 1
  for (let i = 1; i <= r; i++) value = (value * (n - r + i)) / i
  return Math.round(value)
}

export function combinations(n: number, k: number): number[][] {
  if (!Number.isInteger(n) || !Number.isInteger(k) || n < 0 || k < 0 || k > n) return []
  const result: number[][] = []
  const build = (start: number, chosen: number[]) => {
    if (chosen.length === k) {
      result.push([...chosen])
      return
    }
    for (let value = start; value <= n; value++) build(value + 1, [...chosen, value])
  }
  build(1, [])
  return result
}
