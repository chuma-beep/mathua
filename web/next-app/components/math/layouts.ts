// Mathua's virtual keyboards.
//
// These are Mathua's layouts, not MathLive's: five small grids chosen by
// `keyboards.ts` from the concept's domain. A learner on `arith.add.single` sees
// digits, four operations, a fraction, a root and an exponent — no `∫ ∂ Σ lim`,
// because MathLive can do calculus and this exercise cannot.
//
// Two rules the keycaps obey:
//
//  1. Structures are inserted as structures. A fraction key inserts
//     `\frac{#@}{#0}` (with a real placeholder to navigate into), not the
//     character `/`. A root key inserts `\sqrt{#0}`, not `√`. Typing `x²`
//     through MathLive must leave a power with a navigable exponent slot, and it
//     does — but only because the keycap is a structure, not a glyph.
//  2. The visible label is separate from the inserted LaTeX wherever the
//     LaTeX would read badly. The fraction key shows `a⁄b` and inserts a
//     fraction; the exponent key shows `xⁿ` and inserts `#@^{#?}`.
//
// Placeholder tokens are MathLive's: `#@` is the selection or the implicit
// argument to the left of the caret, `#0` the selection or a placeholder, `#?`
// always a placeholder.

import type { MathKeyboardMode } from './keyboards'

// A type-only import: it costs nothing at runtime and keeps the keycaps honest.
// MathLive's `width` is a literal union (`2 | 1 | 0.5 | 1.5 | 5`), which a
// hand-rolled `{ width?: number }` would happily accept and then render wrong, so
// the library's own type is the one worth having.
import type { VirtualKeyboardKeycap } from 'mathlive'

export type Keycap = string | Partial<VirtualKeyboardKeycap>

export type Layout = {
  label: string
  tooltip: string
  layers: { rows: Keycap[][] }[]
}

const BACKSPACE: Keycap = { latex: '[backspace]', width: 2 }
const LEFT: Keycap = '[left]'
const RIGHT: Keycap = '[right]'
const SEPARATOR: Keycap = { label: '[separator]', width: 0.5 }
const HIDE: Keycap = '[hide-keyboard]'

// Structural keycaps. Shared so a fraction looks and behaves identically in
// every mode.
const FRACTION: Keycap = { label: 'a⁄b', latex: '\\frac{#@}{#0}', aside: 'fraction', variants: [{ label: 'n⁄d', latex: '\\frac{#?}{#?}' }] }
const ROOT: Keycap = { label: '√', latex: '\\sqrt{#0}', aside: 'square root', variants: [{ label: '∛', latex: '\\sqrt[#?]{#0}' }] }
const POWER: Keycap = { label: 'xⁿ', latex: '#@^{#?}', aside: 'exponent' }
const TIMES: Keycap = '\\times'
const DIVIDE: Keycap = '\\div'
const MINUS: Keycap = '-'
const PLUS: Keycap = '+'
const EQUALS: Keycap = '='
const OPEN: Keycap = '('
const CLOSE: Keycap = ')'
const DOT: Keycap = '[.]'
const ZERO_WIDE: Keycap = { latex: '[0]', width: 2 }
const SCI: Keycap = { label: '×10ⁿ', latex: '#@\\times 10^{#?}', aside: 'scientific notation' }

const digit = (n: 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9): Keycap => `[${n}]`

/** Digits 7-8-9 downward, the layout a numeric keypad has in every culture. */
function digitColumn(rows: Keycap[][], right: Keycap[][], bottomRight: Keycap[]) {
  rows.push([digit(7), digit(8), digit(9), ...right[0]])
  rows.push([digit(4), digit(5), digit(6), ...right[1]])
  rows.push([digit(1), digit(2), digit(3), ...right[2]])
  rows.push([ZERO_WIDE, DOT, OPEN, CLOSE, EQUALS, ...bottomRight])
}

// A mode with no letters (elementary arithmetic) still needs a second layer for
// the operations, so the learner is never trapped on one grid.
const OPERATORS: Keycap[] = [PLUS, MINUS, TIMES, DIVIDE]

// alphabetic + symbols layer, appended to every mode that has letters.
// Two rows, not one: MathLive's guidance is at most ten keycaps per row, and a
// fourteen-keycap row stops fitting a 320px viewport.
const LETTER_ROW: Keycap[] = [
  { latex: 'x', variants: ['y', 'z', 'a', 'b', 'c', 'n'] },
  { latex: 'y', variants: ['z', 'a', 'b', 'c', 'n', 'x'] },
  { latex: 'n', variants: ['k', 'm', 'r', 's', 't'] },
  { latex: '\\pi', variants: ['\\theta', '\\alpha', '\\beta', '\\infty'] },
  ...OPERATORS,
]

const STRUCTURE_ROW: Keycap[] = [FRACTION, POWER, ROOT, LEFT, RIGHT, BACKSPACE]

const LETTERS: Keycap[][] = [LETTER_ROW, STRUCTURE_ROW]

// Elementary arithmetic gets its own grid rather than the shared `digitColumn`.
//
// The reason is measurable: MathLive's panel is about 330px wide on a 412px phone,
// so a row of eight keycaps can only be 37px each. Six or fewer per row clears the
// 44px touch-target floor. That costs one extra row and buys keys a thumb can hit
// without a second attempt — worth it on the screen where mistakes cost the most.
const ARITHMETIC: Layout = {
  label: 'Arithmetic',
  tooltip: 'Numbers, the four operations, fractions and powers',
  layers: [
    {
      rows: [
        [digit(7), digit(8), digit(9), DIVIDE, ROOT],
        [digit(4), digit(5), digit(6), TIMES, POWER],
        [digit(1), digit(2), digit(3), MINUS, PLUS],
        [ZERO_WIDE, DOT, OPEN, CLOSE, EQUALS],
        [FRACTION, LEFT, RIGHT, { ...BACKSPACE, width: 2 }],
      ],
    },
    {
      rows: [[...OPERATORS, SEPARATOR, FRACTION, POWER, ROOT, SCI, HIDE]],
    },
  ],
}

const FRACTIONS: Layout = {
  label: 'Fractions',
  tooltip: 'Numerator and denominator, mixed numbers',
  layers: [
    {
      rows: (() => {
        const rows: Keycap[][] = []
        digitColumn(
          rows,
          [[FRACTION, BACKSPACE], [FRACTION, LEFT, RIGHT], [MINUS, PLUS, SEPARATOR]],
          [{ label: 'n⁄d', latex: '\\frac{#?}{#?}', aside: 'empty fraction' }, HIDE],
        )
        return rows
      })(),
    },
    { rows: [[...OPERATORS, SEPARATOR, ROOT, POWER, HIDE]] },
  ],
}

const ALGEBRA: Layout = {
  label: 'Algebra',
  tooltip: 'Letters, numbers, the four operations, structures',
  layers: [
    {
      rows: (() => {
        const rows: Keycap[][] = []
        digitColumn(
          rows,
          [[FRACTION, ROOT, BACKSPACE], [TIMES, POWER, LEFT, RIGHT], [DIVIDE, PLUS, MINUS, SEPARATOR]],
          [{ latex: 'x', variants: ['y', 'z'] }, { latex: 'y' }, HIDE],
        )
        return rows
      })(),
    },
    { rows: LETTERS },
  ],
}

const GEOMETRY: Layout = {
  label: 'Geometry',
  tooltip: 'Measures, angles, roots and powers',
  layers: [
    {
      rows: (() => {
        const rows: Keycap[][] = []
        digitColumn(
          rows,
          [[ROOT, BACKSPACE], [POWER, LEFT, RIGHT], [TIMES, PLUS, MINUS, SEPARATOR]],
          [{ latex: '\\pi', variants: ['2\\pi'] }, { latex: '^\\circ', variants: ['\\pi/180'] }, HIDE],
        )
        return rows
      })(),
    },
    { rows: [[...OPERATORS, '\\pi', SEPARATOR, FRACTION, POWER, ROOT, HIDE]] },
  ],
}

const CALCULUS: Layout = {
  label: 'Calculus',
  tooltip: 'Limits, derivatives, integrals, the four operations',
  layers: [
    {
      rows: (() => {
        const rows: Keycap[][] = []
        digitColumn(
          rows,
          [[FRACTION, ROOT, BACKSPACE], [TIMES, POWER, LEFT, RIGHT], [DIVIDE, PLUS, MINUS, SEPARATOR]],
          [{ latex: '\\int', variants: ['\\int_#@^{#?}', '\\oint'] }, HIDE],
        )
        return rows
      })(),
    },
    {
      rows: [
        [
          { latex: '\\lim', variants: ['\\lim_{#@\\to #?}', '\\lim_{#@\\to \\infty}'], aside: 'limit' },
          { latex: '\\frac{d#@}{d#?}', aside: 'derivative' },
          { latex: '\\int', variants: ['\\int_#@^{#?}', '\\oint'] },
          { latex: '\\sin', variants: ['\\cos', '\\tan', '\\arcsin'] },
          '\\ln',
          '\\log',
          { label: '|x|', latex: '|#@|', aside: 'absolute value', variants: [{ label: 'ln|x|', latex: '\\ln|#@|' }] },
        ],
        STRUCTURE_ROW,
      ],
    },
  ],
}

const LAYOUTS: Record<MathKeyboardMode, Layout> = {
  arithmetic: ARITHMETIC,
  fractions: FRACTIONS,
  algebra: ALGEBRA,
  geometry: GEOMETRY,
  calculus: CALCULUS,
}

export function layoutForMode(mode: MathKeyboardMode): Layout {
  return LAYOUTS[mode] ?? LAYOUTS.algebra
}