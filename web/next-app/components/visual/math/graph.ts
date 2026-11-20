// Small undirected graphs for traversal, shortest paths and spanning trees.
//
// The algorithms return the edges and order they actually used. A renderer may
// draw a highlight, but the order in the readout is the algorithm's own order,
// not a second implementation.

export type GraphNodeId = string

export interface GraphNode {
  id: GraphNodeId
  label: string
  x: number
  y: number
}

export interface GraphEdge {
  from: GraphNodeId
  to: GraphNodeId
  weight: number
}

export interface Graph {
  nodes: GraphNode[]
  edges: GraphEdge[]
}

export type GraphAlgorithm =
  | 'bfs'
  | 'dfs'
  | 'dijkstra'
  | 'kruskal'
  | 'topological'
  | 'maxflow'
  | 'bipartite'
  | 'colouring'

export interface EdgeFlow {
  from: GraphNodeId
  to: GraphNodeId
  capacity: number
  flow: number
}

export interface TraversalStep {
  node: GraphNodeId
  /** BFS layer or DFS discovery index. */
  order: number
  via?: GraphNodeId
}

export interface TraversalResult {
  order: TraversalStep[]
  /** Edges used to reach a node, or all edges chosen by Kruskal. */
  tree: GraphEdge[]
  /** Dijkstra distance from the source; absent for other algorithms. */
  distances?: Record<GraphNodeId, number>
  /** Maximum flow value; present only for `maxflow`. */
  flow?: number
  /** Source-side and sink-side vertices of the minimum cut. */
  cut?: { source: GraphNodeId[]; sink: GraphNodeId[] }
  /** Flow on every original edge; present only for `maxflow`. */
  edgeFlows?: EdgeFlow[]
  /** The two sides of a bipartition; present only for `bipartite`. */
  sides?: { left: GraphNodeId[]; right: GraphNodeId[] }
  /** Colour index per vertex; present only for `colouring`. */
  colours?: Record<GraphNodeId, number>
  /** Number of colours used; present only for `colouring`. */
  colourCount?: number
  status: 'ok' | 'invalid' | 'disconnected'
  note?: string
}

export const SAMPLE_GRAPH: Graph = {
  nodes: [
    { id: 'A', label: 'A', x: 0, y: 0 },
    { id: 'B', label: 'B', x: 1, y: 0 },
    { id: 'C', label: 'C', x: 2, y: 0 },
    { id: 'D', label: 'D', x: 0, y: 1 },
    { id: 'E', label: 'E', x: 1, y: 1 },
    { id: 'F', label: 'F', x: 2, y: 1 },
  ],
  edges: [
    { from: 'A', to: 'B', weight: 1 },
    { from: 'B', to: 'C', weight: 1 },
    { from: 'A', to: 'D', weight: 1 },
    { from: 'D', to: 'E', weight: 1 },
    { from: 'B', to: 'E', weight: 1 },
    { from: 'E', to: 'F', weight: 1 },
    { from: 'C', to: 'F', weight: 1 },
  ],
}

/** A directed acyclic graph with one edge A→B→D→E and two shortcuts. */
export const DAG_GRAPH: Graph = {
  nodes: [
    { id: 'A', label: 'A', x: 0, y: 0.5 },
    { id: 'B', label: 'B', x: 1, y: 0 },
    { id: 'C', label: 'C', x: 1, y: 1 },
    { id: 'D', label: 'D', x: 2, y: 0.5 },
    { id: 'E', label: 'E', x: 2, y: 1.5 },
  ],
  edges: [
    { from: 'A', to: 'B', weight: 1 },
    { from: 'A', to: 'C', weight: 1 },
    { from: 'B', to: 'D', weight: 1 },
    { from: 'C', to: 'D', weight: 1 },
    { from: 'C', to: 'E', weight: 1 },
    { from: 'D', to: 'E', weight: 1 },
  ],
}

/** A small capacitated network. The maximum flow is 5 and the minimum cut is 5. */
export const FLOW_GRAPH: Graph = {
  nodes: [
    { id: 'S', label: 'S', x: 0, y: 0.5 },
    { id: 'A', label: 'A', x: 1, y: 0 },
    { id: 'B', label: 'B', x: 1, y: 1 },
    { id: 'T', label: 'T', x: 2, y: 0.5 },
  ],
  edges: [
    { from: 'S', to: 'A', weight: 3 },
    { from: 'S', to: 'B', weight: 2 },
    { from: 'A', to: 'B', weight: 1 },
    { from: 'A', to: 'T', weight: 2 },
    { from: 'B', to: 'T', weight: 3 },
  ],
}

export function graphNodes(graph: Graph): GraphNode[] {
  return graph.nodes
}

export function neighbours(graph: Graph, id: GraphNodeId): Array<{ node: GraphNodeId; weight: number; edge: GraphEdge }> {
  const result: Array<{ node: GraphNodeId; weight: number; edge: GraphEdge }> = []
  for (const edge of graph.edges) {
    if (edge.from === id) result.push({ node: edge.to, weight: edge.weight, edge })
    if (edge.to === id) result.push({ node: edge.from, weight: edge.weight, edge })
  }
  return result
}

function findNode(graph: Graph, id: GraphNodeId): GraphNode | undefined {
  return graph.nodes.find((node) => node.id === id)
}

function treeEdge(edge: GraphEdge): GraphEdge {
  return { ...edge }
}

export function breadthFirst(graph: Graph, start: GraphNodeId): TraversalResult {
  if (!findNode(graph, start)) return { order: [], tree: [], status: 'invalid', note: 'The start node is not in the graph.' }
  const seen = new Set<GraphNodeId>([start])
  const queue: Array<{ node: GraphNodeId; layer: number }> = [{ node: start, layer: 0 }]
  const order: TraversalStep[] = []
  const tree: GraphEdge[] = []
  while (queue.length > 0) {
    const current = queue.shift()!
    order.push(current.node === start ? { node: current.node, order: 0 } : { node: current.node, order: current.layer })
    for (const next of neighbours(graph, current.node).sort((a, b) => a.node.localeCompare(b.node))) {
      if (seen.has(next.node)) continue
      seen.add(next.node)
      tree.push(treeEdge(next.edge))
      queue.push({ node: next.node, layer: current.layer + 1 })
    }
  }
  const disconnected = seen.size !== graph.nodes.length
  return {
    order,
    tree,
    status: disconnected ? 'disconnected' : 'ok',
    note: disconnected ? 'The graph is not connected from this start node.' : undefined,
  }
}

export function depthFirst(graph: Graph, start: GraphNodeId): TraversalResult {
  if (!findNode(graph, start)) return { order: [], tree: [], status: 'invalid', note: 'The start node is not in the graph.' }
  const seen = new Set<GraphNodeId>()
  const order: TraversalStep[] = []
  const tree: GraphEdge[] = []
  const visit = (id: GraphNodeId, parent?: GraphNodeId) => {
    seen.add(id)
    order.push({ node: id, order: order.length, via: parent })
    for (const next of neighbours(graph, id).sort((a, b) => a.node.localeCompare(b.node))) {
      if (seen.has(next.node)) continue
      tree.push(treeEdge(next.edge))
      visit(next.node, id)
    }
  }
  visit(start)
  const disconnected = seen.size !== graph.nodes.length
  return {
    order,
    tree,
    status: disconnected ? 'disconnected' : 'ok',
    note: disconnected ? 'The graph is not connected from this start node.' : undefined,
  }
}

export function dijkstra(graph: Graph, start: GraphNodeId): TraversalResult {
  if (!findNode(graph, start)) return { order: [], tree: [], status: 'invalid', note: 'The start node is not in the graph.' }
  const distances: Record<GraphNodeId, number> = {}
  for (const node of graph.nodes) distances[node.id] = Infinity
  distances[start] = 0
  const unvisited = new Set(graph.nodes.map((node) => node.id))
  const order: TraversalStep[] = []
  const tree: GraphEdge[] = []
  while (unvisited.size > 0) {
    let current: GraphNodeId | null = null
    for (const id of unvisited) {
      if (current === null || distances[id] < distances[current]) current = id
    }
    if (current === null || !Number.isFinite(distances[current])) break
    unvisited.delete(current)
    order.push({ node: current, order: order.length })
    for (const next of neighbours(graph, current)) {
      if (!unvisited.has(next.node)) continue
      const candidate = distances[current] + next.weight
      if (candidate < distances[next.node]) {
        distances[next.node] = candidate
        const existing = tree.findIndex((edge) => edge.to === next.node || edge.from === next.node)
        if (existing >= 0) tree.splice(existing, 1)
        tree.push(treeEdge(next.edge))
      }
    }
  }
  const disconnected = order.length !== graph.nodes.length
  return {
    order,
    tree,
    distances,
    status: disconnected ? 'disconnected' : 'ok',
    note: disconnected ? 'Some nodes are unreachable from this start node.' : undefined,
  }
}

export function kruskal(graph: Graph): TraversalResult {
  const parent = new Map<GraphNodeId, GraphNodeId>(graph.nodes.map((node) => [node.id, node.id]))
  const find = (id: GraphNodeId): GraphNodeId => {
    let current = id
    while (parent.get(current) !== current) current = parent.get(current)!
    while (parent.get(id) !== id) {
      const next = parent.get(id)!
      parent.set(id, current)
      id = next
    }
    return current
  }
  const union = (a: GraphNodeId, b: GraphNodeId) => {
    const rootA = find(a)
    const rootB = find(b)
    if (rootA === rootB) return false
    parent.set(rootA, rootB)
    return true
  }
  const sorted = [...graph.edges].sort((a, b) => a.weight - b.weight || a.from.localeCompare(b.from))
  const tree: GraphEdge[] = []
  for (const edge of sorted) {
    if (union(edge.from, edge.to)) tree.push(treeEdge(edge))
  }
  const connected = graph.nodes.length > 0 && tree.length === graph.nodes.length - 1
  return {
    order: tree.map((edge, index) => ({ node: edge.to, order: index, via: edge.from })),
    tree,
    status: connected ? 'ok' : 'disconnected',
    note: connected ? undefined : 'No spanning tree covers every node.',
  }
}

export function topologicalSort(graph: Graph): TraversalResult {
  const indegree = new Map<GraphNodeId, number>(graph.nodes.map((node) => [node.id, 0]))
  for (const edge of graph.edges) {
    if (indegree.has(edge.to)) indegree.set(edge.to, (indegree.get(edge.to) ?? 0) + 1)
  }
  const ready = graph.nodes
    .filter((node) => (indegree.get(node.id) ?? 0) === 0)
    .map((node) => node.id)
    .sort()
  const order: TraversalStep[] = []
  while (ready.length > 0) {
    const current = ready.shift()!
    order.push({ node: current, order: order.length })
    for (const edge of graph.edges) {
      if (edge.from !== current) continue
      const next = (indegree.get(edge.to) ?? 0) - 1
      indegree.set(edge.to, next)
      if (next === 0) {
        ready.push(edge.to)
        ready.sort()
      }
    }
  }
  if (order.length !== graph.nodes.length) {
    return { order, tree: [], status: 'invalid', note: 'The graph contains a directed cycle, so no topological order exists.' }
  }
  return { order, tree: [], status: 'ok' }
}

/** Edmonds–Karp maximum flow. Edge weights are capacities. */
export function maxFlow(graph: Graph, source: GraphNodeId, sink: GraphNodeId): TraversalResult {
  if (!findNode(graph, source) || !findNode(graph, sink) || source === sink) {
    return { order: [], tree: [], status: 'invalid', note: 'Choose two distinct vertices in the graph.' }
  }
  const ids = graph.nodes.map((node) => node.id)
  const indexOf = new Map(ids.map((id, index) => [id, index]))
  const size = ids.length
  const capacity = Array.from({ length: size }, () => Array<number>(size).fill(0))
  for (const edge of graph.edges) {
    const from = indexOf.get(edge.from)
    const to = indexOf.get(edge.to)
    if (from === undefined || to === undefined || from === to) continue
    capacity[from][to] += Math.max(0, edge.weight)
  }
  const residual = capacity.map((row) => [...row])
  const sourceIndex = indexOf.get(source)!
  const sinkIndex = indexOf.get(sink)!
  let total = 0

  for (let guard = 0; guard < 1000; guard++) {
    const parent = Array<number>(size).fill(-1)
    parent[sourceIndex] = sourceIndex
    const queue = [sourceIndex]
    while (queue.length > 0 && parent[sinkIndex] === -1) {
      const current = queue.shift()!
      for (let next = 0; next < size; next++) {
        if (parent[next] !== -1 || residual[current][next] <= 1e-12) continue
        parent[next] = current
        queue.push(next)
        if (next === sinkIndex) break
      }
    }
    if (parent[sinkIndex] === -1) break
    let bottleneck = Number.POSITIVE_INFINITY
    for (let node = sinkIndex; node !== sourceIndex; node = parent[node]) {
      bottleneck = Math.min(bottleneck, residual[parent[node]][node])
    }
    if (!Number.isFinite(bottleneck) || bottleneck <= 1e-12) break
    for (let node = sinkIndex; node !== sourceIndex; node = parent[node]) {
      residual[parent[node]][node] -= bottleneck
      residual[node][parent[node]] += bottleneck
    }
    total += bottleneck
  }

  const reachable = new Set<number>([sourceIndex])
  const queue = [sourceIndex]
  while (queue.length > 0) {
    const current = queue.shift()!
    for (let next = 0; next < size; next++) {
      if (reachable.has(next) || residual[current][next] <= 1e-12) continue
      reachable.add(next)
      queue.push(next)
    }
  }
  const cutSource = ids.filter((_, index) => reachable.has(index))
  const cutSink = ids.filter((_, index) => !reachable.has(index))
  const edgeFlows: EdgeFlow[] = graph.edges.map((edge) => {
    const from = indexOf.get(edge.from)
    const to = indexOf.get(edge.to)
    const flow = from === undefined || to === undefined ? 0 : capacity[from][to] - residual[from][to]
    return { from: edge.from, to: edge.to, capacity: edge.weight, flow: Math.max(0, flow) }
  })
  return {
    order: ids.map((id, index) => ({ node: id, order: index })),
    tree: graph.edges.filter((edge) => (edgeFlows.find((item) => item.from === edge.from && item.to === edge.to)?.flow ?? 0) > 0),
    flow: total,
    cut: { source: cutSource, sink: cutSink },
    edgeFlows,
    status: 'ok',
  }
}

/** A six-cycle: the smallest even cycle, so a two-colouring exists. */
export const BIPARTITE_GRAPH: Graph = {
  nodes: [
    { id: 'A', label: 'A', x: 0, y: 0 },
    { id: 'B', label: 'B', x: 1, y: 0 },
    { id: 'C', label: 'C', x: 2, y: 0 },
    { id: 'D', label: 'D', x: 2, y: 1 },
    { id: 'E', label: 'E', x: 1, y: 1 },
    { id: 'F', label: 'F', x: 0, y: 1 },
  ],
  edges: [
    { from: 'A', to: 'B', weight: 1 },
    { from: 'B', to: 'C', weight: 1 },
    { from: 'C', to: 'D', weight: 1 },
    { from: 'D', to: 'E', weight: 1 },
    { from: 'E', to: 'F', weight: 1 },
    { from: 'F', to: 'A', weight: 1 },
  ],
}

/** A five-cycle with one chord: an odd cycle, so three colours are needed. */
export const COLOUR_GRAPH: Graph = {
  nodes: [
    { id: 'A', label: 'A', x: 0, y: 0.5 },
    { id: 'B', label: 'B', x: 1, y: 0 },
    { id: 'C', label: 'C', x: 2, y: 0.5 },
    { id: 'D', label: 'D', x: 1.5, y: 1.5 },
    { id: 'E', label: 'E', x: 0.5, y: 1.5 },
  ],
  edges: [
    { from: 'A', to: 'B', weight: 1 },
    { from: 'B', to: 'C', weight: 1 },
    { from: 'C', to: 'D', weight: 1 },
    { from: 'D', to: 'E', weight: 1 },
    { from: 'E', to: 'A', weight: 1 },
    { from: 'A', to: 'C', weight: 1 },
  ],
}

function undirectedNeighbours(graph: Graph): Map<GraphNodeId, GraphNodeId[]> {
  const map = new Map<GraphNodeId, Set<GraphNodeId>>(graph.nodes.map((node) => [node.id, new Set()]))
  for (const edge of graph.edges) {
    if (!map.has(edge.from) || !map.has(edge.to) || edge.from === edge.to) continue
    map.get(edge.from)!.add(edge.to)
    map.get(edge.to)!.add(edge.from)
  }
  return new Map([...map].map(([id, set]) => [id, [...set].sort()]))
}

/** Breadth-first two-colouring. A conflict is a proof that the graph is not bipartite. */
export function bipartiteCheck(graph: Graph): TraversalResult {
  const neighbours = undirectedNeighbours(graph)
  const colour = new Map<GraphNodeId, 0 | 1>()
  for (const node of graph.nodes) {
    if (colour.has(node.id)) continue
    colour.set(node.id, 0)
    const queue = [node.id]
    while (queue.length > 0) {
      const current = queue.shift()!
      for (const next of neighbours.get(current) ?? []) {
        const wanted = colour.get(current) === 0 ? 1 : 0
        const seen = colour.get(next)
        if (seen === undefined) {
          colour.set(next, wanted)
          queue.push(next)
        } else if (seen !== wanted) {
          return {
            order: graph.nodes.map((item, index) => ({ node: item.id, order: index })),
            tree: [],
            status: 'invalid',
            note: 'Two neighbours received the same colour, so no bipartition exists.',
          }
        }
      }
    }
  }
  const left = graph.nodes.filter((node) => colour.get(node.id) === 0).map((node) => node.id)
  const right = graph.nodes.filter((node) => colour.get(node.id) === 1).map((node) => node.id)
  return {
    order: graph.nodes.map((item, index) => ({ node: item.id, order: index })),
    tree: [],
    sides: { left, right },
    status: 'ok',
  }
}

/**
 * DSATUR: colour the vertex with the most coloured neighbours first, so the
 * count is a real greedy result and not a hand-written three-colouring.
 */
export function greedyColour(graph: Graph): TraversalResult {
  const neighbours = undirectedNeighbours(graph)
  const colours = new Map<GraphNodeId, number>()
  const order: GraphNodeId[] = []
  while (colours.size < graph.nodes.length) {
    let best: GraphNodeId | null = null
    let bestScore = -1
    let bestDegree = -1
    for (const node of graph.nodes) {
      if (colours.has(node.id)) continue
      const seen = new Set(
        (neighbours.get(node.id) ?? [])
          .map((id) => colours.get(id))
          .filter((value): value is number => value !== undefined),
      )
      const degree = (neighbours.get(node.id) ?? []).length
      if (seen.size > bestScore || (seen.size === bestScore && degree > bestDegree)) {
        best = node.id
        bestScore = seen.size
        bestDegree = degree
      }
    }
    if (best === null) break
    const used = new Set(
      (neighbours.get(best) ?? [])
        .map((id) => colours.get(id))
        .filter((value): value is number => value !== undefined),
    )
    let colour = 0
    while (used.has(colour)) colour++
    colours.set(best, colour)
    order.push(best)
  }
  const values = [...colours.values()]
  return {
    order: order.map((node, index) => ({ node, order: index })),
    tree: [],
    colours: Object.fromEntries(colours),
    colourCount: values.length > 0 ? Math.max(...values) + 1 : 0,
    status: 'ok',
  }
}

export function runGraph(graph: Graph, algorithm: GraphAlgorithm, start: GraphNodeId): TraversalResult {
  switch (algorithm) {
    case 'bfs':
      return breadthFirst(graph, start)
    case 'dfs':
      return depthFirst(graph, start)
    case 'dijkstra':
      return dijkstra(graph, start)
    case 'kruskal':
      return kruskal(graph)
    case 'topological':
      return topologicalSort(graph)
    case 'maxflow':
      return maxFlow(graph, 'S', 'T')
    case 'bipartite':
      return bipartiteCheck(graph)
    case 'colouring':
      return greedyColour(graph)
  }
}
