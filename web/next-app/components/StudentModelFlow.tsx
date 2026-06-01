'use client'

import FlowDiagram, { type FlowNodeDef, type FlowEdgeDef } from './FlowDiagram'

export default function StudentModelFlow() {
  const nodes: FlowNodeDef[] = [
    { id: 'unseen', label: 'UNSEEN' },
    { id: 'learning', label: 'LEARNING' },
    { id: 'practicing', label: 'PRACTICING' },
    { id: 'mastered', label: 'MASTERED', variant: 'terminal' },
    { id: 'decaying', label: 'DECAYING', variant: 'error' },
  ]

  const edges: FlowEdgeDef[] = [
    { id: 'e-unseen-learning', source: 'unseen', target: 'learning', animated: true },
    { id: 'e-learning-practicing', source: 'learning', target: 'practicing', label: 'first\ncorrect', color: 'green', animated: true },
    { id: 'e-practicing-mastered', source: 'practicing', target: 'mastered', label: 'streak\nreached', color: 'green', animated: true },
    { id: 'e-mastered-decaying', source: 'mastered', target: 'decaying', label: 'interval\nelapsed', color: 'blue', animated: true },
    { id: 'e-decaying-practicing', source: 'decaying', target: 'practicing', label: 'review\ncorrect', dashed: true, color: 'green' },
  ]

  return <FlowDiagram nodes={nodes} edges={edges} direction="TB" height={320} rankSep={150} nodeSep={80} allowZoom />
}
