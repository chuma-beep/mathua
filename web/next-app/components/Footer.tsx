import Link from 'next/link'

interface FooterProps {
  className?: string
}

const linkStyle: React.CSSProperties = {
  color: 'inherit',
  textDecoration: 'none',
  fontFamily: "'IBM Plex Mono', monospace",
  fontSize: '12px',
}

const linkItemStyle: React.CSSProperties = {
  ...linkStyle,
  color: 'var(--text-muted)',
}

const footerTitleStyle: React.CSSProperties = {
  fontFamily: "'IBM Plex Serif', serif",
  fontSize: '1rem',
  color: 'var(--text-secondary)',
  marginBottom: '0.75rem',
}

const linkContainerStyle: React.CSSProperties = {
  fontFamily: "'IBM Plex Mono', monospace",
  fontSize: '12px',
  color: 'var(--text-muted)',
  marginBottom: '1rem',
  display: 'flex',
  flexWrap: 'wrap',
  justifyContent: 'center',
  gap: '4px 12px',
}

const licenseStyle: React.CSSProperties = {
  fontFamily: "'IBM Plex Mono', monospace",
  fontStyle: 'italic',
  fontSize: '0.9rem',
  color: 'var(--text-muted)',
}

export default function Footer({ className = '' }: FooterProps) {
  return (
    <footer
      className={`text-center pt-8 pb-8 max-sm:pb-6 px-6 mt-10 ${className}`}
      style={{
        borderTop: '0.5px solid var(--border)',
        background: 'transparent',
      }}
    >
      <div style={footerTitleStyle}>
        Mathua: Math Understanding Agent
      </div>

      <div style={linkContainerStyle}>
        <span className="link-underline" style={linkItemStyle}>Web App</span>
        <a href="https://github.com/chuma-beep/mathua" className="link-underline" style={linkStyle}>GitHub</a>
        <Link href="/docs" className="link-underline" style={linkStyle}>Docs</Link>
        <Link href="/docs/contributing" className="link-underline" style={linkStyle}>Contributing</Link>
        <span className="link-underline" style={linkItemStyle}>Roadmap</span>
      </div>

      <div style={licenseStyle}>
        MIT License
      </div>
    </footer>
  )
}
