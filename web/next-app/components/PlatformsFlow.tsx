'use client'

import { useMemo } from 'react'
import {
  ReactFlow,
  Background,
  Handle,
  Position,
  type Node,
  type Edge,
} from '@xyflow/react'
import '@xyflow/react/dist/base.css'
import { useTheme } from '../hooks/useTheme'
import { themeColors } from './FlowDiagram'

const monoFont = "'IBM Plex Mono', monospace"

const GROUP_W = 280
const GROUP_HEADER_H = 26
const CHILD_H = 36
const CHILD_GAP = 6
const GROUP_PAD_BOTTOM = 12

const GROUP_H = GROUP_HEADER_H + 4 * CHILD_H + 3 * CHILD_GAP + GROUP_PAD_BOTTOM
const GROUP_GAP = 40
const LEFT_X = 4
const RIGHT_X = LEFT_X + GROUP_W + GROUP_GAP
const ENGINE_W = 170
const ENGINE_H = 44
const ENGINE_Y = GROUP_H + 32

function GroupNode({ data }: { data: { label: string } }) {
  const { theme } = useTheme()
  const isDark = theme === 'dark'
  const c = themeColors[isDark ? 'dark' : 'light']

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
      }}
    >
      <span>{data.label}</span>
      <Handle type="source" position={Position.Bottom} style={{ opacity: 0, pointerEvents: 'none' }} />
    </div>
  )
}

function DetailNode({ data }: { data: { label: string } }) {
  const { theme } = useTheme()
  const isDark = theme === 'dark'
  const c = themeColors[isDark ? 'dark' : 'light']

  return (
    <div
      style={{
        background: c.surface,
        border: '0.5px solid ' + c.border,
        fontFamily: monoFont,
        fontSize: '10px',
        color: c.textSecondary,
        padding: '8px 10px',
        lineHeight: 1.3,
        borderRadius: 0,
        height: '100%',
        display: 'flex',
        alignItems: 'center',
      }}
    >
      <Handle type="target" position={Position.Left} style={{ opacity: 0, pointerEvents: 'none' }} />
      {data.label}
      <Handle type="source" position={Position.Right} style={{ opacity: 0, pointerEvents: 'none' }} />
    </div>
  )
}

function EngineNode({ data }: { data: { label: string } }) {
  const { theme } = useTheme()
  const isDark = theme === 'dark'
  const c = themeColors[isDark ? 'dark' : 'light']

  return (
    <div
      style={{
        background: c.surfaceElevated,
        border: '0.5px solid ' + c.accentBlue,
        fontFamily: monoFont,
        fontSize: '11px',
        color: c.accentBlue,
        fontWeight: 500,
        padding: '12px 24px',
        textAlign: 'center',
        borderRadius: 0,
        width: '100%',
        height: '100%',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        letterSpacing: '0.06em',
        textTransform: 'uppercase',
      }}
    >
      <Handle type="target" position={Position.Top} style={{ opacity: 0, pointerEvents: 'none' }} />
      {data.label}
      <Handle type="source" position={Position.Bottom} style={{ opacity: 0, pointerEvents: 'none' }} />
    </div>
  )
}

const nodeTypes = {
  groupNode: GroupNode,
  detailNode: DetailNode,
  engineNode: EngineNode,
}

const webItems = [
  { id: 'w-react', label: 'React + KaTeX' },
  { id: 'w-account', label: 'Account required' },
  { id: 'w-lboard', label: 'Global leaderboard' },
  { id: 'w-graph', label: 'Graph view' },
]

const desktopItems = [
  { id: 'd-tui', label: 'Bubble Tea TUI' },
  { id: 'd-noacct', label: 'No account needed' },
  { id: 'd-offline', label: 'Fully offline' },
  { id: 'd-sqlite', label: 'SQLite storage' },
]

export default function PlatformsFlow() {
  const { theme, mounted } = useTheme()
  const isDark = theme === 'dark'
  const c = themeColors[isDark ? 'dark' : 'light']

  const { nodes, edges } = useMemo(() => {
    const ws: Node[] = []

    ws.push({
      id: 'web', type: 'groupNode', position: { x: LEFT_X, y: 0 },
      data: { label: 'Web browser' }, style: { width: GROUP_W, height: GROUP_H },
      draggable: false,
    })

    webItems.forEach((item, i) => {
      ws.push({
        id: item.id, type: 'detailNode', parentId: 'web',
        position: { x: 10, y: GROUP_HEADER_H + i * (CHILD_H + CHILD_GAP) },
        data: { label: item.label },
        style: { width: GROUP_W - 20, height: CHILD_H },
        draggable: false,
      })
    })

    ws.push({
      id: 'desktop', type: 'groupNode', position: { x: RIGHT_X, y: 0 },
      data: { label: 'Desktop terminal' }, style: { width: GROUP_W, height: GROUP_H },
      draggable: false,
    })

    desktopItems.forEach((item, i) => {
      ws.push({
        id: item.id, type: 'detailNode', parentId: 'desktop',
        position: { x: 10, y: GROUP_HEADER_H + i * (CHILD_H + CHILD_GAP) },
        data: { label: item.label },
        style: { width: GROUP_W - 20, height: CHILD_H },
        draggable: false,
      })
    })

    const engineX = (LEFT_X + GROUP_W + RIGHT_X + GROUP_W) / 2 - ENGINE_W / 2
    ws.push({
      id: 'engine', type: 'engineNode', position: { x: engineX, y: ENGINE_Y },
      data: { label: 'Engine' },
      style: { width: ENGINE_W, height: ENGINE_H },
      draggable: false,
    })

    const edgeLabelStyle: React.CSSProperties = {
      fontFamily: monoFont, fontSize: '10px',
      fontWeight: 500,
    }

    const es: Edge[] = [
      {
        id: 'e-web-engine', source: 'web', target: 'engine',
        label: 'REST API',
        style: { stroke: c.accentBlue, strokeWidth: 1.5 },
        labelStyle: { ...edgeLabelStyle, fill: c.accentBlue },
        markerEnd: { type: 'arrowclosed' as any, color: c.accentBlue, width: 16, height: 16 },
      },
      {
        id: 'e-desktop-engine', source: 'desktop', target: 'engine',
        label: 'direct calls',
        style: { stroke: c.accentBlue, strokeWidth: 1.5 },
        labelStyle: { ...edgeLabelStyle, fill: c.accentBlue },
        markerEnd: { type: 'arrowclosed' as any, color: c.accentBlue, width: 16, height: 16 },
      },
    ]

    return { nodes: ws, edges: es }
  }, [theme])

  const totalW = LEFT_X + GROUP_W + GROUP_GAP + GROUP_W + 4
  const totalH = ENGINE_Y + ENGINE_H + 24

  if (!mounted) return <div style={{ height: totalH, width: '100%' }} />

  return (
    <div
      style={{
        height: totalH, width: '100%', maxWidth: totalW,
        border: '0.5px solid ' + c.border,
        background: 'transparent',
        margin: '0 auto',
        overflow: 'hidden',
      }}
    >
      <ReactFlow
        nodes={nodes}
        edges={edges}
        fitView
        proOptions={{ hideAttribution: true }}
      >
        <Background color={c.border} gap={24} size={0.5} />
      </ReactFlow>
    </div>
  )
}
