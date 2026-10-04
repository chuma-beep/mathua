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

/**
 * One page of keycaps.
 *
 * A page is a MathLive *layout*, not a MathLive *layer*, and that distinction is
 * load-bearing. MathLive's keypad toolbar builds its switcher by iterating `kb.layouts`
 * and emitting `data-layer="${layout.layers[0].id}"` per entry — so it navigates
 * between layouts and there is no navigation between the `layers` of a single layout.
 * Verified rather than assumed: with a three-layer layout installed, the toolbar offered
 * exactly one switcher entry (MathLive's `alphabetic`), and pressing a keycap that
 * carried a `layer` property did nothing at all, because `renderKeycap` only adds a
 * `layer-switch` class to it.
 *
 * So each page is exported as its own single-layer `Layout` and `layoutsForMode` returns
 * them together. `label` is what the toolbar shows for each.
 */
export type MathuaLayer = { id: string; label: string; rows: Keycap[][] }

export type Layout = {
  label: string
  tooltip: string
  layers: MathuaLayer[]
}

/**
 * Layer ids.
 *
 * Layers used to be unreachable: `layers: [{rows}, {rows}]` renders both, but nothing
 * pointed at the second, so the operators and scientific-notation layer had no way in --
 * the same class of bug as the keycap that typed its own name.
 *
 * The fix is not a keycap. A keycap with a `layer` property only gets a `layer-switch`
 * *class* from `renderKeycap` and does nothing when pressed; the handler reads a
 * `data-layer` attribute from an ancestor, and the only thing that emits one is the
 * toolbar switcher MathLive builds itself from the layers' labels. So the layers need
 * ids and labels, and navigation is the toolbar. Verified in the e2e suite by pressing
 * the switcher and asserting the keypad changed.
 */
export const LAYER_MAIN = 'mathua-main'
export const LAYER_SYMBOLS = 'mathua-symbols'
export const LAYER_LETTERS = 'mathua-letters'

const MAIN: MathuaLayer = { id: LAYER_MAIN, label: '123', rows: [] }
const SYMBOLS: MathuaLayer = { id: LAYER_SYMBOLS, label: '±÷', rows: [] }
const LAYER_LETTERS_TEMPLATE: MathuaLayer = { id: LAYER_LETTERS, label: 'abc', rows: [] }

// `label`, not `latex`, and that distinction is the whole reason these two work.
//
// MathLive resolves a keycap in this order (`executeKeycapCommand`): `command`,
// then `insert`, then `key`, then **`latex` — which it runs as an insert**, and
// only then `typedText(label)`. Separately, `normalizeKeycap` merges the
// library's own definition for a keycap in only two cases: the keycap is a bare
// string, or it carries a `label` or `key` that names a known shortcut.
//
// These two used `latex`, which is neither. So they missed the shortcut table
// *and* fell through to the insert branch, which meant backspace typed the
// literal text `[backspace]` into the field and the zero key typed `[0]` — in all
// five layouts. Both are fixed by naming the shortcut in `label`, which is also
// what MathLive's own layouts do (`{ label: "[backspace]", width: 1 }`).
//
// What the merge brings in, from KEYCAP_SHORTCUTS: `command:
// "performWithFeedback(deleteBackward)"` for backspace, `latex: "0"` for zero,
// an SVG delete glyph as the label, and `class: "action bottom right hide-shift"`
// plus `width: 1.5` from the library. We keep our own `width`, because the
// 320px layout needs the extra room.
const BACKSPACE: Keycap = { label: '[backspace]', width: 2 }
const LEFT: Keycap = '[left]'
const RIGHT: Keycap = '[right]'
const SEPARATOR: Keycap = { label: '[separator]', width: 0.5 }
const HIDE: Keycap = '[hide-keyboard]'

// A spacebar that works, for the mode that needs one most.
//
// `mathModeSpace` (set in `MathLiveField`) is what repairs MathLive's *own* spacebar --
// the one in its `alphabetic` layout, which is appended to every mode and reachable from
// the keypad toolbar. Upstream that key is `{ label: " ", width: 1.5 }`: a label and a
// width, and no `command`, `insert`, `latex` or `key`, so `executeKeycapCommand` falls
// through to `typedText(" ")`, which reaches the `[Space]` handler and, with
// `mathModeSpace` empty, only calls `moveAfterParent()`. The caret moved and nothing was
// inserted. It is an upstream defect, not a Mathua regression.
//
// This is the same key on the arithmetic grid, which is where it earns its place: the
// remainder format `Q R` and mixed numbers are arithmetic answers, and reaching them
// through `alphabetic` first means two taps before a space. One unit wide, on the one
// row that has a spare unit to give -- see the width budget in `test/keyboards.test.ts`,
// which is why this is arithmetic-only and not on every grid.
//
// `latex: '\\,'` rather than a label-only key, and this is the ADR-029 trap twice over: a
// label of `' '` would work (it maps to `[Space]`) but renders as a blank face, which is
// half of why the key read as broken; a visible label of `'␣'` renders the glyph but, with
// no payload beside it, *types* the glyph `␣` into the field. So the payload is stated.
// `\\,` is the proven round-trip value -- the recorded corpus holds
// `expression 5 R 3 <- 5\ R\ 3`, which comes back as the plain text `5 R 3` and grades.
const SPACE: Keycap = { label: '␣', aside: 'space', latex: '\\,' }

// Structural keycaps. Shared so a fraction looks and behaves identically in
// every mode.
const FRACTION: Keycap = { label: 'a⁄b', latex: '\\frac{#@}{#0}', aside: 'fraction', variants: [{ label: 'n⁄d', latex: '\\frac{#?}{#?}' }] }
const ROOT: Keycap = { label: '√', latex: '\\sqrt{#0}', aside: 'square root', variants: [{ label: '∛', latex: '\\sqrt[#?]{#0}' }] }
const POWER: Keycap = { label: 'xⁿ', latex: '#@^{#?}', aside: 'exponent' }
// The separator. 14 concepts cannot be answered without it -- every `tuple` and
// `ordering` answer is a list -- and 209 more are typed alongside their tap targets.
// Absent until the editor went universal, at which point it stopped being optional.
const COMMA: Keycap = ','
// `frac.ops.compare` and `dec.basics.compare` grade a comparison, and typing `<`
// or `≥` on a phone keyboard means reaching for a symbol row that may not be there.
const LESS: Keycap = { label: '<', latex: '<' }
const GREATER: Keycap = { label: '>', latex: '>' }
const LEQ: Keycap = { label: '≤', latex: '\\le' }
const GEQ: Keycap = { label: '≥', latex: '\\ge' }
const TIMES: Keycap = '\\times'
const DIVIDE: Keycap = '\\div'
const MINUS: Keycap = '-'
const PLUS: Keycap = '+'
const EQUALS: Keycap = '='
const OPEN: Keycap = '('
const CLOSE: Keycap = ')'
const DOT: Keycap = '[.]'
// Same story as BACKSPACE: `[0]` is a shortcut whose payload is `latex: "0"`, so
// it has to be named in `label` or it inserts the two-character string `[0]`.
const ZERO_WIDE: Keycap = { label: '[0]', width: 2 }
const SCI: Keycap = { label: '×10ⁿ', latex: '#@\\times 10^{#?}', aside: 'scientific notation' }

const digit = (n: 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9): Keycap => `[${n}]`

/** Digits 7-8-9 downward, the layout a numeric keypad has in every culture. */
function digitColumn(rows: Keycap[][], right: Keycap[][], bottomRight: Keycap[]) {
  rows.push([digit(7), digit(8), digit(9), ...right[0]])
  rows.push([digit(4), digit(5), digit(6), ...right[1]])
  rows.push([digit(1), digit(2), digit(3), ...right[2]])
  // COMMA is here rather than only on a second layer: it is the separator every tuple and
  // ordering answer needs, and a key the learner has to discover is a key they will not
  // find mid-question.
  rows.push([ZERO_WIDE, DOT, COMMA, OPEN, CLOSE, ...bottomRight])
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

// Every layout ends here. A comma and the comparisons are on this row because a
// second layer the learner has to discover is a worse answer than one more key.
const STRUCTURE_ROW: Keycap[] = [FRACTION, POWER, ROOT, LEFT, RIGHT, BACKSPACE]

// Six keycaps, deliberately. A learner answering "which is larger" needs one of these
// and nothing else on the screen, and the row stays inside the ten-keycap guidance so
// it still fits a 320px viewport at the 44px touch-target floor.
const COMPARISON_ROW: Keycap[] = [LESS, GREATER, LEQ, GEQ, COMMA, HIDE]

const LETTER_LAYERS: MathuaLayer[] = [
  { ...MAIN, rows: [LETTER_ROW] },
  { ...SYMBOLS, rows: [STRUCTURE_ROW] },
  { ...LAYER_LETTERS_TEMPLATE, rows: [COMPARISON_ROW] },
]

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
    { ...MAIN, rows: [
        [digit(7), digit(8), digit(9), DIVIDE, ROOT],
        [digit(4), digit(5), digit(6), TIMES, POWER],
        [digit(1), digit(2), digit(3), MINUS, PLUS, SPACE],
        [ZERO_WIDE, DOT, OPEN, CLOSE, EQUALS],
        [FRACTION, COMMA, LEFT, RIGHT, { ...BACKSPACE, width: 2 }],
      ],
    },
    { ...SYMBOLS, rows: [[...OPERATORS, COMMA, FRACTION, POWER, ROOT]] },
    { ...LAYER_LETTERS_TEMPLATE, rows: [COMPARISON_ROW] },
  ],
}

const FRACTIONS: Layout = {
  label: 'Fractions',
  tooltip: 'Numerator and denominator, mixed numbers',
  layers: [
    {
      ...MAIN,
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
    { ...SYMBOLS, rows: [[...OPERATORS, COMMA, ROOT, POWER]] },
    { ...LAYER_LETTERS_TEMPLATE, rows: [COMPARISON_ROW] },
  ],
}

const ALGEBRA: Layout = {
  label: 'Algebra',
  tooltip: 'Letters, numbers, the four operations, structures',
  layers: [
    {
      ...MAIN,
      rows: (() => {
        const rows: Keycap[][] = []
        digitColumn(
          rows,
          [[FRACTION, ROOT, BACKSPACE], [TIMES, POWER, LEFT, RIGHT, { latex: 'y' }], [DIVIDE, PLUS, MINUS, SEPARATOR]],
          [{ latex: 'x', variants: ['y', 'z'] }, HIDE],
        )
        return rows
      })(),
    },
    ...LETTER_LAYERS,
  ],
}

const GEOMETRY: Layout = {
  label: 'Geometry',
  tooltip: 'Measures, angles, roots and powers',
  layers: [
    {
      ...MAIN,
      rows: (() => {
        const rows: Keycap[][] = []
        digitColumn(
          rows,
          [[ROOT, BACKSPACE], [POWER, LEFT, RIGHT, { latex: '^\\circ', variants: ['\\pi/180'] }], [TIMES, PLUS, MINUS, SEPARATOR]],
          [{ latex: '\\pi', variants: ['2\\pi'] }, HIDE],
        )
        return rows
      })(),
    },
    { ...SYMBOLS, rows: [[...OPERATORS, '\\pi', COMMA, FRACTION, POWER]] },
    { ...LAYER_LETTERS_TEMPLATE, rows: [COMPARISON_ROW] },
  ],
}

const CALCULUS: Layout = {
  label: 'Calculus',
  tooltip: 'Limits, derivatives, integrals, the four operations',
  layers: [
    {
      ...MAIN,
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
    { ...MAIN, rows: [
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

/**
 * Every page of a mode, in toolbar order.
 *
 * Handed to MathLive as `kb.layouts`, so the toolbar renders one switcher entry per page
 * and the learner can actually reach the symbols and comparisons. Returns a copy rather
 * than the originals because MathLive freezes the array it is given.
 */
export function layoutsForMode(mode: MathKeyboardMode): Layout[] {
  return layoutForMode(mode).layers.map(layer => ({
    label: layer.label,
    tooltip: layoutForMode(mode).tooltip,
    layers: [{ ...layer }],
  }))
}