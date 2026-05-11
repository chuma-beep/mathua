interface PipelineState {
  label: string
  status: 'unseen' | 'learning' | 'practicing' | 'mastered' | 'decaying'
}

interface PipelineProps {
  states: PipelineState[]
  className?: string
}

const statusStyle: Record<string, React.CSSProperties> = {
  mastered: { color: 'var(--accent-gold)' },
  decaying: { color: 'var(--accent-teal)' },
}

export default function Pipeline({ states, className = '' }: PipelineProps) {
  return (
    <div className={`flex items-center justify-center gap-1 flex-wrap ${className}`}>
      {states.map((state, i) => (
        <span key={i} style={{ display: 'inline-flex', alignItems: 'center', gap: '0.25rem' }}>
          {i > 0 && (
            <span style={{ color: 'var(--border-strong)', fontFamily: "'JetBrains Mono', monospace", fontSize: '13px' }}>
              {' → '}
            </span>
          )}
          <span
            style={{
              fontFamily: "'JetBrains Mono', 'Fira Code', monospace",
              fontSize: '13px',
              letterSpacing: '0.08em',
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
