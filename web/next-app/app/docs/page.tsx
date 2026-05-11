import Link from 'next/link'

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
  return (
    <div className="kami-container">
      <section style={{ paddingBottom: 'var(--kami-space-3xl)' }}>
        <span className="kami-overline">Documentation</span>
        <h1 style={{
          fontFamily: 'var(--kami-serif)',
          fontWeight: 500,
          fontSize: 'var(--kami-display-size)',
          lineHeight: 'var(--kami-lh-display)',
          color: 'var(--kami-text)',
          marginBottom: 'var(--kami-space-md)',
        }}>
          Good content deserves
          <br />
          good paper.
        </h1>
        <p style={{
          fontFamily: 'var(--kami-serif)',
          fontSize: 'var(--kami-h3-size)',
          lineHeight: 'var(--kami-lh-h3)',
          color: 'var(--kami-text-secondary)',
          maxWidth: '36rem',
          marginBottom: 'var(--kami-space-lg)',
        }}>
          Technical documentation for Mathua — the adaptive math learning engine.
          Every page follows the Kami design system: warm parchment, ink blue accent,
          serif-led hierarchy.
        </p>
      </section>

      {/* Metric Row */}
      <section style={{
        display: 'grid',
        gridTemplateColumns: 'repeat(auto-fit, minmax(140px, 1fr))',
        gap: 'var(--kami-space-lg)',
        marginBottom: 'var(--kami-space-3xl)',
      }}>
        <div className="kami-card">
          <div className="kami-metric">
            <div className="kami-metric-number">284</div>
            <div className="kami-metric-label">Concepts</div>
          </div>
        </div>
        <div className="kami-card">
          <div className="kami-metric">
            <div className="kami-metric-number">16</div>
            <div className="kami-metric-label">Domains</div>
          </div>
        </div>
        <div className="kami-card">
          <div className="kami-metric">
            <div className="kami-metric-number">2</div>
            <div className="kami-metric-label">Doc Pages</div>
          </div>
        </div>
        <div className="kami-card">
          <div className="kami-metric">
            <div className="kami-metric-number">1</div>
            <div className="kami-metric-label">Engine</div>
          </div>
        </div>
      </section>

      {/* Doc cards */}
      <section>
        <span className="kami-overline">Explore</span>
        <div style={{
          display: 'grid',
          gridTemplateColumns: 'repeat(auto-fit, minmax(320px, 1fr))',
          gap: 'var(--kami-space-lg)',
          marginTop: 'var(--kami-space-lg)',
        }}>
          {docPages.map((doc) => (
            <Link
              key={doc.href}
              href={doc.href}
              className="kami-feature"
              style={{ textDecoration: 'none', display: 'block', transition: 'box-shadow 0.2s' }}
            >
              <h2 style={{
                fontFamily: 'var(--kami-serif)',
                fontWeight: 500,
                fontSize: 'var(--kami-h2-size)',
                lineHeight: 'var(--kami-lh-h2)',
                color: 'var(--kami-text)',
                margin: 0,
                marginBottom: 'var(--kami-space-sm)',
                borderBottom: 'none',
                paddingBottom: 0,
              }}>
                {doc.title}
              </h2>
              <p style={{
                fontFamily: 'var(--kami-serif)',
                fontSize: 'var(--kami-caption-size)',
                lineHeight: 'var(--kami-lh-caption)',
                color: 'var(--kami-text-secondary)',
                marginBottom: 'var(--kami-space-md)',
              }}>
                {doc.description}
              </p>
              <div style={{ display: 'flex', gap: 'var(--kami-space-xs)' }}>
                {doc.tags.map((tag) => (
                  <span key={tag} className="kami-tag kami-tag--weak">{tag}</span>
                ))}
              </div>
            </Link>
          ))}
        </div>
      </section>

      <hr />

      {/* Quickstart */}
      <section>
        <span className="kami-overline">Quickstart</span>
        <h2>Build from source</h2>
        <pre>{`git clone https://github.com/chuma-beep/mathua.git
cd mathua
go build ./cmd/mathua

# Run the desktop TUI
./mathua

# Run the web server (requires Postgres)
DATABASE_URL=postgres://... ./mathua --serve --port 8080`}</pre>
        <blockquote>
          Requirements: Go 1.21+. For symbolic math grading (algebra concepts): Python 3.10+ with SymPy.
        </blockquote>
      </section>
    </div>
  )
}
