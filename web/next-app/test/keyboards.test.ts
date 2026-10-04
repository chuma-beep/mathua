// The contextual keyboards.
//
// Two things are worth protecting here. The first is the mapping: a concept's
// grading_type and domain decide which control it gets and which keycaps it shows,
// and getting that wrong is invisible until a learner is staring at an integral key
// on a single-digit addition problem. The second is that structural keys insert
// structures — a fraction key that inserted the character `/` would look identical in
// a screenshot and behave completely differently.

import { describe, it, expect } from 'vitest'
import {
  presentationFor,
  presentationForQuestion,
  keyboardModeForDomain,
  editorForGradingType,
  DEFAULT_KEYBOARD_MODE,
  type MathKeyboardMode,
} from '../components/math/keyboards'
import { layoutForMode } from '../components/math/layouts'

const MODES: MathKeyboardMode[] = ['arithmetic', 'algebra', 'fractions', 'geometry', 'calculus']

function keycaps(mode: MathKeyboardMode): string[] {
  return layoutForMode(mode)
    .layers.flatMap(l => l.rows.flat())
    .map(k => (typeof k === 'string' ? k : (k.latex ?? k.label ?? '')))
}

// LaTeX inserted by a keycap, wherever it appears including variants.
function insertedLatex(mode: MathKeyboardMode): string[] {
  const out: string[] = []
  for (const layer of layoutForMode(mode).layers) {
    for (const row of layer.rows) {
      for (const key of row) {
        const collect = (k: unknown) => {
          if (typeof k === 'string') out.push(k)
          else if (k && typeof k === 'object') {
            const cap = k as { latex?: string; variants?: unknown[] }
            if (cap.latex) out.push(cap.latex)
            if (cap.variants) cap.variants.forEach(collect)
          }
        }
        collect(key)
      }
    }
  }
  return out
}

describe('editor vs text', () => {
  // Every grading type gets the math editor. This was the five expression-valued types
  // only, on the reasoning that MathLive earns its place where the grader parses an
  // expression — which left 223 of 657 concepts (209 multiple_choice, 11 tuple,
  // 2 comparison, 1 ordering) typing mathematics into a plain input with no fraction
  // key, no exponent and no parentheses.
  //
  // The grading risk that motivated the old split is real and is now covered by
  // measurement rather than by avoidance: `internal/grader/mathlive_compat_test.go`
  // grades a recording of what MathLive actually emits, including the multiple-choice
  // shapes that broke first (`yes` → `y e s`). See `fixImplicitLetterSpacing`.
  it.each([
    ['numeric', 'math'],
    ['symbolic', 'math'],
    ['expression', 'math'],
    ['polynomial', 'math'],
    ['complex', 'math'],
    ['multiple_choice', 'math'],
    ['tuple', 'math'],
    ['ordering', 'math'],
    ['comparison', 'math'],
  ])('grading_type %s → %s editor', (gradingType, editor) => {
    expect(editorForGradingType(gradingType)).toBe(editor)
  })

  it('still degrades to the text editor for a concept the corpus does not know', () => {
    // A fixture id, or a concept added server-side before the corpus rebuild. The plain
    // input is the safe default because it grades correctly everywhere, and it is the
    // only thing that degrades now.
    expect(presentationForQuestion(undefined, 'not-a-real-concept').editor).toBe('text')
  })

  it('prefers the server grading_type over the corpus, either way it points', () => {
    // ADR-009: every served question carries grading_type, and the server is the
    // authority for it. With the editor universal this cannot change the *outcome* for a
    // known type, but it still must not be the corpus that decides.
    expect(presentationForQuestion('multiple_choice', 'arith.add.single').editor).toBe('math')
    expect(presentationForQuestion('numeric', 'not-a-real-concept').editor).toBe('math')
  })
})

describe('domain → keyboard mode', () => {
  it.each([
    ['arithmetic', 'arithmetic'],
    ['prealgebra', 'arithmetic'],
    ['number_theory', 'arithmetic'],
    ['fractions', 'fractions'],
    ['algebra', 'algebra'],
    ['abstract_algebra', 'algebra'],
    ['linear_algebra', 'algebra'],
    ['complex_numbers', 'algebra'],
    ['trigonometry', 'algebra'],
    ['geometry', 'geometry'],
    ['calculus', 'calculus'],
    ['differential_equations', 'calculus'],
  ])('%s → %s', (domain, mode) => {
    expect(keyboardModeForDomain(domain)).toBe(mode)
  })

  it('falls back to algebra for an unmapped or absent domain', () => {
    // Algebra is the general case: it has the digits, the four operations, a
    // fraction and an exponent. Falling back to arithmetic would hide the letters a
    // linear-algebra or number-theory question needs, which is the worse failure.
    expect(keyboardModeForDomain('topology')).toBe(DEFAULT_KEYBOARD_MODE)
    expect(keyboardModeForDomain(undefined)).toBe(DEFAULT_KEYBOARD_MODE)
  })

  it('resolves a real concept end to end', () => {
    expect(presentationForQuestion(undefined, 'arith.add.single')).toEqual({ editor: 'math', keyboard: 'arithmetic' })
    expect(presentationForQuestion(undefined, 'frac.add.diff')).toEqual({ editor: 'math', keyboard: 'fractions' })
    expect(presentationForQuestion(undefined, 'calc.limit.concept')).toEqual({ editor: 'math', keyboard: 'calculus' })
  })
})

describe('layouts', () => {
  it('no keycap inserts a bracketed command token as literal text', () => {
    // The bug class that shipped broken keys in all five layouts. MathLive resolves
    // a keycap as `command` → `insert` → `key` → `latex` → `typedText(label)`, and it
    // merges its own KEYCAP_SHORTCUTS definition only for a bare-string keycap or one
    // carrying a `label`/`key`. So `{ latex: '[backspace]' }` matched nothing and fell
    // through to the insert branch: the key typed the literal string `[backspace]`,
    // and `{ latex: '[0]' }` typed `[0]`. Both look correct in a screenshot.
    //
    // A bracketed token in `latex` is therefore always a mistake. Name the shortcut
    // in `label` instead, which is what MathLive's own layouts do.
    const BRACKETED = /^\s*\[/
    for (const mode of MODES) {
      for (const layer of layoutForMode(mode).layers) {
        for (const row of layer.rows) {
          for (const key of row) {
            if (typeof key === 'string') continue // resolved via KEYCAP_SHORTCUTS
            expect(key.command, `${mode}: use label, not command, for shortcuts`).toBeUndefined()
            if (typeof key.latex !== 'string') continue
            expect(
              BRACKETED.test(key.latex),
              `${mode}: keycap latex ${JSON.stringify(key.latex)} would be typed as literal text`,
            ).toBe(false)
            for (const v of [key.shift, ...(key.variants ?? [])]) {
              const l = typeof v === 'string' ? v : v?.latex
              if (typeof l !== 'string') continue
              expect(BRACKETED.test(l), `${mode}: variant ${JSON.stringify(l)}`).toBe(false)
            }
          }
        }
      }
    }
  })

  it('every mode has a label, a tooltip and at least one layer', () => {
    for (const mode of MODES) {
      const l = layoutForMode(mode)
      expect(l.label).toBeTruthy()
      expect(l.tooltip).toBeTruthy()
      expect(l.layers.length).toBeGreaterThan(0)
    }
  })

  it('every row has keycaps and none is wider than the recommended ten', () => {
    // MathLive's own guidance: more than ten keycaps in a row and the grid stops
    // fitting a 320px viewport.
    for (const mode of MODES) {
      for (const layer of layoutForMode(mode).layers) {
        expect(layer.rows.length).toBeGreaterThan(0)
        for (const row of layer.rows) {
          expect(row.length).toBeGreaterThan(0)
          expect(row.length, `${mode} row of ${row.length}`).toBeLessThanOrEqual(10)
        }
      }
    }
  })

  it('gives every mode the digits, four operations, a fraction, a root and an exponent', () => {
    // A learner must never be unable to enter the answer the hint asks for.
    for (const mode of MODES) {
      const keys = insertedLatex(mode)
      const has = (needle: string) => keys.some(k => k.includes(needle))
      expect(has('\\frac'), `${mode} has no fraction`).toBe(true)
      expect(has('\\sqrt'), `${mode} has no root`).toBe(true)
      expect(has('^{'), `${mode} has no exponent`).toBe(true)
      expect(has('\\times'), `${mode} has no multiplication`).toBe(true)
      expect(has('\\div'), `${mode} has no division`).toBe(true)
      expect(keys.some(k => k === '-' || k === '+')).toBe(true)
      // Digits 1-9 are bare strings; zero is a shortcut we name in `label`, because
      // `[0]`'s payload (`latex: "0"`) lives in MathLive's own table. Presence is
      // therefore a question about the declared name, not about what we insert --
      // `insertedLatex` is the right tool for the LaTeX assertions above and the
      // wrong one here.
      const names = keycaps(mode)
      for (const d of ['[0]', '[1]', '[2]', '[3]', '[4]', '[5]', '[6]', '[7]', '[8]', '[9]']) {
        expect(names, `${mode} is missing ${d}`).toContain(d)
      }
    }
  })

  it('inserts structures, not glyphs', () => {
    // The whole point of the fraction key. A `/` character would render the same
    // and behave nothing like a fraction.
    const keys = insertedLatex('algebra')
    const fraction = keys.find(k => k.includes('\\frac'))
    expect(fraction).toBeDefined()
    expect(fraction).toContain('#') // navigable placeholder
    expect(keys.some(k => k.includes('\\sqrt{'))).toBe(true)
    expect(keys.some(k => /#@\^\{#\?\}/.test(k))).toBe(true)
  })

  it('withholds calculus keys from elementary arithmetic', () => {
    // The requirement in one assertion: an arithmetic exercise must not offer ∫ ∂ Σ
    // lim sin cos tan just because MathLive can render them.
    const arithmetic = insertedLatex('arithmetic').join(' ')
    for (const forbidden of ['\\int', '\\lim', '\\partial', '\\sum', '\\sin', '\\cos', '\\tan', '\\ln']) {
      expect(arithmetic, `arithmetic keyboard exposes ${forbidden}`).not.toContain(forbidden)
    }
  })

  it('offers calculus keys on the calculus keyboard', () => {
    const calculus = insertedLatex('calculus').join(' ')
    expect(calculus).toContain('\\int')
    expect(calculus).toContain('\\lim')
  })

  it('offers pi and degrees on the geometry keyboard but not on arithmetic', () => {
    expect(insertedLatex('geometry').join(' ')).toContain('\\pi')
    expect(insertedLatex('geometry').join(' ')).toContain('^\\circ')
    expect(insertedLatex('arithmetic').join(' ')).not.toContain('^\\circ')
  })

  it('offers bare letter keys only where a bare variable is the answer', () => {
    // A letter *key* is a keycap whose insert is a single variable. Matching /[a-z]/
    // against the whole set would flag \\frac and \\sqrt, which every keyboard has.
    const isLetterKey = (k: string) => /^[a-z]$/.test(k)

    // Algebra answers can be a bare variable, so it offers x, y, n.
    expect(insertedLatex('algebra').filter(isLetterKey)).toEqual(expect.arrayContaining(['x', 'y', 'n']))

    // Elementary arithmetic and fractions have no variables to enter. A stray `x` on
    // an addition problem invites an answer the question cannot have.
    expect(insertedLatex('arithmetic').filter(isLetterKey)).toEqual([])
    expect(insertedLatex('fractions').filter(isLetterKey)).toEqual([])

    // Calculus deliberately has no bare letter key either: its variable goes into a
    // placeholder (\\frac{d#@}{d#?} supplies the d and the two slots), and the learner
    // types the variable on the physical keyboard. Verified in e2e/math-input.spec.ts.
    expect(insertedLatex('calculus').filter(isLetterKey)).toEqual([])
    expect(insertedLatex('calculus').some(k => k.includes('#@'))).toBe(true)
  })

  it('gives every mode a way to delete and a way to dismiss the keyboard', () => {
    for (const mode of MODES) {
      // Neither of these is an insertion, so neither belongs in `insertedLatex`:
      // backspace runs `deleteBackward`, and the dismiss key hides the panel. Both
      // are asserted by the *name* we declare, which is what has to match
      // MathLive's KEYCAP_SHORTCUTS for the library to supply the behaviour. The
      // e2e suite then proves the key actually deletes.
      const names = keycaps(mode)
      expect(names, `${mode} has no backspace`).toContain('[backspace]')
      expect(names, `${mode} cannot dismiss the keyboard`).toContain('[hide-keyboard]')
    }
  })

  it('keeps cursor movement available', () => {
    for (const mode of MODES) {
      const keys = insertedLatex(mode)
      expect(keys, `${mode} has no arrow keys`).toContain('[left]')
      expect(keys, `${mode} has no arrow keys`).toContain('[right]')
    }
  })
})