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

const stepData: Record<string, { label: string; desc: string }> = {
  '1-client': { label: 'Client', desc: 'POST /api/answer\n{session_id, answer, elapsed}' },
  '2-auth': { label: 'Auth Middleware', desc: 'Validate JWT\nExtract studentID\nInject into context' },
  '3-handler': { label: 'Handler', desc: 'handleAnswer()\nCalls Engine.SubmitAnswer()' },
  '4-grade': { label: 'Grader Router', desc: 'Dispatch by grading_type\nNumeric / SymPy / Choice / etc' },
  '5-mastery': { label: 'Mastery Machine', desc: 'State transition check\nUNSEEN -> LEARNING -> ...' },
  '6-sm2': { label: 'SM-2 Compute', desc: 'Compute quality (0-5)\nUpdate repetitions, interval, efactor' },
  '7-upsert': { label: 'Upsert Progress', desc: 'ON CONFLICT DO UPDATE\nSave all SM-2 fields' },
  '8-weakness': { label: 'Propagate Weakness', desc: 'If weakness > 0.3\nPropagate w*0.3 to dependents' },
  '9-attempt': { label: 'Record Attempt', desc: 'INSERT INTO attempts\n{session, concept, answer, correct, time}' },
  '10-xp': { label: 'Compute XP', desc: 'base * timeMult * streakMult\nrepo.AddXP()' },
  '11-next': { label: 'Next Question', desc: 'sched.Next() -> generator\nNew problem for next concept' },
  '12-response': { label: 'Response', desc: 'JSON to client\n{result, next_question, done}' },
}

function FlowNode({ data }: { data: { label: string; stepId: string } }) {
  const { theme } = useTheme()
  const c = themeColors[theme === 'dark' ? 'dark' : 'light']
  return (
    <div style={{
      background: c.surface,
      border: '0.5px solid ' + c.border,
      borderLeft: '2px solid ' + c.accentBlue,
      fontFamily: monoFont,
      fontSize: '11px',
      color: c.textSecondary,
      padding: '8px 12px',
      textAlign: 'center',
      lineHeight: 1.3,
      borderRadius: 0,
      whiteSpace: 'pre-wrap',
      minWidth: 160,
    }}>
      <Handle type="target" position={Position.Left} style={{ opacity: 0, pointerEvents: 'none' }} />
      <span style={{ color: c.accentBlue, fontWeight: 500 }}>{data.label}</span>
      <Handle type="source" position={Position.Right} style={{ opacity: 0, pointerEvents: 'none' }} />
    </div>
  )
}

const nodeTypes = { flowNode: FlowNode }

function StepDetail({ stepId, onClose }: { stepId: string | null; onClose: () => void }) {
  const { theme } = useTheme()
  const c = themeColors[theme === 'dark' ? 'dark' : 'light']
  if (!stepId) return null
  const step = stepData[stepId]
  if (!step) return null

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
        <strong style={{ fontFamily: monoFont, color: c.accentBlue, fontSize: 12 }}>{step.label}</strong>
        <button type="button" onClick={onClose} aria-label="Close panel" style={{ background: 'none', border: 'none', color: c.textMuted, cursor: 'pointer', fontSize: 14, padding: '2px 6px' }}>x</button>
      </div>
      <pre style={{
        fontFamily: monoFont,
        fontSize: 11,
        color: c.textMuted,
        background: c.surface,
        padding: 8,
        borderRadius: 0,
        whiteSpace: 'pre-wrap',
        margin: 0,
        lineHeight: 1.5,
      }}>{step.desc}</pre>
    </div>
  )
}

const steps = Object.keys(stepData)

export default function RequestFlow() {
  const { theme, mounted } = useTheme()
  const c = themeColors[theme === 'dark' ? 'dark' : 'light']

  const initialNodes: Node[] = steps.map((id, i) => ({
    id,
    type: 'flowNode',
    position: { x: 50, y: i * 52 },
    data: { label: stepData[id].label, stepId: id },
    style: { width: 150 },
  }))

  const edgeDefs: Edge[] = useMemo(() => steps.slice(0, -1).map((id, i) => ({
    id: `e-${id}-${steps[i + 1]}`,
    source: id,
    target: steps[i + 1],
    animated: true,
  })), [])

  const styledEdges = useMemo(() => {
    return edgeDefs.map(e => ({
      ...e,
      style: { stroke: c.borderStrong, strokeWidth: 1.5 },
    }))
  }, [c, edgeDefs])

  const [nodes, , onNodesChange] = useNodesState(initialNodes)
  const [edges, setEdges, onEdgesChange] = useEdgesState(styledEdges)
  const [selectedStep, setSelectedStep] = useState<string | null>(null)

  useEffect(() => { setEdges(styledEdges) }, [styledEdges, setEdges])

  const onNodeClick = useCallback((_: React.MouseEvent, node: Node) => {
    setSelectedStep((prev) => (prev === node.id ? null : node.id))
  }, [])

  const onPaneClick = useCallback(() => { setSelectedStep(null) }, [])

  if (!mounted) return <div style={{ height: 720, width: '100%' }} />

  return (
    <div style={{ height: 720, width: '100%', border: '0.5px solid ' + c.border, background: 'transparent', position: 'relative' }}>
      <div style={{
        position: 'absolute', top: 8, left: 8, zIndex: 10,
        fontFamily: monoFont, fontSize: 10, color: c.textMuted,
        background: c.surface, padding: '4px 8px', borderRadius: 0,
        border: '0.5px solid ' + c.border,
      }}>
        Click a step for details
      </div>
      <ReactFlow
        nodes={nodes}
        edges={edges}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        nodeTypes={nodeTypes}
        onNodeClick={onNodeClick}
        onPaneClick={onPaneClick}
        fitView
        fitViewOptions={{ padding: 0.15 }}
        nodesConnectable={false}
        preventScrolling={false}
        proOptions={{ hideAttribution: true }}
      >
        <Background color={c.border} gap={24} size={0.5} />
        <Controls />
      </ReactFlow>
      <StepDetail stepId={selectedStep} onClose={() => setSelectedStep(null)} />
    </div>
  )
}
