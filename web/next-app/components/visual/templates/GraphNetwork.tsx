'use client'

import { useId, useMemo, useState } from 'react'
import FigureFrame from '../FigureFrame'
import { Toggle } from '../Controls'
import { graphSchema } from '../schemas'
import {
  BIPARTITE_GRAPH,
  COLOUR_GRAPH,
  DAG_GRAPH,
  FLOW_GRAPH,
  SAMPLE_GRAPH,
  runGraph,
  type Graph,
  type GraphAlgorithm,
  type GraphNodeId,
} from '../math/graph'
import { scaleLinear } from '../scale'
import { fmt } from '../format'
import type { TemplateProps } from '../types'

const W = 360
const H = 300
const PAD = 28

const ALGORITHMS: Array<{ value: GraphAlgorithm; label: string }> = [
  { value: 'bfs', label: 'Breadth-first' },
  { value: 'dfs', label: 'Depth-first' },
  { value: 'dijkstra', label: 'Dijkstra' },
  { value: 'kruskal', label: 'Kruskal' },
  { value: 'topological', label: 'Topological sort' },
  { value: 'maxflow', label: 'Maximum flow' },
  { value: 'bipartite', label: 'Bipartition' },
  { value: 'colouring', label: 'Colouring' },
]

const COLOUR_FILL = ['var(--accent-blue)', 'var(--accent-green)', 'var(--accent-red)'] as const

const NODE_IDS: GraphNodeId[] = ['A', 'B', 'C', 'D', 'E', 'F']

function graphFor(algorithm: GraphAlgorithm): Graph {
  if (algorithm === 'topological') return DAG_GRAPH
  if (algorithm === 'maxflow') return FLOW_GRAPH
  if (algorithm === 'bipartite') return BIPARTITE_GRAPH
  if (algorithm === 'colouring') return COLOUR_GRAPH
  return SAMPLE_GRAPH
}

export default function GraphNetwork({ params, title, caption }: TemplateProps) {
  const defaults = useMemo(() => graphSchema.parse(params), [params])
  const [algorithm, setAlgorithm] = useState<GraphAlgorithm>(defaults.algorithm)
  const [start, setStart] = useState<GraphNodeId>(defaults.start as GraphNodeId)
  const [showTree, setShowTree] = useState(defaults.showTree)
  const markerId = useId()

  const graph = useMemo(() => graphFor(algorithm), [algorithm])
  const result = useMemo(() => runGraph(graph, algorithm, start), [graph, algorithm, start])
  const treeKeys = useMemo(
    () =>
      new Set(
        result.tree.flatMap((edge) => [`${edge.from}-${edge.to}`, `${edge.to}-${edge.from}`]),
      ),
    [result.tree],
  )
  const orderIndex = new Map(result.order.map((step) => [step.node, step.order]))
  const coords = useMemo(
    () => Object.fromEntries(graph.nodes.map((node) => [node.id, [node.x, node.y] as const])),
    [graph],
  )
  const X = useMemo(() => scaleLinear(-0.35, 2.35, PAD, W - PAD), [])
  const Y = useMemo(() => scaleLinear(1.8, -0.35, PAD, H - PAD), [])

  const reset = () => {
    setAlgorithm(defaults.algorithm)
    setStart(defaults.start as GraphNodeId)
    setShowTree(defaults.showTree)
  }

  const noStart =
    algorithm === 'kruskal' ||
    algorithm === 'topological' ||
    algorithm === 'maxflow' ||
    algorithm === 'bipartite' ||
    algorithm === 'colouring'
  const disabledStart = noStart
  const startText =
    algorithm === 'maxflow' ? 'S → T' : noStart ? 'not used' : start

  const nodeValue = (id: string) => {
    if (algorithm === 'dijkstra') return result.distances?.[id] ?? '∞'
    if (algorithm === 'kruskal' || algorithm === 'maxflow' || algorithm === 'bipartite' || algorithm === 'colouring') {
      return ''
    }
    return orderIndex.get(id) ?? ''
  }

  const flowFor = (from: string, to: string): { flow: number; capacity: number } | null => {
    const edge = result.edgeFlows?.find((item) => item.from === from && item.to === to)
    if (edge) return edge
    const reverse = result.edgeFlows?.find((item) => item.from === to && item.to === from)
    return reverse ? { flow: reverse.flow, capacity: reverse.capacity } : null
  }

  return (
    <FigureFrame
      title={title ?? 'The algorithm decides the path through the graph'}
      hint="Change the algorithm and the start node. The highlighted edges are the ones the algorithm actually keeps, and the labels are its own order, distance or flow values."
      caption={
        caption ??
        'A graph is a set of vertices and weighted edges. Traversal grows in layers or follows one branch, Dijkstra keeps the cheapest distance, Kruskal adds the lightest edge that does not make a cycle, a topological sort removes zero-indegree vertices, maximum flow saturates a cut, a bipartition separates the two sides of an even cycle, and a greedy colouring never gives adjacent vertices the same colour.'
      }
      onReset={reset}
      readout={[
        { label: 'algorithm', value: ALGORITHMS.find((item) => item.value === algorithm)?.label ?? algorithm },
        { label: 'start', value: startText },
        {
          label: 'order',
          value:
            algorithm === 'maxflow'
              ? '—'
              : result.order.length > 0
                ? result.order.map((step) => step.node).join(', ')
                : '—',
        },
        {
          label: 'tree edges',
          value: result.tree.length > 0 ? result.tree.map((edge) => `${edge.from}–${edge.to}`).join(', ') : '—',
        },
        {
          label: 'distances',
          value:
            algorithm === 'dijkstra'
              ? NODE_IDS.map((id) => `${id}: ${fmt(result.distances?.[id] ?? Number.NaN, 2)}`).join(', ')
              : '—',
        },
        {
          label: 'flow',
          value: algorithm === 'maxflow' ? fmt(result.flow ?? 0, 2) : '—',
        },
        {
          label: 'minimum cut',
          value:
            algorithm === 'maxflow' && result.cut
              ? `${result.cut.source.join(', ')} | ${result.cut.sink.join(', ')}`
              : '—',
        },
        {
          label: 'bipartition',
          value: result.sides
            ? `${result.sides.left.join(', ')} | ${result.sides.right.join(', ')}`
            : '—',
        },
        {
          label: 'colours',
          value: result.colours
            ? Object.entries(result.colours)
                .sort(([a], [b]) => a.localeCompare(b))
                .map(([id, colour]) => `${id}: ${colour + 1}`)
                .join(', ')
            : '—',
        },
        {
          label: 'colour count',
          value: result.colourCount === undefined ? '—' : String(result.colourCount),
        },
        { label: 'status', value: result.status },
        ...(result.note ? [{ label: 'note', value: result.note }] : []),
      ]}
    >
      <svg
        viewBox={`0 0 ${W} ${H}`}
        className="mx-auto block h-auto w-full max-w-[420px] font-mono"
        role="img"
        aria-label={`${algorithm} on ${algorithm === 'maxflow' ? 'the flow network' : algorithm === 'topological' ? 'the directed acyclic graph' : 'the sample graph'} from ${startText}`}
      >
        <defs>
          <marker id={markerId} viewBox="0 0 8 8" refX="7" refY="4" markerWidth="5" markerHeight="5" orient="auto">
            <path d="M0,0 L8,4 L0,8 z" fill="var(--text-muted)" />
          </marker>
        </defs>

        {graph.edges.map((edge) => {
          const [ax, ay] = coords[edge.from]
          const [bx, by] = coords[edge.to]
          const highlighted = showTree && treeKeys.has(`${edge.from}-${edge.to}`)
          const flow = flowFor(edge.from, edge.to)
          const directed = algorithm === 'topological' || algorithm === 'maxflow'
          return (
            <g key={`${edge.from}-${edge.to}`}>
              <line
                x1={X(ax)}
                y1={Y(ay)}
                x2={X(bx)}
                y2={Y(by)}
                stroke={highlighted ? 'var(--accent-blue)' : 'var(--text-muted)'}
                strokeWidth={highlighted ? 2.4 : 1}
                markerEnd={directed ? `url(#${markerId})` : undefined}
              />
              <text
                x={(X(ax) + X(bx)) / 2}
                y={(Y(ay) + Y(by)) / 2 - 3}
                textAnchor="middle"
                fontSize="8"
                fill={flow && flow.flow > 0 ? 'var(--accent-blue)' : 'var(--text-muted)'}
              >
                {flow ? `${fmt(flow.flow, 0)}/${fmt(flow.capacity, 0)}` : edge.weight}
              </text>
            </g>
          )
        })}

        {graph.nodes.map((node) => {
          const [x, y] = coords[node.id]
          const label = nodeValue(node.id)
          const current =
            algorithm !== 'kruskal' &&
            algorithm !== 'maxflow' &&
            algorithm !== 'bipartite' &&
            algorithm !== 'colouring' &&
            orderIndex.has(node.id)
          const colour = result.colours?.[node.id]
          const coloured = typeof colour === 'number'
          const inSource = result.cut?.source.includes(node.id)
          const inLeft = result.sides?.left.includes(node.id)
          const inRight = result.sides?.right.includes(node.id)
          return (
            <g key={node.id}>
              <circle
                cx={X(x)}
                cy={Y(y)}
                r={current ? 15 : 13}
                fill={coloured ? COLOUR_FILL[colour % COLOUR_FILL.length] : current ? 'var(--accent-blue)' : 'var(--bg)'}
                stroke={
                  coloured
                    ? COLOUR_FILL[colour % COLOUR_FILL.length]
                    : current
                      ? 'var(--accent-blue)'
                      : inSource || inLeft
                        ? 'var(--accent-teal)'
                        : inRight
                          ? 'var(--accent-red)'
                          : 'var(--text-muted)'
                }
                strokeWidth={inSource || inLeft || inRight ? 2 : 1.2}
              />
              <text
                x={X(x)}
                y={Y(y) + 3}
                textAnchor="middle"
                fontSize="10"
                fill={current || coloured ? 'var(--bg)' : 'var(--text-primary)'}
              >
                {node.label}
              </text>
              {label !== '' ? (
                <text x={X(x)} y={Y(y) - 20} textAnchor="middle" fontSize="8" fill="var(--accent-blue)">
                  {label}
                </text>
              ) : null}
            </g>
          )
        })}

        <text x={PAD} y={PAD + 8} fontSize="9" fill="var(--text-muted)">
          {algorithm === 'maxflow'
            ? `flow ${fmt(result.flow ?? 0, 0)} · cut ${result.cut?.source.join(',') ?? '—'} | ${result.cut?.sink.join(',') ?? '—'}`
            : algorithm === 'topological'
              ? 'source-free order'
              : algorithm === 'bipartite'
                ? result.status === 'ok'
                  ? 'two sides, no edge inside a side'
                  : 'an odd cycle blocks the split'
                : algorithm === 'colouring'
                  ? `${result.colourCount ?? 0} colours from the greedy order`
                  : algorithm === 'kruskal'
                    ? 'the selected forest'
                    : `order from ${start}`}
        </text>
      </svg>

      <div className="mt-3 space-y-2">
        <label className="flex items-center gap-2 font-mono text-[10px] text-mathua-muted">
          Algorithm
          <select
            aria-label="Algorithm"
            value={algorithm}
            onChange={(event) => setAlgorithm(event.target.value as GraphAlgorithm)}
            className="h-7 border border-mathua-border bg-transparent px-1 font-mono text-[10px] text-mathua-secondary"
          >
            {ALGORITHMS.map((item) => (
              <option key={item.value} value={item.value}>
                {item.label}
              </option>
            ))}
          </select>
        </label>
        <label className="flex items-center gap-2 font-mono text-[10px] text-mathua-muted">
          Start
          <select
            aria-label="Start node"
            value={start}
            onChange={(event) => setStart(event.target.value as GraphNodeId)}
            disabled={disabledStart}
            className="h-7 border border-mathua-border bg-transparent px-1 font-mono text-[10px] text-mathua-secondary disabled:opacity-40"
          >
            {NODE_IDS.map((id) => (
              <option key={id} value={id}>
                {id}
              </option>
            ))}
          </select>
        </label>
        <Toggle label="Show tree" checked={showTree} onChange={setShowTree} />
      </div>
    </FigureFrame>
  )
}
