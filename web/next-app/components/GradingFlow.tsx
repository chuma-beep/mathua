'use client'

import FlowDiagram, { type FlowNodeDef, type FlowEdgeDef } from './FlowDiagram'

export default function GradingFlow() {
  const nodes: FlowNodeDef[] = [
    { id: 'input', label: 'Student\nAnswer' },
    { id: 'router', label: 'Router', variant: 'decision' },
    { id: 'numeric', label: 'Numeric\nGo' },
    { id: 'sympy', label: 'SymPy\nPython' },
    { id: 'choice', label: 'Choice\nGo' },
    { id: 'comp', label: 'Comparison\nGo' },
    { id: 'order', label: 'Ordering\nGo' },
    { id: 'fallback', label: 'Symbolic Go\nfallback', variant: 'optional' },
  ]

  const edges: FlowEdgeDef[] = [
    { id: 'e-input-router', source: 'input', target: 'router', animated: true },
    { id: 'e-router-numeric', source: 'router', target: 'numeric', label: 'numeric', color: 'blue', animated: true },
    { id: 'e-router-sympy', source: 'router', target: 'sympy', label: 'polynomial /\nexpression', color: 'blue', animated: true },
    { id: 'e-router-choice', source: 'router', target: 'choice', label: 'multiple_choice', color: 'blue', animated: true },
    { id: 'e-router-comp', source: 'router', target: 'comp', label: 'comparison', color: 'blue', animated: true },
    { id: 'e-router-order', source: 'router', target: 'order', label: 'ordering', color: 'blue', animated: true },
    { id: 'e-sympy-fallback', source: 'sympy', target: 'fallback', label: 'no Python', dashed: true, color: 'red' },
  ]

  return <FlowDiagram nodes={nodes} edges={edges} direction="LR" height={300} />
}
