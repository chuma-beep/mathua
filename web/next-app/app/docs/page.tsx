import Link from 'next/link'

const monoFont = "'JetBrains Mono', 'Fira Code', monospace"
const bodyFont = "'Lora', serif"

export default function DocsIndexPage() {
  const textStyle: React.CSSProperties = {
    fontFamily: bodyFont,
    fontSize: '0.95rem',
    color: 'var(--text-secondary)',
    lineHeight: 1.7,
  }

  return (
    <div className="max-w-[720px] mx-auto px-6 max-sm:px-4 py-20 max-sm:py-12">
      <h1 style={{
        fontFamily: "'Lora', serif",
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
          style={{
            display: 'flex',
            gap: '1rem',
            padding: '0.75rem 0',
            borderBottom: '0.5px solid var(--border)',
            textDecoration: 'none',
            color: 'inherit',
          }}
        >
          <span style={{
            fontFamily: monoFont,
            fontSize: '13px',
            color: 'var(--accent-gold)',
            minWidth: '1.5rem',
            paddingTop: '1px',
          }}>
            1.
          </span>
          <div>
            <div style={{
              fontFamily: "'Lora', serif",
              fontWeight: 400,
              fontSize: '1.1rem',
              color: 'var(--text-primary)',
              marginBottom: '2px',
            }}>
              Architecture
            </div>
            <p style={{ ...textStyle, margin: 0, fontSize: '0.9rem', color: 'var(--text-muted)' }}>
              How the engine, scheduler, generators, graders, and storage layers compose into a single Go binary.
            </p>
          </div>
        </Link>

        <Link
          href="/docs/contributing"
          style={{
            display: 'flex',
            gap: '1rem',
            padding: '0.75rem 0',
            borderBottom: '0.5px solid var(--border)',
            textDecoration: 'none',
            color: 'inherit',
          }}
        >
          <span style={{
            fontFamily: monoFont,
            fontSize: '13px',
            color: 'var(--accent-gold)',
            minWidth: '1.5rem',
            paddingTop: '1px',
          }}>
            2.
          </span>
          <div>
            <div style={{
              fontFamily: "'Lora', serif",
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
          fontFamily: "'Lora', serif",
          fontWeight: 400,
          fontSize: '1.3rem',
          color: 'var(--text-primary)',
          marginBottom: '1rem',
        }}>
          Quickstart
        </h2>

        <pre style={{
          fontFamily: monoFont,
          fontSize: '13px',
          color: 'var(--text-secondary)',
          whiteSpace: 'pre',
          overflowX: 'auto',
          lineHeight: 1.6,
          margin: 0,
        }}>
{`git clone https://github.com/chuma-beep/mathua.git
cd mathua
go build ./cmd/mathua

# Run the desktop TUI
./mathua

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
          Requirements: Go 1.21+
        </p>
      </div>
    </div>
  )
}
