'use client'

import { useTheme } from '../../hooks/useTheme'
import Link from 'next/link'

export default function DocsLayout({ children }: { children: React.ReactNode }) {
  const { theme, mounted, toggleTheme } = useTheme()

  return (
    <>
      {/* Nav Bar */}
      <nav style={{
        position: 'sticky',
        top: 0,
        zIndex: 100,
        background: theme === 'dark' ? 'rgba(11, 15, 26, 0.95)' : 'rgba(254, 252, 244, 0.95)',
        backdropFilter: 'blur(8px)',
        borderBottom: '0.5px solid var(--border)',
        padding: '12px 24px',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'space-between',
      }}>
        <div style={{ display: 'flex', gap: '24px', alignItems: 'center' }}>
          <Link
            href="/docs"
            style={{
              fontFamily: "'EB Garamond', Garamond, Georgia, serif",
              fontSize: '1.1rem',
              color: 'var(--accent-gold)',
              textDecoration: 'none',
              fontWeight: 400,
            }}
          >
            Mathua Docs
          </Link>
          <Link
            href="/docs/architecture"
            style={{
              fontFamily: "'JetBrains Mono', 'Fira Code', monospace",
              fontSize: '12px',
              color: 'var(--text-muted)',
              textDecoration: 'none',
            }}
          >
            Architecture
          </Link>
          <Link
            href="/docs/contributing"
            style={{
              fontFamily: "'JetBrains Mono', 'Fira Code', monospace",
              fontSize: '12px',
              color: 'var(--text-muted)',
              textDecoration: 'none',
            }}
          >
            Contributing
          </Link>
        </div>
        <div style={{ display: 'flex', gap: '12px', alignItems: 'center' }}>
          <Link
            href="/"
            style={{
              fontFamily: "'JetBrains Mono', 'Fira Code', monospace",
              fontSize: '12px',
              color: 'var(--text-muted)',
              textDecoration: 'none',
            }}
          >
            ← Home
          </Link>
          {mounted && (
            <button
              onClick={toggleTheme}
              aria-label="Toggle theme"
              style={{
                border: '0.5px solid var(--border-strong)',
                color: 'var(--text-muted)',
                padding: '4px 10px',
                cursor: 'pointer',
                fontFamily: "'JetBrains Mono', 'Fira Code', monospace",
                fontSize: '12px',
                background: 'var(--bg)',
                borderRadius: 0,
              }}
            >
              {theme === 'dark' ? '\u2600' : '\u263E'}
            </button>
          )}
        </div>
      </nav>

      {/* Content */}
      <main>{children}</main>

      {/* Footer */}
      <footer style={{
        borderTop: '0.5px solid var(--border)',
        padding: '32px 24px',
        textAlign: 'center',
        fontFamily: "'JetBrains Mono', 'Fira Code', monospace",
        fontSize: '11px',
        color: 'var(--text-muted)',
        lineHeight: 1.6,
      }}>
        Mathua ·{' '}
        <a href="#" style={{ color: 'var(--text-muted)', textDecoration: 'none' }}>Web App</a>
        {' · '}
        <a href="https://github.com/chuma-beep/mathua" style={{ color: 'var(--text-muted)', textDecoration: 'none' }}>GitHub</a>
        {' · '}
        MIT License
      </footer>
    </>
  )
}
