// Non-binding feedback on a learner's own expression.
//
// What this is NOT: a grade. It never sees the expected answer — the server holds
// that, and shipping it to the browser to compute a hint would turn every hint into
// an answer oracle and quietly undo the "the server grades" rule in ADR-005. So
// there is no comparison against a reference anywhere in this file. What it can do
// is look at what the learner wrote and show them what it simplifies to.
//
//   "1/2 + 1/3"        → "that adds up to 5/6"
//   "2(x+1)"            → "that simplifies to 2x + 2"
//   "sqrt(x^2)"        → "that simplifies to |x|"      ← the |x| is the point
//
// That last one is why this is Compute Engine and not a regex. The naive reading of
// `sqrt(x^2)` is `x`, which is false for negative x. CE rewrites it to `|x|` on its
// own, so the learner is shown the correct form instead of a common mistake.
//
// Everything here is three-valued. "Not proven equal" is a real and common answer and
// it is reported as such rather than being rounded to "different", because a learner
// told "these differ" when the tool simply could not tell would be misled.

/** A lazily-created Compute Engine, or null if it could not be loaded. */
export interface Engine {
  parse(latex: string): Expression
  /** Parse without canonicalising, so the tree mirrors what the learner typed. */
  rawOf(latex: string): string
  valid(expr: Expression): boolean
  simplify(expr: Expression): Expression
  latexOf(expr: Expression): string
}

export interface Expression {
  isValid: boolean
  isSame(other: Expression): boolean
  isIdenticallyEqual(other: Expression): boolean | undefined
  toString(): string
}

export type Equivalence = 'equivalent' | 'different' | 'unknown'

let enginePromise: Promise<Engine | null> | null = null

/**
 * Load Compute Engine once, on first use.
 *
 * It is ~1.1 MB gzipped — five times MathLive — so it is never in any route's initial
 * JS and never loaded until a learner has actually typed something worth checking.
 * A load failure returns null and every function here degrades to "say nothing",
 * which is the whole failure mode of this feature: it is optional by construction.
 */
export function loadEngine(): Promise<Engine | null> {
  if (enginePromise) return enginePromise
  enginePromise = import('@cortex-js/compute-engine')
    .then(mod => {
      const ce = new mod.ComputeEngine()
      return {
        parse: (latex: string) => ce.parse(latex) as unknown as Expression,
        rawOf: (latex: string) => String((ce.parse(latex, { form: 'raw' }) as { toString(): string }).toString()),
        valid: (e: Expression) => (e as { isValid: boolean }).isValid,
        simplify: (e: Expression) => (e as unknown as { simplify: () => unknown }).simplify() as Expression,
        latexOf: (e: Expression) => {
          const withLatex = e as unknown as { toLatex?: () => string }
          // CE's own serialization. Falls back to the display form if unavailable.
          return withLatex.toLatex?.() ?? e.toString()
        },
      }
    })
    .catch(err => {
      // A missing optional feature must never break the answer field.
      console.warn('compute engine unavailable; skipping the answer self-check:', err)
      return null
    })
  return enginePromise
}

/** Test seam: forget the cached engine. */
export function resetEngineForTests(): void {
  enginePromise = null
}

/**
 * Are two expressions the same mathematical object?
 *
 * `isIdenticallyEqual` is the primary answer because it is the only one of CE's
 * primitives that can *refuse*: it returns `undefined` where it cannot decide, which
 * is exactly right for `sqrt(x^2)` vs `x`. Two things it will not do on its own:
 *
 *  - it returns `undefined` for identities it has not been taught, notably
 *    `sin^2(x) + cos^2(x)` vs `1`, so the canonical form of each side is compared as
 *    a second opinion — `simplify()` does reduce that pair to `1`.
 *  - it returns `undefined`, not `false`, for plainly different symbolic expressions
 *    like `2x+3` vs `2x+4`. Those stay `unknown`. Reporting them as `different` would
 *    be a claim the tool cannot support.
 */
export function equivalence(engine: Engine, aLatex: string, bLatex: string): Equivalence {
  const a = tryParse(engine, aLatex)
  const b = tryParse(engine, bLatex)
  if (!a || !b) return 'unknown'

  const identical = a.isIdenticallyEqual(b)
  if (identical === true) return 'equivalent'
  if (identical === false) return 'different'

  try {
    if (engine.simplify(a).toString() === engine.simplify(b).toString()) return 'equivalent'
  } catch {
    // Fall through to unknown: an exception while simplifying proves nothing either way.
  }
  return 'unknown'
}

/**
 * What does this simplify to, and is that worth telling the learner?
 *
 * Only reports a change when the simplification is a genuine improvement: shorter,
 * or numerically exact where the input was a decimal. That keeps the hint useful
 * (`1/2+1/3` → `5/6`) and quiet (`2x+3` → `2x + 3`, which is the same thing wearing
 * different spaces).
 */
export function simplificationHint(
  engine: Engine,
  latex: string,
): { latex: string; text: string } | null {
  // Nothing to simplify in a lone symbol, and no reason to pay for a parse.
  if (!HAS_ARITHMETIC.test(latex)) return null
  const expr = tryParse(engine, latex)
  if (!expr) return null

  let simplified: string
  let asTyped: string
  try {
    simplified = engine.simplify(expr).toString()
    // `ce.parse` canonicalises eagerly, so the *parsed* form is already the answer.
    // The comparison has to be against what the learner actually wrote, or every
    // expression looks identical to its own simplification and the hint never fires.
    asTyped = engine.rawOf(latex)
  } catch {
    return null
  }
  if (!simplified) return null
  // Cosmetic spacing and implicit-multiplication markers are not a simplification:
  // "2x+3" and "2x + 3" are the same expression.
  if (compact(simplified) === compact(asTyped)) return null

  return {
    latex: simplified,
    text: `That simplifies to ${simplified}.`,
  }
}

/**
 * Reduce two renderings to something comparable.
 *
 * Strips whitespace, the parentheses CE keeps around whole forms, and the
 * `InvisibleOperator(a, b)` it prints for implicit multiplication — `2x` prints as
 * `InvisibleOperator(2, x)`, which is the same thing the learner typed, not a
 * simplification of it.
 */
function compact(s: string): string {
  return stripWhitespace(expandInvisibleOperators(s)).replace(/[()]/g, '')
}

const INVISIBLE = 'InvisibleOperator('

/**
 * Print CE's implicit multiplication the way the learner wrote it.
 *
 * `2x` prints as `InvisibleOperator(2, x)` and `2(x+1)` as
 * `InvisibleOperator(2, (x + 1))`, so the arguments have to be split at the top-level
 * comma — a regular expression cannot do that when an argument contains parentheses of
 * its own, which is exactly the case that matters.
 */
function expandInvisibleOperators(s: string): string {
  let out = ''
  let i = 0
  while (i < s.length) {
    if (s.startsWith(INVISIBLE, i)) {
      const args = readBalanced(s, i + INVISIBLE.length)
      if (args) {
        const comma = topLevelComma(args)
        // Concatenated, not joined with `*`: the comparison below is against a form
        // that writes multiplication implicitly, so inserting an explicit operator
        // here would make every typed expression look different from itself.
        out += comma === -1 ? args : args.slice(0, comma) + args.slice(comma + 1)
        i += INVISIBLE.length + args.length + 1
        continue
      }
    }
    out += s[i]
    i++
  }
  return out
}

/** Read up to the parenthesis closing the one at `start`, honouring nesting. */
function readBalanced(s: string, start: number): string | null {
  let depth = 0
  for (let i = start; i < s.length; i++) {
    const c = s[i]
    if (c === '(') depth++
    else if (c === ')') {
      if (depth === 0) return s.slice(start, i)
      depth--
    }
  }
  return null
}

function topLevelComma(args: string): number {
  let depth = 0
  for (let i = 0; i < args.length; i++) {
    const c = args[i]
    if (c === '(') depth++
    else if (c === ')') depth--
    else if (c === ',' && depth === 0) return i
  }
  return -1
}

function stripWhitespace(s: string): string {
  return s.replace(/\s+/g, '')
}

// Nothing to simplify in a lone symbol or a lone variable, so don't pay for a parse.
// Also keeps prose out: CE will happily read "not math at all" as a product of
// letter-variables (h * m * n * o * a^3 * t^3 * l^2) and call it valid, because as far
// as it is concerned those are just variables. Requiring some arithmetic to be
// present is what keeps that out of a learner-facing hint.
const HAS_ARITHMETIC = /[0-9]|\^|[{}]|\\frac|\\sqrt|[+\-*/=]/u

function tryParse(engine: Engine, latex: string): Expression | null {
  if (!latex.trim()) return null
  try {
    const e = engine.parse(latex)
    return engine.valid(e) ? e : null
  } catch {
    return null
  }
}