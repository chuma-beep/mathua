export interface ConceptLayoutInput {
  id: string
  prerequisites: string[]
  domain?: string
}

export interface LayoutOptions {
  targetRadius?: number
  seed?: number
  /** Curriculum order for domain sectors; unknown domains append after. */
  domainOrder?: string[]
}

export interface LayoutResult {
  positions: Record<string, [number, number, number]>
  layers: Record<string, number>
  topologicalOrder: string[]
  domainCentroids: Record<string, [number, number, number]>
  /** How many concepts list this node as a prerequisite (hub weight). */
  importance: Record<string, number>
}

const DEFAULT_DOMAIN = 'other'

// Deterministic per-id jitter so identical layers do not read as a perfect
// grid while builds stay byte-stable across machines.
function hashJitter(id: string, salt: number): number {
  let h = salt >>> 0
  for (let i = 0; i < id.length; i++) {
    h = Math.imul(h ^ id.charCodeAt(i), 0x01000193)
  }
  return ((h >>> 0) % 1000) / 1000 - 0.5
}

/**
 * Deterministic, near-linear concept layout.
 *
 * Layers come from a Kahn topological pass, mapped onto latitude rings of a
 * sphere (y axis) so the cloud always reads as a ball rather than a column.
 * Each domain owns a contiguous angular sector proportional to its node
 * count, so clusters read as wedges while prerequisite flow stays vertical.
 * Nodes are placed on a per-layer, per-domain arc: no O(n^2) force loop, no
 * overlap within a layer, and the whole 5k-node graph lays out in
 * milliseconds.
 */
export function layoutDAG3D(concepts: ConceptLayoutInput[], options?: LayoutOptions): LayoutResult {
  const targetRadius = options?.targetRadius ?? 15

  const nodeIds = Array.from(new Set(concepts.map(c => c.id)))
  if (nodeIds.length === 0) {
    return { positions: {}, layers: {}, topologicalOrder: [], domainCentroids: {}, importance: {} }
  }

  const byId = new Map(concepts.map(c => [c.id, c]))
  const nodeIdSet = new Set(nodeIds)
  const children: Record<string, string[]> = {}
  const inDegree: Record<string, number> = {}
  const importance: Record<string, number> = {}

  for (const id of nodeIds) {
    children[id] = []
    inDegree[id] = 0
    importance[id] = 0
  }
  for (const c of concepts) {
    for (const prereq of c.prerequisites) {
      if (!nodeIdSet.has(prereq) || prereq === c.id) continue
      inDegree[c.id] += 1
      children[prereq].push(c.id)
      importance[prereq] += 1
    }
  }

  // Kahn layering: layer(v) = longest path from a root.
  const layers: Record<string, number> = {}
  const remaining = { ...inDegree }
  const queue: string[] = []
  for (const id of nodeIds) {
    if (remaining[id] === 0) {
      queue.push(id)
      layers[id] = 0
    }
  }
  const topologicalOrder: string[] = []
  let head = 0
  while (head < queue.length) {
    const id = queue[head++]
    topologicalOrder.push(id)
    const nextLayer = layers[id] + 1
    for (const child of children[id]) {
      if ((layers[child] ?? 0) < nextLayer) layers[child] = nextLayer
      remaining[child] -= 1
      if (remaining[child] === 0) queue.push(child)
    }
  }
  // Cyclic input (should not happen; validator rejects it) still renders.
  for (const id of nodeIds) {
    if (layers[id] === undefined) {
      layers[id] = 0
      topologicalOrder.push(id)
    }
  }

  // Domain sectors, proportional to node count, in curriculum order.
  const domainIds = new Map<string, string[]>()
  for (const id of nodeIds) {
    const domain = byId.get(id)?.domain || DEFAULT_DOMAIN
    const list = domainIds.get(domain)
    if (list) list.push(id)
    else domainIds.set(domain, [id])
  }
  const requested = options?.domainOrder ?? []
  const domainList = [
    ...requested.filter(d => domainIds.has(d)),
    ...[...domainIds.keys()].filter(d => !requested.includes(d)).sort(),
  ]
  const gap = 0.05
  const available = Math.PI * 2 - gap * domainList.length
  const sectorStart = new Map<string, number>()
  const sectorWidth = new Map<string, number>()
  let cursor = 0
  for (const domain of domainList) {
    const count = domainIds.get(domain)!.length
    const width = available * (count / nodeIds.length)
    sectorStart.set(domain, cursor)
    sectorWidth.set(domain, width)
    cursor += width + gap
  }

  const maxLayer = nodeIds.reduce((m, id) => Math.max(m, layers[id]), 0)
  const layerBuckets = new Map<string, Map<number, string[]>>()
  for (const domain of domainList) {
    const buckets = new Map<number, string[]>()
    for (const id of domainIds.get(domain)!) {
      const layer = layers[id]
      const list = buckets.get(layer)
      if (list) list.push(id)
      else buckets.set(layer, [id])
    }
    for (const list of buckets.values()) {
      list.sort((a, b) => importance[b] - importance[a] || a.localeCompare(b))
    }
    layerBuckets.set(domain, buckets)
  }

  const positions: Record<string, [number, number, number]> = {}
  for (const domain of domainList) {
    const start = sectorStart.get(domain)!
    const width = sectorWidth.get(domain)!
    for (const [layer, ids] of layerBuckets.get(domain)!) {
      // Latitude rings on a sphere: y spans the poles, radius follows the
      // circle so the silhouette stays round at any layer count.
      const h = maxLayer === 0 ? 0 : (layer / maxLayer) * 2 - 1
      const y = h * targetRadius
      const radius = targetRadius * Math.sqrt(Math.max(0.05, 1 - h * h))
      const arc = Math.max(width - gap, 0.001)
      const count = ids.length
      for (let i = 0; i < count; i++) {
        const id = ids[i]
        const t = count === 1 ? 0.5 : (i + 0.5) / count
        const angle = start + gap / 2 + t * arc + hashJitter(id, 0x9e37) * 0.02
        const r = radius * (1 + hashJitter(id, 0x85eb) * 0.08)
        positions[id] = [Math.cos(angle) * r, y, Math.sin(angle) * r]
      }
    }
  }

  // Center and scale to the target radius (keeps camera framing stable).
  let cx = 0
  let cy = 0
  let cz = 0
  for (const id of nodeIds) {
    cx += positions[id][0]
    cy += positions[id][1]
    cz += positions[id][2]
  }
  cx /= nodeIds.length
  cy /= nodeIds.length
  cz /= nodeIds.length

  let maxSqDist = 0
  for (const id of nodeIds) {
    const p = positions[id]
    const dx = p[0] - cx
    const dy = p[1] - cy
    const dz = p[2] - cz
    const sq = dx * dx + dy * dy + dz * dz
    if (sq > maxSqDist) maxSqDist = sq
  }
  const maxDist = Math.sqrt(maxSqDist)
  const scale = maxDist > 0 ? targetRadius / maxDist : 1

  for (const id of nodeIds) {
    const p = positions[id]
    p[0] = (p[0] - cx) * scale
    p[1] = (p[1] - cy) * scale
    p[2] = (p[2] - cz) * scale
  }

  const domainCentroids: Record<string, [number, number, number]> = {}
  for (const domain of domainList) {
    let x = 0
    let y = 0
    let z = 0
    const ids = domainIds.get(domain)!
    for (const id of ids) {
      x += positions[id][0]
      y += positions[id][1]
      z += positions[id][2]
    }
    domainCentroids[domain] = [x / ids.length, y / ids.length, z / ids.length]
  }

  return { positions, layers, topologicalOrder, domainCentroids, importance }
}
