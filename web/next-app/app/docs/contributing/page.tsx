'use client'

import { useTheme } from '../../../hooks/useTheme'
import SectionHeader from '../../../components/SectionHeader'
import AsciiDivider from '../../../components/AsciiDivider'
import MermaidDiagram from '../../../components/MermaidDiagram'

const headingFont = "'Cormorant Garamond', serif"
const bodyFont = "'Inter', -apple-system, sans-serif"
const monoFont = "'JetBrains Mono', 'Fira Code', monospace"

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
  padding: '20px 24px',
  fontFamily: monoFont,
  fontSize: '13px',
  color: 'var(--text-secondary)',
  whiteSpace: 'pre',
  overflowX: 'auto',
  lineHeight: 1.6,
  margin: '1rem 0',
}

const tableHeaderStyle: React.CSSProperties = {
  fontFamily: monoFont,
  fontWeight: 400,
  fontSize: '0.85rem',
  textTransform: 'uppercase',
  letterSpacing: '0.08em',
  color: 'var(--accent-gold)',
  padding: '10px 14px 10px 0',
  borderBottom: '1px solid var(--accent-gold)',
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

export default function ContributingPage() {
  const { theme } = useTheme()

  return (
    <div className="max-w-container mx-auto px-6 max-sm:px-4">
      <section className="pt-8">
        <SectionHeader label="Community" title="Contributing to Mathua" />
        <p style={{ ...bodyStyle, textAlign: 'center', maxWidth: '640px', margin: '0 auto 2rem' }}>
          The most impactful contributions are new concepts and improved generators. Adding
          a concept requires exactly three things: a JSON entry in the concept graph, a Go
          generator function, and a fuzz test.
        </p>
      </section>

      {/* Contributor flow diagram */}
      <MermaidDiagram
        theme={theme}
        chart={`flowchart TD
    Fork["Fork the repo"]
    Concept["Add concept to\\nconcepts.json"]
    Gen["Write generator\\nGo function"]
    Fuzz["Write fuzz test\\n1000 samples"]
    Valid["Run validator\\nvalidate_graph.go"]
    PR["Open pull request"]
    CI{"CI passes?"}
    Review["Maintainer review"]
    Merge["Merged!"]
    Fix["Fix issues"]

    Fork --> Concept
    Concept --> Gen
    Gen --> Fuzz
    Fuzz --> Valid
    Valid --> PR
    PR --> CI
    CI -->|"Yes"| Review
    CI -->|"No"| Fix
    Fix --> PR
    Review --> Merge`}
      />

      {/* Metric cards */}
      <div style={{
        display: 'grid',
        gridTemplateColumns: 'repeat(auto-fit, minmax(140px, 1fr))',
        gap: '16px',
        marginBottom: '48px',
        marginTop: '24px',
      }}>
        {[
          { num: '3', label: 'Steps' },
          { num: '1000', label: 'Fuzz Samples' },
          { num: '0', label: 'Cycles Tolerated' },
        ].map(({ num, label }) => (
          <div
            key={label}
            style={{
              background: 'var(--surface)',
              border: '0.5px solid var(--border)',
              padding: '20px',
              textAlign: 'center',
            }}
          >
            <div style={{ fontFamily: monoFont, fontSize: '2rem', fontWeight: 400, color: 'var(--accent-gold)', lineHeight: 1 }}>
              {num}
            </div>
            <div style={{ fontFamily: monoFont, fontSize: '11px', color: 'var(--text-muted)', textTransform: 'uppercase', letterSpacing: '0.05em', marginTop: '4px' }}>
              {label}
            </div>
          </div>
        ))}
      </div>

      <AsciiDivider pattern="wave" />

      {/* Step 1 */}
      <section className="py-20 max-sm:py-12">
        <h2 style={h2Style}>I. Add a node to the concept graph</h2>
        <p style={bodyStyle}>
          Every concept in Mathua lives in{' '}
          <code style={{ fontFamily: monoFont, fontSize: '0.9em', color: 'var(--accent-gold)' }}>data/concepts.json</code>.
          Adding a concept means adding a JSON object with an ID, label, domain, prerequisite
          list, and mastery thresholds.
        </p>
        <pre style={codeBlockStyle}>
{`{
  "id":                "arith.mult.tables",
  "label":             "Multiplication tables 1-12",
  "domain":            "arithmetic",
  "subdomain":         "arithmetic.multiplication",
  "grading_type":      "numeric",
  "prerequisites":     ["arith.add.multi"],
  "mastery_threshold": {
    "streak":           7,
    "avg_time_seconds": 6.0
  }
}`}
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
                ['grading_type', 'enum', 'Yes', 'numeric for arithmetic, polynomial for algebra'],
                ['prerequisites', 'string[]', 'Yes', 'List of concept IDs that must be mastered first'],
                ['streak', 'number', 'Yes', 'Consecutive correct answers required'],
                ['avg_time_seconds', 'number', 'Yes', 'Maximum acceptable average response time'],
              ].map(([field, type, req, desc], i) => (
                <tr key={i} style={{ background: 'transparent' }}>
                  <td style={{ ...tableCellStyle, fontFamily: monoFont, fontSize: '0.8rem' }}>{field}</td>
                  <td style={tableCellStyle}>{type}</td>
                  <td style={tableCellStyle}>
                    <span style={{
                      fontFamily: monoFont,
                      fontSize: '10px',
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

        <div style={{
          borderLeft: '2px solid var(--accent-gold)',
          paddingLeft: '1.5rem',
          fontFamily: bodyFont,
          fontStyle: 'italic',
          fontSize: '0.9rem',
          color: 'var(--text-muted)',
          lineHeight: 1.7,
          marginTop: '1rem',
        }}>
          The prerequisite list is the most important field. Think carefully — what must a student
          absolutely know before attempting this concept? If in doubt, add the prerequisite.
          The graph validator will catch cycles.
        </div>
      </section>

      <AsciiDivider pattern="dash" />

      {/* Step 2 */}
      <section className="py-20 max-sm:py-12">
        <h2 style={h2Style}>II. Write the generator function</h2>
        <p style={bodyStyle}>
          Each concept needs a generator — a Go function that produces a unique problem every
          time it is called. Generators live in{' '}
          <code style={{ fontFamily: monoFont, fontSize: '0.9em', color: 'var(--accent-gold)' }}>internal/generator/[domain]/</code> and
          implement the Generator interface.
        </p>
        <pre style={codeBlockStyle}>
{`package arithmetic

import "mathua/internal/generator"

type AddSingleGen struct{}

func (g *AddSingleGen) Generate(difficulty float64) generator.Problem {
    max := int(5 + difficulty*4)
    a   := rand.Intn(max) + 2
    b   := rand.Intn(max) + 2
    ans := a + b
    return generator.Problem{
        Question:    fmt.Sprintf("%d + %d = ?", a, b),
        Answer:      strconv.Itoa(ans),
        Explanation: fmt.Sprintf("%d + %d = %d", a, b, ans),
    }
}`}
        </pre>
        <h3 style={h3Style}>Guidelines for good generators</h3>
        {[
          'Use the difficulty parameter to scale operand sizes. At 0.0 the problem should be trivial; at 1.0 it should be challenging for a student at that level.',
          'Always provide an Explanation field. It is shown when a student asks to see the solution.',
          'Use crypto/rand or math/rand with a seeded source. Do not hardcode problems — the generator must produce unique output each call.',
          'Keep the generator deterministic with respect to difficulty. A student should not get a meaningfully harder problem at the same difficulty.',
        ].map((rule, i) => (
          <div key={i} style={{ ...bodyStyle, marginBottom: '0.4rem' }}>
            <span style={{ color: 'var(--border-strong)', marginRight: '0.25rem', fontFamily: monoFont }}>·</span>
            {rule}
          </div>
        ))}
      </section>

      <AsciiDivider pattern="wave" />

      {/* Step 3 */}
      <section className="py-20 max-sm:py-12">
        <h2 style={h2Style}>III. Write a fuzz test</h2>
        <p style={bodyStyle}>
          Every generator must have a fuzz test that asserts 1000 valid samples. This catches
          edge cases — division by zero, negative operand ranges, malformed output — before
          they reach a student.
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
        <div style={{
          borderLeft: '2px solid var(--accent-gold)',
          paddingLeft: '1.5rem',
          fontFamily: bodyFont,
          fontStyle: 'italic',
          fontSize: '0.9rem',
          color: 'var(--text-muted)',
          lineHeight: 1.7,
        }}>
          Run generator fuzz tests with{' '}
          <code style={{ fontFamily: monoFont, fontSize: '0.9em', color: 'var(--accent-gold)' }}>go test ./internal/generator/... -run TestFuzz -count 1000</code>.
          All tests must pass before the PR is merged.
        </div>
      </section>

      <AsciiDivider pattern="dash" />

      {/* Validator */}
      <section className="py-20 max-sm:py-12">
        <h2 style={h2Style}>Graph validator — automatic on every PR</h2>
        <p style={bodyStyle}>
          A graph validator runs on every pull request. It checks two invariants:
        </p>
        {[
          'No cycles. The DAG must remain acyclic. A cycle means concept A requires B, which requires A — creating an impossible prerequisite chain.',
          'No orphans. Every prerequisite referenced must exist in the graph. A reference to a non-existent concept ID is rejected.',
        ].map((rule, i) => (
          <div key={i} style={{ ...bodyStyle, marginBottom: '0.4rem' }}>
            <span style={{ color: 'var(--accent-gold)', fontFamily: monoFont, fontSize: '13px' }}>
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

      {/* PR Checklist */}
      <section className="py-20 max-sm:py-12">
        <h2 style={h2Style}>Pull request checklist</h2>
        {[
          'Add the concept to data/concepts.json with correct prerequisites.',
          'Write the generator in the appropriate domain subdirectory.',
          'Write the fuzz test with 1000 samples.',
          'Run go test ./... and go run scripts/validate_graph.go locally.',
          'Open a PR. The CI pipeline runs the validator and all tests.',
          'A maintainer reviews the concept ordering, threshold values, and generator quality.',
        ].map((step, i) => (
          <div key={i} style={{ ...bodyStyle, marginBottom: '0.4rem' }}>
            <span style={{ color: 'var(--accent-gold)', fontFamily: monoFont, fontSize: '13px' }}>
              {i + 1}.
            </span>{' '}
            {step}
          </div>
        ))}
      </section>

      <AsciiDivider pattern="dash" />

      {/* Design conventions */}
      <section className="py-20 max-sm:py-12">
        <h2 style={h2Style}>Design conventions</h2>
        {[
          'Concept IDs follow the pattern domain.subdomain.descriptor. Use snake_case, no spaces.',
          'Mastery thresholds are pragmatic. A single-digit addition should require faster response (6–8 seconds) than multi-digit multiplication (15–20 seconds).',
          'Subdomains group related concepts. If a domain grows beyond 15 concepts, consider introducing subdomains.',
          'Difficulty scaling should be linear where possible. The difference between difficulty 0.0 and 1.0 should feel meaningful but not extreme.',
        ].map((rule, i) => (
          <div key={i} style={{ ...bodyStyle, marginBottom: '0.4rem' }}>
            <span style={{ color: 'var(--border-strong)', marginRight: '0.25rem', fontFamily: monoFont }}>·</span>
            {rule}
          </div>
        ))}
        <div style={{
          borderLeft: '2px solid var(--accent-gold)',
          paddingLeft: '1.5rem',
          fontFamily: bodyFont,
          fontStyle: 'italic',
          fontSize: '0.9rem',
          color: 'var(--text-muted)',
          lineHeight: 1.7,
          marginTop: '1.5rem',
        }}>
          Read the full design system in{' '}
          <a href="https://github.com/chuma-beep/mathua/blob/main/DESIGN.md" style={{ color: 'var(--accent-gold)', textDecoration: 'underline' }}>
            DESIGN.md
          </a>{' '}
          for color palette, typography, spacing, and component patterns used across the project.
        </div>
      </section>
    </div>
  )
}
