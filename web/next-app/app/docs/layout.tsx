'use client'

import './kami.css'
import { useTheme } from '../../hooks/useTheme'
import Link from 'next/link'

export default function DocsLayout({ children }: { children: React.ReactNode }) {
  const { theme, mounted, toggleTheme } = useTheme()

  return (
    <div className={`kami-doc${theme === 'dark' ? ' dark' : ''}`}>
      {/* Navigation */}
      <nav style={{
        position: 'sticky',
        top: 0,
        zIndex: 100,
        background: 'var(--kami-bg)',
        borderBottom: '0.5px solid var(--kami-border)',
        padding: 'var(--kami-space-sm) var(--kami-space-lg)',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'space-between',
        fontFamily: 'var(--kami-mono)',
        fontSize: 'var(--kami-label-size)',
      }}>
        <div style={{ display: 'flex', gap: 'var(--kami-space-lg)', alignItems: 'center' }}>
          <Link
            href="/docs"
            style={{
              color: 'var(--kami-accent)',
              fontWeight: 500,
              textDecoration: 'none',
              fontFamily: 'var(--kami-serif)',
              fontSize: 'var(--kami-h3-size)',
            }}
          >
            Mathua Docs
          </Link>
          <Link href="/docs/architecture" style={{ color: 'var(--kami-text-muted)', textDecoration: 'none' }}>
            Architecture
          </Link>
          <Link href="/docs/contributing" style={{ color: 'var(--kami-text-muted)', textDecoration: 'none' }}>
            Contributing
          </Link>
        </div>
        <div style={{ display: 'flex', gap: 'var(--kami-space-md)', alignItems: 'center' }}>
          <Link href="/" style={{ color: 'var(--kami-text-muted)', textDecoration: 'none' }}>
            ← Home
          </Link>
          {mounted && (
            <button
              onClick={toggleTheme}
              aria-label="Toggle theme"
              style={{
                background: 'transparent',
                border: '0.5px solid var(--kami-border-strong)',
                color: 'var(--kami-text-muted)',
                padding: '2px var(--kami-space-sm)',
                cursor: 'pointer',
                fontFamily: 'var(--kami-mono)',
                fontSize: 'var(--kami-label-size)',
                borderRadius: 'var(--kami-radius-tight)',
              }}
            >
              {theme === 'dark' ? '\u2600' : '\u263E'}
            </button>
          )}
        </div>
      </nav>

      {/* Content */}
      {children}

      {/* Footer */}
      <footer style={{
        borderTop: '0.5px solid var(--kami-border)',
        padding: 'var(--kami-space-2xl) var(--kami-space-lg)',
        textAlign: 'center',
        fontFamily: 'var(--kami-mono)',
        fontSize: 'var(--kami-label-size)',
        color: 'var(--kami-text-muted)',
        lineHeight: 'var(--kami-lh-caption)',
      }}>
        <p style={{ margin: 0 }}>
          Mathua &middot; MIT License &middot;{' '}
          <a href="https://github.com/chuma-beep/mathua" style={{ color: 'var(--kami-accent)' }}>
            GitHub
          </a>
        </p>
      </footer>
    </div>
  )
}
