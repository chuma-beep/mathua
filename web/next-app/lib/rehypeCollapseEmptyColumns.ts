// rehype plugin: collapse tables left with wholly-empty columns.
//
// Corpus history: the OpenStax ingest strips figures (`ingest_openstax.py`
// decomposes img/svg), leaving layout tables whose image column is entirely
// `<td></td>`. Those render as wide vertical whitespace. This plugin repairs
// them on the parsed tree (never regexes), so it is robust and unit-testable.
//
// Rules:
// 1. Drop a column only if every cell in it — header included — is empty or
//    whitespace (a lone <br> counts as empty). Table 1.1's corner cell
//    survives because its column holds other content.
// 2. A table that is all-empty is removed entirely.
// 3. Width/style/colgroup are stripped from survivors so the remaining
//    column is not pinned at 50%.
// 4. A surviving single column with no <th> was clearly a layout table —
//    unwrap it to plain content. With a <th>, keep the table.
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
      if (keep.every(Boolean)) return

      const remainingCols = keep.filter(Boolean).length
      const keptCells = rows.map((r) => r.filter((_, i) => keep[i])).flat()
      const hasHeader = keptCells.some((c) => isElement(c, 'th'))

      stripSizing(node)
      if (remainingCols === 1 && !hasHeader) {
        const content: HastNode[] = []
        for (const cell of keptCells) content.push(...(cell.children ?? []))
        if (content.every((n) => n.type === 'text' && /^\s*$/.test(n.value ?? ''))) {
          parent.children.splice(index, 1)
        } else {
          parent.children.splice(index, 1, ...content)
        }
        return
      }
      dropCellsFromRows(node, keep)
    }
    walk(tree, null, -1)
  }
}
