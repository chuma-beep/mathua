import { describe, it, expect } from 'vitest'
import { prepareLessonMath } from '../lib/lessonMath'

function unescapeHtml(s: string): string {
  // Mirrors what rehype-raw does before KaTeX: named entities plus the
  // numeric refs prepareLessonMath uses to keep spans opaque to markdown
  // (`&#92;` backslash, `&#10;` newline).
  return s
    .replace(/&amp;/g, '&')
    .replace(/&lt;/g, '<')
    .replace(/&gt;/g, '>')
    .replace(/&quot;/g, '"')
    .replace(/&#39;/g, "'")
    .replace(/&#(\d+);/g, (_, n: string) => String.fromCharCode(Number(n)))
}

function mathSpans(s: string) {
  const display: string[] = []
  const inline: string[] = []
  const re = /<span class="math-(display|inline)">([\s\S]*?)<\/span>/g
  let m: RegExpExecArray | null
  while ((m = re.exec(s)) !== null) {
    ;(m[1] === 'display' ? display : inline).push(unescapeHtml(m[2]))
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

  it('splits fused adjacent inline maths instead of cascading display pairing', () => {
    // Regression: alg.linear.slope.md $\frac{2}{3}$$\text{gallon}/$ glued a
    // $$ junction that the display pass misread as an opener, swallowing
    // 100+ lines into one span and leaking later math into prose.
    const out = prepareLessonMath('rate $\\frac{2}{3}$$\\text{gallon}/$ here')
    const { display, inline } = mathSpans(out)
    expect(display).toEqual([])
    expect(inline).toHaveLength(2)
    expect(out).not.toContain('$$')
  })

  it('keeps later display blocks intact after a fused junction', () => {
    const out = prepareLessonMath(
      'rate $\\frac{2}{3}$$\\text{gallon}/$ here\n$$\n\\begin{aligned}x &= 1\\end{aligned}\n$$'
    )
    const { display } = mathSpans(out)
    expect(display).toHaveLength(1)
    expect(display[0]).toContain('\\begin{aligned}')
  })

  it('wraps bare environment blocks with no delimiters in display spans', () => {
    // Regression: roots-of-unity.md naked \begin{align}...\end{align} leaked
    // as raw LaTeX in prose (struct gate).
    const out = prepareLessonMath('before\n\\begin{aligned}a &= 1 \\\\ b &= 2\\end{aligned}\nafter')
    const { display } = mathSpans(out)
    expect(display).toHaveLength(1)
    expect(display[0]).toContain('\\begin{aligned}')
    const prose = out.replace(/<span class="math-(display|inline)">[\s\S]*?<\/span>/g, ' ')
    expect(prose).not.toMatch(/\\begin\s*\{/)
  })

  it('does not double-wrap display blocks that already carry $$ delimiters', () => {
    const out = prepareLessonMath('$$\n\\begin{aligned}x &= 1\\end{aligned}\n$$')
    const { display } = mathSpans(out)
    expect(display).toHaveLength(1)
  })

  it('keeps table-cell displays glued to HTML tags intact', () => {
    // Regression: frac.add.word.md <td>$$...$$</td> closers glued to '<'
    // must not be split as fused-inline junctions.
    const out = prepareLessonMath('<td data-align="center">$$\\frac{1}{2} + \\frac{1}{3}$$</td>')
    const { display } = mathSpans(out)
    expect(display).toHaveLength(1)
    const prose = out.replace(/<span class="math-(display|inline)">[\s\S]*?<\/span>/g, ' ')
    expect(prose).not.toMatch(/\\frac/)
  })

  it('wraps bare array environments in table cells', () => {
    const out = prepareLessonMath('<td> \\begin{array}{r}\n89 \\\\\n-61\n\\end{array} </td>')
    const { display } = mathSpans(out)
    expect(display).toHaveLength(1)
  })

  it('wraps nested same-name environments as one span without truncating', () => {
    // Regression: arith.mult tables nest array inside \underset; a
    // non-greedy match ended at the inner \end{array}, producing a KaTeX
    // macro-argument error plus a prose \end leak.
    const out = prepareLessonMath(
      '<td> \\begin{array}{r}\n{47} \\\\\n\\underset{A}{\\begin{array}{r}\n000 \\\\\n{4700}\n\\end{array}} \\\\\n{4,700}\n\\end{array} </td>'
    )
    const { display } = mathSpans(out)
    expect(display).toHaveLength(1)
    expect(display[0]).toContain('{4,700}')
    const prose = out.replace(/<span class="math-(display|inline)">[\s\S]*?<\/span>/g, ' ')
    expect(prose).not.toMatch(/\\(begin|end)\s*\{/)
  })

  it('converts ```math fenced blocks to display math', () => {
    const out = prepareLessonMath('Step 1.\n``` math\n\\text{mean} = \\frac{a}{n}\n```\nDone.')
    const { display } = mathSpans(out)
    expect(display).toHaveLength(1)
    expect(display[0]).toContain('\\frac')
    expect(out).not.toContain('```')
  })

  it('leaves untagged code fences alone', () => {
    const out = prepareLessonMath('```\nconst x = 1\n```')
    const { display, inline } = mathSpans(out)
    expect(display).toEqual([])
    expect(inline).toEqual([])
  })

  it('preserves two-backslash row-break spacing inside wrapped blocks', () => {
    // Regression: collapsing \\[6pt] to \[6pt] broke KaTeX (roots-of-unity).
    const out = prepareLessonMath('\\begin{aligned}\nz_0 &= 1 \\\\[6pt]\nz_1 &= 2\n\\end{aligned}')
    const { display } = mathSpans(out)
    expect(display).toHaveLength(1)
    expect(display[0]).toContain('\\\\[6pt]')
  })
})

describe('stripMathDelimiters', () => {
  it('strips doubled, single, and dollar delimiters to inner text', async () => {
    const { stripMathDelimiters } = await import('../lib/lessonMath')
    expect(stripMathDelimiters('The substitution \\\\( t = \\\\tan(x/2) \\\\)')).toBe(
      'The substitution t = \\tan(x/2)'
    )
    expect(stripMathDelimiters('The $p$-Integral Test')).toBe('The p-Integral Test')
    expect(stripMathDelimiters('plain title')).toBe('plain title')
    expect(stripMathDelimiters('costs \\$5 today')).toBe('costs \\$5 today')
  })
})
