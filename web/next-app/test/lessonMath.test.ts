import { describe, it, expect } from 'vitest'
import { prepareLessonMath } from '../lib/lessonMath'

function mathSpans(s: string) {
  const display: string[] = []
  const inline: string[] = []
  const re = /<span class="math-(display|inline)">([\s\S]*?)<\/span>/g
  let m: RegExpExecArray | null
  while ((m = re.exec(s)) !== null) {
    ;(m[1] === 'display' ? display : inline).push(m[2])
  }
  return { display, inline }
}

describe('prepareLessonMath', () => {
  it('never treats escaped currency amounts as math delimiters', () => {
    // The Go pipeline escapes prose currency at ingestion; the frontend must
    // honor the escapes rather than pairing the dollars as math.
    const out = prepareLessonMath('From \\$3.99 each and \\$2.50, the total is \\$20\\$.')
    const { inline, display } = mathSpans(out)
    expect(inline).toEqual([])
    expect(display).toEqual([])
    expect(out).toContain('\\$3.99')
  })

  it('keeps real inline math intact', () => {
    const out = prepareLessonMath('where $x^2 + 5x$ grows')
    expect(mathSpans(out).inline).toEqual(['x^2 + 5x'])
  })

  it('respects escaped dollars as literal text', () => {
    const out = prepareLessonMath('cost \\$5 today')
    expect(mathSpans(out).inline).toEqual([])
    expect(out).toContain('\\$5')
  })

  it('converts single-line display blocks with align to inline-safe aligned', () => {
    const out = prepareLessonMath('\\[\\begin{align} a &= 1 \\tag{2}\\end{align}\\]')
    const { inline } = mathSpans(out)
    expect(inline).toHaveLength(1)
    expect(inline[0]).toContain('\\begin{aligned}')
    expect(inline[0]).not.toContain('\\tag')
    expect(inline[0]).not.toContain('{align}')
  })

  it('protects display math while escaping currency around it', () => {
    const out = prepareLessonMath('$$x = 2$$ then pay \\$4.50\\$ each')
    const { display, inline } = mathSpans(out)
    expect(display).toEqual(['x = 2'])
    expect(inline).toEqual([])
  })
})
