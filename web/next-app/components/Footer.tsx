interface FooterProps {
  className?: string
}

export default function Footer({ className = '' }: FooterProps) {
  return (
    <footer
      className={`text-center pt-8 pb-10 px-6 mt-10 ${className}`}
      style={{
        borderTop: '0.5px solid var(--border)',
        background: 'transparent',
      }}
    >
      <div
        style={{
          fontFamily: "'EB Garamond', Garamond, Georgia, serif",
          fontSize: '1rem',
          color: 'var(--text-secondary)',
          marginBottom: '0.75rem',
        }}
      >
        Mathua — Math Understanding Agent
      </div>

      <div
        style={{
          fontFamily: "'JetBrains Mono', 'Fira Code', monospace",
          fontSize: '11px',
          color: 'var(--text-muted)',
          marginBottom: '1rem',
        }}
      >
        <a href="#" style={{ color: 'inherit', textDecoration: 'none' }}>Web App</a>
        {'  ·  '}
        <a href="https://github.com/chuma-beep/mathua" style={{ color: 'inherit', textDecoration: 'none' }}>GitHub</a>
        {'  ·  '}
        <a href="https://github.com/chuma-beep/mathua/blob/main/docs/architecture.md" style={{ color: 'inherit', textDecoration: 'none' }}>Docs</a>
        {'  ·  '}
        <a href="https://github.com/chuma-beep/mathua/blob/main/CONTRIBUTING.md" style={{ color: 'inherit', textDecoration: 'none' }}>Contributing</a>
        {'  ·  '}
        <a href="#" style={{ color: 'inherit', textDecoration: 'none' }}>Roadmap</a>
      </div>

      <div
        style={{
          fontFamily: "'Source Serif 4', Georgia, serif",
          fontStyle: 'italic',
          fontSize: '0.9rem',
          color: 'var(--text-muted)',
          marginBottom: '0.5rem',
        }}
      >
        MIT License · Built in public by Wisdom Chukwuemeka
      </div>

      <div
        style={{
          fontFamily: "'Source Serif 4', Georgia, serif",
          fontStyle: 'italic',
          fontSize: '0.9rem',
          color: 'var(--text-muted)',
        }}
      >
        Inspired by the mastery-gating philosophy of Math Academy.
        <br />
        No content or code from Math Academy is used.
      </div>
    </footer>
  )
}
