// Build the landing graph artifacts from the source concept corpus.
//
// Expensive work (layout) happens here once per build, never on the client:
//   public/data/concepts.positions.json  quantized [x,y,z] map (poster + scene)
//   public/data/concepts.graph.json      nodes + indexed edges + domain counts
//   data/graph.meta.json                 tiny static stats for page sections
//
// Usage: npm run graph:build
import { readFileSync, writeFileSync, mkdirSync } from 'node:fs'
import { dirname, join } from 'node:path'
import { fileURLToPath } from 'node:url'
import { layoutDAG3D } from '../lib/layoutDAG3D.ts'
import { buildGraphMeta, buildGraphPayload } from '../lib/graphPayload.ts'
import { DOMAIN_ORDER } from '../lib/graphDomains.ts'

const root = join(dirname(fileURLToPath(import.meta.url)), '..')
const concepts = JSON.parse(readFileSync(join(root, 'data', 'concepts.json'), 'utf8'))

const t0 = performance.now()
const input = concepts.map(c => ({
  id: c.id,
  prerequisites: c.prerequisites ?? [],
  domain: c.domain,
}))
const { positions, importance } = layoutDAG3D(input, { domainOrder: [...DOMAIN_ORDER] })
const layoutMs = performance.now() - t0

const round2 = n => Math.round(n * 100) / 100
const quantized = {}
for (const [id, p] of Object.entries(positions)) {
  quantized[id] = [round2(p[0]), round2(p[1]), round2(p[2])]
}

const payload = buildGraphPayload(concepts, importance)
const meta = buildGraphMeta(concepts)

mkdirSync(join(root, 'public', 'data'), { recursive: true })
writeFileSync(join(root, 'public', 'data', 'concepts.positions.json'), JSON.stringify(quantized) + '\n')
writeFileSync(join(root, 'public', 'data', 'concepts.graph.json'), JSON.stringify(payload) + '\n')
writeFileSync(join(root, 'data', 'graph.meta.json'), JSON.stringify(meta, null, 2) + '\n')

console.log(
  `graph:build ${concepts.length} nodes / ${payload.edges.length} edges in ${layoutMs.toFixed(0)}ms ` +
    `(layout), ${Object.keys(quantized).length} positions -> public/data/ + data/graph.meta.json`
)
