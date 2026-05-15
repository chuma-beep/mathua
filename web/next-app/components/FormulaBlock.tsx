interface FormulaBlockProps {
  code: string
  label?: string
  className?: string
}

const formulaPreStyle: React.CSSProperties = {
  background: 'transparent',
  border: 'none',
  borderLeft: '2px solid var(--accent-blue)',
  borderRadius: 0,
  padding: '0.5rem 0 0.5rem 1.5rem',
  fontFamily: "'IBM Plex Mono', monospace",
  fontSize: '13px',
  color: 'var(--text-secondary)',
  whiteSpace: 'pre-wrap',
  overflowX: 'auto',
  lineHeight: 1.6,
  textAlign: 'left',
  display: 'block',
  maxWidth: '100%',
}

const formulaLabelStyle: React.CSSProperties = {
  fontFamily: "'IBM Plex Mono', monospace",
  fontSize: '12px',
  color: 'var(--text-muted)',
  marginTop: '0.25rem',
}

export default function FormulaBlock({ code, label, className = '' }: FormulaBlockProps) {
  return (
    <div className={`flex flex-col items-center ${className}`}>
      <pre style={formulaPreStyle}>
        {code}
      </pre>
      {label && (
        <span style={formulaLabelStyle}>
          {label}
        </span>
      )}
    </div>
  )
}
