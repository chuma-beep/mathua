// Gaussian elimination on a 2×3 augmented matrix, recorded operation by
// operation.
//
// This is deliberately not a linear map: `R2 ← R2 − 3R1` rewrites one equation,
// and the reason it is allowed is that it preserves the solution set rather than
// the plane. Each step therefore carries a description, because the description
// is the lesson: a learner should be able to say what changed and why the answer
// did not.

export type Row = [number, number, number]

export type Step = {
  /** The augmented matrix after the operation. */
  matrix: Row[]
  /** What was done, in the language a learner would use. */
  description: string
}

export type Outcome =
  | { kind: 'unique'; values: [number, number]; steps: Step[] }
  | { kind: 'infinite'; free: number[]; steps: Step[] }
  | { kind: 'none'; reason: string; steps: Step[] }

/** Entries below this count as zero when choosing a pivot. */
const EPS = 1e-9

const round = (value: number): number => (Math.abs(value) < 1e-12 ? 0 : value)

function clone(rows: Row[]): Row[] {
  return rows.map((row) => [row[0], row[1], row[2]] as Row)
}

function label(index: number): string {
  return `R${index + 1}`
}

/**
 * Reduced row echelon form, with every elementary operation recorded.
 *
 * Partial pivoting picks the largest available entry in the column, which is
 * also worth saying out loud: it is the difference between a stable computation
 * and one that loses digits when a pivot is tiny.
 */
export function eliminate(input: Row[]): Outcome {
  let rows = clone(input)
  const steps: Step[] = []
  // Which column each row's pivot sits in. A column without a pivot is free, and
  // counting pivots is the only reliable test: after elimination a multiple row
  // reads 0 = 0 and its column still holds the other row's pivot.
  const pivotColumn: number[] = [-1, -1]

  for (let col = 0; col < 2; col++) {
    let pivot = -1
    let best = EPS
    for (let r = col; r < 2; r++) {
      const magnitude = Math.abs(rows[r][col])
      if (magnitude > best) {
        best = magnitude
        pivot = r
      }
    }

    if (pivot < 0) continue

    if (pivot !== col) {
      const swap = rows[col]
      rows[col] = rows[pivot]
      rows[pivot] = swap
      steps.push({ matrix: clone(rows), description: `Swap ${label(col)} and ${label(pivot)} to put a nonzero entry in the pivot position.` })
    }

    pivotColumn[col] = col
    const p = rows[col][col]
    if (Math.abs(p - 1) > EPS) {
      rows[col] = [round(rows[col][0] / p), round(rows[col][1] / p), round(rows[col][2] / p)]
      steps.push({ matrix: clone(rows), description: `Divide ${label(col)} by ${round(p)} so the pivot becomes 1.` })
    }

    for (let r = 0; r < 2; r++) {
      if (r === col) continue
      const factor = rows[r][col]
      if (Math.abs(factor) > EPS) {
        rows[r] = [
          round(rows[r][0] - factor * rows[col][0]),
          round(rows[r][1] - factor * rows[col][1]),
          round(rows[r][2] - factor * rows[col][2]),
        ]
        steps.push({
          matrix: clone(rows),
          description: `Replace ${label(r)} by ${label(r)} − ${round(factor)}${label(col)}. This subtracts a multiple of one equation from another, so every solution of the old pair still satisfies the new one.`,
        })
      }
    }
  }

  // A row of zeros with a nonzero right-hand side is a contradiction.
  for (const row of rows) {
    if (Math.abs(row[0]) < EPS && Math.abs(row[1]) < EPS) {
      if (Math.abs(row[2]) > EPS) {
        return {
          kind: 'none',
          reason: `Row ${label(rows.indexOf(row))} says 0 = ${round(row[2])}, which is false, so no pair of values satisfies both equations.`,
          steps,
        }
      }
    }
  }

  const pivots = pivotColumn.filter((c) => c >= 0)
  const free: number[] = []
  for (let col = 0; col < 2; col++) if (!pivots.includes(col)) free.push(col)
  if (free.length > 0) {
    return { kind: 'infinite', free, steps }
  }

  const values: [number, number] = [0, 0]
  pivotColumn.forEach((col, row) => {
    if (col >= 0) values[col] = round(rows[row][2])
  })
  return { kind: 'unique', values, steps }
}

/** Text for the outcome, so the template and its tests cannot disagree. */
export function outcomeText(outcome: Outcome): string {
  switch (outcome.kind) {
    case 'unique':
      return `x = ${outcome.values[0]}, y = ${outcome.values[1]}`
    case 'infinite':
      return `infinitely many solutions: ${outcome.free.map((c) => (c === 0 ? 'x' : 'y')).join(' and ')} is free`
    case 'none':
      return 'no solution'
  }
}
