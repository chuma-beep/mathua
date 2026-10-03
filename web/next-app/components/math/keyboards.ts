// The one place `grading_type` and `conceptData.domain` become a keyboard.
//
// Everything downstream — which keycaps a learner sees, and whether the math
// editor or the plain text input is used at all — is derived here. No React
// component branches on a domain or a grading type; they ask this module.
//
// Two independent decisions come out of one lookup:
//
//	editor   'math' | 'text'  → does this concept want a math editor at all?
//	keyboard MathKeyboardMode → which keycaps does it get?
//
// The split matters because a multiple-choice question wants tap targets, not a
// math editor: 213 of 657 concepts are `multiple_choice`, and a math field there
// would be a downgrade.

import { concepts } from '@/lib/conceptData'

export type MathKeyboardMode = 'arithmetic' | 'algebra' | 'fractions' | 'geometry' | 'calculus'

export type AnswerEditor = 'math' | 'text'

/** grading_type → editor. Mirrors internal/grader/router.go. */
const MATH_GRADING_TYPES = new Set(['numeric', 'symbolic', 'expression', 'polynomial', 'complex'])

/**
 * domain → keyboard mode. Domains not listed fall back to 'algebra', which is
 * the most general of the five: it has numbers, the four operations, a fraction
 * and an exponent, so nothing a learner needs is missing. Falling back to
 * 'arithmetic' would hide the letters that a linear-algebra or number-theory
 * question needs, which is the worse failure.
 */
const DOMAIN_MODES: Record<string, MathKeyboardMode> = {
  arithmetic: 'arithmetic',
  prealgebra: 'arithmetic',
  number_theory: 'arithmetic',
  statistics: 'arithmetic',
  fractions: 'fractions',
  algebra: 'algebra',
  abstract_algebra: 'algebra',
  linear_algebra: 'algebra',
  complex_numbers: 'algebra',
  trigonometry: 'algebra',
  precalculus: 'algebra',
  discrete_math: 'algebra',
  topology: 'algebra',
  machine_learning: 'algebra',
  geometry: 'geometry',
  calculus: 'calculus',
  differential_equations: 'calculus',
}

export const DEFAULT_KEYBOARD_MODE: MathKeyboardMode = 'algebra'

export function keyboardModeForDomain(domain?: string): MathKeyboardMode {
  if (!domain) return DEFAULT_KEYBOARD_MODE
  return DOMAIN_MODES[domain] ?? DEFAULT_KEYBOARD_MODE
}

export function editorForGradingType(gradingType?: string): AnswerEditor {
  return gradingType && MATH_GRADING_TYPES.has(gradingType) ? 'math' : 'text'
}

export interface AnswerPresentation {
  editor: AnswerEditor
  keyboard: MathKeyboardMode
}

/**
 * Resolve both decisions from what the server already sends. `engine.Question`
 * carries `concept_id` and `grading_type`; the domain comes from the bundled
 * corpus via `conceptData`. A concept id the corpus does not know (a unit test
 * fixture, a concept added server-side without a corpus rebuild) degrades to the
 * text input rather than guessing — the plain input is the safe default because
 * it is what every concept grades correctly against today.
 */
export function presentationFor(gradingType?: string, domain?: string): AnswerPresentation {
  const editor = editorForGradingType(gradingType)
  return { editor, keyboard: editor === 'math' ? keyboardModeForDomain(domain) : DEFAULT_KEYBOARD_MODE }
}
/**
 * Resolve presentation for a served question.
 *
 * `conceptId` is the only required input: `grading_type` and `domain` both live in
 * the bundled corpus, which is the same place the Go server reads `grading_type`
 * from. A host that already holds the server's `grading_type` passes it and it
 * wins — the server is the authority (ADR-009: every served question carries one).
 *
 * An unknown concept id resolves to the text input, which grades correctly for
 * every concept today, so a corpus gap degrades to the old behaviour rather than
 * to a wrong guess.
 */
export function presentationForQuestion(gradingType?: string, conceptId?: string): AnswerPresentation {
  const concept = conceptId ? concepts.find(c => c.id === conceptId) : undefined
  const type = gradingType ?? concept?.grading_type
  const { editor, keyboard } = presentationFor(type, concept?.domain)
  return { editor, keyboard }
}
