'use client'

import { useMemo, useEffect } from 'react'
import {
  ReactFlow,
  Background,
  Controls,
  Handle,
  Position,
  useNodesState,
  useEdgesState,
  type Node,
  type Edge,
} from '@xyflow/react'
import '@xyflow/react/dist/base.css'
import { useTheme } from '../hooks/useTheme'
import { themeColors } from './FlowDiagram'

const monoFont = "'IBM Plex Mono', monospace"

function LayerNode({ data }: { data: { label: string } }) {
  const { theme } = useTheme()
  const c = themeColors[theme === 'dark' ? 'dark' : 'light']
  return (
    <div
      style={{
        background: 'transparent',
        border: '0.5px solid ' + c.border,
        height: '100%',
        fontFamily: monoFont,
        fontSize: '10px',
        textTransform: 'uppercase',
        letterSpacing: '0.1em',
        color: c.textMuted,
        padding: '6px 10px',
        display: 'flex',
        flexDirection: 'column',
      }}
    >
      <span style={{ marginBottom: 6 }}>{data.label}</span>
    </div>
  )
}

function FlowNode({ data }: { data: { label: string } }) {
  const { theme } = useTheme()
  const c = themeColors[theme === 'dark' ? 'dark' : 'light']
  return (
    <div
      style={{
        background: c.surface,
        border: '0.5px solid ' + c.border,
        fontFamily: monoFont,
        fontSize: '10px',
        color: c.textSecondary,
        padding: '6px 10px',
        textAlign: 'center',
        lineHeight: 1.3,
        borderRadius: 0,
        whiteSpace: 'pre-wrap',
        height: '100%',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
      }}
    >
      <Handle type="target" position={Position.Top} style={{ opacity: 0, pointerEvents: 'none' }} />
      {data.label}
      <Handle type="source" position={Position.Bottom} style={{ opacity: 0, pointerEvents: 'none' }} />
    </div>
  )
}

const nodeTypes = { layerNode: LayerNode, flowNode: FlowNode }

interface LayerItem {
  id: string
  label: string
}

interface LayerDef {
  id: string
  label: string
  items: LayerItem[]
}

interface ArchEdge {
  source: string
  target: string
}

const layers: LayerDef[] = [
  {
    id: 'ui', label: 'UI Layer',
    items: [
      { id: 'web', label: 'Web Browser\nReact + KaTeX' },
    ],
  },
  {
    id: 'api', label: 'API Layer',
    items: [
      { id: 'session', label: 'Session\nService' },
      { id: 'graph', label: 'Graph\nService' },
      { id: 'leaderboard', label: 'Leaderboard\nService' },
      { id: 'diag', label: 'Diagnostic\nService' },
    ],
  },
  {
    id: 'engine', label: 'Core Engine',
    items: [
      { id: 'dag', label: 'DAG\nLoader' },
      { id: 'sched', label: 'Scheduler\nSM-2' },
      { id: 'gen', label: 'Generators' },
      { id: 'grade', label: 'Graders' },
      { id: 'score', label: 'Scoring' },
    ],
  },
  {
    id: 'grading', label: 'Grading Layer',
    items: [
      { id: 'numeric', label: 'Numeric\nGo' },
      { id: 'sympy', label: 'SymPy\nPython' },
      { id: 'symbolic', label: 'Symbolic\nGo' },
      { id: 'choice', label: 'Choice\nGo' },
      { id: 'comp', label: 'Comparison\nGo' },
      { id: 'order', label: 'Ordering\nGo' },
    ],
  },
  {
    id: 'storage', label: 'Storage',
    items: [
      { id: 'sqlite', label: 'SQLite' },
      { id: 'pg', label: 'PostgreSQL\nweb' },
    ],
  },
]

const archEdges: ArchEdge[] = [
  { source: 'web', target: 'session' },
  { source: 'web', target: 'graph' },
  { source: 'web', target: 'leaderboard' },
  { source: 'web', target: 'diag' },
  { source: 'session', target: 'dag' },
  { source: 'session', target: 'sched' },
  { source: 'graph', target: 'dag' },
  { source: 'diag', target: 'sched' },
  { source: 'dag', target: 'sched' },
  { source: 'sched', target: 'gen' },
  { source: 'gen', target: 'grade' },
  { source: 'grade', target: 'score' },
  { source: 'gen', target: 'numeric' },
  { source: 'gen', target: 'sympy' },
  { source: 'gen', target: 'symbolic' },
  { source: 'gen', target: 'choice' },
  { source: 'gen', target: 'comp' },
  { source: 'gen', target: 'order' },
  { source: 'score', target: 'sqlite' },
  { source: 'score', target: 'pg' },
]

const ITEM_W = 110
const ITEM_H = 44
const LAYER_PAD = 8
const LAYER_GAP_X = 8
const LAYER_GAP_Y = 16
const ITEM_GAP_X = 8
const ITEM_GAP_Y = 8

function buildLayout() {
  let totalH = 0
  const layerPositions: Record<string, { x: number; y: number; w: number; h: number }> = {}
  const itemPositions: Record<string, { x: number; y: number }> = {}

  layers.forEach((layer) => {
    const perRow = layer.items.length <= 3 ? layer.items.length : Math.ceil(layer.items.length / 2)
    const rows = Math.ceil(layer.items.length / perRow)
    const w = perRow * ITEM_W + (perRow - 1) * ITEM_GAP_X + LAYER_PAD * 2
    const h = rows * ITEM_H + (rows - 1) * ITEM_GAP_Y + LAYER_PAD * 2 + 16

    const lx = 0
    const ly = totalH
    layerPositions[layer.id] = { x: lx, y: ly, w, h }

    layer.items.forEach((item, idx) => {
      const row = Math.floor(idx / perRow)
      const col = idx % perRow
      itemPositions[item.id] = {
        x: lx + LAYER_PAD + col * (ITEM_W + ITEM_GAP_X),
        y: ly + LAYER_PAD + 16 + row * (ITEM_H + ITEM_GAP_Y),
      }
    })

    totalH += h + LAYER_GAP_Y
  })

  const totalW = Math.max(...Object.values(layerPositions).map(l => l.w))

  return { layerPositions, itemPositions, totalW, totalH }
}

export default function ArchitectureFlow() {
  const { theme, mounted } = useTheme()
  const c = themeColors[theme === 'dark' ? 'dark' : 'light']

  const layout = buildLayout()

  const initialNodes: Node[] = []
  layers.forEach((layer) => {
    const lp = layout.layerPositions[layer.id]
    initialNodes.push({
      id: layer.id,
      type: 'layerNode',
      position: { x: lp.x, y: lp.y },
      data: { label: layer.label },
      style: { width: lp.w, height: lp.h },
    })
    layer.items.forEach((item) => {
      const ip = layout.itemPositions[item.id]
      initialNodes.push({
        id: item.id,
        type: 'flowNode',
        parentId: layer.id,
        position: { x: ip.x - lp.x, y: ip.y - lp.y },
        data: { label: item.label },
      })
    })
  })

  const [nodes, , onNodesChange] = useNodesState(initialNodes)

  const styledEdges = useMemo(() => {
    const resultEdges: Edge[] = []
    const animatedEdges = new Set(['dag-sched', 'sched-gen', 'gen-grade', 'grade-score', 'score-sqlite', 'score-pg'])
    archEdges.forEach((e) => {
      const key = `${e.source}-${e.target}`
      resultEdges.push({
        id: `e-${key}`,
        source: e.source,
        target: e.target,
        animated: animatedEdges.has(key),
        style: { stroke: animatedEdges.has(key) ? c.accentBlue : c.borderStrong, strokeWidth: 1.5 },
      })
    })
    return resultEdges
  }, [c])

  const [edges, setEdges, onEdgesChange] = useEdgesState(styledEdges)

  useEffect(() => {
    setEdges(styledEdges)
  }, [styledEdges, setEdges])

  if (!mounted) return <div style={{ height: 440, width: '100%' }} />

  return (
    <div style={{ height: 440, width: '100%', border: '0.5px solid ' + c.border, background: 'transparent', position: 'relative' }}>
      <ReactFlow
        nodes={nodes}
        edges={edges}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        nodeTypes={nodeTypes}
        fitView
        fitViewOptions={{ padding: 0.3 }}
        nodesConnectable={false}
        preventScrolling={false}
        proOptions={{ hideAttribution: true }}
      >
        <Background color={c.border} gap={24} size={0.5} />
        <Controls />
      </ReactFlow>
    </div>
  )
}
