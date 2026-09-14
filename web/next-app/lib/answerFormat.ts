// Expected-answer-form microcopy for typed inputs. The server passes the
// concept's grading_type alongside every question; this table turns it into
// a one-line hint plus the right mobile keyboard. Unknown types fall back
// to a generic line — never blank, never a guess at the format.

export interface AnswerFormat {
  hint: string
  inputMode: 'numeric' | 'text'
}

export function formatForGradingType(gradingType?: string): AnswerFormat {
  switch (gradingType) {
    case 'numeric':
      return { hint: 'Answer with a number', inputMode: 'numeric' }
    case 'expression':
    case 'polynomial':
      return { hint: 'Answer as an expression', inputMode: 'text' }
    case 'symbolic':
      return { hint: 'Answer in simplest form', inputMode: 'text' }
    case 'multiple_choice':
      return { hint: 'Choose one option', inputMode: 'text' }
    case 'comparison':
      return { hint: 'Answer <, >, or =', inputMode: 'text' }
    case 'ordering':
      return { hint: 'Answer in order, separated by commas', inputMode: 'text' }
    case 'tuple':
      return { hint: 'Answer as a pair, e.g. (2, 3)', inputMode: 'text' }
    case 'complex':
      return { hint: 'Answer as a + bi', inputMode: 'text' }
    default:
      return { hint: 'Answer in the form shown', inputMode: 'text' }
  }
}
