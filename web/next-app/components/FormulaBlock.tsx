interface FormulaBlockProps {
  code: string
  label?: string
  className?: string
}

export default function FormulaBlock({ code, label, className = '' }: FormulaBlockProps) {
  return (
    <div className={`flex flex-col items-center ${className}`}>
      <pre
        style={{
          background: 'transparent',
          border: 'none',
          borderLeft: '2px solid var(--accent-gold)',
          borderRadius: 0,
          padding: '0.5rem 0 0.5rem 1.5rem',
          fontFamily: "'IBM Plex Mono', monospace",
          fontSize: '13px',
          color: 'var(--text-secondary)',
          whiteSpace: 'pre',
          overflowX: 'auto',
          lineHeight: 1.6,
          textAlign: 'left',
          display: 'block',
          maxWidth: '100%',
        }}
      >
        {code}
      </pre>
      {label && (
        <span style={{ fontFamily: "'IBM Plex Mono', monospace", fontSize: '11px', color: 'var(--text-muted)', marginTop: '0.25rem' }}>
          {label}
        </span>
      )}
    </div>
  )
}
