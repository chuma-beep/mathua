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

function escapeHtml(s: string): string {
  return s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')
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

  // --- Convert $...$ and $$...$$ to HTML spans with math-* CSS classes ---
  // This bypasses the remark-math parser entirely, avoiding a stateful
  // tokenizer bug that drops subsequent $...$ expressions in long content.
  // rehype-katex inherently handles elements with math-inline / math-display classes.

  // Display math: $$...$$
  content = content.replace(/\$\$([\s\S]*?)\$\$/g, (_, inner: string) => {
    return '<span class="math-display">' + escapeHtml(unwrapTextCommand(inner.trim())) + '</span>'
  })

  // Inline math: $...$
  content = content.replace(/\$([^$]+?)\$/g, (_, inner: string) => {
    const trimmed = unwrapTextCommand(inner.trim())
    if (trimmed.includes('\n')) {
      return '<span class="math-display">' + escapeHtml(trimmed) + '</span>'
    }
    return '<span class="math-inline">' + escapeHtml(trimmed) + '</span>'
  })

  // Restore escaped literal dollars.
  content = content.split(ESC_DOLLAR).join('\\$')

  return content
}
