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

export type GraphAlgorithm = 'bfs' | 'dfs' | 'dijkstra' | 'kruskal'

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
  }
}
