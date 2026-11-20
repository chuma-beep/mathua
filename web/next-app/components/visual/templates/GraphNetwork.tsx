'use client'

import { useMemo, useState } from 'react'
import FigureFrame from '../FigureFrame'
import { Toggle } from '../Controls'
import { graphSchema } from '../schemas'
import {
  SAMPLE_GRAPH,
  runGraph,
  type GraphAlgorithm,
  type GraphNodeId,
} from '../math/graph'
import { scaleLinear } from '../scale'
import { fmt } from '../format'
import type { TemplateProps } from '../types'

const W = 360
const H = 300
const PAD = 28

const coords = Object.fromEntries(SAMPLE_GRAPH.nodes.map((node) => [node.id, [node.x, node.y] as const]))
const X = scaleLinear(-0.35, 2.35, PAD, W - PAD)
const Y = scaleLinear(1.35, -0.35, PAD, H - PAD)

const ALGORITHMS: Array<{ value: GraphAlgorithm; label: string }> = [
  { value: 'bfs', label: 'Breadth-first' },
  { value: 'dfs', label: 'Depth-first' },
  { value: 'dijkstra', label: 'Dijkstra' },
  { value: 'kruskal', label: 'Kruskal' },
]

const NODE_IDS: GraphNodeId[] = ['A', 'B', 'C', 'D', 'E', 'F']

export default function GraphNetwork({ params, title, caption }: TemplateProps) {
  const defaults = useMemo(() => graphSchema.parse(params), [params])
  const [algorithm, setAlgorithm] = useState<GraphAlgorithm>(defaults.algorithm)
  const [start, setStart] = useState<GraphNodeId>(defaults.start as GraphNodeId)
  const [showTree, setShowTree] = useState(defaults.showTree)

  const result = useMemo(() => runGraph(SAMPLE_GRAPH, algorithm, start), [algorithm, start])
  const treeEdges = new Set(result.tree.map((edge) => `${edge.from}-${edge.to}`))
  const treeKeys = new Set([...treeEdges].flatMap((key) => {
    const [a, b] = key.split('-')
    return [`${a}-${b}`, `${b}-${a}`]
  }))
  const orderIndex = new Map(result.order.map((step) => [step.node, step.order]))

  const reset = () => {
    setAlgorithm(defaults.algorithm)
    setStart(defaults.start as GraphNodeId)
    setShowTree(defaults.showTree)
  }

  const nodeValue = (id: string) => {
    if (algorithm === 'dijkstra') return result.distances?.[id] ?? '∞'
    if (algorithm === 'kruskal') return ''
    return orderIndex.get(id) ?? ''
  }

  return (
    <FigureFrame
      title={title ?? 'The algorithm decides the path through the graph'}
      hint="Change the algorithm and the start node. The highlighted edges are the ones the algorithm actually keeps, and the labels are its own order or distance values."
      caption={
        caption ??
        'A graph is a set of vertices and weighted edges. Breadth-first search grows in layers, depth-first search follows one branch before backtracking, Dijkstra keeps the cheapest known distance, and Kruskal adds the lightest edge that does not make a cycle.'
      }
      onReset={reset}
      readout={[
        { label: 'algorithm', value: ALGORITHMS.find((item) => item.value === algorithm)?.label ?? algorithm },
        { label: 'start', value: algorithm === 'kruskal' ? 'not used' : start },
        {
          label: 'order',
          value: result.order.length > 0 ? result.order.map((step) => step.node).join(', ') : '—',
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
        { label: 'status', value: result.status },
        ...(result.note ? [{ label: 'note', value: result.note }] : []),
      ]}
    >
      <svg
        viewBox={`0 0 ${W} ${H}`}
        className="mx-auto block h-auto w-full max-w-[420px] font-mono"
        role="img"
        aria-label={`${algorithm} on the sample graph from ${start}`}
      >
        {SAMPLE_GRAPH.edges.map((edge) => {
          const [ax, ay] = coords[edge.from]
          const [bx, by] = coords[edge.to]
          const highlighted = showTree && treeKeys.has(`${edge.from}-${edge.to}`)
          return (
            <g key={`${edge.from}-${edge.to}`}>
              <line
                x1={X(ax)}
                y1={Y(ay)}
                x2={X(bx)}
                y2={Y(by)}
                stroke={highlighted ? 'var(--accent-blue)' : 'var(--text-muted)'}
                strokeWidth={highlighted ? 2.4 : 1}
              />
              <text
                x={(X(ax) + X(bx)) / 2}
                y={(Y(ay) + Y(by)) / 2 - 3}
                textAnchor="middle"
                fontSize="8"
                fill="var(--text-muted)"
              >
                {edge.weight}
              </text>
            </g>
          )
        })}

        {SAMPLE_GRAPH.nodes.map((node) => {
          const [x, y] = coords[node.id]
          const label = nodeValue(node.id)
          const current = algorithm !== 'kruskal' && orderIndex.has(node.id)
          return (
            <g key={node.id}>
              <circle
                cx={X(x)}
                cy={Y(y)}
                r={current ? 15 : 13}
                fill={current ? 'var(--accent-blue)' : 'var(--bg)'}
                stroke={current ? 'var(--accent-blue)' : 'var(--text-muted)'}
                strokeWidth="1.2"
              />
              <text
                x={X(x)}
                y={Y(y) + 3}
                textAnchor="middle"
                fontSize="10"
                fill={current ? 'var(--bg)' : 'var(--text-primary)'}
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
          {algorithm === 'kruskal' ? 'the selected forest' : `order from ${start}`}
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
            disabled={algorithm === 'kruskal'}
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
