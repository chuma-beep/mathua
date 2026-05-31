'use client'

import FlowDiagram, { type FlowNodeDef, type FlowEdgeDef } from './FlowDiagram'

export default function PrWorkflow() {
  const nodes: FlowNodeDef[] = [
    { id: 'concept', label: 'Add concept JSON' },
    { id: 'gen', label: 'Write Go generator' },
    { id: 'lesson', label: 'Write lesson\nmarkdown', variant: 'optional' },
    { id: 'test', label: 'Write fuzz test\n1k samples' },
    { id: 'run', label: 'Run go test\n+ validate' },
    { id: 'pr', label: 'Open PR' },
    { id: 'ci', label: 'CI passes?', variant: 'decision' },
    { id: 'review', label: 'Maintainer\nreview' },
    { id: 'fix', label: 'Fix issues', variant: 'error' },
    { id: 'merge', label: 'Merged!', variant: 'terminal' },
  ]

  const edges: FlowEdgeDef[] = [
    { id: 'e-concept-gen', source: 'concept', target: 'gen' },
    { id: 'e-concept-lesson', source: 'concept', target: 'lesson', label: 'optional', dashed: true },
    { id: 'e-gen-test', source: 'gen', target: 'test' },
    { id: 'e-lesson-test', source: 'lesson', target: 'test', dashed: true },
    { id: 'e-test-run', source: 'test', target: 'run', animated: true, color: 'blue' },
    { id: 'e-run-pr', source: 'run', target: 'pr' },
    { id: 'e-pr-ci', source: 'pr', target: 'ci' },
    { id: 'e-ci-review', source: 'ci', target: 'review', label: 'Yes', color: 'green' },
    { id: 'e-ci-fix', source: 'ci', target: 'fix', label: 'No', color: 'red' },
    { id: 'e-fix-pr', source: 'fix', target: 'pr', dashed: true, color: 'red' },
    { id: 'e-review-merge', source: 'review', target: 'merge', animated: true, color: 'green' },
  ]

  return <FlowDiagram nodes={nodes} edges={edges} direction="LR" height={400} />
}
