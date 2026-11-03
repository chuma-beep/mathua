import { describe, expect, it } from 'vitest'
import { prepareLessonMath } from '../lib/lessonMath'

describe('lesson math row breaks', () => {
  it('keeps a cases row break before a single-letter term', () => {
    const source = String.raw`\[\begin{cases} y = x^{2} + 1\\y = x + 1\end{cases}\]`
    const prepared = prepareLessonMath(source)
    expect(prepared).toContain('&#92;&#92;y')
  })
})
