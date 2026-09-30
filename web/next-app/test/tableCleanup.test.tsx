import { describe, it, expect, vi } from 'vitest'
import fs from 'node:fs'
import path from 'node:path'
import { render } from '@testing-library/react'
import KatexContent from '../components/KatexContent'
import rehypeCollapseEmptyColumns, { type HastNode, type HastProperties } from '../lib/rehypeCollapseEmptyColumns'

type Node = HastNode

const el = (tagName: string, children: Node[] = [], properties: HastProperties = {}): Node => ({
  type: 'element',
  tagName,
  properties,
  children,
})
const text = (value: string): Node => ({ type: 'text', value })
const td = (children: Node[] = [], properties: HastProperties = {}): Node =>
  el('td', children, properties)
const th = (children: Node[] = []): Node => el('th', children)
const tr = (cells: Node[]): Node => el('tr', cells)
const tableOf = (rows: Node[][]): Node => el('table', rows.map(tr))

function runPlugin(tree: Node): Node {
  rehypeCollapseEmptyColumns()(tree)
  return tree
}

function tablesOf(tree: Node): Node[] {
  const out: Node[] = []
  const walk = (n: Node): void => {
    if (n.type === 'element' && n.tagName === 'table') out.push(n)
    for (const c of n.children ?? []) walk(c)
  }
  walk(tree)
  return out
}

function cellsOf(table: Node): Node[][] {
  const rows: Node[][] = []
  for (const s of table.children ?? []) {
    if (s.type === 'element' && s.tagName === 'tr') {
      rows.push((s.children ?? []).filter((c) => c.type === 'element' && (c.tagName === 'td' || c.tagName === 'th')))
    } else if (s.type === 'element' && (s.tagName === 'tbody' || s.tagName === 'thead')) {
      for (const r of s.children ?? []) {
        if (r.type === 'element' && r.tagName === 'tr') {
          rows.push((r.children ?? []).filter((c) => c.type === 'element' && (c.tagName === 'td' || c.tagName === 'th')))
        }
      }
    }
  }
  return rows
}

describe('rehypeCollapseEmptyColumns (unit)', () => {
  it('unwraps a layout table whose image column is entirely empty', () => {
    const tree = el('root', [
      tableOf([
        [td([text('We start by modeling the first number with 3 blocks.')]), td([])],
        [td([text('Count the total number of blocks.')]), td([el('br')])],
      ]),
    ])
    runPlugin(tree)
    expect(tablesOf(tree)).toHaveLength(0)
    const body = JSON.stringify(tree)
    expect(body).toContain('We start by modeling')
    expect(body).toContain('Count the total number of blocks.')
  })

  it('keeps columns that hold images, svg, or rendered math', () => {
    const tree = el('root', [
      tableOf([
        [td([text('Model the first number.')]), td([el('img', [], { src: 'blocks.png', alt: 'blocks' })])],
        [td([text('Count the total.')]), td([el('svg')])],
        [td([text('Add.')]), td([el('span', [text('x')], { className: ['katex'] })])],
      ]),
    ])
    runPlugin(tree)
    const tables = tablesOf(tree)
    expect(tables).toHaveLength(1)
    expect(cellsOf(tables[0]).map((r) => r.length)).toEqual([2, 2, 2])
  })

  it('keeps a table with an empty corner cell (column has other content)', () => {
    const tree = el('root', [
      tableOf([
        [th([]), th([text('0')]), th([text('1')])],
        [th([text('0')]), td([text('0')]), td([text('1')])],
      ]),
    ])
    runPlugin(tree)
    const tables = tablesOf(tree)
    expect(tables).toHaveLength(1)
    expect(cellsOf(tables[0]).map((r) => r.length)).toEqual([3, 3])
  })

  it('drops only the empty column of a headed table and keeps the <th>', () => {
    const tree = el('root', [
      el('table', [
        el('thead', [tr([th([text('Step')]), th([])])]),
        el('tbody', [
          tr([td([text('Write the numbers.')]), td([])]),
          tr([td([text('Add the ones.')]), td([el('br')])]),
        ]),
      ]),
    ])
    runPlugin(tree)
    const tables = tablesOf(tree)
    expect(tables).toHaveLength(1)
    const rows = cellsOf(tables[0])
    expect(rows.map((r) => r.length)).toEqual([1, 1, 1])
    expect(rows[0][0].tagName).toBe('th')
  })

  it('removes a table that is entirely empty', () => {
    const tree = el('root', [el('p', [text('before')]), tableOf([[td([]), td([text('  ')])]]), el('p', [text('after')])])
    runPlugin(tree)
    expect(tablesOf(tree)).toHaveLength(0)
    expect(JSON.stringify(tree)).toContain('before')
    expect(JSON.stringify(tree)).toContain('after')
  })

  it('skips tables with colspan (column indexes would be guesses)', () => {
    const tree = el('root', [
      tableOf([
        [td([text('wide')], { colSpan: 2 })],
        [td([text('a')]), td([])],
      ]),
    ])
    runPlugin(tree)
    expect(tablesOf(tree)).toHaveLength(1)
    expect(cellsOf(tablesOf(tree)[0])[1]).toHaveLength(2)
  })

  it('strips width/style/colgroup from survivors', () => {
    const tree = el('root', [
      el('table', [
        el('colgroup', [el('col', [], { style: 'width: 50%' })]),
        el('tbody', [
          tr([td([text('Model the 17.')]), td([])]),
        ]),
      ], { width: '100%' }),
    ])
    runPlugin(tree)
    // Single surviving column without <th> unwraps; nothing sized remains.
    expect(tablesOf(tree)).toHaveLength(0)
    expect(JSON.stringify(tree)).not.toContain('colgroup')
    expect(JSON.stringify(tree)).not.toContain('width')
  })
})

function silence() {
  const err = vi.spyOn(console, 'error').mockImplementation(() => {})
  const warn = vi.spyOn(console, 'warn').mockImplementation(() => {})
  return () => {
    err.mockRestore()
    warn.mockRestore()
  }
}

// No table in the output may keep a wholly-empty column (header included).
function emptyColumnCount(container: HTMLElement): number {
  let bad = 0
  for (const table of container.querySelectorAll('table')) {
    const rows: HTMLElement[][] = []
    const sections = table.querySelectorAll(':scope > thead > tr, :scope > tbody > tr, :scope > tr')
    for (const tr of sections) {
      rows.push(Array.from(tr.querySelectorAll(':scope > th, :scope > td')))
    }
    if (rows.length === 0) continue
    const ncols = Math.max(...rows.map((r) => r.length))
    for (let i = 0; i < ncols; i++) {
      const allEmpty = rows.every((r) => {
        const cell = r[i]
        if (!cell) return true
        return (cell.textContent ?? '').trim() === '' && cell.querySelector('img, svg') === null
      })
      if (allEmpty) bad++
    }
  }
  return bad
}

describe('KatexContent table cleanup (integration)', () => {
  it('renders the arith.add gap sections without empty columns', () => {
    const restore = silence()
    try {
      const md = [
        'Each addend is less than 10, so we can use ones blocks.',
        '',
        '|                                                      |     |',
        '|------------------------------------------------------|-----|',
        '| We start by modeling the first number with 3 blocks. |     |',
        '| Then we model the second number with 4 blocks.       |     |',
        '| Count the total number of blocks.                    |     |',
        '',
        'There are 7 blocks in all, so $3 + 4 = 7.$',
      ].join('\n')
      const { container } = render(<KatexContent>{md}</KatexContent>)
      expect(emptyColumnCount(container)).toBe(0)
      expect(container.textContent).toContain('We start by modeling the first number with 3 blocks.')
      expect(container.querySelector('.katex')).not.toBeNull()
    } finally {
      restore()
    }
  })

  it('keeps a legit corner-cell table intact', () => {
    const restore = silence()
    try {
      const md = ['| \\+ | 0 | 1 |', '|----|---|---|', '| 0  | 0 | 1 |', '| 1  | 1 | 2 |'].join('\n')
      const { container } = render(<KatexContent>{md}</KatexContent>)
      const tables = container.querySelectorAll('table')
      expect(tables).toHaveLength(1)
      expect(tables[0].querySelectorAll('tr')[0].querySelectorAll('th, td')).toHaveLength(3)
      expect(emptyColumnCount(container)).toBe(0)
    } finally {
      restore()
    }
  })

  it('leaves a grading-style explanation body structurally unchanged', () => {
    const restore = silence()
    try {
      const md = [
        'Add the ones: $4 + 6 = 10$. Write 0 in the ones place and carry 1.',
        '',
        '| Place | Action |',
        '|-------|--------|',
        '| Ones  | $4 + 6 = 10$ |',
        '| Tens  | $1 + 2 + 8 = 11$ |',
        '',
        'The sum is $110$.',
      ].join('\n')
      const { container } = render(<KatexContent>{md}</KatexContent>)
      const tables = container.querySelectorAll('table')
      expect(tables).toHaveLength(1)
      // Header + 2 body rows, 2 columns each — nothing dropped.
      expect(tables[0].querySelectorAll('tr')).toHaveLength(3)
      for (const tr of tables[0].querySelectorAll('tr')) {
        expect(tr.querySelectorAll('th, td')).toHaveLength(2)
      }
      expect(container.querySelector('.katex')).not.toBeNull()
      expect(container.textContent).toContain('The sum is')
    } finally {
      restore()
    }
  })

  it('repairs the real arith.add.single corpus sections', () => {
    const restore = silence()
    try {
      const file = path.join(__dirname, '../../../data/lessons/teaching/arith.add.single.md')
      const md = fs.readFileSync(file, 'utf8')
      const { container } = render(<KatexContent>{md}</KatexContent>)
      // The stripped-figure columns at :67-71, :89-109, :129-135, :157-191,
      // :433-460 must not survive as empty columns anywhere in the body.
      expect(emptyColumnCount(container)).toBe(0)
      // Content preserved: prose, Table 1.1, identity property.
      expect(container.textContent).toContain('We start by modeling the first number with 3 blocks.')
      expect(container.textContent).toContain('Table 1.1')
      expect(container.textContent).toContain('Identity Property of Addition')
      // Table 1.1 keeps all 11 columns (corner + 0-9).
      const rows = Array.from(container.querySelectorAll('table tr'))
      const wide = rows.find((r) => r.querySelectorAll('th, td').length >= 11)
      expect(wide).toBeTruthy()
    } finally {
      restore()
    }
  })
})

// Mirror of Go extractSection (internal/lessons/lesson.go): heading up to
// the next heading at the same or higher level, so the test renders exactly
// what GET /api/lessons/{id}/kp serves as worked_example.
function extractSection(body: string, heading: string): string {
  const norm = (s: string): string =>
    s
      .replace(/\\\\[()[\]]|\\[()[\]]|\$\$?/g, '')
      .split(/\s+/)
      .join(' ')
      .trim()
  const want = norm(heading)
  const lines = body.split('\n')
  let start = -1
  let startLevel = 0
  for (let i = 0; i < lines.length; i++) {
    const t = lines[i].trim()
    if (!t || !t.startsWith('#')) continue
    let level = 0
    while (level < t.length && t[level] === '#') level++
    const name = norm(t.slice(level))
    if (start < 0) {
      if (name === want) {
        start = i
        startLevel = level
      }
      continue
    }
    if (level <= startLevel && name !== want) return lines.slice(start, i).join('\n')
  }
  if (start < 0) throw new Error(`section not found: ${heading}`)
  return lines.slice(start).join('\n')
}

// Prod gap guard for /learn?concept=arith.add.single: LearnStepper renders
// each KP worked_example under `whitespace-pre-wrap` (LearnStepper.tsx:325),
// which renders surviving whitespace-only text nodes as blank lines —
// invisible to emptyColumnCount. Covers every KP section in the shard
// (KP 0 "Use Addition Notation" plus the gap-bearing "Model Addition of
// Whole Numbers"), asserting both table cleanliness and no blank-line runs.
describe('Learn KP sections (pre-wrap blank lines)', () => {
  it('renders each arith.add.single KP section with no empty columns and no blank-line runs', () => {
    const restore = silence()
    try {
      const kpFile = path.join(__dirname, '../../../data/lessons/kp/arith.add.single.json')
      const kps = JSON.parse(fs.readFileSync(kpFile, 'utf8')) as { section: string }[]
      expect(kps.length).toBeGreaterThan(0)
      const mdFile = path.join(__dirname, '../../../data/lessons/teaching/arith.add.single.md')
      const md = fs.readFileSync(mdFile, 'utf8')
      for (const kp of kps) {
        const body = extractSection(md, kp.section)
        const { container, unmount } = render(
          <KatexContent className="mt-3 text-sm text-mathua-primary whitespace-pre-wrap">{body}</KatexContent>
        )
        try {
          expect(emptyColumnCount(container)).toBe(0)
          const text = container.textContent ?? ''
          const blankRun = text.split('\n').reduce<{ cur: number; max: number }>(
            (acc, line) => {
              const cur = /^\s*$/.test(line) ? acc.cur + 1 : 0
              return { cur, max: Math.max(acc.max, cur) }
            },
            { cur: 0, max: 0 }
          )
          expect(blankRun.max).toBeLessThan(2)
          for (const p of container.querySelectorAll('p')) {
            const empty = (p.textContent ?? '').trim() === '' && p.querySelector('img, svg') === null
            expect(empty).toBe(false)
          }
          expect(container.innerHTML.match(/<br\s*\/?>(\s*<br\s*\/?>)+/i)).toBeNull()
        } finally {
          unmount()
        }
      }
      // The gap section itself must keep its prose after cleanup.
      const gap = extractSection(md, kps[1].section)
      const { container: gapContainer, unmount: gapUnmount } = render(
        <KatexContent className="mt-3 text-sm text-mathua-primary whitespace-pre-wrap">{gap}</KatexContent>
      )
      try {
        expect(gapContainer.textContent).toContain('We start by modeling the first number with 3 blocks.')
      } finally {
        gapUnmount()
      }
    } finally {
      restore()
    }
  })
})
