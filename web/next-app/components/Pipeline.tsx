interface PipelineState {
  label: string
  status: 'unseen' | 'learning' | 'practicing' | 'mastered' | 'decaying'
}

interface PipelineProps {
  states: PipelineState[]
  className?: string
}

export default function Pipeline({ states, className = '' }: PipelineProps) {
  return (
    <div className={`flex items-center justify-center gap-2 flex-wrap ${className}`}>
      {states.map((state, i) => (
        <span key={i}>
          {i > 0 && <span className="pipeline-arrow mx-1">→</span>}
          <span
            className={`pipeline-state ${
              state.status === 'mastered'
                ? 'mastered'
                : state.status === 'decaying'
                ? 'decaying'
                : ''
            }`}
          >
            {state.label}
          </span>
        </span>
      ))}
    </div>
  )
}
