export interface ConceptLayoutInput {
  id: string
  prerequisites: string[]
}

export interface LayoutResult {
  positions: Record<string, [number, number, number]>
  layers: Record<string, number>
  topologicalOrder: string[]
}

export function layoutDAG3D(
  concepts: ConceptLayoutInput[],
  options?: { ySpacing?: number; targetRadius?: number; iterations?: number }
): LayoutResult {
  const ySpacing = options?.ySpacing ?? 2.8
  const targetRadius = options?.targetRadius ?? 18
  const iterations = options?.iterations ?? 80

  const nodeIdSet = new Set(concepts.map(c => c.id))
  const nodeIds = Array.from(nodeIdSet)

  if (nodeIds.length === 0) {
    return { positions: {}, layers: {}, topologicalOrder: [] }
  }

  const children: Record<string, string[]> = {}
  const inDegree: Record<string, number> = {}

  for (const c of concepts) {
    children[c.id] = []
    inDegree[c.id] = c.prerequisites.length
  }

  for (const c of concepts) {
    for (const prereq of c.prerequisites) {
      if (nodeIdSet.has(prereq)) {
        children[prereq].push(c.id)
      }
    }
  }

  const inDegreeCopy: Record<string, number> = {}
  const layers: Record<string, number> = {}
  const queue: string[] = []

  for (const id of nodeIds) {
    inDegreeCopy[id] = inDegree[id]
    if (inDegree[id] === 0) {
      queue.push(id)
      layers[id] = 0
    }
  }

  const topologicalOrder: string[] = []

  while (queue.length > 0) {
    const id = queue.shift()!
    topologicalOrder.push(id)

    const currentLayer = layers[id] ?? 0

    for (const child of children[id]) {
      layers[child] = Math.max(layers[child] ?? 0, currentLayer + 1)
      inDegreeCopy[child]--
      if (inDegreeCopy[child] === 0) {
        queue.push(child)
      }
    }
  }

  for (const id of nodeIds) {
    if (layers[id] === undefined) {
      layers[id] = 0
      topologicalOrder.push(id)
    }
  }

  const layerGroups: Record<number, string[]> = {}
  for (const id of nodeIds) {
    const layer = layers[id]
    if (!layerGroups[layer]) layerGroups[layer] = []
    layerGroups[layer].push(id)
  }

  const numLayers = Object.keys(layerGroups).length
  const layerRadius = Math.max(8, Math.sqrt(concepts.length) * 1.6)

  const centerY = (numLayers - 1) * ySpacing / 2

  const positions: Record<string, [number, number, number]> = {}

  const layerNodeIndices: Record<number, Map<string, number>> = {}
  for (const id of nodeIds) {
    const layer = layers[id]
    if (!layerNodeIndices[layer]) layerNodeIndices[layer] = new Map()
  }
  for (const [layer, ids] of Object.entries(layerGroups)) {
    const map = layerNodeIndices[Number(layer)]
    ids.forEach((id, idx) => map.set(id, idx))
  }

  for (const id of nodeIds) {
    const layer = layers[id]
    const y = layer * ySpacing - centerY

    const indexInLayer = layerNodeIndices[layer].get(id)!
    const totalInLayer = layerGroups[layer].length

    const angle = (Math.PI * 2 * indexInLayer) / totalInLayer + Math.random() * 0.2
    const r = (0.4 + Math.random() * 0.6) * layerRadius
    const x = Math.cos(angle) * r
    const z = Math.sin(angle) * r

    positions[id] = [x, y, z]
  }

  const forces: Record<string, [number, number, number]> = {}
  for (const id of nodeIds) {
    forces[id] = [0, 0, 0]
  }

  const ids = Array.from(nodeIdSet)

  for (let iter = 0; iter < iterations; iter++) {
    const damping = 1 - iter / iterations

    for (const id of nodeIds) {
      forces[id][0] = 0
      forces[id][1] = 0
      forces[id][2] = 0
    }

    const repulsionStrength = 2.5 * damping

    for (let i = 0; i < ids.length; i++) {
      for (let j = i + 1; j < ids.length; j++) {
        const a = ids[i]
        const b = ids[j]
        const pa = positions[a]
        const pb = positions[b]

        const dx = pa[0] - pb[0]
        const dy = pa[1] - pb[1]
        const dz = pa[2] - pb[2]
        const distSq = dx * dx + dy * dy + dz * dz + 0.01
        const dist = Math.sqrt(distSq)

        const force = repulsionStrength / distSq
        const fx = (dx / dist) * force
        const fy = (dy / dist) * force
        const fz = (dz / dist) * force

        forces[a][0] += fx
        forces[a][1] += fy
        forces[a][2] += fz
        forces[b][0] -= fx
        forces[b][1] -= fy
        forces[b][2] -= fz
      }
    }

    const attractionStrength = 0.008 * damping

    for (const c of concepts) {
      for (const prereq of c.prerequisites) {
        if (!nodeIdSet.has(prereq)) continue
        const pa = positions[prereq]
        const pb = positions[c.id]

        const dx = pb[0] - pa[0]
        const dy = pb[1] - pa[1]
        const dz = pb[2] - pa[2]
        const dist = Math.sqrt(dx * dx + dy * dy + dz * dz)

        if (dist > 0.01) {
          const force = dist * attractionStrength
          const fx = (dx / dist) * force
          const fy = (dy / dist) * force
          const fz = (dz / dist) * force

          forces[prereq][0] += fx
          forces[prereq][1] += fy
          forces[prereq][2] += fz
          forces[c.id][0] -= fx
          forces[c.id][1] -= fy
          forces[c.id][2] -= fz
        }
      }
    }

    const layerStiffness = 0.4 * damping

    for (const id of nodeIds) {
      const targetY = layers[id] * ySpacing - centerY
      const dy = targetY - positions[id][1]
      forces[id][1] += dy * layerStiffness
    }

    for (const id of nodeIds) {
      const f = forces[id]
      positions[id][0] += f[0] * damping
      positions[id][1] += f[1] * damping
      positions[id][2] += f[2] * damping
    }
  }

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
    const dx = positions[id][0] - cx
    const dy = positions[id][1] - cy
    const dz = positions[id][2] - cz
    const sq = dx * dx + dy * dy + dz * dz
    if (sq > maxSqDist) maxSqDist = sq
  }
  const maxDist = Math.sqrt(maxSqDist)

  const scale = maxDist > 0 ? targetRadius / maxDist : 1

  for (const id of nodeIds) {
    positions[id][0] = (positions[id][0] - cx) * scale
    positions[id][1] = (positions[id][1] - cy) * scale
    positions[id][2] = (positions[id][2] - cz) * scale
  }

  return { positions, layers, topologicalOrder }
}
