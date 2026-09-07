import type { Metadata } from 'next'
import Link from 'next/link'

export const metadata: Metadata = {
  title: 'Documentation \u2014 Mathua',
  description: 'Architecture and contributing guide for Mathua, the open-source adaptive math learning engine.',
}

const monoFont = "'IBM Plex Mono', monospace"
const bodyFont = "'IBM Plex Serif', serif"

export default function DocsIndexPage() {
  const textStyle: React.CSSProperties = {
    fontFamily: bodyFont,
    fontSize: '0.95rem',
    color: 'var(--text-secondary)',
    lineHeight: 1.7,
  }

  return (
    <>
    <div className="max-w-[720px] mx-auto px-4 sm:px-6 py-12 sm:py-20 pb-[calc(80px+env(safe-area-inset-bottom))] lg:pb-20 overflow-x-hidden min-w-0">
      <Link href="/" style={{ fontFamily: monoFont, fontSize: '13px', color: 'var(--text-muted)' }}>
        ← Back
      </Link>
      <h1 style={{
        fontFamily: "'IBM Plex Serif', serif",
        fontWeight: 400,
        fontSize: '1.9rem',
        color: 'var(--text-primary)',
        marginBottom: '0.5rem',
        borderBottom: '0.5px solid var(--border)',
        paddingBottom: '0.5rem',
      }}>
        Documentation
      </h1>

      <p style={{ ...textStyle, marginBottom: '2.5rem' }}>
        Architecture and contributing guide for Mathua.
      </p>

      {/* Doc list */}
      <div style={{ marginBottom: '3rem' }}>
        <Link
          href="/docs/architecture"
          className="link-underline"
          style={{
            display: 'flex',
            gap: '1rem',
            padding: '0.75rem 0',
            borderBottom: '0.5px solid var(--border)',
            color: 'inherit',
            textDecoration: 'none',
          }}
        >
          <span style={{
            fontFamily: monoFont,
            fontSize: '13px',
            color: 'var(--accent-blue)',
            minWidth: '1.5rem',
            paddingTop: '1px',
          }}>
            1.
          </span>
          <div>
            <div style={{
              fontFamily: "'IBM Plex Serif', serif",
              fontWeight: 400,
              fontSize: '1.1rem',
              color: 'var(--text-primary)',
              marginBottom: '2px',
            }}>
              Architecture
            </div>
            <p style={{ ...textStyle, margin: 0, fontSize: '0.9rem', color: 'var(--text-muted)' }}>
              Full system design: 5-layer architecture, data model, request flow, grading, CAT diagnostic, and trade-offs.
            </p>
          </div>
        </Link>

        <Link
          href="/docs/contributing"
          className="link-underline"
          style={{
            display: 'flex',
            gap: '1rem',
            padding: '0.75rem 0',
            borderBottom: '0.5px solid var(--border)',
            color: 'inherit',
            textDecoration: 'none',
          }}
        >
          <span style={{
            fontFamily: monoFont,
            fontSize: '13px',
            color: 'var(--accent-blue)',
            minWidth: '1.5rem',
            paddingTop: '1px',
          }}>
            2.
          </span>
          <div>
            <div style={{
              fontFamily: "'IBM Plex Serif', serif",
              fontWeight: 400,
              fontSize: '1.1rem',
              color: 'var(--text-primary)',
              marginBottom: '2px',
            }}>
              Contributing
            </div>
            <p style={{ ...textStyle, margin: 0, fontSize: '0.9rem', color: 'var(--text-muted)' }}>
              Add a concept, write a generator, pass the validator.
            </p>
          </div>
        </Link>
      </div>

      <div style={{
        borderTop: '0.5px solid var(--border)',
        paddingTop: '3rem',
      }}>
        <h2 style={{
          fontFamily: "'IBM Plex Serif', serif",
          fontWeight: 400,
          fontSize: '1.3rem',
          color: 'var(--text-primary)',
          marginBottom: '1rem',
        }}>
          Quickstart <span style={{ fontFamily: monoFont, fontSize: '12px', color: 'var(--text-muted)' }}>(for contributors)</span>
        </h2>

        <pre style={{
          fontFamily: monoFont,
          fontSize: '13px',
          color: 'var(--text-secondary)',
          whiteSpace: 'pre',
          overflowX: 'auto',
          lineHeight: 1.6,
          margin: 0,
          width: '100%',
          maxWidth: '100%',
          minWidth: 0,
        }}>
{`git clone https://github.com/chuma-beep/mathua.git
cd mathua
go build ./cmd/mathua

# Run the web server (requires Postgres)
DATABASE_URL=postgres://... ./mathua --serve --port 8080`}
        </pre>

        <p style={{
          fontFamily: bodyFont,
          fontStyle: 'italic',
          fontSize: '0.9rem',
          color: 'var(--text-muted)',
          marginTop: '1rem',
          lineHeight: 1.7,
        }}>
          Requirements: Go 1.21+ · New to Mathua? Start with <Link href="/how-it-works" className="link-underline" style={{ color: 'var(--accent-blue)' }}>How it works</Link>.
        </p>
      </div>
    </div>
    </>
  )
}
