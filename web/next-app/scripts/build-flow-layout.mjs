// Precompute 2D React Flow (dagre) positions for the concept graph.
//
// Client-side dagre.layout on 657 nodes costs ~2.2s of main-thread time on
// first paint (measured 2026-09-27). Doing it once per build removes that
// cost entirely; the client imports this artifact and falls back to
// client-side dagre only when it is missing or does not cover every concept.
//
// The artifact lives beside graph.meta.json (data/, commit it) rather than in
// public/: it is imported synchronously so the first render already has real
// positions. It is ~22KB — far smaller than the corpus the graph page already
// bundles — and only reaches the dynamically-imported graph chunk.
//
// Canonical input order replicates the Go DAG loader
// (internal/concepts/loader.go: Kahn FIFO seeded in raw file order, with
// prerequisites + encompasses as edges) so the precomputed layout matches
// what the client computes today from /api/graph order. Keep in sync.
//
// Usage: npm run flow:build [-- --print-order]
import { readFileSync, writeFileSync, readdirSync } from 'node:fs'
import { dirname, join } from 'node:path'
import { fileURLToPath } from 'node:url'
import dagre from 'dagre'

const root = join(dirname(fileURLToPath(import.meta.url)), '..')
const conceptsDir = join(root, '..', '..', 'data', 'concepts')

function loadRaw() {
  const files = readdirSync(conceptsDir)
    .filter(f => f.endsWith('.json') && f !== 'enrichment.json')
    .sort()
  const raw = []
  for (const f of files) {
    raw.push(...JSON.parse(readFileSync(join(conceptsDir, f), 'utf8')))
  }
  return raw
}

// Mirror mergeEnrichment (loader.go): id -> encompasses[].
function loadEncompasses() {
  const out = new Map()
  try {
    const entries = JSON.parse(readFileSync(join(conceptsDir, 'enrichment.json'), 'utf8'))
    for (const e of entries) {
      if (e.id && Array.isArray(e.encompasses) && e.encompasses.length > 0) {
        out.set(e.id, e.encompasses)
      }
    }
  } catch {}
  return out
}

// Mirror build() Kahn pass (loader.go:185-217).
function kahnOrder(raw, encompasses) {
  const cmap = new Map(raw.map(c => [c.id, c]))
  const depMap = new Map()
  const encompassedBy = new Map()
  for (const c of raw) {
    for (const pid of c.prerequisites ?? []) {
      if (!cmap.has(pid)) continue
      if (!depMap.has(pid)) depMap.set(pid, [])
      depMap.get(pid).push(c.id)
    }
  }
  for (const c of raw) {
    for (const eid of encompasses.get(c.id) ?? []) {
      if (!cmap.has(eid)) continue
      if (!encompassedBy.has(eid)) encompassedBy.set(eid, [])
      encompassedBy.get(eid).push(c.id)
    }
  }
  const inDegree = new Map()
  for (const c of raw) {
    inDegree.set(c.id, (c.prerequisites ?? []).length + (encompasses.get(c.id) ?? []).length)
  }
  const queue = []
  for (const c of raw) {
    if (inDegree.get(c.id) === 0) queue.push(c.id)
  }
  const order = []
  while (queue.length > 0) {
    const id = queue.shift()
    order.push(id)
    const deps = [...(depMap.get(id) ?? []), ...(encompassedBy.get(id) ?? [])]
    for (const dep of deps) {
      inDegree.set(dep, inDegree.get(dep) - 1)
      if (inDegree.get(dep) === 0) queue.push(dep)
    }
  }
  return { order, byId: cmap }
}

// Mirror layoutWithDagre (components/ConceptGraphFlow.tsx).
function layout(concepts) {
  const g = new dagre.graphlib.Graph()
  g.setGraph({ rankdir: 'LR', nodesep: 16, ranksep: 70, marginx: 20, marginy: 20 })
  g.setDefaultEdgeLabel(() => ({}))
  for (const c of concepts) g.setNode(c.id, { width: 180, height: 40 })
  const known = new Set(concepts.map(c => c.id))
  for (const c of concepts) {
    for (const p of c.prerequisites ?? []) {
      if (known.has(p)) g.setEdge(p, c.id)
    }
  }
  dagre.layout(g)

  const byId = new Map(concepts.map(c => [c.id, c]))
  const domainOrder = []
  const seenDomains = new Set()
  for (const c of concepts) {
    if (!seenDomains.has(c.domain)) {
      seenDomains.add(c.domain)
      domainOrder.push(c.domain)
    }
  }
  const byRank = new Map()
  for (const c of concepts) {
    const n = g.node(c.id)
    if (!n) continue
    const rank = Math.round(n.x)
    if (!byRank.has(rank)) byRank.set(rank, [])
    byRank.get(rank).push(c.id)
  }
  const positions = {}
  byRank.forEach(ids => {
    ids.sort((a, b) => {
      const ca = byId.get(a)
      const cb = byId.get(b)
      const d = domainOrder.indexOf(ca.domain) - domainOrder.indexOf(cb.domain)
      return d !== 0 ? d : ca.label.localeCompare(cb.label)
    })
    let prevDomain = null
    let y = 0
    for (const id of ids) {
      const c = byId.get(id)
      if (prevDomain !== null && c.domain !== prevDomain) y += 28
      positions[id] = [Math.round((g.node(id).x - 90) * 100) / 100, y]
      prevDomain = c.domain
      y += 40 + 16
    }
  })
  return positions
}

const raw = loadRaw()
const encompasses = loadEncompasses()
const { order, byId } = kahnOrder(raw, encompasses)

if (process.argv.includes('--print-order')) {
  console.log(order.join('\n'))
} else {
  const concepts = order.map(id => byId.get(id))
  const t0 = performance.now()
  const positions = layout(concepts)
  const ms = performance.now() - t0
  writeFileSync(join(root, 'data', 'flow.positions.json'), JSON.stringify(positions) + '\n')
  console.log(`flow:build ${Object.keys(positions).length} positions in ${ms.toFixed(0)}ms -> data/flow.positions.json`)
}
