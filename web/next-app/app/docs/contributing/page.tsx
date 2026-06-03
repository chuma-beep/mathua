import type { Metadata } from 'next'
import dynamic from 'next/dynamic'
import SectionHeader from '../../../components/SectionHeader'
import AsciiDivider from '../../../components/AsciiDivider'

export const metadata: Metadata = {
  title: 'Contributing \u2014 Mathua',
  description: 'How to contribute to Mathua — concepts, generators, frontend components, and more.',
}

const PrWorkflow = dynamic(() => import('../../../components/PrWorkflow'), {
  ssr: false,
  loading: () => (
    <div style={{
      height: 400, border: '0.5px solid var(--border)',
      display: 'flex', alignItems: 'center', justifyContent: 'center',
      color: 'var(--text-muted)', fontFamily: "'IBM Plex Mono', monospace", fontSize: '13px',
    }}>
      Loading workflow&hellip;
    </div>
  ),
})

const headingFont = "'IBM Plex Serif', serif"
const bodyFont = "'IBM Plex Serif', serif"
const monoFont = "'IBM Plex Mono', monospace"

const h2Style: React.CSSProperties = {
  fontFamily: headingFont,
  fontWeight: 400,
  fontSize: '1.3rem',
  color: 'var(--text-primary)',
  borderBottom: '0.5px solid var(--border)',
  paddingBottom: '0.5rem',
  marginBottom: '1rem',
}

const h3Style: React.CSSProperties = {
  fontFamily: headingFont,
  fontWeight: 400,
  fontSize: '1rem',
  color: 'var(--text-primary)',
  marginTop: '1.5rem',
  marginBottom: '0.5rem',
}

const bodyStyle: React.CSSProperties = {
  fontFamily: bodyFont,
  fontSize: '0.95rem',
  color: 'var(--text-secondary)',
  lineHeight: 1.85,
  marginBottom: '1rem',
}

const codeBlockStyle: React.CSSProperties = {
  background: 'var(--surface)',
  border: 'none',
  padding: '16px 20px',
  fontFamily: monoFont,
  fontSize: '13px',
  color: 'var(--text-secondary)',
  whiteSpace: 'pre',
  overflowX: 'auto',
  lineHeight: 1.6,
  margin: '1rem 0',
}

const inlineCodeStyle: React.CSSProperties = {
  fontFamily: monoFont,
  fontSize: '0.9em',
  color: 'var(--accent-blue)',
}

const mutedCodeStyle: React.CSSProperties = {
  fontFamily: monoFont,
  fontSize: '0.9em',
  color: 'var(--text-secondary)',
}

const tableHeaderStyle: React.CSSProperties = {
  fontFamily: monoFont,
  fontWeight: 400,
  fontSize: '0.85rem',
  textTransform: 'uppercase',
  letterSpacing: '0.08em',
  color: 'var(--accent-blue)',
  padding: '10px 14px 10px 0',
  borderBottom: '1px solid var(--accent-blue)',
  textAlign: 'left',
  background: 'transparent',
}

const tableCellStyle: React.CSSProperties = {
  fontFamily: bodyFont,
  fontSize: '0.9rem',
  color: 'var(--text-secondary)',
  padding: '10px 14px 10px 0',
  borderBottom: '0.5px solid var(--border)',
  background: 'transparent',
}

const calloutStyle: React.CSSProperties = {
  borderLeft: '2px solid var(--accent-blue)',
  paddingLeft: '1.5rem',
  fontFamily: bodyFont,
  fontStyle: 'italic',
  fontSize: '0.9rem',
  color: 'var(--text-muted)',
  lineHeight: 1.7,
}

export default function ContributingPage() {
  return (
    <div className="max-w-container mx-auto px-6 max-sm:px-4">
      <section className="pt-8">
        <SectionHeader label="Community" title="Contributing to Mathua" />
        <p style={{ ...bodyStyle, textAlign: 'center', maxWidth: '640px', margin: '0 auto 2rem' }}>
          Mathua is a full-stack project — a Go backend with a Next.js frontend,
          a scheduling engine, and a concept graph. Every generator, every diagram,
          every API route was built by someone who wanted to help others learn math.
          Here's how to join them.
        </p>
        <div style={{ textAlign: 'center', marginBottom: '1rem' }}>
          <a
            href="https://github.com/chuma-beep/mathua/blob/main/CONTRIBUTING.md"
            className="link-underline"
            style={{
              fontFamily: monoFont,
              fontSize: '12px',
              color: 'var(--text-muted)',
            }}
          >
            View source on GitHub
          </a>
        </div>
      </section>

      <AsciiDivider pattern="wave" />

      {/* Getting started */}
      <section className="py-20 max-sm:py-12">
        <h2 style={h2Style}>Getting started</h2>
        <p style={bodyStyle}>
          You'll need Go 1.25+ and Node.js 18+. Clone the repo:
        </p>
        <pre style={codeBlockStyle}>
{`git clone https://github.com/chuma-beep/mathua.git
cd mathua`}
        </pre>

        <h3 style={h3Style}>Running the full stack</h3>
        <p style={bodyStyle}>
          The backend serves the API on port 8080. The frontend (Next.js) runs
          on port 3000 and proxies API calls through the{' '}
          <code style={inlineCodeStyle}>NEXT_PUBLIC_API_URL</code> env var:
        </p>
        <pre style={codeBlockStyle}>
{`# Terminal 1 — Go backend
go build ./cmd/mathua
./mathua --serve --port 8080

# Terminal 2 — Next.js frontend
cd web/next-app
npm install
npm run dev`}
        </pre>
        <p style={bodyStyle}>
          For quick local testing without authentication, use{' '}
          <code style={inlineCodeStyle}>./mathua --serve --no-auth</code>.
          In that mode the session page will skip login and route you straight
          to practice.
        </p>

        <h3 style={h3Style}>Project layout</h3>
        <p style={bodyStyle}>
          Key directories you'll work in:
        </p>
        {[
          ['data/concepts/', 'Per-domain JSON files defining the concept graph.'],
          ['internal/generator/', 'Per-domain Go packages that produce problems.'],
          ['internal/server/', 'HTTP API server with route handlers.'],
          ['internal/storage/', 'Repository interface + SQLite implementation.'],
          ['internal/engine/', 'Scheduling engine, session logic, diagnostics.'],
          ['internal/auth/', 'JWT-based authentication service.'],
          ['web/next-app/app/', 'Next.js App Router pages and layouts.'],
          ['web/next-app/components/', 'React components (diagrams, heatmaps, etc).'],
          ['web/next-app/lib/', 'API client, auth helpers, shared utilities.'],
        ].map(([dir, desc]) => (
          <div key={dir} style={{ ...bodyStyle, display: 'flex', gap: '0.75rem', marginBottom: '0.3rem' }}>
            <code style={{ ...inlineCodeStyle, minWidth: '240px', fontSize: '0.8rem' }}>{dir}</code>
            <span style={{ color: 'var(--text-muted)' }}>{desc}</span>
          </div>
        ))}
      </section>

      <AsciiDivider pattern="dash" />

      {/* Step 1: Concepts */}
      <section className="py-20 max-sm:py-12">
        <h2 style={h2Style}>Step 1: Define a new concept</h2>
        <p style={bodyStyle}>
          Concepts live in per-domain JSON files under{' '}
          <code style={inlineCodeStyle}>data/concepts/</code>.
          Add your new concept as an entry in the array for the matching domain
          (e.g., add an arithmetic concept to{' '}
          <code style={inlineCodeStyle}>arithmetic.json</code>):
        </p>
        <pre style={codeBlockStyle}>
{`[
  { "id": "arith.add.single", "label": "Single-digit addition", ... },
  {
    "id":                "arith.mult.tables",
    "label":             "Multiplication tables 1–12",
    "domain":            "arithmetic",
    "subdomain":         "arithmetic.multiplication",
    "grading_type":      "numeric",
    "prerequisites":     ["arith.add.multi"],
    "mastery_threshold": {
      "streak":           7,
      "avg_time_seconds": 6.0
    }
  }
]`}
        </pre>

        <div style={{ overflowX: 'auto', margin: '1.5rem 0' }}>
          <table style={{
            width: '100%',
            borderCollapse: 'collapse',
            border: 'none',
            background: 'transparent',
          }}>
            <thead>
              <tr style={{ background: 'transparent' }}>
                <th style={tableHeaderStyle}>Field</th>
                <th style={tableHeaderStyle}>Type</th>
                <th style={tableHeaderStyle}>Required</th>
                <th style={tableHeaderStyle}>Description</th>
              </tr>
            </thead>
            <tbody>
              {[
                ['id', 'string', 'Yes', 'Unique dot-separated identifier (e.g., arith.add.single)'],
                ['label', 'string', 'Yes', 'Human-readable name displayed in the UI'],
                ['domain', 'string', 'Yes', 'One of the 16 domain categories'],
                ['subdomain', 'string', 'Optional', 'Nested grouping within a domain'],
                ['grading_type', 'enum', 'Yes', 'numeric for arithmetic, polynomial/expression for algebra (uses SymPy), multiple_choice, comparison, ordering'],
                ['prerequisites', 'string[]', 'Yes', 'Concept IDs that must be mastered first'],
                ['mastery_threshold', 'object', 'Yes', 'Object with streak (int) and avg_time_seconds (float) fields'],
                ['mastery_threshold.streak', 'number', 'Yes', 'Consecutive correct answers required for mastery'],
                ['mastery_threshold.avg_time_seconds', 'number', 'Yes', 'Maximum acceptable average response time in seconds'],
              ].map(([field, type, req, desc]) => (
                <tr key={field} style={{ background: 'transparent' }}>
                  <td style={{ ...tableCellStyle, fontFamily: monoFont, fontSize: '0.8rem' }}>{field}</td>
                  <td style={tableCellStyle}>{type}</td>
                  <td style={tableCellStyle}>
                    <span style={{
                      fontFamily: monoFont,
                      fontSize: '12px',
                      color: req === 'Yes' ? 'var(--accent-green)' : 'var(--text-muted)',
                      border: '0.5px solid',
                      borderColor: req === 'Yes' ? 'var(--accent-green)' : 'var(--border)',
                      padding: '1px 6px',
                    }}>
                      {req}
                    </span>
                  </td>
                  <td style={tableCellStyle}>{desc}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        <div style={{ ...calloutStyle, marginTop: '1rem' }}>
          The prerequisites list is the most important field. What must a student
          absolutely know before attempting this? If in doubt, add the prerequisite:
          the graph validator will catch cycles.
        </div>
      </section>

      <AsciiDivider pattern="dash" />

      {/* Step 2: Generators */}
      <section className="py-20 max-sm:py-12">
        <h2 style={h2Style}>Step 2: Teach Mathua to ask questions</h2>
        <p style={bodyStyle}>
          A generator is a Go function that produces a unique problem every time it's
          called. There is no static question bank: every problem is built on demand.
          Generators live in{' '}
          <code style={inlineCodeStyle}>internal/generator/[domain]/</code> and
          implement the Generator interface.
        </p>
        <pre style={codeBlockStyle}>
{`package arithmetic

import (
    "crypto/rand"
    "fmt"
    "math/big"
    "strconv"
    "github.com/chuma-beep/mathua/internal/generator"
)

type AddSingleGen struct{}

func (g *AddSingleGen) Generate(ctx generator.GeneratorContext) generator.Problem {
    max := int(5 + ctx.Difficulty*4)
    a, _ := rand.Int(rand.Reader, big.NewInt(int64(max)))
    b, _ := rand.Int(rand.Reader, big.NewInt(int64(max)))
    ai, bi := int(a.Int64())+2, int(b.Int64())+2
    ans := ai + bi
    return generator.Problem{
        Question:    fmt.Sprintf("%d + %d = ?", ai, bi),
        Answer:      strconv.Itoa(ans),
        Explanation: fmt.Sprintf("%d + %d = %d", ai, bi, ans),
    }
}`}
        </pre>
        <h3 style={h3Style}>Guidelines</h3>
        {[
          'Use the difficulty parameter to scale operand sizes. At 0.0, trivial. At 1.0, challenging for someone at that level.',
          'Always return an Explanation: it is shown when a student asks to see the solution.',
          'Use crypto/rand or math/rand with a seeded source. No hardcoded problems.',
          'Stay deterministic with respect to difficulty. A student should not get a meaningfully harder problem at the same difficulty.',
          'For non-standard answer formats (e.g., "5 R 3" in division), implement the GradedGenerator interface instead.',
        ].map((rule) => (
          <div key={rule} style={{ ...bodyStyle, marginBottom: '0.4rem' }}>
            <span style={{ color: 'var(--border-strong)', marginRight: '0.25rem', fontFamily: monoFont }}>·</span>
            {rule}
          </div>
        ))}
      </section>

      <AsciiDivider pattern="wave" />

      {/* Step 3: Fuzz tests */}
      <section className="py-20 max-sm:py-12">
        <h2 style={h2Style}>Step 3: Prove it works</h2>
        <p style={bodyStyle}>
          Every generator needs a fuzz test that asserts 1&thinsp;000 valid samples.
          This catches edge cases — division by zero, negative operand ranges,
          malformed output — before a student ever sees them.
        </p>
        <pre style={codeBlockStyle}>
{`func TestAddSingleGen(t *testing.T) {
    g := &AddSingleGen{}
    for i := 0; i < 1000; i++ {
        p := g.Generate(rand.Float64())
        assert.NotEmpty(t, p.Question)
        assert.NotEmpty(t, p.Answer)
        assert.NotEmpty(t, p.Explanation)
        assert.True(t, evaluateNumeric(p))
    }
}`}
        </pre>
        <p style={bodyStyle}>
          Run your tests before opening a PR:
        </p>
        <pre style={codeBlockStyle}>
          go test ./internal/generator/... -count 1000
        </pre>
      </section>

      <AsciiDivider pattern="wave" />

      {/* Step 4: Lessons */}
      <section className="py-20 max-sm:py-12">
        <h2 style={h2Style}>Step 4: Write a lesson (optional)</h2>
        <p style={bodyStyle}>
          Mathua includes a built-in lesson system. Each concept can have an
          associated lesson — a markdown file that teaches the material, rendered
          inside the app as a sidebar panel alongside practice problems.
        </p>
        <p style={bodyStyle}>
          Lessons live under{' '}
          <code style={inlineCodeStyle}>data/lessons/</code>{' '}
          and are registered in{' '}
          <code style={inlineCodeStyle}>data/lessons/lessons.json</code>,
          which maps concept IDs to file paths:
        </p>
        <pre style={codeBlockStyle}>
{`{ "concept_id": "calc.deriv.power_rule", "source": "advanced/polynomials/polynomials.md" }`}
        </pre>
        <p style={bodyStyle}>
          Multiple concepts can share a single lesson file. Content is written in
          Markdown with LaTeX via{' '}
          <code style={mutedCodeStyle}>{'\\\\( ... \\\\)'}</code>{' '}
          or{' '}
          <code style={mutedCodeStyle}>{'\\\\[ ... \\\\]'}</code>{' '}
          delimiters.
        </p>
      </section>

      <AsciiDivider pattern="dash" />

      {/* Frontend contributions */}
      <section className="py-20 max-sm:py-12">
        <h2 style={h2Style}>Frontend contributions</h2>
        <p style={bodyStyle}>
          The frontend is a Next.js App Router application under{' '}
          <code style={inlineCodeStyle}>web/next-app/</code>.
          It's designed with a brutalist aesthetic: monospace typography, minimal
          chrome, no rounded corners, and a terminal-calibrated color palette.
        </p>

        <h3 style={h3Style}>Design system</h3>
        {[
          'Typeface: IBM Plex Mono for code/UI, IBM Plex Serif for body/headings. Loaded via next/font as CSS variables (--font-ibm-plex-mono, --font-ibm-plex-serif).',
          'Colors: oklch color space as CSS custom properties (--bg, --surface, --text-primary, --accent-blue, --border, etc). Supports light/dark via class="dark" on <html>.',
          'Borders: 0.5px solid. No rounded corners (rounded-none everywhere). Zero box-shadow on interactive elements.',
          'Spacing: 1rem/4px base grid. Sections use py-20 (5rem). Cards use gap-3 (0.75rem) or gap-4 (1rem).',
          'Font sizes: body 0.95rem, headings scale from 1rem to 1.9rem, mono 13px for code, 12px for labels.',
        ].map((rule) => (
          <div key={rule} style={{ ...bodyStyle, marginBottom: '0.4rem' }}>
            <span style={{ color: 'var(--border-strong)', marginRight: '0.25rem', fontFamily: monoFont }}>·</span>
            {rule}
          </div>
        ))}

        <h3 style={h3Style}>React Flow diagrams</h3>
        <p style={bodyStyle}>
          Architecture, scheduler, and student model diagrams use{' '}
          <code style={inlineCodeStyle}>@xyflow/react</code> (React Flow v12).
          All diagrams share{' '}
          <code style={inlineCodeStyle}>FlowDiagram</code> as a common wrapper
          in{' '}
          <code style={inlineCodeStyle}>components/FlowDiagram.tsx</code>.
        </p>
        {[
          'Use useNodesState + useEdgesState + onNodesChange + onEdgesChange for every diagram. Without these handlers, node dragging is silently ignored.',
          'When adding custom node types, define them outside the component body (otherwise nodes remount on every render).',
          'For diagrams that need pan/zoom, set allowZoom=true and do not pass restrictive interaction props — the Controls lock button toggles the React Flow store directly.',
          'Background dots use var(--border-strong) at size=1 for visibility without distraction.',
          'Node colors come from the shared themeColors object in FlowDiagram.tsx — use those instead of hardcoding.',
        ].map((rule) => (
          <div key={rule} style={{ ...bodyStyle, marginBottom: '0.4rem' }}>
            <span style={{ color: 'var(--border-strong)', marginRight: '0.25rem', fontFamily: monoFont }}>·</span>
            {rule}
          </div>
        ))}

        <h3 style={h3Style}>Components &amp; pages</h3>
        <p style={bodyStyle}>
          Components live in{' '}
          <code style={inlineCodeStyle}>components/</code>. Pages live in{' '}
          <code style={inlineCodeStyle}>app/</code> under their route segment
          (e.g., <code style={inlineCodeStyle}>app/profile/page.tsx</code>).
          The project uses static export (<code style={inlineCodeStyle}>output: 'export'</code>)
          — all pages are client-rendered with <code style={inlineCodeStyle}>'use client'</code>.
        </p>
        {[
          'API calls go through lib/api.ts. Auth headers come from lib/auth.ts (JWT in localStorage).',
          'Use the motion-safe animation classes for transitions (fadeIn, ascii-reveal). Avoid heavy animation libraries.',
          'Profile pages, heatmaps, and stats components consume data from GET /api/activity, /api/scores, and /api/weaknesses.',
          'No third-party UI libraries except React Flow and sonner (toasts). Build everything else from scratch.',
        ].map((rule) => (
          <div key={rule} style={{ ...bodyStyle, marginBottom: '0.4rem' }}>
            <span style={{ color: 'var(--border-strong)', marginRight: '0.25rem', fontFamily: monoFont }}>·</span>
            {rule}
          </div>
        ))}
      </section>

      <AsciiDivider pattern="wave" />

      {/* API server contributions */}
      <section className="py-20 max-sm:py-12">
        <h2 style={h2Style}>API server contributions</h2>
        <p style={bodyStyle}>
          The HTTP server is in{' '}
          <code style={inlineCodeStyle}>internal/server/server.go</code>.
          It's a standard library <code style={inlineCodeStyle}>net/http</code> server
          with a custom <code style={inlineCodeStyle}>ServeMux</code>.
          All routes are registered in{' '}
          <code style={inlineCodeStyle}>Register()</code>.
        </p>

        <h3 style={h3Style}>Adding a new endpoint</h3>
        {[
          'Define your handler method on *Server (e.g., handleSomething).',
          'Register it in Register() with mux.HandleFunc("/api/your-path", logRequest(cors(s.authMiddleware(s.handleSomething)))).',
          'If the endpoint needs auth, studentID is available from r.Context().Value(authStudentKey{}).',
          'If the endpoint should work without auth, check s.auth == nil and handle both paths.',
          'All endpoints set Content-Type: application/json via writeJSON(). Use decodeJSON() and writeError() for request/error handling.',
          'POST endpoints check r.Method != http.MethodPost and return 405.',
          'CORS is handled automatically by the cors middleware wrapper — no extra config needed.',
        ].map((rule) => (
          <div key={rule} style={{ ...bodyStyle, marginBottom: '0.4rem' }}>
            <span style={{ color: 'var(--border-strong)', marginRight: '0.25rem', fontFamily: monoFont }}>·</span>
            {rule}
          </div>
        ))}

        <h3 style={h3Style}>Storage layer</h3>
        <p style={bodyStyle}>
          Data access goes through the{' '}
          <code style={inlineCodeStyle}>storage.Repository</code> interface in{' '}
          <code style={inlineCodeStyle}>internal/storage/store.go</code>.
          The SQLite implementation is in{' '}
          <code style={inlineCodeStyle}>internal/storage/sqlite.go</code>.
          When adding a new query:
        </p>
        {[
          'Add the method signature to the Repository interface first.',
          'Implement it in sqlite.go using the database/sql package.',
          'The server accesses the repo via s.repo.',
          'CGO is required (mattn/go-sqlite3). The Dockerfile sets CGO_ENABLED=1.',
        ].map((rule) => (
          <div key={rule} style={{ ...bodyStyle, marginBottom: '0.4rem' }}>
            <span style={{ color: 'var(--border-strong)', marginRight: '0.25rem', fontFamily: monoFont }}>·</span>
            {rule}
          </div>
        ))}
      </section>

      <AsciiDivider pattern="wave" />

      {/* Deployment */}
      <section className="py-20 max-sm:py-12">
        <h2 style={h2Style}>Deployment</h2>
        <p style={bodyStyle}>
          Mathua is deployed as two services:
        </p>
        {[
          'Backend: Go binary on Fly.io (Dockerfile, fly.toml). Requires a persistent volume for SQLite (mathua_data at /data) and a JWT_SECRET environment variable.',
          'Frontend: Static export on Vercel via GitHub auto-deploy. The NEXT_PUBLIC_API_URL env var points to the Fly.io backend. Set in .env.production.',
          'Local development: Backend on :8080, frontend on :3000. NEXT_PUBLIC_API_URL in .env.local should be http://localhost:8080.',
          'The Dockerfile is multi-stage: golang:1.26-bookworm builder → debian:bookworm-slim runtime. CGO_ENABLED=1 is required for SQLite.',
        ].map((rule) => (
          <div key={rule} style={{ ...bodyStyle, marginBottom: '0.4rem' }}>
            <span style={{ color: 'var(--border-strong)', marginRight: '0.25rem', fontFamily: monoFont }}>·</span>
            {rule}
          </div>
        ))}
      </section>

      <AsciiDivider pattern="dash" />

      {/* Validator */}
      <section className="py-20 max-sm:py-12">
        <h2 style={h2Style}>The graph validator</h2>
        <p style={bodyStyle}>
          A validator runs on every pull request. It checks two invariants
          before any merge can happen:
        </p>
        {[
          'No cycles: concept A cannot require B while B requires A.',
          'No orphans: every prerequisite must exist in the graph.',
        ].map((rule, i) => (
          <div key={rule} style={{ ...bodyStyle, marginBottom: '0.4rem' }}>
            <span style={{ color: 'var(--accent-blue)', fontFamily: monoFont, fontSize: '13px' }}>
              {i + 1}.
            </span>{' '}
            {rule}
          </div>
        ))}
        <pre style={codeBlockStyle}>
          go run scripts/validate_graph.go
        </pre>
      </section>

      <AsciiDivider pattern="wave" />

      {/* PR workflow */}
      <section className="py-20 max-sm:py-12">
        <h2 style={h2Style}>Submitting a pull request</h2>
        <p style={bodyStyle}>
          Once everything passes locally, here's the full checklist:
        </p>
        {[
          'Add the concept to the appropriate domain file in data/concepts/ with correct prerequisites.',
          'Write the Go generator in the appropriate domain subdirectory under internal/generator/.',
          'Write the fuzz test with 1 000 samples.',
          '(Optional) Write a lesson and register it in data/lessons/lessons.json.',
          'Run go test ./... and go run scripts/validate_graph.go locally.',
          'For frontend changes: run npm run build in web/next-app/ to verify the static export compiles.',
          'For API changes: run go vet ./internal/... to check for issues.',
          'Open a PR. The CI pipeline runs the validator and all tests automatically.',
          'A maintainer reviews the concept ordering, thresholds, generator quality, and UI changes.',
        ].map((step, i) => (
          <div key={step} style={{ ...bodyStyle, marginBottom: '0.4rem' }}>
            <span style={{ color: 'var(--accent-blue)', fontFamily: monoFont, fontSize: '13px' }}>
              {i + 1}.
            </span>{' '}
            {step}
          </div>
        ))}
        <PrWorkflow />
        <p style={{ ...bodyStyle, marginTop: '1rem' }}>
          Reviews usually happen within a few days. If a week passes with no
          response, feel free to ping the thread. We read every PR.
        </p>
      </section>

      <AsciiDivider pattern="dash" />

      {/* Design conventions */}
      <section className="py-20 max-sm:py-12">
        <h2 style={h2Style}>Design conventions</h2>

        <h3 style={h3Style}>Backend</h3>
        {[
          'Concept IDs follow domain.subdomain.descriptor: lower case, no spaces.',
          'Mastery thresholds are pragmatic. Single-digit addition should require faster response (6–8 s) than multi-digit multiplication (15–20 s).',
          'Subdomains group related concepts. If a domain grows past 15 concepts, consider introducing subdomains.',
          'Generators accept a generator.GeneratorContext with Difficulty (0.0–1.0) and Seed (int64). Use ctx.Seed for deterministic generation — this enables question replay, regression suites, and A/B testing.',
          'Difficulty scaling should be linear where sensible. The jump from 0.0 to 1.0 should feel meaningful, not extreme.',
          'Grading types: use numeric for arithmetic, polynomial/expression for algebra (routes through SymPy), and comparison/ordering/multiple_choice for structured answers.',
        ].map((rule) => (
          <div key={rule} style={{ ...bodyStyle, marginBottom: '0.4rem' }}>
            <span style={{ color: 'var(--border-strong)', marginRight: '0.25rem', fontFamily: monoFont }}>·</span>
            {rule}
          </div>
        ))}

        <h3 style={h3Style}>Frontend</h3>
        {[
          'Colors only through CSS variables (var(--accent-blue), var(--border), etc). Never hardcode hex/rgb.',
          'Typography: var(--font-ibm-plex-mono) for code, var(--font-ibm-plex-serif) for prose. Loaded via next/font CSS variables.',
          'Borders: border-[0.5px] with var(--border) / var(--border-strong). rounded-none on all interactive elements.',
          'No box-shadow on buttons or inputs. Use border changes or background transitions for hover states.',
          'Animations: use the @keyframes defined in tailwind.config.js (ascii-reveal, fadeIn, progress-fill). Keep transitions under 300ms.',
          'Diagrams: follow the React Flow conventions above. Use dagre for automatic layout, themeColors for node styling.',
          'Read the full design system in DESIGN.md for color palette, typography scales, spacing, and component patterns.',
        ].map((rule) => (
          <div key={rule} style={{ ...bodyStyle, marginBottom: '0.4rem' }}>
            <span style={{ color: 'var(--border-strong)', marginRight: '0.25rem', fontFamily: monoFont }}>·</span>
            {rule}
          </div>
        ))}

        <div style={{ ...calloutStyle, marginTop: '1.5rem' }}>
          Read the full design system in{' '}
          <a href="https://github.com/chuma-beep/mathua/blob/main/DESIGN.md" style={{ color: 'var(--accent-blue)', textDecoration: 'underline' }}>
            DESIGN.md
          </a>{' '}
          for color palette, typography, spacing, and component patterns.
        </div>
      </section>
    </div>
  )
}
