// Finite groups shown by their operation table. The table is the group: no
// geometric analogy is implied unless the group really is a symmetry group.

export type GroupId = 'z4' | 'z5' | 's3'

export interface Group {
  id: GroupId
  label: string
  description: string
  elements: string[]
  table: number[][]
  abelian: boolean
  cyclic: boolean
}

const Z = (n: number): Group => {
  const elements = Array.from({ length: n }, (_, i) => String(i))
  const table = elements.map((_, i) => elements.map((_, j) => (i + j) % n))
  return {
    id: n === 4 ? 'z4' : 'z5',
    label: `Z${n}`,
    description: `Addition modulo ${n}`,
    elements,
    table,
    abelian: true,
    cyclic: true,
  }
}

const PERMS: number[][] = [
  [0, 1, 2],
  [1, 0, 2],
  [0, 2, 1],
  [2, 1, 0],
  [1, 2, 0],
  [2, 0, 1],
]
const PERM_LABELS = ['e', '(12)', '(23)', '(13)', '(123)', '(132)']

const compose = (p: number[], q: number[]) => p.map((value) => q[value])

function s3(): Group {
  const index = new Map(PERMS.map((p, i) => [p.join(','), i]))
  const table = PERMS.map((row) => PERMS.map((col) => index.get(compose(row, col).join(',')) ?? 0))
  return {
    id: 's3',
    label: 'S3',
    description: 'Permutations of three symbols, the symmetries of a triangle',
    elements: PERM_LABELS,
    table,
    abelian: false,
    cyclic: false,
  }
}

export function groupById(id: GroupId): Group {
  if (id === 'z5') return Z(5)
  if (id === 's3') return s3()
  return Z(4)
}

export function elementOrder(group: Group, index: number): number {
  let current = 0
  for (let k = 1; k <= group.elements.length; k++) {
    current = group.table[current][index]
    if (current === 0) return k
  }
  return 0
}

export function inverseOf(group: Group, index: number): number {
  for (let j = 0; j < group.elements.length; j++) {
    if (group.table[index][j] === 0) return j
  }
  return 0
}

export function generatedSubgroup(group: Group, index: number): number[] {
  const seen = new Set<number>()
  let current = 0
  for (let k = 0; k < group.elements.length; k++) {
    seen.add(current)
    if (current === 0 && k > 0) break
    current = group.table[current][index]
  }
  return [...seen].sort((a, b) => a - b)
}

export function isAbelian(group: Group): boolean {
  for (let i = 0; i < group.elements.length; i++) {
    for (let j = 0; j < group.elements.length; j++) {
      if (group.table[i][j] !== group.table[j][i]) return false
    }
  }
  return true
}
