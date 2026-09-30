// rehype plugin: repair tables and whitespace left behind by figure stripping.
//
// Corpus history: the OpenStax ingest strips figures (`ingest_openstax.py`
// decomposes img/svg), leaving layout tables whose image column is entirely
// `<td></td>`. Those render as wide vertical whitespace. A second residue is
// inter-block newline runs: rehype-raw hoists newlines around (and out of)
// table markup into long whitespace-only text runs (a clean 2x2 markdown
// table yields 14 newline siblings; Table 1.1 yields 148), which render as
// stacked blank lines wherever an ancestor sets `white-space: pre-wrap`
// (Learn, Review, Diagnostic, Quiz). This plugin repairs both on the parsed
// tree (never regexes), so it is robust and unit-testable.
//
// Rules:
// 1. Drop a column only if every cell in it — header included — is empty or
//    whitespace (a lone <br> counts as empty). Table 1.1's corner cell
//    survives because its column holds other content.
// 2. Drop rows that are entirely empty (header rows included), keeping at
//    least one content row; drop a thead/tbody/tfoot section left with no
//    rows. A layout table's empty `<th></th>` header must not pin the table.
// 3. A table that is all-empty is removed entirely.
// 4. Width/style/colgroup are stripped from survivors so the remaining
//    column is not pinned at 50%.
// 5. A surviving single column with no <th> was clearly a layout table —
//    unwrap it to one <p> per cell (inline-only cells) so captions keep
//    their line breaks instead of running together. With a <th>, keep table.
// 6. Remove paragraphs that are visually empty (whitespace/<br> only).
// 7. Collapse inter-block whitespace-only runs to a single newline (drop
//    leading/trailing runs). A single newline is today's baseline spacing
//    under pre-wrap and is preserved exactly; math spans, pre/code,
//    textarea, script, and style subtrees are never touched.
//
// Runs in KatexContent between rehypeRaw and rehypeSanitize, so Study and
// Learn are both fixed while grading and non-web consumers are untouched.
// Tables using colspan/rowspan are skipped (column indexes would be guesses).

export interface HastProperties {  [key: string]: string | number | boolean | Array<string | number> | undefined
}

export interface HastNode {
  type: string
  tagName?: string
  properties?: HastProperties
  children?: HastNode[]
  value?: string
}

function isElement(node: HastNode, tag?: string): boolean {
  return node.type === 'element' && (tag === undefined || node.tagName === tag)
}

function isCell(node: HastNode): boolean {
  return isElement(node, 'td') || isElement(node, 'th')
}

// Effectively empty: no children, or only whitespace text, <br>, comments.
function isEmptyCell(cell: HastNode): boolean {
  if (!cell.children || cell.children.length === 0) return true
  return cell.children.every((child) => {
    if (child.type === 'text') return /^\s*$/.test(child.value ?? '')
    if (child.type === 'comment') return true
    return isElement(child, 'br')
  })
}

function rowsOf(table: HastNode): HastNode[][] {
  const rows: HastNode[][] = []
  for (const section of table.children ?? []) {
    if (!isElement(section)) continue
    if (section.tagName === 'tr') {
      rows.push((section.children ?? []).filter(isCell))
    } else if (section.tagName === 'thead' || section.tagName === 'tbody' || section.tagName === 'tfoot') {
      for (const tr of section.children ?? []) {
        if (isElement(tr, 'tr')) rows.push((tr.children ?? []).filter(isCell))
      }
    }
  }
  return rows
}

function hasSpan(cells: HastNode[]): boolean {
  return cells.some((c) => {
    const p = c.properties ?? {}
    return Number(p.colSpan ?? p.colspan ?? 1) > 1 || Number(p.rowSpan ?? p.rowspan ?? 1) > 1
  })
}

function stripSizing(node: HastNode): void {
  if (node.properties) {
    delete node.properties.width
    delete node.properties.style
  }
  node.children = (node.children ?? []).filter((child) => {
    if (isElement(child, 'colgroup') || isElement(child, 'col')) return false
    if (isCell(child)) stripSizing(child)
    return true
  })
}

function dropCellsFromRows(table: HastNode, keep: boolean[]): void {
  const rewrite = (tr: HastNode): void => {
    const cells = (tr.children ?? []).filter(isCell)
    const kept = cells.filter((_, i) => keep[i])
    let ki = 0
    tr.children = (tr.children ?? []).map((child) => {
      if (!isCell(child)) return child
      return kept[ki++] ?? null
    }).filter((c): c is HastNode => c !== null)
  }
  for (const section of table.children ?? []) {
    if (!isElement(section)) continue
    if (section.tagName === 'tr') rewrite(section)
    else if (section.tagName === 'thead' || section.tagName === 'tbody' || section.tagName === 'tfoot') {
      for (const tr of section.children ?? []) {
        if (isElement(tr, 'tr')) rewrite(tr)
      }
    }
  }
}

// Drop rows whose cells are all empty, keeping at least one content row.
// Sections (thead/tbody/tfoot) left with no <tr> are removed, so an empty
// `<th></th>` header row cannot pin a layout table into a bordered shell.
function dropEmptyRows(table: HastNode): void {
  const allRows = rowsOf(table)
  if (allRows.filter((r) => r.some((c) => !isEmptyCell(c))).length === 0) return
  const prune = (holder: HastNode): void => {
    holder.children = (holder.children ?? []).filter((child) => {
      if (!isElement(child, 'tr')) return true
      const cells = (child.children ?? []).filter(isCell)
      if (cells.length === 0) return false
      return cells.some((c) => !isEmptyCell(c))
    })
  }
  const sections = (table.children ?? []).filter((c): c is HastNode => isElement(c))
  if (sections.some((s) => s.tagName === 'thead' || s.tagName === 'tbody' || s.tagName === 'tfoot')) {
    for (const s of sections) {
      if (s.tagName === 'thead' || s.tagName === 'tbody' || s.tagName === 'tfoot') prune(s)
    }
    table.children = (table.children ?? []).filter((child) => {
      if (isElement(child) && (child.tagName === 'thead' || child.tagName === 'tbody' || child.tagName === 'tfoot')) {
        return (child.children ?? []).some((c) => isElement(c, 'tr'))
      }
      return true
    })
  } else {
    prune(table)
  }
}

// Block-level tags: cell content containing any of these is spliced as-is;
// inline-only cell content is wrapped in <p> so unwrapped captions keep
// their line breaks instead of running together.
const BLOCK_TAGS = new Set([
  'address', 'article', 'aside', 'blockquote', 'colgroup', 'col', 'details',
  'dialog', 'dd', 'div', 'dl', 'dt', 'fieldset', 'figcaption', 'figure',
  'footer', 'form', 'h1', 'h2', 'h3', 'h4', 'h5', 'h6', 'header', 'hgroup',
  'hr', 'li', 'main', 'nav', 'ol', 'p', 'pre', 'section', 'table',
  'thead', 'tbody', 'tfoot', 'tr', 'td', 'th', 'ul',
])

function isBlock(node: HastNode): boolean {
  return node.type === 'element' && BLOCK_TAGS.has(node.tagName ?? '')
}

function unwrapCellsToContent(cells: HastNode[]): HastNode[] {
  const content: HastNode[] = []
  for (const cell of cells) {
    const kids = cell.children ?? []
    const meaningful = kids.filter((k) => !(k.type === 'text' && /^\s*$/.test(k.value ?? '')))
    if (meaningful.length === 0) continue
    if (meaningful.some(isBlock)) content.push(...kids)
    else content.push({ type: 'element', tagName: 'p', properties: {}, children: kids })
  }
  return content
}

// Visually empty paragraph: whitespace text, <br>, comments only. Such
// paragraphs render as vertical gaps (margins under normal flow, blank
// lines under pre-wrap) and are figure-stripping residue, not content.
function isEmptyParagraph(p: HastNode): boolean {
  if (!isElement(p, 'p')) return false
  const kids = p.children ?? []
  if (kids.length === 0) return true
  return kids.every((k) => {
    if (k.type === 'text') return /^\s*$/.test(k.value ?? '')
    if (k.type === 'comment') return true
    return isElement(k, 'br')
  })
}

const SKIP_SUBTREE_TAGS = new Set(['pre', 'code', 'textarea', 'script', 'style'])

function skipSubtree(node: HastNode): boolean {
  if (node.type === 'element' && SKIP_SUBTREE_TAGS.has(node.tagName ?? '')) return true
  if (node.type !== 'element') return false
  const classes: unknown = node.properties?.className
  if (Array.isArray(classes)) {
    for (const c of classes) {
      const s = String(c)
      if (s === 'math-display' || s === 'math-inline' || s === 'katex') return true
    }
  }
  return false
}

function removeEmptyParagraphs(node: HastNode, inSkip: boolean): void {
  const skip = inSkip || skipSubtree(node)
  if (!skip && Array.isArray(node.children)) {
    for (let i = node.children.length - 1; i >= 0; i--) {
      removeEmptyParagraphs(node.children[i], false)
    }
    if (node.type === 'element') {
      node.children = node.children.filter((c) => !isEmptyParagraph(c))
    }
  } else if (Array.isArray(node.children)) {
    for (const c of node.children) removeEmptyParagraphs(c, true)
  }
}

// Collapse runs of whitespace-only text nodes containing a newline to a
// single "\n" (drop leading/trailing runs). Single newlines are today's
// baseline inter-block spacing and are preserved exactly — only the
// pathological 14–148-newline runs hoisted by rehype-raw are repaired.
function isWsRunNode(node: HastNode): boolean {
  return (
    node.type === 'text' &&
    /^[ \t\r\n]*$/.test(node.value ?? '') &&
    (node.value ?? '').includes('\n')
  )
}

function collapseInterBlockWhitespace(node: HastNode, inSkip: boolean): void {
  const skip = inSkip || skipSubtree(node)
  if (!skip && Array.isArray(node.children)) {
    const kids = node.children
    const out: HastNode[] = []
    let i = 0
    while (i < kids.length) {
      if (isWsRunNode(kids[i])) {
        let j = i
        while (j < kids.length && isWsRunNode(kids[j])) j++
        if (out.length > 0 && j < kids.length) out.push({ type: 'text', value: '\n' })
        i = j
      } else {
        out.push(kids[i])
        i++
      }
    }
    node.children = out
  }
  for (const c of node.children ?? []) collapseInterBlockWhitespace(c, skip)
}

export type RehypeTransformer = (tree: HastNode) => void

export default function rehypeCollapseEmptyColumns(): RehypeTransformer {
  return (tree: HastNode): void => {
    const walk = (node: HastNode, parent: HastNode | null, index: number): void => {
      if (!node || typeof node !== 'object') return
      // Post-order: children first so nested tables resolve before parents.
      if (Array.isArray(node.children)) {
        for (let i = node.children.length - 1; i >= 0; i--) {
          walk(node.children[i], node, i)
        }
      }
      if (!isElement(node, 'table') || !parent || !Array.isArray(parent.children)) return

      const rows = rowsOf(node)
      if (rows.length === 0 || rows.flat().length === 0) {
        parent.children.splice(index, 1)
        return
      }
      if (rows.flat().some((c) => hasSpan([c]))) return

      const ncols = Math.max(...rows.map((r) => r.length))
      const emptyCols: boolean[] = []
      for (let i = 0; i < ncols; i++) {
        emptyCols.push(rows.every((r) => r[i] === undefined || isEmptyCell(r[i])))
      }
      if (emptyCols.every(Boolean)) {
        parent.children.splice(index, 1)
        return
      }
      const keep = emptyCols.map((e) => !e)

      stripSizing(node)
      if (keep.some((k) => !k)) dropCellsFromRows(node, keep)
      dropEmptyRows(node)

      const remaining = rowsOf(node)
      const contentRows = remaining.filter((r) => r.some((c) => !isEmptyCell(c)))
      if (contentRows.length === 0) {
        parent.children.splice(index, 1)
        return
      }
      const remainingCols = Math.max(...remaining.map((r) => r.length))
      const hasHeader = remaining.flat().some((c) => isElement(c, 'th'))

      if (remainingCols === 1 && !hasHeader) {
        const content = unwrapCellsToContent(remaining.flat())
        if (content.every((n) => n.type === 'text' && /^\s*$/.test(n.value ?? ''))) {
          parent.children.splice(index, 1)
        } else {
          parent.children.splice(index, 1, ...content)
        }
        return
      }
    }
    walk(tree, null, -1)
    // Removal of a table can strand adjacent whitespace; empty paragraphs
    // and runs are swept after tables so both orders are covered.
    removeEmptyParagraphs(tree, false)
    collapseInterBlockWhitespace(tree, false)
  }
}
