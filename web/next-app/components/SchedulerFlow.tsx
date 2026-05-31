'use client'

import FlowDiagram, { type FlowNodeDef, type FlowEdgeDef } from './FlowDiagram'

export default function SchedulerFlow() {
  const nodes: FlowNodeDef[] = [
    { id: 'submit', label: 'Question\nsubmitted' },
    { id: 'score', label: 'Compute priority\nfor each concept' },
    { id: 'enforce', label: 'Enforce\nhard rules', variant: 'decision' },
    { id: 'r1', label: 'Prerequisites\nmust be MASTERED' },
    { id: 'r2', label: 'No same concept\ntwice in a row' },
    { id: 'r3', label: '70% practice\n30% review' },
    { id: 'select', label: 'Select highest\npriority concept', variant: 'terminal' },
  ]

  const edges: FlowEdgeDef[] = [
    { id: 'e-submit-score', source: 'submit', target: 'score', animated: true },
    { id: 'e-score-enforce', source: 'score', target: 'enforce', animated: true },
    { id: 'e-enforce-r1', source: 'enforce', target: 'r1' },
    { id: 'e-enforce-r2', source: 'enforce', target: 'r2' },
    { id: 'e-enforce-r3', source: 'enforce', target: 'r3' },
    { id: 'e-r1-select', source: 'r1', target: 'select', animated: true },
    { id: 'e-r2-select', source: 'r2', target: 'select', animated: true },
    { id: 'e-r3-select', source: 'r3', target: 'select', animated: true },
  ]

  return <FlowDiagram nodes={nodes} edges={edges} direction="TB" height={320} />
}
