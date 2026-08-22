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

const tableDetails: Record<string, { columns: string[]; desc: string }> = {
  students: {
    desc: 'User accounts with XP tracking and auth. UUID primary key.',
    columns: ['id TEXT PK', 'name TEXT', 'username TEXT', 'password_hash TEXT', 'course_id TEXT', 'xp_total INTEGER', 'xp_today INTEGER', 'settings TEXT (JSON)'],
  },
  concept_progress: {
    desc: 'Per-concept student state with SM-2 spaced repetition.',
    columns: ['student_id TEXT PK FK', 'concept_id TEXT PK', 'status TEXT', 'streak INTEGER', 'sm2_repetitions INTEGER', 'sm2_interval INTEGER', 'sm2_efactor REAL', 'next_review_due TEXT', 'weakness_score REAL'],
  },
  sessions: {
    desc: 'Practice sessions. UUID primary key.',
    columns: ['id TEXT PK', 'student_id TEXT FK', 'started_at TEXT'],
  },
  attempts: {
    desc: 'Individual answer records. Auto-increment PK.',
    columns: ['id INTEGER PK AUTO', 'session_id TEXT FK', 'student_id TEXT FK', 'concept_id TEXT', 'answer TEXT', 'expected TEXT', 'correct INTEGER', 'elapsed_seconds REAL', 'timestamp TEXT'],
  },
  questions: {
    desc: 'Optional pre-seeded question bank.',
    columns: ['id INTEGER PK AUTO', 'concept_id TEXT', 'question TEXT', 'answer TEXT', 'explanation TEXT', 'difficulty REAL'],
  },
}

function TableNode({ data }: { data: { label: string } }) {
  const { theme } = useTheme()
  const c = themeColors[theme === 'dark' ? 'dark' : 'light']
  return (
    <div style={{
      background: c.surface,
      border: '0.5px solid ' + c.border,
      fontFamily: monoFont,
      fontSize: '11px',
      color: c.textPrimary,
      padding: '10px 16px',
      textAlign: 'center',
      lineHeight: 1.3,
      borderRadius: 0,
      fontWeight: 500,
    }}>
      <Handle type="target" position={Position.Left} style={{ opacity: 0, pointerEvents: 'none' }} />
      {data.label}
      <Handle type="source" position={Position.Right} style={{ opacity: 0, pointerEvents: 'none' }} />
    </div>
  )
}

const nodeTypes = { tableNode: TableNode }

const initialNodes: Node[] = [
  { id: 'students', type: 'tableNode', position: { x: 0, y: 0 }, data: { label: 'students' }, style: { width: 150 } },
  { id: 'concept_progress', type: 'tableNode', position: { x: 300, y: 0 }, data: { label: 'concept_progress' }, style: { width: 150 } },
  { id: 'sessions', type: 'tableNode', position: { x: 0, y: 100 }, data: { label: 'sessions' }, style: { width: 150 } },
  { id: 'attempts', type: 'tableNode', position: { x: 300, y: 100 }, data: { label: 'attempts' }, style: { width: 150 } },
  { id: 'questions', type: 'tableNode', position: { x: 550, y: 0 }, data: { label: 'questions' }, style: { width: 150 } },
]

const initialEdges: Edge[] = [
  { id: 'e-students-progress', source: 'students', target: 'concept_progress', animated: true, style: { strokeWidth: 1.5 } },
  { id: 'e-students-sessions', source: 'students', target: 'sessions', animated: true, style: { strokeWidth: 1.5 } },
  { id: 'e-students-attempts', source: 'students', target: 'attempts', animated: true, style: { strokeWidth: 1.5 } },
  { id: 'e-sessions-attempts', source: 'sessions', target: 'attempts', animated: true, style: { strokeWidth: 1.5 } },
]

function TableDetail({ tableId, onClose }: { tableId: string | null; onClose: () => void }) {
  const { theme } = useTheme()
  const c = themeColors[theme === 'dark' ? 'dark' : 'light']
  if (!tableId) return null
  const detail = tableDetails[tableId]
  if (!detail) return null

  return (
    <div style={{
      position: 'absolute', top: 10, right: 10, width: 320, zIndex: 10,
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
        <strong style={{ fontFamily: monoFont, color: c.textPrimary, fontSize: 13 }}>{tableId}</strong>
        <button onClick={onClose} style={{ background: 'none', border: 'none', color: c.textMuted, cursor: 'pointer', fontSize: 14, padding: '2px 6px' }}>x</button>
      </div>
      <p style={{ margin: '0 0 10px', fontSize: 12, fontStyle: 'italic' }}>{detail.desc}</p>
      <table style={{ width: '100%', borderCollapse: 'collapse', fontFamily: monoFont, fontSize: 11 }}>
        <tbody>
          {detail.columns.map((col) => (
            <tr key={col} style={{ borderBottom: '0.5px solid ' + c.border }}>
              <td style={{ padding: '3px 4px', color: c.textPrimary }}>{col}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

export default function DataModelFlow() {
  const { theme, mounted } = useTheme()
  const c = themeColors[theme === 'dark' ? 'dark' : 'light']

  const styledEdges = useMemo(() => {
    return initialEdges.map(e => ({
      ...e,
      style: { ...e.style, stroke: c.borderStrong },
    }))
  }, [c])

  const [nodes, , onNodesChange] = useNodesState(initialNodes)
  const [edges, setEdges, onEdgesChange] = useEdgesState(styledEdges)
  const [selectedTable, setSelectedTable] = useState<string | null>(null)

  useEffect(() => { setEdges(styledEdges) }, [styledEdges, setEdges])

  const onNodeClick = useCallback((_: React.MouseEvent, node: Node) => {
    setSelectedTable((prev) => (prev === node.id ? null : node.id))
  }, [])

  const onPaneClick = useCallback(() => { setSelectedTable(null) }, [])

  if (!mounted) return <div style={{ height: 320, width: '100%' }} />

  return (
    <div style={{ height: 320, width: '100%', border: '0.5px solid ' + c.border, background: 'transparent', position: 'relative' }}>
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
      <TableDetail tableId={selectedTable} onClose={() => setSelectedTable(null)} />
    </div>
  )
}
