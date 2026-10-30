import { describe, it, expect } from 'vitest'
import { prepareLessonMath } from '../lib/lessonMath'

// The Algebrica corpus doubles backslashes inconsistently. A doubled command
// reached KaTeX as a row break followed by the word, which is why the matrix
// page read "m \\ times n" instead of m x n.
describe('doubled LaTeX commands', () => {
  it('collapses a doubled command inside inline math', () => {
    const out = prepareLessonMath(String.raw`dimensions \( m \\times n \) here`)
    expect(out).toContain('&#92;times')
    expect(out).not.toContain('&#92;&#92;times')
  })

  it('keeps a real row break with spacing', () => {
    const out = prepareLessonMath(String.raw`$$\begin{aligned} a &= b \\[6pt] c &= d \end{aligned}$$`)
    // Two backslashes survive; the bracket is HTML-encoded by the span.
    expect(out).toContain('&#92;&#92;')
    expect(out).toContain('6pt')
  })
})
