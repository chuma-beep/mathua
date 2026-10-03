// What the rest of Mathua sees. Deliberately small: a caller needs a LaTeX string
// for display, a plain-text string for the existing graders, and nothing else.
//
// `mathJson` is reserved and currently always undefined. MathLive 0.111 cannot
// produce MathJSON — `getValue('mathjson')` and `convertLatexToMathJson` were
// both removed from the library, leaving Compute Engine as the only parser of
// record. Commit ② fills this in; until then it is absent rather than wrong, and
// no call site depends on it.
export type MathInputValue = {
  latex: string
  plainAnswer: string
  mathJson?: unknown
}

export type MathInputStatus = 'default' | 'correct' | 'incorrect' | 'disabled'