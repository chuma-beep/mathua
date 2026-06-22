'use client'

import { useMemo, useEffect, useCallback, useState } from 'react'
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
import { useTheme } from '../../hooks/useTheme'
import { themeColors } from '../FlowDiagram'

const monoFont = "'IBM Plex Mono', monospace"

const detailInfo: Record<string, string> = {
  'web-app': 'Next.js 14 static export with React + KaTeX. Served directly by the Go binary. Dark/light theme via CSS custom properties.',
  'api-router': 'net/http standard library. 25+ REST endpoints. Middleware stack: CORS, logging, JWT auth, rate limiting.',
  'auth': 'JWT-based (HS256, 30-day expiry). bcrypt password hashing. Rate-limited by token bucket (5/min).',
  'session': 'Practice session management. Creates sessions, routes answers through the engine.',
  'dag-loader': 'Loads data/concepts/*.json. Validates no cycles, no orphaned prereqs. Topological sort via Kahn.',
  'sm2-scheduler': 'Priority-based concept selection. Recency (0.5), mastery (0.2), weakness (0.3). 70/30 new/review balance.',
  'generators': '18 domain-specific generator registrations. All problems procedurally generated -- no static bank.',
  'mastery': 'State machine: UNSEEN -> LEARNING -> PRACTICING -> MASTERED -> (DECAYING). Streak + time threshold.',
  'scoring': 'XP = base(10/5) * timeMult(0.5-1.5) * streakMult(1 + min(s,10)*0.1). Levels every 32 mastered.',
  'cat-diagnostic': 'Binary search on topological concept order. ~20-35 questions vs 284. Optimistic merge on retake.',
  'leaderboard': 'Weekly scores. Score = weekly_mastered * 100. Ranked descending. Resets Monday.',
  'numeric-grader': 'Go native. Handles integers, decimals, fractions (big.Rat), mixed numbers, scientific notation.',
  'sympy-grader': 'Long-lived Python subprocess. JSON over stdin/stdout. 10s timeout. Fallback to Go symbolic grader.',
  'multiple-choice': 'Case-insensitive. Single-letter matching (B matches Option B).',
  'sqlite': 'Local/dev. WAL mode, foreign_keys ON, busy_timeout 5000. File-based, zero config.',
  'postgresql': 'Production web. Shared state, concurrent connections. Via Repository interface.',
}

function LayerNode({ data }: { data: { label: string } }) {
  const { theme } = useTheme()
  const c = themeColors[theme === 'dark' ? 'dark' : 'light']
  return (
    <div style={{
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
    }}>
      <span style={{ marginBottom: 6 }}>{data.label}</span>
    </div>
  )
}

function FlowNode({ data }: { data: { label: string } }) {
  const { theme } = useTheme()
  const c = themeColors[theme === 'dark' ? 'dark' : 'light']
  return (
    <div style={{
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
    }}>
      <Handle type="target" position={Position.Top} style={{ opacity: 0, pointerEvents: 'none' }} />
      {data.label}
      <Handle type="source" position={Position.Bottom} style={{ opacity: 0, pointerEvents: 'none' }} />
    </div>
  )
}

const nodeTypes = { layerNode: LayerNode, flowNode: FlowNode }

interface LayerItem { id: string; label: string }
interface LayerDef { id: string; label: string; items: LayerItem[] }
interface ArchEdge { source: string; target: string }

const layers: LayerDef[] = [
  { id: 'ui', label: 'UI Layer', items: [{ id: 'web-app', label: 'Web App\nReact + KaTeX' }] },
  { id: 'api', label: 'API Layer', items: [
    { id: 'api-router', label: 'API Router' },
    { id: 'auth', label: 'Auth' },
    { id: 'session', label: 'Session' },
  ]},
  { id: 'engine', label: 'Core Engine', items: [
    { id: 'dag-loader', label: 'DAG Loader' },
    { id: 'sm2-scheduler', label: 'SM-2 Scheduler' },
    { id: 'generators', label: 'Generators' },
    { id: 'mastery', label: 'Mastery' },
    { id: 'scoring', label: 'Scoring' },
    { id: 'cat-diagnostic', label: 'CAT Diagnostic' },
    { id: 'leaderboard', label: 'Leaderboard' },
  ]},
  { id: 'grading', label: 'Grading', items: [
    { id: 'numeric-grader', label: 'Numeric' },
    { id: 'sympy-grader', label: 'SymPy' },
    { id: 'multiple-choice', label: 'Multiple Choice' },
  ]},
  { id: 'storage', label: 'Storage', items: [
    { id: 'sqlite', label: 'SQLite' },
    { id: 'postgresql', label: 'PostgreSQL' },
  ]},
]

const archEdges: ArchEdge[] = [
  { source: 'web-app', target: 'api-router' },
  { source: 'api-router', target: 'session' },
  { source: 'session', target: 'dag-loader' },
  { source: 'session', target: 'sm2-scheduler' },
  { source: 'dag-loader', target: 'sm2-scheduler' },
  { source: 'sm2-scheduler', target: 'generators' },
  { source: 'generators', target: 'mastery' },
  { source: 'mastery', target: 'scoring' },
  { source: 'generators', target: 'numeric-grader' },
  { source: 'generators', target: 'sympy-grader' },
  { source: 'scoring', target: 'sqlite' },
  { source: 'mastery', target: 'sqlite' },
]

const ITEM_W = 110
const ITEM_H = 40
const LAYER_PAD = 8
const LAYER_GAP_Y = 12
const ITEM_GAP_X = 8
const ITEM_GAP_Y = 6

function buildLayout() {
  let totalH = 0
  const layerPositions: Record<string, { x: number; y: number; w: number; h: number }> = {}
  const itemPositions: Record<string, { x: number; y: number }> = {}

  layers.forEach((layer) => {
    const perRow = layer.items.length <= 4 ? layer.items.length : Math.ceil(layer.items.length / 2)
    const rows = Math.ceil(layer.items.length / perRow)
    const w = perRow * ITEM_W + (perRow - 1) * ITEM_GAP_X + LAYER_PAD * 2
    const h = rows * ITEM_H + (rows - 1) * ITEM_GAP_Y + LAYER_PAD * 2 + 16

    layerPositions[layer.id] = { x: 0, y: totalH, w, h }

    layer.items.forEach((item, idx) => {
      const row = Math.floor(idx / perRow)
      const col = idx % perRow
      itemPositions[item.id] = {
        x: LAYER_PAD + col * (ITEM_W + ITEM_GAP_X),
        y: totalH + LAYER_PAD + 16 + row * (ITEM_H + ITEM_GAP_Y),
      }
    })

    totalH += h + LAYER_GAP_Y
  })

  const totalW = Math.max(...Object.values(layerPositions).map(l => l.w))
  return { layerPositions, itemPositions, totalW, totalH }
}

function DetailPanel({ nodeId, onClose }: { nodeId: string | null; onClose: () => void }) {
  const { theme } = useTheme()
  const c = themeColors[theme === 'dark' ? 'dark' : 'light']
  if (!nodeId) return null
  const info = detailInfo[nodeId]
  if (!info) return null

  return (
    <div style={{
      position: 'absolute', top: 10, right: 10, width: 280, zIndex: 10,
      background: c.surfaceElevated,
      border: '0.5px solid ' + c.border,
      borderRadius: 0,
      padding: '12px 16px',
      fontFamily: "'IBM Plex Serif', serif",
      fontSize: '13px',
      color: c.textSecondary,
      lineHeight: 1.6,
    }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 8 }}>
        <strong style={{ fontFamily: monoFont, color: c.textPrimary, fontSize: 12 }}>{nodeId}</strong>
        <button onClick={onClose} style={{ background: 'none', border: 'none', color: c.textMuted, cursor: 'pointer', fontSize: 14, padding: '2px 6px' }}>x</button>
      </div>
      <p style={{ margin: 0 }}>{info}</p>
    </div>
  )
}

export default function ArchitectureFlow() {
  const { theme, mounted } = useTheme()
  const c = themeColors[theme === 'dark' ? 'dark' : 'light']
  const layout = buildLayout()

  const initialNodes: Node[] = []
  layers.forEach((layer) => {
    const lp = layout.layerPositions[layer.id]
    initialNodes.push({ id: layer.id, type: 'layerNode', position: { x: lp.x, y: lp.y }, data: { label: layer.label }, style: { width: lp.w, height: lp.h } })
    layer.items.forEach((item) => {
      const ip = layout.itemPositions[item.id]
      initialNodes.push({ id: item.id, type: 'flowNode', parentId: layer.id, position: { x: ip.x - lp.x, y: ip.y - lp.y }, data: { label: item.label } })
    })
  })

  const [nodes, , onNodesChange] = useNodesState(initialNodes)

  const animatedEdges = new Set(['dag-loader-sm2-scheduler', 'sm2-scheduler-generators', 'generators-mastery', 'mastery-scoring', 'scoring-sqlite', 'mastery-sqlite'])

  const styledEdges = useMemo(() => {
    return archEdges.map((e) => {
      const key = `${e.source}-${e.target}`
      return {
        id: `e-${key}`,
        source: e.source,
        target: e.target,
        animated: animatedEdges.has(key),
        style: { stroke: animatedEdges.has(key) ? c.accentBlue : c.borderStrong, strokeWidth: 1.5 },
      }
    })
  }, [theme])

  const [edges, setEdges, onEdgesChange] = useEdgesState(styledEdges)
  const [selectedNode, setSelectedNode] = useState<string | null>(null)

  useEffect(() => { setEdges(styledEdges) }, [styledEdges, setEdges])

  const onNodeClick = useCallback((_: React.MouseEvent, node: Node) => {
    setSelectedNode((prev) => (prev === node.id ? null : node.id))
  }, [])

  const onPaneClick = useCallback(() => { setSelectedNode(null) }, [])

  if (!mounted) return <div style={{ height: 520, width: '100%' }} />

  return (
    <div style={{ height: 520, width: '100%', border: '0.5px solid ' + c.border, background: 'transparent', position: 'relative' }}>
      <ReactFlow
        nodes={nodes}
        edges={edges}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        nodeTypes={nodeTypes}
        onNodeClick={onNodeClick}
        onPaneClick={onPaneClick}
        fitView
        fitViewOptions={{ padding: 0.3 }}
        nodesConnectable={false}
        preventScrolling={false}
        proOptions={{ hideAttribution: true }}
      >
        <Background color={c.border} gap={24} size={0.5} />
        <Controls />
      </ReactFlow>
      <DetailPanel nodeId={selectedNode} onClose={() => setSelectedNode(null)} />
    </div>
  )
}
