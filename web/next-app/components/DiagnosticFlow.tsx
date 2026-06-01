'use client'

import { useMemo, useEffect } from 'react'
import {
  ReactFlow,
  Background,
  Controls,
  Handle,
  Position,
  MarkerType,
  useNodesState,
  useEdgesState,
  type Node,
  type Edge,
} from '@xyflow/react'
import '@xyflow/react/dist/base.css'
import { useTheme } from '../hooks/useTheme'
import { themeColors } from './FlowDiagram'

const monoFont = "'IBM Plex Mono', monospace"

const NODE_W = 150
const CX = 200

function BaseNode({ data, variant }: { data: { label: string }; variant: 'start' | 'decision' | 'process' | 'terminal' }) {
  const { theme } = useTheme()
  const isDark = theme === 'dark'
  const c = themeColors[isDark ? 'dark' : 'light']

  if (variant === 'start') {
    return (
      <div style={{
        width: 28, height: 28, borderRadius: '50%', position: 'relative',
        background: isDark ? c.accentBlue : c.accentBlue,
        border: 'none',
        display: 'flex', alignItems: 'center', justifyContent: 'center',
        boxShadow: '0 2px 6px rgba(0,0,0,0.15)',
      }}>
        <Handle type="source" position={Position.Bottom} style={{ opacity: 0, pointerEvents: 'none' }} />
        <span style={{ color: '#fff', fontSize: '14px', lineHeight: 1 }}>●</span>
      </div>
    )
  }

  const style: React.CSSProperties = {
    fontFamily: monoFont,
    fontSize: '11px',
    padding: '10px 14px',
    textAlign: 'center',
    lineHeight: 1.5,
    borderRadius: 0,
    minWidth: 120,
    maxWidth: 180,
    boxShadow: '0 1px 3px rgba(0,0,0,0.06)',
  }

  if (variant === 'decision') {
    Object.assign(style, {
      background: isDark ? '#1e1e24' : '#fafafa',
      border: '0.5px solid ' + (isDark ? '#3f3f46' : '#d4d4d8'),
      borderLeft: '2px solid ' + c.accentBlue,
      color: c.textPrimary,
      fontWeight: 500,
    })
  } else if (variant === 'terminal') {
    Object.assign(style, {
      background: isDark ? '#1a1a2e' : '#f0f4ff',
      border: '0.5px solid ' + c.accentBlue,
      color: c.accentBlue,
      fontWeight: 500,
      fontSize: '10px',
      letterSpacing: '0.05em',
    })
  } else {
    Object.assign(style, {
      background: c.surface,
      border: '0.5px solid ' + c.border,
      color: c.textSecondary,
    })
  }

  return (
    <div style={style}>
      <Handle type="target" position={Position.Top} style={{ opacity: 0, pointerEvents: 'none' }} />
      <span style={{ whiteSpace: 'pre-wrap' }}>{data.label}</span>
      {variant !== 'terminal' && (
        <Handle type="source" position={Position.Bottom} style={{ opacity: 0, pointerEvents: 'none' }} />
      )}
    </div>
  )
}

const nodeTypes = {
  startNode: ({ data, ...props }: { data: { label: string } }) => <BaseNode data={data} {...props} variant="start" />,
  decisionNode: ({ data, ...props }: { data: { label: string } }) => <BaseNode data={data} {...props} variant="decision" />,
  processNode: ({ data, ...props }: { data: { label: string } }) => <BaseNode data={data} {...props} variant="process" />,
  terminalNode: ({ data, ...props }: { data: { label: string } }) => <BaseNode data={data} {...props} variant="terminal" />,
}

export default function DiagnosticFlow() {
  const { theme, mounted } = useTheme()
  const isDark = theme === 'dark'
  const c = themeColors[isDark ? 'dark' : 'light']

  const gap = 64
  const yStart = 4
  const yMid = yStart + 28 + 6
  const yBranch = yMid + 44 + gap
  const yCheck = yBranch + 44 + gap
  const yLocked = yCheck + 44 + gap
  const midX = CX - NODE_W / 2 + 40

  const initialNodes: Node[] = [
    { id: 'start', type: 'startNode', position: { x: CX - 14, y: yStart }, data: { label: '' } },
    { id: 'midpoint', type: 'decisionNode', position: { x: midX, y: yMid }, data: { label: 'Test midpoint\nconcept?' } },
    { id: 'forward', type: 'processNode', position: { x: CX - NODE_W - 30, y: yBranch }, data: { label: 'Move forward:\nharder concepts' } },
    { id: 'backward', type: 'processNode', position: { x: CX + 30, y: yBranch }, data: { label: 'Move backward:\nfoundational' } },
    { id: 'check', type: 'decisionNode', position: { x: CX - NODE_W / 2, y: yCheck }, data: { label: '3 correct\nin a row?' } },
    { id: 'locked', type: 'terminalNode', position: { x: CX - NODE_W / 2, y: yLocked }, data: { label: 'FRONTIER\nLOCKED' } },
  ]

  const [nodes, , onNodesChange] = useNodesState(initialNodes)

  const styledEdges = useMemo(() => {
    const edgeStyle = (color: string) => ({ stroke: color, strokeWidth: 1.5 })
    const labelStyle = (color: string) => ({ fontFamily: monoFont, fontSize: '10px', fill: color, fontWeight: 500 })

    return [
      { id: 'e-start-mid', source: 'start', target: 'midpoint', style: edgeStyle(c.borderStrong) },
      {
        id: 'e-mid-forward', source: 'midpoint', target: 'forward',
        label: 'Correct, fast',
        labelStyle: labelStyle(c.accentGreen),
        style: edgeStyle(c.accentGreen),
        markerEnd: { type: MarkerType.ArrowClosed, color: c.accentGreen },
      },
      {
        id: 'e-mid-backward', source: 'midpoint', target: 'backward',
        label: 'Incorrect, slow',
        labelStyle: labelStyle(c.accentRed),
        style: edgeStyle(c.accentRed),
        markerEnd: { type: MarkerType.ArrowClosed, color: c.accentRed },
      },
      {
        id: 'e-forward-check', source: 'forward', target: 'check',
        style: edgeStyle(c.borderStrong),
        type: 'smoothstep',
      },
      {
        id: 'e-backward-check', source: 'backward', target: 'check',
        style: edgeStyle(c.borderStrong),
        type: 'smoothstep',
      },
      {
        id: 'e-check-locked', source: 'check', target: 'locked',
        label: 'Yes', animated: true,
        labelStyle: labelStyle(c.accentGreen),
        style: edgeStyle(c.accentGreen),
      },
      {
        id: 'e-check-mid', source: 'check', target: 'midpoint',
        label: 'No',
        labelStyle: labelStyle(c.accentRed),
        style: { ...edgeStyle(c.accentRed), strokeDasharray: '5 5' },
        type: 'smoothstep',
      },
    ]
  }, [theme])

  const [edges, setEdges, onEdgesChange] = useEdgesState(styledEdges)

  useEffect(() => {
    setEdges(styledEdges)
  }, [styledEdges, setEdges])

  if (!mounted) return <div style={{ height: 440, width: '100%' }} />

  return (
    <div style={{
      height: 440, width: '100%',
      border: '0.5px solid ' + c.border,
      background: isDark ? 'rgba(0,0,0,0.12)' : 'rgba(0,0,0,0.02)',
    }}>
      <ReactFlow
        nodes={nodes}
        edges={edges}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        nodeTypes={nodeTypes}
        fitView
        fitViewOptions={{ padding: 0.25 }}
        nodesConnectable={false}
        panOnDrag={true}
        zoomOnScroll={true}
        zoomOnDoubleClick={true}
        preventScrolling={false}
        proOptions={{ hideAttribution: true }}
      >
        <Background color={c.border} gap={20} size={0.5} />
        <Controls showInteractive={false} />
      </ReactFlow>
    </div>
  )
}
