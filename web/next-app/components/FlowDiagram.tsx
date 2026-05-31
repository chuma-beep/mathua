'use client'

import { useMemo } from 'react'
import {
  ReactFlow,
  Background,
  Controls,
  Handle,
  Position,
  type Node,
  type Edge,
} from '@xyflow/react'
import '@xyflow/react/dist/base.css'
import dagre from 'dagre'
import { useTheme } from '../hooks/useTheme'

const monoFont = "'IBM Plex Mono', monospace"

export const themeColors = {
  light: {
    surface: '#f2f2f3',
    surfaceElevated: '#f2f2f3',
    textPrimary: 'oklch(0.16 0.004 264)',
    textSecondary: 'oklch(0.54 0 0)',
    textMuted: 'oklch(0.54 0 0)',
    accentBlue: 'oklch(0.6 0.16 255)',
    accentGreen: '#059669',
    accentRed: '#dc2626',
    border: 'oklch(0.85 0.02 264)',
    borderStrong: 'oklch(0.75 0.03 264)',
  },
  dark: {
    surface: 'oklch(0.21 0 0 / 0.8)',
    surfaceElevated: 'oklch(0.21 0 0)',
    textPrimary: 'oklch(1 0 0)',
    textSecondary: 'oklch(0.66 0 0)',
    textMuted: 'oklch(0.66 0 0)',
    accentBlue: 'oklch(0.72 0.145 230)',
    accentGreen: '#10b981',
    accentRed: '#ef4444',
    border: 'oklch(0.35 0.02 264)',
    borderStrong: 'oklch(0.45 0.03 264)',
  },
}

type NodeVariant = 'default' | 'optional' | 'decision' | 'terminal' | 'error' | 'start' | 'process'
type EdgeColor = 'default' | 'blue' | 'green' | 'red'

export interface FlowNodeDef {
  id: string
  label: string
  variant?: NodeVariant
}

export interface FlowEdgeDef {
  id: string
  source: string
  target: string
  label?: string
  animated?: boolean
  dashed?: boolean
  color?: EdgeColor
}

interface FlowDiagramProps {
  nodes: FlowNodeDef[]
  edges: FlowEdgeDef[]
  direction?: 'LR' | 'TB'
  height?: number
}

const NODE_WIDTH = 150
const NODE_HEIGHT = 44

function buildVariantStyles(c: typeof themeColors.light) {
  return {
    default: {
      background: c.surface,
      border: '0.5px solid ' + c.border,
      color: c.textPrimary,
    },
    optional: {
      background: c.surface,
      border: '1px dashed ' + c.borderStrong,
      color: c.textSecondary,
    },
    decision: {
      background: c.surface,
      border: '0.5px solid ' + c.border,
      borderLeft: '3px solid ' + c.accentBlue,
      color: c.textPrimary,
    },
    terminal: {
      background: c.surfaceElevated,
      border: '0.5px solid ' + c.accentBlue,
      color: c.accentBlue,
    },
    error: {
      background: c.surface,
      border: '0.5px solid ' + c.accentRed,
      borderRight: '3px solid ' + c.accentRed,
      color: c.accentRed,
    },
    start: {
      background: c.accentBlue,
      color: '#fff',
      border: 'none',
      borderRadius: '50%',
      width: 28,
      height: 28,
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center',
      fontSize: '16px',
      lineHeight: 1,
      fontWeight: 500,
    },
    process: {
      background: c.surface,
      border: '0.5px solid ' + c.border,
      color: c.textSecondary,
    },
  } as Record<string, React.CSSProperties>
}

function buildEdgeStyle(c: typeof themeColors.light, color?: EdgeColor, dashed?: boolean) {
  const stroke = color === 'green' ? c.accentGreen
    : color === 'red' ? c.accentRed
    : color === 'blue' ? c.accentBlue
    : c.borderStrong
  return {
    stroke,
    ...(dashed ? { strokeDasharray: '4 4' } : {}),
  }
}

function buildEdgeLabelStyle(c: typeof themeColors.light, color?: EdgeColor) {
  const fill = color === 'green' ? c.accentGreen
    : color === 'red' ? c.accentRed
    : color === 'blue' ? c.accentBlue
    : c.textMuted
  return { fontFamily: monoFont, fontSize: '10px', fill }
}

function FlowNode({ data }: { data: { label: string; variant?: string; dir?: string } }) {
  const { theme } = useTheme()
  const c = themeColors[theme === 'dark' ? 'dark' : 'light']
  const variantStyles = buildVariantStyles(c)
  const variant = data.variant ?? 'default'
  const style = variantStyles[variant]
  const isTB = data.dir === 'TB'
  const targetPos = isTB ? Position.Top : Position.Left
  const sourcePos = isTB ? Position.Bottom : Position.Right

  if (variant === 'start') {
    return (
      <div style={{ position: 'relative', ...style }}>
        <Handle type="source" position={sourcePos} style={{ opacity: 0, pointerEvents: 'none' }} />
        ●
      </div>
    )
  }

  return (
    <div
      style={{
        ...style,
        fontFamily: monoFont,
        fontSize: '12px',
        padding: '10px 16px',
        textAlign: 'center',
        lineHeight: 1.4,
        borderRadius: 0,
        minWidth: 120,
        maxWidth: 180,
      }}
    >
      <Handle type="target" position={targetPos} style={{ opacity: 0, pointerEvents: 'none' }} />
      <span style={{ whiteSpace: 'pre-wrap' }}>{data.label}</span>
      <Handle type="source" position={sourcePos} style={{ opacity: 0, pointerEvents: 'none' }} />
    </div>
  )
}

const nodeTypes = { flowNode: FlowNode }

export default function FlowDiagram({ nodes: nodeDefs, edges: edgeDefs, direction = 'LR', height = 400 }: FlowDiagramProps) {
  const { theme, mounted } = useTheme()
  const c = themeColors[theme === 'dark' ? 'dark' : 'light']

  const { nodes, edges } = useMemo(() => {
    const initialNodes: Node[] = nodeDefs.map(n => ({
      id: n.id,
      type: 'flowNode',
      position: { x: 0, y: 0 },
      data: { label: n.label, variant: n.variant ?? 'default', dir: direction },
    }))

    const initialEdges: Edge[] = edgeDefs.map(e => ({
      id: e.id,
      source: e.source,
      target: e.target,
      animated: e.animated,
      label: e.label,
      style: buildEdgeStyle(c, e.color, e.dashed),
      labelStyle: buildEdgeLabelStyle(c, e.color),
    }))

    const g = new dagre.graphlib.Graph()
    g.setGraph({ rankdir: direction, nodesep: 20, ranksep: 60, marginx: 20, marginy: 20 })
    g.setDefaultEdgeLabel(() => ({}))

    initialNodes.forEach(node => {
      const w = node.data.variant === 'start' ? 28 : NODE_WIDTH
      const h = node.data.variant === 'start' ? 28 : NODE_HEIGHT
      g.setNode(node.id, { width: w, height: h })
    })
    initialEdges.forEach(edge => g.setEdge(edge.source, edge.target))

    dagre.layout(g)

    const layoutedNodes = initialNodes.map(node => {
      const pos = g.node(node.id)
      const w = node.data.variant === 'start' ? 28 : NODE_WIDTH
      const h = node.data.variant === 'start' ? 28 : NODE_HEIGHT
      return { ...node, position: { x: pos.x - w / 2, y: pos.y - h / 2 } }
    })

    return { nodes: layoutedNodes, edges: initialEdges }
  }, [nodeDefs, edgeDefs, direction, c])

  if (!mounted) return <div style={{ height, width: '100%' }} />

  return (
    <div style={{ height, width: '100%', border: '0.5px solid ' + c.border, background: 'transparent' }}>
      <ReactFlow
        nodes={nodes}
        edges={edges}
        nodeTypes={nodeTypes}
        fitView
        fitViewOptions={{ padding: 0.3 }}
        panActivationKeyCode={null}
        nodesConnectable={false}
        panOnDrag={false}
        zoomOnScroll={false}
        zoomOnDoubleClick={false}
        preventScrolling={false}
        proOptions={{ hideAttribution: true }}
      >
        <Background color={c.border} gap={24} size={0.5} />
        <Controls showInteractive={false} />
      </ReactFlow>
    </div>
  )
}
