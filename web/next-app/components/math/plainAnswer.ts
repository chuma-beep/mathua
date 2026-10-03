// Turn what MathLive actually produces into what Mathua's existing graders
// already accept.
//
// MathLive's `getValue('plain-text')` gets most of the way there for free:
// fractions come out as `(1)/(2)`, exponents as `x^2`, roots as `sqrt(9)`, mixed
// numbers as `4(1)/(10)`. Four quirks stand between that and a string the server
// can grade, and each is a real false miss if left alone:
//
//	1. `\div` serializes to the ASCIIMath token `-:`. `7 \div 2` becomes
//	   `7 -: 2`, which SymPy reads as subtraction — silently grading 5 as wrong.
//	2. `\frac` on a whole number serializes as `4(1)/(10)`, so the mixed number
//	   `4 1/10` is indistinguishable from `41/10`. The numeric grader wants the
//	   space.
//	3. `(1)/(2)` is MathLive's parenthesisation, not the learner's; the numeric
//	   grader accepts `1/2` and nothing else.
//	4. `sin (x)` carries a typographic space before the argument list. SymPy's
//	   implicit-multiplication pass only matches `sin(`, so with the space it
//	   parses as a symbol named `sin` multiplied by `x`.
//
// The rules stay separate and exported so a test can name the one that mattered.
// This is a compatibility shim over a third-party serializer, not a parser: it
// never reinterprets an expression, it only removes four documented artifacts.

/** ASCIIMath writes division as `-:`. Restore the operator the server expects. */
export function fixDivision(s: string): string {
  return s.split('-:').join('/')
}

const MIXED_RE = /(-?\d+)\((\d+)\)\/\((\d+)\)/
const FRACTION_PARENS_RE = /\((-?\d+)\)\/\((-?\d+)\)/g

/** `4(1)/(10)` → `4 1/10`. Must run before `fixFractionParens`. */
export function fixMixedNumbers(s: string): string {
  return s.replace(MIXED_RE, (_m, whole, num, den) => `${whole} ${num}/${den}`)
}

/** `(1)/(2)` → `1/2`. */
export function fixFractionParens(s: string): string {
  return s.replace(FRACTION_PARENS_RE, '$1/$2')
}

// Space between a name and its argument list or exponent.
const SPACE_BEFORE_CALL_OR_POW = /(?<=[\w)])\s+(?=[(^])/g

export function fixNameSpacing(s: string): string {
  return s.replace(SPACE_BEFORE_CALL_OR_POW, '')
}

/**
 * The plain-text answer submitted to the existing endpoints. Identical wire
 * format to what a plain `<input>` produced; only the fidelity improves, so a
 * structured fraction or exponent now grades instead of erroring.
 */
export function toPlainAnswer(plainText: string): string {
  return fixNameSpacing(fixFractionParens(fixMixedNumbers(fixDivision(plainText))))
}