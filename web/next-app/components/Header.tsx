'use client'

import Link from 'next/link'
import { useEffect, useState } from 'react'
import { useTheme } from '../hooks/useTheme'

interface HeaderLink {
  label: string
  href: string
}

interface HeaderProps {
  links?: HeaderLink[]
}

const monoFont = "'IBM Plex Mono', monospace"
const barStyle: React.CSSProperties = {
  position: 'sticky',
  top: 0,
  zIndex: 100,
  backdropFilter: 'blur(8px)',
  borderBottom: '0.5px solid var(--border)',
  padding: '12px 24px',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  transition: 'background 0.3s ease',
}

const linkStyle: React.CSSProperties = {
  fontFamily: monoFont,
  fontSize: '12px',
  color: 'var(--text-muted)',
  textDecoration: 'none',
  transition: 'color 0.2s',
}

const brandStyle: React.CSSProperties = {
  ...linkStyle,
  fontSize: '13px',
  color: 'var(--accent-gold)',
  marginRight: '1.5rem',
}

export default function Header({ links }: HeaderProps) {
  const { theme, mounted, toggleTheme } = useTheme()
  const [isMobile, setIsMobile] = useState(false)

  useEffect(() => {
    const check = () => setIsMobile(window.innerWidth < 640)
    check()
    window.addEventListener('resize', check)
    return () => window.removeEventListener('resize', check)
  }, [])

  return (
    <header
      className="max-sm:px-4"
      style={{
        ...barStyle,
        background: theme === 'dark' ? 'rgba(11, 15, 26, 0.95)' : 'rgba(254, 252, 244, 0.95)',
      }}
    >
      <div style={{
        display: 'flex',
        alignItems: 'center',
        gap: isMobile ? '0.5rem' : '1rem',
        overflowX: isMobile ? 'auto' : 'visible',
        flex: 1,
        marginRight: isMobile ? '0.5rem' : '0',
      }}>
        <Link href="/" style={{
          ...brandStyle,
          marginRight: isMobile ? '0.5rem' : '1.5rem',
          whiteSpace: 'nowrap',
        }}>Mathua</Link>
        {links?.map((link) => (
          <Link key={link.href} href={link.href} style={{ ...linkStyle, whiteSpace: 'nowrap' }}>
            {link.label}
          </Link>
        ))}
      </div>

      <div style={{ display: 'flex', alignItems: 'center', gap: '12px' }}>
        {mounted && (
          <button
            onClick={toggleTheme}
            aria-label="Toggle theme"
            style={{
              fontFamily: monoFont,
              fontSize: '12px',
              color: 'var(--text-muted)',
              background: 'transparent',
              border: '0.5px solid var(--border-strong)',
              padding: '4px 10px',
              cursor: 'pointer',
              borderRadius: 0,
              transition: 'color 0.2s, border-color 0.2s',
            }}
            onMouseEnter={(e) => {
              e.currentTarget.style.color = 'var(--accent-gold)'
              e.currentTarget.style.borderColor = 'var(--accent-gold)'
            }}
            onMouseLeave={(e) => {
              e.currentTarget.style.color = 'var(--text-muted)'
              e.currentTarget.style.borderColor = 'var(--border-strong)'
            }}
          >
            {theme === 'dark' ? '\u2600' : '\u263E'}
          </button>
        )}
      </div>
    </header>
  )
}
