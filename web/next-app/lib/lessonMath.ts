// Lesson math preparation: the exact string pipeline KatexContent applies
// before handing content to ReactMarkdown/rehype-katex. Extracted so the
// corpus gate (test/latexCorpus.test.ts) exercises the identical code path.

// rehype-katex macros for ORCCA/PreTeXt custom commands used in lesson content,
// plus authoring shortcuts and blackboard letters KaTeX does not define.
// \highlight preserves math mode: its arguments are sometimes math
// (\highlight{\leq}) which \text{} would reject.
// Single source of truth: data/katex-macros.json — imported here and mirrored
// in Go ingestion expectations. Keep in sync.
import macrosJson from '../data/katex-macros.json'
export const lessonMacros: Record<string, string> = macrosJson as Record<string, string>

// Placeholder standing in for an escaped literal \$ so the delimiter regexes
// below never mistake it for a math boundary. Restored at the end.
const ESC_DOLLAR = '\u0000MU-ESC-DOLLAR\u0000'

function escapeHtmlMath(s: string): string {
  // Inside math, & is the aligned/array column separator — must NOT be
  // escaped to &amp; (breaks \begin{aligned} &=). Only escape < > for HTML.
  return s.replace(/</g, '&lt;').replace(/>/g, '&gt;')
}

export function prepareLessonMath(input: string): string {
  let content = input

  // Decode common HTML entities before any LaTeX processing
  content = content
    .replace(/&amp;/g, '&')
    .replace(/&lt;/g, '<')
    .replace(/&gt;/g, '>')
    .replace(/&quot;/g, '"')
    .replace(/&#39;/g, "'")

  // Convert ORCCA-style \amp alignment markers to proper & for LaTeX.
  content = content.replace(/\\amp/g, '&')

  // Escaped literal dollars must never act as math delimiters.
  content = content.replace(/\\\$/g, ESC_DOLLAR)

  // --- Inline math ---
  // Algebrica style: \\(...\\) — unescape \\ back to \ inside
  content = content.replace(/\\\\\(([\s\S]*?)\\\\\)/g, (_, inner) => {
    const clean = inner.replace(/\\\\/g, '\\')
    if (clean.includes('\n')) return '\n$$\n' + clean.trim() + '\n$$\n'
    return '$' + clean + '$'
  })
  // Wikipedia style: \(...\) — backslashes inside are already single
  content = content.replace(/\\\(([\s\S]*?)\\\)/g, (_, inner) => {
    if (inner.includes('\n')) return '\n$$\n' + inner.trim() + '\n$$\n'
    return '$' + inner + '$'
  })

  // --- Display math ---
  // Single-line display blocks become inline $\displaystyle ...$ (e.g. inside
  // table cells). KaTeX forbids {align} and \tag in inline mode, so rewrite
  // align → aligned (inline-legal) and drop \tag — otherwise these blocks
  // hard-fail at render time.
  const inlineSafe = (s: string) =>
    s
      .replace(/\\begin\{align\*?\}/g, '\\begin{aligned}')
      .replace(/\\end\{align\*?\}/g, '\\end{aligned}')
      .replace(/\\tag\{[^{}]*\}/g, '')

  // Algebrica style: \\[...\\]
  content = content.replace(/\\\\\[([\s\S]*?)\\\\\]/g, (_, inner) => {
    const clean = inner.replace(/\\\\/g, '\\')
    if (inner.includes('\n')) {
      return '\n$$\n' + clean.trim() + '\n$$\n'
    }
    return '$\\displaystyle ' + inlineSafe(clean.trim()) + '$'
  })
  // Wikipedia style: \[...\]
  content = content.replace(/\\\[([\s\S]*?)\\\]/g, (_, inner) => {
    if (inner.includes('\n')) {
      return '\n$$\n' + inner.trim() + '\n$$\n'
    }
    return '$\\displaystyle ' + inlineSafe(inner.trim()) + '$'
  })

  // ORCCA content wraps multi-line display math (aligned, array) in single
  // $...$ where the $ sits ALONE on its opening and closing lines. The line
  // anchoring must be strict: a loose ^\s*\$ opener also matches inline
  // $...$ pairs that merely start a line, pairing them with a far-away $
  // and cascading mispaired spans across the file.
  content = content.replace(/^[ \t]*\$(?!\$)\n([\s\S]*?)\n[ \t]*\$(?!\$)[ \t]*$/gm, '$$\n$1\n$$')

  // Scrape artifacts: math commands wrongly wrapped in \text{} — e.g.
  // \text{\pm}, \text{\sin}, \text{-} inside $...$ (KaTeX rejects \text{\pm}).
  const unwrapTextCommand = (s: string) =>
    s.replace(/\\text\{\\([a-zA-Z]+)\}/g, '\\$1').replace(/\\text\{([/\\-])\}/g, '$1')

  // ```math fenced blocks (OpenStax/PreTeXt extracts) are unambiguous display
  // math that carries no $ delimiters. Convert to $$ before the display pass
  // so they render as math instead of <pre><code>. Only the explicit `math`
  // tag converts; untagged fences stay code. Indented fences inside numbered
  // lists (e.g. stat.prob How-To steps) match via leading-space tolerance.
  content = content.replace(/^[ \t]*```[ \t]*math[ \t]*\n([\s\S]*?)\n[ \t]*```/gm, (_, inner: string) => {
    return '$$\n' + inner.trim() + '\n$$'
  })

  // --- Convert $...$ and $$...$$ to HTML spans with math-* CSS classes ---
  // This bypasses the remark-math parser entirely, avoiding a stateful
  // tokenizer bug that drops subsequent $...$ expressions in long content.
  // rehype-katex inherently handles elements with math-inline / math-display classes.

  // Fused-inline defense: two adjacent inline maths ($...$$...$) leave a
  // $$ junction that the display pass below would misread as a display
  // opener, cascading every later pairing (e.g. slope.md
  // $\frac{2}{3}$$\text{gallon}/$ here swallowed 100+ lines into one span).
  // Split only genuine junctions: close-side is math content (never an HTML
  // tag bracket, so <td>$$...$$</td> displays survive) and open-side starts
  // new math (backslash, word char, or paren — never '<', '.', ',', ...).
  // NOTE: '$$' in the replacement pattern is a literal '$', so the
  // replacement emits close-$ + space + open-$ (dollars preserved, separated).
  content = content.replace(/([^\s$><])\$\$([\\\w(])/g, '$1$ $$$2')

  // Display math: $$...$$
  content = content.replace(/\$\$([\s\S]*?)\$\$/g, (_, inner: string) => {
    return '<span class="math-display">' + escapeHtmlMath(unwrapTextCommand(inner.trim())) + '</span>'
  })

  // Inline math: $...$
  content = content.replace(/\$([^$]+?)\$/g, (_, inner: string) => {
    const trimmed = unwrapTextCommand(inner.trim())
    if (trimmed.includes('\n')) {
      return '<span class="math-display">' + escapeHtmlMath(trimmed) + '</span>'
    }
    return '<span class="math-inline">' + escapeHtmlMath(trimmed) + '</span>'
  })

  // Bare-environment defense: standalone \begin{...}...\end{...} blocks with
  // no $ delimiters (scraper/Go-envOutput output, e.g. roots-of-unity.md
  // \begin{align}) never become spans above and leak as raw LaTeX in prose.
  // Wrap them in display spans, operating only on prose outside the spans
  // emitted above so $$-wrapped blocks are never double-wrapped.
  // Matching counts nested same-name envs (arith tables nest array inside
  // \underset) so spans are never truncated mid-argument.
  const bareEnvNames = ['align', 'align*', 'aligned', 'equation', 'equation*', 'gather', 'gather*', 'array', 'matrix', 'pmatrix', 'bmatrix', 'vmatrix', 'cases']
  const wrapBareEnv = (prose: string): string => {
    const openRe = /\\begin\{([a-z]+\*?)\}/g
    let out = ''
    let pos = 0
    let m: RegExpExecArray | null
    while ((m = openRe.exec(prose)) !== null) {
      if (!bareEnvNames.includes(m[1])) continue
      const name = m[1]
      const tokRe = new RegExp('\\\\(begin|end)\\{' + name.replace('*', '\\*') + '\\}', 'g')
      tokRe.lastIndex = m.index
      let depth = 0
      let closeEnd = -1
      let t: RegExpExecArray | null
      while ((t = tokRe.exec(prose)) !== null) {
        depth += t[1] === 'begin' ? 1 : -1
        if (depth === 0) {
          closeEnd = t.index + t[0].length
          break
        }
      }
      if (closeEnd < 0) continue // unbalanced: leave untouched
      // Over-escaped scrape runs (3+ backslashes) collapse to a real row
      // break (\\). Exactly-two sequences are valid KaTeX row-break spacing
      // (\\[6pt]) and must survive untouched.
      const block = unwrapTextCommand(prose.slice(m.index, closeEnd).trim()).replace(/\\{3,}/g, '\\\\')
      out += prose.slice(pos, m.index) + '<span class="math-display">' + escapeHtmlMath(block) + '</span>'
      pos = closeEnd
      openRe.lastIndex = closeEnd
    }
    return out + prose.slice(pos)
  }
  content = content
    .split(/(<span class="math-(?:display|inline)">[\s\S]*?<\/span>)/g)
    .map((part) => (part.startsWith('<span class="math-') ? part : wrapBareEnv(part)))
    .join('')

  // Restore escaped literal dollars.
  content = content.split(ESC_DOLLAR).join('\\$')

  return content
}
