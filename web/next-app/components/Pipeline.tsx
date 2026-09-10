interface PipelineState {
  label: string
  status: 'unseen' | 'learning' | 'practicing' | 'mastered' | 'decaying'
}

interface PipelineProps {
  states: PipelineState[]
  className?: string
}

const statusStyle = {
  mastered: { color: 'var(--accent-blue)' },
  decaying: { color: 'var(--accent-teal)' },
} satisfies Record<string, React.CSSProperties>

export default function Pipeline({ states, className = '' }: PipelineProps) {
  return (
    <div className={`flex items-center justify-center gap-1 flex-wrap ${className}`}>
      {states.map((state, i) => (
        <span key={state.status} style={{ display: 'inline-flex', alignItems: 'center', gap: '0.25rem' }}>
          {i > 0 && (
            <span style={{ color: 'var(--border-strong)', fontFamily: "'IBM Plex Mono', monospace", fontSize: '13px' }}>
              {' → '}
            </span>
          )}
          <span
            style={{
              fontFamily: "'IBM Plex Mono', monospace",
              fontSize: '13px',
              letterSpacing: '0.05em',
              color: 'var(--text-secondary)',
              ...statusStyle[state.status],
            }}
          >
            {state.label}
          </span>
        </span>
      ))}
    </div>
  )
}
