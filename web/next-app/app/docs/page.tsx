'use client'

import Link from 'next/link'
import { useTheme } from '../../hooks/useTheme'
import SectionHeader from '../../components/SectionHeader'
import AsciiDivider from '../../components/AsciiDivider'

const docPages = [
  {
    href: '/docs/architecture',
    title: 'Architecture',
    description: 'From counting to calculus — how the engine, scheduler, generators, graders, and storage layers compose into a single Go binary.',
    tags: ['Engine', 'Design', 'System'],
  },
  {
    href: '/docs/contributing',
    title: 'Contributing',
    description: 'Add a concept, write a generator, pass the validator. Everything you need to contribute to the concept graph.',
    tags: ['Community', 'Guide', 'Concepts'],
  },
]

export default function DocsIndexPage() {
  const { theme } = useTheme()

  return (
    <div className="max-w-container mx-auto px-6 max-sm:px-4">
      <section className="py-20 max-sm:py-12">
        <div className="text-center">
          <SectionHeader
            label="Documentation"
            title="How Mathua is built"
          />
          <p
            style={{
              fontFamily: "'Source Serif 4', Georgia, serif",
              fontSize: '1rem',
              color: 'var(--text-secondary)',
              maxWidth: '600px',
              margin: '0 auto 2rem',
              lineHeight: 1.85,
            }}
          >
            Technical documentation for the adaptive math learning engine — architecture,
            contributing guide, and design decisions. Docs inherit the same warm palette
            and typography as the app.
          </p>
        </div>

        {/* Metric Row */}
        <div style={{
          display: 'grid',
          gridTemplateColumns: 'repeat(auto-fit, minmax(140px, 1fr))',
          gap: '16px',
          marginBottom: '80px',
          marginTop: '24px',
        }}>
          {[
            { num: '284', label: 'Concepts' },
            { num: '16', label: 'Domains' },
            { num: '2', label: 'Doc Pages' },
            { num: '1', label: 'Engine' },
          ].map(({ num, label }) => (
            <div
              key={label}
              style={{
                background: 'var(--surface)',
                border: '0.5px solid var(--border)',
                padding: '24px',
                textAlign: 'center',
              }}
            >
              <div style={{
                fontFamily: "'EB Garamond', Garamond, Georgia, serif",
                fontSize: '2rem',
                fontWeight: 400,
                color: 'var(--accent-gold)',
                lineHeight: 1,
                fontVariantNumeric: 'tabular-nums',
              }}>
                {num}
              </div>
              <div style={{
                fontFamily: "'JetBrains Mono', 'Fira Code', monospace",
                fontSize: '11px',
                color: 'var(--text-muted)',
                textTransform: 'uppercase' as const,
                letterSpacing: '0.05em',
                marginTop: '4px',
              }}>
                {label}
              </div>
            </div>
          ))}
        </div>

        {/* Doc cards */}
        <div style={{
          display: 'grid',
          gridTemplateColumns: 'repeat(auto-fit, minmax(320px, 1fr))',
          gap: '24px',
        }}>
          {docPages.map((doc) => (
            <Link
              key={doc.href}
              href={doc.href}
              style={{
                textDecoration: 'none',
                background: 'var(--surface)',
                border: '0.5px solid var(--border)',
                padding: '32px',
                display: 'block',
                transition: 'border-color 0.2s',
              }}
              onMouseEnter={(e) => { e.currentTarget.style.borderColor = 'var(--accent-gold)' }}
              onMouseLeave={(e) => { e.currentTarget.style.borderColor = 'var(--border)' }}
            >
              <h2 style={{
                fontFamily: "'EB Garamond', Garamond, Georgia, serif",
                fontWeight: 400,
                fontSize: '1.25rem',
                color: 'var(--text-primary)',
                margin: 0,
                marginBottom: '8px',
                lineHeight: 1.3,
              }}>
                {doc.title}
              </h2>
              <p style={{
                fontFamily: "'Source Serif 4', Georgia, serif",
                fontSize: '0.95rem',
                lineHeight: 1.7,
                color: 'var(--text-secondary)',
                marginBottom: '16px',
              }}>
                {doc.description}
              </p>
              <div style={{ display: 'flex', gap: '6px' }}>
                {doc.tags.map((tag) => (
                  <span
                    key={tag}
                    style={{
                      fontFamily: "'JetBrains Mono', 'Fira Code', monospace",
                      fontSize: '10px',
                      color: 'var(--text-muted)',
                      border: '0.5px solid var(--border)',
                      padding: '2px 8px',
                    }}
                  >
                    {tag}
                  </span>
                ))}
              </div>
            </Link>
          ))}
        </div>
      </section>

      <AsciiDivider pattern="dash" />

      <section className="py-20 max-sm:py-12">
        <SectionHeader
          label="Quickstart"
          title="Build from source"
        />
        <pre style={{
          background: 'var(--surface)',
          border: '0.5px solid var(--border)',
          padding: '24px',
          fontFamily: "'JetBrains Mono', 'Fira Code', monospace",
          fontSize: '13px',
          color: 'var(--text-secondary)',
          whiteSpace: 'pre',
          overflowX: 'auto',
          lineHeight: 1.6,
          maxWidth: '600px',
          margin: '0 auto',
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
          fontFamily: "'Source Serif 4', Georgia, serif",
          fontStyle: 'italic',
          fontSize: '0.9rem',
          color: 'var(--text-muted)',
          textAlign: 'center',
          marginTop: '16px',
          lineHeight: 1.7,
        }}>
          Requirements: Go 1.21+. For symbolic math grading (algebra concepts in future):
          Python 3.10+ with SymPy.
        </p>
      </section>
    </div>
  )
}
