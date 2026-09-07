import Link from 'next/link'

interface FooterProps {
  className?: string
}

const linkStyle: React.CSSProperties = {
  textDecoration: 'none',
  fontFamily: "'IBM Plex Mono', monospace",
  fontSize: '12px',
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
      <div style={linkContainerStyle}>
        <a href="https://github.com/chuma-beep/mathua" className="link-underline" style={linkStyle}>GitHub</a>
        <Link href="/docs" className="link-underline" style={linkStyle}>Docs</Link>
        <Link href="/docs/contributing" className="link-underline" style={linkStyle}>Contributing</Link>
        <Link href="/note" className="link-underline" style={linkStyle}>Creator&apos;s Note</Link>
        <a href="https://www.dicebear.com" target="_blank" rel="noreferrer" className="link-underline" style={linkStyle}>Avatars by DiceBear</a>
      </div>

      <div style={licenseStyle}>
        MIT License
      </div>
    </footer>
  )
}
