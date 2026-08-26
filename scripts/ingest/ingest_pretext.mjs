/**
 * Ingest ORCCA PreTeXt sources into CommonMark using the official
 * @pretextbook JS stack: PTX -> xast -> mdast -> markdown.
 *
 * Math becomes typed math nodes serialized as $...$/$$...$$ by
 * mdast-util-math; ordinary text dollars are escaped by the serializer,
 * so delimiter ambiguity cannot arise.
 *
 * A thin pre-transform expands constructs the stock converter does not
 * model (quantity units, xref text, webwork statements, figures).
 */
import fs from 'node:fs'
import path from 'node:path'
import { fromXml } from 'xast-util-from-xml'
import { ptxastToMdast } from '@pretextbook/ptxast-util-to-mdast'
import { toMarkdown } from 'mdast-util-to-markdown'
import { mathToMarkdown } from 'mdast-util-math'
import { gfmTableToMarkdown } from 'mdast-util-gfm-table'
import { directiveToMarkdown } from 'mdast-util-directive'

const ORCCA_BASE = 'https://raw.githubusercontent.com/PCCMathSAC/orcca/edition/src'
const ROOT = path.resolve(import.meta.dirname, '../..')
const OUTPUT_DIR = path.join(ROOT, 'data/lessons/teaching')
const MAPPING_PATH = path.join(ROOT, 'data/lessons/mapping.json')

const UNIT_NAMES = {
  'centi meter': 'cm', 'milli meter': 'mm', 'kilo meter': 'km',
  meter: 'm', inch: 'in', foot: 'ft', yard: 'yd', mile: 'mi',
  'centi gram': 'cg', 'milli gram': 'mg', 'kilo gram': 'kg', gram: 'g',
  'centi liter': 'cL', 'milli liter': 'mL', 'kilo liter': 'kL', liter: 'L',
  second: 's', minute: 'min', hour: 'hr', dollar: '\\$',
}

/** Remove \macro{...} calls (balanced braces), keeping the argument. */
function fixBlankUnderscores(v) {
  v = v.replace(/_{2,}/g, '\\;\\;')
  v = v.replace(/(?<!\\)_(?!\\_|[{A-Za-z0-9])/g, '\\_')
  v = v.replace(/\\_\\_/g, '\\;\\;')
  return v
}

function stripMacroCall(s, name) {
  let i
  while ((i = s.indexOf('\\' + name)) !== -1) {
    let j = i + name.length + 1
    while (s[j] === ' ') j++
    if (s[j] !== '{') {
      s = s.slice(0, i) + s.slice(j)
      continue
    }
    let depth = 0
    let k = j
    for (; k < s.length; k++) {
      if (s[k] === '{') depth++
      else if (s[k] === '}') {
        depth--
        if (depth === 0) break
      }
    }
    s = s.slice(0, i) + s.slice(j + 1, k) + s.slice(k + 1)
  }
  return s
}
const el = (name, value, attrs = {}) =>
  value === undefined ? { type: 'element', name, attributes: attrs } :
  { type: 'element', name, attributes: attrs, children: [{ type: 'text', value: String(value) }] }

const text = (v) => ({ type: 'text', value: v })

function localName(node) {
  return String(node?.name || '').replace(/^.*:/, '')
}

/** Build LaTeX for a <quantity> (mag/unit[/per...]) exactly like PreTeXt renders it. */
function quantityLatex(qEl) {
  const pieces = []
  for (const child of qEl.children ?? []) {
    const tag = localName(child)
    if (tag === 'mag') {
      pieces.push((child.children ?? []).map((c) => c.value ?? '').join('').trim())
    } else if (tag === 'unit') {
      const prefix = child.attributes?.prefix
      const base = child.attributes?.base
      const exp = child.attributes?.exp
      let name = UNIT_NAMES[`${prefix} ${base}`] ?? UNIT_NAMES[base] ?? base
      if (!name) name = 'unit'
      let out = `\\text{${name}}`
      if (exp) out += `^{${exp}}`
      if (pieces.length && pieces[pieces.length - 1] !== '/') pieces.push('\\,')
      pieces.push(out)
    } else if (tag === 'per') {
      while (pieces.length && pieces[pieces.length - 1] === '\\,') pieces.pop()
      pieces.push('/')
    }
  }
  const body = pieces.join('').trim()
  return el('m', body)
}

/** Recursively expand ORCCA constructs into elements the stock converter models. */
const BLOCK_IN_P = new Set(['md', 'me', 'men', 'mdn', 'ol', 'ul', 'tabular', 'figure', 'list', 'webwork'])

/** PreTeXt allows <md> inside <p>; the stock converter only models it at
 * block level. Split such paragraphs so display math survives. */
function splitParagraphs(node) {
  if (!node.children) return node
  const out = []
  let run = []
  const flush = () => {
    if (run.length) {
      out.push({ type: 'element', name: 'p', attributes: {}, children: run })
      run = []
    }
  }
  for (const child of node.children) {
    if (BLOCK_IN_P.has(localName(child))) {
      flush()
      out.push(child)
    } else {
      run.push(child)
    }
  }
  flush()
  if (out.length === 1 && out[0].name === 'p') return out[0]
  return out.map((n) => ({ ...n, children: n.children }))
}

function transform(node) {
  if (!node || typeof node !== 'object') return node
  if (node.type === 'text') {
    // Collapse pretty-printed XML whitespace; drop nodes that become empty
    // so mdast-util-to-markdown has no leading-space entities to escape.
    let v = node.value.replace(/\s+/g, ' ')
    // ORCCA step-highlight macro: keep the content (balanced braces).
    v = stripMacroCall(v, 'nextoperation')
    // ORCCA alignment marker anywhere is invalid KaTeX.
    v = v.replace(/\\amp(?![a-zA-Z])/g, '&')
    // Fill-in-the-blank underscores are not subscripts.
    v = fixBlankUnderscores(v)
    return v.trim() === '' ? null : { ...node, value: v }
  }
  if (node.type === 'root') {
    if (node.children) node.children = node.children.flatMap(transform).filter(Boolean)
    return node
  }
  if (node.type !== 'element') return node

  const tag = localName(node)

  // Strip non-content machinery outright.
  if (['idx', 'var', 'pg-code', 'instruction', 'objectives', 'shortdescription', 'latex-image', 'image', 'video', 'sage'].includes(tag)) {
    return null
  }

  // Units: <quantity> -> inline math.
  if (tag === 'quantity') {
    return transform(quantityLatex(node))
  }

  // Cross-references carry their inner word ("Figure", "Appendix"); when the
  // source omits it, derive a generic noun from the target id.
  if (tag === 'xref') {
    const inner = (node.children ?? []).map((c) => c.value ?? '').join('').trim()
    if (inner) return text(` ${inner} `)
    const ref = String(node.attributes?.ref ?? '')
    const kind = ref.split('-')[0]
    const nouns = { figure: 'the figure', table: 'the table', listing: 'the listing', appendix: 'the appendix', subsection: 'the subsection', section: 'the section', chapter: 'the chapter', example: 'the example', exercise: 'the exercise' }
    return text(` ${nouns[kind] ?? ''} `)
  }

  // ORCCA alignment marker is not valid KaTeX, and multi-row <md> needs an
  // aligned envelope around the joined rows.
  if (['md', 'me', 'men', 'mdn', 'mrow'].includes(tag)) {
    const fixed = structuredClone(node)
    const walkText = (n) => {
      if (n.type === 'text') {
        n.value = n.value.replace(/\\amp(?![a-zA-Z])/g, '&')
        n.value = stripMacroCall(n.value, 'nextoperation')
        // Authors sometimes nest $inline$ markers inside <md> TeX; drop them.
        n.value = n.value.replace(/(?<!\\)\$([^$\n]*)(?<!\\)\$/g, '$1')
      }
      for (const c of n.children ?? []) walkText(c)
    }
    walkText(fixed)
    const rows = (fixed.children ?? []).filter((c) => localName(c) === 'mrow')
    const texts = rows.map((r) => (r.children ?? []).map((c) => c.value ?? '').join('').trim())
    const joined = texts.join(' \\\\ ')
    if (tag === 'md' && joined.includes('&')) {
      fixed.children = [text(`\\begin{aligned}${joined}\\end{aligned}`)]
    }
    if (fixed.children) fixed.children = fixed.children.flatMap(transform).filter(Boolean)
    return fixed
  }

  if (tag === 'term') return { ...el('em'), children: (node.children ?? []).flatMap(transform).filter(Boolean) }
  if (tag === 'q') {
    const kids = (node.children ?? []).flatMap(transform).filter(Boolean)
    return [text('"'), ...kids, text('"')]
  }
  if (tag === 'mdash') return text('---')
  if (tag === 'nbsp') return text(' ')
  if (tag === 'times') return text('\\times')

  // Figures: keep only the caption, italicized.
  if (tag === 'figure' || tag === 'aside') {
    const cap = (node.children ?? []).find((c) => localName(c) === 'caption')
    if (cap && node.name === 'figure') {
      const t = (cap.children ?? []).map((c) => c.value ?? '').join('').trim().toLowerCase()
      if (t && t !== 'alternative video lesson' && t !== 'interactive') {
        return { ...el('p'), children: [{ ...el('em'), children: cap.children }] }
      }
    }
    return null
  }

  // Webwork: surface each task's statement and solution as plain blocks.
  if (tag === 'webwork') {
    const out = []
    for (const task of node.children ?? []) {
      if (localName(task) !== 'task') continue
      for (const partTag of ['statement', 'solution']) {
        const part = (task.children ?? []).find((c) => localName(c) === partTag)
        if (!part) continue
        if (partTag === 'solution') out.push(el('p', 'Solution'))
        for (const child of part.children ?? []) {
          const t = transform(child)
          if (t) out.push(t)
        }
      }
    }
    return out
  }

  // PreTeXt <list><dl>: definition-ish list -> markdown list items.
  if (tag === 'list') {
    const titleChild = (node.children ?? []).find((c) => localName(c) === 'title')
    const dl = (node.children ?? []).find((c) => localName(c) === 'dl')
    const items = []
    if (dl) {
      for (const li of dl.children ?? []) {
        if (localName(li) !== 'li') continue
        const liTitle = (li.children ?? []).find((c) => localName(c) === 'title')
        const label = liTitle ? `**${(liTitle.children ?? []).map((c) => c.value ?? '').join('').trim()}**: ` : ''
        const bodyKids = (li.children ?? []).filter((c) => localName(c) !== 'title').flatMap(transform).filter(Boolean)
        items.push({ ...el('li'), children: label ? [text(label), ...bodyKids] : bodyKids })
      }
    }
    if (!items.length) return null
    const out = []
    if (titleChild) out.push(el('p', (titleChild.children ?? []).map((c) => c.value ?? '').join('').trim()))
    out.push({ ...el('ul'), children: items })
    return out
  }

  // Side-by-side panels and task/statement wrappers: splice children through.
  if (['sidebyside', 'sbsgroup', 'statement', 'solution', 'introduction', 'task'].includes(tag)) {
    return (node.children ?? []).flatMap(transform).filter(Boolean)
  }

  // Paragraphs that embed display math / lists need splitting first.
  if (tag === 'p') {
    node.children = (node.children ?? []).flatMap(transform).filter(Boolean)
    return splitParagraphs({ ...node })
  }

  if (node.children) {
    node.children = node.children.flatMap(transform).filter(Boolean)
  }
  return node
}

async function fetchSection(slug) {
  const url = `${ORCCA_BASE}/${slug}.ptx`
  const resp = await fetch(url)
  if (!resp.ok) throw new Error(`${resp.status} for ${url}`)
  return fromXml(await resp.text())
}

const DIRECTIVE_LABELS = {
  example: 'Example', definition: 'Definition', theorem: 'Theorem', fact: 'Fact',
  remark: 'Remark', note: 'Note', warning: 'Warning', proof: 'Proof', exercise: 'Exercise',
  problem: 'Problem', solution: 'Solution', hint: 'Hint', answer: 'Answer',
  lemma: 'Lemma', corollary: 'Corollary', proposition: 'Proposition', conjecture: 'Conjecture',
  axiom: 'Axiom', principle: 'Principle', exploration: 'Exploration', investigation: 'Investigation',
  project: 'Project', activity: 'Activity', task: 'Task', claim: 'Claim', observation: 'Observation',
  insight: 'Insight', demonstration: 'Demonstration',
}

/** Turn container directives (:::example) into plain markdown blocks. */
function expandDirectives(mdast) {
  const out = []
  for (const child of mdast.children ?? []) {
    if (child.type === 'containerDirective') {
      const label = DIRECTIVE_LABELS[child.name]
      if (label) out.push({ type: 'paragraph', children: [{ type: 'text', value: `**${label}**` }] })
      const inner = { type: 'root', children: child.children ?? [] }
      expandDirectives(inner)
      out.push(...inner.children)
      continue
    }
    if (child.children) {
      const inner = { type: 'root', children: child.children }
      expandDirectives(inner)
      child.children = inner.children
    }
    out.push(child)
  }
  mdast.children = out
}

/** Merge adjacent text nodes, trim block edges, drop empties. */
function normalizeMdast(mdast) {
  for (const child of mdast.children ?? []) {
    if (child.children) normalizeMdast(child)
  }
  // Drop empty math nodes (e.g. <m> that only held a randomized <var>).
  const prune = (n) => {
    if (n.children) n.children = n.children.map(prune).filter(Boolean)
    if ((n.type === 'math' || n.type === 'inlineMath') && !String(n.value ?? '').trim()) return null
    return n
  }
  const pruned = prune({ type: 'root', children: mdast.children ?? [] })
  mdast.children = pruned.children

  const merged = []
  for (const child of mdast.children ?? []) {
    const prev = merged[merged.length - 1]
    if (child.type === 'text' && prev?.type === 'text') {
      prev.value += child.value
      continue
    }
    merged.push(child)
  }
  const isBlock = (n) => ['paragraph', 'heading', 'tableCell'].includes(n.type)
  if (isBlock(mdast)) {
    const first = merged[0]
    if (first?.type === 'text') first.value = first.value.replace(/^\s+/, '')
    const last = merged[merged.length - 1]
    if (last?.type === 'text') last.value = last.value.replace(/\s+$/, '')
  }
  mdast.children = merged.filter((c) => !(c.type === 'text' && c.value === ''))
}

function sectionMarkdown(tree) {
  // ptxastToMarkdown expects a root whose children are PreTeXt divisions;
  // the parsed document also carries <xml?> declarations and license comments.
  const section = (tree.children ?? []).find((c) => c.type === 'element' && localName(c) === 'section')
    ?? { type: 'root', children: (tree.children ?? []).filter((c) => c.type === 'element') }
  const root = section.name ? { type: 'root', children: [section] } : section
  const mdast = ptxastToMdast(root)
  expandDirectives(mdast)
  normalizeMdast(mdast)
  let md = toMarkdown(mdast, {
    extensions: [mathToMarkdown(), gfmTableToMarkdown(), directiveToMarkdown()],
    bullet: '-',
  })
  // Pandoc-escaped emphasis markers leak as \* inside math; KaTeX has no \*.
  md = md.replace(/\\\*/g, '*')
  // Aligned content that ended up in an inline span (list items, cells) is
  // invalid KaTeX; promote any $...$ carrying alignment to display.
  md = md.replace(/\$([^$\n]*&[^$\n]*)\$/g, (_m, inner) => `$$\n${inner.trim()}\n$$`)
  return md
}

const ATTRIBUTION =
  '> Content sourced from [ORCCA](https://pcc.edu/orcca) ' +
  '(Open Resources for Community College Algebra) — CC BY 4.0\n'

async function main() {
  const mapping = JSON.parse(fs.readFileSync(MAPPING_PATH, 'utf8'))
  const slugs = new Map()
  for (const [cid, info] of Object.entries(mapping)) {
    if (info.source !== 'orcca') continue
    if (!slugs.has(info.slug)) slugs.set(info.slug, [])
    slugs.get(info.slug).push(cid)
  }

  console.log(`ORCCA sections to extract: ${slugs.size}`)
  fs.mkdirSync(OUTPUT_DIR, { recursive: true })
  let failures = 0

  for (const [slug, conceptIds] of slugs) {
    try {
      const tree = await fetchSection(slug)
      transform(tree)
      const md = ATTRIBUTION + '\n' + sectionMarkdown(tree) + '\n'
      for (const cid of conceptIds) {
        fs.writeFileSync(path.join(OUTPUT_DIR, `${cid}.md`), md)
      }
      console.log(`ok ${slug} (${conceptIds.length})`)
    } catch (err) {
      failures++
      console.error(`X ${slug}: ${err.message}`)
    }
  }
  process.exitCode = failures ? 1 : 0
}

await main()
