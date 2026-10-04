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
// A single token, not an arbitrary expression: MathLive parenthesises the numerator and
// denominator of every fraction, so `\frac{\pi}{2}` comes back as `(pi)/(2)` exactly as
// `\frac{1}{2}` comes back as `(1)/(2)`. The original pattern only matched digits, so it
// handled the first and missed the second. Restricting each side to one token is what
// keeps this from mangling a real expression: `(x + 1)/(y)` must keep its parentheses,
// because unwrapping it would read as `x + 1/y`.
const FRACTION_PARENS_RE = /\((-?[A-Za-z0-9.]+)\)\/\((-?[A-Za-z0-9.]+)\)/g

/** `4(1)/(10)` → `4 1/10`. Must run before `fixFractionParens`. */
export function fixMixedNumbers(s: string): string {
  return s.replace(MIXED_RE, (_m, whole, num, den) => `${whole} ${num}/${den}`)
}

/** `(1)/(2)` → `1/2`, and `(pi)/(2)` → `pi/2`. */
export function fixFractionParens(s: string): string {
  return s.replace(FRACTION_PARENS_RE, '$1/$2')
}

// Space between a name and its argument list or exponent.
const SPACE_BEFORE_CALL_OR_POW = /(?<=[\w)])\s+(?=[(^])/g

export function fixNameSpacing(s: string): string {
  return s.replace(SPACE_BEFORE_CALL_OR_POW, '')
}

// A run of adjacent single letters separated by single spaces.
const LETTER_RUN = /[A-Za-z](?: [A-Za-z])+/g

/**
 * Rejoin a word the learner's keyboard split into letters.
 *
 * MathLive reads bare adjacent letters as implicit multiplication and serialises them
 * with spaces, so a learner who types `yes` gets back `y e s`, `odd` gets `o d d`, and
 * `addition` gets `a d d i t i o n`. That is a real false miss, not a cosmetic one, and
 * it is the single largest answer shape in the corpus: `yes` is the most common expected
 * answer of all (673 occurrences) and `no` is second (171).
 *
 * It matters because the multiple-choice grader is `strings.EqualFold(expected, answer)`
 * -- an exact comparison -- so `y e s` never matches `yes`.
 *
 * The guard is what keeps this from eating real products. `2pi r` also contains a
 * letter-space-letter, and collapsing it gives `2pir`, which is wrong: the run is part of
 * a product, not a word. So the run is first extended left over its whole adjacent letter
 * sequence, and only collapsed when what precedes *that* is another letter or nothing at
 * all. `y e s` starts the string and collapses; `2pi r` is preceded by a digit and does
 * not.
 *
 * This is not speculation about MathLive: it is measured. `scripts/record-mathlive-corpus.mjs`
 * records what the library actually returns for each of these, and
 * `internal/grader/mathlive_compat_test.go` grades the recording.
 */
export function fixImplicitLetterSpacing(s: string): string {
  return s.replace(LETTER_RUN, (run, offset: number) => {
    let start = offset
    while (start > 0 && /[A-Za-z]/.test(s[start - 1]!)) start--
    let before = start - 1
    while (before >= 0 && s[before] === ' ') before--
    if (before >= 0 && !/[A-Za-z]/.test(s[before]!)) return run
    return run.replace(/ /g, '')
  })
}

/**
 * The plain-text answer submitted to the existing endpoints. Identical wire
 * format to what a plain `<input>` produced; only the fidelity improves, so a
 * structured fraction or exponent now grades instead of erroring.
 */
export function toPlainAnswer(plainText: string): string {
  // Last, so it sees the text the other rules produced. Running it first would work on
  // a string `fixMixedNumbers` was about to rewrite.
  return fixImplicitLetterSpacing(
    fixNameSpacing(fixFractionParens(fixMixedNumbers(fixDivision(plainText)))),
  )
}